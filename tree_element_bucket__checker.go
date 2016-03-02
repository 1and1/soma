package somatree

import (
	"sync"

)

//
// Interface: Checker
func (teb *SomaTreeElemBucket) SetCheck(c Check) {
	c.InheritedFrom = teb.Id
	c.Inherited = false
	teb.storeCheck(c)

	f := Check{}
	f = c
	f.Inherited = true
	teb.inheritCheckDeep(f)
}

func (teb *SomaTreeElemBucket) inheritCheck(c Check) {
	teb.storeCheck(c)
	teb.inheritCheckDeep(c)
}

func (teb *SomaTreeElemBucket) inheritCheckDeep(c Check) {
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		ch := child
		go func(stc Check) {
			defer wg.Done()
			teb.Children[ch].(Checker).inheritCheck(stc)
		}(c)
	}
	wg.Wait()
}

func (teb *SomaTreeElemBucket) storeCheck(c Check) {
	teb.Checks[c.Id.String()] = c

	teb.Action <- &Action{
		Action: "create_check",
		Type:   teb.Type,
		Bucket: somaproto.ProtoBucket{
			Id:          teb.Id.String(),
			Name:        teb.Name,
			Repository:  teb.Repository.String(),
			Team:        teb.Team.String(),
			Environment: teb.Environment,
			IsDeleted:   teb.Deleted,
			IsFrozen:    teb.Frozen,
		},
		CheckId:         c.Id.String(),
		CheckSource:     c.InheritedFrom.String(),
		CheckCapability: c.CapabilityId.String(),
	}
}

func (teb *SomaTreeElemBucket) syncCheck(childId string) {
	for check, _ := range teb.Checks {
		if !teb.Checks[check].Inheritance {
			continue
		}
		f := Check{}
		f = teb.Checks[check]
		f.Inherited = true
		teb.Children[childId].(Checker).inheritCheck(f)
	}
}

func (teb *SomaTreeElemBucket) checkCheck(checkId string) bool {
	if _, ok := teb.Checks[checkId]; ok {
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
