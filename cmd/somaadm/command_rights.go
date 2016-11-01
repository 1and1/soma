package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerRights(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  "rights",
				Usage: "SUBCOMMANDS for rights",
				Subcommands: []cli.Command{
					{
						Name:  "grant",
						Usage: "SUBCOMMANDS for rights grant",
						Subcommands: []cli.Command{
							{
								Name:         "global",
								Usage:        "Grant a global permission",
								Action:       runtime(cmdRightGrantGlobal),
								BashComplete: cmpl.User,
							},
							{
								Name:         "system",
								Usage:        "Grant a system permission",
								Action:       runtime(cmdRightGrantSystem),
								BashComplete: cmpl.User,
							},
						},
					},
					{
						Name:  "revoke",
						Usage: "SUBCOMMANDS for rights revoke",
						Subcommands: []cli.Command{
							{
								Name:         "global",
								Usage:        "Revoke a global permission",
								Action:       runtime(cmdRightRevokeGlobal),
								BashComplete: cmpl.User,
							},
							{
								Name:         "system",
								Usage:        "Revoke a system permission",
								Action:       runtime(cmdRightRevokeSystem),
								BashComplete: cmpl.User,
							},
						},
					},
				},
			},
		}...,
	)
	return &app
}

func cmdRightGrantGlobal(c *cli.Context) error {
	return cmdRightGrant(c, `global`)
}

func cmdRightGrantSystem(c *cli.Context) error {
	return cmdRightGrant(c, `system`)
}

func cmdRightGrant(c *cli.Context, cat string) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`user`},
		[]string{`user`},
		c.Args().Tail()); err != nil {
		return err
	}

	req := proto.NewGrantRequest()
	req.Grant.RecipientType = `user`
	req.Grant.RecipientId = utl.TryGetUserByUUIDOrName(Client,
		opts[`user`][0])
	req.Grant.PermissionId = utl.TryGetPermissionByUUIDOrName(Client,
		c.Args().First())
	req.Grant.Category = cat

	path := fmt.Sprintf("/grant/%s/%s/%s/", req.Grant.Category,
		req.Grant.RecipientType, req.Grant.RecipientId)
	resp := utl.PostRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdRightRevokeGlobal(c *cli.Context) error {
	return cmdRightRevoke(c, `global`)
}

func cmdRightRevokeSystem(c *cli.Context) error {
	return cmdRightRevoke(c, `system`)
}

func cmdRightRevoke(c *cli.Context, cat string) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`user`},
		[]string{`user`},
		c.Args().Tail()); err != nil {
		return err
	}

	permId := utl.TryGetPermissionByUUIDOrName(Client, c.Args().First())
	userId := utl.TryGetUserByUUIDOrName(Client, opts[`user`][0])
	grantId := utl.TryResolveGrantId(Client, `user`, userId, permId, cat)

	path := fmt.Sprintf("/grant/%s/%s/%s/%s", cat, `user`, userId, grantId)
	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
