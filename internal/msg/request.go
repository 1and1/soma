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
	Bucket      proto.Bucket
	Category    proto.Category
	Environment proto.Environment
	Grant       proto.Grant
	Instance    proto.Instance
	Job         proto.Job
	Permission  proto.Permission
	Repository  proto.Repository
	SectionObj  proto.Section
	State       proto.State
	System      proto.SystemOperation
	Tree        proto.Tree
	Workflow    proto.Workflow
}

type Filter struct {
	IsDetailed bool
	Job        proto.JobFilter
}

type UpdateData struct {
	Environment proto.Environment
	State       proto.State
}

type Flags struct {
	JobDetail bool
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
