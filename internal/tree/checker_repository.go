/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	//	"sync"

	"github.com/satori/go.uuid"
)

// Implementation of the `Checker` interface

//
// Checker:> Add Check
func (ter *Repository) SetCheck(c Check) {
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
	f := c.Clone()
	f.Inherited = true
	f.Id = uuid.Nil
	ter.setCheckOnChildren(f)
	// scrub checkitem startup information prior to storing
	c.Items = nil
	ter.addCheck(c)
}

func (ter *Repository) setCheckInherited(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.Clone()
	f.Id = f.GetItemId(ter.Type, ter.Id)
	if uuid.Equal(f.Id, uuid.Nil) {
		f.Id = uuid.NewV4()
	}
	// send original check downwards
	c.Id = uuid.Nil
	ter.setCheckOnChildren(c)
	f.Items = nil
	ter.addCheck(f)
}

func (ter *Repository) setCheckOnChildren(c Check) {
	/*	var wg sync.WaitGroup
		for child, _ := range ter.Children {
			wg.Add(1)
			ch := child
			go func(stc Check) {
				defer wg.Done()
				ter.Children[ch].(Checker).setCheckInherited(stc)
			}(c)
		}
		wg.Wait() */

	// buckets
	for i := 0; i < ter.ordNumChildBck; i++ {
		if child, ok := ter.ordChildrenBck[i]; ok {
			ter.Children[child].(Checker).setCheckInherited(c)
		}
	}
}

func (ter *Repository) addCheck(c Check) {
	ter.Checks[c.Id.String()] = c
	ter.actionCheckNew(c.MakeAction())
}

//
// Checker:> Remove Check

func (ter *Repository) DeleteCheck(c Check) {
	ter.deleteCheckOnChildren(c)
	ter.rmCheck(c)
}

func (ter *Repository) deleteCheckInherited(c Check) {
	ter.deleteCheckOnChildren(c)
	ter.rmCheck(c)
}

func (ter *Repository) deleteCheckOnChildren(c Check) {
	/*	var wg sync.WaitGroup
		for child, _ := range ter.Children {
			wg.Add(1)
			go func(stc Check, ch string) {
				defer wg.Done()
				ter.Children[ch].(Checker).deleteCheckInherited(stc)
			}(c, child)
		}
		wg.Wait() */

	// buckets
	for i := 0; i < ter.ordNumChildBck; i++ {
		if child, ok := ter.ordChildrenBck[i]; ok {
			ter.Children[child].(Checker).deleteCheckInherited(c)
		}
	}
}

func (ter *Repository) rmCheck(c Check) {
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

func (ter *Repository) syncCheck(childId string) {
	for check, _ := range ter.Checks {
		if !ter.Checks[check].Inheritance {
			continue
		}
		// build a pristine version for inheritance
		f := ter.Checks[check]
		c := f.Clone()
		c.Inherited = true
		c.Id = uuid.Nil
		c.Items = nil
		ter.Children[childId].(Checker).setCheckInherited(c)
	}
}

func (ter *Repository) checkCheck(checkId string) bool {
	if _, ok := ter.Checks[checkId]; ok {
		return true
	}
	return false
}

// XXX
func (ter *Repository) LoadInstance(i CheckInstance) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
