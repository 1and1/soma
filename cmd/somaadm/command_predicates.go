package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerPredicates(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// predicates
			{
				Name:  "predicates",
				Usage: "SUBCOMMANDS for threshold predicates",
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Add a predicate",
						Action: runtime(cmdPredicateCreate),
					},
					{
						Name:   "delete",
						Usage:  "Delete a predicate",
						Action: runtime(cmdPredicateDelete),
					},
					{
						Name:   "list",
						Usage:  "List predicates",
						Action: runtime(cmdPredicateList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a predicate",
						Action: runtime(cmdPredicateShow),
					},
				},
			}, // end predicates
		}...,
	)
	return &app
}

func cmdPredicateCreate(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.Request{}
	req.Predicate = &proto.Predicate{}
	req.Predicate.Symbol = c.Args().First()

	resp := utl.PostRequestWithBody(Client, req, "/predicates/")
	fmt.Println(resp)
	return nil
}

func cmdPredicateDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/predicates/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPredicateList(c *cli.Context) error {
	resp := utl.GetRequest(Client, "/predicates/")
	fmt.Println(resp)
	return nil
}

func cmdPredicateShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/predicates/%s", c.Args().First())

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
