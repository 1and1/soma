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
	unique := []string{"mode", "contact", "team", "callback"}
	required := []string{"mode", "contact", "team"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewMonitoringRequest()
	req.Monitoring.Name = c.Args().First()
	req.Monitoring.Mode = opts["mode"][0]
	var err error
	if req.Monitoring.Contact, err = adm.LookupUserId(
		opts[`contact`][0]); err != nil {
		return err
	}
	req.Monitoring.TeamId, err = adm.LookupTeamId(opts[`team`][0])
	if err != nil {
		return err
	}
	if strings.Contains(req.Monitoring.Name, `.`) {
		return fmt.Errorf(
			`Monitoring system names must not contain` +
				` the character '.'`)
	}

	// optional arguments
	if _, ok := opts["callback"]; ok {
		req.Monitoring.Callback = opts["callback"][0]
	}

	return adm.Perform(`postbody`, `/monitoring/`, `command`, req, c)
}

func cmdMonitoringDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	monitoringId, err := adm.LookupMonitoringId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/monitoring/%s", monitoringId)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdMonitoringList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/monitoring/`, `list`, nil, c)
}

func cmdMonitoringShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	monitoringId, err := adm.LookupMonitoringId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/monitoring/%s", monitoringId)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
