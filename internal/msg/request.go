/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg

import (
	"github.com/1and1/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

type Request struct {
	Section    string
	Action     string
	RemoteAddr string
	AuthUser   string
	Reply      chan Result
	JobID      uuid.UUID
	Search     Filter
	Update     UpdateData
	Flag       Flags

	Super *Supervisor

	ActionObj   proto.Action
	Attribute   proto.Attribute
	Bucket      proto.Bucket
	Capability  proto.Capability
	Category    proto.Category
	CheckConfig proto.CheckConfig
	Cluster     proto.Cluster
	Datacenter  proto.Datacenter
	Entity      proto.Entity
	Environment proto.Environment
	Grant       proto.Grant
	Group       proto.Group
	Instance    proto.Instance
	Job         proto.Job
	Level       proto.Level
	Metric      proto.Metric
	Mode        proto.Mode
	Monitoring  proto.Monitoring
	Node        proto.Node
	Oncall      proto.Oncall
	Permission  proto.Permission
	Predicate   proto.Predicate
	Property    proto.Property
	Provider    proto.Provider
	Repository  proto.Repository
	SectionObj  proto.Section
	Server      proto.Server
	State       proto.State
	Status      proto.Status
	System      proto.SystemOperation
	Team        proto.Team
	Tree        proto.Tree
	Unit        proto.Unit
	User        proto.User
	Validity    proto.Validity
	View        proto.View
	Workflow    proto.Workflow

	CacheRequest *Request
}

type Filter struct {
	IsDetailed bool
	Job        proto.JobFilter
	Server     proto.Server
}

type UpdateData struct {
	Datacenter  proto.Datacenter
	Entity      proto.Entity
	Environment proto.Environment
	State       proto.State
	View        proto.View
}

type Flags struct {
	JobDetail bool
	Unscoped  bool
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
