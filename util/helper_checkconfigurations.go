package util

import (
	"bytes"
	"encoding/json"
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

func (u SomaUtil) CleanThresholds(thresholds []somaproto.CheckConfigurationThreshold) []somaproto.CheckConfigurationThreshold {
	clean := []somaproto.CheckConfigurationThreshold{}

	for _, thr := range thresholds {
		c := somaproto.CheckConfigurationThreshold{
			Value: thr.Value,
			Predicate: somaproto.ProtoPredicate{
				Predicate: thr.Predicate.Predicate,
			},
			Level: somaproto.ProtoLevel{
				Name: u.TryGetLevelNameByNameOrShort(thr.Level.Name),
			},
		}
		clean = append(clean, c)
	}
	return clean
}

func (u SomaUtil) CleanConstraints(constraints []somaproto.CheckConfigurationConstraint, repoId string, teamId string) []somaproto.CheckConfigurationConstraint {
	clean := []somaproto.CheckConfigurationConstraint{}

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
			resp := u.GetRequest(fmt.Sprintf("/attributes/%s", prop.Attribute.Attribute))
			_ = u.DecodeProtoResultAttributeFromResponse(resp) // aborts on 404
			clean = append(clean, prop)
		case "oncall":
			oc := somaproto.TreePropertyOncall{}
			if prop.Oncall.Name != "" {
				oc.OncallId = u.TryGetOncallByUUIDOrName(prop.Oncall.Name)
			} else if prop.Oncall.OncallId != "" {
				oc.OncallId = u.TryGetOncallByUUIDOrName(prop.Oncall.OncallId)
			}
			clean = append(clean, somaproto.CheckConfigurationConstraint{
				ConstraintType: prop.ConstraintType,
				Oncall:         &oc,
			})
		case "service":
			so := somaproto.TreePropertyService{
				Name:   u.TryGetServicePropertyByUUIDOrName(prop.Service.Name, teamId),
				TeamId: teamId,
			}
			clean = append(clean, somaproto.CheckConfigurationConstraint{
				ConstraintType: prop.ConstraintType,
				Service:        &so,
			})
		case "custom":
			co := somaproto.TreePropertyCustom{
				RepositoryId: repoId,
				CustomId:     u.TryGetCustomPropertyByUUIDOrName(prop.Custom.Name, repoId),
				Value:        prop.Custom.Value,
			}
			clean = append(clean, somaproto.CheckConfigurationConstraint{
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
	req := somaproto.CheckConfigurationRequest{
		Filter: &somaproto.CheckConfigurationFilter{
			Name: c,
		},
	}

	path := fmt.Sprintf("/filter/checks/%s/", repo)
	resp := u.PostRequestWithBody(req, path)
	checkResult := u.DecodeCheckConfigurationResultFromResponse(resp)

	if c != checkResult.CheckConfigurations[0].Name {
		u.Abort("Received result set for incorrect check configuration")
	}
	return checkResult.CheckConfigurations[0].Id
}

func (u *SomaUtil) DecodeCheckConfigurationResultFromResponse(resp *resty.Response) *somaproto.CheckConfigurationResult {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	res := somaproto.CheckConfigurationResult{}
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
