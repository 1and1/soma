package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetOncallByUUIDOrName(s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetOncallIdByName(s)
	}
	return id.String()
}

func (u SomaUtil) GetOncallIdByName(oncall string) string {
	url := u.ApiUrl
	url.Path = "/oncall/"

	var req somaproto.ProtoRequestOncall
	req.Filter.Name = oncall

	resp := u.GetRequestWithBody(req, url.String())
	oncallResult := u.DecodeProtoResultOncallFromResponse(resp)

	if oncall != oncallResult.Oncalls[0].Name {
		u.Abort("Received result set for incorrect oncall duty")
	}
	return oncallResult.Oncalls[0].Id
}

func (u SomaUtil) DecodeProtoResultOncallFromResponse(resp *resty.Response) *somaproto.ProtoResultOncall {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body))
	var res somaproto.ProtoResultOncall
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
