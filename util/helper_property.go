package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetCustomPropertyByUUIDOrName(s string, r string) uuid.UUID {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		id = u.GetPropertyIdByName("custom", s, r)
	}
	return id
}

func (u SomaUtil) TryGetServicePropertyByUUIDOrName(s string, t string) uuid.UUID {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		id = u.GetPropertyIdByName("service", s, t)
	}
	return id
}

func (u SomaUtil) TryGetSystemPropertyByUUIDOrName(s string) uuid.UUID {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		id = u.GetPropertyIdByName("system", s, "none")
	}
	return id
}

func (u SomaUtil) TryGetTemplatePropertyByUUIDOrName(s string) uuid.UUID {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		id = u.GetPropertyIdByName("template", s, "none")
	}
	return id
}

func (u SomaUtil) GetPropertyIdByName(pType string, prop string, ctx string) uuid.UUID {
	var (
		req   somaproto.ProtoRequestProperty
		ctxId uuid.UUID
		path  string
	)
	req.Filter.Name = prop

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
		ctxId = u.TryGetTeamByUUIDOrName(ctx)
		path = fmt.Sprintf("/property/service/team/%s/", ctxId.String())
	default:
		u.Abort("Unsupported property type in util.GetPropertyIdByName()")
	}

	resp := u.GetRequestWithBody(req, path)
	propResult := u.DecodeProtoResultPropertyFromResponse(resp)

	switch prop {
	case propResult.Custom[0].Property:
		return propResult.Custom[0].Id
	case propResult.System[0].Property:
		return propResult.System[0].Id
	case propResult.Service[0].Property:
		return propResult.Service[0].Id
	default:
		u.Abort("Received result set for incorrect property")
	}

	// required to silence the compiler, since ending in a switch is not
	// analyzed to always return:
	// http://code.google.com/p/go/issues/detail?id=65
	panic("unreachable")
}

func (u SomaUtil) DecodeProtoResultPropertyFromResponse(resp *resty.Response) *somaproto.ProtoResultProperty {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body))
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
