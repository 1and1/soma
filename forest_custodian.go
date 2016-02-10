package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/satori/go.uuid"
)

type forestCustodian struct {
	input     chan somaRepositoryRequest
	shutdown  chan bool
	conn      *sql.DB
	add_stmt  *sql.Stmt
	load_stmt *sql.Stmt
}

func (f *forestCustodian) run() {
	var err error

	log.Println("Prepare: repository/create")
	f.add_stmt, err = f.conn.Prepare(`
INSERT INTO soma.repositories (
	repository_id,
	repository_name,
	repository_active,
	repository_deleted,
	organizational_team_id)
SELECT $1::uuid, $2::varchar, $3::boolean, $4::boolean, $5::uuid
WHERE NOT EXISTS (
	SELECT repository_id
	FROM   soma.repositories
	WHERE  repository_id = $1::uuid
	OR     repository_name = $2::varchar);`)
	if err != nil {
		log.Fatal("repository/add: ", err)
	}
	defer f.add_stmt.Close()

	log.Println("Prepare: repository/load")
	f.load_stmt, err = f.conn.Prepare(`
SELECT repository_id,
	   repository_name,
	   repository_deleted,
	   repository_active,
	   organizational_team_id
FROM   soma.repositories;`)
	if err != nil {
		log.Fatal("repository/load: ", err)
	}
	defer f.load_stmt.Close()

	f.initialLoad()

runloop:
	for {
		select {
		case <-f.shutdown:
			break runloop
		case req := <-f.input:
			f.process(&req)
		}
	}
}

func (f *forestCustodian) process(q *somaRepositoryRequest) {
	var (
		res        sql.Result
		err        error
		sTree      *somatree.SomaTree
		actionChan chan *somatree.Action
		errChan    chan *somatree.Error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: repository/add for %s", q.Repository.Name)
		actionChan = make(chan *somatree.Action, 1024000)
		errChan = make(chan *somatree.Error, 1024000)

		id := uuid.NewV4()
		q.Repository.Id = id.String()

		sTree = somatree.New(somatree.TreeSpec{
			Id:     uuid.NewV4().String(),
			Name:   fmt.Sprintf("root_%s", q.Repository.Name),
			Action: actionChan,
		})
		somatree.NewRepository(somatree.RepositorySpec{
			Id:      q.Repository.Id,
			Name:    q.Repository.Name,
			Team:    q.Repository.Team,
			Deleted: false,
			Active:  q.Repository.IsActive,
		}).Attach(somatree.AttachRequest{
			Root:       sTree,
			ParentType: "root",
			ParentId:   sTree.GetID(),
			ChildType:  "repository",
			ChildName:  q.Repository.Name,
		})
		sTree.SetError(errChan)

		for i := 0; i < len(actionChan); i++ {
			action := <-actionChan
			switch action.Action {
			case "create":
				if action.Type == "fault" {
					continue
				}
				if action.Type == "repository" {
					res, err = f.add_stmt.Exec(
						action.Repository.Id,
						action.Repository.Name,
						action.Repository.IsActive,
						false,
						action.Repository.Team,
					)
				}
			case "attached":
			}
		}
	default:
		log.Printf("R: unimplemented repository/%s", q.action)
		result.SetNotImplemented()
		q.reply <- result
		return
	}
	if result.SetRequestError(err) {
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		result.Append(errors.New("No rows affected"), &somaRepositoryResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaRepositoryResult{})
	default:
		result.Append(nil, &somaRepositoryResult{
			Repository: q.Repository,
		})
		if q.action == "add" {
			f.spawnTreeKeeper(q, sTree, errChan, actionChan)
		}
	}
	q.reply <- result
}

func (f *forestCustodian) initialLoad() {
	var (
		rows                     *sql.Rows
		err                      error
		repoId, repoName, teamId string
		repoActive, repoDeleted  bool
	)
	log.Printf("Loading existing repositories")
	rows, err = f.load_stmt.Query()
	if err != nil {
		log.Fatal("Error loading repositories: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&repoId,
			&repoName,
			&repoDeleted,
			&repoActive,
			&teamId,
		)
		if err != nil {
			log.Printf("Error: %s", err.Error())
		}
		f.loadSomaTree(&somaRepositoryRequest{
			Repository: somaproto.ProtoRepository{
				Id:        repoId,
				Name:      repoName,
				Team:      teamId,
				IsDeleted: repoDeleted,
				IsActive:  repoActive,
			},
		})
	}
}

func (f *forestCustodian) loadSomaTree(q *somaRepositoryRequest) {
	actionChan := make(chan *somatree.Action, 1024000)
	errChan := make(chan *somatree.Error, 1024000)

	sTree := somatree.New(somatree.TreeSpec{
		Id:     uuid.NewV4().String(),
		Name:   fmt.Sprintf("root_%s", q.Repository.Name),
		Action: actionChan,
	})
	somatree.NewRepository(somatree.RepositorySpec{
		Id:      q.Repository.Id,
		Name:    q.Repository.Name,
		Team:    q.Repository.Team,
		Deleted: q.Repository.IsDeleted,
		Active:  q.Repository.IsActive,
	}).Attach(somatree.AttachRequest{
		Root:       sTree,
		ParentType: "root",
		ParentId:   sTree.GetID(),
		ChildType:  "repository",
		ChildName:  q.Repository.Name,
	})
	sTree.SetError(errChan)
	for i := 0; i < len(actionChan); i++ {
		// discard actions on initial load
		<-actionChan
	}
	f.spawnTreeKeeper(q, sTree, errChan, actionChan)
}

func (f *forestCustodian) spawnTreeKeeper(q *somaRepositoryRequest, s *somatree.SomaTree,
	ec chan *somatree.Error, ac chan *somatree.Action) {
	var tK treeKeeper
	tK.input = make(chan treeRequest, 1024)
	tK.shutdown = make(chan bool)
	tK.conn = f.conn
	tK.tree = s
	tK.errChan = ec
	tK.actionChan = ac
	tK.repoId = q.Repository.Id
	tK.repoName = q.Repository.Name
	keeperName := fmt.Sprintf("repository_%s", q.Repository.Name)
	handlerMap[keeperName] = tK
	go tK.run()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
