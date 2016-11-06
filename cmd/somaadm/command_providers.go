package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerProviders(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// providers
			{
				Name:  "providers",
				Usage: "SUBCOMMANDS for metric providers",
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new metric provider",
						Action: runtime(cmdProviderCreate),
					},
					{
						Name:   "delete",
						Usage:  "Delete a metric provider",
						Action: runtime(cmdProviderDelete),
					},
					{
						Name:   "list",
						Usage:  "List metric providers",
						Action: runtime(cmdProviderList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a metric provider",
						Action: runtime(cmdProviderShow),
					},
				},
			}, // end providers
		}...,
	)
	return &app
}

func cmdProviderCreate(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.Request{}
	req.Provider = &proto.Provider{}
	req.Provider.Name = c.Args().First()

	return adm.Perform(`postbody`, `/providers/`, `command`, req, c)
}

func cmdProviderDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/providers/%s", c.Args().First())
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdProviderList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/providers/`, `list`, nil, c)
}

func cmdProviderShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/providers/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
