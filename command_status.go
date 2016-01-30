package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func cmdStatusCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	req := somaproto.ProtoRequestStatus{}
	req.Status = &somaproto.ProtoStatus{}
	req.Status.Status = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/status/")
	fmt.Println(resp)
}

func cmdStatusDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/status/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdStatusList(c *cli.Context) {
	resp := utl.GetRequest("/status/")
	fmt.Println(resp)
}

func cmdStatusShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/status/%s", c.Args().First())

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
