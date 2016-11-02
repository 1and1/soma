package util

import (
	"gopkg.in/resty.v0"
)

// GET
func (u SomaUtil) GetRequest(c *resty.Client, p string) *resty.Response {
	resp, err := c.R().Get(p)
	u.abortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

func (u SomaUtil) getRequestWithBody(c *resty.Client, body interface{}, p string) *resty.Response {
	resp, err := c.R().SetBody(body).SetContentLength(true).Get(p)
	u.abortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

// PUT
func (u SomaUtil) putRequest(c *resty.Client, p string) *resty.Response {
	resp, err := c.R().Put(p)
	u.abortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

func (u SomaUtil) PutRequestWithBody(c *resty.Client, body interface{}, p string) *resty.Response {
	resp, err := c.R().SetBody(body).SetContentLength(true).Put(p)
	u.abortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

// PATCH
func (u SomaUtil) patchRequestWithBody(c *resty.Client, body interface{}, p string) *resty.Response {
	resp, err := c.R().SetBody(body).SetContentLength(true).Patch(p)
	u.abortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

// POST
func (u SomaUtil) PostRequestWithBody(c *resty.Client, body interface{}, p string) *resty.Response {
	resp, err := c.R().SetBody(body).SetContentLength(true).Post(p)
	u.abortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

// DELETE
func (u SomaUtil) DeleteRequest(c *resty.Client, p string) *resty.Response {
	resp, err := c.R().Delete(p)
	u.abortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

func (u SomaUtil) DeleteRequestWithBody(c *resty.Client, body interface{}, p string) *resty.Response {
	resp, err := c.R().SetBody(body).SetContentLength(true).Delete(p)
	u.abortOnError(err)
	u.CheckRestyResponse(resp)
	return resp
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
