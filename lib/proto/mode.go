/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Mode struct {
	Mode    string       `json:"mode,omitempty"`
	Details *ModeDetails `json:"details,omitempty"`
}

type ModeDetails struct {
	DetailsCreation
}

func NewModeRequest() Request {
	return Request{
		Flags: &Flags{},
		Mode:  &Mode{},
	}
}

func NewModeResult() Result {
	return Result{
		Errors: &[]string{},
		Modes:  &[]Mode{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
