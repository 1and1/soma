package somatree

import (
	"sync"

	"github.com/satori/go.uuid"
)

//
// Interface: Checker
func (teb *SomaTreeElemBucket) SetCheck(c Check) {
	c.Id = c.GetItemId(teb.Type, teb.Id)
	if uuid.Equal(c.Id, uuid.Nil) {
		c.Id = uuid.NewV4()
	}
	// this check is the source check
	c.InheritedFrom = teb.Id
	c.Inherited = false
	c.SourceId, _ = uuid.FromString(c.Id.String())
	c.SourceType = teb.Type
	// send a scrubbed copy downward
	f := c.clone()
	f.Inherited = true
	f.Id = uuid.Nil
	teb.inheritCheckDeep(f)
	// scrub checkitem startup information prior to storing
	c.Items = nil
	teb.storeCheck(c)
}

func (teb *SomaTreeElemBucket) inheritCheck(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.clone()
	f.Id = f.GetItemId(teb.Type, teb.Id)
	if uuid.Equal(f.Id, uuid.Nil) {
		f.Id = uuid.NewV4()
	}
	f.Items = nil
	teb.storeCheck(f)
	// send original check downwards
	c.Id = uuid.Nil
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
	teb.actionCheckNew(teb.setupCheckAction(c))
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
