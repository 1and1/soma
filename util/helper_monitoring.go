package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetMonitoringByUUIDOrName(s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetMonitoringIdByName(s)
	}
	return id.String()
}

func (u *SomaUtil) GetMonitoringIdByName(monitoring string) string {
	req := somaproto.ProtoRequestMonitoring{}
	req.Filter = &somaproto.ProtoMonitoringFilter{}
	req.Filter.Name = monitoring

	resp := u.PostRequestWithBody(req, "/filter/monitoring/")
	monitoringResult := u.DecodeProtoResultMonitoringFromResponse(resp)

	if monitoring != monitoringResult.Systems[0].Name {
		u.Abort("Received result set for incorrect monitoring system")
	}
	return monitoringResult.Systems[0].Id
}

func (u *SomaUtil) DecodeProtoResultMonitoringFromResponse(resp *resty.Response) *somaproto.ProtoResultMonitoring {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	var res somaproto.ProtoResultMonitoring
	err := decoder.Decode(&res)
	u.AbortOnError(err, "Error decoding server response body")
	if res.Code > 299 {
		s := fmt.Sprintf("Request failed: %d - %s", res.Code, res.Status)
		msgs := []string{s}
		msgs = append(msgs, res.Text...)
		u.Abort(msgs...)
	}
	return &res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
