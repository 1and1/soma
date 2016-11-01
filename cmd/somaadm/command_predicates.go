package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerPredicates(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// predicates
			{
				Name:  "predicates",
				Usage: "SUBCOMMANDS for threshold predicates",
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Add a predicate",
						Action: runtime(cmdPredicateCreate),
					},
					{
						Name:   "delete",
						Usage:  "Delete a predicate",
						Action: runtime(cmdPredicateDelete),
					},
					{
						Name:   "list",
						Usage:  "List predicates",
						Action: runtime(cmdPredicateList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a predicate",
						Action: runtime(cmdPredicateShow),
					},
				},
			}, // end predicates
		}...,
	)
	return &app
}

func cmdPredicateCreate(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.Request{}
	req.Predicate = &proto.Predicate{}
	req.Predicate.Symbol = c.Args().First()

	if resp, err := adm.PostReqBody(req, `/predicates/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdPredicateDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/predicates/%s", c.Args().First())
	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdPredicateList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	if resp, err := adm.GetReq(`/predicates/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdPredicateShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/predicates/%s", c.Args().First())
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `show`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
