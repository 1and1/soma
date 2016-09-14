package cmpl

import "github.com/codegangsta/cli"

func AttributeCreate(c *cli.Context) {
	Generic(c, []string{`cardinality`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
