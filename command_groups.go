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
				Name:   "groups",
				Usage:  "SUBCOMMANDS for groups",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new group",
						Action: cmdGroupCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a group",
						Action: cmdGroupDelete,
					},
					{
						Name:   "rename",
						Usage:  "Rename a group",
						Action: cmdGroupRename,
					},
					{
						Name:   "list",
						Usage:  "List all groups",
						Action: cmdGroupList,
					},
					{
						Name:   "show",
						Usage:  "Show details about a group",
						Action: cmdGroupShow,
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
										Action: cmdGroupMemberAddGroup,
									},
									{
										Name:   "cluster",
										Usage:  "Add a cluster to a group",
										Action: cmdGroupMemberAddCluster,
									},
									{
										Name:   "node",
										Usage:  "Add a node to a group",
										Action: cmdGroupMemberAddNode,
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
										Action: cmdGroupMemberDeleteGroup,
									},
									{
										Name:   "cluster",
										Usage:  "Delete a cluster from a group",
										Action: cmdGroupMemberDeleteCluster,
									},
									{
										Name:   "node",
										Usage:  "Delete a node from a group",
										Action: cmdGroupMemberDeleteNode,
									},
								},
							},
							{
								Name:   "list",
								Usage:  "List all members of a group",
								Action: cmdGroupMemberList,
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
										Action: cmdGroupSystemPropertyAdd,
									},
									{
										Name:   "service",
										Usage:  "Add a service property to a group",
										Action: cmdGroupServicePropertyAdd,
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

func cmdGroupCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])

	var req proto.Request
	req.Group = &proto.Group{}
	req.Group.Name = c.Args().First()
	req.Group.BucketId = bucketId

	resp := utl.PostRequestWithBody(req, "/groups/")
	fmt.Println(resp)
}

func cmdGroupDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	groupId := utl.TryGetGroupByUUIDOrName(
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/groups/%s", groupId)

	_ = utl.DeleteRequest(path)
}

func cmdGroupRename(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	groupId := utl.TryGetGroupByUUIDOrName(
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/groups/%s", groupId)

	var req proto.Request
	req.Group = &proto.Group{}
	req.Group.Name = opts["to"][0]

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdGroupList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)
	resp := utl.GetRequest("/groups/")
	fmt.Println(resp)
}

func cmdGroupShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	groupId := utl.TryGetGroupByUUIDOrName(
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/groups/%s", groupId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdGroupMemberAddGroup(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	mGroupId := utl.TryGetGroupByUUIDOrName(
		c.Args().First(),
		bucketId)
	groupId := utl.TryGetGroupByUUIDOrName(
		opts["to"][0],
		bucketId)

	var req proto.Request
	var group proto.Group
	group.Id = mGroupId
	req.Group = &proto.Group{}
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	req.Group.MemberGroups = append(req.Group.MemberGroups, group)

	path := fmt.Sprintf("/groups/%s/members/", groupId)

	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdGroupMemberAddCluster(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	mClusterId := utl.TryGetClusterByUUIDOrName(
		c.Args().First(),
		bucketId)
	groupId := utl.TryGetGroupByUUIDOrName(
		opts["to"][0],
		bucketId)

	var req proto.Request
	var cluster proto.Cluster
	cluster.Id = mClusterId
	req.Group = &proto.Group{}
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	req.Group.MemberClusters = append(req.Group.MemberClusters, cluster)

	path := fmt.Sprintf("/groups/%s/members/", groupId)

	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdGroupMemberAddNode(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	mNodeId := utl.TryGetNodeByUUIDOrName(c.Args().First())
	groupId := utl.TryGetGroupByUUIDOrName(
		opts["to"][0],
		bucketId)

	var req proto.Request
	var node proto.Node
	node.Id = mNodeId
	req.Group = &proto.Group{}
	req.Group.Id = groupId
	req.Group.BucketId = bucketId
	req.Group.MemberNodes = append(req.Group.MemberNodes, node)

	path := fmt.Sprintf("/groups/%s/members/", groupId)

	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdGroupMemberDeleteGroup(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"from", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	mGroupId := utl.TryGetGroupByUUIDOrName(
		c.Args().First(),
		bucketId)
	groupId := utl.TryGetGroupByUUIDOrName(
		opts["from"][0],
		bucketId)

	path := fmt.Sprintf("/groups/%s/members/%s", groupId,
		mGroupId)

	_ = utl.DeleteRequest(path)
}

func cmdGroupMemberDeleteCluster(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"from", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	mClusterId := utl.TryGetClusterByUUIDOrName(
		c.Args().First(),
		bucketId)
	groupId := utl.TryGetGroupByUUIDOrName(
		opts["from"][0],
		bucketId)

	path := fmt.Sprintf("/groups/%s/members/%s", groupId,
		mClusterId)

	_ = utl.DeleteRequest(path)
}

func cmdGroupMemberDeleteNode(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"from", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	mNodeId := utl.TryGetNodeByUUIDOrName(c.Args().First())
	groupId := utl.TryGetGroupByUUIDOrName(
		opts["from"][0],
		bucketId)

	path := fmt.Sprintf("/groups/%s/members/%s", groupId,
		mNodeId)

	_ = utl.DeleteRequest(path)
}

func cmdGroupMemberList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	groupId := utl.TryGetGroupByUUIDOrName(
		c.Args().First(),
		bucketId)

	path := fmt.Sprintf("/groups/%s/members/", groupId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdGroupSystemPropertyAdd(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 9)
	multiple := []string{}
	required := []string{"to", "in", "value", "view"}
	unique := []string{"to", "in", "value", "view", "inheritance", "childrenonly"}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(opts["in"][0])
	groupId := utl.TryGetGroupByUUIDOrName(opts["to"][0], bucketId)
	utl.CheckStringIsSystemProperty(c.Args().First())

	sprop := proto.PropertySystem{
		Name:  c.Args().First(),
		Value: opts["value"][0],
	}

	tprop := proto.Property{
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
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdGroupServicePropertyAdd(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	required := []string{"to", "in", "view"}
	unique := []string{"to", "in", "view", "inheritance", "childrenonly"}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(opts["in"][0])
	groupId := utl.TryGetGroupByUUIDOrName(opts["to"][0], bucketId)
	teamId := utl.TeamIdForBucket(bucketId)

	// no reason to fill out the attributes, client-provided
	// attributes are discarded by the server
	tprop := proto.Property{
		PropertyType: "service",
		View:         opts["view"][0],
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
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
