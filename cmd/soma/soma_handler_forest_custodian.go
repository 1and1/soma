package main

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/internal/tree"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
	"github.com/client9/reopen"
	uuid "github.com/satori/go.uuid"
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
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (f *forestCustodian) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ForestAddRepository:          f.add_stmt,
		stmt.ForestLoadRepository:         f.load_stmt,
		stmt.ForestRepoNameById:           f.name_stmt,
		stmt.ForestRebuildDeleteChecks:    f.rbck_stmt,
		stmt.ForestRebuildDeleteInstances: f.rbci_stmt,
	} {
		if prepStmt, err = f.conn.Prepare(statement); err != nil {
			f.errLog.Fatal(`forestcustodian`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

	f.initialLoad()

	if SomaCfg.Observer {
		// XXX restart repository should be possible in observer mode
		f.appLog.Println(`ForestCustodian entered observer mode`)
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
		f.reqLog.Printf(LogStrReq,
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
			Log:    f.appLog,
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
				f.errLog.Printf("R: Unhandled action during tree creation: %s", q.action)
				result.SetNotImplemented()
				q.reply <- result
				return
			}
		}
	default:
		f.errLog.Printf("R: unimplemented repository/%s", q.action)
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
			err = f.spawnTreeKeeper(q, sTree, errChan, actionChan, team)
			result.SetRequestError(err)
		}
	}
	q.reply <- result
}

func (f *forestCustodian) sysprocess(q *msg.Request) {
	var (
		repoId, repoName, teamId, keeper string
		err                              error
	)
	result := msg.FromRequest(q)
	result.System = []proto.SystemOperation{q.System}

	switch q.System.Request {
	case `repository_rebuild`, `repository_restart`:
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
		// remove handler from lookup table
		delete(handlerMap, keeper)

		// stop the handler before shut down to give it a chance to
		// drain the input channel
		if !handler.isStopped() {
			handler.stopchan <- true
		}
		handler.shutdown <- true
	}

	if q.System.Request == `repository_rebuild` {
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
		if err = f.loadSomaTree(&somaRepositoryRequest{
			rebuild: true,
			rbLevel: q.System.RebuildLevel,
			Repository: proto.Repository{
				Id:        repoId,
				Name:      repoName,
				TeamId:    teamId,
				IsDeleted: false,
				IsActive:  true,
			},
		}); err != nil {
			result.ServerError(err)
			goto exit
		}
	}

	// rebuild has finished, restart the tree. If the rebuild did not
	// work, this will simply be a broken tree once more
	if err = f.loadSomaTree(&somaRepositoryRequest{
		rebuild: false,
		rbLevel: "",
		Repository: proto.Repository{
			Id:        repoId,
			Name:      repoName,
			TeamId:    teamId,
			IsDeleted: false,
			IsActive:  true,
		},
	}); err != nil {
		result.ServerError(err)
		goto exit
	}
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
	f.appLog.Printf("Loading existing repositories")
	rows, err = f.load_stmt.Query()
	if err != nil {
		f.errLog.Fatal("Error loading repositories: ", err)
	}
	defer rows.Close()

treeloop:
	for rows.Next() {
		err = rows.Scan(
			&repoId,
			&repoName,
			&repoDeleted,
			&repoActive,
			&teamId,
		)
		if err != nil {
			f.errLog.Printf("Error: %s", err.Error())
			err = nil
			continue treeloop
		}
		err = f.loadSomaTree(&somaRepositoryRequest{
			Repository: proto.Repository{
				Id:        repoId,
				Name:      repoName,
				TeamId:    teamId,
				IsDeleted: repoDeleted,
				IsActive:  repoActive,
			},
		})
		if err != nil {
			f.errLog.Printf("fc.loadSomaTree(), error: %s", err.Error())
			err = nil
		}
	}
	if err = rows.Err(); err != nil {
		f.errLog.Printf("fc.initialLoad(), error: %s", err.Error())
	}
}

func (f *forestCustodian) loadSomaTree(q *somaRepositoryRequest) error {
	actionChan := make(chan *tree.Action, 1024000)
	errChan := make(chan *tree.Error, 1024000)

	sTree := tree.New(tree.TreeSpec{
		Id:     uuid.NewV4().String(),
		Name:   fmt.Sprintf("root_%s", q.Repository.Name),
		Action: actionChan,
		Log:    f.appLog,
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
	return f.spawnTreeKeeper(q, sTree, errChan, actionChan, q.Repository.TeamId)
}

func (f *forestCustodian) spawnTreeKeeper(q *somaRepositoryRequest, s *tree.Tree,
	ec chan *tree.Error, ac chan *tree.Action, team string) error {

	// only start the single requested repo
	if SomaCfg.ObserverRepo != `` && q.Repository.Name != SomaCfg.ObserverRepo {
		return nil
	}
	var (
		err      error
		db       *sql.DB
		lfh, sfh *reopen.FileWriter
	)

	if db, err = newDatabaseConnection(); err != nil {
		return err
	}

	keeperName := fmt.Sprintf("repository_%s", q.Repository.Name)
	if lfh, err = reopen.NewFileWriter(filepath.Join(
		SomaCfg.LogPath,
		`repository`,
		fmt.Sprintf("%s.log", keeperName),
	)); err != nil {
		return err
	}
	if sfh, err = reopen.NewFileWriter(filepath.Join(
		SomaCfg.LogPath,
		`repository`,
		fmt.Sprintf("startup_%s.log", keeperName),
	)); err != nil {
		return err
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
	tK.appLog = f.appLog
	tK.log = log.New()
	tK.log.Out = lfh
	tK.startLog = log.New()
	tK.startLog.Out = sfh
	// startup logs are not rotated, the logrotate map is just used
	// to keep acccess to the filehandle
	logFileMap[fmt.Sprintf("%s", keeperName)] = lfh
	logFileMap[fmt.Sprintf("startup_%s", keeperName)] = sfh

	// during rebuild the treekeeper will not run in background
	if tK.rebuild {
		tK.run()
	} else {
		// non-rebuild, register TK and detach
		handlerMap[keeperName] = tK
		go tK.run()
	}
	return nil
}

/* Ops Access
 */
func (f *forestCustodian) shutdownNow() {
	f.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
