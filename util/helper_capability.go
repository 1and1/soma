package util

import (
	"strings"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetCapabilityByUUIDOrName(s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetCapabilityIdByName(s)
	}
	return id.String()
}

func (u *SomaUtil) GetCapabilityIdByName(capability string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Capability: &proto.CapabilityFilter{},
		},
	}

	split := strings.SplitN(capability, ".", 3)
	if len(split) != 3 {
		u.Abort("Split failed, Capability name invalid")
	}
	req.Filter.Capability.MonitoringId = u.TryGetMonitoringByUUIDOrName(split[0])
	req.Filter.Capability.View = split[1]
	req.Filter.Capability.Metric = split[2]

	resp := u.PostRequestWithBody(req, "/filter/capability/")
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
