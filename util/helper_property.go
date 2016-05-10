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
		req         proto.Request
		ctxIdString string
		path        string
	)
	req = proto.Request{
		Filter: &proto.Filter{
			Property: &proto.PropertyFilter{
				Name: prop,
			},
		},
	}

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
	case (*propResult.Properties)[0].Custom.Name:
		return (*propResult.Properties)[0].Custom.Id
	case (*propResult.Properties)[0].System.Name:
		return (*propResult.Properties)[0].System.Name
	case (*propResult.Properties)[0].Service.Name:
		return (*propResult.Properties)[0].Service.Name
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

	for _, prop := range *res.Properties {
		if prop.System.Name == s {
			return
		}
	}
	u.Abort(fmt.Sprintf("Invalid system property requested: %s", s))
}

func (u SomaUtil) DecodeProtoResultPropertyFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
