package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/resty.v0"
)

func (u SomaUtil) DecodeProtoResultAttributeFromResponse(resp *resty.Response) *somaproto.ProtoResultAttribute {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	var res somaproto.ProtoResultAttribute
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
