package cmpl

import "github.com/codegangsta/cli"

func BucketCreate(c *cli.Context) {
	Generic(c, []string{`repository`, `environment`})
}

func BucketRename(c *cli.Context) {
	Generic(c, []string{`to`, `repository`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
