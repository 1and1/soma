package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func cmdOnCallAdd(c *cli.Context) {
	url := getApiUrl()
	url.Path = "/oncall"

	validateCliArgumentCount(c, 4)
	// possible keys
	keySlice := []string{"name", "phone"}
	// required keys
	reqSlice := []string{"name", "phone"}
	// received arguments
	argSlice := []string{c.Args().First()}
	argSlice = append(argSlice, c.Args().Tail()...)
	// discard list of optional arguments
	options, _ := parseVariableArguments(keySlice, reqSlice, argSlice)

	// validate phone number as numeric(4,0) --> 1-9999
	phoneNumber, err := strconv.Atoi(options["phone"])
	abortOnError(err, "Syntax error, phone argument not a number")
	if phoneNumber <= 0 || phoneNumber > 9999 {
		abort("Phone number must be 4-digit extension")
	}

	var req somaproto.ProtoRequestOncall
	req.OnCall.Name = options["name"]
	req.OnCall.Number = options["phone"]

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Post(url.String())
	abortOnError(err)
	checkRestyResponse(resp)
}

func cmdOnCallDel(c *cli.Context) {
	url := getApiUrl()
	var (
		id  uuid.UUID
		err error
	)

	switch getCliArgumentCount(c) {
	case 1:
		id, err = uuid.FromString(c.Args().First())
		abortOnError(err, "Syntax error, argument not a uuid")
	case 2:
		validateCliArgument(c, 1, "by-name")
		id = getOncallIdByName(c.Args().Get(1))
	default:
		abort("Syntax error, unexpected argument count")
	}
	url.Path = fmt.Sprintf("/oncall/%s", id.String())

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Delete(url.String())
	abortOnError(err)
	checkRestyResponse(resp)
}

func cmdOnCallRename(c *cli.Context) {
	url := getApiUrl()
	var (
		id     uuid.UUID
		err    error
		oncall string
	)

	switch getCliArgumentCount(c) {
	case 3:
		validateCliArgument(c, 2, "to")
		id, err = uuid.FromString(c.Args().First())
		abortOnError(err, "Syntax error, argument not a uuid")
		oncall = c.Args().Get(2)
	case 4:
		validateCliArgument(c, 1, "by-name")
		validateCliArgument(c, 3, "to")
		id = getOncallIdByName(c.Args().Get(1))
		oncall = c.Args().Get(3)
	default:
		abort("Syntax error, unexpected argument count")
	}
	url.Path = fmt.Sprintf("/oncall/%s", id.String())

	var req somaproto.ProtoRequestOncall
	req.OnCall.Name = oncall

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	abortOnError(err)
	checkRestyResponse(resp)
}

func cmdOnCallUpdate(c *cli.Context) {
	url := getApiUrl()
	var (
		id    uuid.UUID
		err   error
		phone string
	)

	switch getCliArgumentCount(c) {
	case 3:
		validateCliArgument(c, 2, "phone")
		id, err = uuid.FromString(c.Args().First())
		abortOnError(err, "Syntax error, argument not a uuid")
		phone = c.Args().Get(2)
	case 4:
		validateCliArgument(c, 1, "by-name")
		validateCliArgument(c, 3, "phone")
		id = getOncallIdByName(c.Args().Get(1))
		phone = c.Args().Get(3)
	default:
		abort("Syntax error, unexpected argument count")
	}
	num, err := strconv.Atoi(phone)
	abortOnError(err, "Syntax error, argument is not a number")
	if num <= 0 || num > 9999 {
		abort("Phone number must be 4-digit extension")
	}

	url.Path = fmt.Sprintf("/oncall/%s", id.String())
	var req somaproto.ProtoRequestOncall
	req.OnCall.Number = phone

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	abortOnError(err)
	checkRestyResponse(resp)
}

func cmdOnCallList(c *cli.Context) {
	url := getApiUrl()
	url.Path = "/oncall"

	validateCliArgumentCount(c, 0)

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	abortOnError(err)
	checkRestyResponse(resp)
	// TODO print list
}

func cmdOnCallShow(c *cli.Context) {
	url := getApiUrl()
	var (
		id  uuid.UUID
		err error
	)

	switch getCliArgumentCount(c) {
	case 1:
		id, err = uuid.FromString(c.Args().First())
		abortOnError(err, "Syntax error, argument not a uuid")
	case 2:
		validateCliArgument(c, 1, "by-name")
		id = getOncallIdByName(c.Args().Get(1))
	default:
		abort("Syntax error, unexpected argument count")
	}
	url.Path = fmt.Sprintf("/oncall/%s", id.String())

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	abortOnError(err)
	checkRestyResponse(resp)
}

func cmdOnCallMemberAdd(c *cli.Context) {
	url := getApiUrl()
	var (
		id  uuid.UUID
		err error
	)

	switch getCliArgumentCount(c) {
	case 3:
		validateCliArgument(c, 2, "from")
		id, err = uuid.FromString(c.Args().Get(2))
	case 4:
		validateCliArgument(c, 2, "from")
		validateCliArgument(c, 3, "by-name")
		id = getOncallIdByName(c.Args().Get(4))
	default:
		abort("Syntax error, unexpected argument count")
	}
	user := c.Args().Get(0)
	url.Path = fmt.Sprintf("/oncall/%s/members", id.String())

	var req somaproto.ProtoRequestOncall
	var member somaproto.ProtoOncallMember
	member.UserName = user
	reqMembers := []somaproto.ProtoOncallMember{member}
	req.Members = reqMembers

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	abortOnError(err)
	checkRestyResponse(resp)
}

func cmdOnCallMemberDel(c *cli.Context) {
}

func cmdOnCallMemberList(c *cli.Context) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
