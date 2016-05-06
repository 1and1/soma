package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func registerProviders(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// providers
			{
				Name:   "providers",
				Usage:  "SUBCOMMANDS for metric providers",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new metric provider",
						Action: cmdProviderCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a metric provider",
						Action: cmdProviderDelete,
					},
					{
						Name:   "list",
						Usage:  "List metric providers",
						Action: cmdProviderList,
					},
					{
						Name:   "show",
						Usage:  "Show details about a metric provider",
						Action: cmdProviderShow,
					},
				},
			}, // end providers
		}...,
	)
	return &app
}

func cmdProviderCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	req := somaproto.ProtoRequestProvider{}
	req.Provider = &somaproto.ProtoProvider{}
	req.Provider.Provider = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/providers/")
	fmt.Println(resp)
}

func cmdProviderDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/providers/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdProviderList(c *cli.Context) {
	resp := utl.GetRequest("/providers/")
	fmt.Println(resp)
}

func cmdProviderShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/providers/%s", c.Args().First())

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
