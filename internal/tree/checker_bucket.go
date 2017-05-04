/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"sync"

	"github.com/satori/go.uuid"
)

// Implementation of the `Checker` interface

//
// Checker:> Add Check

func (teb *Bucket) SetCheck(c Check) {
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
	f := c.Clone()
	f.Inherited = true
	f.Id = uuid.Nil
	teb.setCheckOnChildren(f)
	// scrub checkitem startup information prior to storing
	c.Items = nil
	teb.addCheck(c)
}

func (teb *Bucket) setCheckInherited(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.Clone()
	f.Id = f.GetItemId(teb.Type, teb.Id)
	if uuid.Equal(f.Id, uuid.Nil) {
		f.Id = uuid.NewV4()
	}
	f.Items = nil
	teb.addCheck(f)
	// send original check downwards
	c.Id = uuid.Nil
	teb.setCheckOnChildren(c)
}

func (teb *Bucket) setCheckOnChildren(c Check) {
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		ch := child
		go func(stc Check) {
			defer wg.Done()
			teb.Children[ch].(Checker).setCheckInherited(stc)
		}(c)
	}
	wg.Wait()
}

func (teb *Bucket) addCheck(c Check) {
	teb.Checks[c.Id.String()] = c
	teb.actionCheckNew(teb.setupCheckAction(c))
}

//
// Checker:> Remove Check

func (teb *Bucket) DeleteCheck(c Check) {
	teb.rmCheck(c)
	teb.deleteCheckOnChildren(c)
}

func (teb *Bucket) deleteCheckInherited(c Check) {
	teb.rmCheck(c)
	teb.deleteCheckOnChildren(c)
}

func (teb *Bucket) deleteCheckOnChildren(c Check) {
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		go func(stc Check, ch string) {
			defer wg.Done()
			teb.Children[ch].(Checker).deleteCheckInherited(stc)
		}(c, child)
	}
	wg.Wait()
}

func (teb *Bucket) rmCheck(c Check) {
	for id, _ := range teb.Checks {
		if uuid.Equal(teb.Checks[id].SourceId, c.SourceId) {
			teb.actionCheckRemoved(teb.setupCheckAction(teb.Checks[id]))
			delete(teb.Checks, id)
			return
		}
	}
}

//
// Checker:> Meta

func (teb *Bucket) syncCheck(childId string) {
	for check, _ := range teb.Checks {
		if !teb.Checks[check].Inheritance {
			continue
		}
		// build a pristine version for inheritance
		f := teb.Checks[check]
		c := f.Clone()
		c.Inherited = true
		c.Id = uuid.Nil
		c.Items = nil
		teb.Children[childId].(Checker).setCheckInherited(c)
	}
}

func (teb *Bucket) checkCheck(checkId string) bool {
	if _, ok := teb.Checks[checkId]; ok {
		return true
	}
	return false
}

// XXX
func (teb *Bucket) LoadInstance(i CheckInstance) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
