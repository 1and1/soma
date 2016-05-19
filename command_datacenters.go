package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/codegangsta/cli"
	"gopkg.in/resty.v0"
)

func registerDatacenters(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// datacenters
			{
				Name:  "datacenters",
				Usage: "SUBCOMMANDS for datacenters",
				Subcommands: []cli.Command{
					{
						Name:   "add",
						Usage:  "Register a new datacenter",
						Action: runtime(cmdDatacentersAdd),
					},
					{
						Name:   "remove",
						Usage:  "Remove an existing datacenter",
						Action: runtime(cmdDatacentersRemove),
					},
					{
						Name:   "rename",
						Usage:  "Rename an existing datacenter",
						Action: runtime(cmdDatacentersRename),
					},
					{
						Name:   "list",
						Usage:  "List all datacenters",
						Action: runtime(cmdDatacentersList),
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific datacenter",
						Action: runtime(cmdDatacentersShow),
					},
					{
						Name:   "groupadd",
						Usage:  "Add a datacenter to a datacenter group",
						Action: runtime(cmdDatacentersAddToGroup),
					},
					{
						Name:   "groupdel",
						Usage:  "Remove a datacenter from a datacenter group",
						Action: runtime(cmdDatacentersRemoveFromGroup),
					},
					{
						Name:   "grouplist",
						Usage:  "List all datacenter groups",
						Action: runtime(cmdDatacentersListGroups),
					},
					{
						Name:   "groupshow",
						Usage:  "Show information about a datacenter group",
						Action: runtime(cmdDatacentersShowGroup),
					},
				},
			}, // end datacenters
		}...,
	)
	return &app
}

func cmdDatacentersAdd(c *cli.Context) error {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}
	url.Path = "/datacenters"

	a := c.Args()
	datacenter := a.First()
	if datacenter == "" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: add datacenter [%s]", datacenter)

	var req proto.Request
	req.Datacenter = &proto.Datacenter{}
	req.Datacenter.Locode = datacenter

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

func cmdDatacentersAddToGroup(c *cli.Context) error {
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
	if a.Get(1) != "group" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: add datacenter [%s] to group [%s]", a.Get(0), a.Get(2))

	var req proto.Request
	req.Datacenter = &proto.Datacenter{}
	req.Datacenter.Locode = a.Get(0)
	url.Path = fmt.Sprintf("/datacentergroups/%s", a.Get(2))

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s", resp.Status())
	return nil
}

func cmdDatacentersRemove(c *cli.Context) error {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}

	a := c.Args()
	datacenter := a.First()
	if datacenter == "" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: remove datacenter [%s]", datacenter)
	url.Path = fmt.Sprintf("/datacenters/%s", datacenter)

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

func cmdDatacentersRemoveFromGroup(c *cli.Context) error {
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
	if a.Get(1) != "group" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: remove datacenter [%s] from group [%s]", a.Get(0), a.Get(2))

	var req proto.Request
	req.Datacenter = &proto.Datacenter{}
	req.Datacenter.Locode = a.Get(0)
	url.Path = fmt.Sprintf("/datacentergroups/%s", a.Get(2))

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Delete(url.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s", resp.Status())
	return nil
}

func cmdDatacentersRename(c *cli.Context) error {
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
	log.Printf("Command: rename datacenter [%s] to [%s]", a.Get(0), a.Get(2))

	var req proto.Request
	req.Datacenter = &proto.Datacenter{}
	req.Datacenter.Locode = a.Get(2)
	url.Path = fmt.Sprintf("/datacenters/%s", a.Get(0))

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

func cmdDatacentersList(c *cli.Context) error {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}
	url.Path = "/datacenters"

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

func cmdDatacentersListGroups(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest("/datacentergroups/")
	fmt.Println(resp)
	return nil
}

func cmdDatacentersShowGroup(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/datacentergroups/%s", c.Args().First())
	resp := utl.GetRequest(path)
	fmt.Println(resp)
	return nil
}

func cmdDatacentersShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/datacenters/%s", c.Args().First())
	resp := utl.GetRequest(path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
