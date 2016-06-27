package main

import (
	"fmt"
	"os"

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
						BashComplete: cmpl.Bucket,
					},
					{
						Name:         "restore",
						Usage:        "Restore a bucket marked as deleted",
						Action:       runtime(cmdBucketRestore),
						BashComplete: cmpl.Bucket,
					},
					{
						Name:         "purge",
						Usage:        "Remove a deleted bucket",
						Action:       runtime(cmdBucketPurge),
						BashComplete: cmpl.Bucket,
					},
					{
						Name:         "freeze",
						Usage:        "Freeze a bucket",
						Action:       runtime(cmdBucketFreeze),
						BashComplete: cmpl.Bucket,
					},
					{
						Name:         "thaw",
						Usage:        "Thaw a frozen bucket",
						Action:       runtime(cmdBucketThaw),
						BashComplete: cmpl.Bucket,
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
								Name:  "add",
								Usage: "SUBCOMMANDS for property add",
								Subcommands: []cli.Command{
									{
										Name:         "system",
										Usage:        "Add a system property to a bucket",
										Action:       runtime(cmdBucketSystemPropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         "service",
										Usage:        "Add a service property to a bucket",
										Action:       runtime(cmdBucketServicePropertyAdd),
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
	utl.ValidateCliArgumentCount(c, 5)
	multKeys := []string{"repository", "environment"}
	uniqKeys := []string{}

	opts := utl.ParseVariadicArguments(multKeys,
		uniqKeys,
		multKeys, // as reqKeys
		c.Args().Tail())

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
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().Get(2))
	buckId := utl.TryGetBucketByUUIDOrName(Client,
		c.Args().Get(0),
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
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().Get(2))
	buckId := utl.TryGetBucketByUUIDOrName(Client,
		c.Args().Get(0),
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
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().Get(2))
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
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().Get(2))
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
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().Get(2))
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
	utl.ValidateCliArgumentCount(c, 5)
	utl.ValidateCliArgument(c, 2, "to")
	utl.ValidateCliArgument(c, 4, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().Get(4))
	buckId := utl.TryGetBucketByUUIDOrName(Client,
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	req := proto.Request{
		Bucket: &proto.Bucket{
			Name: c.Args().Get(2),
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
	utl.ValidateCliArgumentCount(c, 0)

	if resp, err := adm.GetReq(`/buckets/`); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `list`)
	}
}

func cmdBucketShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	bucketId := utl.BucketByUUIDOrName(Client, c.Args().First())

	path := fmt.Sprintf("/buckets/%s", bucketId)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `show`)
	}
}

func cmdBucketTree(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	bucketId := utl.BucketByUUIDOrName(Client, c.Args().First())

	path := fmt.Sprintf("/buckets/%s/tree", bucketId)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, `tree`)
	}
}

func cmdBucketSystemPropertyAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	required := []string{"to", "value", "view"}
	unique := []string{"to", "in", "value", "view", "inheritance", "childrenonly"}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	if _, ok := opts["in"]; ok {
		fmt.Fprintln(os.Stderr, "Hint: Keyword `in` is DEPRECATED for buckets, since they are global objects. Ignoring.")
	}

	bucketId := utl.BucketByUUIDOrName(Client, opts["to"][0])
	utl.CheckStringIsSystemProperty(Client, c.Args().First())

	prop := proto.Property{
		Type: "system",
		View: opts["view"][0],
		System: &proto.PropertySystem{
			Name:  c.Args().First(),
			Value: opts["value"][0],
		},
	}
	if _, ok := opts["inheritance"]; ok {
		prop.Inheritance = utl.GetValidatedBool(opts["inheritance"][0])
	} else {
		prop.Inheritance = true
	}
	if _, ok := opts["childrenonly"]; ok {
		prop.ChildrenOnly = utl.GetValidatedBool(opts["childrenonly"][0])
	} else {
		prop.ChildrenOnly = false
	}

	req := proto.Request{
		Bucket: &proto.Bucket{
			Id: bucketId,
			Properties: &[]proto.Property{
				prop,
			},
		},
	}

	path := fmt.Sprintf("/buckets/%s/property/system/", bucketId)
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdBucketServicePropertyAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 5)
	multiple := []string{}
	required := []string{"to", "view"}
	unique := []string{"to", "in", "view", "inheritance", "childrenonly"}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	if _, ok := opts["in"]; ok {
		fmt.Fprintln(os.Stderr, "Hint: Keyword `in` is DEPRECATED for buckets, since they are global objects. Ignoring.")
	}
	bucketId := utl.BucketByUUIDOrName(Client, opts["to"][0])
	teamId := utl.TeamIdForBucket(Client, bucketId)
	// no reason to fill out the attributes, client-provided
	// attributes are discarded by the server
	prop := proto.Property{
		Type: "service",
		View: opts["view"][0],
		Service: &proto.PropertyService{
			Name:       c.Args().First(),
			TeamId:     teamId,
			Attributes: []proto.ServiceAttribute{},
		},
	}
	if _, ok := opts["inheritance"]; ok {
		prop.Inheritance = utl.GetValidatedBool(opts["inheritance"][0])
	} else {
		prop.Inheritance = true
	}
	if _, ok := opts["childrenonly"]; ok {
		prop.ChildrenOnly = utl.GetValidatedBool(opts["childrenonly"][0])
	} else {
		prop.ChildrenOnly = false
	}

	req := proto.Request{
		Bucket: &proto.Bucket{
			Id: bucketId,
			Properties: &[]proto.Property{
				prop,
			},
		},
	}

	path := fmt.Sprintf("/buckets/%s/property/service/", bucketId)
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
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
	utl.ValidateCliMinArgumentCount(c, 5)
	multiple := []string{}
	unique := []string{`from`, `view`}
	required := []string{`from`, `view`}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	bucketId := utl.BucketByUUIDOrName(Client, opts[`from`][0])

	if pType == `system` {
		utl.CheckStringIsSystemProperty(Client, c.Args().First())
	}
	sourceId := utl.FindSourceForBucketProperty(Client, pType, c.Args().First(),
		opts[`view`][0], bucketId)
	if sourceId == `` {
		utl.Abort(`Could not find locally set requested property.`)
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
