/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

// Package soma implements the application handlers of the SOMA
// service.
package soma

type Soma struct {
	handlerMap   *map[string]interface{}
}

func New(
	appHandlerMap *map[string]interface{},
) *Soma {
	s := Soma{}
	s.handlerMap = appHandlerMap
	return &s
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
