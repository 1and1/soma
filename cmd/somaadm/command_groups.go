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
										Action:       runtime(cmdGroupServicePropertyDelete),
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
	uniqKeys := []string{"in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	bucketId, err := adm.LookupBucketId(opts["in"][0])
	if err != nil {
		return err
	}

	var req proto.Request
	req.Group = &proto.Group{}
	req.Group.Name = c.Args().First()
	req.Group.BucketId = bucketId

	if err := adm.ValidateRuneCountRange(req.Group.Name, 4, 256); err != nil {
		return err
	}

	if resp, err := adm.PostReqBody(req, "/groups/"); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupDelete(c *cli.Context) error {
	uniqKeys := []string{"in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	var (
		err               error
		bucketId, groupId string
	)
	if bucketId, err = adm.LookupBucketId(opts["in"][0]); err != nil {
		return err
	}
	if groupId, err = adm.LookupGroupId(c.Args().First(),
		bucketId); err != nil {
		return err
	}
	path := fmt.Sprintf("/groups/%s", groupId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupRename(c *cli.Context) error {
	uniqKeys := []string{"to", "in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	var (
		err               error
		bucketId, groupId string
	)
	if bucketId, err = adm.LookupBucketId(opts["in"][0]); err != nil {
		return err
	}
	if groupId, err = adm.LookupGroupId(c.Args().First(),
		bucketId); err != nil {
		return err
	}

	var req proto.Request
	req.Group = &proto.Group{}
	req.Group.Name = opts["to"][0]

	path := fmt.Sprintf("/groups/%s", groupId)
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdGroupList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/groups/`, `list`, nil, c)
}

func cmdGroupShow(c *cli.Context) error {
	uniqKeys := []string{"in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	var (
		err               error
		bucketId, groupId string
	)
	if bucketId, err = adm.LookupBucketId(opts[`in`][0]); err != nil {
		return err
	}
	if groupId, err = adm.LookupGroupId(c.Args().First(),
		bucketId); err != nil {
		return err
	}

	path := fmt.Sprintf("/groups/%s", groupId)
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdGroupTree(c *cli.Context) error {
	uniqKeys := []string{"in"}
	opts := make(map[string][]string)

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	var (
		err               error
		bucketId, groupId string
	)
	if bucketId, err = adm.LookupBucketId(opts[`in`][0]); err != nil {
		return err
	}
	if groupId, err = adm.LookupGroupId(c.Args().First(),
		bucketId); err != nil {
		return err
	}

	path := fmt.Sprintf("/groups/%s/tree/tree", groupId)
	return adm.Perform(`get`, path, `tree`, nil, c)
}

func cmdGroupMemberAddGroup(c *cli.Context) error {
	uniqKeys := []string{"to", "in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	var (
		err                         error
		bucketId, mGroupId, groupId string
		req                         proto.Request
		group                       proto.Group
	)
	if bucketId, err = adm.LookupBucketId(
		opts["in"][0]); err != nil {
		return err
	}
	if mGroupId, err = adm.LookupGroupId(c.Args().First(),
		bucketId); err != nil {
		return err
	}
	if groupId, err = adm.LookupGroupId(opts["to"][0],
		bucketId); err != nil {
		return err
	}

	group.Id = mGroupId
	req.Group = &proto.Group{}
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	*req.Group.MemberGroups = append(*req.Group.MemberGroups, group)

	path := fmt.Sprintf("/groups/%s/members/", groupId)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

func cmdGroupMemberAddCluster(c *cli.Context) error {
	uniqKeys := []string{"to", "in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	var (
		err                           error
		bucketId, mClusterId, groupId string
		req                           proto.Request
		cluster                       proto.Cluster
	)
	if bucketId, err = adm.LookupBucketId(
		opts["in"][0]); err != nil {
		return err
	}
	if mClusterId, err = adm.LookupGroupId(c.Args().First(),
		bucketId); err != nil {
		return err
	}
	if groupId, err = adm.LookupClusterId(opts["to"][0],
		bucketId); err != nil {
		return err
	}

	cluster.Id = mClusterId
	req.Group = &proto.Group{}
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	*req.Group.MemberClusters = append(
		*req.Group.MemberClusters, cluster)

	path := fmt.Sprintf("/groups/%s/members/", groupId)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

func cmdGroupMemberAddNode(c *cli.Context) error {
	uniqKeys := []string{"to", "in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	var (
		err                        error
		bucketId, groupId, mNodeId string
		req                        proto.Request
		node                       proto.Node
	)
	if bucketId, err = adm.LookupBucketId(opts[`in`][0]); err != nil {
		return err
	}
	if mNodeId, err = adm.LookupNodeId(c.Args().First()); err != nil {
		return err
	}
	if groupId, err = adm.LookupGroupId(opts[`to`][0],
		bucketId); err != nil {
		return err
	}

	node.Id = mNodeId
	req.Group = &proto.Group{}
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	*req.Group.MemberNodes = append(*req.Group.MemberNodes, node)

	path := fmt.Sprintf("/groups/%s/members/", groupId)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

func cmdGroupMemberDeleteGroup(c *cli.Context) error {
	uniqKeys := []string{"from", "in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	var (
		err                         error
		bucketId, mGroupId, groupId string
	)
	if bucketId, err = adm.LookupBucketId(opts[`in`][0]); err != nil {
		return err
	}
	if mGroupId, err = adm.LookupGroupId(c.Args().First(),
		bucketId); err != nil {
		return err
	}
	if groupId, err = adm.LookupGroupId(opts[`from`][0],
		bucketId); err != nil {
		return err
	}

	path := fmt.Sprintf("/groups/%s/members/%s", groupId,
		mGroupId)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdGroupMemberDeleteCluster(c *cli.Context) error {
	uniqKeys := []string{"from", "in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	var (
		err                           error
		bucketId, mClusterId, groupId string
	)
	if bucketId, err = adm.LookupBucketId(opts[`in`][0]); err != nil {
		return err
	}
	if mClusterId, err = adm.LookupClusterId(c.Args().First(),
		bucketId); err != nil {
		return err
	}
	if groupId, err = adm.LookupGroupId(opts[`from`][0],
		bucketId); err != nil {
		return err
	}

	path := fmt.Sprintf("/groups/%s/members/%s", groupId,
		mClusterId)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdGroupMemberDeleteNode(c *cli.Context) error {
	uniqKeys := []string{"from", "in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	var (
		err                        error
		bucketId, groupId, mNodeId string
	)
	if bucketId, err = adm.LookupBucketId(opts[`in`][0]); err != nil {
		return err
	}
	if mNodeId, err = adm.LookupNodeId(c.Args().First()); err != nil {
		return err
	}
	if groupId, err = adm.LookupGroupId(opts[`from`][0],
		bucketId); err != nil {
		return err
	}

	path := fmt.Sprintf("/groups/%s/members/%s", groupId,
		mNodeId)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdGroupMemberList(c *cli.Context) error {
	uniqKeys := []string{"in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	var (
		err               error
		bucketId, groupId string
	)
	if bucketId, err = adm.LookupBucketId(opts["in"][0]); err != nil {
		return err
	}
	if groupId, err = adm.LookupGroupId(c.Args().First(),
		bucketId); err != nil {
		return err
	}

	path := fmt.Sprintf("/groups/%s/members/", groupId)
	return adm.Perform(`get`, path, `list`, nil, c)
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
	var (
		err               error
		bucketId, groupId string
	)
	if bucketId, err = adm.LookupBucketId(opts["in"][0]); err != nil {
		return err
	}
	if groupId, err = adm.LookupGroupId(opts[`from`][0],
		bucketId); err != nil {
		return err
	}

	if pType == `system` {
		if err := adm.ValidateSystemProperty(
			c.Args().First()); err != nil {
			return err
		}
	}
	var sourceId string
	if err := adm.FindGroupPropSrcId(pType, c.Args().First(),
		opts[`view`][0], groupId, &sourceId); err != nil {
		return err
	}

	req := proto.NewGroupRequest()
	req.Group.Id = groupId
	req.Group.BucketId = bucketId

	path := fmt.Sprintf("/groups/%s/property/%s/%s",
		groupId, pType, sourceId)
	return adm.Perform(`deletebody`, path, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
