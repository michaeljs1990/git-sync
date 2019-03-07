package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	// Sorry the naming here got kinda bad of note gconfig is for configuration
	// of what branches and refs you want to push. gitconfig is for playing with
	// the actual git config file.
	git "gopkg.in/src-d/go-git.v4"
	gconfig "gopkg.in/src-d/go-git.v4/config"
	gitconfig "gopkg.in/src-d/go-git.v4/plumbing/format/config"
	gittrans "gopkg.in/src-d/go-git.v4/plumbing/transport"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	github "github.com/google/go-github/v21/github"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
	jsonschema "github.com/santhosh-tekuri/jsonschema"
	log "github.com/sirupsen/logrus"
	ssh "golang.org/x/crypto/ssh"
	oauth2 "golang.org/x/oauth2"
	yaml "gopkg.in/yaml.v2"
)

const (
	FetchMirror = "fetch_mirror"
	PushMirror  = "push_mirror"

	GitProtocol  = "git"
	SshProtocol  = "ssh"
	HttpProtocol = "http"

	CgitMetadata = "cgit"
)

var (
	NoGithubConfigFound  = errors.New("No github config was found for this repo")
	HttpProtocolRegex, _ = regexp.Compile("^http(s)?://")
	GitProtocolRegex, _  = regexp.Compile("^git://")
	SshProtocolRegex, _  = regexp.Compile("^ssh://|.*:.*")
)

var (
	globalConf = Config{}
	configLock = sync.RWMutex{}
)

// Config is the parent level of our YAML data that
// all other config should tie back into.
type Config struct {
	Interval int            `json:"interval"`
	Path     string         `json:"path"`
	Repos    []RepoConfig   `json:"repos"`
	Github   []GithubConfig `json:"github"`
}

// Extras struct is loaded with pretty much anything you want
// that might be useful for people performing metadata operations
type Extras struct {
	Username    string `json:"username"`
	CgitSection string `json:"cgitsection"`
	CgitOwner   string `json:"cgitowner"`
	Description string `json:"description"`
}

// RepoConfig is for adhoc servers that may not live
// in a place with a common api such as the linux kernel
type RepoConfig struct {
	URL        string     `json:"url"`
	Type       string     `json:"type"`
	Path       string     `json:"path"`
	Remote     string     `json:"remote"`
	Refs       []string   `json:"refs"`
	Metadata   []string   `json:"metadata"`
	Extras     Extras     `json:"extras"`
	HTTPAuth   HTTPAuth   `json:"httpauth"`
	SSHAuth    SSHAuth    `json:"sshauth"`
	SSHKeyAuth SSHKeyAuth `json:"sshkeyauth"`
}

// GithubConfig holds a single github user and will pull down all
// repos locally. This should likely allow for some caching since
// github has some strict rate limiting and they shouldn't change
// often in reality.
type GithubConfig struct {
	Username   string     `json:"username"`
	RepoType   string     `json:"repotype"`
	Repos      bool       `json:"repos"`
	Protocol   string     `json:"protocol"`
	Metadata   []string   `json:"metadata"`
	Extras     Extras     `json:"extras"`
	HTTPAuth   HTTPAuth   `json:"httpauth"`
	SSHAuth    SSHAuth    `json:"sshauth"`
	SSHKeyAuth SSHKeyAuth `json:"sshkeyauth"`
}

type HTTPAuth struct {
	User  string `json:"user"`
	Token string `json:"token"`
}

type SSHAuth struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type SSHKeyAuth struct {
	User    string `json:"user"`
	KeyPath string `json:"keypath"`
}

// repo maps very closely to the RepoConfig structure in almost
// all cases however repo is what all other configs must conform
// to for being processed
type repo struct {
	URL        string
	Type       string
	Path       string
	Remote     string
	Refs       []string
	Metadata   []string
	Extras     Extras
	HTTPAuth   HTTPAuth
	SSHAuth    SSHAuth
	SSHKeyAuth SSHKeyAuth
}

func (r RepoConfig) toRepo(c Config) repo {
	return repo{
		URL:        r.URL,
		Refs:       r.Refs,
		Type:       r.Type,
		Path:       r.Path,
		Remote:     r.Remote,
		Metadata:   r.Metadata,
		Extras:     r.Extras,
		HTTPAuth:   r.HTTPAuth,
		SSHAuth:    r.SSHAuth,
		SSHKeyAuth: r.SSHKeyAuth,
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
	opt := &github.RepositoryListOptions{
		Type: repoType,
		ListOptions: github.ListOptions{
			PerPage: 10000,
		},
	}
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

		// Setup our extras struct
		extras := Extras{
			Username:    r.Username,
			CgitSection: r.Username,
			CgitOwner:   r.Username,
		}

		// Need to add support here for using different kinds of urls but currently
		// only HTTP/HTTPS clones are. In addition to this we need to do sanity checking
		// such as not passing SSH auth info to HTTP methods and so on.
		ret = append(ret, repo{
			URL:        url,
			Type:       FetchMirror,
			Path:       path.Join(c.Path, *v.FullName),
			Remote:     "origin",
			Metadata:   r.Metadata,
			Extras:     extras,
			HTTPAuth:   r.HTTPAuth,
			SSHAuth:    r.SSHAuth,
			SSHKeyAuth: r.SSHKeyAuth,
		})
	}

	return ret
}

// Determine the git url protocol type from the URL
// that has been set for it.
func (r repo) protocolType() string {
	url := []byte(r.URL)
	switch {
	case HttpProtocolRegex.Match(url):
		return HttpProtocol
	case GitProtocolRegex.Match(url):
		return GitProtocol
	case SshProtocolRegex.Match(url):
		return SshProtocol
	default:
		log.WithField("repo", r.URL).Fatal("The URL for this repo doesn't look to be a valid git URL")
		return ""
	}
}

// determine the type of auth needed to clone down this repo
// if multiple are supplied it just picks the first one. If we
// can't find anything to use as an auth method we return an
// empty one.
func (r repo) getAuth() gittrans.AuthMethod {

	switch {
	case r.HTTPAuth != (HTTPAuth{}) && r.protocolType() == HttpProtocol:
		return &githttp.BasicAuth{
			Username: r.HTTPAuth.User,
			Password: r.HTTPAuth.Token,
		}
	case r.SSHKeyAuth != (SSHKeyAuth{}) && r.protocolType() == SshProtocol:
		sshKey, err := ioutil.ReadFile(r.SSHKeyAuth.KeyPath)
		if err != nil {
			log.WithField("repo", r.URL).Fatal("Unable to read SSH Key : ", err.Error())
		}
		signer, err := ssh.ParsePrivateKey([]byte(sshKey))
		if err != nil {
			log.WithField("repo", r.URL).Fatal("Unable to parse private key : ", err.Error())
		}
		return &gitssh.PublicKeys{
			User:   r.SSHKeyAuth.User,
			Signer: signer,
		}
	case r.SSHAuth != (SSHAuth{}) && r.protocolType() == SshProtocol:
		return &gitssh.Password{
			User:     r.SSHAuth.User,
			Password: r.SSHAuth.Password,
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
		URL:        r.URL,
		RemoteName: r.Remote,
		Auth:       r.getAuth(),
		Progress:   w,
	})

	if err != nil {
		switch err {
		case git.ErrRepositoryAlreadyExists:
			log.WithField("repo", r.URL).Info("This is expected and good : ", err.Error())
			return err
		default:
			log.WithField("repo", r.URL).Error("This is very bad : ", err.Error())
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
		log.WithField("repo", r.URL).Error("Was trying to sync repo at path ", r.Path, " : ", err.Error())
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
		refspec = append(refspec, gconfig.RefSpec("refs/heads/*:refs/heads/*"))
		refspec = append(refspec, gconfig.RefSpec("refs/tags/*:refs/tags/*"))
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

// Sync a local repo to some other source. This won't try to do any fancy stuff so if
// the repos become out of sync it will just fail.
func syncToRepo(r repo) error {
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

func cgitMetadataOperation(r repo) error {
	gitConfigFile, err := os.OpenFile(path.Join(r.Path, "config"), os.O_RDWR, 0644)
	if err != nil {
		return errors.New("Unable to read this repos git config file : " + err.Error())
	}

	decoder := gitconfig.NewDecoder(gitConfigFile)
	config := &gitconfig.Config{}
	decoder.Decode(config)

	// Seek to the start of the file so we overwrite it instead of appending
	_, err = gitConfigFile.Seek(0, 0)
	if err != nil {
		return errors.New("Unable to seek to the start of the file before writing it out : " + err.Error())
	}

	// Try and be clever about what metadata we should be writing out
	if r.Extras != (Extras{}) {
		config.SetOption("cgit", "", "owner", r.Extras.CgitOwner)
		config.SetOption("cgit", "", "section", r.Extras.CgitSection)
	}

	encoder := gitconfig.NewEncoder(gitConfigFile)
	err = encoder.Encode(config)
	if err != nil {
		return errors.New("Unable to write this repos git config file : " + err.Error())
	}

	err = gitConfigFile.Close()
	if err != nil {
		return errors.New("Unable to close the file we wrote to something went wrong: " + err.Error())
	}

	return nil
}

// This checks if the user has defined any metadata operations to be performed on the repo.
// An example of that is doing some custom configuration to the git config file that may
// be specific to the type of repository browser that you are using.
func runMetadataOperations(r repo) {
	for _, m := range r.Metadata {
		switch m {
		case CgitMetadata:
			err := cgitMetadataOperation(r)
			if err != nil {
				log.WithField("repo", r.URL).Error("Issue writing cgit metadata : ", err.Error())
			} else {
				log.WithField("repo", r.URL).Info("Cgit metadata has been configured for this repo")
			}
		default:
			log.WithField("repo", r.URL).Error("The meta data option (", m, ") has not been defined")
		}
	}
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
		case PushMirror:
			log.WithField("repo", r.URL).Info("Push mirror job launched")
			syncToRepo(r)
		default:
			log.WithField("repo", r.URL).Info(r.Type, " is not a supported job type")
		}

		// These should log information but never cause a crash
		runMetadataOperations(r)
	}
}

// Read config makes a copy of the current config so that if the structure
// is modified by a user signal we don't have race conditions in the goroutines
// accessing the returned variable.
func readConfig() Config {
	var n bytes.Buffer
	enc := gob.NewEncoder(&n)
	dec := gob.NewDecoder(&n)

	dupe := Config{}

	configLock.RLock()
	defer configLock.RUnlock()

	enc.Encode(globalConf)
	dec.Decode(&dupe)
	return dupe
}

// Feed the channel all the information that it needs
func feedChannel(jobs chan repo, oneshot bool) {
	for {
		makeJobs(jobs)

		if oneshot == true {
			close(jobs)
			return
		}

		c := readConfig()

		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

func makeJobs(jobs chan repo) {
	// Get the latest config that might have been updated
	c := readConfig()

	for _, v := range c.Repos {
		jobs <- v.toRepo(c)
	}

	for _, v := range c.Github {
		for _, r := range v.toRepos(c) {
			jobs <- r
		}
	}
}

func launchWorkers(workers int, wg *sync.WaitGroup, rc <-chan repo) {
	for i := 1; i <= workers; i++ {
		log.Trace("Worker ", i, " has been started")
		wg.Add(1)
		go repoWorker(wg, rc)
	}
}

// I'm not in love with what is going on here but in practice it should work just fine.
// I think the one edge case is that this will make it so we can't validate a required
// field that could be empty. I can't think of any reason I would need that though. So
// what we are doing, is taking YAML, something that is half sane for human config and
// populating the Config struct. We then turn it into JSON so we can pass it to the JSON
// schema validator. Modern problems require modern solutions.
func loadConfig(config string) (Config, error) {
	c := Config{}

	// Load our config into memory and do some error checking
	dat, err := ioutil.ReadFile(config)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(dat, &c)
	if err != nil {
		return c, err
	}

	b, err := json.Marshal(c)
	if err != nil {
		return c, err
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(schema)); err != nil {
		return c, err
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return c, err
	}

	if err := schema.Validate(strings.NewReader(string(b))); err != nil {
		return c, err
	}

	return c, nil
}

func handleSignals(signals <-chan os.Signal, config string, jobs chan repo) {
	for signal := range signals {
		switch signal {
		case syscall.SIGUSR1:
			log.Warn("Config reload requested")
			newConfig, err := loadConfig(config)
			if err != nil {
				log.Error("Config not reloaded :", err.Error())
			} else {
				configLock.Lock()
				globalConf = newConfig
				configLock.Unlock()
				log.Warn("Config reload finished")
			}
		case syscall.SIGUSR2:
			log.Warn("Manual sync job was fired off")
			makeJobs(jobs)
		}
	}
}

func main() {
	conf := flag.String("config", "config.yaml", "Config file to start the application with")
	workers := flag.Int("workers", 1, "The number of workers trying to update repos")
	oneshot := flag.Bool("oneshot", false, "Only run the script once and then exit upon completion")
	validate := flag.Bool("validate", false, "Validate the config file that the user passes to us and then stop")
	prom_addr := flag.String("prom-listen-address", ":8080", "The address to listen on for HTTP requests.")

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

	config, err := loadConfig(*conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Set the global conf no mutex is needed at this time
	globalConf = config

	// If we have gotten this far and not failed we think it's all good
	log.Info("Configuration has been validated")
	if *validate {
		fmt.Println("Your config looks good to us!")
		os.Exit(0)
	}

	// Export prometheus metrics about the service
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(*prom_addr, nil)

	// Setup everything needed to run this as a daemon or optionally
	// as a one time service so you can use it in a cron if desired
	var wg sync.WaitGroup
	queue := make(chan repo, *workers*2)
	go feedChannel(queue, *oneshot)
	launchWorkers(*workers, &wg, queue)

	// Signal handling
	signals := make(chan os.Signal, 10)
	signal.Notify(signals, syscall.SIGUSR1, syscall.SIGUSR2)
	go handleSignals(signals, *conf, queue)

	// Hang for now but later this should do some checking for
	// signals that may be sent to the processes as well as managing
	// when the queue should have jobs pushed onto it again
	wg.Wait()
}
