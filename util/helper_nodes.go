package util

import (
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetNodeByUUIDOrName(s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetNodeIdByName(s)
	}
	return id.String()
}

func (u SomaUtil) GetNodeIdByName(node string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Node: &proto.NodeFilter{
				Name: node,
			},
		},
	}

	resp := u.PostRequestWithBody(req, "/filter/nodes/")
	nodeResult := u.DecodeProtoResultNodeFromResponse(resp)

	if node != (*nodeResult.Nodes)[0].Name {
		u.Abort("Received result set for incorrect oncall duty")
	}
	return (*nodeResult.Nodes)[0].Id
}

func (u SomaUtil) GetNodeConfigById(node string) *proto.NodeConfig {
	if _, err := uuid.FromString(node); err != nil {
		node = u.GetNodeIdByName(node)
	}
	path := fmt.Sprintf("/nodes/%s/config", node)
	resp := u.GetRequest(path)
	nodeResult := u.DecodeResultFromResponse(resp)
	return (*nodeResult.Nodes)[0].Config
}

func (u SomaUtil) TeamIdForNode(node string) string {
	nodeId := u.TryGetNodeByUUIDOrName(node)
	resp := u.GetRequest(fmt.Sprintf("/nodes/%s", nodeId))
	res := u.DecodeResultFromResponse(resp)
	if (*res.Nodes)[0].Id != nodeId {
		u.Abort(`Received result for incorrect node`)
	}
	return (*res.Nodes)[0].TeamId
}

func (u SomaUtil) DecodeProtoResultNodeFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
