/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package db

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

func (d *DB) SaveToken(expires, token string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(`tokens`))
		return b.Put([]byte(expires), []byte(token))
	})
}

func (d *DB) GetActiveToken() (string, error) {
	var token []byte
	if err := d.db.View(func(tx *bolt.Tx) error {
		var k, v []byte
		min := []byte(time.Now().UTC().Format(rfc3339Milli))

		c := tx.Bucket([]byte(`tokens`)).Cursor()
		k, v = c.Seek(min)
		if k != nil {
			token = make([]byte, len(v))
			copy(token, v)
		}
		return nil
	}); err != nil {
		return "", err
	}
	if token != nil {
		return string(token), nil
	}
	return "", fmt.Errorf(`Not found`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
