package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func cmdStatusCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	var req somaproto.ProtoRequestStatus
	req.Status.Status = c.Args().First()

	_ = utl.PostRequestWithBody(req, "/status/")
}

func cmdStatusDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/status/%s", c.Args().First())

	_ = utl.DeleteRequest(path)
}

func cmdStatusList(c *cli.Context) {
	_ = utl.GetRequest("/status/")
}

func cmdStatusShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/status/%s", c.Args().First())

	_ = utl.GetRequest(path)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
