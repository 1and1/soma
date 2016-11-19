package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

type somaUserRequest struct {
	action string
	User   proto.User
	reply  chan somaResult
}

type somaUserResult struct {
	ResultError error
	User        proto.User
}

func (a *somaUserResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Users = append(r.Users, somaUserResult{ResultError: err})
	}
}

func (a *somaUserResult) SomaAppendResult(r *somaResult) {
	r.Users = append(r.Users, *a)
}

/* Read Access
 */
type somaUserReadHandler struct {
	input     chan somaUserRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	sync_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaUserReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ListUsers: r.list_stmt,
		stmt.ShowUsers: r.show_stmt,
		stmt.SyncUsers: r.sync_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`user`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

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

func (r *somaUserReadHandler) process(q *somaUserRequest) {
	var (
		userId, userName, firstName, lastName, mailAddr, team string
		employeeNr                                            int
		isActive, isSystem, isDeleted                         bool
		rows                                                  *sql.Rows
		err                                                   error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		r.appLog.Printf("R: users/list")
		rows, err = r.list_stmt.Query()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(
				&userId,
				&userName,
			)
			result.Append(err, &somaUserResult{
				User: proto.User{
					Id:       userId,
					UserName: userName,
				},
			})
		}
		if err = rows.Err(); err != nil {
			result.Append(err, &somaUserResult{})
		}
	case `sync`:
		r.appLog.Printf(`R: users/sync`)
		rows, err = r.sync_stmt.Query()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(
				&userId,
				&userName,
				&firstName,
				&lastName,
				&employeeNr,
				&mailAddr,
				&isDeleted,
				&team,
			)

			result.Append(err, &somaUserResult{
				User: proto.User{
					Id:             userId,
					UserName:       userName,
					FirstName:      firstName,
					LastName:       lastName,
					EmployeeNumber: strconv.Itoa(employeeNr),
					MailAddress:    mailAddr,
					IsDeleted:      isDeleted,
					TeamId:         team,
				},
			})
		}
		if err = rows.Err(); err != nil {
			result.Append(err, &somaUserResult{})
		}
	case "show":
		r.appLog.Printf("R: users/show for %s", q.User.Id)
		err = r.show_stmt.QueryRow(q.User.Id).Scan(
			&userId,
			&userName,
			&firstName,
			&lastName,
			&employeeNr,
			&mailAddr,
			&isActive,
			&isSystem,
			&isDeleted,
			&team,
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

		result.Append(err, &somaUserResult{
			User: proto.User{
				Id:             userId,
				UserName:       userName,
				FirstName:      firstName,
				LastName:       lastName,
				EmployeeNumber: strconv.Itoa(employeeNr),
				MailAddress:    mailAddr,
				IsActive:       isActive,
				IsSystem:       isSystem,
				IsDeleted:      isDeleted,
				TeamId:         team,
			},
		})
	default:
		r.errLog.Printf("R: unimplemented users/%s", q.action)
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaUserWriteHandler struct {
	input    chan somaUserRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	prg_stmt *sql.Stmt
	upd_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaUserWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.UserAdd:    w.add_stmt,
		stmt.UserUpdate: w.upd_stmt,
		stmt.UserDel:    w.del_stmt,
		stmt.UserPurge:  w.prg_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`user`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-w.shutdown:
			break runloop
		case req := <-w.input:
			go func() {
				w.process(&req)
			}()
		}
	}
}

func (w *somaUserWriteHandler) process(q *somaUserRequest) {
	var (
		res    sql.Result
		err    error
		super  *supervisor
		notify msg.Request
	)
	result := somaResult{}
	super = handlerMap[`supervisor`].(*supervisor)
	notify = msg.Request{Section: `map`, Action: `update`,
		Super: &msg.Supervisor{
			Object: `user`,
			User:   q.User,
		},
	}

	switch q.action {
	case "add":
		w.appLog.Printf("R: users/add for %s", q.User.UserName)
		id := uuid.NewV4()
		res, err = w.add_stmt.Exec(
			id.String(),
			q.User.UserName,
			q.User.FirstName,
			q.User.LastName,
			q.User.EmployeeNumber,
			q.User.MailAddress,
			false,
			q.User.IsSystem,
			false,
			q.User.TeamId,
		)
		q.User.Id = id.String()
		notify.Action = `add`
	case `update`:
		w.appLog.Printf("R: users/update for %s", q.User.Id)
		res, err = w.upd_stmt.Exec(
			q.User.UserName,
			q.User.FirstName,
			q.User.LastName,
			q.User.EmployeeNumber,
			q.User.MailAddress,
			q.User.IsDeleted,
			q.User.TeamId,
			q.User.Id,
		)
		notify.Action = `update`
	case "delete":
		w.appLog.Printf("R: users/delete for %s", q.User.Id)
		res, err = w.del_stmt.Exec(
			q.User.Id,
		)
		notify.Action = `delete`
	case "purge":
		w.appLog.Printf("R: user/purge for %s", q.User.Id)
		res, err = w.prg_stmt.Exec(
			q.User.Id,
		)
	default:
		w.errLog.Printf("R: unimplemented users/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaUserResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaUserResult{})
	default:
		result.Append(nil, &somaUserResult{
			User: q.User,
		})
		// send update to supervisor
		super.input <- notify
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaUserReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaUserWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
