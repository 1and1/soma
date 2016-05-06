package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func registerBuckets(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// buckets
			{
				Name:   "buckets",
				Usage:  "SUBCOMMANDS for buckets",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new bucket inside a repository",
						Action: cmdBucketCreate,
					},
					{
						Name:   "delete",
						Usage:  "Mark an existing bucket as deleted",
						Action: cmdBucketDelete,
					},
					{
						Name:   "restore",
						Usage:  "Restore a bucket marked as deleted",
						Action: cmdBucketRestore,
					},
					{
						Name:   "purge",
						Usage:  "Remove a deleted bucket",
						Action: cmdBucketPurge,
					},
					{
						Name:   "freeze",
						Usage:  "Freeze a bucket",
						Action: cmdBucketFreeze,
					},
					{
						Name:   "thaw",
						Usage:  "Thaw a frozen bucket",
						Action: cmdBucketThaw,
					},
					{
						Name:   "rename",
						Usage:  "Rename an existing bucket",
						Action: cmdBucketRename,
					},
					{
						Name:   "list",
						Usage:  "List existing buckets",
						Action: cmdBucketList,
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific bucket",
						Action: cmdBucketShow,
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
										Action: cmdBucketSystemPropertyAdd,
									},
									{
										Name:   "service",
										Usage:  "Add a service property to a bucket",
										Action: cmdBucketServicePropertyAdd,
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

func cmdBucketCreate(c *cli.Context) {
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
	envResponse := utl.GetRequest("/environments/")
	envs := somaproto.ProtoResultEnvironmentList{}
	if err := json.Unmarshal(envResponse.Body(), &envs); err != nil {
		utl.Abort("Failed to unmarshal Environment data")
	}
	utl.ValidateStringInSlice(opts["environment"][0], envs.Environments)

	var req somaproto.ProtoRequestBucket
	req.Bucket = &somaproto.ProtoBucket{}
	req.Bucket.Name = c.Args().First()
	req.Bucket.Repository = repoId
	req.Bucket.Environment = opts["environment"][0]

	resp := utl.PostRequestWithBody(req, "/buckets/")
	fmt.Println(resp)
}

func cmdBucketDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))
	buckId := utl.TryGetBucketByUUIDOrName(
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	_ = utl.DeleteRequest(path)
}

func cmdBucketRestore(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))
	buckId := utl.TryGetBucketByUUIDOrName(
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	var req somaproto.ProtoRequestBucket
	req.Restore = true

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdBucketPurge(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))
	buckId := utl.TryGetBucketByUUIDOrName(
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	var req somaproto.ProtoRequestBucket
	req.Purge = true

	_ = utl.DeleteRequestWithBody(req, path)
}

func cmdBucketFreeze(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))
	buckId := utl.TryGetBucketByUUIDOrName(
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	var req somaproto.ProtoRequestBucket
	req.Freeze = true

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdBucketThaw(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))
	buckId := utl.TryGetBucketByUUIDOrName(
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	var req somaproto.ProtoRequestBucket
	req.Thaw = true

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdBucketRename(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 5)
	utl.ValidateCliArgument(c, 2, "to")
	utl.ValidateCliArgument(c, 4, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(4))
	buckId := utl.TryGetBucketByUUIDOrName(
		c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/buckets/%s", buckId)

	var req somaproto.ProtoRequestBucket
	req.Bucket = &somaproto.ProtoBucket{}
	req.Bucket.Name = c.Args().Get(2)

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdBucketList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest("/buckets/")
	fmt.Println(resp)
}

func cmdBucketShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	bucketId := utl.BucketByUUIDOrName(c.Args().First())

	path := fmt.Sprintf("/buckets/%s", bucketId)
	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdBucketSystemPropertyAdd(c *cli.Context) {
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

	sprop := somaproto.TreePropertySystem{
		Name:  c.Args().First(),
		Value: opts["value"][0],
	}

	tprop := somaproto.TreeProperty{
		PropertyType: "system",
		View:         opts["view"][0],
		System:       &sprop,
	}
	if _, ok := opts["inheritance"]; ok {
		tprop.Inheritance = utl.GetValidatedBool(opts["inheritance"][0])
	} else {
		tprop.Inheritance = true
	}
	if _, ok := opts["childrenonly"]; ok {
		tprop.ChildrenOnly = utl.GetValidatedBool(opts["childrenonly"][0])
	} else {
		tprop.ChildrenOnly = false
	}

	propList := []somaproto.TreeProperty{tprop}

	bucket := somaproto.ProtoBucket{
		Id:         bucketId,
		Properties: &propList,
	}

	req := somaproto.ProtoRequestBucket{
		Bucket: &bucket,
	}

	path := fmt.Sprintf("/buckets/%s/property/system/", bucketId)
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdBucketServicePropertyAdd(c *cli.Context) {
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
	tprop := somaproto.TreeProperty{
		PropertyType: "service",
		View:         opts["view"][0],
		Service: &somaproto.TreePropertyService{
			Name:       c.Args().First(),
			TeamId:     teamId,
			Attributes: []somaproto.TreeServiceAttribute{},
		},
	}
	if _, ok := opts["inheritance"]; ok {
		tprop.Inheritance = utl.GetValidatedBool(opts["inheritance"][0])
	} else {
		tprop.Inheritance = true
	}
	if _, ok := opts["childrenonly"]; ok {
		tprop.ChildrenOnly = utl.GetValidatedBool(opts["childrenonly"][0])
	} else {
		tprop.ChildrenOnly = false
	}

	req := somaproto.ProtoRequestBucket{
		Bucket: &somaproto.ProtoBucket{
			Id: bucketId,
			Properties: &[]somaproto.TreeProperty{
				tprop,
			},
		},
	}

	path := fmt.Sprintf("/buckets/%s/property/service/", bucketId)
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
