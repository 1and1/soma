/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

/*
 * Fault Handler Special Sauce
 *
 * Elemnts return pointers to the Fault Handler instead of nil pointers
 * when asked for something they do not have.
 *
 * This makes these chains safe:
 *		<foo>.Parent.(Receiver).GetBucket().Unlink()
 *
 * Instead of nil, the parent returns the Fault handler which implements
 * Receiver and Unlinker. Due to the information in the
 * Receive-/UnlinkRequest, it can log what went wrong.
 *
 */

//
// Interface: Unlinker
func (tef *Fault) Unlink(u UnlinkRequest) {
	panic(`Fault.Unlink`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
