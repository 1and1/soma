package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
	"strconv"
)

func configSetup(c *cli.Context) error {

	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not determine home directory: %s\n", err.Error())
		os.Exit(1)
	}
	var confPath string
	if c.GlobalIsSet("config") {
		confPath = path.Join(home, ".somaadm", c.GlobalString("config"))
	} else {
		confPath = path.Join(home, ".somaadm", "somaadm.conf")
	}

	// try loading a configuration file
	err = Cfg.populateFromFile(confPath)
	if err != nil {
		if c.GlobalIsSet("config") {
			// missing configuration file is only a error if set via cli
			fmt.Fprintf(os.Stderr, "Error opening config: %s\n", confPath)
			os.Exit(1)
		}
	}

	// finish setting up runtime configuration
	params := []string{"api", "timeout", "user", "logdir", "jobdb"}

	for p := range params {
		// update configuration with cli argument overrides
		if c.GlobalIsSet(params[p]) {
			switch params[p] {
			case "api":
				Cfg.Api = c.GlobalString(params[p])
			case "timeout":
				Cfg.Timeout = strconv.Itoa(c.GlobalInt(params[p]))
			case "user":
				Cfg.Auth.User = c.GlobalString(params[p])
			case "logdir":
				Cfg.LogDir = c.GlobalString(params[p])
			case "jobdb":
				Cfg.JobDb = c.GlobalString(params[p])
			}
			continue
		}
		// set default values for unset configuration parameters
		switch params[p] {
		case "api":
			if Cfg.Api == "" {
				Cfg.Api = "http://localhost.my.domain:9876/"
			}
		case "timeout":
			if Cfg.Timeout == "" {
				Cfg.Timeout = strconv.Itoa(5)
			}
		case "user":
			if Cfg.Auth.User == "" {
				Cfg.Auth.User = "admin_fooname"
			}
		case "logdir":
			if Cfg.LogDir == "" {
				Cfg.LogDir = "logs"
			}
		case "jobdb":
			if Cfg.JobDb == "" {
				Cfg.JobDb = "jobs"
			}
		}
	}

	Cfg.Run.PathLogs = path.Join(home, ".somaadm", Cfg.LogDir)
	Cfg.Run.PathLevelDB = path.Join(home, ".somaadm", Cfg.JobDb)

	// TODO prompt for Password
	if Cfg.Auth.Pass == "" {
		return errors.New("Password required")
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
