package main

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
)

func registerMetrics(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// metrics
			{
				Name:   "metrics",
				Usage:  "SUBCOMMANDS for metrics",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new metric",
						Action: cmdMetricCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a metric",
						Action: cmdMetricDelete,
					},
					{
						Name:   "list",
						Usage:  "List metrics",
						Action: cmdMetricList,
					},
					{
						Name:   "show",
						Usage:  "Show details about a metric",
						Action: cmdMetricShow,
					},
				},
			}, // end metrics
		}...,
	)
	return &app
}

func cmdMetricCreate(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 5)
	multiple := []string{"package"}
	unique := []string{"unit", "description"}
	required := []string{"unit", "description"}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

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
				utl.Abort(fmt.Sprintf("Syntax error, contains no :: %s",
					p))
			}
			pkgs = append(pkgs, proto.MetricPackage{
				Provider: split[0],
				Name:     split[1],
			})
		}
		req.Metric.Packages = &pkgs
	}

	resp := utl.PostRequestWithBody(req, "/metrics/")
	fmt.Println(resp)
}

func cmdMetricDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/metrics/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdMetricList(c *cli.Context) {
	resp := utl.GetRequest("/metrics/")
	fmt.Println(resp)
}

func cmdMetricShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/metrics/%s", c.Args().First())

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
