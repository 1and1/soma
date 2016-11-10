/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Flags struct {
	Restore  bool `json:"restore"`
	Purge    bool `json:"purge"`
	Freeze   bool `json:"freeze"`
	Thaw     bool `json:"thaw"`
	Clear    bool `json:"clear"`    // repository
	Activate bool `json:"activate"` // repository
	Detailed bool `json:"detailed"` // jobs
	Forced   bool `json:"forced"`   // workflow
	Add      bool `json:"add"`      // permission map
	Remove   bool `json:"remove"`   // permission unmap
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
