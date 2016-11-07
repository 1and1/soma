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

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type workflowRead struct {
	input        chan msg.Request
	shutdown     chan bool
	conn         *sql.DB
	stmt_summary *sql.Stmt
	appLog       *log.Logger
	reqLog       *log.Logger
	errLog       *log.Logger
}

func (r *workflowRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.WorkflowSummary: r.stmt_summary,
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
		err    error
		status string
		count  int64
		rows   *sql.Rows
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
