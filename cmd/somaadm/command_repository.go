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
					/*
						{
							Name:   "clone",
							Usage:  "Create a clone of an existing repository",
							Action: cmdRepositoryClone,
						},
					*/
					{
						Name:   "activate",
						Usage:  "Activate a cloned repository",
						Action: runtime(cmdRepositoryActivate),
					},
					/*
						{
							Name:   "wipe",
							Usage:  "Clear all repository contents",
							Action: cmdRepositoryWipe,
						},
					*/
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

	teamId := utl.TryGetTeamByUUIDOrName(Client, opts[`team`][0])

	var req proto.Request
	req.Repository = &proto.Repository{}
	req.Repository.Name = c.Args().Get(0)
	req.Repository.TeamId = teamId

	utl.ValidateRuneCountRange(req.Repository.Name, 4, 128)

	if resp, err := adm.PostReqBody(req, "/repository/"); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdRepositoryDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdRepositoryRestore(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	req := proto.Request{
		Flags: &proto.Flags{
			Restore: true,
		},
	}

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdRepositoryPurge(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	req := proto.Request{
		Flags: &proto.Flags{
			Purge: true,
		},
	}

	if resp, err := adm.DeleteReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdRepositoryClear(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	req := proto.Request{
		Flags: &proto.Flags{
			Clear: true,
		},
	}

	if resp, err := adm.PutReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
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
	id := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	var req proto.Request
	req.Repository = &proto.Repository{}
	req.Repository.Name = opts[`to`][0]

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
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
	id := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().First())
	_ = utl.TryGetTeamByUUIDOrName(Client, opts[`team`][0])
	path := fmt.Sprintf("/repository/%s", id)

	var req proto.Request
	req.Repository = &proto.Repository{}
	req.Repository.TeamId = opts[`to`][0]

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdRepositoryClone(c *cli.Context) error {
	utl.NotImplemented()
	return nil
}

func cmdRepositoryActivate(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	req := proto.Request{
		Flags: &proto.Flags{
			Activate: true,
		},
	}

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdRepositoryWipe(c *cli.Context) error {
	utl.NotImplemented()
	return nil
}

func cmdRepositoryList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}
	if resp, err := adm.GetReq("/repository/"); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdRepositoryShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `show`)
	}
}

func cmdRepositoryTree(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/repository/%s/tree/tree", id)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `tree`)
	}
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
	multiple := []string{}
	unique := []string{`from`, `view`}
	required := []string{`from`, `view`}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(opts, multiple, unique, required,
		c.Args().Tail()); err != nil {
		return err
	}
	repositoryId := utl.TryGetRepositoryByUUIDOrName(Client, opts[`from`][0])

	if pType == `system` {
		utl.CheckStringIsSystemProperty(Client, c.Args().First())
	}
	sourceId := utl.FindSourceForRepoProperty(Client, pType, c.Args().First(),
		opts[`view`][0], repositoryId)
	if sourceId == `` {
		return fmt.Errorf(`Could not find locally set requested property.`)
	}

	path := fmt.Sprintf("/repository/%s/property/%s/%s",
		repositoryId, pType, sourceId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `delete`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
