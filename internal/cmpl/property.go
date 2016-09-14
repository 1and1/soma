package cmpl

import "github.com/codegangsta/cli"

func Team(c *cli.Context) {
	Generic(c, []string{`team`})
}

func Repository(c *cli.Context) {
	Generic(c, []string{`repository`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
