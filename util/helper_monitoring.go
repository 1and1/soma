package util

import (
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetMonitoringByUUIDOrName(c *resty.Client, s string) string {
	if u.IsUUID(s) {
		return s
	}
	return u.GetMonitoringIdByName(c, s)
}

func (u *SomaUtil) GetMonitoringIdByName(c *resty.Client, monitoring string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Monitoring: &proto.MonitoringFilter{
				Name: monitoring,
			},
		},
	}
	resp := u.PostRequestWithBody(c, req, "/filter/monitoring/")
	monitoringResult := u.DecodeProtoResultMonitoringFromResponse(resp)

	if monitoring != (*monitoringResult.Monitorings)[0].Name {
		u.Abort("Received result set for incorrect monitoring system")
	}
	return (*monitoringResult.Monitorings)[0].Id
}

func (u *SomaUtil) DecodeProtoResultMonitoringFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
