package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func cmdTeamAdd(c *cli.Context) {
	url := getApiUrl()
	url.Path = "/teams"
	keys := [...]string{"ldap", "system"}
	keySlice := make([]string, 0)
	var err error
	var req somaproto.ProtoRequestTeam

	switch getCliArgumentCount(c) {
	case 3:
		keySlice = keys[1:1]
	case 5:
		keySlice = keys[1:2]
	default:
		abort("Syntax error, unexpected argument count")
	}

	options := *parseVariableArguments(keySlice, c.Args().Tail())
	req.Team.TeamName = c.Args().Get(0)
	req.Team.LdapId = options["ldap"]
	if options["system"] == "" {
		req.Team.System = false
	} else {
		req.Team.System, err = strconv.ParseBool(options["system"])
		abortOnError(err, "Syntax error, system argument not boolean")
	}

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Post(url.String())
	abortOnError(err)
	checkRestyResponse(resp)
}

func cmdTeamDel(c *cli.Context) {
	url := getApiUrl()
	var (
		id  uuid.UUID
		err error
	)

	switch getCliArgumentCount(c) {
	case 1:
		id, err = uuid.FromString(c.Args().First())
	case 2:
		validateCliArgument(c, 1, "by-name")
		id = getTeamIdByName(c.Args().Get(1))
	default:
		abort("Syntax error, unexpected argument count")
	}
	url.Path = fmt.Sprintf("/teams/%s", id.String())

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Delete(url.String())
	abortOnError(err)
	checkRestyResponse(resp)
}

func cmdTeamRename(c *cli.Context) {
	url := getApiUrl()
	var (
		id       uuid.UUID
		err      error
		teamName string
	)

	switch getCliArgumentCount(c) {
	case 3:
		validateCliArgument(c, 2, "to")
		id, err = uuid.FromString(c.Args().First())
		abortOnError(err, "Could not parse argument as uuid")
		teamName = c.Args().Get(1)
	case 4:
		validateCliArgument(c, 1, "by-name")
		validateCliArgument(c, 3, "to")
		id = getTeamIdByName(c.Args().Get(1))
		teamName = c.Args().Get(3)
	default:
		abort("Syntax error, unexpected argument count")
	}
	url.Path = fmt.Sprintf("/teams/%s", id.String())

	var req somaproto.ProtoRequestTeam
	req.Team.TeamName = teamName

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	abortOnError(err)
	checkRestyResponse(resp)
}

func cmdTeamMigrate(c *cli.Context) {
	// XXX
	abort("Not implemented")
}

func cmdTeamList(c *cli.Context) {
	url := getApiUrl()
	url.Path = "/teams"

	validateCliArgumentCount(0)

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	abortOnError(err)
	checkRestyResponse(resp)
	// TODO print list
}

func cmdTeamShow(c *cli.Context) {
	url := getApiUrl()
	var (
		id       uuid.UUID
		err      error
		teamName string
	)

	switch getCliArgumentCount(c) {
	case 1:
		id, err = uuid.FromString(c.Args().First())
		abortOnError(err, "Could not parse argument as uuid")
	case 2:
		validateCliArgument(c, 1, "by-name")
		id = getTeamIdByName(c.Args().Get(1))
	default:
		abortOnError(err)
	}
	url.Path = fmt.Sprintf("/teams/%s", id.String())

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	abortOnError(err)
	checkRestyResponse(resp)
	// TODO print record
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
