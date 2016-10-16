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

func (tk *treeKeeper) addProperty(q *treeRequest) {
	prop, id := tk.convProperty(`add`, q)
	tk.tree.Find(tree.FindRequest{
		ElementType: q.RequestType,
		ElementId:   id,
	}, true).(tree.Propertier).SetProperty(prop)
}

func (tk *treeKeeper) rmProperty(q *treeRequest) {
	prop, id := tk.convProperty(`rm`, q)
	tk.tree.Find(tree.FindRequest{
		ElementType: q.RequestType,
		ElementId:   id,
	}, true).(tree.Propertier).DeleteProperty(prop)
}

func (tk *treeKeeper) convProperty(task string, q *treeRequest) (
	tree.Property, string) {

	var prop tree.Property
	var id string

	switch q.RequestType {
	case `node`:
		id = q.Node.Node.Id
		prop = tk.pTT(task, (*q.Node.Node.Properties)[0])
	case `cluster`:
		id = q.Cluster.Cluster.Id
		prop = tk.pTT(task, (*q.Cluster.Cluster.Properties)[0])
	case `group`:
		id = q.Group.Group.Id
		prop = tk.pTT(task, (*q.Group.Group.Properties)[0])
	case `bucket`:
		id = q.Bucket.Bucket.Id
		prop = tk.pTT(task, (*q.Bucket.Bucket.Properties)[0])
	case `repository`:
		id = q.Repository.Repository.Id
		prop = tk.pTT(task, (*q.Repository.Repository.Properties)[0])
	}
	return prop, id
}

func (tk *treeKeeper) pTT(task string, pp proto.Property) tree.Property {
	switch pp.Type {
	case `custom`:
		customId, _ := uuid.FromString(pp.Custom.Id)
		switch task {
		case `add`:
			return &tree.PropertyCustom{
				Id:           uuid.NewV4(),
				CustomId:     customId,
				Inheritance:  pp.Inheritance,
				ChildrenOnly: pp.ChildrenOnly,
				View:         pp.View,
				Key:          pp.Custom.Name,
				Value:        pp.Custom.Value,
			}
		case `rm`:
			srcUUID, _ := uuid.FromString(pp.SourceInstanceId)
			return &tree.PropertyCustom{
				SourceId: srcUUID,
				CustomId: customId,
				View:     pp.View,
				Key:      pp.Custom.Name,
				Value:    pp.Custom.Value,
			}
		}
	case `oncall`:
		oncallId, _ := uuid.FromString(pp.Oncall.Id)
		switch task {
		case `add`:
			return &tree.PropertyOncall{
				Id:           uuid.NewV4(),
				OncallId:     oncallId,
				Inheritance:  pp.Inheritance,
				ChildrenOnly: pp.ChildrenOnly,
				View:         pp.View,
				Name:         pp.Oncall.Name,
				Number:       pp.Oncall.Number,
			}
		case `rm`:
			srcUUID, _ := uuid.FromString(pp.SourceInstanceId)
			return &tree.PropertyOncall{
				SourceId: srcUUID,
				OncallId: oncallId,
				View:     pp.View,
				Name:     pp.Oncall.Name,
				Number:   pp.Oncall.Number,
			}
		}
	case `service`:
		switch task {
		case `add`:
			return &tree.PropertyService{
				Id:           uuid.NewV4(),
				Inheritance:  pp.Inheritance,
				ChildrenOnly: pp.ChildrenOnly,
				View:         pp.View,
				Service:      pp.Service.Name,
				Attributes:   pp.Service.Attributes,
			}
		case `rm`:
			srcUUID, _ := uuid.FromString(pp.SourceInstanceId)
			return &tree.PropertyService{
				SourceId: srcUUID,
				View:     pp.View,
				Service:  pp.Service.Name,
			}
		}
	case `system`:
		switch task {
		case `add`:
			return &tree.PropertySystem{
				Id:           uuid.NewV4(),
				Inheritance:  pp.Inheritance,
				ChildrenOnly: pp.ChildrenOnly,
				View:         pp.View,
				Key:          pp.System.Name,
				Value:        pp.System.Value,
			}
		case `rm`:
			srcUUID, _ := uuid.FromString(pp.SourceInstanceId)
			return &tree.PropertySystem{
				SourceId: srcUUID,
				View:     pp.View,
				Key:      pp.System.Name,
				Value:    pp.System.Value,
			}
		}
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
