package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
)

func registerNodes(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// nodes
			{
				Name:  "nodes",
				Usage: "SUBCOMMANDS for nodes",
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Register a new node",
						Action: runtime(cmdNodeAdd),
					},
					{
						Name:   "delete",
						Usage:  "Mark a node as deleted",
						Action: runtime(cmdNodeDel),
					},
					{
						Name:   "purge",
						Usage:  "Purge a node marked as deleted",
						Action: runtime(cmdNodePurge),
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "all, a",
								Usage: "Purge all deleted nodes",
							},
						},
					},
					{
						Name:   "restore",
						Usage:  "Restore a node marked as deleted",
						Action: runtime(cmdNodeRestore),
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "all, a",
								Usage: "Restore all deleted nodes",
							},
						},
					},
					{
						Name:   "rename",
						Usage:  "Rename a node",
						Action: runtime(cmdNodeRename),
					},
					{
						Name:   "repossess",
						Usage:  "Repossess a node to a different team",
						Action: runtime(cmdNodeRepo),
					},
					{
						Name:   "relocate",
						Usage:  "Relocate a node to a different server",
						Action: runtime(cmdNodeMove),
					},
					{
						Name:   "online",
						Usage:  "Set a nodes to online",
						Action: runtime(cmdNodeOnline),
					},
					{
						Name:   "offline",
						Usage:  "Set a node to offline",
						Action: runtime(cmdNodeOffline),
					},
					{
						Name:   "assign",
						Usage:  "Assign a node to configuration bucket",
						Action: runtime(cmdNodeAssign),
					},
					{
						Name:   "list",
						Usage:  "List all nodes",
						Action: runtime(cmdNodeList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a node",
						Action: runtime(cmdNodeShow),
					},
					{
						Name:   "config",
						Usage:  "Show which bucket a node is assigned to",
						Action: runtime(cmdNodeConfig),
					},
					{
						Name:  "property",
						Usage: "SUBCOMMANDS for node properties",
						Subcommands: []cli.Command{
							{
								Name:  "add",
								Usage: "SUBCOMMANDS for property add",
								Subcommands: []cli.Command{
									{
										Name:   "system",
										Usage:  "Add a system property to a node",
										Action: runtime(cmdNodeSystemPropertyAdd),
									},
									{
										Name:   "service",
										Usage:  "Add a service property to a node",
										Action: runtime(cmdNodeServicePropertyAdd),
									},
								},
							},
							/*
								{
									Name:   "add",
									Usage:  "Assign a property to a node",
									Action: cmdNodePropertyAdd,
								},
								{
									Name:   "get",
									Usage:  "Get the value of a node's specific property",
									Action: cmdNodePropertyGet,
								},
								{
									Name:   "delete",
									Usage:  "Delete a property from a node",
									Action: cmdNodePropertyDel,
								},
								{
									Name:   "list",
									Usage:  "List a nodes' local properties",
									Action: cmdNodePropertyList,
									Flags: []cli.Flag{
										cli.BoolFlag{
											Name:  "all, a",
											Usage: "List a nodes full properties (incl. inherited)",
										},
									},
								},
								{
									Name:   "show",
									Usage:  "Show details about a nodes properties",
									Action: cmdNodePropertyShow,
								},
							*/
						},
					}, // end nodes property
				},
			}, // end nodes
		}...,
	)
	return &app
}

func cmdNodeAdd(c *cli.Context) error {
	keySlice := []string{"assetid", "name", "team", "server", "online"}
	reqSlice := []string{"assetid", "name", "team"}

	switch utl.GetCliArgumentCount(c) {
	case 6, 8, 10:
		break
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	argSlice := utl.GetFullArgumentSlice(c)

	options, optional := utl.ParseVariableArguments(keySlice, reqSlice, argSlice)
	req := proto.Request{}
	req.Node = &proto.Node{}

	utl.ValidateStringAsNodeAssetId(options["assetid"])
	if utl.SliceContainsString("online", optional) {
		utl.ValidateStringAsBool(options["online"])
		req.Node.IsOnline, _ = strconv.ParseBool(options["online"])
	} else {
		req.Node.IsOnline = true
	}
	if utl.SliceContainsString("server", optional) {
		req.Node.ServerId = utl.TryGetServerByUUIDOrName(Client, options["server"])
	}
	req.Node.AssetId, _ = strconv.ParseUint(options["assetid"], 10, 64)
	req.Node.Name = options["name"]
	req.Node.TeamId = utl.TryGetTeamByUUIDOrName(Client, options["team"])

	resp := utl.PostRequestWithBody(Client, req, "/nodes/")
	fmt.Println(resp)
	return nil
}

func cmdNodeDel(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdNodePurge(c *cli.Context) error {
	var (
		path string
		req  proto.Request
	)
	if c.Bool("all") {
		utl.ValidateCliArgumentCount(c, 0)
		path = "/nodes/"
	} else {
		utl.ValidateCliArgumentCount(c, 1)
		id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
		path = fmt.Sprintf("/nodes/%s", id)
	}

	req = proto.Request{
		Flags: &proto.Flags{
			Purge: true,
		},
	}

	resp := utl.DeleteRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdNodeRestore(c *cli.Context) error {
	var (
		path string
		req  proto.Request
	)
	if c.Bool("all") {
		utl.ValidateCliArgumentCount(c, 0)
		path = "/nodes/"
	} else {
		utl.ValidateCliArgumentCount(c, 1)
		id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
		path = fmt.Sprintf("/nodes/%s", id)
	}

	req = proto.Request{
		Flags: &proto.Flags{
			Restore: true,
		},
	}

	resp := utl.DeleteRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdNodeRename(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.Name = c.Args().Get(2)

	resp := utl.PatchRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdNodeRepo(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().Get(0))
	team := c.Args().Get(2)
	// try resolving team name to uuid as name validation
	_ = utl.GetTeamIdByName(Client, team)
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.TeamId = team

	resp := utl.PatchRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdNodeMove(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().Get(0))
	server := c.Args().Get(2)
	// try resolving server name to uuid as name validation
	_ = utl.GetServerAssetIdByName(Client, server)
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.ServerId = server

	resp := utl.PatchRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdNodeOnline(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.IsOnline = true

	resp := utl.PatchRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdNodeOffline(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.IsOnline = false

	resp := utl.PatchRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdNodeAssign(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	multiple := []string{}
	unique := []string{"to"}
	required := []string{"to"}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(Client, opts["to"][0])
	repoId := utl.GetRepositoryIdForBucket(Client, bucketId)
	nodeId := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.Id = nodeId
	req.Node.Config = &proto.NodeConfig{}
	req.Node.Config.RepositoryId = repoId
	req.Node.Config.BucketId = bucketId

	path := fmt.Sprintf("/nodes/%s/config", nodeId)
	resp := utl.PutRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdNodeList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest(Client, "/nodes/")
	fmt.Println(resp)
	return nil
}

func cmdNodeShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdNodeConfig(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/nodes/%s/config", id)

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdNodeSystemPropertyAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	required := []string{"to", "value", "view"}
	unique := []string{"to", "in", "value", "view", "inheritance", "childrenonly"}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	if _, ok := opts["in"]; ok {
		fmt.Fprintln(os.Stderr, "Hint: Keyword `in` is DEPRECATED for nodes, since they are global objects. Ignoring.")
	}

	nodeId := utl.TryGetNodeByUUIDOrName(Client, opts["to"][0])
	utl.CheckStringIsSystemProperty(Client, c.Args().First())

	config := utl.GetNodeConfigById(Client, nodeId)

	tprop := proto.Property{
		Type: "system",
		View: opts["view"][0],
		System: &proto.PropertySystem{
			Name:  c.Args().First(),
			Value: opts["value"][0],
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

	propList := []proto.Property{tprop}

	node := proto.Node{
		Id:         nodeId,
		Properties: &propList,
		Config:     config,
	}

	req := proto.Request{
		Node: &node,
	}

	path := fmt.Sprintf("/nodes/%s/property/system/", nodeId)
	resp := utl.PostRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdNodeServicePropertyAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 5)
	multiple := []string{}
	required := []string{"to", "view"}
	unique := []string{"to", "in", "view", "inheritance", "childrenonly"}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	if _, ok := opts["in"]; ok {
		fmt.Fprintln(os.Stderr, "Hint: Keyword `in` is DEPRECATED for nodes, since they are global objects. Ignoring.")
	}

	nodeId := utl.TryGetNodeByUUIDOrName(Client, opts["to"][0])
	config := utl.GetNodeConfigById(Client, nodeId)
	teamId := utl.TeamIdForBucket(Client, config.BucketId)

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
		Node: &proto.Node{
			Id:     nodeId,
			Config: config,
			Properties: &[]proto.Property{
				tprop,
			},
		},
	}

	path := fmt.Sprintf("/nodes/%s/property/service/", nodeId)
	resp := utl.PostRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
