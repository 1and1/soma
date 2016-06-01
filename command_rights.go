package main

import (
	"fmt"

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
								Name:   "global",
								Usage:  "Grant a global permission",
								Action: runtime(cmdRightGrantGlobal),
							},
							{
								Name:   "system",
								Usage:  "Grant a system permission",
								Action: runtime(cmdRightGrantSystem),
							},
						},
					},
					{
						Name:  "revoke",
						Usage: "SUBCOMMANDS for rights revoke",
						Subcommands: []cli.Command{
							{
								Name:   "global",
								Usage:  "Revoke a global permission",
								Action: runtime(cmdRightRevokeGlobal),
							},
							{
								Name:   "system",
								Usage:  "Revoke a system permission",
								Action: runtime(cmdRightRevokeSystem),
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
	utl.ValidateCliArgumentCount(c, 3)
	opts := utl.ParseVariadicArguments(
		[]string{}, []string{`user`},
		[]string{`user`}, c.Args().Tail())

	req := proto.NewGrantRequest()
	req.Grant.RecipientType = `user`
	req.Grant.RecipientId = utl.TryGetUserByUUIDOrName(Client, opts[`user`][0])
	req.Grant.PermissionId = utl.TryGetPermissionByUUIDOrName(Client, c.Args().First())
	req.Grant.Category = cat

	path := fmt.Sprintf("/grant/%s/%s/%s/", req.Grant.Category, req.Grant.RecipientType, req.Grant.RecipientId)
	resp := utl.PostRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdRightRevokeGlobal(c *cli.Context) error {
	return nil
}

func cmdRightRevokeSystem(c *cli.Context) error {
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
