package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	log.SetFlags(0)
	app := cli.NewApp()
	app.Name = "git changelog"
	app.Usage = "generate CHANGELOG.md from commit history"
	app.Version = "0.0.1"
	app.Action = changelogCmd
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "debug git commands",
		},
	}
	app.Run(os.Args)

}
