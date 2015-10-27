package util

import (
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) PutRequestWithBody(body interface{}, p string) *resty.Response {
	u.ApiUrl.Path = p
	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(body).
		Put(u.ApiUrl.String())
	u.AbortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

func (u *SomaUtil) GetRequest(p string) *resty.Response {
	u.ApiUrl.Path = p
	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(u.ApiUrl.String())
	u.AbortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
