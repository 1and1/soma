package somatree

import "sync"

//
// Interface: SomaTreeChecker
func (teb *SomaTreeElemBucket) SetCheck(c SomaTreeCheck) {
	c.InheritedFrom = teb.Id
	c.Inherited = false
	teb.storeCheck(c)

	f := SomaTreeCheck{}
	f = c
	f.Inherited = true
	teb.inheritCheckDeep(f)
}

func (teb *SomaTreeElemBucket) inheritCheck(c SomaTreeCheck) {
	teb.storeCheck(c)
	teb.inheritCheckDeep(c)
}

func (teb *SomaTreeElemBucket) inheritCheckDeep(c SomaTreeCheck) {
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		ch := child
		go func(stc SomaTreeCheck) {
			defer wg.Done()
			teb.Children[ch].(SomaTreeChecker).inheritCheck(stc)
		}(c)
	}
	wg.Wait()
}

func (teb *SomaTreeElemBucket) storeCheck(c SomaTreeCheck) {
	teb.Checks[c.Id.String()] = c
}

func (teb *SomaTreeElemBucket) syncCheck(childId string) {
	for check, _ := range teb.Checks {
		if !teb.Checks[check].Inheritance {
			continue
		}
		f := SomaTreeCheck{}
		f = teb.Checks[check]
		f.Inherited = true
		teb.Children[childId].(SomaTreeChecker).inheritCheck(f)
	}
}

func (teb *SomaTreeElemBucket) checkCheck(checkId string) bool {
	if _, ok := teb.Checks[checkId]; ok {
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
