package main

import (
	"fmt"
	"strconv"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/internal/help"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerCapability(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// capability
			{
				Name:  "capabilities",
				Usage: "SUBCOMMANDS for monitoring capability declarations",
				Subcommands: []cli.Command{
					{
						Name:         "declare",
						Usage:        "Declare a new monitoring system capability",
						Description:  help.Text(`CapabilitiesDeclare`),
						Action:       runtime(cmdCapabilityDeclare),
						BashComplete: cmpl.CapabilityDeclare,
					},
					{
						Name:   "revoke",
						Usage:  "Revoke a monitoring system capability",
						Action: runtime(cmdCapabilityRevoke),
					},
					{
						Name:   "list",
						Usage:  "List monitoring system capabilities",
						Action: runtime(cmdCapabilityList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a monitoring system capability",
						Action: runtime(cmdCapabilityShow),
					},
				},
			}, // end capability
		}...,
	)
	return &app
}

func cmdCapabilityDeclare(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 7)
	multiple := []string{}
	unique := []string{"metric", "view", "thresholds"}
	required := []string{"metric", "view", "thresholds"}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

	thresholds, err := strconv.ParseUint(opts["thresholds"][0], 10, 64)
	adm.AbortOnError(err, "Syntax error, threshold argument not numeric")

	req := proto.Request{
		Capability: &proto.Capability{
			MonitoringId: utl.TryGetMonitoringByUUIDOrName(Client,
				c.Args().First(),
			),
			Metric:     opts["metric"][0],
			View:       opts["view"][0],
			Thresholds: thresholds,
		},
	}

	resp := utl.PostRequestWithBody(Client, req, "/capability/")
	fmt.Println(resp)
	return nil
}

func cmdCapabilityRevoke(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	id := utl.TryGetCapabilityByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/capability/%s", id)

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdCapabilityList(c *cli.Context) error {
	resp := utl.GetRequest(Client, "/capability/")
	fmt.Println(resp)
	return nil
}

func cmdCapabilityShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	id := utl.TryGetCapabilityByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/capability/%s", id)

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
