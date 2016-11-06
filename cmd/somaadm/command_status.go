package main

import (
	"fmt"
	"net/url"

	"github.com/1and1/soma/internal/adm"
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
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.Request{}
	req.Status = &proto.Status{}
	req.Status.Name = c.Args().First()

	return adm.Perform(`postbody`, `/status/`, `command`, req, c)
}

func cmdStatusDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/status/%s", esc)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdStatusList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/status/`, `list`, nil, c)
}

func cmdStatusShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/status/%s", esc)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
