/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Oncall struct {
	Id      string          `json:"id,omitempty"`
	Name    string          `json:"name,omitempty"`
	Number  string          `json:"number,omitempty"`
	Members *[]OncallMember `json:"members,omitempty"`
	Details *OncallDetails  `json:"details,omitempty"`
}

type OncallDetails struct {
	DetailsCreation
}

type OncallMember struct {
	UserName string `json:"userName,omitempty"`
	UserId   string `json:"userId,omitempty"`
}

type OncallFilter struct {
	Name   string `json:"name,omitempty"`
	Number string `json:"number,omitempty"`
}

func (p *Oncall) DeepCompare(a *Oncall) bool {
	if p.Id != a.Id || p.Name != a.Name || p.Number != a.Number {
		return false
	}
	return true
}

func NewOncallRequest() Request {
	return Request{
		Flags:  &Flags{},
		Oncall: &Oncall{},
	}
}

func NewOncallFilter() Request {
	return Request{
		Filter: &Filter{
			Oncall: &OncallFilter{},
		},
	}
}

func NewOncallResult() Result {
	return Result{
		Errors:  &[]string{},
		Oncalls: &[]Oncall{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
