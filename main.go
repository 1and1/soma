package main

import (
	"database/sql"
	"github.com/codegangsta/cli"
	"os"
)

var Cfg Config
var db *sql.DB

func main() {
	app := cli.NewApp()
	app.Name = "somadbctl"
	app.Usage = "SOMA Database Control Utility"
	app.Version = "0.0.16"

	app = registerCommands(*app)
	app = registerFlags(*app)

	app.Run(os.Args)
}
