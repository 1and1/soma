package util

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	req := somaproto.ProtoRequestCapability{}
	req.Filter = &somaproto.ProtoCapabilityFilter{}

	split := strings.SplitN(capability, ".", 3)
	if len(split) != 3 {
		u.Abort("Split failed, Capability name invalid")
	}
	req.Filter.Monitoring = u.TryGetMonitoringByUUIDOrName(split[0])
	req.Filter.View = split[1]
	req.Filter.Metric = split[2]

	resp := u.PostRequestWithBody(req, "/filter/capability/")
	capabilityResult := u.DecodeProtoResultCapabilityFromResponse(resp)

	if capability != capabilityResult.Capabilities[0].Name {
		u.Abort("Received result set for incorrect capability")
	}
	return capabilityResult.Capabilities[0].Id
}

func (u *SomaUtil) DecodeProtoResultCapabilityFromResponse(resp *resty.Response) *somaproto.ProtoResultCapability {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	var res somaproto.ProtoResultCapability
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
