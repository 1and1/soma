/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"database/sql"
	"encoding/json"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"
)

type instance struct {
	input     chan msg.Request
	shutdown  chan bool
	conn      *sql.DB
	stmt_show *sql.Stmt
	stmt_list *sql.Stmt
	stmt_vers *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (i *instance) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.InstanceScopedList: i.stmt_list,
		stmt.InstanceShow:       i.stmt_show,
		stmt.InstanceVersions:   i.stmt_vers,
	} {
		if prepStmt, err = i.conn.Prepare(statement); err != nil {
			i.errLog.Fatal(`instance`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-i.shutdown:
			break runloop
		case req := <-i.input:
			go func() {
				i.process(&req)
			}()
		}
	}
}

//
func (i *instance) process(q *msg.Request) {
	result := msg.FromRequest(q)
	var (
		err                                           error
		version                                       int64
		isInherited                                   bool
		rows                                          *sql.Rows
		nullRepositoryId, nullBucketId                *sql.NullString
		instanceId, checkId, configId, details        string
		objectId, objectType, status, nextStatus      string
		repositoryId, bucketId, instanceConfigId      string
		createdNull, activatedNull, deprovisionedNull pq.NullTime
		updatedNull, notifiedNull                     pq.NullTime
	)

	nullRepositoryId.String = ``
	nullRepositoryId.Valid = false
	nullBucketId.String = ``
	nullBucketId.Valid = false

	switch q.Action {
	case `show`:
		i.reqLog.Printf(LogStrSRq, q.Section, q.Action, q.User,
			q.RemoteAddr)

		if err = i.stmt_show.QueryRow(q.Instance.Id).Scan(
			&instanceId,
			&version,
			&checkId,
			&configId,
			&instanceConfigId,
			&nullRepositoryId,
			&nullBucketId,
			&objectId,
			&objectType,
			&status,
			&nextStatus,
			&isInherited,
			&details,
		); err == sql.ErrNoRows {
			result.NotFound(err)
			goto dispatch
		} else if err != nil {
			result.ServerError(err)
			goto dispatch
		}

		// technically repositoryId can not be null, be we
		// may query it as NULL. This avoid "switching" types
		// between cases
		if nullRepositoryId.Valid {
			repositoryId = nullRepositoryId.String
		}
		if nullBucketId.Valid {
			bucketId = nullBucketId.String
		}

		// unmarhal JSONB deployment details
		depl := proto.Deployment{}
		if err = json.Unmarshal([]byte(details), &depl); err != nil {
			result.ServerError(err)
			goto dispatch
		}

		result.Instance = []proto.Instance{proto.Instance{
			Id:               instanceId,
			Version:          uint64(version),
			CheckId:          checkId,
			ConfigId:         configId,
			InstanceConfigId: instanceConfigId,
			RepositoryId:     repositoryId,
			BucketId:         bucketId,
			ObjectId:         objectId,
			ObjectType:       objectType,
			CurrentStatus:    status,
			NextStatus:       nextStatus,
			IsInherited:      isInherited,
			Deployment:       &depl,
		}}
		result.OK()
	case `versions`:
		i.reqLog.Printf(LogStrSRq, q.Section, q.Action, q.User,
			q.RemoteAddr)

		if rows, err = i.stmt_vers.Query(q.Instance.Id); err != nil {
			result.ServerError(err)
			goto dispatch
		}
		for rows.Next() {
			if err = rows.Scan(
				&instanceConfigId,
				&version,
				&instanceId,
				&createdNull,
				&activatedNull,
				&deprovisionedNull,
				&updatedNull,
				&notifiedNull,
				&status,
				&nextStatus,
				&isInherited,
			); err != nil {
				rows.Close()
				result.ServerError(err)
				result.Clear(q.Section)
				goto dispatch
			}
			inst := proto.Instance{
				InstanceConfigId: instanceConfigId,
				Version:          uint64(version),
				Id:               instanceId,
				CurrentStatus:    status,
				NextStatus:       nextStatus,
				IsInherited:      isInherited,
				Info: &proto.InstanceVersionInfo{
					// created timestamp is a not null column
					CreatedAt: createdNull.Time.Format(rfc3339Milli),
				},
			}
			if activatedNull.Valid {
				inst.Info.ActivatedAt = activatedNull.Time.Format(
					rfc3339Milli)
			}
			if deprovisionedNull.Valid {
				inst.Info.DeprovisionedAt = deprovisionedNull.Time.
					Format(rfc3339Milli)
			}
			if updatedNull.Valid {
				inst.Info.StatusLastUpdatedAt = updatedNull.Time.
					Format(rfc3339Milli)
			}
			if notifiedNull.Valid {
				inst.Info.NotifiedAt = notifiedNull.Time.Format(
					rfc3339Milli)
			}
			result.Instance = append(result.Instance, inst)
		}
		if err = rows.Err(); err != nil {
			result.ServerError(err)
			result.Clear(q.Section)
			goto dispatch
		}
		result.OK()
	case `list`:
		switch q.Instance.ObjectType {
		case `repository`:
			nullRepositoryId.String = q.Instance.ObjectId
			nullRepositoryId.Valid = true
		case `bucket`:
			nullBucketId.String = q.Instance.ObjectId
			nullBucketId.Valid = true
		}
		fallthrough
	case `instance_list_all`:
		// section: runtime
		i.reqLog.Printf(LogStrSRq, q.Section, q.Action, q.User,
			q.RemoteAddr)

		if rows, err = i.stmt_list.Query(
			nullRepositoryId,
			nullBucketId,
		); err != nil {
			result.ServerError(err)
			goto dispatch
		}

		for rows.Next() {
			if err = rows.Scan(
				&instanceId,
				&version,
				&checkId,
				&configId,
				&instanceConfigId,
				&nullRepositoryId,
				&nullBucketId,
				&objectId,
				&objectType,
				&status,
				&nextStatus,
				&isInherited,
			); err != nil {
				rows.Close()
				result.ServerError(err)
				result.Clear(q.Section)
				goto dispatch
			}

			if nullRepositoryId.Valid {
				repositoryId = nullRepositoryId.String
			}
			if nullBucketId.Valid {
				bucketId = nullBucketId.String
			}

			result.Instance = append(result.Instance, proto.Instance{
				Id:               instanceId,
				Version:          uint64(version),
				CheckId:          checkId,
				ConfigId:         configId,
				InstanceConfigId: instanceConfigId,
				RepositoryId:     repositoryId,
				BucketId:         bucketId,
				ObjectId:         objectId,
				ObjectType:       objectType,
				CurrentStatus:    status,
				NextStatus:       nextStatus,
				IsInherited:      isInherited,
			})
		}
		if err = rows.Err(); err != nil {
			result.ServerError(err)
			result.Clear(q.Section)
			goto dispatch
		}
		result.OK()
	default:
		result.UnknownRequest(q)
	}

dispatch:
	q.Reply <- result
}

/* Ops Access
 */
func (i *instance) shutdownNow() {
	i.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
