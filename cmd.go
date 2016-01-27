package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"time"

	"github.com/Masterminds/semver"
	"github.com/codegangsta/cli"
)

var (
	debugGit   bool
	timeLayout = `Mon, 02 Jan 2006 15:04:05 -0700`
	re         = regexp.MustCompile(`Subject: \[PATCH\] (.+)`)
	re2        = regexp.MustCompile(`\[(.+)\] (.+)`)
)

func changelogCmd(c *cli.Context) {
	debugGit = c.Bool("debug")
	// parse commit messages
	// render CHANGELOG.md

	vs := getSemverTags()
	context := make([]Version, 0)

	for i, v := range vs {
		if i == len(vs)-1 {
			// skip last tag
			continue
		}
		tagShow := git("show", `--pretty=%cD`, v.Original())
		tagTime := bytes.Split(tagShow, []byte("\n"))[0]
		tagdatetime, _ := time.Parse(timeLayout, string(tagTime))
		version := Version{
			Version: v.Original(),
			Date:    tagdatetime,
		}
		tagRange := fmt.Sprintf("%s...%s", vs[i+1].Original(), v.Original())
		merge := git("log", "--pretty=email", "--merges", tagRange)
		version.Actions = parseActions(merge)
		context = append(context, version)
	}

	render(context)
}

func git(params ...string) []byte {
	cmd := exec.Command("git", params...)
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	if debugGit {
		log.Printf("%q %q", cmd.Path, cmd.Args)
		log.Printf("%s", output)
	}
	return output
}

func parseActions(commits []byte) []Action {
	a := []Action{}
	actsMap := map[string][]string{}

	for _, rez := range re.FindAllSubmatch(commits, -1) {
		// log.Printf("%s", rez[1])
		for _, inner := range re2.FindAllSubmatch(rez[1], -1) {
			key := string(inner[1])
			actsMap[key] = append(actsMap[key], string(inner[2]))
		}
	}

	for k, v := range actsMap {
		a = append(a, Action{
			Action:  k,
			Changes: v,
		})
	}

	return a
}

func getSemverTags() []*semver.Version {
	rowTags := git("tag")

	tags := bytes.Split(rowTags, []byte("\n"))
	tags = tags[:len(tags)-1]

	vs := make([]*semver.Version, 0)

	for _, r := range tags {
		v, err := semver.NewVersion(string(r))
		if err == nil {
			vs = append(vs, v)
		}
	}

	sort.Sort(sort.Reverse(semver.Collection(vs)))
	return vs
}

func render(c []Version) {
	if err := changeLogTpl.Execute(os.Stdout, c); err != nil {
		log.Fatal(err)
	}
}
