package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
)

func cmdOnCallAdd(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	key := []string{"phone"}

	opts := utl.ParseVariadicArguments(
		key, //allowed
		key, //unique
		key, //required
		c.Args().Tail())
	utl.ValidatePhoneNumber(opts["phone"][0])

	var req somaproto.ProtoRequestOncall
	req.OnCall.Name = c.Args().First()
	req.OnCall.Number = opts["phone"][0]

	resp := utl.PostRequestWithBody(req, "/oncall/")
	fmt.Println(resp)
}

func cmdOnCallDel(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetOncallByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/oncall/%s", id)

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdOnCallRename(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	key := []string{"to"}
	opts := utl.ParseVariadicArguments(key, key, key, c.Args().Tail())

	id := utl.TryGetOncallByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/oncall/%s", id)

	var req somaproto.ProtoRequestOncall
	req.OnCall.Name = opts["to"][0]

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdOnCallUpdate(c *cli.Context) {
	allowed := []string{"phone", "name"}
	unique := []string{"phone", "name"}
	required := []string{}
	utl.ValidateCliMinArgumentCount(c, 3)
	opts := utl.ParseVariadicArguments(allowed, unique, required, c.Args().Tail())

	id := utl.TryGetOncallByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/oncall/%s", id)

	var req somaproto.ProtoRequestOncall
	validUpdate := false
	if len(opts["phone"]) > 0 {
		utl.ValidatePhoneNumber(opts["phone"][0])
		req.OnCall.Number = opts["phone"][0]
		validUpdate = true
	}
	if len(opts["name"]) > 0 {
		req.OnCall.Name = opts["name"][0]
		validUpdate = true
	}
	if !validUpdate {
		utl.Abort("Syntax error, specify name or phone to update")
	}

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdOnCallList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest("/oncall/")
	fmt.Println(resp)
}

func cmdOnCallShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	id := utl.TryGetOncallByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/oncall/%s", id)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdOnCallMemberAdd(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	userId := utl.TryGetUserByUUIDOrName(c.Args().Get(0))
	oncallId := utl.TryGetOncallByUUIDOrName(c.Args().Get(2))

	var req somaproto.ProtoRequestOncall
	var member somaproto.ProtoOncallMember
	member.UserId = userId.String()
	reqMembers := []somaproto.ProtoOncallMember{member}
	req.Members = reqMembers
	path := fmt.Sprintf("/oncall/%s/members", oncallId)

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdOnCallMemberDel(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "from")
	userId := utl.TryGetUserByUUIDOrName(c.Args().Get(0))
	oncallId := utl.TryGetOncallByUUIDOrName(c.Args().Get(2))

	path := fmt.Sprintf("/oncall/%s/members/%s", oncallId, userId.String())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdOnCallMemberList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	oncallId := utl.TryGetOncallByUUIDOrName(c.Args().Get(0))

	path := fmt.Sprintf("/oncall/%s/members/", oncallId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
