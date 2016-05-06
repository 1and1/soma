package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func registerUnits(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:   "units",
				Usage:  "SUBCOMMANDS for metric units",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new metric unit",
						Action: cmdUnitCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a metric unit",
						Action: cmdUnitDelete,
					},
					{
						Name:   "list",
						Usage:  "List metric units",
						Action: cmdUnitList,
					},
					{
						Name:   "show",
						Usage:  "Show details about a metric unit",
						Action: cmdUnitShow,
					},
				},
			},
		}...,
	)
	return &app
}

func cmdUnitCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	key := []string{"name"}

	opts := utl.ParseVariadicArguments(key, key, key, c.Args().Tail())

	req := somaproto.ProtoRequestUnit{}
	req.Unit = &somaproto.ProtoUnit{}
	req.Unit.Unit = c.Args().First()
	req.Unit.Name = opts["name"][0]

	resp := utl.PostRequestWithBody(req, "/units/")
	fmt.Println(resp)
}

func cmdUnitDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/units/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdUnitList(c *cli.Context) {
	resp := utl.GetRequest("/units/")
	fmt.Println(resp)
}

func cmdUnitShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/units/%s", c.Args().First())

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
