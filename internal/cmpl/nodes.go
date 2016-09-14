package cmpl

import "github.com/codegangsta/cli"

func NodeAdd(c *cli.Context) {
	Generic(c, []string{`assetid`, `name`, `team`, `server`, `online`})
}

func NodeUpdate(c *cli.Context) {
	Generic(c, []string{`name`, `assetid`, `server`, `team`, `online`, `deleted`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
