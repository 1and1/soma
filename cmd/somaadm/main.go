package main

import (
	"fmt"
	"os"

	"github.com/1and1/soma/internal/db"
	"github.com/1and1/soma/internal/util"
	"github.com/codegangsta/cli"
)

//go:generate go run ../../script/render_markdown.go ../../docs/somaadm/command_reference ../../internal/help/rendered
//go:generate go-bindata -pkg help -ignore .gitignore -o ../../internal/help/bindata.go -prefix "../../internal/help/rendered/" ../../internal/help/rendered/...

var Cfg Config
var utl util.SomaUtil
var store db.DB

const rfc3339Milli string = "2006-01-02T15:04:05.000Z07:00"

func main() {
	cli.CommandHelpTemplate = `{{.Description}}`

	app := cli.NewApp()
	app.Name = "somaadm"
	app.Usage = "SOMA Administrative Interface"
	app.Version = "0.8.2"
	app.EnableBashCompletion = true

	app = registerCommands(*app)
	app = registerFlags(*app)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
