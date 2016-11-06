package util

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/lib/proto"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) tryGetLevelNameByNameOrShort(c *resty.Client, s string) string {
	req := proto.NewLevelFilter()
	req.Filter.Level.Name = s
	req.Filter.Level.ShortName = s

	resp, err := adm.PostReqBody(req, `/filter/levels/`)
	if err != nil {
		u.abort(fmt.Sprintf("Level lookup request error: %s", err.Error()))
	}
	result, err := u.resultFromResponse(resp)
	if se, ok := err.(SomaError); ok {
		if se.RequestError() {
			u.abort(fmt.Sprintf("Level lookup request error: %s", se.Error()))
		}
		if se.Code() == 404 {
			u.abort(fmt.Sprintf(
				"Could not find notification level with name %s",
				s,
			))
		}
		u.abort(fmt.Sprintf("Level lookup application error: %s", err.Error()))
	}

	if s != (*result.Levels)[0].Name && s != (*result.Levels)[0].ShortName {
		u.abort(fmt.Sprintf(
			"Notification level lookup failed. Wanted %s, received %s/%s",
			s,
			(*result.Levels)[0].Name,
			(*result.Levels)[0].ShortName,
		))
	}
	return (*result.Levels)[0].Name
}

func (u *SomaUtil) decodeProtoResultLevelFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
