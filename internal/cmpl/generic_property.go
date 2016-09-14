package cmpl

import "github.com/codegangsta/cli"

func PropertyAdd(c *cli.Context) {
	genericPropertyAdd(c, []string{})
}

func PropertyAddValue(c *cli.Context) {
	genericPropertyAdd(c, []string{`value`})
}

func genericPropertyAdd(c *cli.Context, args []string) {
	Generic(c, append([]string{`to`, `in`, `value`, `view`, `inheritance`, `childrenonly`}, args...))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
