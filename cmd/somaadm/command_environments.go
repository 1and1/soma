package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerEnvironments(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// environments
			{
				Name:  "environments",
				Usage: "SUBCOMMANDS for environments",
				Subcommands: []cli.Command{
					{
						Name:   "add",
						Usage:  "Register a new view",
						Action: runtime(cmdEnvironmentsAdd),
					},
					{
						Name:   "remove",
						Usage:  "Remove an existing unused environment",
						Action: runtime(cmdEnvironmentsRemove),
					},
					{
						Name:         "rename",
						Usage:        "Rename an existing environment",
						Action:       runtime(cmdEnvironmentsRename),
						BashComplete: cmpl.To,
					},
					{
						Name:   "list",
						Usage:  "List all available environments",
						Action: runtime(cmdEnvironmentsList),
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific environment",
						Action: runtime(cmdEnvironmentsShow),
					},
				},
			}, // end environments
		}...,
	)
	return &app
}

func cmdEnvironmentsAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewEnvironmentRequest()
	req.Environment.Name = c.Args().First()

	return adm.Perform(`postbody`, `/environments/`, `command`, req, c)
}

func cmdEnvironmentsRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/environments/%s", c.Args().First())
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdEnvironmentsRename(c *cli.Context) error {
	key := []string{`to`}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(opts, []string{}, key, key,
		c.Args().Tail()); err != nil {
		return err
	}

	req := proto.NewEnvironmentRequest()
	req.Environment.Name = opts[`to`][0]

	path := fmt.Sprintf("/environments/%s", c.Args().First())
	return adm.Perform(`putbody`, path, `command`, req, c)
}

func cmdEnvironmentsList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/environments/`, `list`, nil, c)
}

func cmdEnvironmentsShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/environments/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
