package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/mitchellh/go-homedir"
	"github.com/peterh/liner"
)

func configSetup(c *cli.Context) error {
	if c.GlobalBool(`no-execute`) {
		return nil
	}

	var (
		configFile, home, wd string
		err                  error
	)

	if home, err = homedir.Dir(); err != nil {
		log.Fatal(err)
	}
	if wd, err = os.Getwd(); err != nil {
		log.Fatal(err)
	}

	// try loading a configuration file
	if c.GlobalIsSet(`config`) {
		if path.IsAbs(c.GlobalString(`config`)) {
			configFile = c.GlobalString(`config`)
		} else {
			configFile = path.Join(wd, c.GlobalString(`config`))
		}
	} else {
		configFile = path.Join(home, `.soma`, `dbctl`, `somadbctl.conf`)
	}

	if err = Cfg.populateFromFile(configFile); err != nil {
		if c.GlobalIsSet(`config`) {
			// missing cli argument file is fatal error
			log.Fatal(err)
		}
		log.Print(fmt.Sprintf("No configuration file found: %s", configFile))
	}

	// finish setting up runtime configuration
	params := []string{"host", "port", "database", "user", "timeout", "tls"}

	for p := range params {
		// update configuration with cli argument overrides
		if c.GlobalIsSet(params[p]) {
			log.Printf("Setting cli value override for %s", params[p])
			switch params[p] {
			case "host":
				Cfg.Database.Host = c.GlobalString(params[p])
			case "port":
				Cfg.Database.Port = strconv.Itoa(c.GlobalInt(params[p]))
			case "database":
				Cfg.Database.Name = c.GlobalString(params[p])
			case "user":
				Cfg.Database.User = c.GlobalString(params[p])
			case "timeout":
				Cfg.Timeout = strconv.Itoa(c.GlobalInt(params[p]))
			case "tls":
				Cfg.TlsMode = c.GlobalString(params[p])
			}
			continue
		}
		// set default values for unset configuration parameters
		switch params[p] {
		case "host":
			if Cfg.Database.Host == "" {
				log.Printf("Setting default value for %s", params[p])
				Cfg.Database.Host = "localhost"
			}
		case "port":
			if Cfg.Database.Port == "" {
				log.Printf("Setting default value for %s", params[p])
				Cfg.Database.Port = strconv.Itoa(5432)
			}
		case "database":
			if Cfg.Database.Name == "" {
				log.Printf("Setting default value for %s", params[p])
				Cfg.Database.Name = "soma"
			}
		case "user":
			if Cfg.Database.User == "" {
				log.Printf("Setting default value for %s", params[p])
				Cfg.Database.User = "soma_dba"
			}
		case "timeout":
			if Cfg.Timeout == "" {
				log.Printf("Setting default value for %s", params[p])
				Cfg.Timeout = strconv.Itoa(3)
			}
		case "tls":
			if Cfg.TlsMode == "" {
				log.Printf("Setting default value for %s", params[p])
				Cfg.TlsMode = "verify-full"
			}
		}
	}

	// prompt for password if the cli flag was set
	if c.GlobalBool(`password`) || Cfg.Database.Pass == "" {
		line := liner.NewLiner()
		defer line.Close()
		line.SetCtrlCAborts(true)

		Cfg.Database.Pass, err = line.PasswordPrompt(`Enter database password: `)
		if err == liner.ErrPromptAborted {
			os.Exit(0)
		} else if err != nil {
			log.Fatal(err)
		}
	}

	// abort if we have no connection password at this point
	if Cfg.Database.Pass == "" {
		log.Fatal(`Can not continue without database connection password`)
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
