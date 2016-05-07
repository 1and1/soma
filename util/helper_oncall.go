package util

import (
	"strconv"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetOncallByUUIDOrName(s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetOncallIdByName(s)
	}
	return id.String()
}

func (u SomaUtil) GetOncallIdByName(oncall string) string {
	req := somaproto.ProtoRequestOncall{}
	req.Filter = &somaproto.ProtoOncallFilter{}
	req.Filter.Name = oncall

	resp := u.PostRequestWithBody(req, "/filter/oncall/")
	oncallResult := u.DecodeProtoResultOncallFromResponse(resp)

	if oncall != oncallResult.Oncalls[0].Name {
		u.Abort("Received result set for incorrect oncall duty")
	}
	return oncallResult.Oncalls[0].Id
}

func (u SomaUtil) DecodeProtoResultOncallFromResponse(resp *resty.Response) *somaproto.Result {
	return DecodeResultFromResponse(resp)
}

func (u SomaUtil) ValidatePhoneNumber(n string) {
	num, err := strconv.Atoi(n)
	u.AbortOnError(err, "Syntax error, argument is not a number")
	if num <= 0 || num > 9999 {
		u.Abort("Phone number must be 4-digit extension")
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
