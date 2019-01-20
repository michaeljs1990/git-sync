Git Sync
========

Git sync is a tool that I have written to keep my custom hosted git repository in sync
with multiple different sources. It currently is somewhat limited because I have only
added proper support for pulling over http/https. You can use the git or ssh protocol as
well with some features but it's at your own risk.

## Features

|Feature                               |Done|Kinda|Planned|Wishful|
|--------------------------------------|:--:|:---:|:-----:|:-----:|
|Clone repositories from github        |X   |     |       |       |
|Clone repositories from bitbucket     |    |     |       |X      |
|Clone repositories from gitlab        |    |     |       |X      |
|Push repo over http auth              |X   |     |       |       |
|Pull repo over http auth              |X   |     |       |       |
|Push repo over git protocol           |    |X    |       |       |
|Pull repo over git protocol           |    |X    |       |       |
|Push repo over SSH protocol           |    |     |X      |       |
|SSH user/pass auth                    |    |     |X      |       |
|SSH private key auth                  |    |     |X      |       |
|Pull repo over SSH protocol           |    |     |X      |       |
|Config validation                     |    |     |X      |       |
|Prometheus metrics                    |    |     |X      |       |
|Hot config reloading                  |    |     |X      |       |
|Force syncing without process restart |    |     |X      |       |

## Config

Here is an example config file. I tried my best to bubble up error messages but for now
config validation is somewhat limited.

```
---
interval: 3600
path: /var/lib/repos
repos:
- url: https://github.com/michaeljs1990/git-sync.git
  type: github_push_mirror
  path: /var/lib/repos/git-sync
  httpauth:
    user: username
    token: sometoken
  remote: github_mirror
github:
- username: michaeljs1990
  httpauth:
    user: username
    token: sometoken
  repos: true
```

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

## Development

You can use nix-shell to setup the environment that I used however using golang 1.11 and up should work fine.
