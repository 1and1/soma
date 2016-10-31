/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type State struct {
	Name    string        `json:"Name,omitempty"`
	Details *StateDetails `json:"details,omitempty"`
}

type StateDetails struct {
	DetailsCreation
}

func NewStateRequest() Request {
	return Request{
		Flags: &Flags{},
		State: &State{},
	}
}

func NewStateResult() Result {
	return Result{
		Errors: &[]string{},
		States: &[]State{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
