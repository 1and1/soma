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
	check := Check{
		Id:           uuid.NewV4(),
		SourceId:     uuid.NewV4(),
		SourceType:   `sourceType`,
		Inherited:    false,
		CapabilityId: uuid.NewV4(),
		ConfigId:     uuid.NewV4(),
		Inheritance:  false,
		ChildrenOnly: false,
		View:         `any`,
		Interval:     0,
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

	clone := check.Clone()

	if !uuid.Equal(check.Id, clone.Id) {
		t.Errorf(`Illegal clone`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
