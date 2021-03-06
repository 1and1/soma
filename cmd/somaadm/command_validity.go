package main

import (
	"fmt"
	"net/url"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerValidity(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  "validity",
				Usage: "SUBCOMMANDS for system property validity",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new system property validity",
						Action:       runtime(cmdValidityCreate),
						BashComplete: cmpl.ValidityCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a system property validity",
						Action: runtime(cmdValidityDelete),
					},
					{
						Name:   "list",
						Usage:  "List system property validity records",
						Action: runtime(cmdValidityList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a system property validity",
						Action: runtime(cmdValidityShow),
					},
				},
			},
		}...,
	)
	return &app
}

func cmdValidityCreate(c *cli.Context) error {
	unique := []string{"on", "direct", "inherited"}
	required := []string{"on", "direct", "inherited"}
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

	req := proto.Request{}
	req.Validity = &proto.Validity{
		SystemProperty: c.Args().First(),
		ObjectType:     opts["on"][0],
	}

	if err := adm.ValidateBool(opts[`direct`][0],
		&req.Validity.Direct); err != nil {
		return err
	}
	if err := adm.ValidateBool(opts[`inherited`][0],
		&req.Validity.Inherited); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/validity/`, `command`, req, c)
}

func cmdValidityDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/validity/%s", esc)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdValidityList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/validity/`, `list`, nil, c)
}

func cmdValidityShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/validity/%s", esc)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
