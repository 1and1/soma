package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
)

func cmdLevelCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"shortname", "numeric"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	req := somaproto.ProtoRequestLevel{}
	req.Level = &somaproto.ProtoLevel{}
	req.Level.Name = c.Args().First()
	req.Level.ShortName = opts["shortname"][0]
	l, err := strconv.ParseUint(opts["numeric"][0], 10, 16)
	utl.AbortOnError(err, "Syntax error, numeric argument not numeric")
	req.Level.Numeric = uint16(l)

	resp := utl.PostRequestWithBody(req, "/levels/")
	fmt.Println(resp)
}

func cmdLevelDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/levels/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdLevelList(c *cli.Context) {
	resp := utl.GetRequest("/levels/")
	fmt.Println(resp)
}

func cmdLevelShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/levels/%s", c.Args().First())

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
