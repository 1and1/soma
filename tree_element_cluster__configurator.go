package somatree

import "github.com/satori/go.uuid"

func (tec *SomaTreeElemCluster) updateCheckInstances() {
	// object has no checks
	if len(tec.Checks) == 0 {
		return
	}

	// process checks
checksloop:
	for i, _ := range tec.Checks {
		if tec.Checks[i].Inherited == false && tec.Checks[i].ChildrenOnly == true {
			continue checksloop
		}
		if tec.Checks[i].View == "local" {
			continue checksloop
		}
		hasBrokenConstraint := false
		hasServiceConstraint := false
		hasAttributeConstraint := false
		view := tec.Checks[i].View

		attributes := []CheckConstraint{}
		oncallC := ""                                  // Id
		systemC := map[string]string{}                 // Id->Value
		nativeC := map[string]string{}                 // Property->Value
		serviceC := map[string]string{}                // Id->Value
		customC := map[string]string{}                 // Id->Value
		attributeC := map[string]map[string][]string{} // svcId->attr->[ value, ... ]

		newInstances := map[string]CheckInstance{}
		newCheckInstances := []string{}

		// these constaint types must always match for the instance to
		// be valid. defer service and attribute
	constraintcheck:
		for _, c := range tec.Checks[i].Constraints {
			switch c.Type {
			case "native":
				if tec.evalNativeProp(c.Key, c.Value) {
					nativeC[c.Key] = c.Value
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "system":
				if id, hit := tec.evalSystemProp(c.Key, c.Value, view); hit {
					systemC[id] = c.Value
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "oncall":
				if id, hit := tec.evalOncallProp(c.Key, c.Value, view); hit {
					oncallC = id
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "custom":
				if id, hit := tec.evalCustomProp(c.Key, c.Value, view); hit {
					customC[id] = c.Value
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "service":
				hasServiceConstraint = true
				if id, hit := tec.evalServiceProp(c.Key, c.Value, view); hit {
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

		/* if the check has both service and attribute constraints,
		* then for the check to hit, the tree element needs to have
		* all the services, and each of them needs to match all
		* attribute constraints
		 */
		if hasServiceConstraint && hasAttributeConstraint {
		svcattrloop:
			for id, _ := range serviceC {
				for _, attr := range attributes {
					hit := tec.evalAttributeOfService(id, view, attr.Key, attr.Value)
					if hit {
						attributeC[id][attr.Key] = append(attributeC[id][attr.Key], attr.Value)
					} else {
						hasBrokenConstraint = true
						break svcattrloop
					}
				}
			}
			/* if the check has only attribute constraints and no
			* service constraint, then we pull in every service that
			* matches all attribute constraints and generate a check
			* instance for it
			 */
		} else if hasAttributeConstraint {
			attrCount := len(attributes)
			for _, attr := range attributes {
				hit, svcIdMap := tec.evalAttributeProp(view, attr.Key, attr.Value)
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
			// declare service constraints in effect if we found a
			// service that bound all attribute constraints
			if len(serviceC) > 0 {
				hasServiceConstraint = true
			}
		}
		if hasBrokenConstraint {
			continue checksloop
		}
		// check triggered, create instances

		/* if there are no service constraints, one check instance is
		* created for this check
		 */
		if !hasServiceConstraint {
			inst := CheckInstance{
				InstanceId: uuid.UUID{},
				CheckId: func(id string) uuid.UUID {
					f, _ := uuid.FromString(id)
					return f
				}(i),
				InstanceConfigId:      uuid.NewV4(),
				ConstraintOncall:      oncallC,
				ConstraintService:     serviceC,
				ConstraintSystem:      systemC,
				ConstraintCustom:      customC,
				ConstraintNative:      nativeC,
				ConstraintAttribute:   attributeC,
				InstanceService:       "",
				InstanceServiceConfig: nil,
				InstanceSvcCfgHash:    "",
			}
			inst.calcConstraintHash()
			inst.calcConstraintValHash()

		nosvcinstanceloop:
			for _, exInstId := range tec.CheckInstances[i] {
				exInst := tec.Instances[exInstId]
				// ignore instances with service constraints
				if exInst.InstanceSvcCfgHash != "" {
					continue nosvcinstanceloop
				}
				// check if an instance exists bound against the same
				// constraints
				if exInst.ConstraintHash == inst.ConstraintHash {
					inst.InstanceId, _ = uuid.FromString(exInst.InstanceId.String())
					inst.Version = exInst.Version + 1
					break nosvcinstanceloop
				}
			}
			if uuid.Equal(uuid.Nil, inst.InstanceId) {
				// no match was found during nosvcinstanceloop, this
				// is a new instance
				inst.Version = 0
				inst.InstanceId = uuid.NewV4()
			}
			newInstances[inst.InstanceId.String()] = inst
			newCheckInstances = append(newCheckInstances, inst.InstanceId.String())
		}

		/* if service constraints are in effect, then we generate
		* instances for every service that bound.
		* Since service attributes can be specified more than once,
		* but the semantics are unclear what the expected behaviour of
		* for example a file age check is that is specified against
		* more than one file path; all possible attribute value
		* permutations for each service are built and then one check
		* instance is built for each of these service config
		* permutations.
		 */
	serviceconstraintloop:
		for svcId, _ := range serviceC {
			if !hasServiceConstraint {
				break serviceconstraintloop
			}

			svcCfg := tec.getServiceMap(svcId)

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
			// build a CheckInstance for every result
			for _, y := range results {
				// ensure we have a full copy and not a header copy
				cfg := map[string]string{}
				for k, v := range y {
					cfg[k] = v
				}
				inst := CheckInstance{
					InstanceId: uuid.UUID{},
					CheckId: func(id string) uuid.UUID {
						f, _ := uuid.FromString(id)
						return f
					}(i),
					InstanceConfigId:      uuid.NewV4(),
					ConstraintOncall:      oncallC,
					ConstraintService:     serviceC,
					ConstraintSystem:      systemC,
					ConstraintCustom:      customC,
					ConstraintNative:      nativeC,
					ConstraintAttribute:   attributeC,
					InstanceService:       svcId,
					InstanceServiceConfig: cfg,
				}
				inst.calcConstraintHash()
				inst.calcConstraintValHash()
				inst.calcInstanceSvcCfgHash()

				// lookup existing instance ids for check in tec.CheckInstances
				// to determine if this is an update
			instanceloop:
				for _, exInstId := range tec.CheckInstances[i] {
					exInst := tec.Instances[exInstId]
					// this existing instance is for the same service
					// configuration -> this is an update
					if exInst.InstanceSvcCfgHash == inst.InstanceSvcCfgHash {
						inst.InstanceId, _ = uuid.FromString(exInst.InstanceId.String())
						inst.Version = exInst.Version + 1
						break instanceloop
					}
				}
				if uuid.Equal(uuid.Nil, inst.InstanceId) {
					// no match was found during instanceloop, this is
					// a new instance
					inst.Version = 0
					inst.InstanceId = uuid.NewV4()
				}
				newInstances[inst.InstanceId.String()] = inst
				newCheckInstances = append(newCheckInstances, inst.InstanceId.String())
			}
		} // LOOPEND: range serviceC
		// all new check instances have been built, check which
		// existing instances did not get an update and need to be
		// deleted
		for _, oldInstanceId := range tec.CheckInstances[i] {
			if _, ok := newInstances[oldInstanceId]; !ok {
				// there is no new version for this instance id
				tec.actionCheckInstanceDelete(tec.Instances[oldInstanceId].MakeAction())
				delete(tec.Instances, oldInstanceId)
				continue
			}
			delete(tec.Instances, oldInstanceId)
			tec.Instances[oldInstanceId] = newInstances[oldInstanceId]
			tec.actionCheckInstanceUpdate(tec.Instances[oldInstanceId].MakeAction())
		}
		for _, newInstanceId := range newCheckInstances {
			if _, ok := tec.Instances[newInstanceId]; !ok {
				// this instance is new, not an update
				tec.Instances[newInstanceId] = newInstances[newInstanceId]
				tec.actionCheckInstanceCreate(tec.Instances[newInstanceId].MakeAction())
			}
		}
		delete(tec.CheckInstances, i)
		tec.CheckInstances[i] = newCheckInstances
	} // LOOPEND: range tec.Checks
}

func (tec *SomaTreeElemCluster) evalNativeProp(
	prop string, val string) bool {
	switch prop {
	case "environment":
		env := tec.Parent.(Bucketeer).GetEnvironment()
		if val == env {
			return true
		}
	case "object_type":
		if val == "node" {
			return true
		}
	case "object_state":
		if val == tec.State {
			return true
		}
	}
	return false
}

func (tec *SomaTreeElemCluster) evalSystemProp(
	prop string, val string, view string) (string, bool) {
	for _, v := range tec.PropertySystem {
		t := v.(*PropertySystem)
		if t.Key == prop && t.Value == val && t.View == view {
			return t.Key, true
		}
	}
	return "", false
}

func (tec *SomaTreeElemCluster) evalOncallProp(
	prop string, val string, view string) (string, bool) {
	for _, v := range tec.PropertyOncall {
		t := v.(*PropertyOncall)
		if "name" == prop && t.Name == val && t.View == view {
			return t.Name, true
		}
	}
	return "", false
}

func (tec *SomaTreeElemCluster) evalCustomProp(
	prop string, val string, view string) (string, bool) {
	for _, v := range tec.PropertyCustom {
		t := v.(*PropertyCustom)
		if t.Key == prop && t.Value == val && t.View == view {
			return t.Key, true
		}
	}
	return "", false
}

func (tec *SomaTreeElemCluster) evalServiceProp(
	prop string, val string, view string) (string, bool) {
	for _, v := range tec.PropertyService {
		t := v.(*PropertyService)
		if prop == "name" && t.Service == val && t.View == view {
			return t.View, true
		}
	}
	return "", false
}

func (tec *SomaTreeElemCluster) evalAttributeOfService(
	svcName string, view string, attribute string, value string) bool {
	for _, v := range tec.PropertyService {
		t := v.(*PropertyService)
		if t.Service != svcName {
			continue
		}
		for _, a := range t.Attributes {
			if a.Attribute == attribute && t.View == view && a.Value == value {
				return true
			}
		}
	}
	return false
}

func (tec *SomaTreeElemCluster) evalAttributeProp(
	view string, attr string, value string) (bool, map[string]string) {
	f := map[string]string{}
	for _, v := range tec.PropertyService {
		t := v.(*PropertyService)
		for _, a := range t.Attributes {
			if a.Attribute == attr && a.Value == value && t.View == view {
				f[t.Service] = a.Attribute
			}
		}
	}
	if len(f) > 0 {
		return true, f
	}
	return false, f
}

func (tec *SomaTreeElemCluster) getServiceMap(serviceId string) map[string][]string {
	svc := new(PropertyService)
	svc = tec.PropertyService[serviceId].(*PropertyService)

	res := map[string][]string{}
	for _, v := range svc.Attributes {
		res[v.Attribute] = append(res[v.Attribute], v.Value)
	}
	return res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
