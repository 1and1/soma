package main

import (
	"github.com/codegangsta/cli"
)

func dbClose(c *cli.Context) error {
	defer db.Close()

	return nil
}
