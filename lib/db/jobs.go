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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"


	"github.com/boltdb/bolt"
)

func (d *DB) SaveJob(jid, jtype string) error {
	if err := d.Open(); err != nil {
		return err
	}
	defer d.Close()
	now := time.Now().UTC().Format(rfc3339Milli)

	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(`jobs`)).Bucket([]byte(`active`))
		id, _ := b.NextSequence()
		return b.Put(
			uitob(id),
			[]byte(fmt.Sprintf("%s|%s|%s", jid, now, jtype)),
		)
	})
}

func (d *DB) FinishJob(id uint64, job *proto.Job) error {
	if err := d.Open(); err != nil {
		return err
	}
	defer d.Close()

	jid := job.Id
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return d.db.Update(func(tx *bolt.Tx) error {
		// delete from active list
		bActive := tx.Bucket([]byte(`jobs`)).Bucket([]byte(`active`))
		if err := bActive.Delete(uitob(id)); err != nil {
			return err
		}

		// insert into finished list
		bFinished := tx.Bucket([]byte(`jobs`)).Bucket([]byte(`finished`))
		if err := bFinished.Put(uitob(id), []byte(jid)); err != nil {
			return err
		}

		// save job data
		bData := tx.Bucket([]byte(`jobs`)).Bucket([]byte(`data`))
		return bData.Put([]byte(jid), data)
	})
}

// Slice of slice of strings, with the inner slice contents as follows:
// 0: storage key
// 1: JobID
// 2: Timestamp
// 3: JobType
func (d *DB) ActiveJobs() ([][]string, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}
	defer d.Close()

	count := 0
	jobs := [][]string{}
	if err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(`jobs`)).Bucket([]byte(`active`))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			count++
			r := []string{bots(k)}
			r = append(r, strings.Split(string(v), `|`)...)
			jobs = append(jobs, r)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if count > 0 {
		return jobs, nil
	}
	return nil, bolt.ErrBucketNotFound
}

func (d *DB) FinishedJobs() ([]proto.Job, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}
	defer d.Close()

	count := 0
	jobs := []proto.Job{}
	if err := d.db.View(func(tx *bolt.Tx) error {
		bF := tx.Bucket([]byte(`jobs`)).Bucket([]byte(`finished`))
		bD := tx.Bucket([]byte(`jobs`)).Bucket([]byte(`data`))
		c := bF.Cursor()

		// iterate over index bucket jobs.finished
		for k, v := c.First(); k != nil; k, v = c.Next() {
			// jobs.data bucket key == jobs.finished value
			val := bD.Get(v)
			if val == nil {
				continue
			}
			count++
			j := proto.Job{}
			if err := json.Unmarshal(val, &j); err != nil {
				return err
			}
			jobs = append(jobs, j)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if count > 0 {
		return jobs, nil
	}
	return nil, bolt.ErrBucketNotFound
}

func uitob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func botui(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func bots(b []byte) string {
	return strconv.FormatUint(binary.BigEndian.Uint64(b), 10)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
