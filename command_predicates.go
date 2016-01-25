package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func cmdPredicateCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	var req somaproto.ProtoRequestPredicate
	req.Predicate.Predicate = c.Args().First()

	_ = utl.PostRequestWithBody(req, "/predicates/")
}

func cmdPredicateDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/predicates/%s", c.Args().First())

	_ = utl.DeleteRequest(path)
}

func cmdPredicateList(c *cli.Context) {
	_ = utl.GetRequest("/predicates/")
}

func cmdPredicateShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/predicates/%s", c.Args().First())

	_ = utl.GetRequest(path)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
