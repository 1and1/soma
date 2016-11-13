/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Action struct {
	Id          string           `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	SectionId   string           `json:"sectionId,omitempty"`
	SectionName string           `json:"sectionName,omitempty"`
	Category    string           `json:"category,omitempty"`
	Details     *DetailsCreation `json:"details,omitempty"`
}

type ActionFilter struct {
	Name      string `json:"name,omitempty"`
	SectionId string `json:"sectionId,omitempty"`
}

func NewActionRequest() Request {
	return Request{
		Flags:  &Flags{},
		Action: &Action{},
	}
}

func NewActionFilter() Request {
	return Request{
		Filter: &Filter{
			Action: &ActionFilter{},
		},
	}
}

func NewActionResult() Result {
	return Result{
		Errors:  &[]string{},
		Actions: &[]Action{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
