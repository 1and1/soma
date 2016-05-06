package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func registerMonitoring(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// monitoring
			{
				Name:   "monitoring",
				Usage:  "SUBCOMMANDS for monitoring systems",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new monitoring system",
						Action: cmdMonitoringCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a monitoring system",
						Action: cmdMonitoringDelete,
					},
					{
						Name:   "list",
						Usage:  "List monitoring systems",
						Action: cmdMonitoringList,
					},
					{
						Name:   "show",
						Usage:  "Show details about a monitoring system",
						Action: cmdMonitoringShow,
					},
				},
			}, // end monitoring
		}...,
	)
	return &app
}

func cmdMonitoringCreate(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	unique := []string{"mode", "contact", "team", "callback"}
	required := []string{"mode", "contact", "team"}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

	req := somaproto.ProtoRequestMonitoring{}
	req.Monitoring = &somaproto.ProtoMonitoring{}
	req.Monitoring.Name = c.Args().First()
	req.Monitoring.Mode = opts["mode"][0]
	req.Monitoring.Contact = utl.TryGetUserByUUIDOrName(opts["contact"][0])
	req.Monitoring.Team = utl.TryGetTeamByUUIDOrName(opts["team"][0])

	// optional arguments
	if _, ok := opts["callback"]; ok {
		req.Monitoring.Callback = opts["callback"][0]
	}

	resp := utl.PostRequestWithBody(req, "/monitoring/")
	fmt.Println(resp)
}

func cmdMonitoringDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	userId := utl.TryGetMonitoringByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/monitoring/%s", userId)

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdMonitoringList(c *cli.Context) {
	resp := utl.GetRequest("/monitoring/")
	fmt.Println(resp)
}

func cmdMonitoringShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetMonitoringByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/monitoring/%s", id)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
