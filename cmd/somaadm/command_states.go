package main

import (
	"fmt"
	"net/url"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerStates(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// states
			{
				Name:  "states",
				Usage: "SUBCOMMANDS for states",
				Subcommands: []cli.Command{
					{
						Name:   "add",
						Usage:  "Add a new object state",
						Action: runtime(cmdObjectStatesAdd),
					},
					{
						Name:   "remove",
						Usage:  "Remove an existing object state",
						Action: runtime(cmdObjectStatesRemove),
					},
					{
						Name:         "rename",
						Usage:        "Rename an existing object state",
						Action:       runtime(cmdObjectStatesRename),
						BashComplete: cmpl.To,
					},
					{
						Name:   "list",
						Usage:  "List all object states",
						Action: runtime(cmdObjectStatesList),
					},
					{
						Name:   "show",
						Usage:  "Show information about an object states",
						Action: runtime(cmdObjectStatesShow),
					},
				},
			}, // end states
		}...,
	)
	return &app
}

func cmdObjectStatesAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewStateRequest()
	req.State.Name = c.Args().First()

	return adm.Perform(`postbody`, `/objstates/`, `command`, req, c)
}

func cmdObjectStatesRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/objstates/%s", esc)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdObjectStatesRename(c *cli.Context) error {
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

	req := proto.NewStateRequest()
	req.State.Name = opts[`to`][0]

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/objstates/%s", esc)
	return adm.Perform(`putbody`, path, `command`, req, c)
}

func cmdObjectStatesList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/objstates/`, `list`, nil, c)
}

func cmdObjectStatesShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/objstates/%s", esc)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
