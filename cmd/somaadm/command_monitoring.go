package main

import (
	"fmt"
	"strings"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerMonitoring(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// monitoring
			{
				Name:  "monitoring",
				Usage: "SUBCOMMANDS for monitoring systems",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new monitoring system",
						Action:       runtime(cmdMonitoringCreate),
						BashComplete: cmpl.MonitoringCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a monitoring system",
						Action: runtime(cmdMonitoringDelete),
					},
					{
						Name:   "list",
						Usage:  "List monitoring systems",
						Action: runtime(cmdMonitoringList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a monitoring system",
						Action: runtime(cmdMonitoringShow),
					},
				},
			}, // end monitoring
		}...,
	)
	return &app
}

func cmdMonitoringCreate(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	unique := []string{"mode", "contact", "team", "callback"}
	required := []string{"mode", "contact", "team"}

	opts := adm.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

	req := proto.Request{}
	req.Monitoring = &proto.Monitoring{}
	req.Monitoring.Name = c.Args().First()
	req.Monitoring.Mode = opts["mode"][0]
	req.Monitoring.Contact = utl.TryGetUserByUUIDOrName(Client, opts["contact"][0])
	req.Monitoring.TeamId = utl.TryGetTeamByUUIDOrName(Client, opts["team"][0])
	if strings.Contains(req.Monitoring.Name, `.`) {
		adm.Abort(`Monitoring system names must not contain the character '.'`)
	}

	// optional arguments
	if _, ok := opts["callback"]; ok {
		req.Monitoring.Callback = opts["callback"][0]
	}

	resp := utl.PostRequestWithBody(Client, req, "/monitoring/")
	fmt.Println(resp)
	return nil
}

func cmdMonitoringDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	userId := utl.TryGetMonitoringByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/monitoring/%s", userId)

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdMonitoringList(c *cli.Context) error {
	resp := utl.GetRequest(Client, "/monitoring/")
	fmt.Println(resp)
	return nil
}

func cmdMonitoringShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetMonitoringByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/monitoring/%s", id)

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
