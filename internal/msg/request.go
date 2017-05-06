/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg

import "github.com/1and1/soma/lib/proto"

type Request struct {
	Section    string
	Action     string
	RemoteAddr string
	AuthUser   string
	Reply      chan Result
	Search     Filter
	Update     UpdateData
	Flag       Flags

	Super *Supervisor

	ActionObj   proto.Action
	Attribute   proto.Attribute
	Bucket      proto.Bucket
	Category    proto.Category
	Cluster     proto.Cluster
	Entity      proto.Entity
	Environment proto.Environment
	Grant       proto.Grant
	Group       proto.Group
	Instance    proto.Instance
	Job         proto.Job
	Monitoring  proto.Monitoring
	Node        proto.Node
	Permission  proto.Permission
	Property    proto.Property
	Repository  proto.Repository
	SectionObj  proto.Section
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
}

type UpdateData struct {
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
