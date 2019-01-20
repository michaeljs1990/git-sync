package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	gconfig "gopkg.in/src-d/go-git.v4/config"
	gittrans "gopkg.in/src-d/go-git.v4/plumbing/transport"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	github "github.com/google/go-github/v21/github"
	log "github.com/sirupsen/logrus"
	oauth2 "golang.org/x/oauth2"
	yaml "gopkg.in/yaml.v2"
)

const (
	FetchMirror      = "fetch_mirror"
	GithubPushMirror = "github_push_mirror"
)

var (
	NoGithubConfigFound = errors.New("No github config was found for this repo")
)

// Config is the parent level of our YAML data that
// all other config should tie back into.
type Config struct {
	Interval int
	Path     string
	Repos    []RepoConfig
	Github   []GithubConfig
}

// RepoConfig is for adhoc servers that may not live
// in a place with a common api such as the linux kernel
type RepoConfig struct {
	URL      string
	Type     string
	Path     string
	Remote   string
	Refs     []string
	Metadata []string
	HTTPAuth HTTPAuth
}

// GithubConfig holds a single github user and will pull down all
// repos locally. This should likely allow for some caching since
// github has some strict rate limiting and they shouldn't change
// often in reality.
type GithubConfig struct {
	Username string
	RepoType string
	Repos    bool
	Protocol string
	Metadata []string
	HTTPAuth HTTPAuth
}

type HTTPAuth struct {
	User  string
	Token string
}

// repo maps very closely to the RepoConfig structure in almost
// all cases however repo is what all other configs must conform
// to for being processed
type repo struct {
	URL      string
	Type     string
	Path     string
	Remote   string
	Refs     []string
	Metadata []string
	HTTPAuth HTTPAuth
}

func (r RepoConfig) toRepo(c Config) repo {
	return repo{
		URL:      r.URL,
		Refs:     r.Refs,
		Type:     r.Type,
		Path:     r.Path,
		Remote:   r.Remote,
		Metadata: r.Metadata,
		HTTPAuth: r.HTTPAuth,
	}
}

func (r GithubConfig) toRepos(c Config) []repo {
	oauth, _ := makeGithubClient(r.HTTPAuth)
	client := github.NewClient(oauth)

	// If the user has not set the repo type provide a sane default
	// and try to get all public repos available options are all,
	// public, private, forks, sources
	repoType := r.RepoType
	if r.RepoType == "" {
		repoType = "public"
	}

	ret := []repo{}
	opt := &github.RepositoryListOptions{Type: repoType}
	repos, _, _ := client.Repositories.List(context.Background(), r.Username, opt)

	for _, v := range repos {
		// Pick the type of protocol that we are going to use to clone over if none is
		// provided we just clone over HTTP. Github also supports SVN clone but adding
		// support for SVN/Mg is a fairly large change.
		var url string
		switch r.Protocol {
		case "git":
			url = *v.GitURL
		case "http":
			url = *v.HTMLURL
		case "ssh":
			url = *v.SSHURL
		case "":
			url = *v.HTMLURL
		default:
			log.Error("The protocol you provided for github is not supported")
			os.Exit(1)
		}

		// Need to add support here for using different kinds of urls but currently
		// only HTTP/HTTPS clones are. In addition to this we need to do sanity checking
		// such as not passing SSH auth info to HTTP methods and so on.
		ret = append(ret, repo{
			URL:      url,
			Type:     FetchMirror,
			Path:     path.Join(c.Path, *v.FullName),
			Remote:   "origin",
			Metadata: r.Metadata,
			HTTPAuth: r.HTTPAuth,
		})
	}

	return ret
}

// determine the type of auth needed to clone down this repo
// if multiple are supplied it just picks the first one. If we
// can't find anything to use as an auth method we return an
// empty one.
func (r repo) getAuth() gittrans.AuthMethod {
	switch {
	case r.HTTPAuth != (HTTPAuth{}):
		return &githttp.BasicAuth{
			Username: r.HTTPAuth.User,
			Password: r.HTTPAuth.Token,
		}
	default:
		var emptyAuth gittrans.AuthMethod
		return emptyAuth
	}
}

// Check if the current repo is able to create a github client and
// return one if it's able to do so.
func makeGithubClient(ga HTTPAuth) (*http.Client, error) {
	if ga.Token == "" {
		return nil, NoGithubConfigFound
	}

	ctx := context.Background()
	token := &oauth2.Token{AccessToken: ga.Token}
	ts := oauth2.StaticTokenSource(token)

	return oauth2.NewClient(ctx, ts), nil
}

func bareMirrorClone(r repo) error {
	// Ensure that the path leading up to the repo is created properly mostly
	// for ensuring that repos with multiple slashes and github style name/repo
	// get put into sub folders. If stat returns an error we assume it means the
	// path doesn't exist yet and we try to make it.
	dir, _ := path.Split(r.Path)
	stat, err := os.Stat(dir)
	_, pathError := err.(*os.PathError)
	if pathError || stat.IsDir() == false {
		log.WithField("repo", r.URL).Info("Creating sub paths before cloning repo")
		if os.MkdirAll(dir, os.FileMode(0755)) != nil {
			log.WithField("repo", r.URL).Error("Was unable to create the file path needed for clone : ", err.Error())
		}
	}

	// Setup io.Writter logger with some extra fields to make this info
	// useful if we encounter any issues around clones hanging.
	w := log.StandardLogger().WithField("repo", r.URL).WriterLevel(log.TraceLevel)
	defer w.Close()

	log.WithField("repo", r.URL).Info("Trying to bare clone repo")
	_, err = git.PlainClone(r.Path, true, &git.CloneOptions{
		URL:      r.URL,
		Auth:     r.getAuth(),
		Progress: w,
	})

	if err != nil {
		switch err {
		case git.ErrRepositoryAlreadyExists:
			log.WithField("repo", r.URL).Info("This is expected and good : ", err.Error())
			return err
		default:
			log.WithField("repo", r.URL).Error(err.Error())
			return err
		}
	}

	log.WithField("repo", r.URL).Info("Finished bare clone of repo")
	return err
}

// Try as best as we can to sync the repository with it's remote.
// wc in this context stands for working_copy but go hates underscores.
func syncRepo(r repo) error {
	wc, err := git.PlainOpen(r.Path)
	if err != nil {
		log.WithField("repo", r.URL).Error("Was trying to create a repo for syncing at path ", r.Path, " : ", err.Error())
		return err
	}

	log.WithField("repo", r.URL).Trace("Trying to sync repo")

	wl := log.StandardLogger().WithField("repo", r.URL).WriterLevel(log.TraceLevel)
	defer wl.Close()

	// Setup refspec if you customized it or not. You can easily break this by putting
	// in an invalid refspec but i'm really not sure how to guard against it.
	refspec := []gconfig.RefSpec{}
	if len(r.Refs) != 0 {
		for _, ref := range r.Refs {
			refspec = append(refspec, gconfig.RefSpec(ref))
		}
	} else {
		refspec = append(refspec, gconfig.RefSpec("refs/heads/*:refs/remotes/heads/*"))
		refspec = append(refspec, gconfig.RefSpec("refs/tags/*:refs/remotes/tags/*"))
	}

	// This is a little silly since if it's up to date it still returns
	// that information via an error message.
	err = wc.Fetch(&git.FetchOptions{
		RemoteName: r.Remote,
		Progress:   wl,
		Force:      true,
		RefSpecs:   refspec,
		Auth:       r.getAuth(),
	})

	if err != nil {
		switch err {
		case git.NoErrAlreadyUpToDate:
			log.WithField("repo", r.URL).Info(err.Error())
			return nil
		default:
			log.WithField("repo", r.URL).Error("Was trying to fetch : ", err.Error())
			return err
		}
	}

	log.WithField("repo", r.URL).Info("Repository has been updated to the latest state")

	return err
}

// Sync a local repo to github. This won't try to do any fancy stuff so if
// the repos become out of sync it will just fail.
func syncToGithub(r repo) error {
	wc, err := git.PlainOpen(r.Path)
	if err != nil {
		log.WithField("repo", r.URL).Error("Was trying to create a repo for syncing at path ", r.Path, " : ", err.Error())
		return err
	}

	wl := log.StandardLogger().WithField("repo", r.URL).WriterLevel(log.TraceLevel)
	defer wl.Close()

	// Setup refspec if you customized it or not. You can easily break this by putting
	// in an invalid refspec but i'm really not sure how to guard against it.
	refspec := []gconfig.RefSpec{}
	if len(r.Refs) != 0 {
		for _, ref := range r.Refs {
			refspec = append(refspec, gconfig.RefSpec(ref))
		}
	} else {
		refspec = append(refspec, gconfig.RefSpec("refs/heads/*:refs/heads/*"))
		refspec = append(refspec, gconfig.RefSpec("refs/tags/*:refs/tags/*"))
	}

	// Ensure the remote exists that we want to be syncing to
	_, err = wc.CreateRemote(&gconfig.RemoteConfig{
		Name: r.Remote,
		URLs: []string{r.URL},
	})

	// This needs a username and password or better yet a username and auth token.
	// In the case of an auth token the username can actually be anything that you want.
	err = wc.Push(&git.PushOptions{
		RemoteName: r.Remote,
		Auth:       r.getAuth(),
		RefSpecs:   refspec,
		Progress:   wl,
	})

	if err != nil {
		switch err {
		case git.NoErrAlreadyUpToDate:
			log.WithField("repo", r.URL).Info(err.Error())
			return nil
		default:
			log.WithField("repo", r.URL).Error("Was trying to Push : ", err.Error())
			return err
		}
	}

	return nil
}

func repoWorker(wg *sync.WaitGroup, rc <-chan repo) {
	defer wg.Done()

	for r := range rc {
		switch r.Type {
		case FetchMirror:
			log.WithField("repo", r.URL).Info("Fetch mirror job launched")
			err := bareMirrorClone(r)
			if err == git.ErrRepositoryAlreadyExists {
				syncRepo(r)
			}
		case GithubPushMirror:
			log.WithField("repo", r.URL).Info("Github push mirror job launched")
			syncToGithub(r)
		default:
			log.WithField("repo", r.URL).Info(r.Type, " is not a supported job type")
		}
	}
}

// Feed the channel all the information that it needs
func feedChannel(jobs chan repo, c Config, oneshot bool) {
	for {
		for _, v := range c.Repos {
			jobs <- v.toRepo(c)
		}

		for _, v := range c.Github {
			for _, r := range v.toRepos(c) {
				jobs <- r
			}
		}

		if oneshot == true {
			close(jobs)
			return
		}

		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

func launchWorkers(workers int, wg *sync.WaitGroup, rc <-chan repo) {
	for i := 1; i <= workers; i++ {
		log.Trace("Worker ", i, " has been started")
		wg.Add(1)
		go repoWorker(wg, rc)
	}
}

func main() {
	conf := flag.String("config", "config.yaml", "Config file to start the application with")
	workers := flag.Int("workers", 1, "The number of workers trying to update repos")
	oneshot := flag.Bool("oneshot", false, "Only run the script once and then exit upon completion")

	verbose := flag.Bool("verbose", false, "Control the level of logging you would like output")
	log_file := flag.String("log_file", "", "Set the file to log to by default this just sends data to stdout")
	log_format := flag.String("log_format", "", "Set the the log format that you would like to be output")
	flag.Parse()

	// Logging setup
	if *log_format == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}

	if *log_file == "" {
		log.SetOutput(os.Stdout)
	} else {
		fh, err := os.Create(*log_file)
		if err != nil {
			fmt.Printf("Error creating log file (%s)", err.Error())
			os.Exit(1)
		}
		log.SetOutput(fh)
	}

	if *verbose {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	// Load our config into memory and do some error checking
	dat, err := ioutil.ReadFile(*conf)
	if err != nil {
		log.Fatal(err)
	}

	// This should be validated properly
	config := Config{}
	err = yaml.Unmarshal(dat, &config)
	if err != nil {
		log.Fatal(err)
	}

	// Setup everything needed to run this as a daemon or optionally
	// as a one time service so you can use it in a cron if desired
	var wg sync.WaitGroup
	queue := make(chan repo, *workers*2)
	go feedChannel(queue, config, *oneshot)
	launchWorkers(*workers, &wg, queue)

	// Hang for now but later this should do some checking for
	// signals that may be sent to the processes as well as managing
	// when the queue should have jobs pushed onto it again
	wg.Wait()
}
