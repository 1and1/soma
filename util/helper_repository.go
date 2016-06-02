package util

import (
	"fmt"

	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetRepositoryByUUIDOrName(c *resty.Client, s string) string {
	if u.IsUUID(s) {
		return s
	}
	return u.GetRepositoryIdByName(c, s)
}

func (u SomaUtil) GetRepositoryIdByName(c *resty.Client, repo string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Repository: &proto.RepositoryFilter{
				Name: repo,
			},
		},
	}

	resp := u.PostRequestWithBody(c, req, "/filter/repository/")
	repoResult := u.DecodeProtoResultRepositoryFromResponse(resp)

	if repo != (*repoResult.Repositories)[0].Name {
		u.Abort("Received result set for incorrect repository")
	}
	return (*repoResult.Repositories)[0].Id
}

func (u SomaUtil) GetTeamIdByRepositoryId(c *resty.Client, repo string) string {
	repoId := u.TryGetRepositoryByUUIDOrName(c, repo)

	resp := u.GetRequest(c, fmt.Sprintf("/repository/%s", repoId))
	repoResult := u.DecodeResultFromResponse(resp)
	return (*repoResult.Repositories)[0].TeamId
}

func (u SomaUtil) DecodeProtoResultRepositoryFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
