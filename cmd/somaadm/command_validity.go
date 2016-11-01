package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerValidity(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  "validity",
				Usage: "SUBCOMMANDS for system property validity",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new system property validity",
						Action:       runtime(cmdValidityCreate),
						BashComplete: cmpl.ValidityCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a system property validity",
						Action: runtime(cmdValidityDelete),
					},
					{
						Name:   "list",
						Usage:  "List system property validity records",
						Action: runtime(cmdValidityList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a system property validity",
						Action: runtime(cmdValidityShow),
					},
				},
			},
		}...,
	)
	return &app
}

func cmdValidityCreate(c *cli.Context) error {
	multiple := []string{}
	unique := []string{"on", "direct", "inherited"}
	required := []string{"on", "direct", "inherited"}

	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		multiple,
		unique,
		required,
		c.Args().Tail()); err != nil {
		return err
	}

	req := proto.Request{}
	req.Validity = &proto.Validity{
		SystemProperty: c.Args().First(),
		ObjectType:     opts["on"][0],
		Direct:         utl.GetValidatedBool(opts["direct"][0]),
		Inherited:      utl.GetValidatedBool(opts["inherited"][0]),
	}

	resp := utl.PostRequestWithBody(Client, req, "/validity/")
	fmt.Println(resp)
	return nil
}

func cmdValidityDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/validity/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdValidityList(c *cli.Context) error {
	resp := utl.GetRequest(Client, "/validity/")
	fmt.Println(resp)
	return nil
}

func cmdValidityShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/validity/%s", c.Args().First())

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
