package util

import (
	"fmt"
	"strings"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/lib/proto"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetCapabilityByUUIDOrName(c *resty.Client, s string) string {
	if u.IsUUID(s) {
		return s
	}
	return u.GetCapabilityIdByName(c, s)
}

func (u *SomaUtil) GetCapabilityIdByName(c *resty.Client, capability string) string {
	req := proto.NewCapabilityFilter()

	split := strings.SplitN(capability, ".", 3)
	if len(split) != 3 {
		u.abort("Capability split failed, name invalid")
	}
	req.Filter.Capability.MonitoringId = u.TryGetMonitoringByUUIDOrName(c, split[0])
	req.Filter.Capability.View = split[1]
	req.Filter.Capability.Metric = split[2]

	resp, err := adm.PostReqBody(req, `/filter/capability/`)
	if err != nil {
		u.abort(fmt.Sprintf("Capability lookup request error: %s", err.Error()))
	}
	result, err := u.ResultFromResponse(resp)
	if se, ok := err.(SomaError); ok {
		if se.RequestError() {
			u.abort(fmt.Sprintf("Capability lookup request error: %s", se.Error()))
		}
		if se.Code() == 404 {
			u.abort(fmt.Sprintf(
				"Could not find capability with name %s",
				capability,
			))
		}
		u.abort(fmt.Sprintf("Capability lookup application error: %s", err.Error()))
	}

	if capability != (*result.Capabilities)[0].Name {
		u.abort(fmt.Sprintf(
			"Capability lookup failed. Wanted %s, received %s",
			capability,
			(*result.Capabilities)[0].Name,
		))
	}
	return (*result.Capabilities)[0].Id
}

func (u *SomaUtil) DecodeProtoResultCapabilityFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
