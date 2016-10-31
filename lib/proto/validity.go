/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Validity struct {
	SystemProperty string           `json:"systemProperty,omitempty"`
	ObjectType     string           `json:"objectType,omitempty"`
	Direct         bool             `json:"direct,string,omitempty"`
	Inherited      bool             `json:"inherited,string,omitempty"`
	Details        *ValidityDetails `json:"details,omitempty"`
}

type ValidityDetails struct {
	DetailsCreation
}

func NewValidityRequest() Request {
	return Request{
		Flags:    &Flags{},
		Validity: &Validity{},
	}
}

func NewValidityResult() Result {
	return Result{
		Errors:     &[]string{},
		Validities: &[]Validity{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
