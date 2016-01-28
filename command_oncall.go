package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/satori/go.uuid"
)

func cmdOnCallAdd(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 4)
	// possible keys
	keySlice := []string{"name", "phone"}
	// required keys
	reqSlice := []string{"name", "phone"}
	// received arguments
	argSlice := []string{c.Args().First()}
	argSlice = append(argSlice, c.Args().Tail()...)
	// discard list of optional arguments
	options, _ := utl.ParseVariableArguments(keySlice, reqSlice, argSlice)

	// validate phone number as numeric(4,0) --> 1-9999
	phoneNumber, err := strconv.Atoi(options["phone"])
	utl.AbortOnError(err, "Syntax error, phone argument not a number")
	if phoneNumber <= 0 || phoneNumber > 9999 {
		utl.Abort("Phone number must be 4-digit extension")
	}

	var req somaproto.ProtoRequestOncall
	req.OnCall.Name = options["name"]
	req.OnCall.Number = options["phone"]

	resp := utl.PostRequestWithBody(req, "/oncall/")
	fmt.Println(resp)
}

func cmdOnCallDel(c *cli.Context) {
	var (
		id string
	)

	switch utl.GetCliArgumentCount(c) {
	case 1:
		uid, err := uuid.FromString(c.Args().First())
		utl.AbortOnError(err, "Syntax error, argument not a uuid")
		id = uid.String()
	case 2:
		utl.ValidateCliArgument(c, 1, "by-name")
		id = utl.GetOncallIdByName(c.Args().Get(1))
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	path := fmt.Sprintf("/oncall/%s", id)

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdOnCallRename(c *cli.Context) {
	var (
		id     string
		oncall string
	)

	switch utl.GetCliArgumentCount(c) {
	case 3:
		utl.ValidateCliArgument(c, 2, "to")
		uid, err := uuid.FromString(c.Args().First())
		utl.AbortOnError(err, "Syntax error, argument not a uuid")
		id = uid.String()
		oncall = c.Args().Get(2)
	case 4:
		utl.ValidateCliArgument(c, 1, "by-name")
		utl.ValidateCliArgument(c, 3, "to")
		id = utl.GetOncallIdByName(c.Args().Get(1))
		oncall = c.Args().Get(3)
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	path := fmt.Sprintf("/oncall/%s", id)

	var req somaproto.ProtoRequestOncall
	req.OnCall.Name = oncall

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdOnCallUpdate(c *cli.Context) {
	var (
		id    string
		err   error
		phone string
	)

	switch utl.GetCliArgumentCount(c) {
	case 3:
		utl.ValidateCliArgument(c, 2, "phone")
		uid, err := uuid.FromString(c.Args().First())
		utl.AbortOnError(err, "Syntax error, argument not a uuid")
		id = uid.String()
		phone = c.Args().Get(2)
	case 4:
		utl.ValidateCliArgument(c, 1, "by-name")
		utl.ValidateCliArgument(c, 3, "phone")
		id = utl.GetOncallIdByName(c.Args().Get(1))
		phone = c.Args().Get(3)
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	num, err := strconv.Atoi(phone)
	utl.AbortOnError(err, "Syntax error, argument is not a number")
	if num <= 0 || num > 9999 {
		utl.Abort("Phone number must be 4-digit extension")
	}

	path := fmt.Sprintf("/oncall/%s", id)
	var req somaproto.ProtoRequestOncall
	req.OnCall.Number = phone

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdOnCallList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest("/oncall/")
	fmt.Println(resp)
}

func cmdOnCallShow(c *cli.Context) {
	var (
		id string
	)

	switch utl.GetCliArgumentCount(c) {
	case 1:
		uid, err := uuid.FromString(c.Args().First())
		utl.AbortOnError(err, "Syntax error, argument not a uuid")
		id = uid.String()
	case 2:
		utl.ValidateCliArgument(c, 1, "by-name")
		id = utl.GetOncallIdByName(c.Args().Get(1))
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
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
