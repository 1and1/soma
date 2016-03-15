package main

import (
	"github.com/codegangsta/cli"
)

func registerFlags(app cli.App) *cli.App {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host, H",
			Value: "localhost",
			Usage: "database server host",
		},
		cli.IntFlag{
			Name:  "port, p",
			Value: 5432,
			Usage: "database server port",
		},
		cli.StringFlag{
			Name:  "database, d",
			Value: "soma",
			Usage: "database name",
		},
		cli.StringFlag{
			Name:  "user, u",
			Value: "soma_dba",
			Usage: "database user name",
		},
		cli.BoolFlag{
			Name:  "password, P",
			Usage: "prompt for password",
		},
		cli.IntFlag{
			Name:  "timeout, t",
			Value: 3,
			Usage: "connect timeout in seconds",
		},
		cli.StringFlag{
			Name:  "tls, T",
			Value: "verify-full",
			Usage: "TLS connection mode setting",
		},
		cli.BoolFlag{
			Name:  "no-execute, n",
			Usage: "print SQL statements",
		},
		cli.StringFlag{
			Name:  "config, c",
			Value: "${HOME}/.soma/somadbctl.conf",
			Usage: "configuration file location",
		},
	}
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
