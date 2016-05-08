package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func registerStatus(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// status
			{
				Name:   "status",
				Usage:  "SUBCOMMANDS for check instance status",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Add a check instance status",
						Action: cmdStatusCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a check instance status",
						Action: cmdStatusDelete,
					},
					{
						Name:   "list",
						Usage:  "List check instance status",
						Action: cmdStatusList,
					},
					{
						Name:   "show",
						Usage:  "Show details about a check instance status",
						Action: cmdStatusShow,
					},
				},
			}, // end status
		}...,
	)
	return &app
}

func cmdStatusCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	req := proto.Request{}
	req.Status = &proto.Status{}
	req.Status.Name = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/status/")
	fmt.Println(resp)
}

func cmdStatusDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/status/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdStatusList(c *cli.Context) {
	resp := utl.GetRequest("/status/")
	fmt.Println(resp)
}

func cmdStatusShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/status/%s", c.Args().First())

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
