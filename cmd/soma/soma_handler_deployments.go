package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type somaDeploymentRequest struct {
	action     string
	Deployment string
	reply      chan somaResult
}

type somaDeploymentResult struct {
	ResultError error
	ListEntry   string
	Deployment  proto.Deployment
}

func (a *somaDeploymentResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Deployments = append(r.Deployments, somaDeploymentResult{ResultError: err})
	}
}

func (a *somaDeploymentResult) SomaAppendResult(r *somaResult) {
	r.Deployments = append(r.Deployments, *a)
}

type somaDeploymentHandler struct {
	input    chan somaDeploymentRequest
	shutdown chan bool
	conn     *sql.DB
	get_stmt *sql.Stmt
	upd_stmt *sql.Stmt
	sta_stmt *sql.Stmt
	act_stmt *sql.Stmt
	lst_stmt *sql.Stmt
	all_stmt *sql.Stmt
	clr_stmt *sql.Stmt
	dpr_stmt *sql.Stmt
	sty_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (self *somaDeploymentHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.DeploymentGet:              self.get_stmt,
		stmt.DeploymentUpdate:           self.upd_stmt,
		stmt.DeploymentStatus:           self.sta_stmt,
		stmt.DeploymentActivate:         self.act_stmt,
		stmt.DeploymentList:             self.lst_stmt,
		stmt.DeploymentListAll:          self.all_stmt,
		stmt.DeploymentClearFlag:        self.clr_stmt,
		stmt.DeploymentDeprovision:      self.dpr_stmt,
		stmt.DeploymentDeprovisionStyle: self.sty_stmt,
	} {
		if prepStmt, err = self.conn.Prepare(statement); err != nil {
			self.errLog.Fatal(`deployment`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-self.shutdown:
			break runloop
		case req := <-self.input:
			go func() {
				self.process(&req)
			}()
		}
	}
}

func (self *somaDeploymentHandler) process(q *somaDeploymentRequest) {
	var (
		instanceConfigID, instanceID, status string
		next, details, nextNG, deprStyle     string
		err                                  error
		list                                 *sql.Rows
		updated, blocksRollout               bool
	)
	result := somaResult{}

	switch q.action {
	case "get":
		self.appLog.Printf("R: deployment/get for %s", q.Deployment)
		if err = self.get_stmt.QueryRow(q.Deployment).Scan(
			&instanceConfigID,
			&status,
			&next,
			&details,
		); err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
				q.reply <- result
				return
			}
			result.SetRequestError(err)
			q.reply <- result
			return
		}

		depl := proto.Deployment{}
		if err = json.Unmarshal([]byte(details), &depl); err != nil {
			result.SetRequestError(err)
			q.reply <- result
			return
		}

		// returns true if there is a updated version blocked, ie.
		// after this deprovisioning a new version will be rolled out
		if err = self.sty_stmt.QueryRow(q.Deployment).Scan(
			&blocksRollout,
		); result.SetRequestError(err) {
			q.reply <- result
			return
		}
		deprStyle = `deprovision`
		if !blocksRollout {
			deprStyle = `delete`
		}

		switch status {
		case "awaiting_rollout":
			next = "rollout_in_progress"
			nextNG = "active"
			depl.Task = "rollout"
			updated = true
		case "rollout_in_progress":
			depl.Task = "rollout"
			updated = false
		case "active":
			depl.Task = "rollout"
			updated = false
		case "rollout_failed":
			next = "rollout_in_progress"
			nextNG = "active"
			depl.Task = "rollout"
			updated = true
		case "awaiting_deprovision":
			next = "deprovision_in_progress"
			nextNG = "deprovisioned"
			depl.Task = deprStyle
			updated = true
		case "deprovision_in_progress":
			depl.Task = deprStyle
			updated = false
		case "deprovision_failed":
			next = "deprovision_in_progress"
			nextNG = "deprovisioned"
			depl.Task = deprStyle
			updated = true
		case `deprovisioned`:
			depl.Task = deprStyle
			updated = false
		}

		result.Append(err, &somaDeploymentResult{
			Deployment: depl,
		})
		if updated {
			self.upd_stmt.Exec(next, nextNG, instanceConfigID)
		}
	case "update/success":
		self.appLog.Printf("R: deployment/update/success for %s", q.Deployment)
		if err = self.sta_stmt.QueryRow(q.Deployment).Scan(
			&instanceConfigID,
			&status,
			&next,
		); err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
				q.reply <- result
				return
			}
			result.SetRequestError(err)
			q.reply <- result
			return
		}
		switch status {
		case "rollout_in_progress":
			self.act_stmt.Exec(
				next,
				"none",
				time.Now().UTC(),
				instanceConfigID,
			)
			result.Append(nil, &somaDeploymentResult{
				Deployment: proto.Deployment{
					Task: "rollout",
				},
			})
		case "deprovision_in_progress":
			self.dpr_stmt.Exec(
				next,
				"none",
				time.Now().UTC(),
				instanceConfigID,
			)
			result.Append(nil, &somaDeploymentResult{
				Deployment: proto.Deployment{
					Task: "deprovision",
				},
			})
		default:
			result.SetRequestError(fmt.Errorf("Illegal current state for state update"))
		}
	case "update/failed":
		self.appLog.Printf("R: deployment/update/failed for %s", q.Deployment)
		if err = self.sta_stmt.QueryRow(q.Deployment).Scan(
			&instanceConfigID,
			&status,
			&next,
		); err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
				q.reply <- result
				return
			}
			result.SetRequestError(err)
			q.reply <- result
			return
		}
		switch status {
		case "rollout_in_progress":
			self.upd_stmt.Exec(
				"rollout_failed",
				"none",
				instanceConfigID,
			)
			result.Append(nil, &somaDeploymentResult{
				Deployment: proto.Deployment{
					Task: "rollout",
				},
			})
		case "deprovision_in_progress":
			self.upd_stmt.Exec(
				"deprovision_failed",
				"none",
				instanceConfigID,
			)
			result.Append(nil, &somaDeploymentResult{
				Deployment: proto.Deployment{
					Task: "deprovision",
				},
			})
		default:
			result.SetRequestError(fmt.Errorf("Illegal current state for state update"))
		}
	case "list":
		self.appLog.Printf("R: deployment/list for %s", q.Deployment)
		if list, err = self.lst_stmt.Query(q.Deployment); err != nil {
			result.SetRequestError(err)
			q.reply <- result
			return
		}

		for list.Next() {
			if err = list.Scan(
				&instanceID,
			); err != nil {
				if err == sql.ErrNoRows {
					result.SetNotFound()
					q.reply <- result
					return
				}
				result.SetRequestError(err)
				q.reply <- result
				return
			}

			result.Append(nil, &somaDeploymentResult{
				ListEntry: instanceID,
			})
			self.clr_stmt.Exec(instanceID)
		}
	case "listall":
		self.appLog.Printf("R: deployment/listall for %s", q.Deployment)
		if list, err = self.all_stmt.Query(q.Deployment); err != nil {
			result.SetRequestError(err)
			q.reply <- result
			return
		}

		for list.Next() {
			if err = list.Scan(
				&instanceID,
			); err != nil {
				if err == sql.ErrNoRows {
					result.SetNotFound()
					q.reply <- result
					return
				}
				result.SetRequestError(err)
				q.reply <- result
				return
			}

			result.Append(nil, &somaDeploymentResult{
				ListEntry: instanceID,
			})
			self.clr_stmt.Exec(instanceID)
		}
	default:
		self.errLog.Printf("R: unimplemented deployment/%s", q.action)
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Ops Access
 */
func (dh *somaDeploymentHandler) shutdownNow() {
	dh.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
