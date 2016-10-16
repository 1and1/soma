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

	if r.list_stmt, err = r.conn.Prepare(stmt.ListUsers); err != nil {
		log.Fatal("user/list: ", err)
	}
	defer r.list_stmt.Close()

	if r.show_stmt, err = r.conn.Prepare(stmt.ShowUsers); err != nil {
		log.Fatal("user/show: ", err)
	}
	defer r.show_stmt.Close()

	if r.sync_stmt, err = r.conn.Prepare(stmt.SyncUsers); err != nil {
		log.Fatal("user/sync: ", err)
	}
	defer r.sync_stmt.Close()

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
		log.Printf("R: users/list")
		rows, err = r.list_stmt.Query()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(
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
			err = nil
		}
	case `sync`:
		log.Printf(`R: users/sync`)
		rows, err = r.sync_stmt.Query()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(
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
			err = nil
		}
	case "show":
		log.Printf("R: users/show for %s", q.User.Id)
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

	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO inventory.users (
	user_id,
	user_uid,
	user_first_name,
	user_last_name,
	user_employee_number,
	user_mail_address,
	user_is_active,
	user_is_system,
	user_is_deleted,
	organizational_team_id)
SELECT $1::uuid, $2::varchar, $3::varchar, $4::varchar, $5::numeric,
	   $6::text, $7::boolean, $8::boolean, $9::boolean, $10::uuid
WHERE NOT EXISTS (
	SELECT user_id
	FROM   inventory.users
	WHERE  user_id = $1::uuid
    OR     user_uid = $2::varchar
    OR     user_employee_number = $5::numeric);`)
	if err != nil {
		log.Fatal("user/add: ", err)
	}
	defer w.add_stmt.Close()

	if w.upd_stmt, err = w.conn.Prepare(`
UPDATE inventory.users
SET    user_uid = $1::varchar,
       user_first_name = $2::varchar,
       user_last_name = $3::varchar,
       user_employee_number = $4::numeric,
       user_mail_address = $5::text,
       user_is_deleted = $6::boolean,
       organizational_team_id = $7::uuid
WHERE  user_id = $8::uuid;`); err != nil {
		log.Fatal(`user/update: `, err)
	}
	defer w.upd_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
UPDATE inventory.users
SET    user_is_deleted = 'yes',
       user_is_active = 'no'
WHERE  user_id = $1::uuid;`)
	if err != nil {
		log.Fatal("user/delete: ", err)
	}
	defer w.del_stmt.Close()

	w.prg_stmt, err = w.conn.Prepare(`
DELETE FROM inventory.users
WHERE  user_id = $1::uuid
AND    user_is_deleted;`)
	if err != nil {
		log.Fatal("user/purge: ", err)
	}
	defer w.prg_stmt.Close()

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
	notify = msg.Request{Type: `supervisor`, Action: `update_map`,
		Super: &msg.Supervisor{
			Object: `user`,
			User:   q.User,
		},
	}

	switch q.action {
	case "add":
		log.Printf("R: users/add for %s", q.User.UserName)
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
		notify.Super.Action = `add`
	case `update`:
		log.Printf("R: users/update for %s", q.User.Id)
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
		notify.Super.Action = `update`
	case "delete":
		log.Printf("R: users/delete for %s", q.User.Id)
		res, err = w.del_stmt.Exec(
			q.User.Id,
		)
		notify.Super.Action = `delete`
	case "purge":
		log.Printf("R: user/purge for %s", q.User.Id)
		res, err = w.prg_stmt.Exec(
			q.User.Id,
		)
	default:
		log.Printf("R: unimplemented users/%s", q.action)
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
