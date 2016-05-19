package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func registerChecks(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:  "checks",
				Usage: "SUBCOMMANDS for check configurations",
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new check configuration",
						Action: runtime(cmdCheckAdd),
					},
					{
						Name:   "list",
						Usage:  "List check configurations",
						Action: runtime(cmdCheckList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a check configuration",
						Action: runtime(cmdCheckShow),
					},
				},
			},
		}...,
	)
	return &app
}

func cmdCheckAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 8)

	opts, constraints, thresholds := utl.ParseVariadicCheckArguments(c.Args().Tail())

	req := proto.Request{}
	req.CheckConfig = &proto.CheckConfig{
		Name:         utl.ValidateRuneCount(c.Args().First(), 256),
		Interval:     utl.GetValidatedUint64(opts["interval"][0], 1),
		BucketId:     utl.BucketByUUIDOrName(opts["in"][0]),
		CapabilityId: utl.TryGetCapabilityByUUIDOrName(opts["with"][0]),
		ObjectType:   opts["on/type"][0],
	}
	req.CheckConfig.RepositoryId = utl.GetRepositoryIdForBucket(
		req.CheckConfig.BucketId)
	req.CheckConfig.ObjectId = utl.GetObjectIdForCheck(
		opts["on/type"][0],
		opts["on/object"][0],
		req.CheckConfig.BucketId)

	// clear bucketid if check is on a repository
	if req.CheckConfig.ObjectType == "repository" {
		req.CheckConfig.BucketId = ""
	}

	// optional argument: inheritance
	if iv, ok := opts["inheritance"]; ok {
		req.CheckConfig.Inheritance = utl.GetValidatedBool(iv[0])
	} else {
		// inheritance defaults to true
		req.CheckConfig.Inheritance = true
	}

	// optional argument: childrenonly
	if co, ok := opts["childrenonly"]; ok {
		req.CheckConfig.ChildrenOnly = utl.GetValidatedBool(co[0])
	} else {
		// childrenonly defaults to false
		req.CheckConfig.ChildrenOnly = false
	}

	// optional argument: extern
	if ex, ok := opts["extern"]; ok {
		req.CheckConfig.ExternalId = utl.ValidateRuneCount(ex[0], 64)
	}

	teamId := utl.GetTeamIdByRepositoryId(req.CheckConfig.RepositoryId)

	req.CheckConfig.Thresholds = utl.CleanThresholds(thresholds)
	req.CheckConfig.Constraints = utl.CleanConstraints(
		constraints,
		req.CheckConfig.RepositoryId,
		teamId)

	path := fmt.Sprintf("/checks/%s/", req.CheckConfig.RepositoryId)
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
	return nil
}

func cmdCheckList(c *cli.Context) error {
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
	return nil
}

func cmdCheckShow(c *cli.Context) error {
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
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
