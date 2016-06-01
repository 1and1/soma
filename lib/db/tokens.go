/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package db

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

func (d *DB) SaveToken(user, valid, expires, token string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(`tokens`))
		u, err := b.CreateBucketIfNotExists([]byte(user))
		if err != nil {
			return err
		}
		mapdata := map[string]string{
			"valid":   valid,
			"expires": expires,
			"token":   token,
		}
		data, _ := json.Marshal(&mapdata)
		return u.Put([]byte(expires), data)
	})
}

func (d *DB) GetActiveToken(user string) (string, error) {
	var token string
	if err := d.db.View(func(tx *bolt.Tx) error {
		// build cursor seek position
		var k, v []byte
		min := []byte(time.Now().UTC().Format(rfc3339Milli))

		// open bucket for that user
		b := tx.Bucket([]byte(`tokens`)).Bucket([]byte(user))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		// seek an entry
		c := b.Cursor()
		k, v = c.Seek(min)
		if k != nil {
			data := make(map[string]string)
			json.Unmarshal(v, &data)
			token = data["token"]
		}
		return nil
	}); err != nil {
		return "", err
	}
	if token != "" {
		return token, nil
	}
	return "", bolt.ErrBucketNotFound
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
