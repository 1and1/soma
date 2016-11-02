package main

import (
	"fmt"
	"strconv"

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
		bl, err := strconv.ParseBool(opts["system"][0])
		if err != nil {
			return fmt.Errorf("Argument to system parameter must be boolean")
		}
		req.Team.IsSystem = bl
	}

	if resp, err := adm.PostReqBody(req, `/teams/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
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
		var err error
		req.Team.IsSystem, err = strconv.ParseBool(opts[`system`][0])
		if err != nil {
			return fmt.Errorf("Argument to system parameter must be boolean")
		}
	}
	path := fmt.Sprintf("/teams/%s", teamid)
	if resp, err := adm.PutReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
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
	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
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
	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
}

func cmdTeamMigrate(c *cli.Context) error {
	return fmt.Errorf(`Not implemented.`)
}

func cmdTeamList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	if resp, err := adm.GetReq(`/teams/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdTeamSync(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	if resp, err := adm.GetReq(`/sync/teams/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `command`)
	}
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
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `show`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
