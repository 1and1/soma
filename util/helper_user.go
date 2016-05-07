package util

import (
	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetUserByUUIDOrName(s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetUserIdByName(s)
	}
	return id.String()
}

func (u *SomaUtil) GetUserIdByName(user string) string {
	req := somaproto.ProtoRequestUser{}
	req.Filter = &somaproto.ProtoUserFilter{}
	req.Filter.UserName = user

	resp := u.PostRequestWithBody(req, "/filter/users/")
	userResult := u.DecodeProtoResultUserFromResponse(resp)

	if user != userResult.Users[0].UserName {
		u.Abort("Received result set for incorrect user")
	}
	return userResult.Users[0].Id
}

func (u *SomaUtil) DecodeProtoResultUserFromResponse(resp *resty.Response) *somaproto.Result {
	return DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
