package main

import (
	"fmt"
	"strings"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/lib/proto"
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
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.Request{}
	req.View = &proto.View{}
	req.View.Name = c.Args().First()
	if strings.Contains(req.View.Name, `.`) {
		adm.Abort(`Views must not contain the character '.'`)
	}

	resp := utl.PostRequestWithBody(Client, req, "/views/")
	fmt.Println(resp)
	return nil
}

func cmdViewsRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/views/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdViewsRename(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	key := []string{"to"}

	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		key,
		key,
		c.Args().Tail()); err != nil {
		return err
	}

	req := proto.Request{}
	req.View = &proto.View{}
	req.View.Name = opts["to"][0]
	path := fmt.Sprintf("/views/%s", c.Args().First())

	resp := utl.PatchRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdViewsList(c *cli.Context) error {
	resp := utl.GetRequest(Client, "/views/")
	fmt.Println(resp)
	return nil
}

func cmdViewsShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/views/%s", c.Args().First())

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
