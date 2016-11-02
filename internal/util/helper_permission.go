package util

import (
	"github.com/1and1/soma/lib/proto"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetPermissionByUUIDOrName(c *resty.Client, s string) string {
	if u.IsUUID(s) {
		return s
	}
	return u.getPermissionIdByName(c, s)
}

func (u *SomaUtil) getPermissionIdByName(c *resty.Client, perm string) string {
	req := proto.NewPermissionFilter()
	req.Filter.Permission.Name = perm

	resp := u.PostRequestWithBody(c, req, `/filter/permission/`)
	res := u.DecodeResultFromResponse(resp)

	if perm != (*res.Permissions)[0].Name {
		u.abort(`Received result set for incorrect permission`)
	}
	return (*res.Permissions)[0].Id
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
