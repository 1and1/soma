package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func cmdRepositoryCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "team")

	teamId := utl.TryGetTeamByUUIDOrName(c.Args().Get(2))

	var req somaproto.ProtoRequestRepository
	req.Repository = &somaproto.ProtoRepository{}
	req.Repository.Name = c.Args().Get(0)
	req.Repository.Team = teamId

	resp := utl.PostRequestWithBody(req, "/repository/")
	fmt.Println(resp)
}

func cmdRepositoryDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	_ = utl.DeleteRequest(path)
}

func cmdRepositoryRestore(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	var req somaproto.ProtoRequestRepository
	req.Restore = true

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdRepositoryPurge(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	var req somaproto.ProtoRequestRepository
	req.Purge = true

	_ = utl.DeleteRequestWithBody(req, path)
}

func cmdRepositoryClear(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	var req somaproto.ProtoRequestRepository
	req.Clear = true

	_ = utl.PutRequestWithBody(req, path)
}

func cmdRepositoryRename(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	var req somaproto.ProtoRequestRepository
	req.Repository = &somaproto.ProtoRepository{}
	req.Repository.Name = c.Args().Get(2)

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdRepositoryRepossess(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	_ = utl.TryGetTeamByUUIDOrName(c.Args().Get(2))
	path := fmt.Sprintf("/repository/%s", id)

	var req somaproto.ProtoRequestRepository
	req.Repository = &somaproto.ProtoRepository{}
	req.Repository.Team = c.Args().Get(2)

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdRepositoryClone(c *cli.Context) {
	utl.NotImplemented()
}

func cmdRepositoryActivate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	var req somaproto.ProtoRequestRepository
	req.Activate = true

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdRepositoryWipe(c *cli.Context) {
	utl.NotImplemented()
}

func cmdRepositoryList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)
	resp := utl.GetRequest("/repository/")
	fmt.Println(resp)
}

func cmdRepositoryShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdRepositorySystemPropertyAdd(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 7)
	multiple := []string{}
	required := []string{"to", "value", "view"}
	unique := []string{"to", "in", "value", "view", "inheritance", "childrenonly"}
	opts := utl.ParseVariadicArguments(multiple, unique, required, c.Args().Tail())
	if _, ok := opts["in"]; ok {
		fmt.Fprintln(os.Stderr, "Hint: Keyword `in` is DEPRECATED for repositories, since they are global objects. Ignoring.")
	}

	repositoryId := utl.TryGetRepositoryByUUIDOrName(opts["to"][0])
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

	repository := somaproto.ProtoRepository{
		Id:         repositoryId,
		Properties: &propList,
	}

	req := somaproto.ProtoRequestRepository{
		Repository: &repository,
	}

	path := fmt.Sprintf("/repository/%s/property/system/", repositoryId)
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
