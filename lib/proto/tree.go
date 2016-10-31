/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Tree struct {
	Id         string      `json:"id,omitempty"`
	Type       string      `json:"type,omitempty"`
	Repository *Repository `json:"repository,omitempty"`
	Bucket     *Bucket     `json:"bucket,omitempty"`
	Group      *Group      `json:"group,omitempty"`
	Cluster    *Cluster    `json:"cluster,omitempty"`
	Node       *Node       `json:"node,omitempty"`
}

func NewTreeResult() Result {
	return Result{
		Errors: &[]string{},
		Tree:   &Tree{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
