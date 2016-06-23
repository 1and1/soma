package main

import (
	"fmt"
	"os"

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
						Name:   "create",
						Usage:  "Create a new repository",
						Action: runtime(cmdRepositoryCreate),
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
						Name:   "rename",
						Usage:  "Rename an existing repository",
						Action: runtime(cmdRepositoryRename),
					},
					{
						Name:   "repossess",
						Usage:  "Change the owner of a repository",
						Action: runtime(cmdRepositoryRepossess),
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
										Name:   "system",
										Usage:  "Add a system property to a repository",
										Action: runtime(cmdRepositorySystemPropertyAdd),
									},
									{
										Name:   "service",
										Usage:  "Add a service property to a repository",
										Action: runtime(cmdRepositoryServicePropertyAdd),
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
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "team")

	teamId := utl.TryGetTeamByUUIDOrName(Client, c.Args().Get(2))

	var req proto.Request
	req.Repository = &proto.Repository{}
	req.Repository.Name = c.Args().Get(0)
	req.Repository.TeamId = teamId

	utl.ValidateRuneCountRange(req.Repository.Name, 4, 128)

	if resp, err := adm.PostReqBody(req, "/repository/"); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdRepositoryDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdRepositoryRestore(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
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
		fmt.Println(resp)
	}
	return nil
}

func cmdRepositoryPurge(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
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
		fmt.Println(resp)
	}
	return nil
}

func cmdRepositoryClear(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
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
		fmt.Println(resp)
	}
	return nil
}

func cmdRepositoryRename(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	var req proto.Request
	req.Repository = &proto.Repository{}
	req.Repository.Name = c.Args().Get(2)

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdRepositoryRepossess(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().First())
	_ = utl.TryGetTeamByUUIDOrName(Client, c.Args().Get(2))
	path := fmt.Sprintf("/repository/%s", id)

	var req proto.Request
	req.Repository = &proto.Repository{}
	req.Repository.TeamId = c.Args().Get(2)

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdRepositoryClone(c *cli.Context) error {
	utl.NotImplemented()
	return nil
}

func cmdRepositoryActivate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
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
		fmt.Println(resp)
	}
	return nil
}

func cmdRepositoryWipe(c *cli.Context) error {
	utl.NotImplemented()
	return nil
}

func cmdRepositoryList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)
	if resp, err := adm.GetReq("/repository/"); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdRepositoryShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdRepositoryTree(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/repository/%s/tree", id)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdRepositorySystemPropertyAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	required := []string{"to", "value", "view"}
	unique := []string{"to", "in", "value", "view", "inheritance", "childrenonly"}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	if _, ok := opts["in"]; ok {
		fmt.Fprintln(os.Stderr, "Hint: Keyword `in` is DEPRECATED for repositories, since they are global objects. Ignoring.")
	}

	repositoryId := utl.TryGetRepositoryByUUIDOrName(Client, opts["to"][0])
	utl.CheckStringIsSystemProperty(Client, c.Args().First())

	sprop := proto.PropertySystem{
		Name:  c.Args().First(),
		Value: opts["value"][0],
	}

	tprop := proto.Property{
		Type:   "system",
		View:   opts["view"][0],
		System: &sprop,
	}
	if _, ok := opts["inheritance"]; ok {
		tprop.Inheritance = utl.GetValidatedBool(opts["inheritance"][0])
	} else {
		tprop.Inheritance = true
	}
	if _, ok := opts["childrenonly"]; ok {
		tprop.ChildrenOnly = utl.GetValidatedBool(opts["childrenonly"][0])
	} else {
		tprop.ChildrenOnly = false
	}

	propList := []proto.Property{tprop}

	repository := proto.Repository{
		Id:         repositoryId,
		Properties: &propList,
	}

	req := proto.Request{
		Repository: &repository,
	}

	path := fmt.Sprintf("/repository/%s/property/system/", repositoryId)
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdRepositoryServicePropertyAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 5)
	multiple := []string{}
	required := []string{"to", "view"}
	unique := []string{"to", "in", "view", "inheritance", "childrenonly"}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	if _, ok := opts["in"]; ok {
		fmt.Fprintln(os.Stderr, "Hint: Keyword `in` is DEPRECATED for repositories, since they are global objects. Ignoring.")
	}

	repositoryId := utl.TryGetRepositoryByUUIDOrName(Client, opts["to"][0])
	teamId := utl.GetTeamIdByRepositoryId(Client, repositoryId)

	// no reason to fill out the attributes, client-provided
	// attributes are discarded by the server
	tprop := proto.Property{
		Type: "service",
		View: opts["view"][0],
		Service: &proto.PropertyService{
			Name:       c.Args().First(),
			TeamId:     teamId,
			Attributes: []proto.ServiceAttribute{},
		},
	}
	if _, ok := opts["inheritance"]; ok {
		tprop.Inheritance = utl.GetValidatedBool(opts["inheritance"][0])
	} else {
		tprop.Inheritance = true
	}
	if _, ok := opts["childrenonly"]; ok {
		tprop.ChildrenOnly = utl.GetValidatedBool(opts["childrenonly"][0])
	} else {
		tprop.ChildrenOnly = false
	}

	req := proto.Request{
		Repository: &proto.Repository{
			Id: repositoryId,
			Properties: &[]proto.Property{
				tprop,
			},
		},
	}

	path := fmt.Sprintf("/repository/%s/property/service/", repositoryId)
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
