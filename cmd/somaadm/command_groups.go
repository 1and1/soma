package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerGroups(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// groups
			{
				Name:  "groups",
				Usage: "SUBCOMMANDS for groups",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new group",
						Action:       runtime(cmdGroupCreate),
						BashComplete: cmpl.In,
					},
					{
						Name:         "delete",
						Usage:        "Delete a group",
						Action:       runtime(cmdGroupDelete),
						BashComplete: cmpl.In,
					},
					{
						Name:         "rename",
						Usage:        "Rename a group",
						Action:       runtime(cmdGroupRename),
						BashComplete: cmpl.InTo,
					},
					{
						Name:   "list",
						Usage:  "List all groups",
						Action: runtime(cmdGroupList),
					},
					{
						Name:         "show",
						Usage:        "Show details about a group",
						Action:       runtime(cmdGroupShow),
						BashComplete: cmpl.In,
					},
					{
						Name:         `tree`,
						Usage:        `Display the group as tree`,
						Action:       runtime(cmdGroupTree),
						BashComplete: cmpl.In,
					},
					{
						Name:  "members",
						Usage: "SUBCOMMANDS for members",
						Subcommands: []cli.Command{
							{
								Name:  "add",
								Usage: "SUBCOMMANDS for members add",
								Subcommands: []cli.Command{
									{
										Name:         "group",
										Usage:        "Add a group to a group",
										Action:       runtime(cmdGroupMemberAddGroup),
										BashComplete: cmpl.InTo,
									},
									{
										Name:         "cluster",
										Usage:        "Add a cluster to a group",
										Action:       runtime(cmdGroupMemberAddCluster),
										BashComplete: cmpl.InTo,
									},
									{
										Name:         "node",
										Usage:        "Add a node to a group",
										Action:       runtime(cmdGroupMemberAddNode),
										BashComplete: cmpl.InTo,
									},
								},
							},
							{
								Name:  "delete",
								Usage: "SUBCOMMANDS for members delete",
								Subcommands: []cli.Command{
									{
										Name:         "group",
										Usage:        "Delete a group from a group",
										Action:       runtime(cmdGroupMemberDeleteGroup),
										BashComplete: cmpl.InFrom,
									},
									{
										Name:         "cluster",
										Usage:        "Delete a cluster from a group",
										Action:       runtime(cmdGroupMemberDeleteCluster),
										BashComplete: cmpl.InFrom,
									},
									{
										Name:         "node",
										Usage:        "Delete a node from a group",
										Action:       runtime(cmdGroupMemberDeleteNode),
										BashComplete: cmpl.InFrom,
									},
								},
							},
							{
								Name:         "list",
								Usage:        "List all members of a group",
								Action:       runtime(cmdGroupMemberList),
								BashComplete: cmpl.In,
							},
						},
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
										Usage:        "Add a system property to a group",
										Action:       runtime(cmdGroupSystemPropertyAdd),
										BashComplete: cmpl.PropertyAddValue,
									},
									{
										Name:         "service",
										Usage:        "Add a service property to a group",
										Action:       runtime(cmdGroupServicePropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         `oncall`,
										Usage:        `Add an oncall property to a group`,
										Action:       runtime(cmdGroupOncallPropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         `custom`,
										Usage:        `Add a custom property to a group`,
										Action:       runtime(cmdGroupCustomPropertyAdd),
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
										Usage:        `Delete a system property from a group`,
										Action:       runtime(cmdGroupSystemPropertyDelete),
										BashComplete: cmpl.InFromView,
									},
									{
										Name:         `service`,
										Usage:        `Delete a service property from a group`,
										Action:       runtime(cmdGroupOncallPropertyDelete),
										BashComplete: cmpl.InFromView,
									},
									{
										Name:         `oncall`,
										Usage:        `Delete an oncall property from a group`,
										Action:       runtime(cmdGroupOncallPropertyDelete),
										BashComplete: cmpl.InFromView,
									},
									{
										Name:         `custom`,
										Usage:        `Delete a custom property from a group`,
										Action:       runtime(cmdGroupCustomPropertyDelete),
										BashComplete: cmpl.InFromView,
									},
								},
							},
						},
					},
				},
			}, // end groups
		}...,
	)
	return &app
}

func cmdGroupCreate(c *cli.Context) error {
	multKeys := []string{"in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail()); err != nil {
		return err
	}

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])

	var req proto.Request
	req.Group = &proto.Group{}
	req.Group.Name = c.Args().First()
	req.Group.BucketId = bucketId

	utl.ValidateRuneCountRange(req.Group.Name, 4, 256)

	if resp, err := adm.PostReqBody(req, "/groups/"); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupDelete(c *cli.Context) error {
	multKeys := []string{"in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail()); err != nil {
		return err
	}

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	groupId := utl.TryGetGroupByUUIDOrName(Client,
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/groups/%s", groupId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupRename(c *cli.Context) error {
	multKeys := []string{"to", "in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail()); err != nil {
		return err
	}

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	groupId := utl.TryGetGroupByUUIDOrName(Client,
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/groups/%s", groupId)

	var req proto.Request
	req.Group = &proto.Group{}
	req.Group.Name = opts["to"][0]

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}
	if resp, err := adm.GetReq("/groups/"); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdGroupShow(c *cli.Context) error {
	multKeys := []string{"in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multKeys,
		multKeys,
		multKeys,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	groupId := utl.TryGetGroupByUUIDOrName(Client,
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/groups/%s", groupId)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `show`)
	}
}

func cmdGroupTree(c *cli.Context) error {
	multKeys := []string{"in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multKeys,
		multKeys,
		multKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	groupId := utl.TryGetGroupByUUIDOrName(Client,
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/groups/%s/tree/tree", groupId)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `tree`)
	}
}

func cmdGroupMemberAddGroup(c *cli.Context) error {
	multKeys := []string{"to", "in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multKeys,
		multKeys,
		multKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	mGroupId := utl.TryGetGroupByUUIDOrName(Client,
		c.Args().First(),
		bucketId)
	groupId := utl.TryGetGroupByUUIDOrName(Client,
		opts["to"][0],
		bucketId)

	var req proto.Request
	var group proto.Group
	group.Id = mGroupId
	req.Group = &proto.Group{}
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	*req.Group.MemberGroups = append(*req.Group.MemberGroups, group)

	path := fmt.Sprintf("/groups/%s/members/", groupId)

	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupMemberAddCluster(c *cli.Context) error {
	multKeys := []string{"to", "in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multKeys,
		multKeys,
		multKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	mClusterId := utl.TryGetClusterByUUIDOrName(Client,
		c.Args().First(),
		bucketId)
	groupId := utl.TryGetGroupByUUIDOrName(Client,
		opts["to"][0],
		bucketId)

	var req proto.Request
	var cluster proto.Cluster
	cluster.Id = mClusterId
	req.Group = &proto.Group{}
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	*req.Group.MemberClusters = append(*req.Group.MemberClusters, cluster)

	path := fmt.Sprintf("/groups/%s/members/", groupId)

	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupMemberAddNode(c *cli.Context) error {
	multKeys := []string{"to", "in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multKeys,
		multKeys,
		multKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	mNodeId := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	groupId := utl.TryGetGroupByUUIDOrName(Client,
		opts["to"][0],
		bucketId)

	var req proto.Request
	var node proto.Node
	node.Id = mNodeId
	req.Group = &proto.Group{}
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	*req.Group.MemberNodes = append(*req.Group.MemberNodes, node)

	path := fmt.Sprintf("/groups/%s/members/", groupId)

	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupMemberDeleteGroup(c *cli.Context) error {
	multKeys := []string{"from", "in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multKeys,
		multKeys,
		multKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	mGroupId := utl.TryGetGroupByUUIDOrName(Client,
		c.Args().First(),
		bucketId)
	groupId := utl.TryGetGroupByUUIDOrName(Client,
		opts["from"][0],
		bucketId)

	path := fmt.Sprintf("/groups/%s/members/%s", groupId,
		mGroupId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupMemberDeleteCluster(c *cli.Context) error {
	multKeys := []string{"from", "in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multKeys,
		multKeys,
		multKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	mClusterId := utl.TryGetClusterByUUIDOrName(Client,
		c.Args().First(),
		bucketId)
	groupId := utl.TryGetGroupByUUIDOrName(Client,
		opts["from"][0],
		bucketId)

	path := fmt.Sprintf("/groups/%s/members/%s", groupId,
		mClusterId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupMemberDeleteNode(c *cli.Context) error {
	multKeys := []string{"from", "in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multKeys,
		multKeys,
		multKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	mNodeId := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	groupId := utl.TryGetGroupByUUIDOrName(Client,
		opts["from"][0],
		bucketId)

	path := fmt.Sprintf("/groups/%s/members/%s", groupId,
		mNodeId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupMemberList(c *cli.Context) error {
	multKeys := []string{"in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multKeys,
		multKeys,
		multKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	groupId := utl.TryGetGroupByUUIDOrName(Client,
		c.Args().First(),
		bucketId)

	path := fmt.Sprintf("/groups/%s/members/", groupId)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupSystemPropertyAdd(c *cli.Context) error {
	return cmdGroupPropertyAdd(c, `system`)
}

func cmdGroupServicePropertyAdd(c *cli.Context) error {
	return cmdGroupPropertyAdd(c, `service`)
}

func cmdGroupOncallPropertyAdd(c *cli.Context) error {
	return cmdGroupPropertyAdd(c, `oncall`)
}

func cmdGroupCustomPropertyAdd(c *cli.Context) error {
	return cmdGroupPropertyAdd(c, `custom`)
}

func cmdGroupPropertyAdd(c *cli.Context, pType string) error {
	return cmdPropertyAdd(c, pType, `group`)
}

func cmdGroupSystemPropertyDelete(c *cli.Context) error {
	return cmdGroupPropertyDelete(c, `system`)
}

func cmdGroupServicePropertyDelete(c *cli.Context) error {
	return cmdGroupPropertyDelete(c, `service`)
}

func cmdGroupOncallPropertyDelete(c *cli.Context) error {
	return cmdGroupPropertyDelete(c, `oncall`)
}

func cmdGroupCustomPropertyDelete(c *cli.Context) error {
	return cmdGroupPropertyDelete(c, `custom`)
}

func cmdGroupPropertyDelete(c *cli.Context, pType string) error {
	multiple := []string{}
	unique := []string{`from`, `view`, `in`}
	required := []string{`from`, `view`, `in`}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(opts, multiple, unique,
		required, c.Args().Tail()); err != nil {
		return err
	}
	bucketId := utl.BucketByUUIDOrName(Client, opts[`in`][0])
	groupId := utl.TryGetGroupByUUIDOrName(Client, opts[`from`][0], bucketId)

	if pType == `system` {
		utl.CheckStringIsSystemProperty(Client, c.Args().First())
	}
	sourceId := utl.FindSourceForGroupProperty(Client, pType, c.Args().First(),
		opts[`view`][0], groupId)
	if sourceId == `` {
		adm.Abort(`Could not find locally set requested property.`)
	}

	req := proto.NewGroupRequest()
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	path := fmt.Sprintf("/groups/%s/property/%s/%s",
		groupId, pType, sourceId)

	if resp, err := adm.DeleteReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `delete`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
