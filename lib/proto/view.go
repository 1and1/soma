/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type View struct {
	Name    string       `json:"name,omitempty"`
	Details *ViewDetails `json:"details,omitempty"`
}

type ViewDetails struct {
	DetailsCreation
}

func NewViewRequest() Request {
	return Request{
		Flags: &Flags{},
		View:  &View{},
	}
}

func NewViewResult() Result {
	return Result{
		Errors: &[]string{},
		Views:  &[]View{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
