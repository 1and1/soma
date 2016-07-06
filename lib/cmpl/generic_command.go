package cmpl

import "github.com/codegangsta/cli"

func In(c *cli.Context) {
	Generic(c, []string{`in`})
}

func InTo(c *cli.Context) {
	Generic(c, []string{`in`, `to`})
}

func InFrom(c *cli.Context) {
	Generic(c, []string{`in`, `from`})
}

func InFromView(c *cli.Context) {
	Generic(c, []string{`in`, `from`, `view`})
}

func FromView(c *cli.Context) {
	Generic(c, []string{`from`, `view`})
}

func ValidityCreate(c *cli.Context) {
	Generic(c, []string{`on`, `direct`, `inherited`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
