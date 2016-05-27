package util

import (
	"gopkg.in/resty.v0"
)

// GET
func (u SomaUtil) GetRequest(p string) *resty.Response {
	u.ApiUrl.Path = p
	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(u.ApiUrl.String())
	u.AbortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

func (u SomaUtil) GetRequestWithBody(body interface{}, p string) *resty.Response {
	u.ApiUrl.Path = p
	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(body).
		Get(u.ApiUrl.String())
	u.AbortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

// PUT
func (u SomaUtil) PutRequest(p string) *resty.Response {
	u.ApiUrl.Path = p
	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Put(u.ApiUrl.String())
	u.AbortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

func (u SomaUtil) PutRequestWithBody(body interface{}, p string) *resty.Response {
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

// PATCH
func (u SomaUtil) PatchRequestWithBody(body interface{}, p string) *resty.Response {
	u.ApiUrl.Path = p
	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(body).
		Patch(u.ApiUrl.String())
	u.AbortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

// POST
func (u SomaUtil) PostRequestWithBody(body interface{}, p string) *resty.Response {
	u.ApiUrl.Path = p
	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(body).
		Post(u.ApiUrl.String())
	u.AbortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

// DELETE
func (u SomaUtil) DeleteRequest(p string) *resty.Response {
	u.ApiUrl.Path = p
	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Delete(u.ApiUrl.String())
	u.AbortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

func (u SomaUtil) DeleteRequestWithBody(body interface{}, p string) *resty.Response {
	u.ApiUrl.Path = p
	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(body).
		Delete(u.ApiUrl.String())
	u.AbortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
