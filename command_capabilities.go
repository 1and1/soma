package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
)

func registerCapability(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// capability
			{
				Name:   "capabilities",
				Usage:  "SUBCOMMANDS for monitoring capability declarations",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "declare",
						Usage:  "Declare a new monitoring system capability",
						Action: cmdCapabilityDeclare,
					},
					{
						Name:   "revoke",
						Usage:  "Revoke a monitoring system capability",
						Action: cmdCapabilityRevoke,
					},
					{
						Name:   "list",
						Usage:  "List monitoring system capabilities",
						Action: cmdCapabilityList,
					},
					{
						Name:   "show",
						Usage:  "Show details about a monitoring system capability",
						Action: cmdCapabilityShow,
					},
				},
			}, // end capability
		}...,
	)
	return &app
}

func cmdCapabilityDeclare(c *cli.Context) {
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
	utl.AbortOnError(err, "Syntax error, threshold argument not numeric")

	req := somaproto.ProtoRequestCapability{}
	req.Capability = &somaproto.ProtoCapability{}
	req.Capability.Monitoring = utl.TryGetMonitoringByUUIDOrName(
		c.Args().First(),
	)
	req.Capability.Metric = opts["metric"][0]
	req.Capability.View = opts["view"][0]
	req.Capability.Thresholds = thresholds

	resp := utl.PostRequestWithBody(req, "/capability/")
	fmt.Println(resp)
}

func cmdCapabilityRevoke(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	id := utl.TryGetCapabilityByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/capability/%s", id)

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdCapabilityList(c *cli.Context) {
	resp := utl.GetRequest("/capability/")
	fmt.Println(resp)
}

func cmdCapabilityShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	id := utl.TryGetCapabilityByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/capability/%s", id)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
