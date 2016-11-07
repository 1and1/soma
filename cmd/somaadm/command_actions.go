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

func registerAction(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  `sections`,
				Usage: `SUBCOMMANDS for action sections`,
				Subcommands: []cli.Command{
					{
						Name:   `add`,
						Usage:  `Add a permission section`,
						Action: runtime(cmdActionAdd),
					},
					{
						Name:   `remove`,
						Usage:  `Remove a permission section`,
						Action: runtime(cmdActionRemove),
					},
					{
						Name:   `list`,
						Usage:  `List permission sections`,
						Action: runtime(cmdActionList),
					},
					{
						Name:   `show`,
						Usage:  `Show details about permission section`,
						Action: runtime(cmdActionShow),
					},
				},
			},
		}...,
	)
	return &app
}

func cmdActionAdd(c *cli.Context) error {
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
	//TODO validate opts[`to`][0] as section
	//TODO lookup section_id
	// somaadm actions add $action to $section

	var sectionId string
	req := proto.Request{}
	// req := proto.NewActionRequest()
	// req.Action.Name = c.Args().First()
	// req.Action.SectionId = section_id
	path := fmt.Sprintf("/sections/%s/actions/", sectionId)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

func cmdActionRemove(c *cli.Context) error {
	var (
		//err                 error
		sectionId, actionId string
	)
	//TODO lookup section_id
	//TODO lookup action_id
	// somaadm actions delete $action from $section

	path := fmt.Sprintf("/sections/%s/actions/%s", sectionId, actionId)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdActionList(c *cli.Context) error {
	var (
		//err       error
		sectionId string
	)
	//TODO lookup section_id
	// ParseVariadic: ./somaadm actions list in $section

	path := fmt.Sprintf("/sections/%s/actions/", sectionId)
	return adm.Perform(`get`, path, `list`, nil, c)
}

func cmdActionShow(c *cli.Context) error {
	var (
		//err                 error
		sectionId, actionId string
	)
	//TODO lookup section_id
	//TODO lookup action_id
	// ParseVariadic: ./somaadm actions show $action in $section

	path := fmt.Sprintf("/sections/%s/actions/%s", sectionId, actionId)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
