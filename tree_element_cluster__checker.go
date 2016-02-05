package somatree

import "sync"

//
// Interface: SomaTreeChecker
func (tec *SomaTreeElemCluster) SetCheck(c SomaTreeCheck) {
	c.InheritedFrom = tec.Id
	c.Inherited = false
	tec.storeCheck(c)

	f := SomaTreeCheck{}
	f = c
	f.Inherited = true
	tec.inheritCheckDeep(f)
}

func (tec *SomaTreeElemCluster) inheritCheck(c SomaTreeCheck) {
	tec.storeCheck(c)
	tec.inheritCheckDeep(c)
}

func (tec *SomaTreeElemCluster) inheritCheckDeep(c SomaTreeCheck) {
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		ch := child
		go func(stc SomaTreeCheck) {
			defer wg.Done()
			tec.Children[ch].(SomaTreeChecker).inheritCheck(stc)
		}(c)
	}
	wg.Wait()
}

func (tec *SomaTreeElemCluster) storeCheck(c SomaTreeCheck) {
	tec.Checks[c.Id.String()] = c
}

func (tec *SomaTreeElemCluster) syncCheck(childId string) {
	for check, _ := range tec.Checks {
		if !tec.Checks[check].Inheritance {
			continue
		}
		f := SomaTreeCheck{}
		f = tec.Checks[check]
		f.Inherited = true
		tec.Children[childId].(SomaTreeChecker).inheritCheck(f)
	}
}

func (tec *SomaTreeElemCluster) checkCheck(checkId string) bool {
	if _, ok := tec.Checks[checkId]; ok {
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
