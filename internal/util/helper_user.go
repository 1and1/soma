package util

import (
	"github.com/1and1/soma/lib/proto"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) tryGetUserByUUIDOrName(c *resty.Client, s string) string {
	if u.IsUUID(s) {
		return s
	}
	return u.getUserIdByName(c, s)
}

func (u *SomaUtil) getUserIdByName(c *resty.Client, user string) string {
	req := proto.NewUserFilter()
	req.Filter.User.UserName = user

	resp := u.PostRequestWithBody(c, req, "/filter/users/")
	res := u.DecodeResultFromResponse(resp)

	if user != (*res.Users)[0].UserName {
		u.abort("Received result set for incorrect user")
	}
	return (*res.Users)[0].Id
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
