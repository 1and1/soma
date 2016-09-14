package main

import (
	"fmt"
	"strconv"

	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerLevels(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// levels
			{
				Name:  "levels",
				Usage: "SUBCOMMANDS for notification levels",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new notification level",
						Action:       runtime(cmdLevelCreate),
						BashComplete: cmpl.LevelCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a notification level",
						Action: runtime(cmdLevelDelete),
					},
					{
						Name:   "list",
						Usage:  "List notification levels",
						Action: runtime(cmdLevelList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a notification level",
						Action: runtime(cmdLevelShow),
					},
				},
			}, // end levels
		}...,
	)
	return &app
}

func cmdLevelCreate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"shortname", "numeric"}

	opts := utl.ParseVariadicArguments(multKeys,
		multKeys,
		multKeys,
		c.Args().Tail())

	req := proto.Request{}
	req.Level = &proto.Level{}
	req.Level.Name = c.Args().First()
	req.Level.ShortName = opts["shortname"][0]
	l, err := strconv.ParseUint(opts["numeric"][0], 10, 16)
	utl.AbortOnError(err, "Syntax error, numeric argument not numeric")
	req.Level.Numeric = uint16(l)

	resp := utl.PostRequestWithBody(Client, req, "/levels/")
	fmt.Println(resp)
	return nil
}

func cmdLevelDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/levels/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdLevelList(c *cli.Context) error {
	resp := utl.GetRequest(Client, "/levels/")
	fmt.Println(resp)
	return nil
}

func cmdLevelShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/levels/%s", c.Args().First())

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
