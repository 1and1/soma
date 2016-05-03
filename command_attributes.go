package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func cmdAttributeCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multiple := []string{}
	unique := []string{"cardinality"}
	required := []string{"cardinality"}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

	switch opts["cardinality"][0] {
	case "once":
	case "multi":
	default:
		utl.Abort("Illegal value for cardinality")
	}

	req := somaproto.AttributeRequest{}
	req.Attribute = &somaproto.Attribute{}
	req.Attribute.Attribute = c.Args().First()
	req.Attribute.Cardinality = opts["cardinality"][0]

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
