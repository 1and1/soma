package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
)

func cmdLevelCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 6)
	multKeys := []string{"shortname", "numeric"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	var req somaproto.ProtoRequestLevel
	req.Level.Name = c.Args().First()
	req.Level.ShortName = opts["shortname"][0]
	l, err := strconv.ParseUint(opts["numeric"][0], 10, 16)
	utl.AbortOnError(err, "Syntax error, numeric argument not numeric")
	req.Level.Numeric = uint16(l)

	_ = utl.PostRequestWithBody(req, "/levels/")
}

func cmdLevelDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/levels/%s", c.Args().First())

	_ = utl.DeleteRequest(path)
}

func cmdLevelList(c *cli.Context) {
	_ = utl.GetRequest("/levels/")
}

func cmdLevelShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/levels/%s", c.Args().First())

	_ = utl.GetRequest(path)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
