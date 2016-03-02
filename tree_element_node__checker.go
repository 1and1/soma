package somatree

//
// Interface: Checker
func (ten *SomaTreeElemNode) SetCheck(c Check) {
	c.InheritedFrom = ten.Id
	c.Inherited = false
	ten.storeCheck(c)
}

func (ten *SomaTreeElemNode) inheritCheck(c Check) {
	ten.storeCheck(c)
}

// noop, satisfy interface
func (ten *SomaTreeElemNode) inheritCheckDeep(c Check) {
}

func (ten *SomaTreeElemNode) storeCheck(c Check) {
	ten.Checks[c.Id.String()] = c

	ten.Action <- &Action{
		Action:          "create_check",
		Type:            "node",
		Id:              ten.Id.String(),
		CheckId:         c.Id.String(),
		CheckSource:     c.InheritedFrom.String(),
		CheckCapability: c.CapabilityId.String(),
	}
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
