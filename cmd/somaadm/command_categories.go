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
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerCategories(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  "category",
				Usage: "SUBCOMMANDS for permission categories",
				Subcommands: []cli.Command{
					{
						Name:   "add",
						Usage:  "Register a new permission category",
						Action: runtime(cmdPermissionCategoryAdd),
					},
					{
						Name:   "remove",
						Usage:  "Remove an existing permission category",
						Action: runtime(cmdPermissionCategoryDel),
					},
					{
						Name:   "list",
						Usage:  "List all permission categories",
						Action: runtime(cmdPermissionCategoryList),
					},
					{
						Name:   "show",
						Usage:  "Show details for a permission category",
						Action: runtime(cmdPermissionCategoryShow),
					},
				},
			},
		}...,
	)
	return &app
}

func cmdPermissionCategoryAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewCategoryRequest()
	req.Category.Name = c.Args().First()

	return adm.Perform(`postbody`, `/category/`, `command`, req, c)
}

func cmdPermissionCategoryDel(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/category/%s", esc)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdPermissionCategoryList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/category/`, `list`, nil, c)
}

func cmdPermissionCategoryShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/category/%s", esc)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
