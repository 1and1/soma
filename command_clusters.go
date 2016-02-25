package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func cmdClusterCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multKeys := []string{"bucket"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys, // as uniqKeys
		multKeys, // as reqKeys
		c.Args().Tail())

	bucketId := utl.BucketByUUIDOrName(opts["bucket"][0])

	var req somaproto.ProtoRequestCluster
	req.Cluster = &somaproto.ProtoCluster{}
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

	var req somaproto.ProtoRequestCluster
	req.Cluster = &somaproto.ProtoCluster{}
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

	req := somaproto.ProtoRequestCluster{}
	conf := somaproto.ProtoNodeConfig{
		BucketId: bucketId,
	}
	node := somaproto.ProtoNode{
		Id:     nodeId,
		Config: &conf,
	}
	req.Cluster = &somaproto.ProtoCluster{
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
