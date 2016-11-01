package main

import (
	"fmt"
	"strings"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerMetrics(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// metrics
			{
				Name:  "metrics",
				Usage: "SUBCOMMANDS for metrics",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new metric",
						Action:       runtime(cmdMetricCreate),
						BashComplete: cmpl.MetricCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a metric",
						Action: runtime(cmdMetricDelete),
					},
					{
						Name:   "list",
						Usage:  "List metrics",
						Action: runtime(cmdMetricList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a metric",
						Action: runtime(cmdMetricShow),
					},
				},
			}, // end metrics
		}...,
	)
	return &app
}

func cmdMetricCreate(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 5)
	multiple := []string{"package"}
	unique := []string{"unit", "description"}
	required := []string{"unit", "description"}

	opts := adm.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

	utl.ValidateUnitExists(Client, opts["unit"][0])
	req := proto.Request{}
	req.Metric = &proto.Metric{}
	req.Metric.Path = c.Args().First()
	req.Metric.Unit = opts["unit"][0]
	req.Metric.Description = opts["description"][0]

	pkgs := make([]proto.MetricPackage, 0)
	if _, ok := opts["package"]; ok {
		for _, p := range opts["package"] {
			split := strings.SplitN(p, "::", 2)
			if len(split) != 2 {
				// coult not split
				adm.Abort(fmt.Sprintf("Syntax error, contains no :: %s",
					p))
			}
			utl.ValidateProviderExists(Client, split[0])
			pkgs = append(pkgs, proto.MetricPackage{
				Provider: split[0],
				Name:     split[1],
			})
		}
		req.Metric.Packages = &pkgs
	}

	resp := utl.PostRequestWithBody(Client, req, "/metrics/")
	fmt.Println(resp)
	return nil
}

func cmdMetricDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/metrics/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdMetricList(c *cli.Context) error {
	resp := utl.GetRequest(Client, "/metrics/")
	fmt.Println(resp)
	return nil
}

func cmdMetricShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/metrics/%s", c.Args().First())

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
