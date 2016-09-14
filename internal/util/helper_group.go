package util

import (
	"fmt"

	"github.com/1and1/soma/lib/proto"
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

func (u SomaUtil) GetGroupDetails(c *resty.Client, groupId string) *proto.Group {
	resp := u.GetRequest(c, fmt.Sprintf("/groups/%s", groupId))
	res := u.DecodeResultFromResponse(resp)
	return &(*res.Groups)[0]
}

func (u SomaUtil) FindSourceForGroupProperty(c *resty.Client, pTyp, pName, view, groupId string) string {
	group := u.GetGroupDetails(c, groupId)
	if group == nil {
		return ``
	}
	for _, prop := range *group.Properties {
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

func (u SomaUtil) DecodeProtoResultGroupFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
