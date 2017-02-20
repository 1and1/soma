/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import "sync"

// HandlerMap is a concurrent map that is used to look up input
// channels for application handlers
type HandlerMap struct {
	hmap map[string]interface{}
	sync.RWMutex
}

// Add registers a new handler
func (h *HandlerMap) Add(key string, value interface{}) {
	h.Lock()
	defer h.Unlock()
	h.hmap[key] = value
}

// Get retrieves a handler
func (h *HandlerMap) Get(key string) interface{} {
	h.RLock()
	defer h.RUnlock()
	return h.hmap[key]
}

// Del removes a handler
func (h *HandlerMap) Del(key string) {
	h.Lock()
	defer h.Unlock()
	delete(h.hmap, key)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
