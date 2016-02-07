package somatree

import "sync"

//
// Interface: SomaTreeChecker
func (teg *SomaTreeElemGroup) SetCheck(c SomaTreeCheck) {
	c.InheritedFrom = teg.Id
	c.Inherited = false
	teg.storeCheck(c)

	f := SomaTreeCheck{}
	f = c
	f.Inherited = true
	teg.inheritCheckDeep(f)
}

func (teg *SomaTreeElemGroup) inheritCheck(c SomaTreeCheck) {
	teg.storeCheck(c)
	teg.inheritCheckDeep(c)
}

func (teg *SomaTreeElemGroup) inheritCheckDeep(c SomaTreeCheck) {
	var wg sync.WaitGroup
	for child, _ := range teg.Children {
		wg.Add(1)
		ch := child
		go func(stc SomaTreeCheck) {
			defer wg.Done()
			teg.Children[ch].(SomaTreeChecker).inheritCheck(stc)
		}(c)
	}
	wg.Wait()
}

func (teg *SomaTreeElemGroup) storeCheck(c SomaTreeCheck) {
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
		f := SomaTreeCheck{}
		f = teg.Checks[check]
		f.Inherited = true
		teg.Children[childId].(SomaTreeChecker).inheritCheck(f)
	}
}

func (teg *SomaTreeElemGroup) checkCheck(checkId string) bool {
	if _, ok := teg.Checks[checkId]; ok {
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
