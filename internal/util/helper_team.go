package util

import (
	"log"

	"github.com/1and1/soma/lib/proto"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) tryGetTeamByUUIDOrName(c *resty.Client, s string) string {
	if u.isUUID(s) {
		return s
	}
	return u.getTeamIdByName(c, s)
}

func (u SomaUtil) getTeamIdByName(c *resty.Client, teamName string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Team: &proto.TeamFilter{
				Name: teamName,
			},
		},
	}

	resp := u.PostRequestWithBody(c, req, "/filter/teams/")
	teamResult := u.decodeProtoResultTeamFromResponse(resp)

	if teamName != (*teamResult.Teams)[0].Name {
		log.Fatal("Received result set for incorrect team")
	}
	return (*teamResult.Teams)[0].Id
}

func (u SomaUtil) decodeProtoResultTeamFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
