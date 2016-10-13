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

func TestTreeKeeperRunPreparedStmt(t *testing.T) {
	var err error
	tk, mock := testSpawnTreeKeeper(t)
	mock.ExpectPrepare(tkStmtStartJob)
	mock.ExpectPrepare(tkStmtGetViewFromCapability)
	mock.ExpectClose()

	tk.start_job, err = tk.conn.Prepare(tkStmtStartJob)
	if err != nil {
		t.Errorf("Error '%s' preparing statement", err)
	}
	if tk.start_job == nil {
		t.Errorf("stmt was expected preparing statement")
	}
	err = nil
	tk.get_view, err = tk.conn.Prepare(tkStmtGetViewFromCapability)
	if err != nil {
		t.Errorf("Error '%s' preparing statement", err)
	}
	if tk.get_view == nil {
		t.Errorf("stmt was expected preparing statement")
	}

	tk.conn.Close()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestTreeKeeperIsReady(t *testing.T) {
	tk, _ := testSpawnTreeKeeper(t)
	defer tk.conn.Close()
	tk.ready = false
	if tk.isReady() {
		t.Errorf("TK should not be ready")
	}
	tk.ready = true
	if !tk.isReady() {
		t.Errorf("TK should be reporting ready")
	}
}

func TestTreeKeeperIsBroken(t *testing.T) {
	tk, _ := testSpawnTreeKeeper(t)
	tk.conn.Close()
	tk.broken = false
	if tk.isBroken() {
		t.Errorf("TK should not be broken")
	}
	tk.broken = true
	if !tk.isBroken() {
		t.Errorf("TK should be reporting broken")
	}
}

func TestTreeKeeperIsStopped(t *testing.T) {
	tk, _ := testSpawnTreeKeeper(t)
	tk.conn.Close()
	tk.stopped = false
	if tk.isStopped() {
		t.Errorf("TK should not be stopped")
	}
	tk.stopped = true
	if !tk.isStopped() {
		t.Errorf("TK should be reporting being stopped")
	}
}

func TestTreeKeeperStop(t *testing.T) {
	tk, _ := testSpawnTreeKeeper(t)
	tk.conn.Close()
	tk.stopped = false
	tk.ready = true
	tk.broken = true
	tk.stop()
	if !tk.isStopped() {
		t.Errorf("TK should be stopped")
	}
	if tk.isReady() {
		t.Errorf("TK should not be reporting ready")
	}
	if tk.isBroken() {
		t.Errorf("TK should not be reporting broken")
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
