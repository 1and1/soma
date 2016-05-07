package util

import (
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
		req         somaproto.PropertyRequest
		ctxIdString string
		path        string
	)
	req = somaproto.PropertyRequest{}
	req.Filter = &somaproto.PropertyFilter{}
	req.Filter.Name = prop

	switch pType {
	case "custom":
		// context ctx is repository
		ctxIdString = u.TryGetRepositoryByUUIDOrName(ctx)
		path = fmt.Sprintf("/property/custom/%s/", ctxIdString)
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
	case propResult.Custom[0].Name:
		return propResult.Custom[0].CustomId
	case propResult.System[0].Name:
		return propResult.System[0].Name
	case propResult.Service[0].Name:
		return propResult.Service[0].Name
	default:
		u.Abort("Received result set for incorrect property")
	}

	// required to silence the compiler, since ending in a switch is not
	// analyzed to always return:
	// http://code.google.com/p/go/issues/detail?id=65
	panic("unreachable")
}

func (u SomaUtil) CheckStringIsSystemProperty(s string) {
	resp := u.GetRequest("/property/system/")
	res := u.DecodeProtoResultPropertyFromResponse(resp)

	for _, prop := range res.System {
		if prop.Name == s {
			return
		}
	}
	u.Abort("Invalid system property requested")
}

func (u SomaUtil) DecodeProtoResultPropertyFromResponse(resp *resty.Response) *somaproto.Result {
	return DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
