package main

import (
	"fmt"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/mitchellh/go-homedir"
)

// This command runs before the config file exists
func cmdClientInit(c *cli.Context) error {
	// get user home directory
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not determine home directory: %s\n", err.Error())
		os.Exit(1)
	}

	// create ~/.soma/ directory
	somaPath := path.Join(home, ".soma")
	err = os.MkdirAll(somaPath, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating path %s: %s\n", somaPath, err.Error())
		os.Exit(1)
	}

	// create ~/.soma/adm/ directory
	somaPath = path.Join(somaPath, "adm")
	err = os.MkdirAll(somaPath, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating path %s: %s\n", somaPath, err.Error())
		os.Exit(1)
	}

	// create ~/.soma/adm/logs directory
	var (
		logsPath string
	)
	if c.GlobalIsSet("logdir") {
		logsPath = path.Join(somaPath, c.GlobalString("logdir"))
	} else {
		logsPath = path.Join(somaPath, "logs")
	}
	err = os.MkdirAll(logsPath, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating path %s: %s\n", logsPath, err.Error())
		os.Exit(1)
	}

	// create ~/.soma/adm/db directory
	var (
		dbPath string
	)
	if c.GlobalIsSet("dbdir") {
		dbPath = path.Join(somaPath, c.GlobalString("dbdir"))
	} else {
		dbPath = path.Join(somaPath, "db")
	}
	err = os.MkdirAll(dbPath, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating path %s: %s\n", dbPath, err.Error())
		os.Exit(1)
	}
	return nil
}
