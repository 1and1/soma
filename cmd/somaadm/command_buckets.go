package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/internal/help"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerBuckets(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// buckets
			{
				Name:  "buckets",
				Usage: "SUBCOMMANDS for buckets",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new bucket inside a repository",
						Action:       runtime(cmdBucketCreate),
						BashComplete: cmpl.BucketCreate,
					},
					{
						Name:         "delete",
						Usage:        "Mark an existing bucket as deleted",
						Action:       runtime(cmdBucketDelete),
						BashComplete: cmpl.Repository,
					},
					{
						Name:         "restore",
						Usage:        "Restore a bucket marked as deleted",
						Action:       runtime(cmdBucketRestore),
						BashComplete: cmpl.Repository,
					},
					{
						Name:         "purge",
						Usage:        "Remove a deleted bucket",
						Action:       runtime(cmdBucketPurge),
						BashComplete: cmpl.Repository,
					},
					{
						Name:         "freeze",
						Usage:        "Freeze a bucket",
						Action:       runtime(cmdBucketFreeze),
						BashComplete: cmpl.Repository,
					},
					{
						Name:         "thaw",
						Usage:        "Thaw a frozen bucket",
						Action:       runtime(cmdBucketThaw),
						BashComplete: cmpl.Repository,
					},
					{
						Name:         "rename",
						Usage:        "Rename an existing bucket",
						Action:       runtime(cmdBucketRename),
						BashComplete: cmpl.BucketRename,
					},
					{
						Name:   "list",
						Usage:  "List existing buckets",
						Action: runtime(cmdBucketList),
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific bucket",
						Action: runtime(cmdBucketShow),
					},
					{
						Name:   `tree`,
						Usage:  `Display the bucket as tree`,
						Action: runtime(cmdBucketTree),
					},
					{
						Name:   `instances`,
						Usage:  `List check instances for a bucket`,
						Action: runtime(cmdBucketInstance),
					},
					{
						Name:  "property",
						Usage: "SUBCOMMANDS for properties",
						Subcommands: []cli.Command{
							{
								Name:        "add",
								Usage:       "SUBCOMMANDS for property add",
								Description: help.Text(`BucketsPropertyAdd`),
								Subcommands: []cli.Command{
									{
										Name:         "system",
										Usage:        "Add a system property to a bucket",
										Action:       runtime(cmdBucketSystemPropertyAdd),
										BashComplete: cmpl.PropertyAddValue,
									},
									{
										Name:         "service",
										Usage:        "Add a service property to a bucket",
										Action:       runtime(cmdBucketServicePropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         "oncall",
										Usage:        "Add an oncall property to a bucket",
										Action:       runtime(cmdBucketOncallPropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         "custom",
										Usage:        "Add a custom property to a bucket",
										Action:       runtime(cmdBucketCustomPropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
								},
							},
							{
								Name:  `delete`,
								Usage: `SUBCOMMANDS for property delete`,
								Subcommands: []cli.Command{
									{
										Name:         `system`,
										Usage:        `Delete a system property from a bucket`,
										Action:       runtime(cmdBucketSystemPropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `service`,
										Usage:        `Delete a service property from a bucket`,
										Action:       runtime(cmdBucketServicePropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `oncall`,
										Usage:        `Delete an oncall property from a bucket`,
										Action:       runtime(cmdBucketOncallPropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `custom`,
										Usage:        `Delete a custom property from a bucket`,
										Action:       runtime(cmdBucketCustomPropertyDelete),
										BashComplete: cmpl.FromView,
									},
								},
							},
						},
					},
				},
			}, // end buckets
		}...,
	)
	return &app
}

func cmdBucketCreate(c *cli.Context) error {
	uniqKeys := []string{`repository`, `environment`}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	repoId, err := adm.LookupRepoId(opts[`repository`][0])
	if err != nil {
		return err
	}

	// fetch list of environments from SOMA to check if a valid
	// environment was requested
	if err := adm.ValidateEnvironment(
		opts["environment"][0]); err != nil {
		return err
	}

	req := proto.Request{
		Bucket: &proto.Bucket{
			Name:         c.Args().First(),
			RepositoryId: repoId,
			Environment:  opts["environment"][0],
		},
	}

	if err := adm.ValidateRuneCountRange(req.Bucket.Name,
		4, 512); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/buckets/`, `command`, req, c)
}

func cmdBucketDelete(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`repository`},
		[]string{`repository`},
		c.Args().Tail()); err != nil {
		return err
	}
	buckId, err := adm.LookupBucketId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/buckets/%s", buckId)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdBucketRestore(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`repository`},
		[]string{`repository`},
		c.Args().Tail()); err != nil {
		return err
	}
	buckId, err := adm.LookupBucketId(c.Args().First())
	if err != nil {
		return err
	}

	req := proto.Request{
		Flags: &proto.Flags{
			Restore: true,
		},
	}

	path := fmt.Sprintf("/buckets/%s", buckId)
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdBucketPurge(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`repository`},
		[]string{`repository`},
		c.Args().Tail()); err != nil {
		return err
	}
	buckId, err := adm.LookupBucketId(c.Args().First())
	if err != nil {
		return err
	}
	req := proto.Request{
		Flags: &proto.Flags{
			Purge: true,
		},
	}

	path := fmt.Sprintf("/buckets/%s", buckId)
	return adm.Perform(`deletebody`, path, `command`, req, c)
}

func cmdBucketFreeze(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`repository`},
		[]string{`repository`},
		c.Args().Tail()); err != nil {
		return err
	}
	buckId, err := adm.LookupBucketId(c.Args().First())
	if err != nil {
		return err
	}

	req := proto.Request{
		Flags: &proto.Flags{
			Freeze: true,
		},
	}

	path := fmt.Sprintf("/buckets/%s", buckId)
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdBucketThaw(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`repository`},
		[]string{`repository`},
		c.Args().Tail()); err != nil {
		return err
	}
	buckId, err := adm.LookupBucketId(c.Args().First())
	if err != nil {
		return err
	}

	req := proto.Request{
		Flags: &proto.Flags{
			Thaw: true,
		},
	}

	path := fmt.Sprintf("/buckets/%s", buckId)
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdBucketRename(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`repository`, `to`},
		[]string{`repository`, `to`},
		c.Args().Tail(),
	); err != nil {
		return err
	}
	buckId, err := adm.LookupBucketId(c.Args().First())
	if err != nil {
		return err
	}

	req := proto.Request{
		Bucket: &proto.Bucket{
			Name: opts[`to`][0],
		},
	}

	path := fmt.Sprintf("/buckets/%s", buckId)
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdBucketList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/buckets/`, `list`, nil, c)
}

func cmdBucketShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	bucketId, err := adm.LookupBucketId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/buckets/%s", bucketId)
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdBucketInstance(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	bucketId, err := adm.LookupBucketId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/buckets/%s/instances/", bucketId)
	return adm.Perform(`get`, path, `list`, nil, c)
}

func cmdBucketTree(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	bucketId, err := adm.LookupBucketId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/buckets/%s/tree/tree", bucketId)
	return adm.Perform(`get`, path, `tree`, nil, c)
}

func cmdBucketSystemPropertyAdd(c *cli.Context) error {
	return cmdBucketPropertyAdd(c, `system`)
}

func cmdBucketServicePropertyAdd(c *cli.Context) error {
	return cmdBucketPropertyAdd(c, `service`)
}

func cmdBucketOncallPropertyAdd(c *cli.Context) error {
	return cmdBucketPropertyAdd(c, `oncall`)
}

func cmdBucketCustomPropertyAdd(c *cli.Context) error {
	return cmdBucketPropertyAdd(c, `custom`)
}

func cmdBucketPropertyAdd(c *cli.Context, pType string) error {
	return cmdPropertyAdd(c, pType, `bucket`)
}

func cmdBucketSystemPropertyDelete(c *cli.Context) error {
	return cmdBucketPropertyDelete(c, `system`)
}

func cmdBucketServicePropertyDelete(c *cli.Context) error {
	return cmdBucketPropertyDelete(c, `service`)
}

func cmdBucketOncallPropertyDelete(c *cli.Context) error {
	return cmdBucketPropertyDelete(c, `oncall`)
}

func cmdBucketCustomPropertyDelete(c *cli.Context) error {
	return cmdBucketPropertyDelete(c, `custom`)
}

func cmdBucketPropertyDelete(c *cli.Context, pType string) error {
	unique := []string{`from`, `view`}
	required := []string{`from`, `view`}
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
	bucketId, err := adm.LookupBucketId(opts[`from`][0])
	if err != nil {
		return err
	}

	if pType == `system` {
		if err := adm.ValidateSystemProperty(
			c.Args().First()); err != nil {
			return err
		}
	}
	var sourceId string
	if err := adm.FindBucketPropSrcId(pType, c.Args().First(),
		opts[`view`][0], bucketId, &sourceId); err != nil {
		return err
	}

	path := fmt.Sprintf("/buckets/%s/property/%s/%s",
		bucketId, pType, sourceId)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
