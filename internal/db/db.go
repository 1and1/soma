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
	db         *bolt.DB
	open       bool
	ensure     bool
	configured bool
	Path       string
	Mode       os.FileMode
	Options    *bolt.Options
}

func (d *DB) Configure(p string, m os.FileMode, o *bolt.Options) {
	d.Path = p
	d.Mode = m
	d.Options = o
	d.configured = true
}

func (d *DB) Open() error {
	var err error
	if d.open {
		return nil
	}
	if !d.configured {
		return fmt.Errorf(`DB is not configured`)
	}

	if d.db, err = bolt.Open(d.Path, d.Mode, d.Options); err == nil {
		d.open = true
	}
	return err
}

func (d *DB) Close() error {
	var err error
	if !d.open {
		return nil
	}

	if err = d.db.Close(); err == nil {
		d.open = false
	}
	return err
}

func (d *DB) EnsureBuckets() error {
	if err := d.Open(); err != nil {
		return err
	}
	defer d.Close()

	for _, buck := range []string{"jobs", "tokens", "idcache"} {
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
			case `idcache`:
				b := tx.Bucket([]byte(buck))
				if _, err := b.CreateBucketIfNotExists([]byte(`team`)); err != nil {
					return fmt.Errorf("Failed to create DB bucket: %s", err)
				}
				if _, err := b.CreateBucketIfNotExists([]byte(`servername`)); err != nil {
					return fmt.Errorf("Failed to create DB bucket: %s", err)
				}
				if _, err := b.CreateBucketIfNotExists([]byte(`serverasset`)); err != nil {
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
