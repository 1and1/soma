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
	keySlice := []string{"ldap", "system"}
	reqSlice := []string{"ldap"}
	var err error
	var req somaproto.ProtoRequestTeam

	switch utl.GetCliArgumentCount(c) {
	case 3, 5:
		break // nop
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}

	options, optArgs := parseVariableArguments(keySlice, reqSlice, c.Args().Tail())
	req.Team.TeamName = c.Args().Get(0)
	req.Team.LdapId = options["ldap"]
	if utl.SliceContainsString("system", optArgs) {
		req.Team.System, err = strconv.ParseBool(options["system"])
		utl.AbortOnError(err, "Syntax error, system argument not boolean")
	} else {
		req.Team.System = false
	}

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Post(url.String())
	utl.AbortOnError(err)
	utl.CheckRestyResponse(resp)
}

func cmdTeamDel(c *cli.Context) {
	url := getApiUrl()
	var (
		id  uuid.UUID
		err error
	)

	switch utl.GetCliArgumentCount(c) {
	case 1:
		id, err = uuid.FromString(c.Args().First())
	case 2:
		utl.ValidateCliArgument(c, 1, "by-name")
		id = getTeamIdByName(c.Args().Get(1))
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	url.Path = fmt.Sprintf("/teams/%s", id.String())

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Delete(url.String())
	utl.AbortOnError(err)
	utl.CheckRestyResponse(resp)
}

func cmdTeamRename(c *cli.Context) {
	url := getApiUrl()
	var (
		id       uuid.UUID
		err      error
		teamName string
	)

	switch utl.GetCliArgumentCount(c) {
	case 3:
		utl.ValidateCliArgument(c, 2, "to")
		id, err = uuid.FromString(c.Args().First())
		utl.AbortOnError(err, "Could not parse argument as uuid")
		teamName = c.Args().Get(2)
	case 4:
		utl.ValidateCliArgument(c, 1, "by-name")
		utl.ValidateCliArgument(c, 3, "to")
		id = getTeamIdByName(c.Args().Get(1))
		teamName = c.Args().Get(3)
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	url.Path = fmt.Sprintf("/teams/%s", id.String())

	var req somaproto.ProtoRequestTeam
	req.Team.TeamName = teamName

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	utl.AbortOnError(err)
	utl.CheckRestyResponse(resp)
}

func cmdTeamMigrate(c *cli.Context) {
	// XXX
	utl.Abort("Not implemented")
}

func cmdTeamList(c *cli.Context) {
	url := getApiUrl()
	url.Path = "/teams"

	utl.ValidateCliArgumentCount(c, 0)

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	utl.AbortOnError(err)
	utl.CheckRestyResponse(resp)
	// TODO print list
}

func cmdTeamShow(c *cli.Context) {
	url := getApiUrl()
	var (
		id  uuid.UUID
		err error
	)

	switch utl.GetCliArgumentCount(c) {
	case 1:
		id, err = uuid.FromString(c.Args().First())
		utl.AbortOnError(err, "Could not parse argument as uuid")
	case 2:
		utl.ValidateCliArgument(c, 1, "by-name")
		id = getTeamIdByName(c.Args().Get(1))
	default:
		utl.AbortOnError(err)
	}
	url.Path = fmt.Sprintf("/teams/%s", id.String())

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	utl.AbortOnError(err)
	utl.CheckRestyResponse(resp)
	// TODO print record
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
