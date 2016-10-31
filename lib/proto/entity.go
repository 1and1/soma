/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

// An Entity is a Type without the golang keyword problem
type Entity struct {
	Name    string         `json:"entity,omitempty"`
	Details *EntityDetails `json:"details,omitempty"`
}

type EntityDetails struct {
	DetailsCreation
}

func NewEntityRequest() Request {
	return Request{
		Flags:  &Flags{},
		Entity: &Entity{},
	}
}

func NewEntityResult() Result {
	return Result{
		Errors:   &[]string{},
		Entities: &[]Entity{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
