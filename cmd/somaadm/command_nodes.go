package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
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
						Name:         "create",
						Usage:        "Register a new node",
						Action:       runtime(cmdNodeAdd),
						BashComplete: cmpl.NodeAdd,
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
						Name:         `update`,
						Usage:        `Update a node's information`,
						Action:       runtime(cmdNodeUpdate),
						BashComplete: cmpl.NodeUpdate,
					},
					{
						Name:         "rename",
						Usage:        "Rename a node",
						Action:       runtime(cmdNodeRename),
						BashComplete: cmpl.To,
					},
					{
						Name:         "repossess",
						Usage:        "Repossess a node to a different team",
						Action:       runtime(cmdNodeRepo),
						BashComplete: cmpl.To,
					},
					{
						Name:         "relocate",
						Usage:        "Relocate a node to a different server",
						Action:       runtime(cmdNodeMove),
						BashComplete: cmpl.To,
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
						Name:         "assign",
						Usage:        "Assign a node to configuration bucket",
						Action:       runtime(cmdNodeAssign),
						BashComplete: cmpl.To,
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
										Name:         "system",
										Usage:        "Add a system property to a node",
										Action:       runtime(cmdNodeSystemPropertyAdd),
										BashComplete: cmpl.PropertyAddValue,
									},
									{
										Name:         "service",
										Usage:        "Add a service property to a node",
										Action:       runtime(cmdNodeServicePropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         `oncall`,
										Usage:        `Add an oncall property to a node`,
										Action:       runtime(cmdNodeOncallPropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         `custom`,
										Usage:        `Add a custom property to a node`,
										Action:       runtime(cmdNodeCustomPropertyAdd),
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
						},
					}, // end nodes property
				},
			}, // end nodes
		}...,
	)
	return &app
}

func cmdNodeAdd(c *cli.Context) error {
	opts := map[string][]string{}
	uniqKeys := []string{`assetid`, `name`, `team`, `server`, `online`}
	reqKeys := []string{`assetid`, `name`, `team`}

	var err error
	if err = adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		reqKeys,
		adm.AllArguments(c),
	); err != nil {
		return err
	}
	req := proto.NewNodeRequest()

	if _, ok := opts[`online`]; ok {
		if err = adm.ValidateBool(opts[`online`][0],
			&req.Node.IsOnline); err != nil {
			return err
		}
	} else {
		req.Node.IsOnline = true
	}
	if _, ok := opts[`server`]; ok {
		if req.Node.ServerId, err = adm.LookupServerId(
			opts[`server`][0]); err != nil {
			return err
		}
	}
	if err = adm.ValidateLBoundUint64(opts[`assetid`][0],
		&req.Node.AssetId, 1); err != nil {
		return err
	}
	req.Node.Name = opts[`name`][0]
	if req.Node.TeamId, err = adm.LookupTeamId(
		opts[`team`][0]); err != nil {
		return nil
	}

	return adm.Perform(`postbody`, `/nodes/`, `command`, req, c)
}

func cmdNodeUpdate(c *cli.Context) error {
	unique := []string{`name`, `assetid`, `server`, `team`,
		`online`, `deleted`}
	required := []string{`name`, `assetid`, `server`, `team`,
		`online`, `deleted`}
	opts := map[string][]string{}

	var err error
	if err = adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewNodeRequest()
	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf(
			`Node/update command requires UUID as first argument`)
	}
	req.Node.Id = c.Args().First()
	req.Node.Name = opts[`name`][0]
	if err = adm.ValidateBool(opts[`online`][0],
		&req.Node.IsOnline); err != nil {
		return err
	}
	if err = adm.ValidateBool(opts[`deleted`][0],
		&req.Node.IsDeleted); err != nil {
		return err
	}
	if req.Node.ServerId, err = adm.LookupServerId(
		opts[`server`][0]); err != nil {
		return err
	}
	if err = adm.ValidateLBoundUint64(opts[`assetid`][0],
		&req.Node.AssetId, 1); err != nil {
		return err
	}
	if req.Node.TeamId, err = adm.LookupTeamId(
		opts[`team`][0]); err != nil {
		return err
	}
	path := fmt.Sprintf("/nodes/%s", req.Node.Id)
	return adm.Perform(`putbody`, path, `command`, req, c)
}

func cmdNodeDel(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	var path string
	if id, err := adm.LookupNodeId(c.Args().First()); err != nil {
		return err
	} else {
		path = fmt.Sprintf("/nodes/%s", id)
	}

	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdNodePurge(c *cli.Context) error {
	var (
		path string
		req  proto.Request
	)
	if c.Bool("all") {
		if err := adm.VerifyNoArgument(c); err != nil {
			return err
		}
		path = "/nodes/"
	} else {
		if err := adm.VerifySingleArgument(c); err != nil {
			return err
		}
		if id, err := adm.LookupNodeId(c.Args().First()); err != nil {
			return err
		} else {
			path = fmt.Sprintf("/nodes/%s", id)
		}
	}

	req = proto.Request{
		Flags: &proto.Flags{
			Purge: true,
		},
	}

	return adm.Perform(`deletebody`, path, `command`, req, c)
}

func cmdNodeRestore(c *cli.Context) error {
	var (
		path string
		req  proto.Request
	)
	if c.Bool("all") {
		if err := adm.VerifyNoArgument(c); err != nil {
			return err
		}
		path = "/nodes/"
	} else {
		if err := adm.VerifySingleArgument(c); err != nil {
			return err
		}
		if id, err := adm.LookupNodeId(c.Args().First()); err != nil {
			return err
		} else {
			path = fmt.Sprintf("/nodes/%s", id)
		}
	}

	req = proto.Request{
		Flags: &proto.Flags{
			Restore: true,
		},
	}

	return adm.Perform(`deletebody`, path, `command`, req, c)
}

func cmdNodeRename(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`to`},
		[]string{`to`},
		c.Args().Tail()); err != nil {
		return err
	}
	var path string
	if id, err := adm.LookupNodeId(c.Args().First()); err != nil {
		return err
	} else {
		path = fmt.Sprintf("/nodes/%s", id)
	}

	req := proto.NewNodeRequest()
	req.Node.Name = opts[`to`][0]

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdNodeRepo(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`to`},
		[]string{`to`},
		c.Args().Tail()); err != nil {
		return err
	}
	var id, teamId string
	{
		var err error
		if id, err = adm.LookupNodeId(c.Args().First()); err != nil {
			return err
		}
		if teamId, err = adm.LookupTeamId(opts[`to`][0]); err != nil {
			return err
		}
	}
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.TeamId = teamId

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdNodeMove(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`to`},
		[]string{`to`},
		c.Args().Tail()); err != nil {
		return err
	}
	var id string
	{
		var err error
		if id, err = adm.LookupNodeId(c.Args().First()); err != nil {
			return err
		}
	}
	server := opts[`to`][0]
	serverId, err := adm.LookupServerId(server)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.ServerId = serverId

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdNodeOnline(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupNodeId(c.Args().First())
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.IsOnline = true

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdNodeOffline(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupNodeId(c.Args().First())
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/nodes/%s", id)

	req := proto.Request{}
	req.Node = &proto.Node{}
	req.Node.IsOnline = false

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdNodeAssign(c *cli.Context) error {
	multiple := []string{}
	unique := []string{"to"}
	required := []string{"to"}

	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		multiple,
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}
	var (
		err                      error
		bucketId, repoId, nodeId string
		bucketTId, nodeTId       string
	)
	if bucketId, err = adm.LookupBucketId(opts["to"][0]); err != nil {
		return err
	}
	if repoId, err = adm.LookupRepoByBucket(bucketId); err != nil {
		return err
	}
	if nodeId, err = adm.LookupNodeId(c.Args().First()); err != nil {
		return err
	}
	if bucketTId, err = adm.LookupTeamByBucket(bucketId); err != nil {
		return err
	}
	if nodeTId, err = adm.LookupTeamByNode(nodeId); err != nil {
		return err
	}
	if bucketTId != nodeTId {
		return fmt.Errorf(
			`Cannot assign node since node and bucket belong to` +
				` different teams.`)
	}

	req := proto.NewNodeRequest()
	req.Node.Id = nodeId
	req.Node.Config = &proto.NodeConfig{}
	req.Node.Config.RepositoryId = repoId
	req.Node.Config.BucketId = bucketId

	path := fmt.Sprintf("/nodes/%s/config", nodeId)
	return adm.Perform(`putbody`, path, `command`, req, c)
}

func cmdNodeList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/nodes/`, `list`, nil, c)
}

func cmdNodeShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupNodeId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/nodes/%s", id)
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdNodeTree(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupNodeId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/nodes/%s/tree/tree", id)
	return adm.Perform(`get`, path, `tree`, nil, c)
}

func cmdNodeSync(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/sync/nodes/`, `list`, nil, c)
}

func cmdNodeConfig(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupNodeId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/nodes/%s/config", id)
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdNodeSystemPropertyAdd(c *cli.Context) error {
	return cmdNodePropertyAdd(c, `system`)
}

func cmdNodeServicePropertyAdd(c *cli.Context) error {
	return cmdNodePropertyAdd(c, `service`)
}

func cmdNodeOncallPropertyAdd(c *cli.Context) error {
	return cmdNodePropertyAdd(c, `oncall`)
}

func cmdNodeCustomPropertyAdd(c *cli.Context) error {
	return cmdNodePropertyAdd(c, `custom`)
}

func cmdNodePropertyAdd(c *cli.Context, pType string) error {
	return cmdPropertyAdd(c, pType, `node`)
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
	var (
		err              error
		nodeId, sourceId string
		config           *proto.NodeConfig
	)
	if nodeId, err = adm.LookupNodeId(c.Args().First()); err != nil {
		return err
	}
	if config, err = adm.LookupNodeConfig(nodeId); err != nil {
		return err
	}
	if pType == `system` {
		if err = adm.ValidateSystemProperty(
			c.Args().First()); err != nil {
			return err
		}
	}
	if err = adm.FindNodePropSrcId(pType, c.Args().First(),
		opts[`view`][0], nodeId, &sourceId); err != nil {
		return err
	}

	req := proto.NewNodeRequest()
	req.Node.Id = nodeId
	req.Node.Config = config

	path := fmt.Sprintf("/nodes/%s/property/%s/%s",
		nodeId, pType, sourceId)
	return adm.Perform(`deletebody`, path, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
