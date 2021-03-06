package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerTypes(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// types
			{
				Name:  "types",
				Usage: "SUBCOMMANDS for object types",
				Subcommands: []cli.Command{
					{
						Name:   "add",
						Usage:  "Add a new object type",
						Action: runtime(cmdObjectTypesAdd),
					},
					{
						Name:   "remove",
						Usage:  "Remove an existing object type",
						Action: runtime(cmdObjectTypesRemove),
					},
					{
						Name:         "rename",
						Usage:        "Rename an existing object type",
						Action:       runtime(cmdObjectTypesRename),
						BashComplete: cmpl.To,
					},
					{
						Name:   "list",
						Usage:  "List all object types",
						Action: runtime(cmdObjectTypesList),
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific object type",
						Action: runtime(cmdObjectTypesShow),
					},
				},
			}, // end types
		}...,
	)
	return &app
}

func cmdObjectTypesAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewEntityRequest()
	req.Entity.Name = c.Args().First()

	return adm.Perform(`postbody`, `/objtypes/`, `command`, req, c)
}

func cmdObjectTypesRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/objtypes/%s", c.Args().First())
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdObjectTypesRename(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`to`},
		[]string{`to`},
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewEntityRequest()
	req.Entity.Name = opts[`to`][0]

	path := fmt.Sprintf("/objtypes/%s", c.Args().First())
	return adm.Perform(`putbody`, path, `command`, req, c)
}

func cmdObjectTypesList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/objtypes/`, `list`, nil, c)
}

func cmdObjectTypesShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/objtypes/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
