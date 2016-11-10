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
	"net/url"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerPermissions(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  "permissions",
				Usage: "SUBCOMMANDS for permissions",
				Subcommands: []cli.Command{
					{
						Name:         "add",
						Usage:        "Register a new permission",
						Action:       runtime(cmdPermissionAdd),
						BashComplete: cmpl.PermissionAdd,
					},
					{
						Name:   "remove",
						Usage:  "Remove a permission",
						Action: runtime(cmdPermissionDel),
					},
					{
						Name:   "list",
						Usage:  "List all permissions",
						Action: runtime(cmdPermissionList),
					},
					{
						Name:   "show",
						Usage:  "Show details for a permission",
						Action: runtime(cmdPermissionShow),
					},
				},
			}, // end permissions
		}...,
	)
	return &app
}

func cmdPermissionAdd(c *cli.Context) error {
	unique := []string{`category`, `grants`}
	required := []string{`category`}
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

	req := proto.NewPermissionRequest()
	req.Permission.Name = c.Args().First()
	req.Permission.Category = opts[`category`][0]
	//XXX requires update, nwo works differently
	//if sl, ok := opts[`grants`]; ok && len(sl) > 0 {
	//	req.Permission.Grants = opts[`grants`][0]
	//}

	return adm.Perform(`postbody`, `/permission/`, `command`, req, c)
}

func cmdPermissionDel(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/permission/%s", esc)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdPermissionList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/permission/`, `list`, nil, c)
}

func cmdPermissionShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/permission/%s", esc)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
