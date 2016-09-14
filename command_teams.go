package main

import (
	"fmt"
	"strconv"

	"github.com/1and1/soma/lib/adm"
	"github.com/1and1/soma/lib/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerTeams(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// teams
			{
				Name:  "teams",
				Usage: "SUBCOMMANDS for teams",
				Subcommands: []cli.Command{
					{
						Name:         "add",
						Usage:        "Register a new team",
						Action:       runtime(cmdTeamAdd),
						BashComplete: cmpl.TeamCreate,
					},
					{
						Name:   "remove",
						Usage:  "Delete an existing team",
						Action: runtime(cmdTeamDel),
					},
					{
						Name:         "rename",
						Usage:        "Rename an existing team",
						Action:       runtime(cmdTeamRename),
						BashComplete: cmpl.To,
					},
					{
						Name:   "migrate",
						Usage:  "Migrate users between teams",
						Action: runtime(cmdTeamMigrate),
					},
					{
						Name:   "list",
						Usage:  "List all teams",
						Action: runtime(cmdTeamList),
					},
					{
						Name:   "synclist",
						Usage:  "Export a list of all teams suitable for sync",
						Action: runtime(cmdTeamSync),
					},
					{
						Name:   "show",
						Usage:  "Show information about a team",
						Action: runtime(cmdTeamShow),
					},
					{
						Name:         "update",
						Usage:        "Update team information",
						Action:       runtime(cmdTeamUpdate),
						BashComplete: cmpl.TeamUpdate,
					},
				},
			}, // end teams
		}...,
	)
	return &app
}

func cmdTeamAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 3)
	switch utl.GetCliArgumentCount(c) {
	case 3, 5:
		break // nop
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	multi := []string{}
	required := []string{"ldap"}
	unique := []string{"ldap", "system"}

	opts := utl.ParseVariadicArguments(
		multi,
		unique,
		required,
		c.Args().Tail())

	req := proto.Request{}
	req.Team = &proto.Team{}
	req.Team.Name = c.Args().First()
	req.Team.LdapId = opts["ldap"][0]
	if len(opts["system"]) > 0 {
		bl, err := strconv.ParseBool(opts["system"][0])
		if err != nil {
			utl.Abort("Argument to system parameter must be boolean")
		}
		req.Team.IsSystem = bl
	}

	resp := utl.PostRequestWithBody(Client, req, "/teams/")
	fmt.Println(resp)
	return nil
}

func cmdTeamUpdate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 5)
	multi := []string{}
	unique := []string{`name`, `ldap`, `system`}
	required := []string{`name`, `ldap`}

	opts := utl.ParseVariadicArguments(multi, unique, required, c.Args().Tail())

	teamid := utl.TryGetTeamByUUIDOrName(Client, c.Args().First())
	req := proto.NewTeamRequest()
	req.Team.Name = opts[`name`][0]
	req.Team.LdapId = opts[`ldap`][0]
	if len(opts[`system`]) > 0 {
		req.Team.IsSystem = utl.GetValidatedBool(opts[`system`][0])
	}
	path := fmt.Sprintf("/teams/%s", teamid)
	resp := utl.PutRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdTeamDel(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetTeamByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/teams/%s", id)

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdTeamRename(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	key := []string{"to"}
	opts := utl.ParseVariadicArguments(key, key, key, c.Args().Tail())

	id := utl.TryGetTeamByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/teams/%s", id)

	req := proto.Request{}
	req.Team = &proto.Team{}
	req.Team.Name = opts["to"][0]

	resp := utl.PatchRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdTeamMigrate(c *cli.Context) error {
	// XXX
	utl.Abort("Not implemented")
	return nil
}

func cmdTeamList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)

	resp, err := adm.GetReq(`/teams/`)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func cmdTeamSync(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)

	resp, err := adm.GetReq(`/sync/teams/`)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func cmdTeamShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	id := utl.TryGetTeamByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/teams/%s", id)

	resp, err := adm.GetReq(path)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
