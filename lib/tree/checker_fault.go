/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

//
// Interface: Checker
func (tef *Fault) SetCheck(c Check) {
}

func (tef *Fault) setCheckInherited(c Check) {
}

func (tef *Fault) setCheckOnChildren(c Check) {
}

func (tef *Fault) addCheck(c Check) {
}

func (tef *Fault) DeleteCheck(c Check) {
}

func (tef *Fault) deleteCheckInherited(c Check) {
}

func (tef *Fault) deleteCheckOnChildren(c Check) {
}

func (tef *Fault) rmCheck(c Check) {
}

func (tef *Fault) syncCheck(childId string) {
}

func (tef *Fault) checkCheck(checkId string) bool {
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
