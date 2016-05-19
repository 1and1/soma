package util

import (
	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetGroupByUUIDOrName(c *resty.Client, g string, b string) string {
	var id string
	bId := u.BucketByUUIDOrName(c, b)

	gId, err := uuid.FromString(g)
	if err != nil {
		id = u.GetGroupIdByName(c, g, bId)
	} else {
		id = gId.String()
	}
	return id
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
