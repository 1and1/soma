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

	if resp, err := adm.PostReqBody(req, `/status/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdStatusDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/status/%s", esc)
	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdStatusList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	if resp, err := adm.GetReq(`/status/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdStatusShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/status/%s", esc)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `show`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
