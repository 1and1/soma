package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetGroupByUUIDOrName(g string, b string) string {
	var id string
	bId := u.BucketByUUIDOrName(b)

	gId, err := uuid.FromString(g)
	if err != nil {
		id = u.GetGroupIdByName(g, bId)
	} else {
		id = gId.String()
	}
	return id
}

func (u SomaUtil) GetGroupIdByName(g string, bId string) string {
	var req somaproto.ProtoRequestGroup
	req.Filter.Name = g
	req.Filter.BucketId = bId

	resp := u.GetRequestWithBody(req, "/groups/")
	groupResult := u.DecodeProtoResultGroupFromResponse(resp)

	return groupResult.Groups[0].Id
}

func (u SomaUtil) DecodeProtoResultGroupFromResponse(resp *resty.Response) *somaproto.ProtoResultGroup {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	var res somaproto.ProtoResultGroup
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
