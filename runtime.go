package main

import (
	"github.com/codegangsta/cli"
)

func runtimePreCmd(c *cli.Context) error {
	err := configSetup(c)
	if err != nil {
		return err
	}

	//
	initLogFile()

	//
	utl.SetUrl(Cfg.Api)
	utl.SetPropertyTypes([]string{"system", "service", "custom", "oncall"})
	utl.SetViews([]string{"internal", "external", "local", "any"})

	//
	return nil
}
