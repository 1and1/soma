package main

import (
	"fmt"
	"net/url"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
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
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`name`},
		[]string{`name`},
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.Request{}
	req.Unit = &proto.Unit{}
	req.Unit.Unit = c.Args().First()
	req.Unit.Name = opts["name"][0]

	return adm.Perform(`postbody`, `/units/`, `command`, req, c)
}

func cmdUnitDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/units/%s", esc)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdUnitList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/units/`, `list`, nil, c)
}

func cmdUnitShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/units/%s", esc)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
