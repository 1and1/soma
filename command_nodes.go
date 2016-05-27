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
				Name:   "nodes",
				Usage:  "SUBCOMMANDS for nodes",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Register a new node",
						Action: cmdNodeAdd,
					},
					{
						Name:   "delete",
						Usage:  "Mark a node as deleted",
						Action: cmdNodeDel,
					},
					{
						Name:   "purge",
						Usage:  "Purge a node marked as deleted",
						Action: cmdNodePurge,
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
						Action: cmdNodeRestore,
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
						Action: cmdNodeRename,
					},
					{
						Name:   "repossess",
						Usage:  "Repossess a node to a different team",
						Action: cmdNodeRepo,
					},
					{
						Name:   "relocate",
						Usage:  "Relocate a node to a different server",
						Action: cmdNodeMove,
					},
					{
						Name:   "online",
						Usage:  "Set a nodes to online",
						Action: cmdNodeOnline,
					},
					{
						Name:   "offline",
						Usage:  "Set a node to offline",
						Action: cmdNodeOffline,
					},
					{
						Name:   "assign",
						Usage:  "Assign a node to configuration bucket",
						Action: cmdNodeAssign,
					},
					{
						Name:   "list",
						Usage:  "List all nodes",
						Action: cmdNodeList,
					},
					{
						Name:   "show",
						Usage:  "Show details about a node",
						Action: cmdNodeShow,
					},
					{
						Name:   "config",
						Usage:  "Show which bucket a node is assigned to",
						Action: cmdNodeConfig,
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
										Action: cmdNodeSystemPropertyAdd,
									},
									{
										Name:   "service",
										Usage:  "Add a service property to a node",
										Action: cmdNodeServicePropertyAdd,
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

func cmdNodeAdd(c *cli.Context) {
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
		req.Node.ServerId = utl.TryGetServerByUUIDOrName(options["server"])
	}
	req.Node.AssetId, _ = strconv.ParseUint(options["assetid"], 10, 64)
	req.Node.Name = options["name"]
	req.Node.TeamId = utl.TryGetTeamByUUIDOrName(options["team"])

	resp := utl.PostRequestWithBody(req, "/nodes/")
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdNodeDel(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdNodePurge(c *cli.Context) {
	var (
		path string
		req  proto.Request
	)
	if c.Bool("all") {
		utl.ValidateCliArgumentCount(c, 0)
		path = "/nodes/"
	} else {
		utl.ValidateCliArgumentCount(c, 1)
		id := utl.TryGetNodeByUUIDOrName(c.Args().First())
		path = fmt.Sprintf("/nodes/%s", id)
	}

	req = proto.Request{
		Flags: &proto.Flags{
			Purge: true,
		},
	}

	resp := utl.DeleteRequestWithBody(req, path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdNodeRestore(c *cli.Context) {
	var (
		path string
		req  proto.Request
	)
	if c.Bool("all") {
		utl.ValidateCliArgumentCount(c, 0)
		path = "/nodes/"
	} else {
		utl.ValidateCliArgumentCount(c, 1)
		id := utl.TryGetNodeByUUIDOrName(c.Args().First())
		path = fmt.Sprintf("/nodes/%s", id)
	}

	req = proto.Request{
		Flags: &proto.Flags{
			Restore: true,
		},
	}

	resp := utl.DeleteRequestWithBody(req, path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdNodeRename(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.Name = c.Args().Get(2)

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdNodeRepo(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetNodeByUUIDOrName(c.Args().Get(0))
	team := c.Args().Get(2)
	// try resolving team name to uuid as name validation
	_ = utl.GetTeamIdByName(team)
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.TeamId = team

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdNodeMove(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetNodeByUUIDOrName(c.Args().Get(0))
	server := c.Args().Get(2)
	// try resolving server name to uuid as name validation
	_ = utl.GetServerAssetIdByName(server)
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.ServerId = server

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdNodeOnline(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.IsOnline = true

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdNodeOffline(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.IsOnline = false

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdNodeAssign(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multiple := []string{}
	unique := []string{"to"}
	required := []string{"to"}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(opts["to"][0])
	repoId := utl.GetRepositoryIdForBucket(bucketId)
	nodeId := utl.TryGetNodeByUUIDOrName(c.Args().First())

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.Id = nodeId
	req.Node.Config = &proto.NodeConfig{}
	req.Node.Config.RepositoryId = repoId
	req.Node.Config.BucketId = bucketId

	path := fmt.Sprintf("/nodes/%s/config", nodeId)
	resp := utl.PutRequestWithBody(req, path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdNodeList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest("/nodes/")
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdNodeShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdNodeConfig(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s/config", id)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdNodeSystemPropertyAdd(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	required := []string{"to", "value", "view"}
	unique := []string{"to", "in", "value", "view", "inheritance", "childrenonly"}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	if _, ok := opts["in"]; ok {
		fmt.Fprintln(os.Stderr, "Hint: Keyword `in` is DEPRECATED for nodes, since they are global objects. Ignoring.")
	}

	nodeId := utl.TryGetNodeByUUIDOrName(opts["to"][0])
	utl.CheckStringIsSystemProperty(c.Args().First())

	config := utl.GetNodeConfigById(nodeId)

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
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdNodeServicePropertyAdd(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 5)
	multiple := []string{}
	required := []string{"to", "view"}
	unique := []string{"to", "in", "view", "inheritance", "childrenonly"}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	if _, ok := opts["in"]; ok {
		fmt.Fprintln(os.Stderr, "Hint: Keyword `in` is DEPRECATED for nodes, since they are global objects. Ignoring.")
	}

	nodeId := utl.TryGetNodeByUUIDOrName(opts["to"][0])
	config := utl.GetNodeConfigById(nodeId)
	teamId := utl.TeamIdForBucket(config.BucketId)

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
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
