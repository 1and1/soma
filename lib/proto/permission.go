/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Permission struct {
	Id       string           `json:"id,omitempty"`
	Name     string           `json:"name,omitempty"`
	Category string           `json:"category,omitempty"`
	Actions  *[]Action        `json:"actions,omitempty"`
	Sections *[]Section       `json:"sections,omitempty"`
	Details  *DetailsCreation `json:"details,omitempty"`
}

type PermissionFilter struct {
	Name     string `json:"name,omitempty"`
	Category string `json:"category,omitempty"`
}

func NewPermissionRequest() Request {
	return Request{
		Flags:      &Flags{},
		Permission: &Permission{},
	}
}

func NewPermissionFilter() Request {
	return Request{
		Filter: &Filter{
			Permission: &PermissionFilter{},
		},
	}
}

func NewPermissionResult() Result {
	return Result{
		Errors:      &[]string{},
		Permissions: &[]Permission{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
