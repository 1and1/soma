package somatree

import (
	"log"

	"github.com/satori/go.uuid"
)

func (teg *SomaTreeElemGroup) updateCheckInstances() {
	repoName := teg.GetRepositoryName()

	// object has no checks
	if len(teg.Checks) == 0 {
		log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, HasChecks=%t",
			repoName,
			`UpdateCheckInstances`,
			`node`,
			teg.Id.String(),
			false,
		)
		return
	}

	// if there are loaded instances, then this is the initial rebuild
	// of the tree
	startupLoad := false
	if len(teg.loadedInstances) > 0 {
		startupLoad = true
	}

	// process checks
checksloop:
	for i, _ := range teg.Checks {
		if teg.Checks[i].Inherited == false && teg.Checks[i].ChildrenOnly == true {
			continue checksloop
		}
		if teg.Checks[i].View == "local" {
			continue checksloop
		}
		hasBrokenConstraint := false
		hasServiceConstraint := false
		hasAttributeConstraint := false
		view := teg.Checks[i].View

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
		for _, c := range teg.Checks[i].Constraints {
			switch c.Type {
			case "native":
				if teg.evalNativeProp(c.Key, c.Value) {
					nativeC[c.Key] = c.Value
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "system":
				if id, hit := teg.evalSystemProp(c.Key, c.Value, view); hit {
					systemC[id] = c.Value
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "oncall":
				if id, hit := teg.evalOncallProp(c.Key, c.Value, view); hit {
					oncallC = id
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "custom":
				if id, hit := teg.evalCustomProp(c.Key, c.Value, view); hit {
					customC[id] = c.Value
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "service":
				hasServiceConstraint = true
				if id, hit := teg.evalServiceProp(c.Key, c.Value, view); hit {
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
			log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, Match=%t",
				repoName,
				`ConstraintEvaluation`,
				`group`,
				teg.Id.String(),
				i,
				false,
			)
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
					hit := teg.evalAttributeOfService(id, view, attr.Key, attr.Value)
					if hit {
						if attributeC[id] == nil {
							// attributeC[id] might still be a nil map
							attributeC[id] = map[string][]string{}
						}
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
				hit, svcIdMap := teg.evalAttributeProp(view, attr.Key, attr.Value)
				if hit {
					for id, _ := range svcIdMap {
						serviceC[id] = svcIdMap[id]
						if attributeC[id] == nil {
							// attributeC[id] might still be a nil map
							attributeC[id] = map[string][]string{}
						}
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
			log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, Match=%t",
				repoName,
				`ConstraintEvaluation`,
				`group`,
				teg.Id.String(),
				i,
				false,
			)
			continue checksloop
		}
		// check triggered, create instances
		log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, Match=%t",
			repoName,
			`ConstraintEvaluation`,
			`group`,
			teg.Id.String(),
			i,
			true,
		)

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
				ConfigId: func(id string) uuid.UUID {
					f, _ := uuid.FromString(teg.Checks[id].ConfigId.String())
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

			if startupLoad {
			nosvcstartinstanceloop:
				for ldInstId, ldInst := range teg.loadedInstances[i] {
					if ldInst.InstanceSvcCfgHash != "" {
						continue nosvcstartinstanceloop
					}
					// check if an instance exists bound against the
					// same constraints
					if ldInst.ConstraintHash == inst.ConstraintHash {
						inst.InstanceId, _ = uuid.FromString(ldInstId)
						if !uuid.Equal(ldInst.ConfigId, inst.ConfigId) {
							panic(`Matched instances loaded for different ConfigId`)
						}
						inst.InstanceConfigId, _ = uuid.FromString(ldInst.InstanceConfigId.String())
						inst.Version = ldInst.Version
						if inst.ConstraintValHash != ldInst.ConstraintValHash {
							panic(`Matched instances loaded for different ConstraintValHash`)
						}
						delete(teg.loadedInstances[i], ldInstId)
						log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s, ServiceConstrained=%t",
							repoName,
							`ComputeInstance`,
							`group`,
							teg.Id.String(),
							i,
							ldInstId,
							false,
						)
						break nosvcstartinstanceloop
					}
					// if we hit here, then we just computed an instance
					// that we could not match to any loaded instances
					// -> something is wrong
					panic(`Failed to match computed instance to loaded instances`)
				}
			} else {
			nosvcinstanceloop:
				for _, exInstId := range teg.CheckInstances[i] {
					exInst := teg.Instances[exInstId]
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
				log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s, ServiceConstrained=%t",
					repoName,
					`ComputeInstance`,
					`group`,
					teg.Id.String(),
					i,
					inst.InstanceId.String(),
					false,
				)
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

			svcCfg := teg.getServiceMap(svcId)

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
					ConfigId: func(id string) uuid.UUID {
						f, _ := uuid.FromString(teg.Checks[id].ConfigId.String())
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

				if startupLoad {
				startinstanceloop:
					for ldInstId, ldInst := range teg.loadedInstances[i] {
						// check for data from loaded instance
						if ldInst.InstanceSvcCfgHash == inst.InstanceSvcCfgHash {
							inst.InstanceId, _ = uuid.FromString(ldInstId)
							if !uuid.Equal(ldInst.ConfigId, inst.ConfigId) {
								panic(`Matched instances loaded for different ConfigId`)
							}
							inst.InstanceConfigId, _ = uuid.FromString(ldInst.InstanceConfigId.String())
							inst.Version = ldInst.Version
							if inst.ConstraintHash != ldInst.ConstraintHash {
								panic(`Matched instances loaded for different ConstraintHash`)
							}
							if inst.ConstraintValHash != ldInst.ConstraintValHash {
								panic(`Matched instances loaded for different ConstraintValHash`)
							}
							if inst.InstanceService != ldInst.InstanceService {
								panic(`Matched instances loaded for different InstanceService`)
							}
							// we can assume InstanceServiceConfig to
							// be equal, since InstanceSvcCfgHash is
							// equal
							delete(teg.loadedInstances[i], ldInstId)
							log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s, ServiceConstrained=%t",
								repoName,
								`ComputeInstance`,
								`group`,
								teg.Id.String(),
								i,
								ldInstId,
								true,
							)
							break startinstanceloop
						}
						// if we hit here, then just computed an
						// instance that we could not match to any
						// loaded instances -> something is wrong
						panic(`Failed to match computed instance to loaded instances`)
					}
				} else {
					// lookup existing instance ids for check in teg.CheckInstances
					// to determine if this is an update
				instanceloop:
					for _, exInstId := range teg.CheckInstances[i] {
						exInst := teg.Instances[exInstId]
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
					log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s, ServiceConstrained=%t",
						repoName,
						`ComputeInstance`,
						`group`,
						teg.Id.String(),
						i,
						inst.InstanceId.String(),
						true,
					)
				}
				newInstances[inst.InstanceId.String()] = inst
				newCheckInstances = append(newCheckInstances, inst.InstanceId.String())
			}
		} // LOOPEND: range serviceC

		// all instances have been built and matched to
		// loaded instances, but there are loaded
		// instances left. why?
		if startupLoad && len(teg.loadedInstances[i]) != 0 {
			panic(`Leftover matched instances after assignment, computed instances missing`)
		}

		// all new check instances have been built, check which
		// existing instances did not get an update and need to be
		// deleted
		for _, oldInstanceId := range teg.CheckInstances[i] {
			if _, ok := newInstances[oldInstanceId]; !ok {
				// there is no new version for this instance id
				teg.actionCheckInstanceDelete(teg.Instances[oldInstanceId].MakeAction())
				log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s",
					repoName,
					`DeleteInstance`,
					`group`,
					teg.Id.String(),
					i,
					oldInstanceId,
				)
				delete(teg.Instances, oldInstanceId)
				continue
			}
			delete(teg.Instances, oldInstanceId)
			teg.Instances[oldInstanceId] = newInstances[oldInstanceId]
			teg.actionCheckInstanceUpdate(teg.Instances[oldInstanceId].MakeAction())
			log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s",
				repoName,
				`UpdateInstance`,
				`group`,
				teg.Id.String(),
				i,
				oldInstanceId,
			)
		}
		for _, newInstanceId := range newCheckInstances {
			if _, ok := teg.Instances[newInstanceId]; !ok {
				// this instance is new, not an update
				teg.Instances[newInstanceId] = newInstances[newInstanceId]
				// no need to send a create action during load; the
				// action channel is drained anyway
				if !startupLoad {
					teg.actionCheckInstanceCreate(teg.Instances[newInstanceId].MakeAction())
					log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s",
						repoName,
						`CreateInstance`,
						`group`,
						teg.Id.String(),
						i,
						newInstanceId,
					)
				} else {
					log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s",
						repoName,
						`RecreateInstance`,
						`group`,
						teg.Id.String(),
						i,
						newInstanceId,
					)
				}
			}
		}
		delete(teg.CheckInstances, i)
		teg.CheckInstances[i] = newCheckInstances
	} // LOOPEND: range teg.Checks
}

func (teg *SomaTreeElemGroup) evalNativeProp(
	prop string, val string) bool {
	switch prop {
	case "environment":
		env := teg.Parent.(Bucketeer).GetEnvironment()
		if val == env {
			return true
		}
	case "object_type":
		if val == "node" {
			return true
		}
	case "object_state":
		if val == teg.State {
			return true
		}
	case "hardware_node":
		// group != hardware
		return false
	}
	return false
}

func (teg *SomaTreeElemGroup) evalSystemProp(
	prop string, val string, view string) (string, bool) {
	for _, v := range teg.PropertySystem {
		t := v.(*PropertySystem)
		if t.Key == prop && (t.Value == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.Key, true
		}
	}
	return "", false
}

func (teg *SomaTreeElemGroup) evalOncallProp(
	prop string, val string, view string) (string, bool) {
	for _, v := range teg.PropertyOncall {
		t := v.(*PropertyOncall)
		if "OncallId" == prop && t.Id.String() == val && (t.View == view || t.View == `any`) {
			return t.Id.String(), true
		}
	}
	return "", false
}

func (teg *SomaTreeElemGroup) evalCustomProp(
	prop string, val string, view string) (string, bool) {
	for _, v := range teg.PropertyCustom {
		t := v.(*PropertyCustom)
		if t.Key == prop && (t.Value == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.Key, true
		}
	}
	return "", false
}

func (teg *SomaTreeElemGroup) evalServiceProp(
	prop string, val string, view string) (string, bool) {
	for _, v := range teg.PropertyService {
		t := v.(*PropertyService)
		if prop == "name" && (t.Service == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.Id.String(), true
		}
	}
	return "", false
}

func (teg *SomaTreeElemGroup) evalAttributeOfService(
	svcName string, view string, attribute string, value string) bool {
	for _, v := range teg.PropertyService {
		t := v.(*PropertyService)
		if t.Service != svcName {
			continue
		}
		for _, a := range t.Attributes {
			if a.Name == attribute && (t.View == view || t.View == `any`) && (a.Value == value || value == `@defined`) {
				return true
			}
		}
	}
	return false
}

func (teg *SomaTreeElemGroup) evalAttributeProp(
	view string, attr string, value string) (bool, map[string]string) {
	f := map[string]string{}
	for _, v := range teg.PropertyService {
		t := v.(*PropertyService)
		for _, a := range t.Attributes {
			if a.Name == attr && (a.Value == value || value == `@defined`) && (t.View == view || view == `any`) {
				f[t.Id.String()] = a.Name
			}
		}
	}
	if len(f) > 0 {
		return true, f
	}
	return false, f
}

func (teg *SomaTreeElemGroup) getServiceMap(serviceId string) map[string][]string {
	svc := new(PropertyService)
	svc = teg.PropertyService[serviceId].(*PropertyService)

	res := map[string][]string{}
	for _, v := range svc.Attributes {
		res[v.Name] = append(res[v.Name], v.Value)
	}
	return res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
