package main

import (
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"
)

func registerJobs(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  `jobs`,
				Usage: `SUBCOMMANDS for job information`,
				Subcommands: []cli.Command{
					{
						Name:   `list`,
						Usage:  `List outstanding jobs`,
						Action: runtime(cmdJobList),
					},
					{
						Name:   `show`,
						Usage:  `Show details about a job`,
						Action: runtime(cmdJobShow),
					},
					{
						Name:  `local`,
						Usage: `SUBCOMMANDS for locally saved jobs`,
						Subcommands: []cli.Command{
							{
								Name:   `outstanding`,
								Usage:  `List outstanding locally saved Jobs`,
								Action: runtime(cmdJobLocalOutstanding),
							},
							{
								Name:   `update`,
								Usage:  `Check and update status of oustanding locally cached jobs`,
								Action: runtime(cmdJobLocalUpdate),
							},
							{
								Name:   `list`,
								Usage:  `List all locally cached jobs`,
								Action: runtime(cmdJobLocalList),
							},
						},
					},
				},
			},
		}...,
	)
	return &app
}

func cmdJobList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)
	if resp, err := adm.GetReq(`/jobs/`); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdJobShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	if !utl.IsUUID(c.Args().First()) {
		return fmt.Errorf("Argument is not a UUID: %s", c.Args().First())
	}
	path := fmt.Sprintf("/jobs/%s", c.Args().First())
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdJobLocalOutstanding(c *cli.Context) error {
	jobs, err := store.GetActiveJobs()
	if err != nil {
		return err
	}

	if enc, err := json.Marshal(&jobs); err != nil {
		return err
	} else {
		fmt.Println(string(enc))
	}
	return nil
}

func cmdJobLocalUpdate(c *cli.Context) error {
	return nil
}

func cmdJobLocalList(c *cli.Context) error {
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
