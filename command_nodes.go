package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
)

func cmdNodeAdd(c *cli.Context) {
	keySlice := []string{"assetid", "name", "team", "server", "online"}
	reqSlice := []string{"assetid", "name", "team", "server"}

	switch utl.GetCliArgumentCount(c) {
	case 8, 10:
		break
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	argSlice := utl.GetFullArgumentSlice(c)

	options, optional := parseVariableArguments(keySlice, reqSlice, argSlice)
	var req somaproto.ProtoRequestNode

	utl.ValidateStringAsNodeAssetId(options["assetid"])
	if utl.SliceContainsString("online", optional) {
		utl.ValidateStringAsBool(options["online"])
		req.Node.IsOnline, _ = strconv.ParseBool(options["online"])
	} else {
		req.Node.IsOnline = true
	}
	req.Node.AssetId, _ = strconv.ParseUint(options["assetid"], 10, 64)
	req.Node.Name = options["name"]
	req.Node.Team = options["team"]
	req.Node.Server = options["server"]

	_ = utl.PostRequestWithBody(req, "/nodes/")
}

func cmdNodeDel(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id.String())

	_ = utl.DeleteRequest(path)
}

func cmdNodePurge(c *cli.Context) {
	var (
		path string
		req  somaproto.ProtoRequestNode
	)
	if c.Bool("all") {
		utl.ValidateCliArgumentCount(c, 0)
		path = "/nodes/"
	} else {
		utl.ValidateCliArgumentCount(c, 1)
		id := utl.TryGetNodeByUUIDOrName(c.Args().First())
		path = fmt.Sprintf("/nodes/%s", id.String())
	}

	req.Purge = true

	_ = utl.DeleteRequestWithBody(req, path)
}

func cmdNodeRestore(c *cli.Context) {
	var (
		path string
		req  somaproto.ProtoRequestNode
	)
	if c.Bool("all") {
		utl.ValidateCliArgumentCount(c, 0)
		path = "/nodes/"
	} else {
		utl.ValidateCliArgumentCount(c, 1)
		id := utl.TryGetNodeByUUIDOrName(c.Args().First())
		path = fmt.Sprintf("/nodes/%s", id.String())
	}

	req.Restore = true

	_ = utl.DeleteRequestWithBody(req, path)
}

func cmdNodeRename(c *cli.Context) {
}

func cmdNodeRepo(c *cli.Context) {
}

func cmdNodeMove(c *cli.Context) {
}

func cmdNodeOnline(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id.String())

	var req somaproto.ProtoRequestNode
	req.Node.IsOnline = true

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdNodeOffline(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id.String())

	var req somaproto.ProtoRequestNode
	req.Node.IsOnline = false

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdNodeAssign(c *cli.Context) {
}

func cmdNodeList(c *cli.Context) {
}

func cmdNodeShow(c *cli.Context) {
}

func cmdNodePropertyAdd(c *cli.Context) {
}

func cmdNodePropertyGet(c *cli.Context) {
}

func cmdNodePropertyDel(c *cli.Context) {
}

func cmdNodePropertyList(c *cli.Context) {
}

func cmdNodePropertyShow(c *cli.Context) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
