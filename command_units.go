package main

import (
	"fmt"
	"net/url"

	"github.com/codegangsta/cli"
)

func registerUnits(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  "units",
				Usage: "SUBCOMMANDS for metric units",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new metric unit",
						Action:       runtime(cmdUnitCreate),
						BashComplete: cmpl.Name,
					},
					{
						Name:   "delete",
						Usage:  "Delete a metric unit",
						Action: runtime(cmdUnitDelete),
					},
					{
						Name:   "list",
						Usage:  "List metric units",
						Action: runtime(cmdUnitList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a metric unit",
						Action: runtime(cmdUnitShow),
					},
				},
			},
		}...,
	)
	return &app
}

func cmdUnitCreate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	key := []string{"name"}

	opts := utl.ParseVariadicArguments(key, key, key, c.Args().Tail())

	req := proto.Request{}
	req.Unit = &proto.Unit{}
	req.Unit.Unit = c.Args().First()
	req.Unit.Name = opts["name"][0]

	resp := utl.PostRequestWithBody(Client, req, "/units/")
	fmt.Println(resp)
	return nil
}

func cmdUnitDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/units/%s", esc)

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdUnitList(c *cli.Context) error {
	resp := utl.GetRequest(Client, "/units/")
	fmt.Println(resp)
	return nil
}

func cmdUnitShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/units/%s", esc)

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
