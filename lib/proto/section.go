/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Section struct {
	Id       string           `json:"id,omitempty"`
	Name     string           `json:"name,omitempty"`
	Category string           `json:"category,omitempty"`
	Details  *DetailsCreation `json:"details,omitempty"`
}

type SectionFilter struct {
	Name     string `json:"name,omitempty"`
	Category string `json:"category,omitempty"`
}

func NewSectionRequest() Request {
	return Request{
		Flags:   &Flags{},
		Section: &Section{},
	}
}

func NewSectionFilter() Request {
	return Request{
		Filter: &Filter{
			Section: &SectionFilter{},
		},
	}
}

func NewSectionResult() Result {
	return Result{
		Errors:   &[]string{},
		Sections: &[]Section{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
