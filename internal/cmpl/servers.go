package cmpl

import "github.com/codegangsta/cli"

func ServerCreate(c *cli.Context) {
	Generic(c, []string{`assetid`, `datacenter`, `location`, `online`})
}

func ServerUpdate(c *cli.Context) {
	Generic(c, []string{`name`, `assetid`, `datacenter`, `location`, `online`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
