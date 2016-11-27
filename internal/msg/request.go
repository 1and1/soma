/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
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
	User       string
	Reply      chan Result
	Search     Filter
	Update     UpdateData
	Flag       Flags

	Super *Supervisor

	ActionObj   proto.Action
	Attribute   proto.Attribute
	Bucket      proto.Bucket
	Category    proto.Category
	Entity      proto.Entity
	Environment proto.Environment
	Grant       proto.Grant
	Instance    proto.Instance
	Job         proto.Job
	Monitoring  proto.Monitoring
	Permission  proto.Permission
	Repository  proto.Repository
	SectionObj  proto.Section
	State       proto.State
	System      proto.SystemOperation
	Team        proto.Team
	Tree        proto.Tree
	UserObj     proto.User
	Workflow    proto.Workflow
}

type Filter struct {
	IsDetailed bool
	Job        proto.JobFilter
}

type UpdateData struct {
	Entity      proto.Entity
	Environment proto.Environment
	State       proto.State
}

type Flags struct {
	JobDetail bool
	Unscoped  bool
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
