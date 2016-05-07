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
				Name:   "datacenters",
				Usage:  "SUBCOMMANDS for datacenters",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "add",
						Usage:  "Register a new datacenter",
						Action: cmdDatacentersAdd,
					},
					{
						Name:   "remove",
						Usage:  "Remove an existing datacenter",
						Action: cmdDatacentersRemove,
					},
					{
						Name:   "rename",
						Usage:  "Rename an existing datacenter",
						Action: cmdDatacentersRename,
					},
					{
						Name:   "list",
						Usage:  "List all datacenters",
						Action: cmdDatacentersList,
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific datacenter",
						Action: cmdDatacentersShow,
					},
					{
						Name:   "groupadd",
						Usage:  "Add a datacenter to a datacenter group",
						Action: cmdDatacentersAddToGroup,
					},
					{
						Name:   "groupdel",
						Usage:  "Remove a datacenter from a datacenter group",
						Action: cmdDatacentersRemoveFromGroup,
					},
					{
						Name:   "grouplist",
						Usage:  "List all datacenter groups",
						Action: cmdDatacentersListGroups,
					},
					{
						Name:   "groupshow",
						Usage:  "Show information about a datacenter group",
						Action: cmdDatacentersShowGroup,
					},
				},
			}, // end datacenters
		}...,
	)
	return &app
}

func cmdDatacentersAdd(c *cli.Context) {
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
}

func cmdDatacentersAddToGroup(c *cli.Context) {
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
}

func cmdDatacentersRemove(c *cli.Context) {
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
}

func cmdDatacentersRemoveFromGroup(c *cli.Context) {
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
}

func cmdDatacentersRename(c *cli.Context) {
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
}

func cmdDatacentersList(c *cli.Context) {
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
}

func cmdDatacentersListGroups(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest("/datacentergroups/")
	fmt.Println(resp)
}

func cmdDatacentersShowGroup(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/datacentergroups/%s", c.Args().First())
	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdDatacentersShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/datacenters/%s", c.Args().First())
	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
