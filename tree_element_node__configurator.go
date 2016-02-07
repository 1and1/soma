package somatree

import "github.com/satori/go.uuid"

func (ten *SomaTreeElemNode) updateCheckInstances() {
	// object has no checks
	if len(ten.Checks) == 0 {
		return
	}

	// process checks
checksloop:
	for i, _ := range ten.Checks {
		if ten.Checks[i].Inherited == false && ten.Checks[i].ChildrenOnly == true {
			continue checksloop
		}
		hasBrokenConstraint := false
		hasServiceConstraint := false
		hasAttributeConstraint := false
		view := ten.Checks[i].View

		attributes := []SomaTreeCheckConstraint{}
		oncallC := ""                                  // Id
		systemC := map[string]string{}                 // Id->Value
		nativeC := map[string]string{}                 // Property->Value
		serviceC := map[string]string{}                // Id->Value
		customC := map[string]string{}                 // Id->Value
		attributeC := map[string]map[string][]string{} // svcId->attr->[ value, ... ]

		// these constaint types must always match for the instance to
		// be valid. defer service and attribute
	constraintcheck:
		for _, c := range ten.Checks[i].Constraints {
			switch c.Type {
			case "native":
				if ten.evalNativeProp(c.Key, c.Value) {
					nativeC[c.Key] = c.Value
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "system":
				if id, hit := ten.evalSystemProp(c.Key, c.Value, view); hit {
					systemC[id] = c.Value
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "oncall":
				if id, hit := ten.evalOncallProp(c.Key, c.Value, view); hit {
					oncallC = id
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "custom":
				if id, hit := ten.evalCustomProp(c.Key, c.Value, view); hit {
					customC[id] = c.Value
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "service":
				hasServiceConstraint = true
				if id, hit := ten.evalServiceProp(c.Key, c.Value, view); hit {
					serviceC[id] = c.Value
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "attribute":
				hasAttributeConstraint = true
				attributes = append(attributes, c)
			}
		}
		if hasBrokenConstraint {
			continue checksloop
		}

		if hasServiceConstraint && hasAttributeConstraint {
		svcattrloop:
			for id, _ := range serviceC {
				for _, attr := range attributes {
					hit := ten.evalAttributeOfService(id, view, attr.Key, attr.Value)
					if hit {
						attributeC[id][attr.Key] = append(attributeC[id][attr.Key], attr.Value)
					} else {
						hasBrokenConstraint = true
						break svcattrloop
					}
				}
			}
		} else if hasAttributeConstraint {
			attrCount := len(attributes)
			for _, attr := range attributes {
				hit, svcIdMap := ten.evalAttributeProp(view, attr.Key, attr.Value)
				if hit {
					for id, _ := range svcIdMap {
						serviceC[id] = svcIdMap[id]
						attributeC[id][attr.Key] = append(attributeC[id][attr.Key], attr.Value)
					}
				}
			}
			// delete all services that did not match all attributes
			for id, _ := range attributeC {
				if len(attributeC[id]) != attrCount {
					delete(serviceC, id)
					delete(attributeC, id)
				}
			}
			if len(serviceC) > 0 {
				hasServiceConstraint = true
			}
		}
		if hasBrokenConstraint {
			continue checksloop
		}
		// check triggered, create instances

		if !hasServiceConstraint {
			// TODO create 1x
			continue checksloop
		}

		for svcId, _ := range serviceC {
			svcCfg := ten.getServiceMap(svcId)

			// calculate how many instances this service spawns
			combinations := 1
			for attr, _ := range svcCfg {
				combinations = combinations * len(svcCfg[attr])
			}

			// build all attribute combinations
			results := make([]map[string]string, 0, combinations)
			for attr, _ := range svcCfg {
				if len(results) == 0 {
					for i, _ := range svcCfg[attr] {
						res := map[string]string{}
						res[attr] = svcCfg[attr][i]
						results = append(results, res)
					}
					continue
				}
				ires := make([]map[string]string, 0, combinations)
				for r, _ := range results {
					for j, _ := range svcCfg[attr] {
						res := map[string]string{}
						for k, v := range results[r] {
							res[k] = v
						}
						res[attr] = svcCfg[attr][j]
						ires = append(ires, res)
					}
				}
				results = ires
			}
			// build a SomaTreeCheckInstance for every result
			for _, y := range results {
				// ensure we have a full copy and not a header copy
				cfg := map[string]string{}
				for k, v := range y {
					cfg[k] = v
				}
				inst := SomaTreeCheckInstance{
					CheckId: func(id string) uuid.UUID {
						f, _ := uuid.FromString(id)
						return f
					}(i),
					ConstraintOncall:    oncallC,
					ConstraintService:   serviceC,
					ConstraintSystem:    systemC,
					ConstraintCustom:    customC,
					ConstraintNative:    nativeC,
					ConstraintAttribute: attributeC,
					InstanceService: func(id string) uuid.UUID {
						f, _ := uuid.FromString(id)
						return f
					}(svcId),
					InstanceServiceConfig: cfg,
				}
				inst.calcConstraintHash()
				inst.calcConstraintValHash()
				inst.calcInstanceSvcCfgHash()
				// TODO lookup existing instance ids for check in ten.CheckInstances
				// TODO check existing for same ConstraintHash
				// TODO ... same ConstraintValHash
				// TODO ... ... same InstanceSvcCfgHash --> instance update
				// TODO ... same InstanceSvcCfgHash     --> instance update (new constraints)
				// TODO add new instances
				// TODO remove old instances
				inst.Version = 0
				inst.CheckId, _ = uuid.FromString(i)
				inst.InstanceId = uuid.NewV4()
				ten.Instances[inst.InstanceId.String()] = inst
				ten.CheckInstances[i] = append(ten.CheckInstances[i], inst.InstanceId.String())
			}
		}
	}
}

func (ten *SomaTreeElemNode) evalNativeProp(
	prop string, val string) bool {
	switch prop {
	case "environment":
		env := ten.Parent.(SomaTreeBucketeer).GetEnvironment()
		if val == env {
			return true
		}
	case "object_type":
		if val == "node" {
			return true
		}
	case "object_state":
		if val == ten.State {
			return true
		}
	}
	return false
}

func (ten *SomaTreeElemNode) evalSystemProp(
	prop string, val string, view string) (string, bool) {
	// TODO
	return "", false
}

func (ten *SomaTreeElemNode) evalOncallProp(
	prop string, val string, view string) (string, bool) {
	// TODO
	return "", false
}

func (ten *SomaTreeElemNode) evalCustomProp(
	prop string, val string, view string) (string, bool) {
	// TODO
	return "", false
}

func (ten *SomaTreeElemNode) evalServiceProp(
	prop string, val string, view string) (string, bool) {
	// TODO
	return "", false
}

func (ten *SomaTreeElemNode) evalAttributeOfService(
	svcName string, view string, attribute string, value string) bool {
	// TODO
	return false
}

func (ten *SomaTreeElemNode) evalAttributeProp(
	view string, attr string, value string) (bool, map[string]string) {
	f := map[string]string{}
	// TODO
	return false, f
}

func (ten *SomaTreeElemNode) getServiceMap(serviceId string) map[string][]string {
	svc := new(PropertyService)
	svc = ten.PropertyService[serviceId].(*PropertyService)

	res := map[string][]string{}
	for _, v := range svc.Attributes {
		res[v.Attribute] = append(res[v.Attribute], v.Value)
	}
	return res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
