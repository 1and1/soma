package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/resty.v0"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
)

var Client *resty.Client

func runtimePreCmd(c *cli.Context) error {
	var (
		err  error
		mode uint64
	)
	err = configSetup(c)
	if err != nil {
		return err
	}
	dbTimeout := time.Duration(Cfg.BoltDB.Timeout)
	if mode, err = strconv.ParseUint(Cfg.BoltDB.Mode, 8, 32); err != nil {
		fmt.Fprintf(os.Stderr,
			"Failed to parse configuration field boltdb.mode: %s\n", err)
		os.Exit(1)
	}
	if err = store.Open(
		Cfg.Run.PathBoltDB,
		os.FileMode(uint32(mode)),
		&bolt.Options{Timeout: dbTimeout * time.Second},
	); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open database: %s\n", err)
		os.Exit(1)
	}

	//
	//initLogFile()

	//
	utl.SetUrl(Cfg.Api)
	utl.SetPropertyTypes([]string{"system", "service", "custom", "oncall"})
	utl.SetViews([]string{"internal", "external", "local", "any"})

	//
	return nil
}

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
		&bolt.Options{Timeout: Cfg.Run.TimeoutBoltDB * time.Second},
	); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open database: %s\n", err)
		os.Exit(1)
	}

	// ensure database content structure is in place
	if err = store.EnsureBuckets(); err != nil {
		fmt.Fprintf(os.Stderr, "Database bucket error: %s\n", err)
		os.Exit(1)
	}

	// set the configured API endpoint
	utl.SetUrl(Cfg.Api)

	// setup our REST client
	Client = resty.New().SetRESTMode().
		SetDisableWarn(true).
		SetHeader(`User-Agent`, `somaadm 0.4.8`).
		SetHostURL(utl.ApiUrl.String())

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
