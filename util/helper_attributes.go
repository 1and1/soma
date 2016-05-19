package util

import (
	"fmt"

	"gopkg.in/resty.v0"
)

func (u SomaUtil) CheckStringIsServiceAttribute(c *resty.Client, s string) {
	resp := u.GetRequest(c, "/attributes/")
	res := u.DecodeResultFromResponse(resp)

	for _, attr := range *res.Attributes {
		if attr.Name == s {
			return
		}
	}
	u.Abort(fmt.Sprintf("Invalid service attribute requested: %s", s))
}

func (u SomaUtil) DecodeProtoResultAttributeFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
