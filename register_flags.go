package main

import (
	"github.com/codegangsta/cli"
)

func registerFlags(app cli.App) *cli.App {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "user, u",
			Usage: "username for authentication",
		},
		cli.IntFlag{
			Name:  "timeout, t",
			Usage: "connect timeout in seconds",
		},
		cli.StringFlag{
			Name:  "host, H",
			Usage: "API URI to connect to",
		},
		cli.StringFlag{
			Name:  "dbdir, d",
			Usage: "name of the db subdirectory",
		},
		cli.StringFlag{
			Name:  "logdir, l",
			Usage: "name of the log subdirectory",
		},
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "configuration file location",
			EnvVar: "SOMA_ADM_CONFIG",
		},
		cli.BoolFlag{
			Name:  "json, J",
			Usage: "output reply as JSON",
		},
	}
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
