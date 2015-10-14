package main

import (
	"github.com/codegangsta/cli"
	"log"
	"strconv"
)

func configSetup(c *cli.Context) error {
	log.Print("Starting runtime config initialization")

	// try loading a configuration file
	err := Cfg.populateFromFile(c.GlobalString("config"))
	if err != nil {
		if c.GlobalIsSet("config") {
			// missing cli argument file is fatal error
			log.Fatal(err)
		}
		log.Print("No configuration file found")
	}

	// finish setting up runtime configuration
	params := []string{"api", "timeout", "user"}

	for p := range params {
		// update configuration with cli argument overrides
		if c.GlobalIsSet(params[p]) {
			log.Printf("Setting cli value override for %s", params[p])
			switch params[p] {
			case "api":
				Cfg.Api = c.GlobalString(params[p])
			case "timeout":
				Cfg.Timeout = strconv.Itoa(c.GlobalInt(params[p]))
			case "user":
				Cfg.Auth.User = c.GlobalString(params[p])
			}
			continue
		}
	}

	// TODO
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
