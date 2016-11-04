package util

import (
	"fmt"

	"github.com/1and1/soma/lib/proto"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) tryGetNodeByUUIDOrName(c *resty.Client, s string) string {
	if u.isUUID(s) {
		return s
	}
	return u.getNodeIdByName(c, s)
}

func (u SomaUtil) getNodeIdByName(c *resty.Client, node string) string {
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
		u.abort("Received result set for incorrect oncall duty")
	}
	return (*nodeResult.Nodes)[0].Id
}

func (u SomaUtil) GetNodeConfigById(c *resty.Client, node string) *proto.NodeConfig {
	if !u.isUUID(node) {
		node = u.getNodeIdByName(c, node)
	}
	path := fmt.Sprintf("/nodes/%s/config", node)
	resp := u.GetRequest(c, path)
	nodeResult := u.unfilteredResultFromResponse(resp)
	if nodeResult.StatusCode == 404 {
		u.abort(`Node is not assigned to a configuration repository yet.`)
	} else if nodeResult.StatusCode > 299 {
		s := fmt.Sprintf("Request failed: %d - %s", nodeResult.StatusCode, nodeResult.StatusText)
		msgs := []string{s}
		if nodeResult.Errors != nil {
			msgs = append(msgs, *nodeResult.Errors...)
		}
		u.abort(msgs...)
	}
	return (*nodeResult.Nodes)[0].Config
}

func (u SomaUtil) TeamIdForNode(c *resty.Client, node string) string {
	nodeId := u.tryGetNodeByUUIDOrName(c, node)
	resp := u.GetRequest(c, fmt.Sprintf("/nodes/%s", nodeId))
	res := u.DecodeResultFromResponse(resp)
	if (*res.Nodes)[0].Id != nodeId {
		u.abort(`Received result for incorrect node`)
	}
	return (*res.Nodes)[0].TeamId
}

func (u SomaUtil) getNodeDetails(c *resty.Client, nodeId string) *proto.Node {
	resp := u.GetRequest(c, fmt.Sprintf("/nodes/%s", nodeId))
	res := u.DecodeResultFromResponse(resp)
	return &(*res.Nodes)[0]
}

func (u SomaUtil) findSourceForNodeProperty(c *resty.Client, pTyp, pName, view, nodeId string) string {
	node := u.getNodeDetails(c, nodeId)
	if node == nil {
		return ``
	}
	for _, prop := range *node.Properties {
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

func (u SomaUtil) DecodeProtoResultNodeFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
