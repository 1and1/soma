/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm // import "github.com/1and1/soma/internal/perm"

import (
	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
)

// isAuthorized implements Cache.IsAuthorized and checks if the
// request is authorized
func (c *Cache) isAuthorized(q *msg.Request) msg.Result {
	result := msg.FromRequest(q)
	// default action is to deny
	result.Super = &msg.Supervisor{
		Verdict:      401,
		VerdictAdmin: false,
	}
	var user *proto.User

	// set readlock on the cache
	c.lock.RLock()
	defer c.lock.RUnlock()

	// look up the user
	if user = c.user.getByName(q.User); user == nil {
		goto dispatch
	}

	// check if the user has omnipotence
	if c.grantGlobal.assess(
		`user`, // TODO lookup usertype
		user.Id,
		`omnipotence`,
		`00000000-0000-0000-0000-000000000000`,
	) {
		result.Super.Verdict = 200
		result.Super.VerdictAdmin = true
		goto dispatch
	}

dispatch:
	return result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
