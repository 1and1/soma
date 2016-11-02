package util

import (
	"fmt"
	"strconv"

	"github.com/1and1/soma/lib/proto"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) tryGetOncallByUUIDOrName(c *resty.Client, s string) string {
	if u.isUUID(s) {
		return s
	}
	return u.getOncallIdByName(c, s)
}

func (u SomaUtil) getOncallIdByName(c *resty.Client, oncall string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Oncall: &proto.OncallFilter{
				Name: oncall,
			},
		},
	}

	resp := u.PostRequestWithBody(c, req, "/filter/oncall/")
	oncallResult := u.decodeProtoResultOncallFromResponse(resp)

	if oncall != (*oncallResult.Oncalls)[0].Name {
		u.abort("Received result set for incorrect oncall duty")
	}
	return (*oncallResult.Oncalls)[0].Id
}

func (u SomaUtil) GetOncallDetailsById(c *resty.Client, oncallid string) (string, string) {
	path := fmt.Sprintf("/oncall/%s", oncallid)
	resp := u.GetRequest(c, path)
	res := u.DecodeResultFromResponse(resp)

	if oncallid != (*res.Oncalls)[0].Id {
		u.abort(`Received result set for incorrect oncall duty`)
	}
	return (*res.Oncalls)[0].Name, (*res.Oncalls)[0].Number
}

func (u SomaUtil) decodeProtoResultOncallFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

func (u SomaUtil) validatePhoneNumber(n string) {
	num, err := strconv.Atoi(n)
	u.abortOnError(err, "Syntax error, argument is not a number")
	if num <= 0 || num > 9999 {
		u.abort("Phone number must be 4-digit extension")
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
