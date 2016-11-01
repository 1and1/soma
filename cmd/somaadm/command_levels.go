package main

import (
	"fmt"
	"strconv"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerLevels(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// levels
			{
				Name:  "levels",
				Usage: "SUBCOMMANDS for notification levels",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new notification level",
						Action:       runtime(cmdLevelCreate),
						BashComplete: cmpl.LevelCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a notification level",
						Action: runtime(cmdLevelDelete),
					},
					{
						Name:   "list",
						Usage:  "List notification levels",
						Action: runtime(cmdLevelList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a notification level",
						Action: runtime(cmdLevelShow),
					},
				},
			}, // end levels
		}...,
	)
	return &app
}

func cmdLevelCreate(c *cli.Context) error {
	uniqKeys := []string{"shortname", "numeric"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	req := proto.Request{}
	req.Level = &proto.Level{}
	req.Level.Name = c.Args().First()
	req.Level.ShortName = opts["shortname"][0]

	l, err := strconv.ParseUint(opts["numeric"][0], 10, 16)
	if err != nil {
		return fmt.Errorf(
			"Syntax error, numeric argument not numeric: %s",
			err.Error())
	}
	req.Level.Numeric = uint16(l)

	if resp, err := adm.PostReqBody(req, `/levels/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdLevelDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/levels/%s", c.Args().First())
	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdLevelList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	if resp, err := adm.GetReq(`/levels/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdLevelShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/levels/%s", c.Args().First())
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `show`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
