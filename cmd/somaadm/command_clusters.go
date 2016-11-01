package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
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
						Name:         "create",
						Usage:        "Create a new cluster",
						Action:       runtime(cmdClusterCreate),
						BashComplete: cmpl.In,
					},
					{
						Name:         "delete",
						Usage:        "Delete a cluster",
						Action:       runtime(cmdClusterDelete),
						BashComplete: cmpl.In,
					},
					{
						Name:         "rename",
						Usage:        "Rename a cluster",
						Action:       runtime(cmdClusterRename),
						BashComplete: cmpl.InTo,
					},
					{
						Name:   "list",
						Usage:  "List all clusters",
						Action: runtime(cmdClusterList),
					},
					{
						Name:         "show",
						Usage:        "Show details about a cluster",
						Action:       runtime(cmdClusterShow),
						BashComplete: cmpl.In,
					},
					{
						Name:         "tree",
						Usage:        "Display the cluster as tree",
						Action:       runtime(cmdClusterTree),
						BashComplete: cmpl.In,
					},
					{
						Name:  "members",
						Usage: "SUBCOMMANDS for cluster members",
						Subcommands: []cli.Command{
							{
								Name:         "add",
								Usage:        "Add a node to a cluster",
								Action:       runtime(cmdClusterMemberAdd),
								BashComplete: cmpl.InTo,
							},
							{
								Name:         "delete",
								Usage:        "Delete a node from a cluster",
								Action:       runtime(cmdClusterMemberDelete),
								BashComplete: cmpl.InFrom,
							},
							{
								Name:         "list",
								Usage:        "List members of a cluster",
								Action:       runtime(cmdClusterMemberList),
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
										Usage:        "Add a system property to a cluster",
										Action:       runtime(cmdClusterSystemPropertyAdd),
										BashComplete: cmpl.PropertyAddValue,
									},
									{
										Name:         "service",
										Usage:        "Add a service property to a cluster",
										Action:       runtime(cmdClusterServicePropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         `oncall`,
										Usage:        `Add an oncall property to a cluster`,
										Action:       runtime(cmdClusterOncallPropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         `custom`,
										Usage:        `Add a custom property to a cluster`,
										Action:       runtime(cmdClusterCustomPropertyAdd),
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
										Usage:        `Delete a system property from a cluster`,
										Action:       runtime(cmdClusterSystemPropertyDelete),
										BashComplete: cmpl.InFromView,
									},
									{
										Name:         `service`,
										Usage:        `Delete a service property from a cluster`,
										Action:       runtime(cmdClusterOncallPropertyDelete),
										BashComplete: cmpl.InFromView,
									},
									{
										Name:         `oncall`,
										Usage:        `Delete an oncall property from a cluster`,
										Action:       runtime(cmdClusterOncallPropertyDelete),
										BashComplete: cmpl.InFromView,
									},
									{
										Name:         `custom`,
										Usage:        `Delete a custom property from a cluster`,
										Action:       runtime(cmdClusterCustomPropertyDelete),
										BashComplete: cmpl.InFromView,
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
	uniqKeys := []string{"in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])

	var req proto.Request
	req.Cluster = &proto.Cluster{}
	req.Cluster.Name = c.Args().First()
	req.Cluster.BucketId = bucketId

	utl.ValidateRuneCountRange(req.Cluster.Name, 4, 256)

	if resp, err := adm.PostReqBody(req, "/clusters/"); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `create`)
	}
}

func cmdClusterDelete(c *cli.Context) error {
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

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client,
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/clusters/%s", clusterId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `create`)
	}
}

func cmdClusterRename(c *cli.Context) error {
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

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
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
		return adm.FormatOut(c, resp, `create`)
	}
}

func cmdClusterList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}
	if resp, err := adm.GetReq(`/clusters/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdClusterShow(c *cli.Context) error {
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

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client,
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/clusters/%s", clusterId)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `show`)
	}
}

func cmdClusterTree(c *cli.Context) error {
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

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client,
		c.Args().First(),
		bucketId)
	path := fmt.Sprintf("/clusters/%s/tree/tree", clusterId)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `tree`)
	}
}

func cmdClusterMemberAdd(c *cli.Context) error {
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

	nodeId := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	//TODO: get bucketId via node
	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
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
	*req.Cluster.Members = append(*req.Cluster.Members, node)

	path := fmt.Sprintf("/clusters/%s/members/", clusterId)

	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdClusterMemberDelete(c *cli.Context) error {
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

	nodeId := utl.TryGetNodeByUUIDOrName(Client, c.Args().First())
	//TODO: get bucketId via node
	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client,
		opts["from"][0], bucketId)

	path := fmt.Sprintf("/clusters/%s/members/%s", clusterId,
		nodeId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdClusterMemberList(c *cli.Context) error {
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

	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client,
		c.Args().First(), bucketId)

	path := fmt.Sprintf("/clusters/%s/members/", clusterId)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

func cmdClusterSystemPropertyAdd(c *cli.Context) error {
	return cmdClusterPropertyAdd(c, `system`)
}

func cmdClusterServicePropertyAdd(c *cli.Context) error {
	return cmdClusterPropertyAdd(c, `service`)
}

func cmdClusterOncallPropertyAdd(c *cli.Context) error {
	return cmdClusterPropertyAdd(c, `oncall`)
}

func cmdClusterCustomPropertyAdd(c *cli.Context) error {
	return cmdClusterPropertyAdd(c, `custom`)
}

func cmdClusterPropertyAdd(c *cli.Context, pType string) error {
	return cmdPropertyAdd(c, pType, `cluster`)
}

func cmdClusterSystemPropertyDelete(c *cli.Context) error {
	return cmdClusterPropertyDelete(c, `system`)
}

func cmdClusterServicePropertyDelete(c *cli.Context) error {
	return cmdClusterPropertyDelete(c, `service`)
}

func cmdClusterOncallPropertyDelete(c *cli.Context) error {
	return cmdClusterPropertyDelete(c, `oncall`)
}

func cmdClusterCustomPropertyDelete(c *cli.Context) error {
	return cmdClusterPropertyDelete(c, `custom`)
}

func cmdClusterPropertyDelete(c *cli.Context, pType string) error {
	multiple := []string{}
	unique := []string{`from`, `view`, `in`}
	required := []string{`from`, `view`, `in`}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		multiple,
		unique,
		required,
		c.Args().Tail()); err != nil {
		return err
	}
	bucketId := utl.BucketByUUIDOrName(Client, opts[`in`][0])
	clusterId := utl.TryGetClusterByUUIDOrName(Client, opts[`from`][0], bucketId)

	if pType == `system` {
		utl.CheckStringIsSystemProperty(Client, c.Args().First())
	}
	sourceId := utl.FindSourceForClusterProperty(Client, pType, c.Args().First(),
		opts[`view`][0], clusterId)
	if sourceId == `` {
		adm.Abort(`Could not find locally set requested property.`)
	}

	req := proto.NewClusterRequest()
	req.Cluster.Id = clusterId
	req.Cluster.BucketId = bucketId
	path := fmt.Sprintf("/clusters/%s/property/%s/%s",
		clusterId, pType, sourceId)

	if resp, err := adm.DeleteReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `delete`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
