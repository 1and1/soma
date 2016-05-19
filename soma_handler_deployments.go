package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

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
}

func (self *somaDeploymentHandler) run() {
	var err error

	if self.get_stmt, err = self.conn.Prepare(stmtGetDeployment); err != nil {
		log.Fatal("deployment/get: ", err)
	}
	defer self.get_stmt.Close()

	if self.upd_stmt, err = self.conn.Prepare(stmtUpdateDeployment); err != nil {
		log.Fatal("deployment/update: ", err)
	}
	defer self.upd_stmt.Close()

	if self.sta_stmt, err = self.conn.Prepare(stmtDeploymentStatus); err != nil {
		log.Fatal("deployment/status: ", err)
	}
	defer self.sta_stmt.Close()

	if self.act_stmt, err = self.conn.Prepare(stmtActivateDeployment); err != nil {
		log.Fatal("deployment/activate: ", err)
	}
	defer self.act_stmt.Close()

	if self.lst_stmt, err = self.conn.Prepare(stmtGetDeploymentList); err != nil {
		log.Fatal("deployment/list: ", err)
	}
	defer self.lst_stmt.Close()

	if self.all_stmt, err = self.conn.Prepare(stmtGetAllDeploymentList); err != nil {
		log.Fatal("deployment/listall: ", err)
	}
	defer self.all_stmt.Close()

	if self.clr_stmt, err = self.conn.Prepare(stmtDeployClearFlag); err != nil {
		log.Fatal("deployment/clearflag: ", err)
	}
	defer self.clr_stmt.Close()

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
		instanceConfigID, instanceID, status, next, details, nextNG string
		err                                                         error
		list                                                        *sql.Rows
	)
	result := somaResult{}

	switch q.action {
	case "get":
		log.Printf("R: deployment/get for %s", q.Deployment)
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

		switch status {
		case "awaiting_rollout":
			next = "rollout_in_progress"
			nextNG = "active"
			depl.Task = "rollout"
		case "rollout_in_progress":
			next = "rollout_in_progress"
			nextNG = "active"
			depl.Task = "rollout"
		case "active":
			next = "active"
			nextNG = "none"
			depl.Task = "rollout"
		case "rollout_failed":
			next = "rollout_in_progress"
			nextNG = "active"
			depl.Task = "rollout"
		case "awaiting_deprovision":
			next = "deprovision_in_progress"
			nextNG = "deprovisioned"
			depl.Task = "deprovision"
		case "deprovision_in_progress":
			next = "deprovision_in_progress"
			nextNG = "deprovisioned"
			depl.Task = "deprovision"
		case "deprovision_failed":
			next = "deprovision_in_progress"
			nextNG = "deprovisioned"
			depl.Task = "deprovision"
		}

		result.Append(err, &somaDeploymentResult{
			Deployment: depl,
		})
		self.upd_stmt.Exec(next, nextNG, instanceConfigID)
	case "update/success":
		log.Printf("R: deployment/update/success for %s", q.Deployment)
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
			self.upd_stmt.Exec(
				next,
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
	case "update/failed":
		log.Printf("R: deployment/update/failed for %s", q.Deployment)
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
		log.Printf("R: deployment/list for %s", q.Deployment)
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
		log.Printf("R: deployment/listall for %s", q.Deployment)
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
		result.SetNotImplemented()
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
