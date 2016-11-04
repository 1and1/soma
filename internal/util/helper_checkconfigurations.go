package util

import (
	"fmt"

	"github.com/1and1/soma/lib/proto"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) GetObjectIdForCheck(c *resty.Client, t string, n string, b string) string {
	switch t {
	case "repository":
		return u.tryGetRepositoryByUUIDOrName(c, n)
	case "bucket":
		return u.bucketByUUIDOrName(c, n)
	case "group":
		return u.tryGetGroupByUUIDOrName(c, n, b)
	case "cluster":
		return u.tryGetClusterByUUIDOrName(c, n, b)
	case "node":
		return u.tryGetNodeByUUIDOrName(c, n)
	default:
		u.abort(fmt.Sprintf("Error, unknown object type: %s", t))
	}
	return ""
}

func (u SomaUtil) CleanThresholds(c *resty.Client, thresholds []proto.CheckConfigThreshold) []proto.CheckConfigThreshold {
	clean := []proto.CheckConfigThreshold{}

	for _, thr := range thresholds {
		c := proto.CheckConfigThreshold{
			Value: thr.Value,
			Predicate: proto.Predicate{
				Symbol: thr.Predicate.Symbol,
			},
			Level: proto.Level{
				Name: u.tryGetLevelNameByNameOrShort(c, thr.Level.Name),
			},
		}
		clean = append(clean, c)
	}
	return clean
}

func (u SomaUtil) CleanConstraints(c *resty.Client, constraints []proto.CheckConfigConstraint, repoId string, teamId string) []proto.CheckConfigConstraint {
	clean := []proto.CheckConfigConstraint{}

	for _, prop := range constraints {
		switch prop.ConstraintType {
		case "native":
			resp := u.GetRequest(c, fmt.Sprintf("/property/native/%s", prop.Native.Name))
			_ = u.DecodeProtoResultPropertyFromResponse(resp) // aborts on 404
			clean = append(clean, prop)
		case "system":
			resp := u.GetRequest(c, fmt.Sprintf("/property/system/%s", prop.System.Name))
			_ = u.DecodeProtoResultPropertyFromResponse(resp) // aborts on 404
			clean = append(clean, prop)
		case "attribute":
			resp := u.GetRequest(c, fmt.Sprintf("/attributes/%s", prop.Attribute.Name))
			_ = u.DecodeProtoResultAttributeFromResponse(resp) // aborts on 404
			clean = append(clean, prop)
		case "oncall":
			oc := proto.PropertyOncall{}
			if prop.Oncall.Name != "" {
				oc.Id = u.tryGetOncallByUUIDOrName(c, prop.Oncall.Name)
			} else if prop.Oncall.Id != "" {
				oc.Id = u.tryGetOncallByUUIDOrName(c, prop.Oncall.Id)
			}
			clean = append(clean, proto.CheckConfigConstraint{
				ConstraintType: prop.ConstraintType,
				Oncall:         &oc,
			})
		case "service":
			so := proto.PropertyService{
				Name:   u.TryGetServicePropertyByUUIDOrName(c, prop.Service.Name, teamId),
				TeamId: teamId,
			}
			clean = append(clean, proto.CheckConfigConstraint{
				ConstraintType: prop.ConstraintType,
				Service:        &so,
			})
		case "custom":
			co := proto.PropertyCustom{
				RepositoryId: repoId,
				Id:           u.TryGetCustomPropertyByUUIDOrName(c, prop.Custom.Name, repoId),
				Value:        prop.Custom.Value,
			}
			clean = append(clean, proto.CheckConfigConstraint{
				ConstraintType: prop.ConstraintType,
				Custom:         &co,
			})
		}
	}
	return clean
}

func (u *SomaUtil) TryGetCheckByUUIDOrName(c *resty.Client, ck string, r string) string {
	if u.isUUID(ck) {
		return ck
	}
	return u.getCheckByName(c, ck, r)
}

func (u *SomaUtil) getCheckByName(c *resty.Client, ck string, r string) string {
	repo := u.tryGetRepositoryByUUIDOrName(c, r)
	req := proto.Request{
		Filter: &proto.Filter{
			CheckConfig: &proto.CheckConfigFilter{
				Name: ck,
			},
		},
	}

	path := fmt.Sprintf("/filter/checks/%s/", repo)
	resp := u.PostRequestWithBody(c, req, path)
	checkResult := u.DecodeCheckConfigurationResultFromResponse(resp)

	if ck != (*checkResult.CheckConfigs)[0].Name {
		u.abort("Received result set for incorrect check configuration")
	}
	return (*checkResult.CheckConfigs)[0].Id
}

func (u *SomaUtil) DecodeCheckConfigurationResultFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
