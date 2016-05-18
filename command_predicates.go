package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func registerPredicates(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// predicates
			{
				Name:   "predicates",
				Usage:  "SUBCOMMANDS for threshold predicates",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Add a predicate",
						Action: cmdPredicateCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a predicate",
						Action: cmdPredicateDelete,
					},
					{
						Name:   "list",
						Usage:  "List predicates",
						Action: cmdPredicateList,
					},
					{
						Name:   "show",
						Usage:  "Show details about a predicate",
						Action: cmdPredicateShow,
					},
				},
			}, // end predicates
		}...,
	)
	return &app
}

func cmdPredicateCreate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	req := proto.Request{}
	req.Predicate = &proto.Predicate{}
	req.Predicate.Symbol = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/predicates/")
	fmt.Println(resp)
	return nil
}

func cmdPredicateDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/predicates/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
	return nil
}

func cmdPredicateList(c *cli.Context) error {
	resp := utl.GetRequest("/predicates/")
	fmt.Println(resp)
	return nil
}

func cmdPredicateShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/predicates/%s", c.Args().First())

	resp := utl.GetRequest(path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
