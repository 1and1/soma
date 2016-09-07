package cmpl

import "github.com/codegangsta/cli"

func OpsRepoRebuild(c *cli.Context) {
	Generic(c, []string{`level`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
