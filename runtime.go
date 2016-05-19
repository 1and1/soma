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

func runtime(action cli.ActionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		fmt.Println("Do stuff here")

		return action(c)
	}
}

func boottime(action cli.ActionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		var (
			err  error
			mode uint64
			resp *resty.Response
		)
		if err := configSetup(c); err != nil {
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
		if err = store.EnsureBuckets(); err != nil {
			fmt.Fprintf(os.Stderr, "Database bucket error: %s\n", err)
			os.Exit(1)
		}

		utl.SetUrl(Cfg.Api)

		Client = resty.New().SetRESTMode().
			SetHeader(`X-Client`, `0.4.8`).
			SetHostURL(utl.ApiUrl.String())

			// check configured API
		if resp, err = Client.R().Head(`/`); err != nil {
			utl.AbortOnError(err)
		} else if resp.StatusCode() != 204 {
			utl.Abort(fmt.Sprintf("API Url returned %d instead of 200. Sure this is SOMA?\n",
				resp.StatusCode()))
		}

		// check who we talked to
		if resp.Header().Get(`X-Powered-By`) != `SOMA Configuration System` {
			utl.Abort(`Just FYI, at the end of that API URL is not SOMA`)
		}

		return action(c)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
