package somatree

import "sync"

//
// Interface: Checker
func (tec *SomaTreeElemCluster) SetCheck(c Check) {
	c.InheritedFrom = tec.Id
	c.Inherited = false
	tec.storeCheck(c)

	f := Check{}
	f = c
	f.Inherited = true
	tec.inheritCheckDeep(f)
}

func (tec *SomaTreeElemCluster) inheritCheck(c Check) {
	tec.storeCheck(c)
	tec.inheritCheckDeep(c)
}

func (tec *SomaTreeElemCluster) inheritCheckDeep(c Check) {
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		ch := child
		go func(stc Check) {
			defer wg.Done()
			tec.Children[ch].(Checker).inheritCheck(stc)
		}(c)
	}
	wg.Wait()
}

func (tec *SomaTreeElemCluster) storeCheck(c Check) {
	tec.Checks[c.Id.String()] = c

	tec.Action <- &Action{
		Action:          "create_check",
		Type:            "cluster",
		Id:              tec.Id.String(),
		CheckId:         c.Id.String(),
		CheckSource:     c.InheritedFrom.String(),
		CheckCapability: c.CapabilityId.String(),
	}
}

func (tec *SomaTreeElemCluster) syncCheck(childId string) {
	for check, _ := range tec.Checks {
		if !tec.Checks[check].Inheritance {
			continue
		}
		f := Check{}
		f = tec.Checks[check]
		f.Inherited = true
		tec.Children[childId].(Checker).inheritCheck(f)
	}
}

func (tec *SomaTreeElemCluster) checkCheck(checkId string) bool {
	if _, ok := tec.Checks[checkId]; ok {
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
