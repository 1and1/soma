package main

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type somaHostDeploymentRequest struct {
	action  string
	system  string
	assetid int64
	idlist  []string
	reply   chan somaResult
}

type somaHostDeploymentResult struct {
	ResultError error
	Delete      bool
	DeleteId    string
	Deployment  proto.Deployment
}

func (h *somaHostDeploymentResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.HostDeployments = append(r.HostDeployments, somaHostDeploymentResult{ResultError: err})
	}
}

func (h *somaHostDeploymentResult) SomaAppendResult(r *somaResult) {
	r.HostDeployments = append(r.HostDeployments, *h)
}

type somaHostDeploymentHandler struct {
	input     chan somaHostDeploymentRequest
	shutdown  chan bool
	conn      *sql.DB
	geti_stmt *sql.Stmt
	last_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (self *somaHostDeploymentHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.DeploymentInstancesForNode:    self.geti_stmt,
		stmt.DeploymentLastInstanceVersion: self.last_stmt,
	} {
		if prepStmt, err = self.conn.Prepare(statement); err != nil {
			self.errLog.Fatal(`hostdeployment`, err, statement)
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

func (self *somaHostDeploymentHandler) process(q *somaHostDeploymentRequest) {
	var (
		checkInstanceID, deploymentDetails, status string
		idList                                     *sql.Rows
		err                                        error
	)
	result := somaResult{}

	switch q.action {
	case "get":
		self.reqLog.Printf("R: hostdeployment/get-for-node for %d", q.assetid)
		if idList, err = self.geti_stmt.Query(
			q.assetid,
			q.system,
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
		defer idList.Close()

	idgetloop:
		for idList.Next() {
			if err = idList.Scan(&checkInstanceID); err != nil {
				if err == sql.ErrNoRows {
					result.SetNotFound()
					q.reply <- result
					return
				}
				result.Append(err, &somaHostDeploymentResult{})
				continue idgetloop
			}
			err = self.last_stmt.QueryRow(checkInstanceID).Scan(&deploymentDetails, &status)
			if err != nil {
				result.Append(err, &somaHostDeploymentResult{})
				continue idgetloop
			}
			depl := proto.Deployment{}
			if err = json.Unmarshal([]byte(deploymentDetails), &depl); err != nil {
				result.Append(err, &somaHostDeploymentResult{})
				continue idgetloop
			}
			switch status {
			case "awaiting_rollout":
				depl.Task = "rollout"
			case "rollout_in_progress":
				depl.Task = "rollout"
			case "active":
				depl.Task = "rollout"
			case "rollout_failed":
				depl.Task = "rollout"
			case "awaiting_deprovision":
				depl.Task = "deprovision"
			case "deprovision_in_progress":
				depl.Task = "deprovision"
			case "deprovision_failed":
				depl.Task = "deprovision"
			default:
				depl.Task = "pending"
			}
			// remove credentials from the hostapi
			for i, _ := range depl.Service.Attributes {
				if strings.Contains(depl.Service.Attributes[i].Name, "credential_") {
					depl.Service.Attributes[i].Value = ""
				}
			}
			result.Append(nil, &somaHostDeploymentResult{
				Deployment: depl,
			})
		}
	case "assemble":
		idMap := map[string]bool{}
		self.reqLog.Printf("R: hostdeployment/get-for-node for %d", q.assetid)
		if idList, err = self.geti_stmt.Query(
			q.assetid,
			q.system,
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
		defer idList.Close()

	assembleloop:
		for idList.Next() {
			if err = idList.Scan(&checkInstanceID); err != nil {
				if err == sql.ErrNoRows {
					result.SetNotFound()
					q.reply <- result
					return
				}
				result.Append(err, &somaHostDeploymentResult{})
				continue assembleloop
			}
			idMap[checkInstanceID] = true
			err = self.last_stmt.QueryRow(checkInstanceID).Scan(&deploymentDetails, &status)
			if err != nil {
				result.Append(err, &somaHostDeploymentResult{})
				continue assembleloop
			}
			depl := proto.Deployment{}
			if err = json.Unmarshal([]byte(deploymentDetails), &depl); err != nil {
				result.Append(err, &somaHostDeploymentResult{})
				continue assembleloop
			}
			switch status {
			case "awaiting_rollout":
				depl.Task = "rollout"
			case "rollout_in_progress":
				depl.Task = "rollout"
			case "active":
				depl.Task = "rollout"
			case "rollout_failed":
				depl.Task = "rollout"
			case "blocked":
				depl.Task = "rollout"
			default:
				// bump this id to the delete list
				delete(idMap, checkInstanceID)
				continue assembleloop
			}
			// remove credentials from the hostapi
			for i, _ := range depl.Service.Attributes {
				if strings.Contains(depl.Service.Attributes[i].Name, "credential_") {
					depl.Service.Attributes[i].Value = ""
				}
			}
			result.Append(nil, &somaHostDeploymentResult{
				Deployment: depl,
			})
		}
		for _, delId := range q.idlist {
			if _, ok := idMap[delId]; !ok {
				result.Append(nil, &somaHostDeploymentResult{
					Delete:   true,
					DeleteId: delId,
				})
			}
		}
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Ops Access
 */
func (hd *somaHostDeploymentHandler) shutdownNow() {
	hd.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
