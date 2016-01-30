package main

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
)

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

	req := somaproto.ProtoRequestMetric{}
	req.Metric = &somaproto.ProtoMetric{}
	req.Metric.Metric = c.Args().First()
	req.Metric.Unit = opts["unit"][0]
	req.Metric.Description = opts["description"][0]

	pkgs := make([]somaproto.ProtoMetricProviderPackage, 0)
	if _, ok := opts["package"]; ok {
		for _, p := range opts["package"] {
			split := strings.SplitN(p, "::", 2)
			if len(split) != 2 {
				// coult not split
				utl.Abort(fmt.Sprintf("Syntax error, contains no :: %s",
					p))
			}
			pkgs = append(pkgs, somaproto.ProtoMetricProviderPackage{
				Provider: split[0],
				Package:  split[1],
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
