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
						Name:   `update`,
						Usage:  `Update a node's information`,
						Action: runtime(cmdNodeUpdate),
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
						Name:   "synclist",
						Usage:  "List all nodes suitable for sync",
						Action: runtime(cmdNodeSync),
					},
					{
						Name:   "show",
						Usage:  "Show details about a node",
						Action: runtime(cmdNodeShow),
					},
					{
						Name:   "tree",
						Usage:  "Display the most uninteresting tree ever",
						Action: runtime(cmdNodeTree),
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
									{
										Name:   `oncall`,
										Usage:  `Add an oncall property to a node`,
										Action: runtime(cmdNodeOncallPropertyAdd),
									},
								},
							},
							{
								Name:  `delete`,
								Usage: `SUBCOMMANDS for property delete`,
								Subcommands: []cli.Command{
									{
										Name:         `system`,
										Usage:        `Delete a system property from a node`,
										Action:       runtime(cmdNodeSystemPropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `service`,
										Usage:        `Delete a service property from a node`,
										Action:       runtime(cmdNodeServicePropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `oncall`,
										Usage:        `Delete an oncall property from a node`,
										Action:       runtime(cmdNodeOncallPropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `custom`,
										Usage:        `Delete a custom property from a node`,
										Action:       runtime(cmdNodeCustomPropertyDelete),
										BashComplete: cmpl.FromView,
									},
								},
							},
							/*
								{
									Name:   "get",
									Usage:  "Get the value of a node's specific property",
									Action: cmdNodePropertyGet,
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
		req.Node.ServerId = utl.TryGetServerByUUIDOrName(&store, Client, options["server"])
	}
	req.Node.AssetId, _ = strconv.ParseUint(options["assetid"], 10, 64)
	req.Node.Name = options["name"]
	req.Node.TeamId = utl.TryGetTeamByUUIDOrName(Client, options["team"])

	if resp, err := adm.PostReqBody(req, "/nodes/"); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdNodeUpdate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 13)
	multiple := []string{}
	unique := []string{`name`, `assetid`, `server`, `team`, `online`, `deleted`}
	required := []string{`name`, `assetid`, `server`, `team`, `online`, `deleted`}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	utl.ValidateStringAsNodeAssetId(opts[`assetid`][0])
	req := proto.NewNodeRequest()
	if !utl.IsUUID(c.Args().First()) {
		return fmt.Errorf(`Node/update command requires UUID as first argument`)
	}
	req.Node.Id = c.Args().First()
	req.Node.Name = opts[`name`][0]
	req.Node.TeamId = utl.TryGetTeamByUUIDOrName(Client, opts[`team`][0])
	req.Node.IsOnline = utl.GetValidatedBool(opts[`online`][0])
	req.Node.IsDeleted = utl.GetValidatedBool(opts[`deleted`][0])
	req.Node.ServerId = utl.TryGetServerByUUIDOrName(&store, Client, opts[`server`][0])
	req.Node.AssetId, _ = strconv.ParseUint(opts[`assetid`][0], 10, 64)
	path := fmt.Sprintf("/nodes/%s", req.Node.Id)
	if resp, err := adm.PutReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdNodeDel(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
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

	if resp, err := adm.DeleteReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
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

	if resp, err := adm.DeleteReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdNodeRename(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.Name = c.Args().Get(2)

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
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

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
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

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdNodeOnline(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.IsOnline = true

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdNodeOffline(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.IsOnline = false

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
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

	bucketTeamId := utl.TeamIdForBucket(Client, bucketId)
	nodeTeamId := utl.TeamIdForNode(Client, nodeId)
	if bucketTeamId != nodeTeamId {
		utl.Abort(`Cannot assign node since node and bucket belong to different teams.`)
	}

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.Id = nodeId
	req.Node.Config = &proto.NodeConfig{}
	req.Node.Config.RepositoryId = repoId
	req.Node.Config.BucketId = bucketId

	path := fmt.Sprintf("/nodes/%s/config", nodeId)
	if resp, err := adm.PutReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdNodeList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)

	if resp, err := adm.GetReq("/nodes/"); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdNodeShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdNodeTree(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/nodes/%s/tree/tree", id)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `tree`)
	}
}

func cmdNodeSync(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)

	if resp, err := adm.GetReq(`/sync/nodes/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdNodeConfig(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/nodes/%s/config", id)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
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
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
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
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdNodeOncallPropertyAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 5)
	multiple := []string{}
	required := []string{"to", "view"}
	unique := []string{"to", "in", "view", "inheritance", "childrenonly"}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	if _, ok := opts["in"]; ok {
		fmt.Fprintln(os.Stderr, "Hint: Keyword `in` is DEPRECATED for nodes, since they are global objects. Ignoring.")
	}

	nodeId := utl.TryGetNodeByUUIDOrName(Client, opts["to"][0])
	oncallId := utl.TryGetOncallByUUIDOrName(Client, c.Args().First())
	oprop := proto.PropertyOncall{
		Id: oncallId,
	}
	oprop.Name, oprop.Number = utl.GetOncallDetailsById(Client, oncallId)

	config := utl.GetNodeConfigById(Client, nodeId)

	tprop := proto.Property{
		Type:   `oncall`,
		View:   opts["view"][0],
		Oncall: &oprop,
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

	path := fmt.Sprintf("/nodes/%s/property/oncall/", nodeId)
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdNodeSystemPropertyDelete(c *cli.Context) error {
	return cmdNodePropertyDelete(c, `system`)
}

func cmdNodeServicePropertyDelete(c *cli.Context) error {
	return cmdNodePropertyDelete(c, `service`)
}

func cmdNodeOncallPropertyDelete(c *cli.Context) error {
	return cmdNodePropertyDelete(c, `oncall`)
}

func cmdNodeCustomPropertyDelete(c *cli.Context) error {
	return cmdNodePropertyDelete(c, `custom`)
}

func cmdNodePropertyDelete(c *cli.Context, pType string) error {
	utl.ValidateCliMinArgumentCount(c, 5)
	multiple := []string{}
	unique := []string{`from`, `view`}
	required := []string{`from`, `view`}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	nodeId := utl.TryGetNodeByUUIDOrName(Client, opts[`from`][0])
	config := utl.GetNodeConfigById(Client, nodeId)

	if pType == `system` {
		utl.CheckStringIsSystemProperty(Client, c.Args().First())
	}
	sourceId := utl.FindSourceForNodeProperty(Client, pType, c.Args().First(),
		opts[`view`][0], nodeId)
	if sourceId == `` {
		return fmt.Errorf(`Could not find locally set requested property.`)
	}

	req := proto.NewNodeRequest()
	req.Node.Id = nodeId
	req.Node.Config = config
	path := fmt.Sprintf("/nodes/%s/property/%s/%s",
		nodeId, pType, sourceId)

	if resp, err := adm.DeleteReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `delete`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
