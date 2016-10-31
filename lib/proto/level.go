/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Level struct {
	Name      string        `json:"name,omitempty"`
	ShortName string        `json:"shortName,omitempty"`
	Numeric   uint16        `json:"numeric,omitempty"`
	Details   *LevelDetails `json:"details,omitempty"`
}

type LevelFilter struct {
	Name      string `json:"name,omitempty"`
	ShortName string `json:"shortName,omitempty"`
	Numeric   uint16 `json:"numeric,omitempty"`
}

type LevelDetails struct {
	DetailsCreation
}

func NewLevelRequest() Request {
	return Request{
		Flags: &Flags{},
		Level: &Level{},
	}
}

func NewLevelFilter() Request {
	return Request{
		Filter: &Filter{
			Level: &LevelFilter{},
		},
	}
}

func NewLevelResult() Result {
	return Result{
		Errors: &[]string{},
		Levels: &[]Level{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
