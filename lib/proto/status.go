/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Status struct {
	Name    string         `json:"name,omitempty"`
	Details *StatusDetails `json:"details,omitempty"`
}

type StatusDetails struct {
	DetailsCreation
}

func NewStatusRequest() Request {
	return Request{
		Flags:  &Flags{},
		Status: &Status{},
	}
}

func NewStatusResult() Result {
	return Result{
		Errors: &[]string{},
		Status: &[]Status{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
