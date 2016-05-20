package main

import (
	"fmt"
	"os"

	"gopkg.in/resty.v0"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
)

var Client *resty.Client

// initCommon provides common startup initialization
func initCommon(c *cli.Context) {
	var (
		err  error
		resp *resty.Response
	)
	if err = configSetup(c); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read the configuration: "+
			"%s\n", err.Error())
		os.Exit(1)
	}

	// open database
	if err = store.Open(
		Cfg.Run.PathBoltDB,
		os.FileMode(uint32(Cfg.Run.ModeBoltDB)),
		&bolt.Options{Timeout: Cfg.Run.TimeoutBoltDB},
	); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open database: %s\n", err)
		os.Exit(1)
	}

	// ensure database content structure is in place
	if err = store.EnsureBuckets(); err != nil {
		fmt.Fprintf(os.Stderr, "Database bucket error: %s\n", err)
		os.Exit(1)
	}

	// setup our REST client
	Client = resty.New().SetRESTMode().
		SetTimeout(Cfg.Run.TimeoutResty).
		SetDisableWarn(true).
		SetHeader(`User-Agent`, `somaadm 0.4.8`).
		SetHostURL(Cfg.Run.SomaAPI.String()).
		SetContentLength(true)

	if Cfg.Run.SomaAPI.Scheme == `https` {
		Client = Client.SetRootCertificate(Cfg.Run.CertPath)
	}

	// check configured API
	if resp, err = Client.R().Head(`/`); err != nil {
		fmt.Fprintf(os.Stderr, "Error tasting the API endpoint: %s\n",
			err.Error())
	} else if resp.StatusCode() != 204 {
		fmt.Fprintf(os.Stderr, "Error, API Url returned %d instead of 204."+
			" Sure this is SOMA?\n", resp.StatusCode())
		os.Exit(1)
	}

	// check who we talked to
	if resp.Header().Get(`X-Powered-By`) != `SOMA Configuration System` {
		fmt.Fprintf(os.Stderr, `Just FYI, at the end of that API URL`+
			` is not SOMA`)
		os.Exit(1)
	}
}

// boottime is the pre-run target for bootstrapping SOMA
func boottime(action cli.ActionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		initCommon(c)

		return action(c)
	}
}

// runtime is the regular pre-run target
func runtime(action cli.ActionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		initCommon(c)

		return action(c)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
