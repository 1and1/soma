package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/resty.v0"
	"log"
	"net/url"
)

func registerTypes(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// types
			{
				Name:  "types",
				Usage: "SUBCOMMANDS for object types",
				Subcommands: []cli.Command{
					{
						Name:   "add",
						Usage:  "Add a new object type",
						Action: runtime(cmdObjectTypesAdd),
					},
					{
						Name:   "remove",
						Usage:  "Remove an existing object type",
						Action: runtime(cmdObjectTypesRemove),
					},
					{
						Name:   "rename",
						Usage:  "Rename an existing object type",
						Action: runtime(cmdObjectTypesRename),
					},
					{
						Name:   "list",
						Usage:  "List all object types",
						Action: runtime(cmdObjectTypesList),
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific object type",
						Action: runtime(cmdObjectTypesShow),
					},
				},
			}, // end types
		}...,
	)
	return &app
}

func cmdObjectTypesAdd(c *cli.Context) error {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}
	url.Path = "/objtypes"

	a := c.Args()
	objectType := a.First()
	if objectType == "" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: add objectType [%s]", objectType)

	var req proto.Request
	req.Entity = &proto.Entity{}
	req.Entity.Name = objectType

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Post(url.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s\n", resp.Status())
	return nil
}

func cmdObjectTypesRemove(c *cli.Context) error {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}

	a := c.Args()
	objectType := a.First()
	if objectType == "" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: remove objectType [%s]", objectType)
	url.Path = fmt.Sprintf("/objtypes/%s", objectType)

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Delete(url.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s\n", resp.Status())
	return nil
}

func cmdObjectTypesRename(c *cli.Context) error {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}

	a := c.Args()
	// we expected exactly 3 arguments
	if len(a) != 3 {
		log.Fatal("Syntax error")
	}
	// second arg must be `to`
	if a.Get(1) != "to" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: rename objectType [%s] to [%s]", a.Get(0), a.Get(2))

	var req proto.Request
	req.Entity = &proto.Entity{}
	req.Entity.Name = a.Get(2)
	url.Path = fmt.Sprintf("/objtypes/%s", a.Get(0))

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Put(url.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s\n", resp.Status())
	return nil
}

func cmdObjectTypesList(c *cli.Context) error {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}
	url.Path = "/objtypes"

	a := c.Args()
	if len(a) != 0 {
		log.Fatal("Syntax error")
	}

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s\n", resp.Status())
	return nil
}

func cmdObjectTypesShow(c *cli.Context) error {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}

	a := c.Args()
	if len(a) != 1 {
		log.Fatal("Syntax error")
	}
	objectType := a.First()
	if objectType == "" {
		log.Fatal("Syntax error")
	}
	url.Path = fmt.Sprintf("/objtypes/%s", objectType)

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s\n", resp.Status())
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
