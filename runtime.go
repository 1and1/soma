package main

import (
	"github.com/codegangsta/cli"
)

func runtimePreCmd(c *cli.Context) error {
	//
	err := configSetup(c)
	if err != nil {
		return err
	}

	//
	initLogFile()

	//
	return nil
}
