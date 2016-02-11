package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	broken     bool
	ready      bool
	input      chan treeRequest
	shutdown   chan bool
	conn       *sql.DB
	tree       *somatree.SomaTree
	errChan    chan *somatree.Error
	actionChan chan *somatree.Action
}

func (tk *treeKeeper) run() {
	log.Printf("Starting TreeKeeper for Repo %s (%s)", tk.repoName, tk.repoId)
	//var err error

	tk.startupLoad()

	if tk.broken {
		tickTack := time.NewTicker(time.Second * 10).C
	hoverloop:
		for {
			select {
			case <-tickTack:
				log.Printf("TK[%s]: BROKEN REPOSITORY %s flying holding patterns!\n",
					tk.repoName, tk.repoId)
			case <-tk.shutdown:
				break hoverloop
			}
		}
		return
	}
	log.Printf("TK[%s]: ready for service!\n", tk.repoName)
	tk.ready = true

runloop:
	for {
		select {
		case <-tk.shutdown:
			break runloop
		case <-tk.input:
			fmt.Printf("TK %s received input request\n", tk.repoName)
		}
	}
}

func (tk *treeKeeper) isReady() bool {
	return tk.ready
}

func (tk *treeKeeper) isBroken() bool {
	return tk.broken
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
