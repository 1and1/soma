package main

import (
  "os"
  "github.com/codegangsta/cli"
  "database/sql"
)

var Cfg Config
var db *sql.DB

func main() {
  app := cli.NewApp()
  app.Name = "somadbctl"
  app.Usage = "SOMA Database Control Utility"
  app.Version = "0.0.1"

  app = registerCommands(*app)
  app = registerFlags(*app)

  app.Run(os.Args)
}
