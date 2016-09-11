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
	system    chan msg.Request
	shutdown  chan bool
	conn      *sql.DB
	add_stmt  *sql.Stmt
	load_stmt *sql.Stmt
	name_stmt *sql.Stmt
	rbck_stmt *sql.Stmt
	rbci_stmt *sql.Stmt
}

func (f *forestCustodian) run() {
	var err error

	if f.add_stmt, err = f.conn.Prepare(
		stmt.ForestAddRepository,
	); err != nil {
		log.Fatal("repository/add: ", err)
	}
	defer f.add_stmt.Close()

	if f.load_stmt, err = f.conn.Prepare(
		stmt.ForestLoadRepository,
	); err != nil {
		log.Fatal("repository/load: ", err)
	}
	defer f.load_stmt.Close()

	if f.name_stmt, err = f.conn.Prepare(
		stmt.ForestRepoNameById,
	); err != nil {
		log.Fatal("forestCustodian/reponame-by-id: ", err)
	}
	defer f.name_stmt.Close()

	if f.rbck_stmt, err = f.conn.Prepare(
		stmt.ForestRebuildDeleteChecks,
	); err != nil {
		log.Fatal("forestCustodian/delete-checks-for-repo: ", err)
	}
	defer f.rbck_stmt.Close()

	if f.rbci_stmt, err = f.conn.Prepare(
		stmt.ForestRebuildDeleteInstances,
	); err != nil {
		log.Fatal("forestCustodian/delete-check-instances-for-repo: ",
			err)
	}
	defer f.rbci_stmt.Close()

	f.initialLoad()

	if SomaCfg.Observer {
		log.Println(`ForestCustodian entered observer mode`)
		<-f.shutdown
		goto exit
	}

runloop:
	for {
		select {
		case <-f.shutdown:
			break runloop
		case req := <-f.input:
			f.process(&req)
		case req := <-f.system:
			f.sysprocess(&req)
		}
	}
exit:
}

func (f *forestCustodian) process(q *somaRepositoryRequest) {
	var (
		res        sql.Result
		err        error
		sTree      *tree.Tree
		actionChan chan *tree.Action
		errChan    chan *tree.Error
		team       string
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf(LogStrReq,
			`ForestCustodian`,
			fmt.Sprintf("%s/%s", `CreateRepository`, q.Repository.Name),
			q.user, q.remoteAddr,
		)
		actionChan = make(chan *tree.Action, 1024000)
		errChan = make(chan *tree.Error, 1024000)

		id := uuid.NewV4()
		q.Repository.Id = id.String()

		sTree = tree.New(tree.TreeSpec{
			Id:     uuid.NewV4().String(),
			Name:   fmt.Sprintf("root_%s", q.Repository.Name),
			Action: actionChan,
		})
		tree.NewRepository(tree.RepositorySpec{
			Id:      q.Repository.Id,
			Name:    q.Repository.Name,
			Team:    q.Repository.TeamId,
			Deleted: false,
			Active:  q.Repository.IsActive,
		}).Attach(tree.AttachRequest{
			Root:       sTree,
			ParentType: "root",
			ParentId:   sTree.GetID(),
		})
		sTree.SetError(errChan)

		for i := len(actionChan); i > 0; i-- {
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
						action.Repository.TeamId,
						q.user,
					)
					team = action.Repository.TeamId
				}
			case "attached":
			default:
				log.Printf("R: Unhandled action during tree creation: %s", q.action)
				result.SetNotImplemented()
				q.reply <- result
				return
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
		if team == "" {
			result.SetRequestError(
				fmt.Errorf("Team has not been set prior to spawning TreeKeeper for repo: %s", q.Repository.Name),
			)
			q.reply <- result
			return
		}
		result.Append(nil, &somaRepositoryResult{
			Repository: q.Repository,
		})
		if q.action == "add" {
			f.spawnTreeKeeper(q, sTree, errChan, actionChan, team)
		}
	}
	q.reply <- result
}

func (f *forestCustodian) sysprocess(q *msg.Request) {
	var (
		repoId, repoName, teamId, keeper string
		err                              error
	)
	result := msg.Result{
		Type:   `forestcustodian`,
		Action: `systemoperation`,
		System: []proto.SystemOperation{q.System},
	}

	switch q.System.Request {
	case `rebuild_repository`:
		repoId = q.System.RepositoryId
	default:
		result.NotImplemented(
			fmt.Errorf("Unknown requested system operation: %s",
				q.System.Request),
		)
		goto exit
	}

	// look up name of the repository
	if err = f.name_stmt.QueryRow(repoId).
		Scan(&repoName, &teamId); err != nil {
		if err == sql.ErrNoRows {
			result.NotFound(fmt.Errorf(`No such repository`))
		} else {
			result.ServerError(err)
		}
		goto exit
	}

	// get the treekeeper for the repository
	keeper = fmt.Sprintf("repository_%s", repoName)
	if handler, ok := handlerMap[keeper].(*treeKeeper); ok {
		// stop the handler
		handler.stopchan <- true

		// remove handler from lookup table
		delete(handlerMap, keeper)
	}

	// mark all existing check instances as deleted - instances
	// are deleted for both rebuild levels checks and instances
	if _, err = f.rbci_stmt.Exec(repoId); err != nil {
		result.ServerError(err)
		goto exit
	}
	// only delete checks for rebuild level checks
	if q.System.RebuildLevel == `checks` {
		if _, err = f.rbck_stmt.Exec(repoId); err != nil {
			result.ServerError(err)
			goto exit
		}
	}

	// load the tree again, with requested rebuild active
	f.loadSomaTree(&somaRepositoryRequest{
		rebuild: true,
		rbLevel: q.System.RebuildLevel,
		Repository: proto.Repository{
			Id:        repoId,
			Name:      repoName,
			TeamId:    teamId,
			IsDeleted: false,
			IsActive:  true,
		},
	})

	// rebuild has finished, restart the tree. If the rebuild did not
	// work, this will simply be a broken tree once more
	f.loadSomaTree(&somaRepositoryRequest{
		rebuild: false,
		rbLevel: "",
		Repository: proto.Repository{
			Id:        repoId,
			Name:      repoName,
			TeamId:    teamId,
			IsDeleted: false,
			IsActive:  true,
		},
	})
	result.OK()

exit:
	q.Reply <- result
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
			Repository: proto.Repository{
				Id:        repoId,
				Name:      repoName,
				TeamId:    teamId,
				IsDeleted: repoDeleted,
				IsActive:  repoActive,
			},
		})
	}
}

func (f *forestCustodian) loadSomaTree(q *somaRepositoryRequest) {
	actionChan := make(chan *tree.Action, 1024000)
	errChan := make(chan *tree.Error, 1024000)

	sTree := tree.New(tree.TreeSpec{
		Id:     uuid.NewV4().String(),
		Name:   fmt.Sprintf("root_%s", q.Repository.Name),
		Action: actionChan,
	})
	tree.NewRepository(tree.RepositorySpec{
		Id:      q.Repository.Id,
		Name:    q.Repository.Name,
		Team:    q.Repository.TeamId,
		Deleted: q.Repository.IsDeleted,
		Active:  q.Repository.IsActive,
	}).Attach(tree.AttachRequest{
		Root:       sTree,
		ParentType: "root",
		ParentId:   sTree.GetID(),
	})
	sTree.SetError(errChan)
	for i := len(actionChan); i > 0; i-- {
		// discard actions on initial load
		<-actionChan
	}
	for i := len(errChan); i > 0; i-- {
		// discard actions on initial load
		<-errChan
	}
	f.spawnTreeKeeper(q, sTree, errChan, actionChan, q.Repository.TeamId)
}

func (f *forestCustodian) spawnTreeKeeper(q *somaRepositoryRequest, s *tree.Tree,
	ec chan *tree.Error, ac chan *tree.Action, team string) {

	// only start the single requested repo in observer mode with
	// set repo flag
	if SomaCfg.Observer && SomaCfg.ObserverRepo != `` && q.Repository.Name != SomaCfg.ObserverRepo {
		return
	}

	db, err := newDatabaseConnection()
	if err != nil {
		return
	}
	tK := new(treeKeeper)
	tK.input = make(chan treeRequest, 1024)
	tK.shutdown = make(chan bool)
	tK.stopchan = make(chan bool)
	tK.conn = db
	tK.tree = s
	tK.errChan = ec
	tK.actionChan = ac
	tK.broken = false
	tK.ready = false
	tK.frozen = false
	tK.stopped = false
	tK.rebuild = q.rebuild
	tK.rbLevel = q.rbLevel
	tK.repoId = q.Repository.Id
	tK.repoName = q.Repository.Name
	tK.team = team
	keeperName := fmt.Sprintf("repository_%s", q.Repository.Name)

	// during rebuild the treekeeper will not run in background
	if tK.rebuild {
		tK.run()
	} else {
		// non-rebuild, register TK and detach
		handlerMap[keeperName] = tK
		go tK.run()
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
