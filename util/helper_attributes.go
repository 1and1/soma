package util

import (
	"gopkg.in/resty.v0"
)

func (u SomaUtil) DecodeProtoResultAttributeFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
