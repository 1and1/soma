package util

import (
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"

)

func (u SomaUtil) GetObjectIdForCheck(t string, n string, b string) string {
	switch t {
	case "repository":
		return u.TryGetRepositoryByUUIDOrName(n)
	case "bucket":
		return u.BucketByUUIDOrName(n)
	case "group":
		return u.TryGetGroupByUUIDOrName(n, b)
	case "cluster":
		return u.TryGetClusterByUUIDOrName(n, b)
	case "node":
		return u.TryGetNodeByUUIDOrName(n)
	default:
		u.Abort(fmt.Sprintf("Error, unknown object type: %s", t))
	}
	return ""
}

func (u SomaUtil) CleanThresholds(thresholds []somaproto.CheckConfigThreshold) []somaproto.CheckConfigThreshold {
	clean := []somaproto.CheckConfigThreshold{}

	for _, thr := range thresholds {
		c := somaproto.CheckConfigThreshold{
			Value: thr.Value,
			Predicate: somaproto.Predicate{
				Symbol: thr.Predicate.Symbol,
			},
			Level: somaproto.Level{
				Name: u.TryGetLevelNameByNameOrShort(thr.Level.Name),
			},
		}
		clean = append(clean, c)
	}
	return clean
}

func (u SomaUtil) CleanConstraints(constraints []somaproto.CheckConfigConstraint, repoId string, teamId string) []somaproto.CheckConfigConstraint {
	clean := []somaproto.CheckConfigConstraint{}

	for _, prop := range constraints {
		switch prop.ConstraintType {
		case "native":
			resp := u.GetRequest(fmt.Sprintf("/property/native/%s", prop.Native.Name))
			_ = u.DecodeProtoResultPropertyFromResponse(resp) // aborts on 404
			clean = append(clean, prop)
		case "system":
			resp := u.GetRequest(fmt.Sprintf("/property/system/%s", prop.System.Name))
			_ = u.DecodeProtoResultPropertyFromResponse(resp) // aborts on 404
			clean = append(clean, prop)
		case "attribute":
			resp := u.GetRequest(fmt.Sprintf("/attributes/%s", prop.Attribute.Name))
			_ = u.DecodeProtoResultAttributeFromResponse(resp) // aborts on 404
			clean = append(clean, prop)
		case "oncall":
			oc := somaproto.PropertyOncall{}
			if prop.Oncall.Name != "" {
				oc.Id = u.TryGetOncallByUUIDOrName(prop.Oncall.Name)
			} else if prop.Oncall.Id != "" {
				oc.Id = u.TryGetOncallByUUIDOrName(prop.Oncall.Id)
			}
			clean = append(clean, somaproto.CheckConfigConstraint{
				ConstraintType: prop.ConstraintType,
				Oncall:         &oc,
			})
		case "service":
			so := somaproto.PropertyService{
				Name:   u.TryGetServicePropertyByUUIDOrName(prop.Service.Name, teamId),
				TeamId: teamId,
			}
			clean = append(clean, somaproto.CheckConfigConstraint{
				ConstraintType: prop.ConstraintType,
				Service:        &so,
			})
		case "custom":
			co := somaproto.PropertyCustom{
				RepositoryId: repoId,
				Id:           u.TryGetCustomPropertyByUUIDOrName(prop.Custom.Name, repoId),
				Value:        prop.Custom.Value,
			}
			clean = append(clean, somaproto.CheckConfigConstraint{
				ConstraintType: prop.ConstraintType,
				Custom:         &co,
			})
		}
	}
	return clean
}

func (u *SomaUtil) TryGetCheckByUUIDOrName(c string, r string) string {
	id, err := uuid.FromString(c)
	if err != nil {
		return u.GetCheckByName(c, r)
	}
	return id.String()
}

func (u *SomaUtil) GetCheckByName(c string, r string) string {
	repo := u.TryGetRepositoryByUUIDOrName(r)
	req := somaproto.Request{
		Filter: &somaproto.Filter{
			CheckConfig: &somaproto.CheckConfigFilter{
				Name: c,
			},
		},
	}

	path := fmt.Sprintf("/filter/checks/%s/", repo)
	resp := u.PostRequestWithBody(req, path)
	checkResult := u.DecodeCheckConfigurationResultFromResponse(resp)

	if c != (*checkResult.CheckConfigs)[0].Name {
		u.Abort("Received result set for incorrect check configuration")
	}
	return (*checkResult.CheckConfigs)[0].Id
}

func (u *SomaUtil) DecodeCheckConfigurationResultFromResponse(resp *resty.Response) *somaproto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
