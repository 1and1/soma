package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var Cfg Config
var utl util.SomaUtil
var store db.DB

func main() {
	app := cli.NewApp()
	app.Name = "somaadm"
	app.Usage = "SOMA Administrative Interface"
	app.Version = "0.4.8"
	app.EnableBashCompletion = true

	app = registerCommands(*app)
	app = registerFlags(*app)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
