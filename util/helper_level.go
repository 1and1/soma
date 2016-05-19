package util

import (
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetLevelNameByNameOrShort(c *resty.Client, s string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Level: &proto.LevelFilter{
				Name:      s,
				ShortName: s,
			},
		},
	}
	resp := u.PostRequestWithBody(c, req, "/filter/levels/")
	levelResult := u.DecodeProtoResultLevelFromResponse(resp)

	if s != (*levelResult.Levels)[0].Name && s != (*levelResult.Levels)[0].ShortName {
		u.Abort("Received result set for incorrect level")
	}
	return (*levelResult.Levels)[0].Name
}

func (u *SomaUtil) DecodeProtoResultLevelFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
