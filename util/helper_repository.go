package util

import (
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetRepositoryByUUIDOrName(s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetRepositoryIdByName(s)
	}
	return id.String()
}

func (u SomaUtil) GetRepositoryIdByName(repo string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Repository: &proto.RepositoryFilter{
				Name: repo,
			},
		},
	}

	resp := u.GetRequestWithBody(req, "/repository/")
	repoResult := u.DecodeProtoResultRepositoryFromResponse(resp)

	if repo != (*repoResult.Repositories)[0].Name {
		u.Abort("Received result set for incorrect repository")
	}
	return (*repoResult.Repositories)[0].Id
}

func (u SomaUtil) GetTeamIdByRepositoryId(repo string) string {
	repoId := u.TryGetRepositoryByUUIDOrName(repo)

	resp := u.GetRequest(fmt.Sprintf("/repository/%s", repoId))
	repoResult := u.DecodeResultFromResponse(resp)
	return (*repoResult.Repositories)[0].TeamId
}

func (u SomaUtil) DecodeProtoResultRepositoryFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
