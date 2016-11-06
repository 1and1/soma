package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/internal/help"
	"github.com/1and1/soma/lib/proto"
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
						Name:         "create",
						Usage:        "Create a new check configuration",
						Description:  help.Text(`ChecksCreate`),
						Action:       runtime(cmdCheckAdd),
						BashComplete: cmpl.CheckAdd,
					},
					{
						Name:         `delete`,
						Usage:        "Delete a check configuration",
						Description:  help.Text(`ChecksDelete`),
						Action:       runtime(cmdCheckDelete),
						BashComplete: cmpl.In,
					},
					{
						Name:         "list",
						Usage:        "List check configurations",
						Description:  help.Text(`ChecksList`),
						Action:       runtime(cmdCheckList),
						BashComplete: cmpl.In,
					},
					{
						Name:         "show",
						Usage:        "Show details about a check configuration",
						Description:  help.Text(`ChecksShow`),
						Action:       runtime(cmdCheckShow),
						BashComplete: cmpl.In,
					},
				},
			},
		}...,
	)
	return &app
}

func cmdCheckAdd(c *cli.Context) error {
	var (
		err    error
		teamId string
	)
	opts := make(map[string][]string)
	constraints := []proto.CheckConfigConstraint{}
	thresholds := []proto.CheckConfigThreshold{}

	if err = adm.ParseVariadicCheckArguments(
		opts,
		constraints,
		thresholds,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewCheckConfigRequest()
	if err = adm.ValidateLBoundUint64(opts["interval"][0],
		&req.CheckConfig.Interval, 1); err != nil {
		return err
	}
	if err = adm.ValidateRuneCount(c.Args().First(), 256); err != nil {
		return err
	}
	if req.CheckConfig.CapabilityId, err = adm.LookupCapabilityId(
		opts[`with`][0]); err != nil {
		return err
	}
	req.CheckConfig.ObjectType = opts[`on/type`][0]
	req.CheckConfig.Name = c.Args().First()
	req.CheckConfig.BucketId, err = adm.LookupBucketId(opts["in"][0])
	if err != nil {
		return err
	}
	if req.CheckConfig.RepositoryId, err = adm.LookupRepoByBucket(
		req.CheckConfig.BucketId); err != nil {
		return err
	}
	if req.CheckConfig.ObjectId, err = adm.LookupCheckObjectId(
		opts[`on/type`][0],
		opts[`on/object`][0],
		req.CheckConfig.BucketId,
	); err != nil {
		return err
	}

	// clear bucketid if check is on a repository
	if req.CheckConfig.ObjectType == "repository" {
		req.CheckConfig.BucketId = ""
	}

	// optional argument: inheritance
	if iv, ok := opts["inheritance"]; ok {
		if err := adm.ValidateBool(iv[0],
			&req.CheckConfig.Inheritance); err != nil {
			return err
		}
	} else {
		// inheritance defaults to true
		req.CheckConfig.Inheritance = true
	}

	// optional argument: childrenonly
	if co, ok := opts["childrenonly"]; ok {
		if err := adm.ValidateBool(co[0],
			&req.CheckConfig.ChildrenOnly); err != nil {
			return err
		}
	} else {
		// childrenonly defaults to false
		req.CheckConfig.ChildrenOnly = false
	}

	// optional argument: extern
	if ex, ok := opts["extern"]; ok {
		if err := adm.ValidateRuneCount(ex[0], 64); err != nil {
			return err
		}
		req.CheckConfig.ExternalId = ex[0]
	}

	if teamId, err = adm.LookupTeamByRepo(
		req.CheckConfig.RepositoryId); err != nil {
		return err
	}

	if req.CheckConfig.Thresholds, err = adm.ValidateThresholds(
		thresholds,
	); err != nil {
		return err
	}

	if req.CheckConfig.Constraints, err = adm.ValidateCheckConstraints(
		req.CheckConfig.RepositoryId,
		teamId,
		constraints,
	); err != nil {
		return err
	}

	path := fmt.Sprintf("/checks/%s/", req.CheckConfig.RepositoryId)
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdCheckDelete(c *cli.Context) error {
	multiple := []string{}
	unique := []string{"in"}
	required := []string{"in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multiple,
		unique,
		required,
		c.Args().Tail()); err != nil {
		return err
	}
	var (
		err                       error
		bucketId, repoId, checkId string
	)
	bucketId, err = adm.LookupBucketId(opts["in"][0])
	if err != nil {
		return err
	}
	if repoId, err = adm.LookupRepoByBucket(bucketId); err != nil {
		return err
	}
	if checkId, err = adm.LookupCheckConfigId(c.Args().First(),
		repoId); err != nil {
		return err
	}

	path := fmt.Sprintf("/checks/%s/%s", repoId, checkId)
	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdCheckList(c *cli.Context) error {
	multiple := []string{}
	unique := []string{"in"}
	required := []string{"in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multiple,
		unique,
		required,
		c.Args()); err != nil {
		return err
	}
	var (
		err              error
		bucketId, repoId string
	)
	bucketId, err = adm.LookupBucketId(opts["in"][0])
	if err != nil {
		return err
	}
	if repoId, err = adm.LookupRepoByBucket(bucketId); err != nil {
		return err
	}

	path := fmt.Sprintf("/checks/%s/", repoId)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdCheckShow(c *cli.Context) error {
	multiple := []string{}
	unique := []string{"in"}
	required := []string{"in"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multiple,
		unique,
		required,
		c.Args().Tail()); err != nil {
		return err
	}
	var (
		err                       error
		bucketId, repoId, checkId string
	)
	if bucketId, err = adm.LookupBucketId(opts["in"][0]); err != nil {
		return err
	}
	if repoId, err = adm.LookupRepoByBucket(bucketId); err != nil {
		return err
	}
	if checkId, err = adm.LookupCheckConfigId(c.Args().First(),
		repoId); err != nil {
		return err
	}

	path := fmt.Sprintf("/checks/%s/%s", repoId, checkId)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
