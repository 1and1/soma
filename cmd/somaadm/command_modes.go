package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerModes(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// modes
			{
				Name:  "modes",
				Usage: "SUBCOMMANDS for monitoring system modes",
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new monitoring system mode",
						Action: runtime(cmdModeCreate),
					},
					{
						Name:   "delete",
						Usage:  "Delete a monitoring system mode",
						Action: runtime(cmdModeDelete),
					},
					{
						Name:   "list",
						Usage:  "List monitoring system modes",
						Action: runtime(cmdModeList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a monitoring mode",
						Action: runtime(cmdModeShow),
					},
				},
			}, // end modes
		}...,
	)
	return &app
}

func cmdModeCreate(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.Request{}
	req.Mode = &proto.Mode{}
	req.Mode.Mode = c.Args().First()

	if resp, err := adm.PostReqBody(req, `/modes/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdModeDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/modes/%s", c.Args().First())
	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdModeList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	if resp, err := adm.GetReq(`/modes/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdModeShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/modes/%s", c.Args().First())
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `show`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
