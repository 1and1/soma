package util

import (
	"github.com/1and1/soma/lib/proto"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryResolveGrantId(c *resty.Client, rtyp, rid, pid, cat string) string {
	req := proto.NewGrantFilter()
	req.Filter.Grant.RecipientType = rtyp
	req.Filter.Grant.RecipientId = rid
	req.Filter.Grant.PermissionId = pid
	req.Filter.Grant.Category = cat

	resp := u.PostRequestWithBody(c, req, `/filter/grant/`)
	res := u.DecodeResultFromResponse(resp)

	if pid != (*res.Grants)[0].PermissionId {
		u.Abort(`Received result set for incorrect grant`)
	}
	return (*res.Grants)[0].Id
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
