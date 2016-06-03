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

func (d *DB) Server(name, id, assetid string) error {
	if err := d.Open(); err != nil {
		return err
	}
	defer d.Close()

	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(`idcache`))
		n, err := b.CreateBucketIfNotExists([]byte(`servername`))
		a, err := b.CreateBucketIfNotExists([]byte(`serverasset`))
		if err != nil {
			return err
		}
		mapdata := map[string]string{
			`name`:    name,
			`id`:      id,
			`assetid`: assetid,
			`expire`:  time.Now().UTC().Add(time.Hour * 48).Format(time.RFC3339),
		}
		data, _ := json.Marshal(&mapdata)
		err = n.Put([]byte(name), data)
		if err != nil {
			return err
		}
		return a.Put([]byte(assetid), data)
	})
}

func (d *DB) ServerByName(name string) (map[string]string, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}
	defer d.Close()

	serverinfo := make(map[string]string)
	if err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(`idcache`)).Bucket([]byte(`servername`))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		data := b.Get([]byte(name))
		if data == nil {
			return bolt.ErrBucketNotFound
		}
		json.Unmarshal(data, &serverinfo)
		expire, _ := time.Parse(time.RFC3339, serverinfo[`expire`])
		if time.Now().UTC().After(expire.UTC()) {
			return bolt.ErrBucketNotFound
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if len(serverinfo) > 0 {
		return serverinfo, nil
	}
	return nil, bolt.ErrBucketNotFound
}

func (d *DB) ServerByAsset(asset string) (map[string]string, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}
	defer d.Close()

	serverinfo := make(map[string]string)
	if err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(`idcache`)).Bucket([]byte(`serverasset`))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		data := b.Get([]byte(asset))
		if data == nil {
			return bolt.ErrBucketNotFound
		}
		json.Unmarshal(data, &serverinfo)
		expire, _ := time.Parse(time.RFC3339, serverinfo[`expire`])
		if time.Now().UTC().After(expire.UTC()) {
			return bolt.ErrBucketNotFound
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if len(serverinfo) > 0 {
		return serverinfo, nil
	}
	return nil, bolt.ErrBucketNotFound
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
