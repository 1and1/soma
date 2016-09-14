package cmpl

import (
	"github.com/codegangsta/cli"
)

func LevelCreate(c *cli.Context) {
	Generic(c, []string{`shortname`, `numeric`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
