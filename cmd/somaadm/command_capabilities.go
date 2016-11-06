package main

import (
	"fmt"

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
	multiple := []string{}
	unique := []string{"metric", "view", "thresholds"}
	required := []string{"metric", "view", "thresholds"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multiple,
		unique,
		required,
		c.Args().Tail()); err != nil {
		return err
	}

	var thresholds uint64
	if err := adm.ValidateLBoundUint64(opts["thresholds"][0],
		&thresholds, 1); err != nil {
		return err
	}

	req := proto.Request{
		Capability: &proto.Capability{
			Metric:     opts["metric"][0],
			View:       opts["view"][0],
			Thresholds: thresholds,
		},
	}
	var err error
	if req.Capability.MonitoringId, err = adm.LookupMonitoringId(
		c.Args().First()); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/capability/`, `command`, req, c)
}

func cmdCapabilityRevoke(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	var path string
	if id, err := adm.LookupCapabilityId(
		c.Args().First()); err != nil {
		return err
	} else {
		path = fmt.Sprintf("/capability/%s", id)
	}

	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdCapabilityList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/capability/`, `list`, nil, c)
}

func cmdCapabilityShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	var path string
	if id, err := adm.LookupCapabilityId(
		c.Args().First()); err != nil {
		return err
	} else {
		path = fmt.Sprintf("/capability/%s", id)
	}

	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
