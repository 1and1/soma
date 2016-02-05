package somatree

import "fmt"

//
// Interface: SomaTreeChecker
func (ten *SomaTreeElemNode) SetCheck(c SomaTreeCheck) {
	c.InheritedFrom = ten.Id
	c.Inherited = false
	ten.storeCheck(c)
}

func (ten *SomaTreeElemNode) inheritCheck(c SomaTreeCheck) {
	ten.storeCheck(c)
}

// noop, satisfy interface
func (ten *SomaTreeElemNode) inheritCheckDeep(c SomaTreeCheck) {
}

func (ten *SomaTreeElemNode) storeCheck(c SomaTreeCheck) {
	ten.Checks[c.Id.String()] = c
}

// noop, satisfy interface
func (ten *SomaTreeElemNode) syncCheck(childId string) {
}

func (ten *SomaTreeElemNode) checkCheck(checkId string) bool {
	if _, ok := ten.Checks[checkId]; ok {
		return true
	}
	return false
}

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
		hasAttributeConstraint := false
		instance := SomaTreeCheckInstance{
			ConstraintOncall:    []string{},
			ConstraintService:   []string{},
			ConstraintSystem:    []string{},
			ConstraintCustom:    []string{},
			ConstraintNative:    map[string]string{},
			ConstraintAttribute: map[string][]string{},
		}
	constraintcheck:
		for _, constr := range ten.Checks[i].Constraints {
			fmt.Printf("%+v\n", constr)
			fmt.Printf("%+v\n", instance)

			view := ten.Checks[i].View
			switch constr.Type {
			case "native":
				if evalNativeProperty(constr.Key, constr.Value) {
					instance.ConstraintNative[constr.Key] = constr.Value
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "system":
				if id, hit := evalSystemProperty(constr.Key, constr.Value, view); hit {
					instance.ConstraintSystem = append(instance.ConstraintSystem, id)
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "service":
				if id, hit := evalServiceProperty(constr.Key, constr.Value, view); hit {
					instance.ConstraintService = append(instance.ConstraintService, id)
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "oncall":
				if id, hit := evalOncallProperty(constr.Key, constr.Value, view); hit {
					instance.ConstraintOncall = append(instance.ConstraintOncall, id)
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "attribute":
				hasAttributeConstraint = true
			}
		}
		if hasBrokenConstraint {
			continue checksloop
		}
		if len(instance.ConstraintService) > 0 {
		serviceattributecheck:
			for svc := range instance.ConstraintService {
				for _, constr := range ten.Checks[i].Constraints {
					switch constr.Type {
					case "attribute":
						if hit := evalServiceAttribute(constr.Key, constr.Value, view, svc); hit {
							instance.ConstraintAttribute[constr.Value] = append(
								instance.ConstraintAttribute[constr.Value], constr.Value)
						} else {
							hasBrokenConstraint = true
							break serviceattributecheck
						}
					}
				}
			}
		} else {
			// XXX TODO left off here
		}
	}
}

func (ten *SomaTreeElemNode) evalNativeProperty(
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
