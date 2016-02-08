package main

type forestCustodian struct {
	input    chan somaRepositoryRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
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
	organizational_team_id,
SELECT $1::uuid, $2::varchar, $3::boolean, $4::boolean, $5::uuid
WHERE NOT EXISTS (
	SELECT repository_id
	FROM   soma.repositories
	WHERE  repository_id = $1::uuid
	OR     repository_name = $2::varchar;`)
	if err != nil {
		log.Fatal("repository/add: ", err)
	}
	defer f.add_stmt.Close()

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

func (f *forestCustodian) process() {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: repository/add for %s", q.Repository.Name)
		id := uuid.NewV4()
		res, err = w.add_stmt.Exec(
			id.String(),
			q.Repository.Name,
			q.Repository.Active,
			false,
			q.Repository.Team,
		)
		q.Repository.Id = id.String()
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
		result.Append(errors.New("No rows affected"), &somaNodeResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaNodeResult{})
	default:
		result.Append(nil, &somaNodeResult{
			Node: q.Node,
		})

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
			Deleted: false,
			Active:  q.Repository.Active,
		}).Attach(somatree.AttachRequest{
			Root:       sTree,
			ParentType: "root",
			ParentId:   sTree.GetID(),
			ChildType:  "repository",
			ChildName:  q.Repository.Name,
		})
		sTree.SetError(errChan)

		var treeKeeper somaTreeKeeper
		treeKeeper.input = make(chan somaTreeRequest, 1024)
		treeKeeper.shutdown = make(chan bool)
		treeKeeper.conn = conn
		treeKeeper.tree = Stree
		treeKeeper.errChan = errChan
		treeKeeper.actionChan = actionChan
		keeperName := fmt.Sprintf("repository_%s", q.Repository.Name)
		handlerMap[keeperName] = treeKeeper
		go treeKeeper.run()
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
