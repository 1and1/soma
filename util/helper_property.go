package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetCustomPropertyByUUIDOrName(s string, r string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetPropertyIdByName("custom", s, r)
	}
	return id.String()
}

func (u SomaUtil) TryGetServicePropertyByUUIDOrName(s string, t string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetPropertyIdByName("service", s, t)
	}
	return id.String()
}

func (u SomaUtil) TryGetSystemPropertyByUUIDOrName(s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetPropertyIdByName("system", s, "none")
	}
	return id.String()
}

func (u SomaUtil) TryGetTemplatePropertyByUUIDOrName(s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetPropertyIdByName("template", s, "none")
	}
	return id.String()
}

func (u SomaUtil) GetPropertyIdByName(pType string, prop string, ctx string) string {
	var (
		req         somaproto.ProtoRequestProperty
		ctxId       uuid.UUID
		ctxIdString string
		path        string
	)
	req.Filter.Property = prop

	switch pType {
	case "custom":
		// context ctx is repository
		ctxId = u.TryGetRepositoryByUUIDOrName(ctx)
		path = fmt.Sprintf("/property/custom/%s/", ctxId.String())
	case "system":
		path = "/property/system/"
	case "template":
		path = "/property/service/global/"
	case "service":
		// context ctx is team
		ctxIdString = u.TryGetTeamByUUIDOrName(ctx)
		path = fmt.Sprintf("/property/service/team/%s/", ctxIdString)
	default:
		u.Abort("Unsupported property type in util.GetPropertyIdByName()")
	}

	resp := u.GetRequestWithBody(req, path)
	propResult := u.DecodeProtoResultPropertyFromResponse(resp)

	switch prop {
	case propResult.Custom[0].Property:
		return propResult.Custom[0].Id
	case propResult.System[0].Property:
		return propResult.System[0].Property
	case propResult.Service[0].Property:
		return propResult.Service[0].Property
	default:
		u.Abort("Received result set for incorrect property")
	}

	// required to silence the compiler, since ending in a switch is not
	// analyzed to always return:
	// http://code.google.com/p/go/issues/detail?id=65
	panic("unreachable")
}

func (u SomaUtil) DecodeProtoResultPropertyFromResponse(resp *resty.Response) *somaproto.ProtoResultProperty {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	var res somaproto.ProtoResultProperty
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
