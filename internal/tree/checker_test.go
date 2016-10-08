/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"testing"

	"github.com/satori/go.uuid"
)

func TestCheckClone(t *testing.T) {
	check := testSpawnCheck(false, false, false)

	clone := check.Clone()

	if !uuid.Equal(check.Id, clone.Id) {
		t.Errorf(`Illegal clone`)
	}
}

func TestCheckGetter(t *testing.T) {
	check := testSpawnCheck(false, false, true)

	if _, err := uuid.FromString(check.GetSourceCheckId()); err != nil {
		t.Errorf(`Received error`, err)
	}

	if _, err := uuid.FromString(check.GetCheckConfigId()); err != nil {
		t.Errorf(`Received error`, err)
	}

	if sourceType := check.GetSourceType(); sourceType == "" {
		t.Errorf(`Received empty Check.SourceType`)
	}

	if _, err := uuid.FromString(check.GetCapabilityId()); err != nil {
		t.Errorf(`Received error`, err)
	}

	if view := check.GetView(); view == "" {
		t.Errorf(`Received empty Check.View`)
	} else {
		switch view {
		case `internal`, `external`, `local`, `any`:
		default:
			t.Errorf(`Received unknown View`)
		}
	}

	if interval := check.GetInterval(); interval == 0 {
		t.Errorf(`Execution interval is every zero seconds`)
	}

	if child := check.GetChildrenOnly(); child == false {
		t.Errorf(`GetChildren received zero value return`)
	}
}

func TestCheckInherited(t *testing.T) {
	check := testSpawnCheck(true, true, false)

	if inherit := check.GetIsInherited(); inherit != true {
		t.Errorf(`Incorrect inheritance`)
	}

	if inheritance := check.GetInheritance(); inheritance == false {
		t.Errorf(`Inherited check can not have inheritance disabled`)
	}

	var id, idFrom uuid.UUID
	var err error

	if id, err = uuid.FromString(check.GetCheckId()); err != nil {
		t.Errorf(`Received error`, err)
	}
	if idFrom, err = uuid.FromString(check.GetInheritedFrom()); err != nil {
		t.Errorf(`Received error`, err)
	}
	if uuid.Equal(id, idFrom) {
		t.Errorf(`Equal id/sourceId for inherited check`)
	}
}

func TestCheckNotInherited(t *testing.T) {
	check := testSpawnCheck(false, false, false)

	if inherit := check.GetIsInherited(); inherit != false {
		t.Errorf(`Incorrect inheritance`)
	}

	var id, idFrom uuid.UUID
	var err error

	if id, err = uuid.FromString(check.GetCheckId()); err != nil {
		t.Errorf(`Received error`, err)
	}
	if idFrom, err = uuid.FromString(check.GetInheritedFrom()); err != nil {
		t.Errorf(`Received error`, err)
	}
	if !uuid.Equal(id, idFrom) {
		t.Errorf(`Unequal id/sourceId for non-inherited check`)
	}
}

func TestCheckAction(t *testing.T) {
	check := testSpawnCheck(false, false, false)

	action := check.MakeAction()

	if action.Check.CheckId != check.GetCheckId() {
		t.Errorf(`Created action is incorrect`)
	}
}

func testSpawnCheck(inherited, inheritance, childrenOnly bool) Check {
	id := uuid.NewV4()
	var idFrom uuid.UUID
	if inherited {
		idFrom = uuid.NewV4()
	} else {
		idFrom, _ = uuid.FromString(id.String())
	}

	return Check{
		Id:            id,
		SourceId:      uuid.NewV4(),
		SourceType:    `sourceType`,
		Inherited:     inherited,
		InheritedFrom: idFrom,
		CapabilityId:  uuid.NewV4(),
		ConfigId:      uuid.NewV4(),
		Inheritance:   inheritance,
		ChildrenOnly:  childrenOnly,
		View:          `any`,
		Interval:      1,
		Thresholds: []CheckThreshold{
			{
				Predicate: `==`,
				Level:     0,
				Value:     0,
			},
		},
		Constraints: []CheckConstraint{
			{
				Type:  `type`,
				Key:   `key`,
				Value: `value`,
			},
		},
		Items: []CheckItem{
			{
				ObjectId:   uuid.NewV4(),
				ObjectType: `objectType`,
				ItemId:     uuid.NewV4(),
			},
		},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
