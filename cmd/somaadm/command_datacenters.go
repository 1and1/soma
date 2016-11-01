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
				},
			}, // end datacenters
		}...,
	)
	return &app
}

func cmdDatacentersAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewDatacenterRequest()
	req.Datacenter.Locode = c.Args().First()

	resp := utl.PostRequestWithBody(Client, req, `/datacenters/`)
	fmt.Println(resp)
	return nil
}

func cmdDatacentersRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/datacenters/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdDatacentersRename(c *cli.Context) error {
	key := []string{`to`}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(opts, key, key, key,
		c.Args().Tail()); err != nil {
		return err
	}

	req := proto.NewDatacenterRequest()
	req.Datacenter.Locode = opts[`to`][0]

	path := fmt.Sprintf("/datacenters/%s", c.Args().First())

	resp := utl.PutRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdDatacentersList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}
	resp := utl.GetRequest(Client, `/datacenters/`)
	fmt.Println(resp)
	return nil
}

func cmdDatacentersSync(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}
	resp := utl.GetRequest(Client, `/sync/datacenters/`)
	fmt.Println(resp)
	return nil
}

func cmdDatacentersShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/datacenters/%s", c.Args().First())
	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// DC:group

func cmdDatacentersAddToGroup(c *cli.Context) error {
	return nil
}

func cmdDatacentersRemoveFromGroup(c *cli.Context) error {
	return nil
}

func cmdDatacentersListGroups(c *cli.Context) error {
	return nil
}

func cmdDatacentersShowGroup(c *cli.Context) error {
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
