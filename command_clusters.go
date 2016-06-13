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
				Name:  "clusters",
				Usage: "SUBCOMMANDS for clusters",
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new cluster",
						Action: runtime(cmdClusterCreate),
					},
					{
						Name:   "delete",
						Usage:  "Delete a cluster",
						Action: runtime(cmdClusterDelete),
					},
					{
						Name:   "rename",
						Usage:  "Rename a cluster",
						Action: runtime(cmdClusterRename),
					},
					{
						Name:   "list",
						Usage:  "List all clusters",
						Action: runtime(cmdClusterList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a cluster",
						Action: runtime(cmdClusterShow),
					},
					{
						Name:  "members",
						Usage: "SUBCOMMANDS for cluster members",
						Subcommands: []cli.Command{
							{
								Name:   "add",
								Usage:  "Add a node to a cluster",
								Action: runtime(cmdClusterMemberAdd),
							},
							{
								Name:   "delete",
								Usage:  "Delete a node from a cluster",
								Action: runtime(cmdClusterMemberDelete),
							},
							{
								Name:   "list",
								Usage:  "List members of a cluster",
								Action: runtime(cmdClusterMemberList),
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
										Action: runtime(cmdClusterSystemPropertyAdd),
									},
									{
										Name:   "service",
										Usage:  "Add a service property to a cluster",
										Action: runtime(cmdClusterServicePropertyAdd),
									},
									{
										Name:   `oncall`,
										Usage:  `Add an oncall property to a cluster`,
										Action: runtime(cmdClusterOncallPropertyAdd),
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

func cmdClusterCreate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(Client, opts["bucket"][0])

	var req proto.Request
	req.Cluster = &proto.Cluster{}
	req.Cluster.Name = c.Args().First()
	req.Cluster.BucketId = bucketId

	if resp, err := adm.PostReqBody(req, "/clusters/"); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdClusterDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(Client, opts["bucket"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client,
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/clusters/%s", clusterId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdClusterRename(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(Client, opts["bucket"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client,
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/clusters/%s", clusterId)

	var req proto.Request
	req.Cluster = &proto.Cluster{}
	req.Cluster.Name = opts["to"][0]

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdClusterList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)
	if resp, err := adm.GetReq(`/clusters/`); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdClusterShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(Client, opts["bucket"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client,
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/clusters/%s", clusterId)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdClusterMemberAdd(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"to", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	nodeId := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	//TODO: get bucketId via node
	bucketId := utl.BucketByUUIDOrName(Client, opts["bucket"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client,
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

	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdClusterMemberDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"from", "bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	nodeId := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	//TODO: get bucketId via node
	bucketId := utl.BucketByUUIDOrName(Client, opts["bucket"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client,
		opts["from"][0], bucketId)

	path := fmt.Sprintf("/clusters/%s/members/%s", clusterId,
		nodeId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdClusterMemberList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(Client, opts["bucket"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client,
		c.Args().First(), bucketId)

	path := fmt.Sprintf("/clusters/%s/members/", clusterId)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdClusterSystemPropertyAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 9)
	multiple := []string{}
	required := []string{"to", "in", "value", "view"}
	unique := []string{"to", "in", "value", "view", "inheritance", "childrenonly"}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client, opts["to"][0], bucketId)
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

	cluster := proto.Cluster{
		Id:         clusterId,
		BucketId:   bucketId,
		Properties: &propList,
	}

	req := proto.Request{
		Cluster: &cluster,
	}

	path := fmt.Sprintf("/clusters/%s/property/system/", clusterId)
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdClusterServicePropertyAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	required := []string{"to", "in", "view"}
	unique := []string{"to", "in", "view", "inheritance", "childrenonly"}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client, opts["to"][0], bucketId)
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
		Cluster: &proto.Cluster{
			Id:       clusterId,
			BucketId: bucketId,
			Properties: &[]proto.Property{
				tprop,
			},
		},
	}

	path := fmt.Sprintf("/clusters/%s/property/service/", clusterId)
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdClusterOncallPropertyAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	required := []string{"to", "in", "view"}
	unique := []string{"to", "in", "view", "inheritance", "childrenonly"}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client, opts["to"][0], bucketId)

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

	cluster := proto.Cluster{
		Id:         clusterId,
		BucketId:   bucketId,
		Properties: &propList,
	}

	req := proto.Request{
		Cluster: &cluster,
	}

	path := fmt.Sprintf("/clusters/%s/property/oncall/", clusterId)
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
