package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
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
	multi := []string{}
	required := []string{"ldap"}
	unique := []string{"ldap", "system"}

	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		multi,
		unique,
		required,
		c.Args().Tail()); err != nil {
		return err
	}

	req := proto.NewTeamRequest()
	req.Team.Name = c.Args().First()
	req.Team.LdapId = opts["ldap"][0]
	if len(opts["system"]) > 0 {
		if err := adm.ValidateBool(opts["system"][0],
			&req.Team.IsSystem); err != nil {
			return err
		}
	}

	return adm.Perform(`postbody`, `/teams/`, `command`, req, c)
}

func cmdTeamUpdate(c *cli.Context) error {
	multi := []string{}
	unique := []string{`name`, `ldap`, `system`}
	required := []string{`name`, `ldap`}

	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(opts, multi, unique, required,
		c.Args().Tail()); err != nil {
		return err
	}

	teamid, err := adm.LookupTeamId(c.Args().First())
	if err != nil {
		return err
	}
	req := proto.NewTeamRequest()
	req.Team.Name = opts[`name`][0]
	req.Team.LdapId = opts[`ldap`][0]
	if len(opts[`system`]) > 0 {
		if err := adm.ValidateBool(opts["system"][0],
			&req.Team.IsSystem); err != nil {
			return fmt.Errorf("Argument to system parameter must" +
				" be boolean")
		}
	}
	path := fmt.Sprintf("/teams/%s", teamid)
	return adm.Perform(`putbody`, path, `command`, req, c)
}

func cmdTeamDel(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupTeamId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/teams/%s", id)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdTeamRename(c *cli.Context) error {
	key := []string{"to"}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(opts, key, key, key,
		c.Args().Tail()); err != nil {
		return err
	}

	id, err := adm.LookupTeamId(c.Args().First())
	if err != nil {
		return err
	}

	req := proto.NewTeamRequest()
	req.Team.Name = opts["to"][0]

	path := fmt.Sprintf("/teams/%s", id)
	return adm.Perform(`patchbody`, path, `command`, nil, c)
}

func cmdTeamMigrate(c *cli.Context) error {
	return fmt.Errorf(`Not implemented.`)
}

func cmdTeamList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/teams/`, `list`, nil, c)
}

func cmdTeamSync(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/sync/teams/`, `list`, nil, c)
}

func cmdTeamShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	id, err := adm.LookupTeamId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/teams/%s", id)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
