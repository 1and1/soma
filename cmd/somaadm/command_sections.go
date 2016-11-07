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
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerSection(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  `sections`,
				Usage: `SUBCOMMANDS for action sections`,
				Subcommands: []cli.Command{
					{
						Name:   `add`,
						Usage:  `Add a permission section`,
						Action: runtime(cmdSectionAdd),
					},
					{
						Name:   `remove`,
						Usage:  `Remove a permission section`,
						Action: runtime(cmdSectionRemove),
					},
					{
						Name:   `list`,
						Usage:  `List permission sections`,
						Action: runtime(cmdSectionList),
					},
					{
						Name:   `show`,
						Usage:  `Show details about permission section`,
						Action: runtime(cmdSectionShow),
					},
				},
			},
		}...,
	)
	return &app
}

func cmdSectionAdd(c *cli.Context) error {
	unique := []string{`to`}
	required := []string{`to`}
	opts := make(map[string][]string)
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	//TODO validate opts[`to`][0] as category

	req := proto.Request{} // NewSectionRequest()
	// req := proto.NewSectionRequest()
	// req.Section.Name = c.Args().First()
	// req.Section.Category = opts[`to`][0]
	return adm.Perform(`postbody`, `/sections/`, `command`, req, c)
}

func cmdSectionRemove(c *cli.Context) error {
	var (
		err       error
		sectionId string
	)
	if err = adm.VerifySingleArgument(c); err != nil {
		return err
	}
	//TODO lookup section_id

	path := fmt.Sprintf("/sections/%s", sectionId)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdSectionList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/sections/`, `list`, nil, c)
}

func cmdSectionShow(c *cli.Context) error {
	var (
		err       error
		sectionId string
	)
	if err = adm.VerifySingleArgument(c); err != nil {
		return err
	}
	//TODO lookup section_id

	path := fmt.Sprintf("/sections/%s", sectionId)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix