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
	req := somaproto.ProtoRequestNode{}
	req.Node = &somaproto.ProtoNode{}

	utl.ValidateStringAsNodeAssetId(options["assetid"])
	if utl.SliceContainsString("online", optional) {
		utl.ValidateStringAsBool(options["online"])
		req.Node.IsOnline, _ = strconv.ParseBool(options["online"])
	} else {
		req.Node.IsOnline = true
	}
	if utl.SliceContainsString("server", optional) {
		req.Node.Server = utl.TryGetServerByUUIDOrName(options["server"])
	}
	req.Node.AssetId, _ = strconv.ParseUint(options["assetid"], 10, 64)
	req.Node.Name = options["name"]
	req.Node.Team = utl.TryGetTeamByUUIDOrName(options["team"])

	resp := utl.PostRequestWithBody(req, "/nodes/")
	fmt.Println(resp)
}

func cmdNodeDel(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	_ = utl.DeleteRequest(path)
}

func cmdNodePurge(c *cli.Context) {
	var (
		path string
		req  somaproto.ProtoRequestNode
	)
	if c.Bool("all") {
		utl.ValidateCliArgumentCount(c, 0)
		path = "/nodes/"
	} else {
		utl.ValidateCliArgumentCount(c, 1)
		id := utl.TryGetNodeByUUIDOrName(c.Args().First())
		path = fmt.Sprintf("/nodes/%s", id)
	}

	req.Purge = true

	_ = utl.DeleteRequestWithBody(req, path)
}

func cmdNodeRestore(c *cli.Context) {
	var (
		path string
		req  somaproto.ProtoRequestNode
	)
	if c.Bool("all") {
		utl.ValidateCliArgumentCount(c, 0)
		path = "/nodes/"
	} else {
		utl.ValidateCliArgumentCount(c, 1)
		id := utl.TryGetNodeByUUIDOrName(c.Args().First())
		path = fmt.Sprintf("/nodes/%s", id)
	}

	req.Restore = true

	_ = utl.DeleteRequestWithBody(req, path)
}

func cmdNodeRename(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	req := somaproto.ProtoRequestNode{}
	req.Node = &somaproto.ProtoNode{}
	req.Node.Name = c.Args().Get(2)

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdNodeRepo(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetNodeByUUIDOrName(c.Args().Get(0))
	team := c.Args().Get(2)
	// try resolving team name to uuid as name validation
	_ = utl.GetTeamIdByName(team)
	path := fmt.Sprintf("/nodes/%s", id)

	req := somaproto.ProtoRequestNode{}
	req.Node = &somaproto.ProtoNode{}
	req.Node.Team = team

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdNodeMove(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetNodeByUUIDOrName(c.Args().Get(0))
	server := c.Args().Get(2)
	// try resolving server name to uuid as name validation
	_ = utl.GetServerAssetIdByName(server)
	path := fmt.Sprintf("/nodes/%s", id)

	req := somaproto.ProtoRequestNode{}
	req.Node = &somaproto.ProtoNode{}
	req.Node.Server = server

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdNodeOnline(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	req := somaproto.ProtoRequestNode{}
	req.Node = &somaproto.ProtoNode{}
	req.Node.IsOnline = true

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdNodeOffline(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	req := somaproto.ProtoRequestNode{}
	req.Node = &somaproto.ProtoNode{}
	req.Node.IsOnline = false

	_ = utl.PatchRequestWithBody(req, path)
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

	req := somaproto.ProtoRequestNode{}
	req.Node = &somaproto.ProtoNode{}
	req.Node.Id = nodeId
	req.Node.Config = &somaproto.ProtoNodeConfig{}
	req.Node.Config.RepositoryId = repoId
	req.Node.Config.BucketId = bucketId

	path := fmt.Sprintf("/nodes/%s/config", nodeId)
	resp := utl.PutRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdNodeList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest("/nodes/")
	fmt.Println(resp)
}

func cmdNodeShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdNodeConfig(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s/config", id)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
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

	sprop := somaproto.TreePropertySystem{
		Name:  c.Args().First(),
		Value: opts["value"][0],
	}

	tprop := somaproto.TreeProperty{
		PropertyType: "system",
		View:         opts["view"][0],
		System:       &sprop,
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

	propList := []somaproto.TreeProperty{tprop}

	node := somaproto.ProtoNode{
		Id:         nodeId,
		Properties: &propList,
		Config:     config,
	}

	req := somaproto.ProtoRequestNode{
		Node: &node,
	}

	path := fmt.Sprintf("/nodes/%s/property/system/", nodeId)
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
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
	tprop := somaproto.TreeProperty{
		PropertyType: "service",
		View:         opts["view"][0],
		Service: &somaproto.TreePropertyService{
			Name:       c.Args().First(),
			TeamId:     teamId,
			Attributes: []somaproto.TreeServiceAttribute{},
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

	req := somaproto.ProtoRequestNode{
		Node: &somaproto.ProtoNode{
			Id:     nodeId,
			Config: config,
			Properties: &[]somaproto.TreeProperty{
				tprop,
			},
		},
	}

	path := fmt.Sprintf("/nodes/%s/property/service/", nodeId)
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

/* XXX 0xDEADC0DE
func cmdNodePropertyAdd(c *cli.Context) {
	// preliminary argv validation
	switch utl.GetCliArgumentCount(c) {
	case 4, 6, 8, 10, 12:
		break
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	utl.ValidateCliArgument(c, 3, "to")
	argSlice := utl.GetFullArgumentSlice(c)

	// XXX THIS WHOLE THING IS BROKEN, NEEDS ADAPTION
	// XXX TO WORK WITH somaproto.TreeProperty
	// XXX IS JUST FIXED UP ENOUGH TO COMPILE

	// get property types
	typSlice := utl.PropertyTypes

	// first argument must be a valid property type
	utl.ValidateStringInSlice(argSlice[0], typSlice)

	// TODO: validate property of that type and name exists
	propertyType := argSlice[0]
	//property := argSlice[1]

	// get which node is being modified
	id := utl.TryGetNodeByUUIDOrName(argSlice[3])
	path := fmt.Sprintf("/nodes/%s/property/%s/", id, argSlice[0])

	// variable key/value part of argv
	argSlice = argSlice[4:]

	// define accepted and required keys
	keySlice := []string{"inheritance", "childrenonly", "view", "value"}
	reqSlice := []string{"view"}
	if propertyType != "service" {
		// non service properties require values, services are
		// predefined and do not
		reqSlice = append(reqSlice, "value")
	}
	options, optional := utl.ParseVariableArguments(keySlice, reqSlice, argSlice)

	// build node property JSON
	var prop somaproto.TreeProperty
	prop.PropertyType = propertyType
	prop.View = options["view"] //required
	//prop.Property = property XXX BROKEN
	// add value if it was required
	if utl.SliceContainsString("value", reqSlice) {
		//prop.Value = options["value"]
	}

	// optional inheritance, default true
	if utl.SliceContainsString("inheritance", optional) {
		utl.ValidateStringAsBool(options["inheritance"])
		prop.Inheritance, _ = strconv.ParseBool(options["inheritance"])
	} else {
		prop.Inheritance = true
	}

	// optional childrenonly, default false
	if utl.SliceContainsString("childrenonly", optional) {
		utl.ValidateStringAsBool(options["childrenonly"])
		prop.ChildrenOnly, _ = strconv.ParseBool(options["childrenonly"])
	} else {
		prop.ChildrenOnly = false
	}

	// build request JSON
	req := somaproto.ProtoRequestNode{}
	req.Node = &somaproto.ProtoNode{}
	*req.Node.Properties = append(*req.Node.Properties, prop)

	_ = utl.PostRequestWithBody(req, path)
	// TODO save jobid locally as outstanding
}

func cmdNodePropertyGet(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 6)
	utl.ValidateCliArgument(c, 3, "from")
	utl.ValidateCliArgument(c, 5, "view")
	argSlice := utl.GetFullArgumentSlice(c)

	// first argument must be a valid property type, sixth a view
	utl.ValidateStringInSlice(argSlice[0], utl.PropertyTypes)
	utl.ValidateStringInSlice(argSlice[5], utl.Views)

	// get which node is being modified
	id := utl.TryGetNodeByUUIDOrName(argSlice[3])

	// TODO: validate property of that type and name exists
	path := fmt.Sprintf("/nodes/%s/property/%s/%s/%s",
		id,
		argSlice[0], // type
		argSlice[5], // view
		argSlice[1], // property
	)

	_ = utl.GetRequest(path)
}

func cmdNodePropertyDel(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 6)
	utl.ValidateCliArgument(c, 3, "from")
	utl.ValidateCliArgument(c, 5, "view")
	argSlice := utl.GetFullArgumentSlice(c)

	// first argument must be a valid property type, sixth a view
	utl.ValidateStringInSlice(argSlice[0], utl.PropertyTypes)
	utl.ValidateStringInSlice(argSlice[5], utl.Views)

	// get which node is being modified
	id := utl.TryGetNodeByUUIDOrName(argSlice[3])

	// TODO: validate property of that type and name exists
	path := fmt.Sprintf("/nodes/%s/property/%s/%s/%s",
		id,
		argSlice[0], //type
		argSlice[5], //view
		argSlice[1], //property
	)

	_ = utl.DeleteRequest(path)
}

func cmdNodePropertyList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s/property/", id)

	req := somaproto.ProtoRequestNode{}
	req.Filter = &somaproto.ProtoNodeFilter{}

	if c.Bool("all") {
		req.Filter.LocalProperty = false
		_ = utl.GetRequest(path)
	} else {
		req.Filter.LocalProperty = true
		_ = utl.GetRequestWithBody(req, path)
	}
}

func cmdNodePropertyShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s/property/", id)

	_ = utl.GetRequest(path)
}
*/

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
