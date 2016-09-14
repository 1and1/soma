package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func cmdExperiment(c *cli.Context) error {

	fmt.Println("This is Experimental!")

	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
