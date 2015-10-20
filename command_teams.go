package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
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
		Slog.Fatal("Syntax error, unexpected argument count")
	}

	options := *parseVariableArguments(keySlice, c.Args().Tail())
	req.Team.TeamName = c.Args().Get(0)
	req.Team.LdapId = options["ldap"]
	if options["system"] == "" {
		req.Team.System = false
	} else {
		req.Team.System, err = strconv.ParseBool(options["system"])
		if err != nil {
			Slog.Fatal("Syntax error, system argument not boolean")
		}
	}

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Post(url.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		Slog.Fatal(err)
	}
	checkRestyResponse(resp)
}

func cmdTeamDel(c *cli.Context) {
}

func cmdTeamRename(c *cli.Context) {
}

func cmdTeamMigrate(c *cli.Context) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
