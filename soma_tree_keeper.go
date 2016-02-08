package main

import (
	"database/sql"

	"github.com/satori/go.uuid"

)

type treeRequest struct {
	RequestType string
	Action      string
	Repository  somaRepositoryRequest
	Bucket      somaBucketRequest
}

type treeResult struct {
	ResultType  string
	ResultError error
	JobId       uuid.UUID
	Repository  somaRepositoryResult
	Bucket      somaRepositoryRequest
}

type treeKeeper struct {
	input      chan treeRequest
	shutdown   chan bool
	conn       *sql.DB
	tree       *somatree.SomaTree
	errChan    chan *somatree.Error
	actionChan chan *somatree.Action
}

func (tk *treeKeeper) run() {
}

type guidePost struct {
	input chan treeRequest
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
