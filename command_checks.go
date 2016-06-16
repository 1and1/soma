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
						Name:         "create",
						Usage:        "Create a new check configuration",
						Description:  help.CmdCheckAdd,
						Action:       runtime(cmdCheckAdd),
						BashComplete: cmdCheckAddBashComp,
					},
					{
						Name:        `delete`,
						Usage:       "Delete a check configuration",
						Description: help.CmdCheckDelete,
						Action:      runtime(cmdCheckDelete),
						BashComplete: func(c *cli.Context) {
							switch {
							case (c.NArg() % 2) == 1:
								for _, t := range []string{`in`} {
									fmt.Println(t)
								}
							}
						},
					},
					{
						Name:        "list",
						Usage:       "List check configurations",
						Description: help.CmdCheckList,
						Action:      runtime(cmdCheckList),
						BashComplete: func(c *cli.Context) {
							switch {
							case (c.NArg() % 2) == 0:
								for _, t := range []string{`in`} {
									fmt.Println(t)
								}
							}
						},
					},
					{
						Name:        "show",
						Usage:       "Show details about a check configuration",
						Description: help.CmdCheckShow,
						Action:      runtime(cmdCheckShow),
						BashComplete: func(c *cli.Context) {
							switch {
							case (c.NArg() % 2) == 1:
								for _, t := range []string{`in`} {
									fmt.Println(t)
								}
							}
						},
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
		BucketId:     utl.BucketByUUIDOrName(Client, opts["in"][0]),
		CapabilityId: utl.TryGetCapabilityByUUIDOrName(Client, opts["with"][0]),
		ObjectType:   opts["on/type"][0],
	}
	req.CheckConfig.RepositoryId = utl.GetRepositoryIdForBucket(
		Client, req.CheckConfig.BucketId)
	req.CheckConfig.ObjectId = utl.GetObjectIdForCheck(
		Client,
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

	teamId := utl.GetTeamIdByRepositoryId(Client, req.CheckConfig.RepositoryId)

	req.CheckConfig.Thresholds = utl.CleanThresholds(Client, thresholds)
	req.CheckConfig.Constraints = utl.CleanConstraints(
		Client,
		constraints,
		req.CheckConfig.RepositoryId,
		teamId)

	path := fmt.Sprintf("/checks/%s/", req.CheckConfig.RepositoryId)
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdCheckDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	multiple := []string{}
	unique := []string{"in"}
	required := []string{"in"}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	repoId := utl.GetRepositoryIdForBucket(Client, bucketId)
	checkId := utl.TryGetCheckByUUIDOrName(Client, c.Args().First(), repoId)

	path := fmt.Sprintf("/checks/%s/%s", repoId, checkId)
	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
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
	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	repoId := utl.GetRepositoryIdForBucket(Client, bucketId)

	path := fmt.Sprintf("/checks/%s/", repoId)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
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
	bucketId := utl.BucketByUUIDOrName(Client, opts["in"][0])
	repoId := utl.GetRepositoryIdForBucket(Client, bucketId)
	checkId := utl.TryGetCheckByUUIDOrName(Client, c.Args().First(), repoId)

	path := fmt.Sprintf("/checks/%s/%s", repoId, checkId)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

// I'm sorry as well.
func cmdCheckAddBashComp(c *cli.Context) {
	topArgs := []string{`in`, `on`, `with`, `interval`, `inheritance`, `childrenonly`, `extern`, `threshold`, `constraint`}
	thrArgs := []string{`predicate`, `level`, `value`}
	ctrArgs := []string{`service`, `oncall`, `attribute`, `system`, `native`, `custom`}
	onArgs := []string{`repository`, `bucket`, `group`, `cluster`, `node`}

	if c.NArg() == 0 {
		return
	}

	if c.NArg() == 1 {
		for _, t := range topArgs {
			fmt.Println(t)
		}
	}

	skipNext := 0
	subON := false
	subTHRESHOLD := false
	subCONSTRAINT := false

	hasIN := false
	hasON := false
	hasWITH := false
	hasINTERVAL := false
	hasINHERITANCE := false
	hasCHILDRENONLY := false
	hasEXTERN := false

	hasTHR_predicate := false
	hasTHR_level := false
	hasTHR_value := false

	hasCTR_service := false
	hasCTR_oncall := false
	hasCTR_attribute := false
	hasCTR_system := false
	hasCTR_native := false
	hasCTR_custom := false
	hasCTR_selected_service := false
	hasCTR_selected_oncall := false

	for _, t := range c.Args().Tail() {
		if skipNext > 0 {
			skipNext--
			continue
		}
		if subON {
			skipNext = 1
			subON = false
		}
		if subTHRESHOLD {
			if hasTHR_predicate && hasTHR_level && hasTHR_value {
				subTHRESHOLD = false
				hasTHR_predicate = false
				hasTHR_level = false
				hasTHR_value = false
			} else {
				switch t {
				case `predicate`:
					skipNext = 1
					hasTHR_predicate = true
					continue
				case `level`:
					skipNext = 1
					hasTHR_level = true
					continue
				case `value`:
					skipNext = 1
					hasTHR_value = true
					continue
				}
			}
		}
		if subCONSTRAINT {
			if hasCTR_selected_service {
				skipNext = 1
				hasCTR_selected_service = false
				continue
			}
			if hasCTR_selected_oncall {
				skipNext = 1
				hasCTR_selected_oncall = false
				continue
			}
			if hasCTR_service || hasCTR_oncall || hasCTR_attribute || hasCTR_system || hasCTR_native || hasCTR_custom {
				subCONSTRAINT = false
				hasCTR_service = false
				hasCTR_oncall = false
				hasCTR_attribute = false
				hasCTR_system = false
				hasCTR_native = false
				hasCTR_custom = false
				hasCTR_selected_service = false
				hasCTR_selected_oncall = false
			} else {
				switch t {
				case `service`:
					hasCTR_selected_service = true
					hasCTR_service = true
					continue
				case `oncall`:
					hasCTR_selected_oncall = true
					hasCTR_oncall = true
					continue
				case `attribute`:
					skipNext = 2
					hasCTR_attribute = true
					continue
				case `system`:
					skipNext = 2
					hasCTR_system = true
					continue
				case `native`:
					skipNext = 2
					hasCTR_native = true
					continue
				case `custom`:
					skipNext = 2
					hasCTR_custom = true
					continue
				}
			}

		}
		switch t {
		case `in`:
			skipNext = 1
			hasIN = true
			continue
		case `on`:
			hasON = true
			subON = true
			continue
		case `with`:
			skipNext = 1
			hasWITH = true
			continue
		case `interval`:
			skipNext = 1
			hasINTERVAL = true
			continue
		case `inheritance`:
			skipNext = 1
			hasINHERITANCE = true
			continue
		case `childrenonly`:
			skipNext = 1
			hasCHILDRENONLY = true
			continue
		case `extern`:
			skipNext = 1
			hasEXTERN = true
			continue
		case `threshold`:
			subTHRESHOLD = true
			continue
		case `constraint`:
			subCONSTRAINT = true
			continue
		}
	}
	// skipNext not yet consumed
	if skipNext > 0 {
		return
	}
	// in subchain: ON
	if subON {
		for _, t := range onArgs {
			fmt.Println(t)
		}
		return
	}
	// in subchain: CONSTRAINT
	if subCONSTRAINT {
		if hasCTR_selected_service || hasCTR_selected_oncall {
			fmt.Println(`name`)
			return
		}
		if !(hasCTR_service || hasCTR_oncall || hasCTR_attribute || hasCTR_system || hasCTR_native || hasCTR_custom) {
			for _, t := range ctrArgs {
				fmt.Println(t)
			}
			return
		}
	}
	// in subchain: THRESHOLD
	if subTHRESHOLD {
		if !(hasTHR_predicate && hasTHR_level && hasTHR_value) {
			for _, t := range thrArgs {
				switch t {
				case `predicate`:
					if !hasTHR_predicate {
						fmt.Println(t)
					}
				case `level`:
					if !hasTHR_level {
						fmt.Println(t)
					}
				case `value`:
					if !hasTHR_value {
						fmt.Println(t)
					}
				}
			}
			return
		}
	}
	// not in any subchain
	for _, t := range topArgs {
		switch t {
		case `in`:
			if !hasIN {
				fmt.Println(t)
			}
		case `on`:
			if !hasON {
				fmt.Println(t)
			}
		case `with`:
			if !hasWITH {
				fmt.Println(t)
			}
		case `interval`:
			if !hasINTERVAL {
				fmt.Println(t)
			}
		case `inheritance`:
			if !hasINHERITANCE {
				fmt.Println(t)
			}
		case `childrenonly`:
			if !hasCHILDRENONLY {
				fmt.Println(t)
			}
		case `extern`:
			if !hasEXTERN {
				fmt.Println(t)
			}
		default:
			fmt.Println(t)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
