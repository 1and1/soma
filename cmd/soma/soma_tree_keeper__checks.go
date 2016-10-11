/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */
package main

import (
	"github.com/1and1/soma/internal/tree"
	"github.com/1and1/soma/lib/proto"
	"github.com/satori/go.uuid"
)

func (tk *treeKeeper) addCheck(config *proto.CheckConfig) error {
	if chk, err := tk.convertCheck(config); err == nil {
		tk.tree.Find(tree.FindRequest{
			ElementType: config.ObjectType,
			ElementId:   config.ObjectId,
		}, true).SetCheck(*chk)
		return nil
	} else {
		return err
	}
}

func (tk *treeKeeper) rmCheck(config *proto.CheckConfig) error {
	if chk, err := tk.convertCheckForDelete(config); err == nil {
		tk.tree.Find(tree.FindRequest{
			ElementType: config.ObjectType,
			ElementId:   config.ObjectId,
		}, true).DeleteCheck(*chk)
		return nil
	} else {
		return err
	}
}

func (tk *treeKeeper) convertCheck(conf *proto.CheckConfig) (*tree.Check, error) {
	treechk := &tree.Check{
		Id:            uuid.Nil,
		SourceId:      uuid.Nil,
		InheritedFrom: uuid.Nil,
		Inheritance:   conf.Inheritance,
		ChildrenOnly:  conf.ChildrenOnly,
		Interval:      conf.Interval,
	}
	treechk.CapabilityId, _ = uuid.FromString(conf.CapabilityId)
	treechk.ConfigId, _ = uuid.FromString(conf.Id)
	if err := tk.get_view.QueryRow(conf.CapabilityId).Scan(&treechk.View); err != nil {
		return &tree.Check{}, err
	}

	treechk.Thresholds = make([]tree.CheckThreshold, len(conf.Thresholds))
	for i, thr := range conf.Thresholds {
		nthr := tree.CheckThreshold{
			Predicate: thr.Predicate.Symbol,
			Level:     uint8(thr.Level.Numeric),
			Value:     thr.Value,
		}
		treechk.Thresholds[i] = nthr
	}

	treechk.Constraints = make([]tree.CheckConstraint, len(conf.Constraints))
	for i, constr := range conf.Constraints {
		ncon := tree.CheckConstraint{
			Type: constr.ConstraintType,
		}
		switch constr.ConstraintType {
		case "native":
			ncon.Key = constr.Native.Name
			ncon.Value = constr.Native.Value
		case "oncall":
			ncon.Key = "OncallId"
			ncon.Value = constr.Oncall.Id
		case "custom":
			ncon.Key = constr.Custom.Id
			ncon.Value = constr.Custom.Value
		case "system":
			ncon.Key = constr.System.Name
			ncon.Value = constr.System.Value
		case "service":
			ncon.Key = "name"
			ncon.Value = constr.Service.Name
		case "attribute":
			ncon.Key = constr.Attribute.Name
			ncon.Value = constr.Attribute.Value
		}
		treechk.Constraints[i] = ncon
	}
	return treechk, nil
}

func (tk *treeKeeper) convertCheckForDelete(conf *proto.CheckConfig) (*tree.Check, error) {
	var err error
	treechk := &tree.Check{
		Id:            uuid.Nil,
		InheritedFrom: uuid.Nil,
	}
	if treechk.SourceId, err = uuid.FromString(conf.ExternalId); err != nil {
		return nil, err
	}
	if treechk.ConfigId, err = uuid.FromString(conf.Id); err != nil {
		return nil, err
	}
	return treechk, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
