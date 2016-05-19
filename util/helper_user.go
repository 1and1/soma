package util

import (
	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetUserByUUIDOrName(c *resty.Client, s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetUserIdByName(c, s)
	}
	return id.String()
}

func (u *SomaUtil) GetUserIdByName(c *resty.Client, user string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			User: &proto.UserFilter{
				UserName: user,
			},
		},
	}

	resp := u.PostRequestWithBody(c, req, "/filter/users/")
	userResult := u.DecodeProtoResultUserFromResponse(resp)

	if user != (*userResult.Users)[0].UserName {
		u.Abort("Received result set for incorrect user")
	}
	return (*userResult.Users)[0].Id
}

func (u *SomaUtil) DecodeProtoResultUserFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
