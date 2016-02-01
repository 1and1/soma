package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func cmdAttributeCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	req := somaproto.ProtoRequestAttribute{}
	req.Attribute = &somaproto.ProtoAttribute{}
	req.Attribute.Attribute = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/attributes/")
	fmt.Println(resp)
}

func cmdAttributeDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/attributes/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdAttributeList(c *cli.Context) {
	resp := utl.GetRequest("/attributes/")
	fmt.Println(resp)
}

func cmdAttributeShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/attributes/%s", c.Args().First())

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
