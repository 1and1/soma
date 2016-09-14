package cmpl

import "github.com/codegangsta/cli"

func UserAdd(c *cli.Context) {
	Generic(c, []string{`firstname`, `lastname`, `employeenr`, `mailaddr`, `team`, `deleted`, `active`, `system`})
}

func UserUpdate(c *cli.Context) {
	Generic(c, []string{`username`, `firstname`, `lastname`, `employeenr`, `mailaddr`, `team`, `deleted`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
