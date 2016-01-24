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
	req.Cluster.Name = c.Args().First()
	req.Cluster.BucketId = bucketId

	_ = utl.PostRequestWithBody(req, "/clusters/")
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
	req.Cluster.Name = opts["to"][0]

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdClusterList(c *cli.Context) {
	multKeys := []string{"bucket"}
	uniqKeys := []string{}

	opts := utl.ParseVariadicArguments(multKeys,
		uniqKeys,
		uniqKeys,
		c.Args().Tail())

	var req somaproto.ProtoRequestCluster
	req.Filter.BucketId = utl.BucketByUUIDOrName(opts["bucket"][0])
	_ = utl.GetRequestWithBody(req, "/clusters/")
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

	_ = utl.GetRequest(path)
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

	var req somaproto.ProtoRequestCluster
	var node somaproto.ProtoNode
	node.Id = nodeId
	req.Cluster.Members = append(req.Cluster.Members, node)

	path := fmt.Sprintf("/clusters/%s/members", clusterId)

	_ = utl.PostRequestWithBody(req, path)
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
		nodeId.String())

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

	_ = utl.GetRequest(path)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
