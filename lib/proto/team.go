/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Team struct {
	Id       string       `json:"id,omitempty"`
	Name     string       `json:"name,omitempty"`
	LdapId   string       `json:"ldapId,omitempty"`
	IsSystem bool         `json:"isSystem,omitempty"`
	Details  *TeamDetails `json:"details,omitempty"`
}

type TeamDetails struct {
	DetailsCreation
}

type TeamFilter struct {
	Name     string `json:"name,omitempty"`
	LdapId   string `json:"ldapId,omitempty"`
	IsSystem bool   `json:"isSystem,omitempty"`
}

func (p *Team) DeepCompare(a *Team) bool {
	if p.Id != a.Id || p.Name != a.Name || p.LdapId != a.LdapId || p.IsSystem != a.IsSystem {
		return false
	}
	return true
}

func NewTeamRequest() Request {
	return Request{
		Flags: &Flags{},
		Team:  &Team{},
	}
}

func NewTeamFilter() Request {
	return Request{
		Filter: &Filter{
			Team: &TeamFilter{},
		},
	}
}

func NewTeamResult() Result {
	return Result{
		Errors: &[]string{},
		Teams:  &[]Team{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
