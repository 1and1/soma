package util

import (
	"bytes"
	"encoding/json"
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
	req := somaproto.ProtoRequestNode{}
	req.Filter = &somaproto.ProtoNodeFilter{}
	req.Filter.Name = node

	resp := u.PostRequestWithBody(req, "/filter/nodes/")
	nodeResult := u.DecodeProtoResultNodeFromResponse(resp)

	if node != nodeResult.Nodes[0].Name {
		u.Abort("Received result set for incorrect oncall duty")
	}
	return nodeResult.Nodes[0].Id
}

func (u SomaUtil) GetNodeConfigById(node string) *somaproto.ProtoNodeConfig {
	if _, err := uuid.FromString(node); err != nil {
		node = u.GetNodeIdByName(node)
	}
	path := fmt.Sprintf("/nodes/%s/config", node)
	resp := u.GetRequest(path)
	nodeResult := u.DecodeProtoResultNodeFromResponse(resp)
	return nodeResult.Nodes[0].Config
}

func (u SomaUtil) DecodeProtoResultNodeFromResponse(resp *resty.Response) *somaproto.ProtoResultNode {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	var res somaproto.ProtoResultNode
	err := decoder.Decode(&res)
	u.AbortOnError(err, "Error decoding server response body")
	if res.Code > 299 {
		s := fmt.Sprintf("Request failed: %d - %s", res.Code, res.Status)
		msgs := []string{s}
		msgs = append(msgs, res.Text...)
		u.Abort(msgs...)
	}
	return &res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
