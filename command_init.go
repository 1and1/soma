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
	const (
		skelConf = `/usr/share/somaadm/skel/somaadm.conf`
		skelCert = `/usr/share/somaadm/skel/ca.cert.pem`
		defConf  = `somaadm.conf`
		defCert  = `ca.cert.pem`
	)
	var conf, cert bool

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

	// if we find skeleton files in /usr/share we copy them to
	if _, err = os.Stat(`/usr/share/somaadm/skel/somaadm.conf`); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error testing for skeleton config: %s\n", err.Error())
		os.Exit(1)
	} else if !os.IsNotExist(err) {
		if _, err = adm.CopyFile(path.Join(somaPath, defConf), skelConf); err != nil && !os.IsExist(err) {
			fmt.Fprintf(os.Stderr, "Error copying skeleton config: %s\n", err.Error())
			os.Exit(1)
		}
		conf = true
		os.Chmod(path.Join(somaPath, defConf), 0600)
	}

	if _, err = os.Stat(skelCert); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error testing for skeleton certificate: %s\n", err.Error())
		os.Exit(1)
	} else if !os.IsNotExist(err) {
		if _, err = adm.CopyFile(path.Join(somaPath, defCert), skelCert); err != nil && !os.IsExist(err) {
			fmt.Fprintf(os.Stderr, "Error copying skeleton certificate: %s\n", err.Error())
			os.Exit(1)
		}
		cert = true
		os.Chmod(path.Join(somaPath, defCert), 0444)
	}

	if conf && cert {
		return boottime(cmdDummyInit)(c)
	}
	return nil
}

func cmdDummyInit(c *cli.Context) error {
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
