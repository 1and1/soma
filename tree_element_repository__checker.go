package somatree

import (
	"sync"

	"github.com/satori/go.uuid"
)

// Implementation of the `Checker` interface

//
// Checker:> Add Check
func (ter *SomaTreeElemRepository) SetCheck(c Check) {
	c.Id = c.GetItemId(ter.Type, ter.Id)
	if uuid.Equal(c.Id, uuid.Nil) {
		c.Id = uuid.NewV4()
	}
	// this check is the source check
	c.InheritedFrom = ter.Id
	c.Inherited = false
	c.SourceId, _ = uuid.FromString(c.Id.String())
	c.SourceType = ter.Type
	// send a scrubbed copy downward
	f := c.clone()
	f.Inherited = true
	f.Id = uuid.Nil
	ter.inheritCheckDeep(f)
	// scrub checkitem startup information prior to storing
	c.Items = nil
	ter.storeCheck(c)
}

func (ter *SomaTreeElemRepository) inheritCheck(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.clone()
	f.Id = f.GetItemId(ter.Type, ter.Id)
	if uuid.Equal(f.Id, uuid.Nil) {
		f.Id = uuid.NewV4()
	}
	f.Items = nil
	ter.storeCheck(f)
	// send original check downwards
	c.Id = uuid.Nil
	ter.inheritCheckDeep(c)
}

func (ter *SomaTreeElemRepository) inheritCheckDeep(c Check) {
	var wg sync.WaitGroup
	for child, _ := range ter.Children {
		wg.Add(1)
		ch := child
		go func(stc Check) {
			defer wg.Done()
			ter.Children[ch].(Checker).inheritCheck(stc)
		}(c)
	}
	wg.Wait()
}

func (ter *SomaTreeElemRepository) storeCheck(c Check) {
	ter.Checks[c.Id.String()] = c
	ter.actionCheckNew(c.MakeAction())
}

//
// Checker:> Remove Check

func (ter *SomaTreeElemRepository) DeleteCheck(c Check) {
	ter.rmCheck(c)
	ter.deleteCheckOnChildren(c)
}

func (ter *SomaTreeElemRepository) deleteCheckInherited(c Check) {
	ter.rmCheck(c)
	ter.deleteCheckOnChildren(c)
}

func (ter *SomaTreeElemRepository) deleteCheckOnChildren(c Check) {
	var wg sync.WaitGroup
	for child, _ := range ter.Children {
		wg.Add(1)
		go func(stc Check, ch string) {
			defer wg.Done()
			ter.Children[ch].(Checker).deleteCheckInherited(stc)
		}(c, child)
	}
	wg.Wait()
}

func (ter *SomaTreeElemRepository) rmCheck(c Check) {
	for id, _ := range ter.Checks {
		if uuid.Equal(ter.Checks[id].SourceId, c.SourceId) {
			ter.actionCheckRemoved(ter.setupCheckAction(ter.Checks[id]))
			delete(ter.Checks, id)
			return
		}
	}
}

//
// Checker:> Meta

func (ter *SomaTreeElemRepository) syncCheck(childId string) {
	for check, _ := range ter.Checks {
		if !ter.Checks[check].Inheritance {
			continue
		}
		f := Check{}
		f = ter.Checks[check]
		f.Inherited = true
		ter.Children[childId].(Checker).inheritCheck(f)
	}
}

func (ter *SomaTreeElemRepository) checkCheck(checkId string) bool {
	if _, ok := ter.Checks[checkId]; ok {
		return true
	}
	return false
}

// XXX
func (ter *SomaTreeElemRepository) LoadInstance(i CheckInstance) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
