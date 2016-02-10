package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/satori/go.uuid"

)

type treeRequest struct {
	RequestType string
	Action      string
	JobId       uuid.UUID
	reply       chan somaResult
	Repository  somaRepositoryRequest
	Bucket      somaBucketRequest
	Group       somaGroupRequest
	Cluster     somaClusterRequest
	Node        somaNodeRequest
}

type treeResult struct {
	ResultType  string
	ResultError error
	JobId       uuid.UUID
	Repository  somaRepositoryResult
	Bucket      somaRepositoryRequest
}

type treeKeeper struct {
	repoId     string
	repoName   string
	input      chan treeRequest
	shutdown   chan bool
	conn       *sql.DB
	tree       *somatree.SomaTree
	errChan    chan *somatree.Error
	actionChan chan *somatree.Action
}

func (tk *treeKeeper) run() {
	log.Printf("Starting TreeKeeper for Repo %s (%s)", tk.repoName, tk.repoId)

runloop:
	for {
		select {
		case <-tk.shutdown:
			break runloop
		case <-tk.input:
			fmt.Printf("TK %s received input request", tk.repoName)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
