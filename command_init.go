package main

import (
	"fmt"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/mitchellh/go-homedir"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// This command runs before the config file exists
func cmdClientInit(c *cli.Context) {
	// get user home directory
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not determine home directory: %s\n", err.Error())
		os.Exit(1)
	}

	// create ~/.somaadm/ directory
	somaPath := path.Join(home, ".somaadm")
	err = os.MkdirAll(somaPath, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating path %s: %s\n", somaPath, err.Error())
		os.Exit(1)
	}

	// create ~/.somaadm/logs directory
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

	// create LevelDB
	var (
		ldbPath string
		ldbOpt  opt.Options
	)
	if c.GlobalIsSet("jobdb") {
		ldbPath = path.Join(somaPath, c.GlobalString("jobdb"))
	} else {
		ldbPath = path.Join(somaPath, "jobs")
	}
	ldbOpt.ErrorIfExist = true
	db, err := leveldb.OpenFile(ldbPath, &ldbOpt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating LevelDB at %s: %s\n", ldbPath, err.Error())
		os.Exit(1)
	}
	defer db.Close()
}
