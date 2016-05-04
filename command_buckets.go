package main

import (
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"
)

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
	utl.ValidateCliMinArgumentCount(c, 9)
	multiple := []string{}
	required := []string{"to", "in", "value", "view"}
	unique := []string{"to", "in", "value", "view", "inheritance", "childrenonly"}

	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	if opts["in"][0] != opts["to"][0] {
		utl.Abort(fmt.Sprintf("Bucket %s can not be in bucket %s", opts["to"][0], opts["in"][0]))
	}
	bucketId := utl.BucketByUUIDOrName(opts["in"][0])
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
