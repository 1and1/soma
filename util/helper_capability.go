package util

import (
	"strings"

	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetCapabilityByUUIDOrName(c *resty.Client, s string) string {
	if u.IsUUID(s) {
		return s
	}
	return u.GetCapabilityIdByName(c, s)
}

func (u *SomaUtil) GetCapabilityIdByName(c *resty.Client, capability string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Capability: &proto.CapabilityFilter{},
		},
	}

	split := strings.SplitN(capability, ".", 3)
	if len(split) != 3 {
		u.Abort("Split failed, Capability name invalid")
	}
	req.Filter.Capability.MonitoringId = u.TryGetMonitoringByUUIDOrName(c, split[0])
	req.Filter.Capability.View = split[1]
	req.Filter.Capability.Metric = split[2]

	resp := u.PostRequestWithBody(c, req, "/filter/capability/")
	capabilityResult := u.DecodeProtoResultCapabilityFromResponse(resp)

	if capability != (*capabilityResult.Capabilities)[0].Name {
		u.Abort("Received result set for incorrect capability")
	}
	return (*capabilityResult.Capabilities)[0].Id
}

func (u *SomaUtil) DecodeProtoResultCapabilityFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
