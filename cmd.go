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
	ghRE       = regexp.MustCompile(`Merge pull request (\#\d+) from (.+)\/`)
	actRE      = regexp.MustCompile(`\[([\w\s\.]+)\]`)
	changeRe   = regexp.MustCompile(`\] ([^\]\n]+)`)
	nullbyte   = []byte{0x00}
)

type actionsMap map[string][]string

func changelogCmd(c *cli.Context) {
	debugGit = c.Bool("debug")

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
		merge := git("log", `--pretty=%B%x00`, "--merges", tagRange)
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
	actsMap := make(actionsMap)

	commitsSplit := bytes.Split(commits, nullbyte)
	if len(commitsSplit) == 0 {
		return a
	}

	for _, c := range commitsSplit {
		if ghRE.Match(c) {
			ghParse(c, actsMap)
		} else if actRE.Match(c) {
			regularParse(c, actsMap)
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

func ghParse(b []byte, m actionsMap) {
	match := ghRE.FindAllSubmatch(b, -1)
	prauthor := match[0]
	ghPR, ghAuthor := prauthor[1], prauthor[2]
	actualTitle := bytes.Split(b, []byte("\n\n"))[1]

	act, change := extractAction(actualTitle)
	change = fmt.Sprintf("%s. PR [%s][] by [@%s][].", change, ghPR, ghAuthor)
	m[act] = append(m[act], change)
}

func extractAction(b []byte) (string, string) {
	act := actRE.FindAllSubmatch(b, -1)[0][1]
	change := changeRe.FindAllSubmatch(b, -1)[0][1]
	return string(act), string(change)
}

func regularParse(b []byte, m actionsMap) {
	act, change := extractAction(b)
	m[act] = append(m[act], change)
}
