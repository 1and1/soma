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

func cmdPermissionCategoryAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewCategoryRequest()
	req.Category.Name = c.Args().First()

	if resp, err := adm.PostReqBody(req, `/category/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdPermissionCategoryDel(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/category/%s", esc)
	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdPermissionCategoryList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	if resp, err := adm.GetReq(`/category/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdPermissionCategoryShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/category/%s", esc)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `show`)
	}
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
	if sl, ok := opts[`grants`]; ok && len(sl) > 0 {
		req.Permission.Grants = opts[`grants`][0]
	}

	if resp, err := adm.PostReqBody(req, `/permission/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdPermissionDel(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/permission/%s", esc)
	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdPermissionList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	if resp, err := adm.GetReq(`/permission/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdPermissionShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/permission/%s", esc)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `show`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
