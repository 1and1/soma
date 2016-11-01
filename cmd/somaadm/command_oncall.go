package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerOncall(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// oncall
			{
				Name:  "oncall",
				Usage: "SUBCOMMANDS for oncall duty teams",
				Subcommands: []cli.Command{
					{
						Name:         "add",
						Usage:        "Register a new oncall duty team",
						Action:       runtime(cmdOnCallAdd),
						BashComplete: cmpl.OnCallAdd,
					},
					{
						Name:   "remove",
						Usage:  "Delete an existing oncall duty team",
						Action: runtime(cmdOnCallDel),
					},
					{
						Name:         "rename",
						Usage:        "Rename an existing oncall duty team",
						Action:       runtime(cmdOnCallRename),
						BashComplete: cmpl.To,
					},
					{
						Name:         "update",
						Usage:        "Update phone number of an existing oncall duty team",
						Action:       runtime(cmdOnCallUpdate),
						BashComplete: cmpl.OnCallUpdate,
					},
					{
						Name:   "list",
						Usage:  "List all registered oncall duty teams",
						Action: runtime(cmdOnCallList),
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific oncall duty team",
						Action: runtime(cmdOnCallShow),
					},
					{
						Name:  "member",
						Usage: "SUBCOMMANDS to manipulate oncall duty members",
						Subcommands: []cli.Command{
							{
								Name:         "add",
								Usage:        "Add a user to an oncall duty team",
								Action:       runtime(cmdOnCallMemberAdd),
								BashComplete: cmpl.To,
							},
							{
								Name:         "remove",
								Usage:        "Remove a member from an oncall duty team",
								Action:       runtime(cmdOnCallMemberDel),
								BashComplete: cmpl.From,
							},
							{
								Name:   "list",
								Usage:  "List the users of an oncall duty team",
								Action: runtime(cmdOnCallMemberList),
							},
						},
					},
				},
			}, // end oncall
		}...,
	)
	return &app
}

func cmdOnCallAdd(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`phone`},
		[]string{`phone`},
		c.Args().Tail()); err != nil {
		return err
	}
	if err := adm.ValidateOncallNumber(opts["phone"][0]); err != nil {
		return err
	}

	req := proto.NewOncallRequest()
	req.Oncall.Name = c.Args().First()
	req.Oncall.Number = opts["phone"][0]

	if resp, err := adm.PostReqBody(req, `/oncall/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdOnCallDel(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	id, err := adm.LookupOncallId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/oncall/%s", id)
	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdOnCallRename(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`to`},
		[]string{`to`},
		c.Args().Tail(),
	); err != nil {
		return err
	}

	id, err := adm.LookupOncallId(c.Args().First())
	if err != nil {
		return err
	}

	req := proto.NewOncallRequest()
	req.Oncall.Name = opts["to"][0]

	path := fmt.Sprintf("/oncall/%s", id)
	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdOnCallUpdate(c *cli.Context) error {
	unique := []string{"phone", "name"}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		[]string{},
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewOncallRequest()
	validUpdate := false
	if len(opts["phone"]) > 0 {
		if err := adm.ValidateOncallNumber(opts["phone"][0]); err != nil {
			return err
		}
		req.Oncall.Number = opts["phone"][0]
		validUpdate = true
	}
	if len(opts["name"]) > 0 {
		req.Oncall.Name = opts["name"][0]
		validUpdate = true
	}
	if !validUpdate {
		return fmt.Errorf("Syntax error, specify name or phone to update")
	}

	id, err := adm.LookupOncallId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/oncall/%s", id)
	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdOnCallList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	if resp, err := adm.GetReq(`/oncall/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdOnCallShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	id, err := adm.LookupOncallId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/oncall/%s", id)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `show`)
	}
}

func cmdOnCallMemberAdd(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`to`},
		[]string{`to`},
		c.Args().Tail()); err != nil {
		return err
	}

	var (
		err              error
		userId, oncallId string
	)
	if userId, err = adm.LookupUserId(c.Args().First()); err != nil {
		return err
	}
	if oncallId, err = adm.LookupOncallId(opts[`to`][0]); err != nil {
		return err
	}

	req := proto.NewOncallRequest()
	req.Oncall.Members = &[]proto.OncallMember{
		proto.OncallMember{UserId: userId},
	}

	path := fmt.Sprintf("/oncall/%s/members", oncallId)
	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdOnCallMemberDel(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`from`},
		[]string{`from`},
		c.Args().Tail(),
	); err != nil {
		return err
	}

	var (
		err              error
		userId, oncallId string
	)
	if userId, err = adm.LookupUserId(c.Args().First()); err != nil {
		return err
	}
	if oncallId, err = adm.LookupOncallId(opts[`from`][0]); err != nil {
		return err
	}

	path := fmt.Sprintf("/oncall/%s/members/%s", oncallId, userId)
	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdOnCallMemberList(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	oncallId, err := adm.LookupOncallId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/oncall/%s/members/", oncallId)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
