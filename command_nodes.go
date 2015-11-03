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

	options, optional := utl.ParseVariableArguments(keySlice, reqSlice, argSlice)
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
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id.String())

	var req somaproto.ProtoRequestNode
	req.Node.Name = c.Args().Get(2)

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdNodeRepo(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetNodeByUUIDOrName(c.Args().Get(0))
	team := c.Args().Get(2)
	// try resolving team name to uuid as name validation
	_ = utl.GetTeamIdByName(team)
	path := fmt.Sprintf("/nodes/%s", id.String())

	var req somaproto.ProtoRequestNode
	req.Node.Team = team

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdNodeMove(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetNodeByUUIDOrName(c.Args().Get(0))
	server := c.Args().Get(2)
	// try resolving server name to uuid as name validation
	_ = utl.GetServerAssetIdByName(server)
	path := fmt.Sprintf("/nodes/%s", id.String())

	var req somaproto.ProtoRequestNode
	req.Node.Server = server

	_ = utl.PatchRequestWithBody(req, path)
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
	utl.ValidateCliArgumentCount(c, 6)
	utl.ValidateCliArgument(c, 2, "to")
	keySlice := []string{"repository", "bucket"}
	argSlice := utl.GetFullArgumentSlice(c)[2:]

	options, _ := utl.ParseVariableArguments(keySlice, keySlice, argSlice)
	var req somaproto.ProtoRequestJob
	req.JobType = "node"
	req.Node.Action = "assign"
	req.Node.Node.Config.RepositoryName = options["repository"]
	req.Node.Node.Config.BucketName = options["bucket"]

	_ = utl.PostRequestWithBody(req, "/jobs/")
	// TODO save jobid locally as outstanding
}

func cmdNodeList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)

	_ = utl.GetRequest("/nodes/")
}

func cmdNodeShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s", id.String())

	_ = utl.GetRequest(path)
}

func cmdNodePropertyAdd(c *cli.Context) {
	// preliminary argv validation
	switch utl.GetCliArgumentCount(c) {
	case 4, 6, 8, 10, 12:
		break
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	utl.ValidateCliArgument(c, 3, "to")
	argSlice := utl.GetFullArgumentSlice(c)

	// get property types
	typSlice := utl.PropertyTypes

	// first argument must be a valid property type
	utl.ValidateStringInSlice(argSlice[0], typSlice)

	// TODO: validate property of that type and name exists
	propertyType := argSlice[0]
	property := argSlice[1]

	// get which node is being modified
	id := utl.TryGetNodeByUUIDOrName(argSlice[3])
	path := fmt.Sprintf("/nodes/%s/property/", id.String())

	// variable key/value part of argv
	argSlice = argSlice[4:]

	// define accepted and required keys
	keySlice := []string{"inheritance", "childrenonly", "view", "value"}
	reqSlice := []string{"view"}
	if propertyType != "service" {
		// non service properties require values, services are
		// predefined and do not
		reqSlice = append(reqSlice, "value")
	}
	options, optional := utl.ParseVariableArguments(keySlice, reqSlice, argSlice)

	// build node property JSON
	var prop somaproto.ProtoNodeProperty
	prop.Type = propertyType
	prop.View = options["view"] //required
	prop.Property = property
	// add value if it was required
	if utl.SliceContainsString("value", reqSlice) {
		prop.Value = options["value"]
	}

	// optional inheritance, default true
	if utl.SliceContainsString("inheritance", optional) {
		utl.ValidateStringAsBool(options["inheritance"])
		prop.Inheritance, _ = strconv.ParseBool(options["inheritance"])
	} else {
		prop.Inheritance = true
	}

	// optional childrenonly, default false
	if utl.SliceContainsString("childrenonly", optional) {
		utl.ValidateStringAsBool(options["childrenonly"])
		prop.ChildrenOnly, _ = strconv.ParseBool(options["childrenonly"])
	} else {
		prop.ChildrenOnly = false
	}

	// build request JSON
	var req somaproto.ProtoRequestNode
	req.Node.Properties = append(req.Node.Properties, prop)

	_ = utl.PostRequestWithBody(req, path)
	// TODO save jobid locally as outstanding
}

func cmdNodePropertyGet(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 6)
	utl.ValidateCliArgument(c, 3, "from")
	utl.ValidateCliArgument(c, 5, "view")
	argSlice := utl.GetFullArgumentSlice(c)

	// first argument must be a valid property type, sixth a view
	utl.ValidateStringInSlice(argSlice[0], utl.PropertyTypes)
	utl.ValidateStringInSlice(argSlice[5], utl.Views)

	// get which node is being modified
	id := utl.TryGetNodeByUUIDOrName(argSlice[3])

	// TODO: validate property of that type and name exists
	path := fmt.Sprintf("/nodes/%s/property/%s/%s/%s",
		id.String(),
		argSlice[0], // type
		argSlice[5], // view
		argSlice[1], // property
	)

	_ = utl.GetRequest(path)
}

func cmdNodePropertyDel(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 6)
	utl.ValidateCliArgument(c, 3, "from")
	utl.ValidateCliArgument(c, 5, "view")
	argSlice := utl.GetFullArgumentSlice(c)

	// first argument must be a valid property type, sixth a view
	utl.ValidateStringInSlice(argSlice[0], utl.PropertyTypes)
	utl.ValidateStringInSlice(argSlice[5], utl.Views)

	propertyType := argSlice[0]
	property := argSlice[1]

	// get which node is being modified
	id := utl.TryGetNodeByUUIDOrName(argSlice[3])

	// TODO: validate property of that type and name exists
	path := fmt.Sprintf("/nodes/%s/property/%s/%s/%s",
		id.String(),
		argSlice[0], //type
		argSlice[5], //view
		argSlice[1], //property
	)

	_ = DeleteRequest(path)
}

func cmdNodePropertyList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s/property/", id.String())

	var req somaproto.ProtoRequestNode

	if c.Bool("all") {
		req.Filter.LocalProperty = false
		_ = GetRequest(path)
	} else {
		req.Filter.LocalProperty = true
		_ = GetRequestWithBody(req, path)
	}
}

func cmdNodePropertyShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetNodeByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/nodes/%s/property/", id.String())

	_ = GetRequest(path)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
