package cmpl

import "github.com/codegangsta/cli"

func Datacenter(c *cli.Context) {
	Generic(c, []string{`datacenter`})
}

func In(c *cli.Context) {
	Generic(c, []string{`in`})
}

func Direct_In(c *cli.Context) {
	GenericDirect(c, []string{`in`})
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

func From(c *cli.Context) {
	Generic(c, []string{`from`})
}

func FromView(c *cli.Context) {
	Generic(c, []string{`from`, `view`})
}

func Name(c *cli.Context) {
	Generic(c, []string{`name`})
}

func To(c *cli.Context) {
	Generic(c, []string{`to`})
}

func Triple_ToOn(c *cli.Context) {
	GenericTriple(c, []string{`to`, `on`})
}

func Triple_FromOn(c *cli.Context) {
	GenericTriple(c, []string{`from`, `on`})
}

func ValidityCreate(c *cli.Context) {
	Generic(c, []string{`on`, `direct`, `inherited`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
