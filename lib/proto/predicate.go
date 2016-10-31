/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Predicate struct {
	Symbol  string            `json:"symbol,omitempty"`
	Details *PredicateDetails `json:"details,omitempty"`
}

type PredicateDetails struct {
	DetailsCreation
}

func NewPredicateRequest() Request {
	return Request{
		Flags:     &Flags{},
		Predicate: &Predicate{},
	}
}

func NewPredicateResult() Result {
	return Result{
		Errors:     &[]string{},
		Predicates: &[]Predicate{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
