package main

import (
	"fmt"
	"strings"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/internal/help"
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
						Name:         "grant",
						Usage:        "Grant a permission",
						Action:       runtime(cmdRightGrant),
						Description:  help.Text(`RightsGrant`),
						BashComplete: cmpl.Triple_ToOn,
					},
					{
						Name:         "revoke",
						Usage:        "Revoke a permission",
						Action:       runtime(cmdRightRevoke),
						Description:  help.Text(`RightsRevoke`),
						BashComplete: cmpl.Triple_FromOn,
					},
					{
						Name:        `list`,
						Usage:       `List all grants of a permission`,
						Action:      runtime(cmdRightList),
						Description: help.Text(`RightsList`),
					},
					{
						Name:   `show`,
						Usage:  `Show a permission grant for a recipient`,
						Action: runtime(cmdRightShow),
						// BashComplete: cmpl.Triple_For,
						Description: help.Text(`RightsShow`),
					},
				},
			},
		}...,
	)
	return &app
}

func cmdRightGrant(c *cli.Context) error {
	opts := map[string][][2]string{}
	if err := adm.ParseVariadicTriples(
		opts,
		[]string{},
		[]string{`to`, `on`},
		[]string{`to`},
		c.Args().Tail(),
	); err != nil {
		return err
	}
	var (
		err error
	)
	req := proto.NewGrantRequest()

	permissionSlice := strings.Split(c.Args().First(), `::`)
	if len(permissionSlice) != 2 {
		return fmt.Errorf("Invalid split of permission into %s",
			permissionSlice)
	}
	// validate category
	req.Grant.Category = permissionSlice[0]
	if err = adm.ValidateCategory(req.Grant.Category); err != nil {
		return err
	}

	// check optional argument chain
	switch req.Grant.Category {
	case `system`, `global`, `permission`, `operations`:
		fallthrough
	case `global:grant`, `permission:grant`, `operations:grant`:
		if len(opts[`on`]) != 0 {
			return fmt.Errorf("Permissions in category %s are global"+
				" and require no 'on' keyword target.",
				req.Grant.Category)
		}
	case `repository`, `team`, `monitoring`:
		fallthrough
	case `repository:grant`, `team:grant`, `monitoring:grant`:
		if len(opts[`on`]) != 1 {
			return fmt.Errorf("Permissions in category %s require a"+
				" target, specified via 'on' keyword.",
				req.Grant.Category)
		}
	}

	// lookup permissionid
	if err = adm.LookupPermIdRef(
		permissionSlice[1],
		req.Grant.Category,
		&req.Grant.PermissionId,
	); err != nil {
		return err
	}

	// check that the permission is granted to a valid entity
	if err = adm.VerifyPermissionTarget(opts[`to`][0][0]); err != nil {
		return err
	}
	switch opts[`to`][0][0] {
	case `user`:
		req.Grant.RecipientType = `user`
		if req.Grant.RecipientId, err = adm.LookupUserId(
			opts[`to`][0][1]); err != nil {
			return err
		}
	case `admin`:
		return fmt.Errorf(`Admin permissions are not implemented.`)
	case `tool`:
		return fmt.Errorf(`Tool permissions are not implemented.`)
	case `team`:
		return fmt.Errorf(`Team permissions are not implemented.`)
	}

	if len(opts[`on`]) == 1 {
		switch req.Grant.Category {
		case `repository`, `repository:grant`:
			switch opts[`on`][0][0] {
			case `repository`:
				req.Grant.ObjectType = `repository`
				if req.Grant.ObjectId, err = adm.LookupRepoId(
					opts[`on`][0][1],
				); err != nil {
					return err
				}
			case `bucket`:
				req.Grant.ObjectType = `bucket`
				if req.Grant.ObjectId, err = adm.LookupBucketId(
					opts[`on`][0][1],
				); err != nil {
					return err
				}
			default:
				return fmt.Errorf(`Invalid`)
			}
		case `team`, `team:grant`:
			switch opts[`on`][0][0] {
			case `team`:
				req.Grant.ObjectType = `team`
				if req.Grant.ObjectId, err = adm.LookupTeamId(
					opts[`on`][0][1],
				); err != nil {
					return err
				}
			default:
				return fmt.Errorf(`Invalid`)
			}
		case `monitoring`, `monitoring:grant`:
			switch opts[`on`][0][0] {
			case `monitoring`:
				req.Grant.ObjectType = `monitoring`
				if req.Grant.ObjectId, err = adm.LookupMonitoringId(
					opts[`on`][0][1],
				); err != nil {
					return err
				}
			default:
				return fmt.Errorf(`Invalid`)
			}
		}
	}

	path := fmt.Sprintf("/category/%s/permissions/%s/grant/",
		req.Grant.Category, req.Grant.PermissionId)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

func cmdRightRevoke(c *cli.Context) error {
	return fmt.Errorf(`Not implemented - TODO`)
	/*
		opts := map[string][][2]string{}
		if err := adm.ParseVariadicTriples(
			opts,
			[]string{},
			[]string{`from`, `on`},
			[]string{`from`},
			c.Args().Tail(),
		); err != nil {
			return err
		}

		var (
			err                     error
			userId, permId, grantId string
		)
		if err = adm.LookupPermIdRef(c.Args().First(),
			`foobar`, // dummy value for new structs
			&permId); err != nil {
			return err
		}
		if userId, err = adm.LookupUserId(opts[`from`][0][1]); err != nil {
			return err
		}
		if err = adm.LookupGrantIdRef(`user`, userId, permId, `category`,
			&grantId); err != nil {
			return err
		}

		path := fmt.Sprintf("/grant/%s/%s/%s/%s", `category`, `user`, userId,
			grantId)
		return adm.Perform(`delete`, path, `command`, nil, c)
	*/
}

func cmdRightList(c *cli.Context) error {
	return fmt.Errorf(`Not implemented - TODO`)
}

func cmdRightShow(c *cli.Context) error {
	return fmt.Errorf(`Not implemented - TODO`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
