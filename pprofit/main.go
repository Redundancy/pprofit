package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var app *cli.App = cli.NewApp()

func main() {
	app.Name = "pprofit"
	app.Author = "Redundancy"
	app.Email = "dan@the-nexus.co.uk"
	app.Run(os.Args)
}
