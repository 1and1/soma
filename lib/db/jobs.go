/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package db

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

func (d *DB) SaveJob(jid, jtype string) error {
	now := time.Now().UTC().Format(rfc3339Milli)

	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(`jobs`)).Bucket([]byte(`active`))
		id, _ := b.NextSequence()
		return b.Put(
			uitob(id),
			[]byte(fmt.Sprintf("id:%s|time:%s|type:%s", jid, now, jtype)),
		)
	})
}

func uitob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func botui(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
