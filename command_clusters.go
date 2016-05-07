package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func registerClusters(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// clusters
			{
				Name:   "clusters",
				Usage:  "SUBCOMMANDS for clusters",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new cluster",
						Action: cmdClusterCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a cluster",
						Action: cmdClusterDelete,
					},
					{
						Name:   "rename",
						Usage:  "Rename a cluster",
						Action: cmdClusterRename,
					},
					{
						Name:   "list",
						Usage:  "List all clusters",
						Action: cmdClusterList,
					},
					{
						Name:   "show",
						Usage:  "Show details about a cluster",
						Action: cmdClusterShow,
					},
					{
						Name:  "members",
						Usage: "SUBCOMMANDS for cluster members",
						Subcommands: []cli.Command{
							{
								Name:   "add",
								Usage:  "Add a node to a cluster",
								Action: cmdClusterMemberAdd,
							},
							{
								Name:   "delete",
								Usage:  "Delete a node from a cluster",
								Action: cmdClusterMemberDelete,
							},
							{
								Name:   "list",
								Usage:  "List members of a cluster",
								Action: cmdClusterMemberList,
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
										Usage:  "Add a system property to a cluster",
										Action: cmdClusterSystemPropertyAdd,
									},
									{
										Name:   "service",
										Usage:  "Add a service property to a cluster",
										Action: cmdClusterServicePropertyAdd,
									},
								},
							},
						},
					},
				},
			}, // end clusters
		}...,
	)
	return &app
}

func cmdClusterCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])

	var req proto.Request
	req.Cluster = &proto.Cluster{}
	req.Cluster.Name = c.Args().First()
	req.Cluster.BucketId = bucketId

	resp := utl.PostRequestWithBody(req, "/clusters/")
	fmt.Println(resp)
}

func cmdClusterDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/clusters/%s", clusterId)

	_ = utl.DeleteRequest(path)
}

func cmdClusterRename(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/clusters/%s", clusterId)

	var req proto.Request
	req.Cluster = &proto.Cluster{}
	req.Cluster.Name = opts["to"][0]

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdClusterList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)
	/*
			multKeys := []string{"bucket"}
			uniqKeys := []string{}

			opts := utl.ParseVariadicArguments(multKeys,
				uniqKeys,
				uniqKeys,
				c.Args())

			req := somaproto.ProtoRequestCluster{}
			req.Filter = &somaproto.ProtoClusterFilter{}
			req.Filter.BucketId = utl.BucketByUUIDOrName(opts["bucket"][0])
		resp := utl.GetRequestWithBody(req, "/clusters/")
	*/
	resp := utl.GetRequest("/clusters/")
	fmt.Println(resp)
}

func cmdClusterShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/clusters/%s", clusterId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdClusterMemberAdd(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	nodeId := utl.TryGetNodeByUUIDOrName(c.Args().First())
	//TODO: get bucketId via node
	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(
		opts["to"][0], bucketId)

	req := proto.Request{}
	conf := proto.NodeConfig{
		BucketId: bucketId,
	}
	node := proto.Node{
		Id:     nodeId,
		Config: &conf,
	}
	req.Cluster = &proto.Cluster{
		Id:       clusterId,
		BucketId: bucketId,
	}
	req.Cluster.Members = append(req.Cluster.Members, node)

	path := fmt.Sprintf("/clusters/%s/members/", clusterId)

	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdClusterMemberDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"from", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	nodeId := utl.TryGetNodeByUUIDOrName(c.Args().First())
	//TODO: get bucketId via node
	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(
		opts["from"][0], bucketId)

	path := fmt.Sprintf("/clusters/%s/members/%s", clusterId,
		nodeId)

	_ = utl.DeleteRequest(path)
}

func cmdClusterMemberList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(
		c.Args().First(), bucketId)

	path := fmt.Sprintf("/clusters/%s/members/", clusterId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdClusterSystemPropertyAdd(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 9)
	multiple := []string{}
	required := []string{"to", "in", "value", "view"}
	unique := []string{"to", "in", "value", "view", "inheritance", "childrenonly"}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(opts["in"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(opts["to"][0], bucketId)
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

	cluster := proto.Cluster{
		Id:         clusterId,
		BucketId:   bucketId,
		Properties: &propList,
	}

	req := proto.Request{
		Cluster: &cluster,
	}

	path := fmt.Sprintf("/clusters/%s/property/system/", clusterId)
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdClusterServicePropertyAdd(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	required := []string{"to", "in", "view"}
	unique := []string{"to", "in", "view", "inheritance", "childrenonly"}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(opts["in"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(opts["to"][0], bucketId)
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
		Cluster: &proto.Cluster{
			Id:       clusterId,
			BucketId: bucketId,
			Properties: &[]proto.Property{
				tprop,
			},
		},
	}

	path := fmt.Sprintf("/clusters/%s/property/service/", clusterId)
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
