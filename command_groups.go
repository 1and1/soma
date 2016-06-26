package main

import (
	"fmt"

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
						Name:   "create",
						Usage:  "Create a new group",
						Action: runtime(cmdGroupCreate),
					},
					{
						Name:   "delete",
						Usage:  "Delete a group",
						Action: runtime(cmdGroupDelete),
					},
					{
						Name:   "rename",
						Usage:  "Rename a group",
						Action: runtime(cmdGroupRename),
					},
					{
						Name:   "list",
						Usage:  "List all groups",
						Action: runtime(cmdGroupList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a group",
						Action: runtime(cmdGroupShow),
					},
					{
						Name:   `tree`,
						Usage:  `Display the group as tree`,
						Action: runtime(cmdGroupTree),
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
										Name:   "group",
										Usage:  "Add a group to a group",
										Action: runtime(cmdGroupMemberAddGroup),
									},
									{
										Name:   "cluster",
										Usage:  "Add a cluster to a group",
										Action: runtime(cmdGroupMemberAddCluster),
									},
									{
										Name:   "node",
										Usage:  "Add a node to a group",
										Action: runtime(cmdGroupMemberAddNode),
									},
								},
							},
							{
								Name:  "delete",
								Usage: "SUBCOMMANDS for members delete",
								Subcommands: []cli.Command{
									{
										Name:   "group",
										Usage:  "Delete a group from a group",
										Action: runtime(cmdGroupMemberDeleteGroup),
									},
									{
										Name:   "cluster",
										Usage:  "Delete a cluster from a group",
										Action: runtime(cmdGroupMemberDeleteCluster),
									},
									{
										Name:   "node",
										Usage:  "Delete a node from a group",
										Action: runtime(cmdGroupMemberDeleteNode),
									},
								},
							},
							{
								Name:   "list",
								Usage:  "List all members of a group",
								Action: runtime(cmdGroupMemberList),
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
										Name:   "system",
										Usage:  "Add a system property to a group",
										Action: runtime(cmdGroupSystemPropertyAdd),
									},
									{
										Name:   "service",
										Usage:  "Add a service property to a group",
										Action: runtime(cmdGroupServicePropertyAdd),
									},
									{
										Name:   `oncall`,
										Usage:  `Add an oncall property to a group`,
										Action: runtime(cmdGroupOncallPropertyAdd),
									},
								},
							},
							{
								Name:  `delete`,
								Usage: `SUBCOMMANDS for property delete`,
								Subcommands: []cli.Command{
									{
										Name:   `system`,
										Usage:  `Delete a system property from a group`,
										Action: runtime(cmdGroupSystemPropertyDelete),
									},
									{
										Name:   `service`,
										Usage:  `Delete a service property from a group`,
										Action: runtime(cmdGroupOncallPropertyDelete),
									},
									{
										Name:   `oncall`,
										Usage:  `Delete an oncall property from a group`,
										Action: runtime(cmdGroupOncallPropertyDelete),
									},
									{
										Name:   `custom`,
										Usage:  `Delete a custom property from a group`,
										Action: runtime(cmdGroupCustomPropertyDelete),
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
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"in"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

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
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"in"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

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
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "in"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

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
	utl.ValidateCliArgumentCount(c, 0)
	if resp, err := adm.GetReq("/groups/"); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdGroupShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"in"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

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
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"in"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

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
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "in"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

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
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "in"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

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
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "in"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

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
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"from", "in"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

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
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"from", "in"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

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
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"from", "in"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

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
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"in"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

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
	utl.ValidateCliMinArgumentCount(c, 9)
	multiple := []string{}
	required := []string{"to", "in", "value", "view"}
	unique := []string{"to", "in", "value", "view", "inheritance", "childrenonly"}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	groupId := utl.TryGetGroupByUUIDOrName(Client, opts["to"][0], bucketId)
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

	group := proto.Group{
		Id:         groupId,
		BucketId:   bucketId,
		Properties: &propList,
	}

	req := proto.Request{
		Group: &group,
	}

	path := fmt.Sprintf("/groups/%s/property/system/", groupId)
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupServicePropertyAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	required := []string{"to", "in", "view"}
	unique := []string{"to", "in", "view", "inheritance", "childrenonly"}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	groupId := utl.TryGetGroupByUUIDOrName(Client, opts["to"][0], bucketId)
	teamId := utl.TeamIdForBucket(Client, bucketId)

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
		Group: &proto.Group{
			Id:       groupId,
			BucketId: bucketId,
			Properties: &[]proto.Property{
				tprop,
			},
		},
	}

	path := fmt.Sprintf("/groups/%s/property/service/", groupId)
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupOncallPropertyAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	required := []string{"to", "in", "view"}
	unique := []string{"to", "in", "view", "inheritance", "childrenonly"}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	groupId := utl.TryGetGroupByUUIDOrName(Client, opts["to"][0], bucketId)

	oncallId := utl.TryGetOncallByUUIDOrName(Client, c.Args().First())
	oprop := proto.PropertyOncall{
		Id: oncallId,
	}
	oprop.Name, oprop.Number = utl.GetOncallDetailsById(Client, oncallId)

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

	group := proto.Group{
		Id:         groupId,
		BucketId:   bucketId,
		Properties: &propList,
	}

	req := proto.Request{
		Group: &group,
	}

	path := fmt.Sprintf("/groups/%s/property/oncall/", groupId)
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdGroupSystemPropertyDelete(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	unique := []string{`from`, `view`, `in`}
	required := []string{`from`, `view`, `in`}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(Client, opts[`in`][0])
	utl.CheckStringIsSystemProperty(Client, c.Args().First())
	groupId := utl.TryGetGroupByUUIDOrName(Client, opts[`from`][0], bucketId)

	sourceId := utl.FindSourceForGroupProperty(Client, `system`, c.Args().First(),
		opts[`view`][0], groupId)
	if sourceId == `` {
		utl.Abort(`Could not find locally set requested property.`)
	}

	req := proto.NewGroupRequest()
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	path := fmt.Sprintf("/groups/%s/property/%s/%s",
		groupId, `system`, sourceId)

	if resp, err := adm.DeleteReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `delete`)
	}
}

func cmdGroupServicePropertyDelete(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	unique := []string{`from`, `view`, `in`}
	required := []string{`from`, `view`, `in`}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(Client, opts[`in`][0])
	groupId := utl.TryGetGroupByUUIDOrName(Client, opts[`from`][0], bucketId)

	sourceId := utl.FindSourceForGroupProperty(Client, `service`, c.Args().First(),
		opts[`view`][0], groupId)
	if sourceId == `` {
		utl.Abort(`Could not find locally set requested property.`)
	}

	req := proto.NewGroupRequest()
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	path := fmt.Sprintf("/groups/%s/property/%s/%s",
		groupId, `service`, sourceId)

	if resp, err := adm.DeleteReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `delete`)
	}
}

func cmdGroupOncallPropertyDelete(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	unique := []string{`from`, `view`, `in`}
	required := []string{`from`, `view`, `in`}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(Client, opts[`in`][0])
	groupId := utl.TryGetGroupByUUIDOrName(Client, opts[`from`][0], bucketId)

	sourceId := utl.FindSourceForGroupProperty(Client, `oncall`, c.Args().First(),
		opts[`view`][0], groupId)
	if sourceId == `` {
		utl.Abort(`Could not find locally set requested property.`)
	}

	req := proto.NewGroupRequest()
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	path := fmt.Sprintf("/groups/%s/property/%s/%s",
		groupId, `oncall`, sourceId)

	if resp, err := adm.DeleteReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `delete`)
	}
}

func cmdGroupCustomPropertyDelete(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	unique := []string{`from`, `view`, `in`}
	required := []string{`from`, `view`, `in`}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(Client, opts[`in`][0])
	groupId := utl.TryGetGroupByUUIDOrName(Client, opts[`from`][0], bucketId)

	sourceId := utl.FindSourceForGroupProperty(Client, `custom`, c.Args().First(),
		opts[`view`][0], groupId)
	if sourceId == `` {
		utl.Abort(`Could not find locally set requested property.`)
	}

	req := proto.NewGroupRequest()
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	path := fmt.Sprintf("/groups/%s/property/%s/%s",
		groupId, `custom`, sourceId)

	if resp, err := adm.DeleteReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `delete`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
