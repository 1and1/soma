package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func registerModes(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// modes
			{
				Name:   "modes",
				Usage:  "SUBCOMMANDS for monitoring system modes",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new monitoring system mode",
						Action: cmdModeCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a monitoring system mode",
						Action: cmdModeDelete,
					},
					{
						Name:   "list",
						Usage:  "List monitoring system modes",
						Action: cmdModeList,
					},
					{
						Name:   "show",
						Usage:  "Show details about a monitoring mode",
						Action: cmdModeShow,
					},
				},
			}, // end modes
		}...,
	)
	return &app
}

func cmdModeCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	req := somaproto.ProtoRequestMode{}
	req.Mode = &somaproto.ProtoMode{}
	req.Mode.Mode = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/modes/")
	fmt.Println(resp)
}

func cmdModeDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/modes/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdModeList(c *cli.Context) {
	resp := utl.GetRequest("/modes/")
	fmt.Println(resp)
}

func cmdModeShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/modes/%s", c.Args().First())

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
