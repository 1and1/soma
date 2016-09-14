package main

import (
	"fmt"

	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerStatus(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// status
			{
				Name:  "status",
				Usage: "SUBCOMMANDS for check instance status",
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Add a check instance status",
						Action: runtime(cmdStatusCreate),
					},
					{
						Name:   "delete",
						Usage:  "Delete a check instance status",
						Action: runtime(cmdStatusDelete),
					},
					{
						Name:   "list",
						Usage:  "List check instance status",
						Action: runtime(cmdStatusList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a check instance status",
						Action: runtime(cmdStatusShow),
					},
				},
			}, // end status
		}...,
	)
	return &app
}

func cmdStatusCreate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	req := proto.Request{}
	req.Status = &proto.Status{}
	req.Status.Name = c.Args().First()

	resp := utl.PostRequestWithBody(Client, req, "/status/")
	fmt.Println(resp)
	return nil
}

func cmdStatusDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/status/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdStatusList(c *cli.Context) error {
	resp := utl.GetRequest(Client, "/status/")
	fmt.Println(resp)
	return nil
}

func cmdStatusShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/status/%s", c.Args().First())

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
