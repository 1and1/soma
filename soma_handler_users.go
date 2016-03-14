package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/satori/go.uuid"

)

type somaUserRequest struct {
	action string
	User   somaproto.ProtoUser
	reply  chan somaResult
}

type somaUserResult struct {
	ResultError error
	User        somaproto.ProtoUser
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
}

func (r *somaUserReadHandler) run() {
	var err error

	log.Println("Prepare: user/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT user_id,
       user_uid
FROM   inventory.users;`)
	if err != nil {
		log.Fatal("user/list: ", err)
	}
	defer r.list_stmt.Close()

	log.Println("Prepare: user/show")
	r.show_stmt, err = r.conn.Prepare(`
SELECT user_id,
       user_uid,
       user_first_name,
	   user_last_name,
	   user_employee_number,
	   user_mail_address,
	   user_is_active,
	   user_is_system,
	   user_is_deleted,
	   organizational_team_id
FROM   inventory.users
WHERE  user_id = $1::uuid;`)
	if err != nil {
		log.Fatal("user/show: ", err)
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
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(
				&userId,
				&userName,
			)
			result.Append(err, &somaUserResult{
				User: somaproto.ProtoUser{
					Id:       userId,
					UserName: userName,
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
			if err.Error() != "sql: no rows in result set" {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		result.Append(err, &somaUserResult{
			User: somaproto.ProtoUser{
				Id:             userId,
				UserName:       userName,
				FirstName:      firstName,
				LastName:       lastName,
				EmployeeNumber: strconv.Itoa(employeeNr),
				MailAddress:    mailAddr,
				IsActive:       isActive,
				IsSystem:       isSystem,
				IsDeleted:      isDeleted,
				Team:           team,
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
}

func (w *somaUserWriteHandler) run() {
	var err error

	log.Println("Prepare: user/add")
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

	log.Println("Prepare: user/delete")
	w.del_stmt, err = w.conn.Prepare(`
UPDATE inventory.users
SET    user_is_deleted = 'yes',
       user_is_active = 'no'
WHERE  user_id = $1::uuid;`)
	if err != nil {
		log.Fatal("user/delete: ", err)
	}
	defer w.del_stmt.Close()

	log.Println("Prepare: user/purge")
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
			w.process(&req)
		}
	}
}

func (w *somaUserWriteHandler) process(q *somaUserRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

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
			q.User.Team,
		)
		q.User.Id = id.String()
	case "delete":
		log.Printf("R: users/delete for %s", q.User.Id)
		res, err = w.del_stmt.Exec(
			q.User.Id,
		)
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
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
