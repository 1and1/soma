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

func (tec *Cluster) SetCheck(c Check) {
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
	f := c.Clone()
	f.Inherited = true
	f.Id = uuid.Nil
	tec.setCheckOnChildren(f)
	// scrub checkitem startup information prior to storing
	c.Items = nil
	tec.addCheck(c)
}

func (tec *Cluster) setCheckInherited(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.Clone()
	f.Id = f.GetItemId(tec.Type, tec.Id)
	if uuid.Equal(f.Id, uuid.Nil) {
		f.Id = uuid.NewV4()
	}
	f.Items = nil
	tec.addCheck(f)
	// send original check downwards
	c.Id = uuid.Nil
	tec.setCheckOnChildren(c)
}

func (tec *Cluster) setCheckOnChildren(c Check) {
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		ch := child
		go func(stc Check) {
			defer wg.Done()
			tec.Children[ch].(Checker).setCheckInherited(stc)
		}(c)
	}
	wg.Wait()
}

func (tec *Cluster) addCheck(c Check) {
	tec.Checks[c.Id.String()] = c
	tec.actionCheckNew(tec.setupCheckAction(c))
}

//
// Checker:> Remove Check

func (tec *Cluster) DeleteCheck(c Check) {
	tec.rmCheck(c)
	tec.deleteCheckOnChildren(c)
}

func (tec *Cluster) deleteCheckInherited(c Check) {
	tec.rmCheck(c)
	tec.deleteCheckOnChildren(c)
}

func (tec *Cluster) deleteCheckOnChildren(c Check) {
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		go func(stc Check, ch string) {
			defer wg.Done()
			tec.Children[ch].(Checker).deleteCheckInherited(stc)
		}(c, child)
	}
	wg.Wait()
}

func (tec *Cluster) rmCheck(c Check) {
	for id, _ := range tec.Checks {
		if uuid.Equal(tec.Checks[id].SourceId, c.SourceId) {
			tec.actionCheckRemoved(tec.setupCheckAction(tec.Checks[id]))
			delete(tec.Checks, id)
			return
		}
	}
}

//
// Checker:> Meta

func (tec *Cluster) syncCheck(childId string) {
	for check, _ := range tec.Checks {
		if !tec.Checks[check].Inheritance {
			continue
		}
		// build a pristine version for inheritance
		f := tec.Checks[check]
		c := f.Clone()
		c.Inherited = true
		c.Id = uuid.Nil
		c.Items = nil
		tec.Children[childId].(Checker).setCheckInherited(c)
	}
}

func (tec *Cluster) checkCheck(checkId string) bool {
	if _, ok := tec.Checks[checkId]; ok {
		return true
	}
	return false
}

func (tec *Cluster) LoadInstance(i CheckInstance) {
	ckId := i.CheckId.String()
	ckInstId := i.InstanceId.String()
	if tec.loadedInstances[ckId] == nil {
		tec.loadedInstances[ckId] = map[string]CheckInstance{}
	}
	tec.loadedInstances[ckId][ckInstId] = i
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
