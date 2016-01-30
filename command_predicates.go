package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func cmdPredicateCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	req := somaproto.ProtoRequestPredicate{}
	req.Predicate = &somaproto.ProtoPredicate{}
	req.Predicate.Predicate = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/predicates/")
	fmt.Println(resp)
}

func cmdPredicateDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/predicates/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdPredicateList(c *cli.Context) {
	resp := utl.GetRequest("/predicates/")
	fmt.Println(resp)
}

func cmdPredicateShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/predicates/%s", c.Args().First())

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
