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
	Deployment  somaproto.DeploymentDetails
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
}

func (self *somaDeploymentHandler) run() {
	var err error

	log.Println("Prepare: deployment/get")
	if self.get_stmt, err = self.conn.Prepare(stmtGetDeployment); err != nil {
		log.Fatal("deployment/get: ", err)
	}
	defer self.get_stmt.Close()

	log.Println("Prepare: deployment/update")
	if self.upd_stmt, err = self.conn.Prepare(stmtUpdateDeployment); err != nil {
		log.Fatal("deployment/update: ", err)
	}
	defer self.upd_stmt.Close()

	log.Println("Prepare: deployment/status")
	if self.sta_stmt, err = self.conn.Prepare(stmtDeploymentStatus); err != nil {
		log.Fatal("deployment/status: ", err)
	}
	defer self.sta_stmt.Close()

	log.Println("Prepare: deployment/activate")
	if self.act_stmt, err = self.conn.Prepare(stmtActivateDeployment); err != nil {
		log.Fatal("deployment/activate: ", err)
	}
	defer self.act_stmt.Close()

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
		instanceConfigID, status, next, details, nextNG string
		err                                             error
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

		depl := somaproto.DeploymentDetails{}
		if err = json.Unmarshal([]byte(details), depl); err != nil {
			result.Append(err, &somaDeploymentResult{})
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
				Deployment: somaproto.DeploymentDetails{
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
				Deployment: somaproto.DeploymentDetails{
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
				Deployment: somaproto.DeploymentDetails{
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
				Deployment: somaproto.DeploymentDetails{
					Task: "deprovision",
				},
			})
		default:
			result.SetRequestError(fmt.Errorf("Illegal current state for state update"))
		}

	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
