package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerEntities(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  "entity",
				Usage: "SUBCOMMANDS for entities (object types)",
				Subcommands: []cli.Command{
					{
						Name:   "add",
						Usage:  "Add a new entity",
						Action: runtime(cmdEntityAdd),
					},
					{
						Name:   "remove",
						Usage:  "Remove an existing entity",
						Action: runtime(cmdEntityRemove),
					},
					{
						Name:         "rename",
						Usage:        "Rename an existing entity",
						Action:       runtime(cmdEntityRename),
						BashComplete: cmpl.To,
					},
					{
						Name:   "list",
						Usage:  "List all entities",
						Action: runtime(cmdEntityList),
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific entity",
						Action: runtime(cmdEntityShow),
					},
				},
			},
		}...,
	)
	return &app
}

func cmdEntityAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewEntityRequest()
	req.Entity.Name = c.Args().First()

	return adm.Perform(`postbody`, `/entity/`, `command`, req, c)
}

func cmdEntityRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/entity/%s", c.Args().First())
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdEntityRename(c *cli.Context) error {
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

	path := fmt.Sprintf("/entity/%s", c.Args().First())
	return adm.Perform(`putbody`, path, `command`, req, c)
}

func cmdEntityList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/entity/`, `list`, nil, c)
}

func cmdEntityShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/entity/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
