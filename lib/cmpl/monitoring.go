package cmpl

import "github.com/codegangsta/cli"

func MonitoringCreate(c *cli.Context) {
	Generic(c, []string{`mode`, `contact`, `team`, `callback`})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
