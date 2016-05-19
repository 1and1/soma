package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/resty.v0"
	"log"
	"net/url"
)

func registerStates(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// states
			{
				Name:  "states",
				Usage: "SUBCOMMANDS for states",
				Subcommands: []cli.Command{
					{
						Name:   "add",
						Usage:  "Add a new object state",
						Action: runtime(cmdObjectStatesAdd),
					},
					{
						Name:   "remove",
						Usage:  "Remove an existing object state",
						Action: runtime(cmdObjectStatesRemove),
					},
					{
						Name:   "rename",
						Usage:  "Rename an existing object state",
						Action: runtime(cmdObjectStatesRename),
					},
					{
						Name:   "list",
						Usage:  "List all object states",
						Action: runtime(cmdObjectStatesList),
					},
					{
						Name:   "show",
						Usage:  "Show information about an object states",
						Action: runtime(cmdObjectStatesShow),
					},
				},
			}, // end states
		}...,
	)
	return &app
}

func cmdObjectStatesAdd(c *cli.Context) error {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}
	url.Path = "/objstates"

	a := c.Args()
	state := a.First()
	if state == "" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: add state [%s]", state)

	req := proto.Request{
		State: &proto.State{
			Name: state,
		},
	}

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

func cmdObjectStatesRemove(c *cli.Context) error {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}

	a := c.Args()
	state := a.First()
	if state == "" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: remove state [%s]", state)
	url.Path = fmt.Sprintf("/objstates/%s", state)

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

func cmdObjectStatesRename(c *cli.Context) error {
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
	log.Printf("Command: rename state [%s] to [%s]", a.Get(0), a.Get(2))

	var req proto.Request
	req.State = &proto.State{}
	req.State.Name = a.Get(2)
	url.Path = fmt.Sprintf("/objstates/%s", a.Get(0))

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

func cmdObjectStatesList(c *cli.Context) error {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}
	url.Path = "/objstates"

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

func cmdObjectStatesShow(c *cli.Context) error {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}

	a := c.Args()
	if len(a) != 1 {
		log.Fatal("Syntax error")
	}
	state := a.First()
	if state == "" {
		log.Fatal("Syntax error")
	}
	url.Path = fmt.Sprintf("/objstates/%s", state)

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
