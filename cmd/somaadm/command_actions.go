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
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/internal/help"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerAction(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  `actions`,
				Usage: `SUBCOMMANDS for permission actions`,
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Add a permission action to a section`,
						Description:  help.Text(`ActionsAdd`),
						Action:       runtime(cmdActionAdd),
						BashComplete: cmpl.To,
					},
					{
						Name:         `remove`,
						Usage:        `Remove a permission action from a section`,
						Description:  help.Text(`ActionsRemove`),
						Action:       runtime(cmdActionRemove),
						BashComplete: cmpl.From,
					},
					{
						Name:         `list`,
						Usage:        `List permission actions in a section`,
						Description:  help.Text(`ActionsList`),
						Action:       runtime(cmdActionList),
						BashComplete: cmpl.Direct_In,
					},
					{
						Name:         `show`,
						Usage:        `Show details about a permission action`,
						Description:  help.Text(`ActionsShow`),
						Action:       runtime(cmdActionShow),
						BashComplete: cmpl.In,
					},
				},
			},
		}...,
	)
	return &app
}

func cmdActionAdd(c *cli.Context) error {
	var (
		err       error
		sectionId string
	)
	unique := []string{`to`}
	required := []string{`to`}
	opts := make(map[string][]string)
	if err = adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if sectionId, err = adm.LookupSectionId(
		opts[`to`][0],
	); err != nil {
		return err
	}

	req := proto.NewActionRequest()
	req.Action.Name = c.Args().First()
	req.Action.SectionId = sectionId
	path := fmt.Sprintf("/sections/%s/actions/", sectionId)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

func cmdActionRemove(c *cli.Context) error {
	var (
		err                 error
		sectionId, actionId string
	)
	unique := []string{`from`}
	required := []string{`from`}
	opts := make(map[string][]string)
	if err = adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}
	if sectionId, err = adm.LookupSectionId(
		opts[`from`][0],
	); err != nil {
		return err
	}
	if actionId, err = adm.LookupActionId(
		c.Args().First(),
		sectionId,
	); err != nil {
		return err
	}

	path := fmt.Sprintf("/sections/%s/actions/%s", sectionId, actionId)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdActionList(c *cli.Context) error {
	var (
		err       error
		sectionId string
	)
	unique := []string{`in`}
	required := []string{`in`}
	opts := make(map[string][]string)
	if err = adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}
	if sectionId, err = adm.LookupSectionId(
		opts[`in`][0],
	); err != nil {
		return err
	}

	path := fmt.Sprintf("/sections/%s/actions/", sectionId)
	return adm.Perform(`get`, path, `list`, nil, c)
}

func cmdActionShow(c *cli.Context) error {
	var (
		err                 error
		sectionId, actionId string
	)
	unique := []string{`in`}
	required := []string{`in`}
	opts := make(map[string][]string)
	if err = adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}
	if sectionId, err = adm.LookupSectionId(
		opts[`in`][0],
	); err != nil {
		return err
	}
	if actionId, err = adm.LookupActionId(
		c.Args().First(),
		sectionId,
	); err != nil {
		return err
	}

	path := fmt.Sprintf("/sections/%s/actions/%s", sectionId, actionId)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
