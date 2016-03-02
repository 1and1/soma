package somatree

import "sync"

//
// Interface: Checker
func (teg *SomaTreeElemGroup) SetCheck(c Check) {
	c.InheritedFrom = teg.Id
	c.Inherited = false
	teg.storeCheck(c)

	f := Check{}
	f = c
	f.Inherited = true
	teg.inheritCheckDeep(f)
}

func (teg *SomaTreeElemGroup) inheritCheck(c Check) {
	teg.storeCheck(c)
	teg.inheritCheckDeep(c)
}

func (teg *SomaTreeElemGroup) inheritCheckDeep(c Check) {
	var wg sync.WaitGroup
	for child, _ := range teg.Children {
		wg.Add(1)
		ch := child
		go func(stc Check) {
			defer wg.Done()
			teg.Children[ch].(Checker).inheritCheck(stc)
		}(c)
	}
	wg.Wait()
}

func (teg *SomaTreeElemGroup) storeCheck(c Check) {
	teg.Checks[c.Id.String()] = c

	teg.Action <- &Action{
		Action:          "create_check",
		Type:            "group",
		Id:              teg.Id.String(),
		CheckId:         c.Id.String(),
		CheckSource:     c.InheritedFrom.String(),
		CheckCapability: c.CapabilityId.String(),
	}
}

func (teg *SomaTreeElemGroup) syncCheck(childId string) {
	for check, _ := range teg.Checks {
		if !teg.Checks[check].Inheritance {
			continue
		}
		f := Check{}
		f = teg.Checks[check]
		f.Inherited = true
		teg.Children[childId].(Checker).inheritCheck(f)
	}
}

func (teg *SomaTreeElemGroup) checkCheck(checkId string) bool {
	if _, ok := teg.Checks[checkId]; ok {
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
