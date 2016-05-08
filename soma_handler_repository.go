package main

import (
	"database/sql"
	"log"

)

type somaRepositoryRequest struct {
	action     string
	Repository proto.Repository
	reply      chan somaResult
}

type somaRepositoryResult struct {
	ResultError error
	Repository  proto.Repository
}

func (a *somaRepositoryResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Repositories = append(r.Repositories, somaRepositoryResult{ResultError: err})
	}
}

func (a *somaRepositoryResult) SomaAppendResult(r *somaResult) {
	r.Repositories = append(r.Repositories, *a)
}

/* Read Access
 */
type somaRepositoryReadHandler struct {
	input     chan somaRepositoryRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
}

func (r *somaRepositoryReadHandler) run() {
	var err error

	log.Println("Prepare: repository/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT repository_id,
       repository_name
FROM   soma.repositories;`)
	if err != nil {
		log.Fatal("repository/list: ", err)
	}
	defer r.list_stmt.Close()

	log.Println("Prepare: repository/show")
	r.show_stmt, err = r.conn.Prepare(`
SELECT repository_id,
       repository_name,
       repository_active,
       organizational_team_id
FROM   soma.repositories
WHERE  repository_id = $1
AND    NOT repository_deleted;`)
	if err != nil {
		log.Fatal("repository/show: ", err)
	}
	defer r.show_stmt.Close()

runloop:
	for {
		select {
		case <-r.shutdown:
			break runloop
		case req := <-r.input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

func (r *somaRepositoryReadHandler) process(q *somaRepositoryRequest) {
	var (
		repoId, repoName, teamId string
		rows                     *sql.Rows
		repoActive               bool
		err                      error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: repository/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&repoId, &repoName)
			result.Append(err, &somaRepositoryResult{
				Repository: proto.Repository{
					Id:   repoId,
					Name: repoName,
				},
			})
		}
	case "show":
		log.Printf("R: repository/show for %s", q.Repository.Id)
		err = r.show_stmt.QueryRow(q.Repository.Id).Scan(
			&repoId,
			&repoName,
			&repoActive,
			&teamId,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		result.Append(err, &somaRepositoryResult{
			Repository: proto.Repository{
				Id:        repoId,
				Name:      repoName,
				TeamId:    teamId,
				IsDeleted: false,
				IsActive:  repoActive,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
