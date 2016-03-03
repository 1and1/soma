package somatree

import (
	"sync"

	"github.com/satori/go.uuid"
)

//
// Interface: Checker
func (tec *SomaTreeElemCluster) SetCheck(c Check) {
	c.Id = c.GetItemId(tec.Type, tec.Id)
	if uuid.Equal(c.Id, uuid.Nil) {
		c.Id = uuid.NewV4()
	}
	// this check is the source check
	c.InheritedFrom = tec.Id
	c.Inherited = false
	c.SourceId, _ = uuid.FromString(c.Id.String())
	c.SourceType = tec.Type
	// send a scrubbed copy downward
	f := c.clone()
	f.Inherited = true
	f.Id = uuid.Nil
	tec.inheritCheckDeep(f)
	// scrub checkitem startup information prior to storing
	c.Items = nil
	tec.storeCheck(c)
}

func (tec *SomaTreeElemCluster) inheritCheck(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.clone()
	f.Id = f.GetItemId(tec.Type, tec.Id)
	if uuid.Equal(f.Id, uuid.Nil) {
		f.Id = uuid.NewV4()
	}
	f.Items = nil
	tec.storeCheck(f)
	// send original check downwards
	c.Id = uuid.Nil
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
	tec.actionCheckNew(tec.setupCheckAction(c))
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
