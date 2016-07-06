package cmpl

import "github.com/codegangsta/cli"

func TeamCreate(c *cli.Context) {
	Generic(c, []string{`ldap`, `system`})
}

func TeamUpdate(c *cli.Context) {
	Generic(c, []string{`name`, `ldap`, `system`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
