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

	repoId := utl.TryGetRepositoryByUUIDOrName(Client, opts["repository"][0])

	// fetch list of environments from SOMA to check if a valid
	// environment was requested
	utl.VerifyEnvironment(Client, opts["environment"][0])

	req := proto.Request{
		Bucket: &proto.Bucket{
			Name:         c.Args().First(),
			RepositoryId: repoId,
			Environment:  opts["environment"][0],
		},
	}

	utl.ValidateRuneCountRange(req.Bucket.Name, 4, 512)

	if resp, err := adm.PostReqBody(req, "/buckets/"); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
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
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, opts[`repository`][0])
	buckId := utl.TryGetBucketByUUIDOrName(Client,
		c.Args().First(),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
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
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, opts[`repository`][0])
	buckId := utl.TryGetBucketByUUIDOrName(Client,
		c.Args().First(),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	req := proto.Request{
		Flags: &proto.Flags{
			Restore: true,
		},
	}

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
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
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, opts[`repository`][0])
	buckId := utl.TryGetBucketByUUIDOrName(Client,
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	req := proto.Request{
		Flags: &proto.Flags{
			Purge: true,
		},
	}

	if resp, err := adm.DeleteReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
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
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, opts[`repository`][0])
	buckId := utl.TryGetBucketByUUIDOrName(Client,
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	req := proto.Request{
		Flags: &proto.Flags{
			Freeze: true,
		},
	}

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
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
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, opts[`repository`][0])
	buckId := utl.TryGetBucketByUUIDOrName(Client,
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	req := proto.Request{
		Flags: &proto.Flags{
			Thaw: true,
		},
	}

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdBucketRename(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`repository`, `to`},
		[]string{`repository`, `to`},
		c.Args().Tail()); err != nil {
		return err
	}
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, opts[`repository`][0])
	buckId := utl.TryGetBucketByUUIDOrName(Client,
		c.Args().First(),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	req := proto.Request{
		Bucket: &proto.Bucket{
			Name: opts[`to`][0],
		},
	}

	if resp, err := adm.PatchReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdBucketList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	if resp, err := adm.GetReq(`/buckets/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdBucketShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	bucketId := utl.BucketByUUIDOrName(Client, c.Args().First())

	path := fmt.Sprintf("/buckets/%s", bucketId)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `show`)
	}
}

func cmdBucketTree(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	bucketId := utl.BucketByUUIDOrName(Client, c.Args().First())

	path := fmt.Sprintf("/buckets/%s/tree/tree", bucketId)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `tree`)
	}
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
	multiple := []string{}
	unique := []string{`from`, `view`}
	required := []string{`from`, `view`}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(opts, multiple, unique, required,
		c.Args().Tail()); err != nil {
		return err
	}
	bucketId := utl.BucketByUUIDOrName(Client, opts[`from`][0])

	if pType == `system` {
		utl.CheckStringIsSystemProperty(Client, c.Args().First())
	}
	sourceId := utl.FindSourceForBucketProperty(Client, pType, c.Args().First(),
		opts[`view`][0], bucketId)
	if sourceId == `` {
		adm.Abort(`Could not find locally set requested property.`)
	}

	path := fmt.Sprintf("/buckets/%s/property/%s/%s",
		bucketId, pType, sourceId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `delete`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
