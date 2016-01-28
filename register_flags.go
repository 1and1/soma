package main

import (
	"github.com/codegangsta/cli"
)

func registerFlags(app cli.App) *cli.App {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "user, u",
			Usage:  "username for authentication",
			EnvVar: "SOMA_ADM_USER,USER",
		},
		cli.BoolFlag{
			Name:  "password, p",
			Usage: "prompt for password",
		},
		cli.IntFlag{
			Name:   "timeout, t",
			Usage:  "connect timeout in seconds",
			EnvVar: "SOMA_ADM_TIMEOUT",
		},
		cli.StringFlag{
			Name:   "api, a",
			Usage:  "API URI to connect to",
			EnvVar: "SOMA_ADM_API",
		},
		cli.StringFlag{
			Name:   "jobdb, j",
			Usage:  "name of the jobs data subdirectory",
			EnvVar: "SOMA_ADM_JOBSDB",
		},
		cli.StringFlag{
			Name:   "logdir, l",
			Usage:  "name of the log subdirectory",
			EnvVar: "SOMA_ADM_LOGDIR",
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
