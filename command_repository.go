package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func cmdRepositoryCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "team")

	_ = utl.TryGetTeamByUUIDOrName(c.Args().Get(2))

	var req somaproto.ProtoRequestRepository
	req.Repository.Name = c.Args().Get(0)
	req.Repository.Team = c.Args().Get(2)

	_ = utl.PostRequestWithBody(req, "/repository/")
}

func cmdRepositoryDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id.String())

	_ = utl.DeleteRequest(path)
}

func cmdRepositoryRestore(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id.String())

	var req somaproto.ProtoRequestRepository
	req.Restore = true

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdRepositoryPurge(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id.String())

	var req somaproto.ProtoRequestRepository
	req.Purge = true

	_ = utl.DeleteRequestWithBody(req, path)
}

func cmdRepositoryClear(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id.String())

	var req somaproto.ProtoRequestRepository
	req.Clear = true

	_ = utl.PutRequestWithBody(req, path)
}

func cmdRepositoryRename(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id.String())

	var req somaproto.ProtoRequestRepository
	req.Repository.Name = c.Args().Get(2)

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdRepositoryRepossess(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to")
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	_ = utl.TryGetTeamByUUIDOrName(c.Args().Get(2))
	path := fmt.Sprintf("/repository/%s", id.String())

	var req somaproto.ProtoRequestRepository
	req.Repository.Team = c.Args().Get(2)

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdRepositoryClone(c *cli.Context) {
	utl.NotImplemented()
}

func cmdRepositoryActivate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id.String())

	var req somaproto.ProtoRequestRepository
	req.Activate = true

	_ = utl.PatchRequestWithBody(req, path)
}

func cmdRepositoryWipe(c *cli.Context) {
	utl.NotImplemented()
}

func cmdRepositoryList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)
	_ = utl.GetRequest("/repository/")
}

func cmdRepositoryShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetRepositoryByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/repository/%s", id.String())

	_ = utl.GetRequest(path)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
