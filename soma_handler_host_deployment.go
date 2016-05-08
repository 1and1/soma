package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"

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
}

func (self *somaHostDeploymentHandler) run() {
	var err error

	log.Println("Prepare: hostdeployment/get-for-node")
	if self.geti_stmt, err = self.conn.Prepare(stmtGetInstancesForNode); err != nil {
		log.Fatal("hostdeployment/get-for-node: ", err)
	}
	defer self.geti_stmt.Close()

	log.Println("Prepare: hostdeployment/last-version")
	if self.last_stmt, err = self.conn.Prepare(stmtGetLastInstanceVersion); err != nil {
		log.Fatal("hostdeployment/last-version: ", err)
	}
	defer self.last_stmt.Close()

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
		checkInstanceID, deploymentDetails string
		idList                             *sql.Rows
		err                                error
	)
	result := somaResult{}

	switch q.action {
	case "get":
		log.Printf("R: hostdeployment/get-for-node for %s", q.assetid)
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
			err = self.last_stmt.QueryRow(checkInstanceID).Scan(&deploymentDetails)
			if err != nil {
				result.Append(err, &somaHostDeploymentResult{})
				continue idgetloop
			}
			depl := proto.Deployment{}
			if err = json.Unmarshal([]byte(deploymentDetails), depl); err != nil {
				result.Append(err, &somaHostDeploymentResult{})
				continue idgetloop
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
		log.Printf("R: hostdeployment/get-for-node for %s", q.assetid)
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
			err = self.last_stmt.QueryRow(checkInstanceID).Scan(&deploymentDetails)
			if err != nil {
				result.Append(err, &somaHostDeploymentResult{})
				continue assembleloop
			}
			depl := proto.Deployment{}
			if err = json.Unmarshal([]byte(deploymentDetails), depl); err != nil {
				result.Append(err, &somaHostDeploymentResult{})
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
