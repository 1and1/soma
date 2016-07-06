package cmpl

import "github.com/codegangsta/cli"

func MetricCreate(c *cli.Context) {
	GenericMulti(c, []string{`unit`, `description`}, []string{`package`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
