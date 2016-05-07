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
				Name:   "repository",
				Usage:  "SUBCOMMANDS for repository",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new repository",
						Action: cmdRepositoryCreate,
					},
					{
						Name:   "delete",
						Usage:  "Mark an existing repository as deleted",
						Action: cmdRepositoryDelete,
					},
					{
						Name:   "restore",
						Usage:  "Restore a repository marked as deleted",
						Action: cmdRepositoryRestore,
					},
					{
						Name:   "purge",
						Usage:  "Remove an unreferenced deleted repository",
						Action: cmdRepositoryPurge,
					},
					{
						Name:   "clear",
						Usage:  "Clear all check instances for this repository",
						Action: cmdRepositoryClear,
					},
					{
						Name:   "rename",
						Usage:  "Rename an existing repository",
						Action: cmdRepositoryRename,
					},
					{
						Name:   "repossess",
						Usage:  "Change the owner of a repository",
						Action: cmdRepositoryRepossess,
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
						Action: cmdRepositoryActivate,
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
						Action: cmdRepositoryList,
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific repository",
						Action: cmdRepositoryShow,
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
										Action: cmdRepositorySystemPropertyAdd,
									},
									{
										Name:   "service",
										Usage:  "Add a service property to a repository",
										Action: cmdRepositoryServicePropertyAdd,
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

func cmdRepositoryCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "team")

	teamId := utl.TryGetTeamByUUIDOrName(c.Args().Get(2))

	var req proto.Request
	req.Repository = &proto.Repository{}
	req.Repository.Name = c.Args().Get(0)
	req.Repository.TeamId = teamId

	resp := utl.PostRequestWithBody(req, "/repository/")
	fmt.Println(resp)
}

func cmdRepositoryDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	_ = utl.DeleteRequest(path)
}

func cmdRepositoryRestore(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	req := proto.Request{
		Flags: &proto.Flags{
			Restore: true,
		},
	}

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdRepositoryPurge(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	req := proto.Request{
		Flags: &proto.Flags{
			Purge: true,
		},
	}

	resp := utl.DeleteRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdRepositoryClear(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	req := proto.Request{
		Flags: &proto.Flags{
			Clear: true,
		},
	}

	resp := utl.PutRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdRepositoryRename(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	var req proto.Request
	req.Repository = &proto.Repository{}
	req.Repository.Name = c.Args().Get(2)

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdRepositoryRepossess(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	_ = utl.TryGetTeamByUUIDOrName(c.Args().Get(2))
	path := fmt.Sprintf("/repository/%s", id)

	var req proto.Request
	req.Repository = &proto.Repository{}
	req.Repository.TeamId = c.Args().Get(2)

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdRepositoryClone(c *cli.Context) {
	utl.NotImplemented()
}

func cmdRepositoryActivate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	req := proto.Request{
		Flags: &proto.Flags{
			Activate: true,
		},
	}

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdRepositoryWipe(c *cli.Context) {
	utl.NotImplemented()
}

func cmdRepositoryList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)
	resp := utl.GetRequest("/repository/")
	fmt.Println(resp)
}

func cmdRepositoryShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdRepositorySystemPropertyAdd(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	required := []string{"to", "value", "view"}
	unique := []string{"to", "in", "value", "view", "inheritance", "childrenonly"}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	if _, ok := opts["in"]; ok {
		fmt.Fprintln(os.Stderr, "Hint: Keyword `in` is DEPRECATED for repositories, since they are global objects. Ignoring.")
	}

	repositoryId := utl.TryGetRepositoryByUUIDOrName(opts["to"][0])
	utl.CheckStringIsSystemProperty(c.Args().First())

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
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdRepositoryServicePropertyAdd(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 5)
	multiple := []string{}
	required := []string{"to", "view"}
	unique := []string{"to", "in", "view", "inheritance", "childrenonly"}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	if _, ok := opts["in"]; ok {
		fmt.Fprintln(os.Stderr, "Hint: Keyword `in` is DEPRECATED for repositories, since they are global objects. Ignoring.")
	}

	repositoryId := utl.TryGetRepositoryByUUIDOrName(opts["to"][0])
	teamId := utl.GetTeamIdByRepositoryId(repositoryId)

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
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
