/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type SystemOperation struct {
	Request      string `json:"request,omitempty"`
	RepositoryId string `json:"repositoryId,omitempty"`
	RebuildLevel string `json:"rebuildLevel,omitempty"`
}

func NewSystemOperationRequest() Request {
	return Request{
		Flags:           &Flags{},
		SystemOperation: &SystemOperation{},
	}
}

func NewSystemOperationResult() Result {
	return Result{
		Errors:           &[]string{},
		SystemOperations: &[]SystemOperation{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
