package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetLevelNameByNameOrShort(s string) string {
	req := somaproto.ProtoRequestLevel{}
	req.Filter = &somaproto.ProtoLevelFilter{
		Name:      s,
		ShortName: s,
	}
	resp := u.PostRequestWithBody(req, "/filter/levels/")
	levelResult := u.DecodeProtoResultLevelFromResponse(resp)

	if s != levelResult.Levels[0].Name && s != levelResult.Levels[0].ShortName {
		u.Abort("Received result set for incorrect level")
	}
	return levelResult.Levels[0].Name
}

func (u *SomaUtil) DecodeProtoResultLevelFromResponse(resp *resty.Response) *somaproto.ProtoResultLevel {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	var res somaproto.ProtoResultLevel
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
