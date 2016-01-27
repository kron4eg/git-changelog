# git changelog - Automatic markdown changelog generator

## Preconditions

This will help you if following preconditiona are met:

* Your project follows (semver)[http://semver.org]
* Semver versions are git tags
* You accept changes exclusevely via PRs (i.e. there are git merges)
* Each PR has tags in the beggining of git title (ex. [fixes] [bug] [new feature] etc)

## Installation

    $ go get github.com/kron4eg/git-changelog
    $ mv $GOPATH/bin/git-changelog /usr/loca/bin/ # actually anywhere in the $PATH

## Usage

    $ cd your/git/repo
    $ git changelog > CHANGELOG.md
    $ git add CHANGELOG.md; git commit

## Limitations

Currently only one tag per merge (PR) is supported
