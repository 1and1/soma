package util

import (
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetGroupByUUIDOrName(c *resty.Client, g string, b string) string {
	if u.IsUUID(g) {
		return g
	}
	bId := u.BucketByUUIDOrName(c, b)
	return u.GetGroupIdByName(c, g, bId)
}

func (u SomaUtil) GetGroupIdByName(c *resty.Client, g string, bId string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Group: &proto.GroupFilter{
				Name:     g,
				BucketId: bId,
			},
		},
	}

	resp := u.PostRequestWithBody(c, req, "/filter/groups/")
	groupResult := u.DecodeProtoResultGroupFromResponse(resp)

	return (*groupResult.Groups)[0].Id
}

func (u SomaUtil) DecodeProtoResultGroupFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
