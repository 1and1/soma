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
	"os"

	"github.com/boltdb/bolt"
)

const rfc3339Milli string = "2006-01-02T15:04:05.000Z07:00"

type DB struct {
	db     *bolt.DB
	open   bool
	ensure bool
}

func (d *DB) Open(f string, m os.FileMode, opt *bolt.Options) error {
	var err error
	if d.open {
		return nil
	}

	if d.db, err = bolt.Open(f, m, opt); err == nil {
		d.open = true
	}
	return err
}

func (d *DB) Close() error {
	if !d.open {
		return nil
	}

	return d.db.Close()
}

func (d *DB) EnsureBuckets() error {
	if !d.open {
		return fmt.Errorf("DB is not open")
	}

	for _, buck := range []string{"jobs", "tokens"} {
		err := d.db.Update(func(tx *bolt.Tx) error {
			if _, err := tx.CreateBucketIfNotExists([]byte(buck)); err != nil {
				return fmt.Errorf("Failed to create DB bucket: %s", err)
			}
			switch buck {
			case `jobs`:
				b := tx.Bucket([]byte(buck))
				if _, err := b.CreateBucketIfNotExists([]byte(`active`)); err != nil {
					return fmt.Errorf("Failed to create DB bucket: %s", err)
				}
				if _, err := b.CreateBucketIfNotExists([]byte(`finished`)); err != nil {
					return fmt.Errorf("Failed to create DB bucket: %s", err)
				}
				if _, err := b.CreateBucketIfNotExists([]byte(`data`)); err != nil {
					return fmt.Errorf("Failed to create DB bucket: %s", err)
				}
			case `tokens`:
				b := tx.Bucket([]byte(buck))
				if _, err := b.CreateBucketIfNotExists([]byte(`user`)); err != nil {
					return fmt.Errorf("Failed to create DB bucket: %s", err)
				}
				if _, err := b.CreateBucketIfNotExists([]byte(`admin`)); err != nil {
					return fmt.Errorf("Failed to create DB bucket: %s", err)
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	d.ensure = true
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
