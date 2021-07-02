Git Sync
========

Git sync is a tool that I have written to keep my custom hosted git repo in sync with
multiple different sources. It can pull from places such as github at a set interval
as well as push to any git, http, https, or ssh repository. Additinally I have added
some features for populating metadata on cgit repositories that you can turn on if you
want.

This works on all OS but hot reloading of config on windows doesn't work.

## Config

Here is an example config file with secrets replaced. To validate your config you can 
run `git-sync -validate -config file.yaml`.

```
---
interval: 3600
path: /var/lib/repos
repos:
- url: https://review.coreboot.org/coreboot.git
  type: fetch_mirror
  path: /var/lib/repos/mirrors/coreboot
  extras:
    cgitowner: coreboot
    cgitsection: mirrors
  metadata:
  - cgit
- url: https://anongit.gentoo.org/git/repo/gentoo.git
  type: fetch_mirror
  path: /var/lib/repos/mirrors/gentoo
  extras:
    cgitowner: gentoo
    cgitsection: mirrors
  metadata:
  - cgit
- url: https://github.com/michaeljs1990/collins-go-cli.git
  type: push_mirror
  path: /var/lib/repos/xrt0x/collins-go-cli
  httpauth:
    user: michaeljs1990
    token: some-token
  remote: github_mirror
github:
- username: michaeljs1990
  httpauth:
    user: michaeljs1990
    token: some-token
  repos: true
  metadata:
  - cgit
```

## Sync Types

`push_mirror` - Push from the local path provided to the given URL to keep them in sync.

`fetch_mirror` - Pull from the url provided and place it at the provided path.

`github` - pull all repos located under the provided username to the top level path.

## Example

Run one time to test out your config

```
git-sync -config config.yaml -verbose -oneshot
```

Start the proces as a long running service
```
git-sync -config config.yaml -log_file git_sync.log
```

Log to file as JSON
```
git-sync -config config.yaml -log_file git_sync.log -log_format json
```

Run multiple workers to speed things up
```
git-sync -config config.yaml -log_file git_sync.log -log_format json -workers 10
```

Validate your config
```
git-sync -validate
```
## Signals

You can reload the config file for your git-sync process with the use of usr1. This will validate
the config before doing the reload so if you have a bad config it will just output an error but
keep running with the old config.

If you are excited to have something pushed out or want a new update you can hit git-sync with usr2.
This will cause it to sync all repositories again that it knows about.

## Development

Test with golang 1.11 - 1.16

## Feature List

|Feature                               |Done|Kinda|Planned|Wishful|
|--------------------------------------|:--:|:---:|:-----:|:-----:|
|Clone repositories from bitbucket     |    |     |       |X      |
|Clone repositories from gitlab        |    |     |       |X      |
|Add unit tests                        |    |     |       |X      |
|Add file watching for push repos      |    |     |       |X      |
|Clone repositories from github        |X   |     |       |       |
|Push repo over http auth              |X   |     |       |       |
|Pull repo over http auth              |X   |     |       |       |
|Push repo over git protocol           |X   |     |       |       |
|Pull repo over git protocol           |X   |     |       |       |
|Push repo over SSH protocol           |X   |     |       |       |
|SSH user/pass auth                    |X   |     |       |       |
|SSH private key auth                  |X   |     |       |       |
|Pull repo over SSH protocol           |X   |     |       |       |
|Config validation                     |X   |     |       |       |
|Prometheus metrics                    |X   |     |       |       |
|Hot config reloading                  |X   |     |       |       |
|Force syncing without process restart |X   |     |       |       |
|Populate some cgit metadata           |X   |     |       |       |
