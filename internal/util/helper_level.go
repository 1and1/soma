package util

import (
	"fmt"

	"github.com/1and1/soma/lib/adm"
	"github.com/1and1/soma/lib/proto"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetLevelNameByNameOrShort(c *resty.Client, s string) string {
	req := proto.NewLevelFilter()
	req.Filter.Level.Name = s
	req.Filter.Level.ShortName = s

	resp, err := adm.PostReqBody(req, `/filter/levels/`)
	if err != nil {
		u.Abort(fmt.Sprintf("Level lookup request error: %s", err.Error()))
	}
	result, err := u.ResultFromResponse(resp)
	if se, ok := err.(SomaError); ok {
		if se.RequestError() {
			u.Abort(fmt.Sprintf("Level lookup request error: %s", se.Error()))
		}
		if se.Code() == 404 {
			u.Abort(fmt.Sprintf(
				"Could not find notification level with name %s",
				s,
			))
		}
		u.Abort(fmt.Sprintf("Level lookup application error: %s", err.Error()))
	}

	if s != (*result.Levels)[0].Name && s != (*result.Levels)[0].ShortName {
		u.Abort(fmt.Sprintf(
			"Notification level lookup failed. Wanted %s, received %s/%s",
			s,
			(*result.Levels)[0].Name,
			(*result.Levels)[0].ShortName,
		))
	}
	return (*result.Levels)[0].Name
}

func (u *SomaUtil) DecodeProtoResultLevelFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
