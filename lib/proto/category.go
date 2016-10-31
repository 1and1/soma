/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Category struct {
	Name    string           `json:"name,omitempty"`
	Details *CategoryDetails `json:"details,omitempty"`
}

type CategoryDetails struct {
	CreatedAt string `json:"createdAt,omitempty"`
	CreatedBy string `json:"createdBy,omitempty"`
}

func NewCategoryRequest() Request {
	return Request{
		Flags:    &Flags{},
		Category: &Category{},
	}
}

func NewCategoryResult() Result {
	return Result{
		Errors:     &[]string{},
		Categories: &[]Category{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
