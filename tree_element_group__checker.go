package somatree

import (
	"sync"

	"github.com/satori/go.uuid"
)

//
// Interface: Checker
func (teg *SomaTreeElemGroup) SetCheck(c Check) {
	c.Id = c.GetItemId(teg.Type, teg.Id)
	if uuid.Equal(c.Id, uuid.Nil) {
		c.Id = uuid.NewV4()
	}
	// this check is the source check
	c.InheritedFrom = teg.Id
	c.Inherited = false
	c.SourceId, _ = uuid.FromString(c.Id.String())
	c.SourceType = teg.Type
	// send a scrubbed copy downward
	f := c.clone()
	f.Inherited = true
	f.Id = uuid.Nil
	teg.inheritCheckDeep(f)
	// scrub checkitem startup information prior to storing
	c.Items = nil
	teg.storeCheck(c)
	teg.actionCheckNew(teg.setupCheckAction(c))
}

func (teg *SomaTreeElemGroup) inheritCheck(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.clone()
	f.Id = f.GetItemId(teg.Type, teg.Id)
	if uuid.Equal(f.Id, uuid.Nil) {
		f.Id = uuid.NewV4()
	}
	f.Items = nil
	teg.storeCheck(f)
	// send original check downwards
	c.Id = uuid.Nil
	teg.inheritCheckDeep(c)
	teg.actionCheckNew(teg.setupCheckAction(f))
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
