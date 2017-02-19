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
										Action:       runtime(cmdClusterServicePropertyDelete),
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

	bucketId, err := adm.LookupBucketId(opts["in"][0])
	if err != nil {
		return err
	}

	var req proto.Request
	req.Cluster = &proto.Cluster{}
	req.Cluster.Name = c.Args().First()
	req.Cluster.BucketId = bucketId

	if err := adm.ValidateRuneCountRange(
		req.Cluster.Name, 4, 256); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/clusters/`, `command`, req, c)
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

	var (
		err                 error
		bucketId, clusterId string
	)
	if bucketId, err = adm.LookupBucketId(opts["in"][0]); err != nil {
		return err
	}
	if clusterId, err = adm.LookupClusterId(c.Args().First(),
		bucketId); err != nil {
		return err
	}

	path := fmt.Sprintf("/clusters/%s", clusterId)
	return adm.Perform(`delete`, path, `command`, nil, c)
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

	var (
		err                 error
		bucketId, clusterId string
		req                 proto.Request
	)
	if bucketId, err = adm.LookupBucketId(opts["in"][0]); err != nil {
		return err
	}
	if clusterId, err = adm.LookupClusterId(c.Args().First(),
		bucketId); err != nil {
		return err
	}

	req.Cluster = &proto.Cluster{}
	req.Cluster.Name = opts["to"][0]

	path := fmt.Sprintf("/clusters/%s", clusterId)
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdClusterList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/clusters/`, `list`, nil, c)
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

	var (
		err                 error
		bucketId, clusterId string
	)
	if bucketId, err = adm.LookupBucketId(opts["in"][0]); err != nil {
		return err
	}
	if clusterId, err = adm.LookupClusterId(c.Args().First(),
		bucketId); err != nil {
		return err
	}

	path := fmt.Sprintf("/clusters/%s", clusterId)
	return adm.Perform(`get`, path, `show`, nil, c)
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

	var (
		err                 error
		bucketId, clusterId string
	)
	if bucketId, err = adm.LookupBucketId(opts["in"][0]); err != nil {
		return err
	}
	if clusterId, err = adm.LookupClusterId(c.Args().First(),
		bucketId); err != nil {
		return err
	}

	path := fmt.Sprintf("/clusters/%s/tree/tree", clusterId)
	return adm.Perform(`get`, path, `tree`, nil, c)
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
	var (
		err                         error
		nodeId, bucketId, clusterId string
	)
	if nodeId, err = adm.LookupNodeId(c.Args().First()); err != nil {
		return err
	}
	//TODO: get bucketId via node
	if bucketId, err = adm.LookupBucketId(opts["in"][0]); err != nil {
		return err
	}
	if clusterId, err = adm.LookupClusterId(opts["to"][0],
		bucketId); err != nil {
		return err
	}

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
	return adm.Perform(`postbody`, path, `command`, req, c)
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

	var (
		err                         error
		nodeId, bucketId, clusterId string
	)
	if nodeId, err = adm.LookupNodeId(c.Args().First()); err != nil {
		return err
	}
	//TODO: get bucketId via node
	if bucketId, err = adm.LookupBucketId(opts["in"][0]); err != nil {
		return err
	}
	if clusterId, err = adm.LookupClusterId(opts["from"][0],
		bucketId); err != nil {
		return err
	}

	path := fmt.Sprintf("/clusters/%s/members/%s", clusterId,
		nodeId)
	return adm.Perform(`delete`, path, `command`, nil, c)
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

	var (
		err                 error
		bucketId, clusterId string
	)
	if bucketId, err = adm.LookupBucketId(opts["in"][0]); err != nil {
		return err
	}
	if clusterId, err = adm.LookupClusterId(c.Args().First(),
		bucketId); err != nil {
		return err
	}

	path := fmt.Sprintf("/clusters/%s/members/", clusterId)
	return adm.Perform(`get`, path, `list`, nil, c)
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
	var (
		err                           error
		bucketId, clusterId, sourceId string
	)
	if bucketId, err = adm.LookupBucketId(opts["in"][0]); err != nil {
		return err
	}
	if clusterId, err = adm.LookupClusterId(opts[`from`][0],
		bucketId); err != nil {
		return err
	}

	if pType == `system` {
		if err := adm.ValidateSystemProperty(
			c.Args().First()); err != nil {
			return err
		}
	}
	if err := adm.FindClusterPropSrcID(pType, c.Args().First(),
		opts[`view`][0], clusterId, &sourceId); err != nil {
		return err
	}

	req := proto.NewClusterRequest()
	req.Cluster.Id = clusterId
	req.Cluster.BucketId = bucketId

	path := fmt.Sprintf("/clusters/%s/property/%s/%s",
		clusterId, pType, sourceId)
	return adm.Perform(`deletebody`, path, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
