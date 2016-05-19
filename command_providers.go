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
				Name:  "providers",
				Usage: "SUBCOMMANDS for metric providers",
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new metric provider",
						Action: runtime(cmdProviderCreate),
					},
					{
						Name:   "delete",
						Usage:  "Delete a metric provider",
						Action: runtime(cmdProviderDelete),
					},
					{
						Name:   "list",
						Usage:  "List metric providers",
						Action: runtime(cmdProviderList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a metric provider",
						Action: runtime(cmdProviderShow),
					},
				},
			}, // end providers
		}...,
	)
	return &app
}

func cmdProviderCreate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	req := proto.Request{}
	req.Provider = &proto.Provider{}
	req.Provider.Name = c.Args().First()

	resp := utl.PostRequestWithBody(Client, req, "/providers/")
	fmt.Println(resp)
	return nil
}

func cmdProviderDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/providers/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdProviderList(c *cli.Context) error {
	resp := utl.GetRequest(Client, "/providers/")
	fmt.Println(resp)
	return nil
}

func cmdProviderShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/providers/%s", c.Args().First())

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
