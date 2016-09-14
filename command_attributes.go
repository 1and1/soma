package main

import (
	"fmt"

	"github.com/1and1/soma/lib/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerAttributes(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// attributes
			{
				Name:  "attributes",
				Usage: "SUBCOMMANDS for service attributes",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new service attribute",
						Action:       runtime(cmdAttributeCreate),
						BashComplete: cmpl.AttributeCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a service attribute",
						Action: runtime(cmdAttributeDelete),
					},
					{
						Name:   "list",
						Usage:  "List service attributes",
						Action: runtime(cmdAttributeList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a service attribute",
						Action: runtime(cmdAttributeShow),
					},
				},
			}, // end attributes
		}...,
	)
	return &app
}

func cmdAttributeCreate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	multiple := []string{}
	unique := []string{"cardinality"}
	required := []string{"cardinality"}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

	switch opts["cardinality"][0] {
	case "once":
	case "multi":
	default:
		utl.Abort("Illegal value for cardinality")
	}

	req := proto.Request{
		Attribute: &proto.Attribute{
			Name:        c.Args().First(),
			Cardinality: opts["cardinality"][0],
		},
	}
	utl.ValidateRuneCount(req.Attribute.Name, 128)

	resp := utl.PostRequestWithBody(Client, req, "/attributes/")
	fmt.Println(resp)
	return nil
}

func cmdAttributeDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/attributes/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdAttributeList(c *cli.Context) error {
	resp := utl.GetRequest(Client, "/attributes/")
	fmt.Println(resp)
	return nil
}

func cmdAttributeShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/attributes/%s", c.Args().First())

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
