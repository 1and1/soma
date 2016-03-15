package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func cmdValidityCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 7)
	multiple := []string{}
	unique := []string{"on", "direct", "inherited"}
	required := []string{"on", "direct", "inherited"}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

	req := somaproto.ValidityRequest{}
	req.Validity = &somaproto.Validity{
		SystemProperty: c.Args().First(),
		ObjectType:     opts["on"][0],
		Direct:         utl.GetValidatedBool(opts["direct"][0]),
		Inherited:      utl.GetValidatedBool(opts["inherited"][0]),
	}

	resp := utl.PostRequestWithBody(req, "/validity/")
	fmt.Println(resp)
}

func cmdValidityDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/validity/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdValidityList(c *cli.Context) {
	resp := utl.GetRequest("/validity/")
	fmt.Println(resp)
}

func cmdValidityShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/validity/%s", c.Args().First())

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
