package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func cmdViewsAdd(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	var req somaproto.ProtoRequestView
	req.View.View = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/views/")
	fmt.Println(resp)
}

func cmdViewsRemove(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/views/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdViewsRename(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	key := []string{"to"}

	opts := utl.ParseVariadicArguments(
		key,
		key,
		key,
		c.Args().Tail())

	var req somaproto.ProtoRequestView
	req.View = &somaproto.ProtoView{}
	req.View.View = opts["to"][0]
	path := fmt.Sprintf("/views/%s", c.Args().First())

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdViewsList(c *cli.Context) {
	resp := utl.GetRequest("/views/")
	fmt.Println(resp)
}

func cmdViewsShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/views/%s", c.Args().First())

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
