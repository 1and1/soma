/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerWorkflow(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  `workflow`,
				Usage: `SUBCOMMANDS for workflow inquiry`,
				Subcommands: []cli.Command{
					{
						Name:   `summary`,
						Usage:  `Show summary of workflow status`,
						Action: runtime(cmdWorkflowSummary),
					},
					{
						Name:   `list`,
						Usage:  `List instances in a specific workflow state`,
						Action: runtime(cmdWorkflowList),
					},
					{
						Name:   `retry`,
						Usage:  `Reschedule an instance in a failed state`,
						Action: runtime(cmdWorkflowRetry),
					},
					{
						Name:   `set`,
						Usage:  `Hard-set an instance's worflow status`,
						Action: runtime(cmdWorkflowSet),
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  `force, f`,
								Usage: `Force is required to break the workflow`,
							},
						},
					},
				},
			},
		}...,
	)
	return &app
}

func cmdWorkflowSummary(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/workflow/summary`, `list`, nil, c)
}

func cmdWorkflowList(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if err := adm.ValidateStatus(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewWorkflowFilter()
	req.Filter.Workflow.Status = c.Args().First()

	return adm.Perform(`postbody`, `/filter/workflow/`, `list`, req, c)
}

func cmdWorkflowRetry(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if err := adm.ValidateInstance(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewWorkflowRequest()
	req.Workflow.InstanceId = c.Args().First()

	path := `/workflow/retry`
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdWorkflowSet(c *cli.Context) error {
	unique := []string{`status`, `next`}
	required := []string{`status`, `next`}
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
	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf("Argument is not a UUID: %s",
			c.Args().First())
	}
	if err := adm.ValidateStatus(opts[`status`][0]); err != nil {
		return err
	}
	if err := adm.ValidateStatus(opts[`next`][0]); err != nil {
		return err
	}
	req := proto.NewWorkflowRequest()
	req.Flags.Forced = c.Bool(`force`)
	req.Workflow.InstanceConfigId = c.Args().First()
	req.Workflow.Status = opts[`status`][0]
	req.Workflow.NextStatus = opts[`next`][0]

	path := fmt.Sprintf("/workflow/instanceconfig/%s",
		c.Args().First())
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
