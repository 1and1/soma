package main

import (
	"fmt"

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
						Name:   "add",
						Usage:  "Register a new oncall duty team",
						Action: runtime(cmdOnCallAdd),
					},
					{
						Name:   "remove",
						Usage:  "Delete an existing oncall duty team",
						Action: runtime(cmdOnCallDel),
					},
					{
						Name:   "rename",
						Usage:  "Rename an existing oncall duty team",
						Action: runtime(cmdOnCallRename),
					},
					{
						Name:   "update",
						Usage:  "Update phone number of an existing oncall duty team",
						Action: runtime(cmdOnCallUpdate),
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
								Name:   "add",
								Usage:  "Add a user to an oncall duty team",
								Action: runtime(cmdOnCallMemberAdd),
							},
							{
								Name:   "remove",
								Usage:  "Remove a member from an oncall duty team",
								Action: runtime(cmdOnCallMemberDel),
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
	utl.ValidateCliArgumentCount(c, 3)
	key := []string{"phone"}

	opts := utl.ParseVariadicArguments(
		key, //allowed
		key, //unique
		key, //required
		c.Args().Tail())
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
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetOncallByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/oncall/%s", id)

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdOnCallRename(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	key := []string{"to"}
	opts := utl.ParseVariadicArguments(key, key, key, c.Args().Tail())

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
	utl.ValidateCliMinArgumentCount(c, 3)
	opts := utl.ParseVariadicArguments(allowed, unique, required, c.Args().Tail())

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
		utl.Abort("Syntax error, specify name or phone to update")
	}

	resp := utl.PatchRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdOnCallList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest(Client, "/oncall/")
	fmt.Println(resp)
	return nil
}

func cmdOnCallShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	id := utl.TryGetOncallByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/oncall/%s", id)

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdOnCallMemberAdd(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	userId := utl.TryGetUserByUUIDOrName(Client, c.Args().Get(0))
	oncallId := utl.TryGetOncallByUUIDOrName(Client, c.Args().Get(2))

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
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "from")
	userId := utl.TryGetUserByUUIDOrName(Client, c.Args().Get(0))
	oncallId := utl.TryGetOncallByUUIDOrName(Client, c.Args().Get(2))

	path := fmt.Sprintf("/oncall/%s/members/%s", oncallId, userId)

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdOnCallMemberList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	oncallId := utl.TryGetOncallByUUIDOrName(Client, c.Args().Get(0))

	path := fmt.Sprintf("/oncall/%s/members/", oncallId)

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
