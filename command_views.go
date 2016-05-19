package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func registerViews(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// views
			{
				Name:  "views",
				Usage: "SUBCOMMANDS for views",
				Subcommands: []cli.Command{
					{
						Name:   "add",
						Usage:  "Register a new view",
						Action: runtime(cmdViewsAdd),
					},
					{
						Name:   "remove",
						Usage:  "Remove an existing view",
						Action: runtime(cmdViewsRemove),
					},
					{
						Name:   "rename",
						Usage:  "Rename an existing view",
						Action: runtime(cmdViewsRename),
					},
					{
						Name:   "list",
						Usage:  "List all registered views",
						Action: runtime(cmdViewsList),
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific view",
						Action: runtime(cmdViewsShow),
					},
				},
			}, // end views
		}...,
	)
	return &app
}

func cmdViewsAdd(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	req := proto.Request{}
	req.View = &proto.View{}
	req.View.Name = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/views/")
	fmt.Println(resp)
	return nil
}

func cmdViewsRemove(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/views/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
	return nil
}

func cmdViewsRename(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	key := []string{"to"}

	opts := utl.ParseVariadicArguments(
		key,
		key,
		key,
		c.Args().Tail())

	req := proto.Request{}
	req.View = &proto.View{}
	req.View.Name = opts["to"][0]
	path := fmt.Sprintf("/views/%s", c.Args().First())

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
	return nil
}

func cmdViewsList(c *cli.Context) error {
	resp := utl.GetRequest("/views/")
	fmt.Println(resp)
	return nil
}

func cmdViewsShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/views/%s", c.Args().First())

	resp := utl.GetRequest(path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
