/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg

import (
	"fmt"

)

type Result struct {
	Type   string
	Action string
	Code   uint16
	Error  error
	JobId  string

	Super      *Supervisor
	Tree       proto.Tree
	Category   []proto.Category
	Permission []proto.Permission
	Grant      []proto.Grant
	Job        []proto.Job
	System     []proto.SystemOperation
}

func (r *Result) RowCnt(i int64, err error) bool {
	if err != nil {
		r.ServerError(err)
		return false
	}
	switch i {
	case 0:
		r.OK()
		r.SetError(fmt.Errorf(`No rows affected`))
	case 1:
		r.OK()
		return true
	default:
		r.ServerError(fmt.Errorf("Too many rows affected: %d", i))
	}
	return false
}

func (r *Result) Clear(s string) {
	switch s {
	case `category`:
		r.Category = []proto.Category{}
	case `permission`:
		r.Permission = []proto.Permission{}
	case `job`:
		r.Job = []proto.Job{}
	}
}

func (r *Result) SetError(err error) {
	if err != nil {
		r.Error = err
	}
}

func (r *Result) OK() {
	r.Code = 200
	r.Error = nil
}

func (r *Result) Accepted() {
	r.Code = 202
	r.Error = nil
}

func (r *Result) Partial() {
	r.Code = 206
	r.Error = nil
}

func (r *Result) BadRequest(err error) {
	r.Code = 400
	r.SetError(err)
}

func (r *Result) Unauthorized(err error) {
	r.Code = 401
	r.SetError(err)
}

func (r *Result) Forbidden(err error) {
	r.Code = 403
	r.SetError(err)
}

func (r *Result) NotFound(err error) {
	r.Code = 404
	r.SetError(err)
}

func (r *Result) Conflict(err error) {
	r.Code = 406
	r.SetError(err)
}

func (r *Result) ServerError(err error) {
	r.Code = 500
	r.SetError(err)
}

func (r *Result) NotImplemented(err error) {
	r.Code = 500
	r.SetError(err)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
