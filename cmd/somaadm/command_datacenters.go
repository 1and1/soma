package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
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
						Name:         "rename",
						Usage:        "Rename an existing datacenter",
						Action:       runtime(cmdDatacentersRename),
						BashComplete: cmpl.To,
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
						Name:   "synclist",
						Usage:  "List all datacenters suitable for sync",
						Action: runtime(cmdDatacentersSync),
					},
					/*
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
					*/
				},
			}, // end datacenters
		}...,
	)
	return &app
}

func cmdDatacentersAdd(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	req := proto.NewDatacenterRequest()
	req.Datacenter.Locode = c.Args().First()

	resp := utl.PostRequestWithBody(Client, req, `/datacenters/`)
	fmt.Println(resp)
	return nil
}

func cmdDatacentersRemove(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/datacenters/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdDatacentersRename(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	key := []string{`to`}

	opts := adm.ParseVariadicArguments(key, key, key, c.Args().Tail())

	req := proto.NewDatacenterRequest()
	req.Datacenter.Locode = opts[`to`][0]

	path := fmt.Sprintf("/datacenters/%s", c.Args().First())

	resp := utl.PutRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdDatacentersList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)
	resp := utl.GetRequest(Client, `/datacenters/`)
	fmt.Println(resp)
	return nil
}

func cmdDatacentersSync(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)
	resp := utl.GetRequest(Client, `/sync/datacenters/`)
	fmt.Println(resp)
	return nil
}

func cmdDatacentersShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/datacenters/%s", c.Args().First())
	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// DC:group

func cmdDatacentersAddToGroup(c *cli.Context) error {
	/*
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
	*/
	return nil
}

func cmdDatacentersRemoveFromGroup(c *cli.Context) error {
	/*
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
	*/
	return nil
}

func cmdDatacentersListGroups(c *cli.Context) error {
	/*
		utl.ValidateCliArgumentCount(c, 0)

		resp := utl.GetRequest(Client, "/datacentergroups/")
		fmt.Println(resp)
	*/
	return nil
}

func cmdDatacentersShowGroup(c *cli.Context) error {
	/*
		utl.ValidateCliArgumentCount(c, 1)

		path := fmt.Sprintf("/datacentergroups/%s", c.Args().First())
		resp := utl.GetRequest(Client, path)
		fmt.Println(resp)
	*/
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
