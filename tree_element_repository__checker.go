package somatree

import "sync"

//
// Interface: SomaTreeChecker
func (ter *SomaTreeElemRepository) SetCheck(c SomaTreeCheck) {
	c.InheritedFrom = ter.Id
	c.Inherited = false
	ter.storeCheck(c)

	f := SomaTreeCheck{}
	f = c
	f.Inherited = true
	ter.inheritCheckDeep(f)
}

func (ter *SomaTreeElemRepository) inheritCheck(c SomaTreeCheck) {
	ter.storeCheck(c)
	ter.inheritCheckDeep(c)
}

func (ter *SomaTreeElemRepository) inheritCheckDeep(c SomaTreeCheck) {
	var wg sync.WaitGroup
	for child, _ := range ter.Children {
		wg.Add(1)
		ch := child
		go func(stc SomaTreeCheck) {
			defer wg.Done()
			ter.Children[ch].(SomaTreeChecker).inheritCheck(stc)
		}(c)
	}
	wg.Wait()
}

func (ter *SomaTreeElemRepository) storeCheck(c SomaTreeCheck) {
	ter.Checks[c.Id.String()] = c
}

func (ter *SomaTreeElemRepository) syncCheck(childId string) {
	for check, _ := range ter.Checks {
		if !ter.Checks[check].Inheritance {
			continue
		}
		f := SomaTreeCheck{}
		f = ter.Checks[check]
		f.Inherited = true
		ter.Children[childId].(SomaTreeChecker).inheritCheck(f)
	}
}

func (ter *SomaTreeElemRepository) checkCheck(checkId string) bool {
	if _, ok := ter.Checks[checkId]; ok {
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
