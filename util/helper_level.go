package util

import (
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

func (u *SomaUtil) DecodeProtoResultLevelFromResponse(resp *resty.Response) *somaproto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
