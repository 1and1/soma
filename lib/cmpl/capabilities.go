package cmpl

import "github.com/codegangsta/cli"

func CapabilityDeclare(c *cli.Context) {
	Generic(c, []string{`metric`, `view`, `thresholds`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
