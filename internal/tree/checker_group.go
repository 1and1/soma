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

func (teg *Group) SetCheck(c Check) {
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
	teg.setCheckOnChildren(f)
	// scrub checkitem startup information prior to storing
	c.Items = nil
	teg.addCheck(c)
}

func (teg *Group) setCheckInherited(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.clone()
	f.Id = f.GetItemId(teg.Type, teg.Id)
	if uuid.Equal(f.Id, uuid.Nil) {
		f.Id = uuid.NewV4()
	}
	// send original check downwards
	c.Id = uuid.Nil
	teg.setCheckOnChildren(c)
	f.Items = nil
	teg.addCheck(f)
}

func (teg *Group) setCheckOnChildren(c Check) {
	/*	var wg sync.WaitGroup
		for child, _ := range teg.Children {
			wg.Add(1)
			ch := child
			go func(stc Check) {
				defer wg.Done()
				teg.Children[ch].(Checker).setCheckInherited(stc)
			}(c)
		}
		wg.Wait() */

	// groups
	for i := 0; i < teg.ordNumChildGrp; i++ {
		if child, ok := teg.ordChildrenGrp[i]; ok {
			teg.Children[child].(Checker).setCheckInherited(c)
		}
	}
	// clusters
	for i := 0; i < teg.ordNumChildClr; i++ {
		if child, ok := teg.ordChildrenClr[i]; ok {
			teg.Children[child].(Checker).setCheckInherited(c)
		}
	}
	// nodes
	for i := 0; i < teg.ordNumChildNod; i++ {
		if child, ok := teg.ordChildrenNod[i]; ok {
			teg.Children[child].(Checker).setCheckInherited(c)
		}
	}
}

func (teg *Group) addCheck(c Check) {
	teg.Checks[c.Id.String()] = c
	teg.actionCheckNew(c.MakeAction())
}

//
// Checker:> Remove Check

func (teg *Group) DeleteCheck(c Check) {
	teg.deleteCheckOnChildren(c)
	teg.rmCheck(c)
}

func (teg *Group) deleteCheckInherited(c Check) {
	teg.deleteCheckOnChildren(c)
	teg.rmCheck(c)
}

func (teg *Group) deleteCheckOnChildren(c Check) {
	/*	var wg sync.WaitGroup
		for child, _ := range teg.Children {
			wg.Add(1)
			go func(stc Check, ch string) {
				defer wg.Done()
				teg.Children[ch].(Checker).deleteCheckInherited(stc)
			}(c, child)
		}
		wg.Wait() */

	// groups
	for i := 0; i < teg.ordNumChildGrp; i++ {
		if child, ok := teg.ordChildrenGrp[i]; ok {
			teg.Children[child].(Checker).deleteCheckInherited(c)
		}
	}
	// clusters
	for i := 0; i < teg.ordNumChildClr; i++ {
		if child, ok := teg.ordChildrenClr[i]; ok {
			teg.Children[child].(Checker).deleteCheckInherited(c)
		}
	}
	// nodes
	for i := 0; i < teg.ordNumChildNod; i++ {
		if child, ok := teg.ordChildrenNod[i]; ok {
			teg.Children[child].(Checker).deleteCheckInherited(c)
		}
	}
}

func (teg *Group) rmCheck(c Check) {
	for id, _ := range teg.Checks {
		if uuid.Equal(teg.Checks[id].SourceId, c.SourceId) {
			teg.actionCheckRemoved(teg.setupCheckAction(teg.Checks[id]))
			delete(teg.Checks, id)
			return
		}
	}
}

//
// Checker:> Meta

func (teg *Group) syncCheck(childId string) {
	for check, _ := range teg.Checks {
		if !teg.Checks[check].Inheritance {
			continue
		}
		f := Check{}
		f = teg.Checks[check]
		f.Inherited = true
		teg.Children[childId].(Checker).setCheckInherited(f)
	}
}

func (teg *Group) checkCheck(checkId string) bool {
	if _, ok := teg.Checks[checkId]; ok {
		return true
	}
	return false
}

//
func (teg *Group) LoadInstance(i CheckInstance) {
	ckId := i.CheckId.String()
	ckInstId := i.InstanceId.String()
	if teg.loadedInstances[ckId] == nil {
		teg.loadedInstances[ckId] = map[string]CheckInstance{}
	}
	teg.loadedInstances[ckId][ckInstId] = i
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
