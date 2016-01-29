package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
)

func cmdTeamAdd(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 3)
	switch utl.GetCliArgumentCount(c) {
	case 3, 5:
		break // nop
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	allowed := []string{"ldap", "system"}
	required := []string{"ldap"}
	unique := []string{"ldap", "system"}

	opts := utl.ParseVariadicArguments(
		allowed,
		unique,
		required,
		c.Args().Tail())

	req := somaproto.ProtoRequestTeam{}
	req.Team = &somaproto.ProtoTeam{}
	req.Team.Name = c.Args().First()
	req.Team.Ldap = opts["ldap"][0]
	if len(opts["system"]) > 0 {
		bl, err := strconv.ParseBool(opts["system"][0])
		if err != nil {
			utl.Abort("Argument to system parameter must be boolean")
		}
		req.Team.System = bl
	}

	resp := utl.PostRequestWithBody(req, "/teams/")
	fmt.Println(resp)
}

func cmdTeamDel(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetTeamByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/teams/%s", id)

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdTeamRename(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	key := []string{"to"}
	opts := utl.ParseVariadicArguments(key, key, key, c.Args().Tail())

	id := utl.TryGetTeamByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/teams/%s", id)

	req := somaproto.ProtoRequestTeam{}
	req.Team = &somaproto.ProtoTeam{}
	req.Team.Name = opts["to"][0]

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdTeamMigrate(c *cli.Context) {
	// XXX
	utl.Abort("Not implemented")
}

func cmdTeamList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest("/teams/")
	fmt.Println(resp)
}

func cmdTeamShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	id := utl.TryGetTeamByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/teams/%s", id)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
