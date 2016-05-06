package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func registerChecks(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:   "checks",
				Usage:  "SUBCOMMANDS for check configurations",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new check configuration",
						Action: cmdCheckAdd,
					},
					{
						Name:   "list",
						Usage:  "List check configurations",
						Action: cmdCheckList,
					},
					{
						Name:   "show",
						Usage:  "Show details about a check configuration",
						Action: cmdCheckShow,
					},
				},
			},
		}...,
	)
	return &app
}

func cmdCheckAdd(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 8)

	opts, constraints, thresholds := utl.ParseVariadicCheckArguments(c.Args().Tail())

	req := somaproto.CheckConfigurationRequest{}
	req.CheckConfiguration = &somaproto.CheckConfiguration{
		Name:         utl.ValidateRuneCount(c.Args().First(), 256),
		Interval:     utl.GetValidatedUint64(opts["interval"][0], 1),
		BucketId:     utl.BucketByUUIDOrName(opts["in"][0]),
		CapabilityId: utl.TryGetCapabilityByUUIDOrName(opts["with"][0]),
		ObjectType:   opts["on/type"][0],
	}
	req.CheckConfiguration.RepositoryId = utl.GetRepositoryIdForBucket(
		req.CheckConfiguration.BucketId)
	req.CheckConfiguration.ObjectId = utl.GetObjectIdForCheck(
		opts["on/type"][0],
		opts["on/object"][0],
		req.CheckConfiguration.BucketId)

	// clear bucketid if check is on a repository
	if req.CheckConfiguration.ObjectType == "repository" {
		req.CheckConfiguration.BucketId = ""
	}

	// optional argument: inheritance
	if iv, ok := opts["inheritance"]; ok {
		req.CheckConfiguration.Inheritance = utl.GetValidatedBool(iv[0])
	} else {
		// inheritance defaults to true
		req.CheckConfiguration.Inheritance = true
	}

	// optional argument: childrenonly
	if co, ok := opts["childrenonly"]; ok {
		req.CheckConfiguration.ChildrenOnly = utl.GetValidatedBool(co[0])
	} else {
		// childrenonly defaults to false
		req.CheckConfiguration.ChildrenOnly = false
	}

	// optional argument: extern
	if ex, ok := opts["extern"]; ok {
		req.CheckConfiguration.ExternalId = utl.ValidateRuneCount(ex[0], 64)
	}

	teamId := utl.GetTeamIdByRepositoryId(req.CheckConfiguration.RepositoryId)

	req.CheckConfiguration.Thresholds = utl.CleanThresholds(thresholds)
	req.CheckConfiguration.Constraints = utl.CleanConstraints(
		constraints,
		req.CheckConfiguration.RepositoryId,
		teamId)

	path := fmt.Sprintf("/checks/%s/", req.CheckConfiguration.RepositoryId)
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdCheckList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 2)
	multiple := []string{}
	unique := []string{"in"}
	required := []string{"in"}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args())
	bucketId := utl.BucketByUUIDOrName(opts["in"][0])
	repoId := utl.GetRepositoryIdForBucket(bucketId)

	path := fmt.Sprintf("/checks/%s/", repoId)
	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdCheckShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multiple := []string{}
	unique := []string{"in"}
	required := []string{"in"}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(opts["in"][0])
	repoId := utl.GetRepositoryIdForBucket(bucketId)
	checkId := utl.TryGetCheckByUUIDOrName(c.Args().First(), repoId)

	path := fmt.Sprintf("/checks/%s/%s", repoId, checkId)
	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
