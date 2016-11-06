/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/codegangsta/cli"
)

func registerInstances(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  `instances`,
				Usage: `SUBCOMMANDS for check instances`,
				Subcommands: []cli.Command{
					{
						Name:   `cascade-delete`,
						Usage:  `Delete check configuration that created the instance`,
						Action: runtime(cmdInstanceCascade),
					},
					{
						Name:   `list`,
						Usage:  `List all check instances`,
						Action: runtime(cmdInstanceList),
					},
					{
						Name:   `show`,
						Usage:  `Show details about a check instance`,
						Action: runtime(cmdInstanceShow),
					},
					{
						Name:   `versions`,
						Usage:  `Show version history of a check instance`,
						Action: runtime(cmdInstanceVersion),
					},
				},
			},
		}...,
	)
	return &app
}

func cmdInstanceCascade(c *cli.Context) error {
	return fmt.Errorf(`Not implemented.`)
}

func cmdInstanceList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/instances/`, `list`, nil, c)
}

func cmdInstanceShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf("Argument is not a UUID: %s",
			c.Args().First())
	}

	path := fmt.Sprintf("/instances/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdInstanceVersion(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf("Argument is not a UUID: %s",
			c.Args().First())
	}

	path := fmt.Sprintf("/instances/%s/versions", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
