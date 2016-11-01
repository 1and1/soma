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
	key := []string{"phone"}

	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		key, //allowed
		key, //unique
		key, //required
		c.Args().Tail()); err != nil {
		return err
	}
	utl.ValidatePhoneNumber(opts["phone"][0])

	req := proto.Request{}
	req.Oncall = &proto.Oncall{}
	req.Oncall.Name = c.Args().First()
	req.Oncall.Number = opts["phone"][0]

	resp := utl.PostRequestWithBody(Client, req, "/oncall/")
	fmt.Println(resp)
	return nil
}

func cmdOnCallDel(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id := utl.TryGetOncallByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/oncall/%s", id)

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdOnCallRename(c *cli.Context) error {
	key := []string{"to"}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(opts, key, key, key,
		c.Args().Tail()); err != nil {
		return err
	}

	id := utl.TryGetOncallByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/oncall/%s", id)

	req := proto.Request{}
	req.Oncall = &proto.Oncall{}
	req.Oncall.Name = opts["to"][0]

	resp := utl.PatchRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdOnCallUpdate(c *cli.Context) error {
	allowed := []string{"phone", "name"}
	unique := []string{"phone", "name"}
	required := []string{}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(opts, allowed, unique, required,
		c.Args().Tail()); err != nil {
		return err
	}

	id := utl.TryGetOncallByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/oncall/%s", id)

	req := proto.Request{}
	req.Oncall = &proto.Oncall{}
	validUpdate := false
	if len(opts["phone"]) > 0 {
		utl.ValidatePhoneNumber(opts["phone"][0])
		req.Oncall.Number = opts["phone"][0]
		validUpdate = true
	}
	if len(opts["name"]) > 0 {
		req.Oncall.Name = opts["name"][0]
		validUpdate = true
	}
	if !validUpdate {
		adm.Abort("Syntax error, specify name or phone to update")
	}

	resp := utl.PatchRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdOnCallList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	resp := utl.GetRequest(Client, "/oncall/")
	fmt.Println(resp)
	return nil
}

func cmdOnCallShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	id := utl.TryGetOncallByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/oncall/%s", id)

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
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
	userId := utl.TryGetUserByUUIDOrName(Client, c.Args().First())
	oncallId := utl.TryGetOncallByUUIDOrName(Client, opts[`to`][0])

	req := proto.Request{}
	req.Oncall = &proto.Oncall{}
	member := proto.OncallMember{}
	member.UserId = userId
	reqMembers := []proto.OncallMember{member}
	req.Oncall.Members = &reqMembers
	path := fmt.Sprintf("/oncall/%s/members", oncallId)

	resp := utl.PatchRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdOnCallMemberDel(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`from`},
		[]string{`from`},
		c.Args().Tail()); err != nil {
		return err
	}
	userId := utl.TryGetUserByUUIDOrName(Client, c.Args().First())
	oncallId := utl.TryGetOncallByUUIDOrName(Client, opts[`from`][0])

	path := fmt.Sprintf("/oncall/%s/members/%s", oncallId, userId)

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdOnCallMemberList(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	oncallId := utl.TryGetOncallByUUIDOrName(Client, c.Args().Get(0))

	path := fmt.Sprintf("/oncall/%s/members/", oncallId)

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
