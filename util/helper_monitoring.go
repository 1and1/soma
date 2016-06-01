package util

import (
	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetMonitoringByUUIDOrName(c *resty.Client, s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetMonitoringIdByName(c, s)
	}
	return id.String()
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
