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
}

func cmdClusterShow(c *cli.Context) {
}

func cmdClusterMemberAdd(c *cli.Context) {
}

func cmdClusterMemberDelete(c *cli.Context) {
}

func cmdClusterMemberList(c *cli.Context) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
