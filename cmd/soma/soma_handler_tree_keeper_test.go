package main

import (
	"fmt"
	"testing"

	"github.com/1and1/soma/internal/tree"
	"github.com/satori/go.uuid"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestTreeKeeperDBClose(t *testing.T) {
	tk, mock := testSpawnTreeKeeper(t)
	mock.ExpectClose()
	tk.conn.Close()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func testSpawnTreeKeeper(t *testing.T) (*treeKeeper, sqlmock.Sqlmock) {
	actionChan := make(chan *tree.Action, 1024000)
	errChan := make(chan *tree.Error, 1024000)

	sTree := tree.New(tree.TreeSpec{
		Id:     uuid.NewV4().String(),
		Name:   fmt.Sprintf("root_%s", `repo_test`),
		Action: actionChan,
	})
	tree.NewRepository(tree.RepositorySpec{
		Id:      uuid.NewV4().String(),
		Name:    `repo_test`,
		Team:    uuid.NewV4().String(),
		Deleted: false,
		Active:  true,
	}).Attach(tree.AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   sTree.GetID(),
	})
	sTree.SetError(errChan)
	for i := len(actionChan); i > 0; i-- {
		<-actionChan
	}
	for i := len(errChan); i > 0; i-- {
		<-errChan
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error '%s' encountered while opening"+
			"stub database connection", err)
	}

	tK := new(treeKeeper)
	tK.input = make(chan treeRequest, 1024)
	tK.shutdown = make(chan bool)
	tK.stopchan = make(chan bool)
	tK.conn = db
	tK.tree = sTree
	tK.errChan = errChan
	tK.actionChan = actionChan
	tK.broken = false
	tK.ready = false
	tK.frozen = false
	tK.stopped = false
	tK.rebuild = false
	tK.rbLevel = ``
	tK.repoId = uuid.NewV4().String()
	tK.repoName = `repo_test`
	tK.team = uuid.NewV4().String()

	return tK, mock
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
