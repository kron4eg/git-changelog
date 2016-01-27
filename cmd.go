package main

import (
	"log"
	"os/exec"
	"sort"

	"github.com/Masterminds/semver"
	"github.com/codegangsta/cli"
)

func changelogCmd(c *cli.Context) {
	// git tag
	// drop non-semver
	// sort tags descending
	// get tag info, git show --pretty=email TAG
	// git log --pretty=email --merges OLDTAG..NEWTAG
	// parse commit messages
	// render CHANGELOG.md

	raw := []string{"v0.1.0", "v0.1.1"}
	vs := make([]*semver.Version, len(raw))

	for i, r := range raw {
		v, err := semver.NewVersion(r)
		if err != nil {
			log.Fatal("Error parsing version: %s", err)
		}

		vs[i] = v
	}

	sort.Sort(semver.Collection(vs))

	for _, v := range vs {
		log.Printf("%#v", v)
	}

}

func git(params ...string) exec.Cmd {
	cmd := exec.Command("git", params...)
	return *cmd
}
