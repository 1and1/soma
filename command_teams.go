package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
)

func registerTeams(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// teams
			{
				Name:   "teams",
				Usage:  "SUBCOMMANDS for teams",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "add",
						Usage:  "Register a new team",
						Action: cmdTeamAdd,
					},
					{
						Name:   "remove",
						Usage:  "Delete an existing team",
						Action: cmdTeamDel,
					},
					{
						Name:   "rename",
						Usage:  "Rename an existing team",
						Action: cmdTeamRename,
					},
					{
						Name:   "migrate",
						Usage:  "Migrate users between teams",
						Action: cmdTeamMigrate,
					},
					{
						Name:   "list",
						Usage:  "List all teams",
						Action: cmdTeamList,
					},
					{
						Name:   "show",
						Usage:  "Show information about a team",
						Action: cmdTeamShow,
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
	allowed := []string{"ldap", "system"}
	required := []string{"ldap"}
	unique := []string{"ldap", "system"}

	opts := utl.ParseVariadicArguments(
		allowed,
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

	resp := utl.PostRequestWithBody(req, "/teams/")
	fmt.Println(resp)
	return nil
}

func cmdTeamDel(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetTeamByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/teams/%s", id)

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
	return nil
}

func cmdTeamRename(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	key := []string{"to"}
	opts := utl.ParseVariadicArguments(key, key, key, c.Args().Tail())

	id := utl.TryGetTeamByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/teams/%s", id)

	req := proto.Request{}
	req.Team = &proto.Team{}
	req.Team.Name = opts["to"][0]

	resp := utl.PatchRequestWithBody(req, path)
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

	resp := utl.GetRequest("/teams/")
	fmt.Println(resp)
	return nil
}

func cmdTeamShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	id := utl.TryGetTeamByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/teams/%s", id)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
