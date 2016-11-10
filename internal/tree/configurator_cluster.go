/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import uuid "github.com/satori/go.uuid"

func (tec *Cluster) updateCheckInstances() {
	repoName := tec.GetRepositoryName()

	// object may have no checks, but there could be instances to mop up
	if len(tec.Checks) == 0 && len(tec.Instances) == 0 {
		tec.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, HasChecks=%t",
			repoName,
			`UpdateCheckInstances`,
			`cluster`,
			tec.Id.String(),
			false,
		)
		// found nothing to do, ensure update flag is unset again
		tec.hasUpdate = false
		return
	}

	// if there are loaded instances, then this is the initial rebuild
	// of the tree
	startupLoad := false
	if len(tec.loadedInstances) > 0 {
		startupLoad = true
	}

	// if this is not the startupLoad and there are no updates, then there
	// is noting to do
	if !startupLoad && !tec.hasUpdate {
		return
	}

	// scan over all current checkinstances if their check still exists.
	// If not the check has been deleted and the spawned instances need
	// a good deletion
	for ck, _ := range tec.CheckInstances {
		if _, ok := tec.Checks[ck]; ok {
			// check still exists
			continue
		}

		// check no longer exists -> cleanup
		inst := tec.CheckInstances[ck]
		for _, i := range inst {
			tec.actionCheckInstanceDelete(tec.Instances[i].MakeAction())
			tec.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s",
				repoName,
				`CleanupInstance`,
				`cluster`,
				tec.Id.String(),
				ck,
				i,
			)
			delete(tec.Instances, i)
		}
		delete(tec.CheckInstances, ck)
	}

	// loop over all checks and test if there is a reason to disable
	// its check instances. And with disable we mean delete.
	for chk, _ := range tec.Checks {
		disableThis := false
		// disable this check if the system property
		// `disable_all_monitoring` is set for the view that the check
		// uses
		if _, hit, _ := tec.evalSystemProp(
			`disable_all_monitoring`,
			`true`,
			tec.Checks[chk].View,
		); hit {
			disableThis = true
		}
		// disable this check if the system property
		// `disable_check_configuration` is set to the
		// check_configuration that spawned this check
		if _, hit, _ := tec.evalSystemProp(
			`disable_check_configuration`,
			tec.Checks[chk].ConfigId.String(),
			tec.Checks[chk].View,
		); hit {
			disableThis = true
		}
		// if there was a reason to disable this check, all instances
		// are deleted
		if disableThis {
			if instanceArray, ok := tec.CheckInstances[chk]; ok {
				for _, i := range instanceArray {
					tec.actionCheckInstanceDelete(tec.Instances[i].MakeAction())
					tec.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s",
						repoName,
						`RemoveDisabledInstance`,
						`cluster`,
						tec.Id.String(),
						chk,
						i,
					)
					delete(tec.Instances, i)
				}
				delete(tec.CheckInstances, chk)
			}
		}
	}

	// process remaining checks
checksloop:
	for i, _ := range tec.Checks {
		if tec.Checks[i].Inherited == false && tec.Checks[i].ChildrenOnly == true {
			continue checksloop
		}
		if tec.Checks[i].View == "local" {
			continue checksloop
		}
		// skip check if its view has `disable_all_monitoring`
		// property set
		if _, hit, _ := tec.evalSystemProp(
			`disable_all_monitoring`,
			`true`,
			tec.Checks[i].View,
		); hit {
			continue checksloop
		}
		// skip check if there is a matching `disable_check_configuration`
		// property
		if _, hit, _ := tec.evalSystemProp(
			`disable_check_configuration`,
			tec.Checks[i].ConfigId.String(),
			tec.Checks[i].View,
		); hit {
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
				if id, hit, bind := tec.evalSystemProp(c.Key, c.Value, view); hit {
					systemC[id] = bind
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
				if id, hit, bind := tec.evalCustomProp(c.Key, c.Value, view); hit {
					customC[id] = bind
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "service":
				hasServiceConstraint = true
				if id, hit, bind := tec.evalServiceProp(c.Key, c.Value, view); hit {
					serviceC[id] = bind
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
			tec.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, Match=%t",
				repoName,
				`ConstraintEvaluation`,
				`cluster`,
				tec.Id.String(),
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
					hit, bind := tec.evalAttributeOfService(id, view, attr.Key, attr.Value)
					if hit {
						if attributeC[id] == nil {
							// attributeC[id] might still be a nil map
							attributeC[id] = map[string][]string{}
						}
						attributeC[id][attr.Key] = append(attributeC[id][attr.Key], bind)
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
					for id, bind := range svcIdMap {
						serviceC[id] = svcIdMap[id]
						if attributeC[id] == nil {
							// attributeC[id] might still be a nil map
							attributeC[id] = make(map[string][]string)
						}
						attributeC[id][attr.Key] = append(attributeC[id][attr.Key], bind)
					}
				}
			}
			// delete all services that did not match all attributes
			//
			// if a check has two attribute constraints on the same
			// attribute, then len(attributeC[id]) != len(attributes)
			for id, _ := range attributeC {
				if tec.countAttribC(attributeC[id]) != attrCount {
					delete(serviceC, id)
					delete(attributeC, id)
				}
			}
			// declare service constraints in effect if we found a
			// service that bound all attribute constraints
			if len(serviceC) > 0 {
				hasServiceConstraint = true
			} else {
				// found no services that fulfilled all constraints
				hasBrokenConstraint = true
			}
		}
		if hasBrokenConstraint {
			tec.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, Match=%t",
				repoName,
				`ConstraintEvaluation`,
				`cluster`,
				tec.Id.String(),
				i,
				false,
			)
			continue checksloop
		}
		// check triggered, create instances
		tec.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, Match=%t",
			repoName,
			`ConstraintEvaluation`,
			`cluster`,
			tec.Id.String(),
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
					f, _ := uuid.FromString(tec.Checks[id].ConfigId.String())
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
				for ldInstId, ldInst := range tec.loadedInstances[i] {
					if ldInst.InstanceSvcCfgHash != "" {
						continue nosvcstartinstanceloop
					}
					// check if an instance exists bound against the same
					// constraints
					if ldInst.ConstraintHash == inst.ConstraintHash &&
						uuid.Equal(ldInst.ConfigId, inst.ConfigId) &&
						ldInst.ConstraintValHash == inst.ConstraintValHash {

						// found a match
						inst.InstanceId, _ = uuid.FromString(ldInstId)
						inst.InstanceConfigId, _ = uuid.FromString(ldInst.InstanceConfigId.String())
						inst.Version = ldInst.Version
						delete(tec.loadedInstances[i], ldInstId)
						tec.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s, ServiceConstrained=%t",
							repoName,
							`ComputeInstance`,
							`cluster`,
							tec.Id.String(),
							i,
							ldInstId,
							false,
						)
						goto nosvcstartinstancematch
					}
				}
				// if we hit here, then we just computed an instance
				// that we could not match to any loaded instances
				// -> something is wrong
				tec.log.Printf("TK[%s]: Failed to match computed instance to loaded instances."+
					" ObjType=%s, ObjId=%s, CheckId=%s", `cluster`, tec.Id.String(), i, repoName)
				tec.Fault.Error <- &Error{Action: `Failed to match a computed instance to loaded data`}
				return
			nosvcstartinstancematch:
			} else {
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
				tec.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s, ServiceConstrained=%t",
					repoName,
					`ComputeInstance`,
					`cluster`,
					tec.Id.String(),
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
					ConfigId: func(id string) uuid.UUID {
						f, _ := uuid.FromString(tec.Checks[id].ConfigId.String())
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
					for ldInstId, ldInst := range tec.loadedInstances[i] {
						// check for data from loaded instance
						if ldInst.InstanceSvcCfgHash == inst.InstanceSvcCfgHash &&
							ldInst.ConstraintHash == inst.ConstraintHash &&
							ldInst.ConstraintValHash == inst.ConstraintValHash &&
							ldInst.InstanceService == inst.InstanceService &&
							uuid.Equal(ldInst.ConfigId, inst.ConfigId) {

							// found a match
							inst.InstanceId, _ = uuid.FromString(ldInstId)
							inst.InstanceConfigId, _ = uuid.FromString(ldInst.InstanceConfigId.String())
							inst.Version = ldInst.Version
							// we can assume InstanceServiceConfig to
							// be equal, since InstanceSvcCfgHash is
							// equal
							delete(tec.loadedInstances[i], ldInstId)
							tec.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s, ServiceConstrained=%t",
								repoName,
								`ComputeInstance`,
								`cluster`,
								tec.Id.String(),
								i,
								ldInstId,
								true,
							)
							goto startinstancematch
						}
					}
					// if we hit here, then just computed an
					// instance that we could not match to any
					// loaded instances -> something is wrong
					tec.log.Printf("TK[%s]: Failed to match computed instance to loaded instances."+
						" ObjType=%s, ObjId=%s, CheckId=%s", `cluster`, tec.Id.String(), i, repoName)
					tec.Fault.Error <- &Error{Action: `Failed to match a computed instance to loaded data`}
					return
				startinstancematch:
				} else {
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
					tec.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s, ServiceConstrained=%t",
						repoName,
						`ComputeInstance`,
						`cluster`,
						tec.Id.String(),
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
		if startupLoad && len(tec.loadedInstances[i]) != 0 {
			tec.Fault.Error <- &Error{Action: `Leftover matched instances after assignment, computed instances missing`}
			return
		}

		// all new check instances have been built, check which
		// existing instances did not get an update and need to be
		// deleted
		for _, oldInstanceId := range tec.CheckInstances[i] {
			if _, ok := newInstances[oldInstanceId]; !ok {
				// there is no new version for this instance id
				tec.actionCheckInstanceDelete(tec.Instances[oldInstanceId].MakeAction())
				tec.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s",
					repoName,
					`DeleteInstance`,
					`cluster`,
					tec.Id.String(),
					i,
					oldInstanceId,
				)
				delete(tec.Instances, oldInstanceId)
				continue
			}
			delete(tec.Instances, oldInstanceId)
			tec.Instances[oldInstanceId] = newInstances[oldInstanceId]
			tec.actionCheckInstanceUpdate(tec.Instances[oldInstanceId].MakeAction())
			tec.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s",
				repoName,
				`UpdateInstance`,
				`cluster`,
				tec.Id.String(),
				i,
				oldInstanceId,
			)
		}
		for _, newInstanceId := range newCheckInstances {
			if _, ok := tec.Instances[newInstanceId]; !ok {
				// this instance is new, not an update
				tec.Instances[newInstanceId] = newInstances[newInstanceId]
				// no need to send a create action during load; the
				// action channel is drained anyway
				if !startupLoad {
					tec.actionCheckInstanceCreate(tec.Instances[newInstanceId].MakeAction())
					tec.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s",
						repoName,
						`CreateInstance`,
						`cluster`,
						tec.Id.String(),
						i,
						newInstanceId,
					)
				} else {
					tec.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s",
						repoName,
						`RecreateInstance`,
						`cluster`,
						tec.Id.String(),
						i,
						newInstanceId,
					)
				}
			}
		}
		delete(tec.CheckInstances, i)
		tec.CheckInstances[i] = newCheckInstances
	} // LOOPEND: range tec.Checks

	// completed the pass, reset update flag
	tec.hasUpdate = false
}

func (tec *Cluster) evalNativeProp(
	prop string, val string) bool {
	switch prop {
	case "environment":
		env := tec.Parent.(Bucketeer).GetEnvironment()
		if val == env {
			return true
		}
	case "object_type":
		if val == "cluster" {
			return true
		}
	case "object_state":
		if val == tec.State {
			return true
		}
	case "hardware_node":
		// cluster != hardware
		return false
	}
	return false
}

func (tec *Cluster) evalSystemProp(
	prop string, val string, view string) (string, bool, string) {
	for _, v := range tec.PropertySystem {
		t := v.(*PropertySystem)
		if t.Key == prop && (t.Value == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.Key, true, t.Value
		}
	}
	return "", false, ""
}

func (tec *Cluster) evalOncallProp(
	prop string, val string, view string) (string, bool) {
	for _, v := range tec.PropertyOncall {
		t := v.(*PropertyOncall)
		if "OncallId" == prop && t.Id.String() == val && (t.View == view || t.View == `any`) {
			return t.Id.String(), true
		}
	}
	return "", false
}

func (tec *Cluster) evalCustomProp(
	prop string, val string, view string) (string, bool, string) {
	for _, v := range tec.PropertyCustom {
		t := v.(*PropertyCustom)
		if t.Key == prop && (t.Value == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.Key, true, t.Value
		}
	}
	return "", false, ""
}

func (tec *Cluster) evalServiceProp(
	prop string, val string, view string) (string, bool, string) {
	for _, v := range tec.PropertyService {
		t := v.(*PropertyService)
		if prop == "name" && (t.Service == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.Id.String(), true, t.Service
		}
	}
	return "", false, ""
}

func (tec *Cluster) evalAttributeOfService(
	svcId string, view string, attribute string, value string) (bool, string) {
	t := tec.PropertyService[svcId].(*PropertyService)
	for _, a := range t.Attributes {
		if a.Name == attribute && (t.View == view || t.View == `any`) && (a.Value == value || value == `@defined`) {
			return true, a.Value
		}
	}
	return false, ""
}

func (tec *Cluster) evalAttributeProp(
	view string, attr string, value string) (bool, map[string]string) {
	f := map[string]string{}
svcloop:
	for _, v := range tec.PropertyService {
		t := v.(*PropertyService)
		for _, a := range t.Attributes {
			if a.Name == attr && (a.Value == value || value == `@defined`) && (t.View == view || t.View == `any`) {
				f[t.Id.String()] = a.Value
				continue svcloop
			}
		}
	}
	if len(f) > 0 {
		return true, f
	}
	return false, f
}

func (tec *Cluster) getServiceMap(serviceId string) map[string][]string {
	svc := new(PropertyService)
	svc = tec.PropertyService[serviceId].(*PropertyService)

	res := map[string][]string{}
	for _, v := range svc.Attributes {
		res[v.Name] = append(res[v.Name], v.Value)
	}
	return res
}

func (tec *Cluster) countAttribC(attributeC map[string][]string) int {
	var count int = 0
	for key, _ := range attributeC {
		count = count + len(attributeC[key])
	}
	return count
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
