package main

import (
	"fmt"

	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerModes(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// modes
			{
				Name:  "modes",
				Usage: "SUBCOMMANDS for monitoring system modes",
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new monitoring system mode",
						Action: runtime(cmdModeCreate),
					},
					{
						Name:   "delete",
						Usage:  "Delete a monitoring system mode",
						Action: runtime(cmdModeDelete),
					},
					{
						Name:   "list",
						Usage:  "List monitoring system modes",
						Action: runtime(cmdModeList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a monitoring mode",
						Action: runtime(cmdModeShow),
					},
				},
			}, // end modes
		}...,
	)
	return &app
}

func cmdModeCreate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	req := proto.Request{}
	req.Mode = &proto.Mode{}
	req.Mode.Mode = c.Args().First()

	resp := utl.PostRequestWithBody(Client, req, "/modes/")
	fmt.Println(resp)
	return nil
}

func cmdModeDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/modes/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdModeList(c *cli.Context) error {
	resp := utl.GetRequest(Client, "/modes/")
	fmt.Println(resp)
	return nil
}

func cmdModeShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/modes/%s", c.Args().First())

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
