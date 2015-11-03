package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetNodeByUUIDOrName(s string) uuid.UUID {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		id = u.GetNodeIdByName(s)
	}
	return id
}

func (u SomaUtil) GetNodeIdByName(node string) uuid.UUID {
	var req somaproto.ProtoRequestNode
	req.Filter.Name = node

	resp := u.GetRequestWithBody(req, "/nodes/")
	nodeResult := u.DecodeProtoResultNodeFromResponse(resp)

	if node != nodeResult.Nodes[0].Name {
		u.Abort("Received result set for incorrect oncall duty")
	}
	return nodeResult.Nodes[0].Id
}

func (u SomaUtil) DecodeProtoResultNodeFromResponse(resp *resty.Response) *somaproto.ProtoResultNode {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body))
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
