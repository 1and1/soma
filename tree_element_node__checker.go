package somatree

import "github.com/satori/go.uuid"

//
// Interface: Checker
func (ten *SomaTreeElemNode) SetCheck(c Check) {
	c.Id = c.GetItemId(ten.Type, ten.Id)
	if uuid.Equal(c.Id, uuid.Nil) {
		c.Id = uuid.NewV4()
	}
	// this check is the source check
	c.InheritedFrom = ten.Id
	c.Inherited = false
	c.SourceId, _ = uuid.FromString(c.Id.String())
	c.SourceType = ten.Type
	// scrub checkitem startup information prior to storing
	c.Items = nil
	ten.storeCheck(c)
}

func (ten *SomaTreeElemNode) inheritCheck(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.clone()
	f.Id = f.GetItemId(ten.Type, ten.Id)
	if uuid.Equal(f.Id, uuid.Nil) {
		f.Id = uuid.NewV4()
	}
	f.Items = nil
	ten.storeCheck(f)
}

// noop, satisfy interface
func (ten *SomaTreeElemNode) inheritCheckDeep(c Check) {
}

func (ten *SomaTreeElemNode) storeCheck(c Check) {
	ten.Checks[c.Id.String()] = c
	ten.actionCheckNew(ten.setupCheckAction(c))
}

// noop, satisfy interface
func (ten *SomaTreeElemNode) syncCheck(childId string) {
}

func (ten *SomaTreeElemNode) checkCheck(checkId string) bool {
	if _, ok := ten.Checks[checkId]; ok {
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
