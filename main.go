package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var Cfg Config
var utl util.SomaUtil

func main() {
	app := cli.NewApp()
	app.Name = "somaadm"
	app.Usage = "SOMA Administrative Interface"
	app.Version = "0.0.43"

	app = registerCommands(*app)
	app = registerFlags(*app)

	app.Run(os.Args)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
