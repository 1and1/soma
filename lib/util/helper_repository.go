package util

import (
	"fmt"

	"github.com/1and1/soma/lib/proto"
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

func (u SomaUtil) GetRepositoryDetails(c *resty.Client, repoId string) *proto.Repository {
	resp := u.GetRequest(c, fmt.Sprintf("/repository/%s", repoId))
	res := u.DecodeResultFromResponse(resp)
	return &(*res.Repositories)[0]
}

func (u SomaUtil) FindSourceForRepoProperty(c *resty.Client, pTyp, pName, view, repoId string) string {
	repo := u.GetRepositoryDetails(c, repoId)
	if repo == nil {
		return ``
	}
	for _, prop := range *repo.Properties {
		// wrong type
		if prop.Type != pTyp {
			continue
		}
		// wrong view
		if prop.View != view {
			continue
		}
		// inherited property
		if prop.InstanceId != prop.SourceInstanceId {
			continue
		}
		switch pTyp {
		case `system`:
			if prop.System.Name == pName {
				return prop.SourceInstanceId
			}
		case `oncall`:
			if prop.Oncall.Name == pName {
				return prop.SourceInstanceId
			}
		case `custom`:
			if prop.Custom.Name == pName {
				return prop.SourceInstanceId
			}
		case `service`:
			if prop.Service.Name == pName {
				return prop.SourceInstanceId
			}
		}
	}
	return ``
}

func (u SomaUtil) DecodeProtoResultRepositoryFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
