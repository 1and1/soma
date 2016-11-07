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
	"fmt"
	"time"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"
)

type workflowRead struct {
	input        chan msg.Request
	shutdown     chan bool
	conn         *sql.DB
	stmt_summary *sql.Stmt
	stmt_list    *sql.Stmt
	appLog       *log.Logger
	reqLog       *log.Logger
	errLog       *log.Logger
}

func (r *workflowRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.WorkflowSummary: r.stmt_summary,
		stmt.WorkflowList:    r.stmt_list,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`workflow_r`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-r.shutdown:
			break runloop
		case req := <-r.input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

func (r *workflowRead) process(q *msg.Request) {
	result := msg.Result{Type: q.Type, Action: q.Action,
		Workflow: []proto.Workflow{}}
	var (
		err                                           error
		status, instanceId, checkId, repoId, configId string
		instanceConfigId                              string
		count, version                                int64
		rows                                          *sql.Rows
		activatedNull, deprovisionedNull              pq.NullTime
		updatedNull, notifiedNull                     pq.NullTime
		created                                       time.Time
		isInherited                                   bool
	)

	switch q.Action {
	case `summary`:
		r.reqLog.Printf(LogStrArg, q.Type, q.Action, q.User,
			q.RemoteAddr, q.Job.Id)
		summary := proto.WorkflowSummary{}

		if rows, err = r.stmt_summary.Query(); err != nil {
			result.ServerError(err)
			goto dispatch
		}
		for rows.Next() {
			if err = rows.Scan(
				&status,
				&count,
			); err != nil {
				rows.Close()
				result.ServerError(err)
				result.Clear(q.Type)
				goto dispatch
			}
			switch status {
			case `awaiting_computation`:
				summary.AwaitingComputation = uint64(count)
			case `computed`:
				summary.Computed = uint64(count)
			case `awaiting_rollout`:
				summary.AwaitingRollout = uint64(count)
			case `rollout_in_progress`:
				summary.RolloutInProgress = uint64(count)
			case `rollout_failed`:
				summary.RolloutFailed = uint64(count)
			case `active`:
				summary.Active = uint64(count)
			case `awaiting_deprovision`:
				summary.AwaitingDeprovision = uint64(count)
			case `deprovision_in_progress`:
				summary.DeprovisionInProgress = uint64(count)
			case `deprovision_failed`:
				summary.DeprovisionFailed = uint64(count)
			case `deprovisioned`:
				summary.Deprovisioned = uint64(count)
			case `awaiting_deletion`:
				summary.AwaitingDeletion = uint64(count)
			case `blocked`:
				summary.Blocked = uint64(count)
			}
		}
		if err = rows.Err(); err != nil {
			result.ServerError(err)
			result.Clear(q.Type)
			goto dispatch
		}
		result.Workflow = append(result.Workflow, proto.Workflow{
			Summary: &summary,
		})
		result.OK()

	case `list`:
		r.reqLog.Printf(LogStrArg, q.Type, q.Action, q.User,
			q.RemoteAddr, q.Job.Id)
		workflow := proto.Workflow{
			Instances: &[]proto.Instance{},
		}

		if rows, err = r.stmt_list.Query(
			q.Workflow.Status,
		); err != nil {
			result.ServerError(err)
			goto dispatch
		}
		for rows.Next() {
			if err = rows.Scan(
				&instanceId,
				&checkId,
				&repoId,
				&configId,
				&instanceConfigId,
				&version,
				&status,
				&created,
				&activatedNull,
				&deprovisionedNull,
				&updatedNull,
				&notifiedNull,
				&isInherited,
			); err != nil {
				rows.Close()
				result.ServerError(err)
				result.Clear(q.Type)
				goto dispatch
			}
			instance := proto.Instance{
				Id:               instanceId,
				CheckId:          checkId,
				RepositoryId:     repoId,
				ConfigId:         configId,
				InstanceConfigId: instanceConfigId,
				Version:          uint64(version),
				CurrentStatus:    status,
				IsInherited:      isInherited,
				Info: &proto.InstanceVersionInfo{
					CreatedAt: created.UTC().Format(rfc3339Milli),
				},
			}
			if activatedNull.Valid {
				instance.Info.ActivatedAt = activatedNull.
					Time.UTC().Format(rfc3339Milli)
			}
			if deprovisionedNull.Valid {
				instance.Info.DeprovisionedAt = deprovisionedNull.
					Time.UTC().Format(rfc3339Milli)
			}
			if updatedNull.Valid {
				instance.Info.StatusLastUpdatedAt = updatedNull.
					Time.UTC().Format(rfc3339Milli)
			}
			if notifiedNull.Valid {
				instance.Info.NotifiedAt = notifiedNull.
					Time.UTC().Format(rfc3339Milli)
			}
			*workflow.Instances = append(*workflow.Instances,
				instance)
		}
		if err = rows.Err(); err != nil {
			result.ServerError(err)
			result.Clear(q.Type)
			goto dispatch
		}
		result.Workflow = append(result.Workflow, workflow)
		result.OK()

	default:
		result.NotImplemented(fmt.Errorf(
			"Unknown requested action: %s/%s", q.Type, q.Action))
	}

dispatch:
	q.Reply <- result
}

func (r *workflowRead) shutdownNow() {
	r.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
