/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package db

import "github.com/boltdb/bolt"

func (d *DB) SaveJob(jid, jtype string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(`jobs`)).Bucket([]byte(`active`))
		return b.Put([]byte(jid), []byte(jtype))
	})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
