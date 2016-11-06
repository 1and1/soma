package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerRepository(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// repository
			{
				Name:  "repository",
				Usage: "SUBCOMMANDS for repository",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new repository",
						Action:       runtime(cmdRepositoryCreate),
						BashComplete: cmpl.Team,
					},
					{
						Name:   "delete",
						Usage:  "Mark an existing repository as deleted",
						Action: runtime(cmdRepositoryDelete),
					},
					{
						Name:   "restore",
						Usage:  "Restore a repository marked as deleted",
						Action: runtime(cmdRepositoryRestore),
					},
					{
						Name:   "purge",
						Usage:  "Remove an unreferenced deleted repository",
						Action: runtime(cmdRepositoryPurge),
					},
					{
						Name:   "clear",
						Usage:  "Clear all check instances for this repository",
						Action: runtime(cmdRepositoryClear),
					},
					{
						Name:         "rename",
						Usage:        "Rename an existing repository",
						Action:       runtime(cmdRepositoryRename),
						BashComplete: cmpl.To,
					},
					{
						Name:         "repossess",
						Usage:        "Change the owner of a repository",
						Action:       runtime(cmdRepositoryRepossess),
						BashComplete: cmpl.To,
					},
					{
						Name:   "activate",
						Usage:  "Activate a cloned repository",
						Action: runtime(cmdRepositoryActivate),
					},
					{
						Name:   "list",
						Usage:  "List all existing repositories",
						Action: runtime(cmdRepositoryList),
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific repository",
						Action: runtime(cmdRepositoryShow),
					},
					{
						Name:   `tree`,
						Usage:  `Display the repository as tree`,
						Action: runtime(cmdRepositoryTree),
					},
					{
						Name:   `instances`,
						Usage:  `List check instances for a repository`,
						Action: runtime(cmdRepositoryInstance),
					},
					{
						Name:  "property",
						Usage: "SUBCOMMANDS for properties",
						Subcommands: []cli.Command{
							{
								Name:  "add",
								Usage: "SUBCOMMANDS for property add",
								Subcommands: []cli.Command{
									{
										Name:         "system",
										Usage:        "Add a system property to a repository",
										Action:       runtime(cmdRepositorySystemPropertyAdd),
										BashComplete: cmpl.PropertyAddValue,
									},
									{
										Name:         "service",
										Usage:        "Add a service property to a repository",
										Action:       runtime(cmdRepositoryServicePropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         "oncall",
										Usage:        "Add an oncall property to a repository",
										Action:       runtime(cmdRepositoryOncallPropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         "custom",
										Usage:        "Add a custom property to a repository",
										Action:       runtime(cmdRepositoryCustomPropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
								},
							},
							{
								Name:  `delete`,
								Usage: `SUBCOMMANDS for property delete`,
								Subcommands: []cli.Command{
									{
										Name:         `system`,
										Usage:        `Delete a system property from a repository`,
										Action:       runtime(cmdRepositorySystemPropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `service`,
										Usage:        `Delete a service property from a repository`,
										Action:       runtime(cmdRepositoryServicePropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `oncall`,
										Usage:        `Delete an oncall property from a repository`,
										Action:       runtime(cmdRepositoryOncallPropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `custom`,
										Usage:        `Delete a custom property from a repository`,
										Action:       runtime(cmdRepositoryCustomPropertyDelete),
										BashComplete: cmpl.FromView,
									},
								},
							},
						},
					},
				},
			}, // end repository
		}...,
	)
	return &app
}

func cmdRepositoryCreate(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`team`},
		[]string{`team`},
		adm.AllArguments(c)); err != nil {
		return err
	}

	teamId, err := adm.LookupTeamId(opts[`team`][0])
	if err != nil {
		return err
	}

	var req proto.Request
	req.Repository = &proto.Repository{}
	req.Repository.Name = c.Args().Get(0)
	req.Repository.TeamId = teamId

	if err := adm.ValidateRuneCountRange(req.Repository.Name,
		4, 128); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/repository/`, `command`, req, c)
}

func cmdRepositoryDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupRepoId(c.Args().First())
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/repository/%s", id)

	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdRepositoryRestore(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupRepoId(c.Args().First())
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/repository/%s", id)

	req := proto.Request{
		Flags: &proto.Flags{
			Restore: true,
		},
	}

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdRepositoryPurge(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupRepoId(c.Args().First())
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/repository/%s", id)

	req := proto.Request{
		Flags: &proto.Flags{
			Purge: true,
		},
	}

	return adm.Perform(`deletebody`, path, `command`, req, c)
}

func cmdRepositoryClear(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupRepoId(c.Args().First())
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/repository/%s", id)

	req := proto.Request{
		Flags: &proto.Flags{
			Clear: true,
		},
	}

	return adm.Perform(`putbody`, path, `command`, req, c)
}

func cmdRepositoryRename(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`to`},
		[]string{`to`},
		c.Args().Tail()); err != nil {
		return err
	}
	id, err := adm.LookupRepoId(c.Args().First())
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/repository/%s", id)

	var req proto.Request
	req.Repository = &proto.Repository{}
	req.Repository.Name = opts[`to`][0]

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdRepositoryRepossess(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`to`},
		[]string{`to`},
		c.Args().Tail()); err != nil {
		return err
	}
	id, err := adm.LookupRepoId(c.Args().First())
	if err != nil {
		return err
	}
	teamId, err := adm.LookupTeamId(opts[`team`][0])
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/repository/%s", id)

	var req proto.Request
	req.Repository = &proto.Repository{}
	req.Repository.TeamId = teamId

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdRepositoryClone(c *cli.Context) error {
	return fmt.Errorf(`Not implemented`)
}

func cmdRepositoryActivate(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupRepoId(c.Args().First())
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/repository/%s", id)

	req := proto.Request{
		Flags: &proto.Flags{
			Activate: true,
		},
	}

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdRepositoryWipe(c *cli.Context) error {
	return fmt.Errorf(`Not implemented`)
}

func cmdRepositoryList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/repository/`, `list`, nil, c)
}

func cmdRepositoryShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupRepoId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s", id)
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdRepositoryInstance(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupRepoId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/instances/", id)
	return adm.Perform(`get`, path, `list`, nil, c)
}

func cmdRepositoryTree(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupRepoId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/tree/tree", id)
	return adm.Perform(`get`, path, `tree`, nil, c)
}

func cmdRepositorySystemPropertyAdd(c *cli.Context) error {
	return cmdRepositoryPropertyAdd(c, `system`)
}

func cmdRepositoryServicePropertyAdd(c *cli.Context) error {
	return cmdRepositoryPropertyAdd(c, `service`)
}

func cmdRepositoryOncallPropertyAdd(c *cli.Context) error {
	return cmdRepositoryPropertyAdd(c, `oncall`)
}

func cmdRepositoryCustomPropertyAdd(c *cli.Context) error {
	return cmdRepositoryPropertyAdd(c, `custom`)
}

func cmdRepositoryPropertyAdd(c *cli.Context, pType string) error {
	return cmdPropertyAdd(c, pType, `repository`)
}

func cmdRepositorySystemPropertyDelete(c *cli.Context) error {
	return cmdRepositoryPropertyDelete(c, `system`)
}

func cmdRepositoryServicePropertyDelete(c *cli.Context) error {
	return cmdRepositoryPropertyDelete(c, `service`)
}

func cmdRepositoryOncallPropertyDelete(c *cli.Context) error {
	return cmdRepositoryPropertyDelete(c, `oncall`)
}

func cmdRepositoryCustomPropertyDelete(c *cli.Context) error {
	return cmdRepositoryPropertyDelete(c, `custom`)
}

func cmdRepositoryPropertyDelete(c *cli.Context, pType string) error {
	unique := []string{`from`, `view`}
	required := []string{`from`, `view`}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}
	repositoryId, err := adm.LookupRepoId(opts[`from`][0])
	if err != nil {
		return err
	}

	if pType == `system` {
		if err := adm.ValidateSystemProperty(
			c.Args().First()); err != nil {
			return err
		}
	}
	var sourceId string
	if err := adm.FindRepoPropSrcId(pType, c.Args().First(),
		opts[`view`][0], repositoryId, &sourceId); err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/property/%s/%s",
		repositoryId, pType, sourceId)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
