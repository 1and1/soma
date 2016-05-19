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
						Name:   "create",
						Usage:  "Create a new bucket inside a repository",
						Action: runtime(cmdBucketCreate),
					},
					{
						Name:   "delete",
						Usage:  "Mark an existing bucket as deleted",
						Action: runtime(cmdBucketDelete),
					},
					{
						Name:   "restore",
						Usage:  "Restore a bucket marked as deleted",
						Action: runtime(cmdBucketRestore),
					},
					{
						Name:   "purge",
						Usage:  "Remove a deleted bucket",
						Action: runtime(cmdBucketPurge),
					},
					{
						Name:   "freeze",
						Usage:  "Freeze a bucket",
						Action: runtime(cmdBucketFreeze),
					},
					{
						Name:   "thaw",
						Usage:  "Thaw a frozen bucket",
						Action: runtime(cmdBucketThaw),
					},
					{
						Name:   "rename",
						Usage:  "Rename an existing bucket",
						Action: runtime(cmdBucketRename),
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
						Name:  "property",
						Usage: "SUBCOMMANDS for properties",
						Subcommands: []cli.Command{
							{
								Name:  "add",
								Usage: "SUBCOMMANDS for property add",
								Subcommands: []cli.Command{
									{
										Name:   "system",
										Usage:  "Add a system property to a bucket",
										Action: runtime(cmdBucketSystemPropertyAdd),
									},
									{
										Name:   "service",
										Usage:  "Add a service property to a bucket",
										Action: runtime(cmdBucketServicePropertyAdd),
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

	repoId := utl.TryGetRepositoryByUUIDOrName(opts["repository"][0])

	// fetch list of environments from SOMA to check if a valid
	// environment was requested
	utl.VerifyEnvironment(opts["environment"][0])

	req := proto.Request{
		Bucket: &proto.Bucket{
			Name:         c.Args().First(),
			RepositoryId: repoId,
			Environment:  opts["environment"][0],
		},
	}

	resp := utl.PostRequestWithBody(req, "/buckets/")
	fmt.Println(resp)
	return nil
}

func cmdBucketDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))
	buckId := utl.TryGetBucketByUUIDOrName(
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
	return nil
}

func cmdBucketRestore(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))
	buckId := utl.TryGetBucketByUUIDOrName(
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	req := proto.Request{
		Flags: &proto.Flags{
			Restore: true,
		},
	}

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
	return nil
}

func cmdBucketPurge(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))
	buckId := utl.TryGetBucketByUUIDOrName(
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	req := proto.Request{
		Flags: &proto.Flags{
			Purge: true,
		},
	}

	resp := utl.DeleteRequestWithBody(req, path)
	fmt.Println(resp)
	return nil
}

func cmdBucketFreeze(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))
	buckId := utl.TryGetBucketByUUIDOrName(
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	req := proto.Request{
		Flags: &proto.Flags{
			Freeze: true,
		},
	}

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
	return nil
}

func cmdBucketThaw(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))
	buckId := utl.TryGetBucketByUUIDOrName(
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	req := proto.Request{
		Flags: &proto.Flags{
			Thaw: true,
		},
	}

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
	return nil
}

func cmdBucketRename(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 5)
	utl.ValidateCliArgument(c, 2, "to")
	utl.ValidateCliArgument(c, 4, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(4))
	buckId := utl.TryGetBucketByUUIDOrName(
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	req := proto.Request{
		Bucket: &proto.Bucket{
			Name: c.Args().Get(2),
		},
	}

	resp := utl.PatchRequestWithBody(req, path)
	fmt.Println(resp)
	return nil
}

func cmdBucketList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest("/buckets/")
	fmt.Println(resp)
	return nil
}

func cmdBucketShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	bucketId := utl.BucketByUUIDOrName(c.Args().First())

	path := fmt.Sprintf("/buckets/%s", bucketId)
	resp := utl.GetRequest(path)
	fmt.Println(resp)
	return nil
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

	bucketId := utl.BucketByUUIDOrName(opts["to"][0])
	utl.CheckStringIsSystemProperty(c.Args().First())

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
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
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
	bucketId := utl.BucketByUUIDOrName(opts["to"][0])
	teamId := utl.TeamIdForBucket(bucketId)
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
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
