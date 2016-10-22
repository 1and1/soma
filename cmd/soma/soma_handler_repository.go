package main

import (
	"database/sql"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type somaRepositoryRequest struct {
	action     string
	remoteAddr string
	user       string
	rbLevel    string
	rebuild    bool
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
	ponc_stmt *sql.Stmt
	psvc_stmt *sql.Stmt
	psys_stmt *sql.Stmt
	pcst_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaRepositoryReadHandler) run() {
	var err error

	if r.list_stmt, err = r.conn.Prepare(stmt.ListAllRepositories); err != nil {
		r.errLog.Fatal("repository/list: ", err)
	}
	defer r.list_stmt.Close()

	if r.show_stmt, err = r.conn.Prepare(stmt.ShowRepository); err != nil {
		r.errLog.Fatal("repository/show: ", err)
	}
	defer r.show_stmt.Close()

	if r.ponc_stmt, err = r.conn.Prepare(stmt.RepoOncProps); err != nil {
		r.errLog.Fatal(`repository/property-oncall: `, err)
	}
	defer r.ponc_stmt.Close()

	if r.psvc_stmt, err = r.conn.Prepare(stmt.RepoSvcProps); err != nil {
		r.errLog.Fatal(`repository/property-service: `, err)
	}
	defer r.psvc_stmt.Close()

	if r.psys_stmt, err = r.conn.Prepare(stmt.RepoSysProps); err != nil {
		r.errLog.Fatal(`repository/property-system: `, err)
	}
	defer r.psys_stmt.Close()

	if r.pcst_stmt, err = r.conn.Prepare(stmt.RepoCstProps); err != nil {
		r.errLog.Fatal(`repository/property-custom: `, err)
	}
	defer r.pcst_stmt.Close()

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
		repoId, repoName, teamId, instanceId, sourceInstanceId string
		view, oncallId, oncallName, serviceName, customId      string
		systemProp, value, customProp                          string
		rows                                                   *sql.Rows
		repoActive                                             bool
		err                                                    error
	)
	result := somaResult{}

	switch q.action {
	case `list`:
		r.reqLog.Printf("R: repository/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			goto dispatch
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
	case `show`:
		r.reqLog.Printf("R: repository/show for %s", q.Repository.Id)
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
			goto dispatch
		}
		repo := proto.Repository{
			Id:        repoId,
			Name:      repoName,
			TeamId:    teamId,
			IsDeleted: false,
			IsActive:  repoActive,
		}
		repo.Properties = &[]proto.Property{}

		// oncall properties
		rows, err = r.ponc_stmt.Query(q.Repository.Id)
		if result.SetRequestError(err) {
			goto dispatch
		}
		for rows.Next() {
			if err := rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&oncallId,
				&oncallName,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*repo.Properties = append(
				*repo.Properties,
				proto.Property{
					Type:             `oncall`,
					RepositoryId:     q.Repository.Id,
					InstanceId:       instanceId,
					SourceInstanceId: sourceInstanceId,
					View:             view,
					Oncall: &proto.PropertyOncall{
						Id:   oncallId,
						Name: oncallName,
					},
				},
			)
		}
		if err = rows.Err(); err != nil {
			result.SetRequestError(err)
			goto dispatch
		}

		// service properties
		rows, err = r.psvc_stmt.Query(q.Repository.Id)
		for rows.Next() {
			if err := rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&serviceName,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*repo.Properties = append(
				*repo.Properties,
				proto.Property{
					Type:             `service`,
					RepositoryId:     q.Repository.Id,
					InstanceId:       instanceId,
					SourceInstanceId: sourceInstanceId,
					View:             view,
					Service: &proto.PropertyService{
						Name: serviceName,
					},
				},
			)
		}
		if err = rows.Err(); err != nil {
			result.SetRequestError(err)
			goto dispatch
		}

		// system properties
		rows, err = r.psys_stmt.Query(q.Repository.Id)
		for rows.Next() {
			if err := rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&systemProp,
				&value,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*repo.Properties = append(
				*repo.Properties,
				proto.Property{
					Type:             `system`,
					RepositoryId:     q.Repository.Id,
					InstanceId:       instanceId,
					SourceInstanceId: sourceInstanceId,
					View:             view,
					System: &proto.PropertySystem{
						Name:  systemProp,
						Value: value,
					},
				},
			)
		}
		if err = rows.Err(); err != nil {
			result.SetRequestError(err)
			goto dispatch
		}

		// custom properties
		rows, err = r.pcst_stmt.Query(q.Repository.Id)
		for rows.Next() {
			if err := rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&customId,
				&value,
				&customProp,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*repo.Properties = append(
				*repo.Properties,
				proto.Property{
					Type:             `custom`,
					RepositoryId:     q.Repository.Id,
					InstanceId:       instanceId,
					SourceInstanceId: sourceInstanceId,
					View:             view,
					Custom: &proto.PropertyCustom{
						Id:    customId,
						Name:  customProp,
						Value: value,
					},
				},
			)
		}
		if err = rows.Err(); err != nil {
			result.SetRequestError(err)
			goto dispatch
		}

		result.Append(err, &somaRepositoryResult{
			Repository: repo,
		})
	default:
		result.SetNotImplemented()
	}

dispatch:
	q.reply <- result
}

/* Ops Access
 */
func (r *somaRepositoryReadHandler) shutdownNow() {
	r.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
