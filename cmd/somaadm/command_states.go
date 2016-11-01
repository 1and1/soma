package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
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
						Name:         "rename",
						Usage:        "Rename an existing object state",
						Action:       runtime(cmdObjectStatesRename),
						BashComplete: cmpl.To,
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
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewStateRequest()
	req.State.Name = c.Args().First()

	resp := utl.PostRequestWithBody(Client, req, `/objstates/`)
	fmt.Println(resp)
	return nil
}

func cmdObjectStatesRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/objstates/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdObjectStatesRename(c *cli.Context) error {
	key := []string{`to`}

	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(opts, []string{}, key, key,
		c.Args().Tail()); err != nil {
		return err
	}

	req := proto.NewStateRequest()
	req.State.Name = opts[`to`][0]

	path := fmt.Sprintf("/objstates/%s", c.Args().First())

	resp := utl.PutRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdObjectStatesList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}
	resp := utl.GetRequest(Client, "/objstates/")
	fmt.Println(resp)
	return nil
}

func cmdObjectStatesShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/objstates/%s", c.Args().First())

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
