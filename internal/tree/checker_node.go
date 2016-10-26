/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import "github.com/satori/go.uuid"

// Implementation of the `Checker` interface

//
// Checker:> Add Check

func (ten *Node) SetCheck(c Check) {
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
	ten.addCheck(c)
}

func (ten *Node) setCheckInherited(c Check) {
	// we keep a local copy, that way we know it is ours....
	f := c.Clone()
	f.Id = f.GetItemId(ten.Type, ten.Id)
	if uuid.Equal(f.Id, uuid.Nil) {
		f.Id = uuid.NewV4()
	}
	f.Items = nil
	ten.addCheck(f)
}

func (ten *Node) setCheckOnChildren(c Check) {
}

func (ten *Node) addCheck(c Check) {
	ten.hasUpdate = true
	ten.Checks[c.Id.String()] = c
	ten.actionCheckNew(ten.setupCheckAction(c))
}

//
// Checker:> Remove Check

func (ten *Node) DeleteCheck(c Check) {
	ten.rmCheck(c)
}

func (ten *Node) deleteCheckInherited(c Check) {
	ten.rmCheck(c)
}

func (ten *Node) deleteCheckOnChildren(c Check) {
}

func (ten *Node) rmCheck(c Check) {
	for id, _ := range ten.Checks {
		if uuid.Equal(ten.Checks[id].SourceId, c.SourceId) {
			ten.hasUpdate = true
			ten.actionCheckRemoved(ten.setupCheckAction(ten.Checks[id]))
			delete(ten.Checks, id)
			return
		}
	}
}

// noop, satisfy interface
func (ten *Node) syncCheck(childId string) {
}

func (ten *Node) checkCheck(checkId string) bool {
	if _, ok := ten.Checks[checkId]; ok {
		return true
	}
	return false
}

//
func (ten *Node) LoadInstance(i CheckInstance) {
	ckId := i.CheckId.String()
	ckInstId := i.InstanceId.String()
	if ten.loadedInstances[ckId] == nil {
		ten.loadedInstances[ckId] = map[string]CheckInstance{}
	}
	ten.loadedInstances[ckId][ckInstId] = i
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
