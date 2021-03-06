package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
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
	unique := []string{`cardinality`}
	required := []string{`cardinality`}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	switch opts["cardinality"][0] {
	case "once":
	case "multi":
	default:
		return fmt.Errorf("Illegal value for cardinality: %s."+
			" Accepted: once, multi", opts["cardinality"][0])
	}

	req := proto.Request{
		Attribute: &proto.Attribute{
			Name:        c.Args().First(),
			Cardinality: opts["cardinality"][0],
		},
	}

	// check attribute length
	if err := adm.ValidateRuneCount(
		req.Attribute.Name,
		128,
	); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/attributes/`, `command`, req, c)
}

func cmdAttributeDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/attributes/%s", c.Args().First())
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdAttributeList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/attributes/`, `list`, nil, c)
}

func cmdAttributeShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/attributes/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
