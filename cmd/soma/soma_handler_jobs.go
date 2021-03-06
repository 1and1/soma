package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type jobsRead struct {
	input         chan msg.Request
	shutdown      chan bool
	conn          *sql.DB
	listall_stmt  *sql.Stmt
	listscp_stmt  *sql.Stmt
	showid_stmt   *sql.Stmt
	showlist_stmt *sql.Stmt
	appLog        *log.Logger
	reqLog        *log.Logger
	errLog        *log.Logger
}

func (j *jobsRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ListAllOutstandingJobs:    j.listall_stmt,
		stmt.ListScopedOutstandingJobs: j.listscp_stmt,
		stmt.JobResultForId:            j.showid_stmt,
		stmt.JobResultsForList:         j.showlist_stmt,
	} {
		if prepStmt, err = j.conn.Prepare(statement); err != nil {
			j.errLog.Fatal(`jobs`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-j.shutdown:
			break runloop
		case req := <-j.input:
			go func() {
				j.process(&req)
			}()
		}
	}
}

func (j *jobsRead) process(q *msg.Request) {
	result := msg.Result{Type: q.Type, Action: q.Action, Job: []proto.Job{}}
	var (
		rows                                                               *sql.Rows
		err                                                                error
		jobId, jobType, jobStatus, jobResult, repositoryId, userId, teamId string
		jobError, jobSpec, idList                                          string
		jobSerial                                                          int
		jobQueued                                                          time.Time
		jobStarted, jobFinished                                            pq.NullTime
	)

	switch q.Action {
	case `list`:
		j.reqLog.Printf(LogStrReq, q.Type, q.Action, q.User, q.RemoteAddr)
		if q.IsAdmin {
			rows, err = j.listall_stmt.Query()
		} else {
			rows, err = j.listscp_stmt.Query(q.User)
		}
		if err != nil {
			result.ServerError(err)
			goto dispatch
		}
		defer rows.Close()

		for rows.Next() {
			if err = rows.Scan(
				&jobId,
				&jobType,
			); err != nil {
				result.ServerError(err)
				result.Clear(q.Type)
				goto dispatch
			}
			result.Job = append(result.Job,
				proto.Job{Id: jobId, Type: jobType})
		}
		if rows.Err() != nil {
			result.ServerError(err)
			result.Clear(q.Type)
			goto dispatch
		}
		result.OK()
	case `show`:
		j.reqLog.Printf(LogStrArg, q.Type, q.Action, q.User, q.RemoteAddr, q.Job.Id)
		if err = j.showid_stmt.QueryRow(q.Job.Id).Scan(
			&jobId,
			&jobStatus,
			&jobResult,
			&jobType,
			&jobSerial,
			&repositoryId,
			&userId,
			&teamId,
			&jobQueued,
			&jobStarted,
			&jobFinished,
			&jobError,
			&jobSpec,
		); err == sql.ErrNoRows {
			result.NotFound(err)
			goto dispatch
		} else if err != nil {
			result.ServerError(err)
			goto dispatch
		}
		job := proto.Job{
			Id:           jobId,
			Status:       jobStatus,
			Result:       jobResult,
			Type:         jobType,
			Serial:       jobSerial,
			RepositoryId: repositoryId,
			UserId:       userId,
			TeamId:       teamId,
			Error:        jobError,
		}
		job.TsQueued = jobQueued.Format(rfc3339Milli)
		if jobStarted.Valid {
			job.TsStarted = jobStarted.Time.Format(rfc3339Milli)
		}
		if jobFinished.Valid {
			job.TsFinished = jobFinished.Time.Format(rfc3339Milli)
		}
		if q.IsAdmin {
			job.Details = &proto.JobDetails{
				Specification: jobSpec,
			}
		}
		result.Job = []proto.Job{job}
		result.OK()
	case `search/idlist`:
		idList = fmt.Sprintf("{%s}", strings.Join(q.Search.Job.IdList, `,`))
		j.reqLog.Printf(LogStrArg, q.Type, q.Action, q.User, q.RemoteAddr, idList)
		if rows, err = j.showlist_stmt.Query(idList); err != nil {
			result.ServerError(err)
			goto dispatch
		}
		defer rows.Close()

		for rows.Next() {
			if err = rows.Scan(
				&jobId,
				&jobStatus,
				&jobResult,
				&jobType,
				&jobSerial,
				&repositoryId,
				&userId,
				&teamId,
				&jobQueued,
				&jobStarted,
				&jobFinished,
				&jobError,
				&jobSpec,
			); err != nil {
				result.ServerError(err)
				result.Clear(q.Type)
				goto dispatch
			}
			job := proto.Job{
				Id:           jobId,
				Status:       jobStatus,
				Result:       jobResult,
				Type:         jobType,
				Serial:       jobSerial,
				RepositoryId: repositoryId,
				UserId:       userId,
				TeamId:       teamId,
				Error:        jobError,
			}
			job.TsQueued = jobQueued.Format(rfc3339Milli)
			if jobStarted.Valid {
				job.TsStarted = jobStarted.Time.Format(rfc3339Milli)
			}
			if jobFinished.Valid {
				job.TsFinished = jobFinished.Time.Format(rfc3339Milli)
			}
			if q.IsAdmin && q.Search.IsDetailed {
				job.Details = &proto.JobDetails{
					Specification: jobSpec,
				}
			}
			result.Job = append(result.Job, job)
		}
		if err = rows.Err(); err != nil {
			result.ServerError(err)
			result.Clear(q.Type)
			goto dispatch
		}
		result.OK()
	default:
		result.NotImplemented(fmt.Errorf("Unknown requested action: %s/%s", q.Type, q.Action))
	}

dispatch:
	q.Reply <- result
}

/* Ops Access
 */
func (j *jobsRead) shutdownNow() {
	j.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
