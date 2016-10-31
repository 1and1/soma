/*-
 * Copyright (c) 2015-2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Repository struct {
	Id         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	TeamId     string      `json:"teamId,omitempty"`
	IsDeleted  bool        `json:"isDeleted,omitempty"`
	IsActive   bool        `json:"isActive,omitempty"`
	Members    *[]Bucket   `json:"members,omitempty"`
	Details    *Details    `json:"details,omitempty"`
	Properties *[]Property `json:"properties,omitempty"`
}

type RepositoryFilter struct {
	Name      string `json:"name,omitempty"`
	TeamId    string `json:"teamId,omitempty"`
	IsDeleted bool   `json:"isDeleted,omitempty"`
	IsActive  bool   `json:"isActive,omitempty"`
}

func NewRepositoryRequest() Request {
	return Request{
		Flags:      &Flags{},
		Repository: &Repository{},
	}
}

func NewRepositoryFilter() Request {
	return Request{
		Filter: &Filter{
			Repository: &RepositoryFilter{},
		},
	}
}

func NewRepositoryResult() Result {
	return Result{
		Errors:       &[]string{},
		Repositories: &[]Repository{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
