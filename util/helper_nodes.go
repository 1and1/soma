package util

import (
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetNodeByUUIDOrName(c *resty.Client, s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetNodeIdByName(c, s)
	}
	return id.String()
}

func (u SomaUtil) GetNodeIdByName(c *resty.Client, node string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Node: &proto.NodeFilter{
				Name: node,
			},
		},
	}

	resp := u.PostRequestWithBody(c, req, "/filter/nodes/")
	nodeResult := u.DecodeProtoResultNodeFromResponse(resp)

	if node != (*nodeResult.Nodes)[0].Name {
		u.Abort("Received result set for incorrect oncall duty")
	}
	return (*nodeResult.Nodes)[0].Id
}

func (u SomaUtil) GetNodeConfigById(c *resty.Client, node string) *proto.NodeConfig {
	if _, err := uuid.FromString(node); err != nil {
		node = u.GetNodeIdByName(c, node)
	}
	path := fmt.Sprintf("/nodes/%s/config", node)
	resp := u.GetRequest(c, path)
	nodeResult := u.DecodeResultFromResponse(resp)
	return (*nodeResult.Nodes)[0].Config
}

func (u SomaUtil) DecodeProtoResultNodeFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
