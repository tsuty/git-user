# git-user

[![Build Status](https://travis-ci.org/tsuty/git-user.svg?branch=master)](https://travis-ci.org/tsuty/git-user)

git client multi-user support command.

We use some git hosting services. But we use not always same user name and email. 


## Install

```bash
go get -u github.com/tsuty/git-user
```

## Example

First, set your user info each hosting services.

```bash
git-user set -u git@github.com:* yourname yourname@example.com
git-user set -u git@gitlab.com:* othername othername@example.com
git-user set -u git@bitbucket.org:* somename somename@example.com
```

If you set user.name user.email to global conf, delete from global conf.

Sync git-user conf to local conf.

```bash
cd your_repository
git-user sync
```

show local conf

```bash
git config --local --get-regexp user*
# or 
git-user local
```

### Useful

If you use
[git-prompt](https://github.com/git/git/blob/master/contrib/completion/git-prompt.sh) ([git completion](https://github.com/git/git/blob/master/contrib/completion))
and setting prompt.

```bash
source ~/.git-prompt.sh
PS1='[\u@\h \W$(__git_ps1 " (%s)")]\$ '  
```

Useful command `git-user print`.

this command show local conf for prompt. and sync git-user conf to local conf automatically.

```bash
PS1='[\u@\h \W$(__git_ps1 " (%s)")]$(git-user print)\$ '
```

## bash complete

```bash
ln -s $GOPATH/src/github.com/tsuty/git-user/autocomplete/bash_autocomplete \
    /etc/bash_completion.d/git-user
```

