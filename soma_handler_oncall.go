package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/satori/go.uuid"
)

type somaOncallRequest struct {
	action string
	Oncall proto.Oncall
	reply  chan somaResult
}

type somaOncallResult struct {
	ResultError error
	Oncall      proto.Oncall
}

func (a *somaOncallResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Oncall = append(r.Oncall, somaOncallResult{ResultError: err})
	}
}

func (a *somaOncallResult) SomaAppendResult(r *somaResult) {
	r.Oncall = append(r.Oncall, *a)
}

/* Read Access
 *
 */
type somaOncallReadHandler struct {
	input     chan somaOncallRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
}

func (r *somaOncallReadHandler) run() {
	var err error

	r.list_stmt, err = r.conn.Prepare(`
SELECT oncall_duty_id,
       oncall_duty_name
FROM   inventory.oncall_duty_teams;`)
	if err != nil {
		log.Fatal("oncall/list: ", err)
	}
	defer r.list_stmt.Close()

	r.show_stmt, err = r.conn.Prepare(`
SELECT oncall_duty_id,
       oncall_duty_name,
	   oncall_duty_phone_number
FROM   inventory.oncall_duty_teams
WHERE  oncall_duty_id = $1;`)
	if err != nil {
		log.Fatal("oncall/show: ", err)
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

func (r *somaOncallReadHandler) process(q *somaOncallRequest) {
	var (
		oncallId, oncallName string
		oncallNumber         int
		rows                 *sql.Rows
		err                  error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: oncall/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&oncallId, &oncallName)
			result.Append(err, &somaOncallResult{
				Oncall: proto.Oncall{
					Id:   oncallId,
					Name: oncallName,
				},
			})
		}
	case "show":
		log.Printf("R: oncall/show for %s", q.Oncall.Id)
		err = r.show_stmt.QueryRow(q.Oncall.Id).Scan(
			&oncallId,
			&oncallName,
			&oncallNumber,
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

		result.Append(err, &somaOncallResult{
			Oncall: proto.Oncall{
				Id:     oncallId,
				Name:   oncallName,
				Number: strconv.Itoa(oncallNumber),
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaOncallWriteHandler struct {
	input    chan somaOncallRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	upd_stmt *sql.Stmt
	del_stmt *sql.Stmt
}

func (w *somaOncallWriteHandler) run() {
	var err error

	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO inventory.oncall_duty_teams (
	oncall_duty_id,
	oncall_duty_name,
	oncall_duty_phone_number)
SELECT $1::uuid, $2::varchar, $3::numeric WHERE NOT EXISTS (
	SELECT oncall_duty_id
	FROM inventory.oncall_duty_teams
	WHERE oncall_duty_id = $1::uuid
	OR oncall_duty_name = $2::varchar
	OR oncall_duty_phone_number = $3::numeric);`)
	if err != nil {
		log.Fatal("oncall/add: ", err)
	}
	defer w.add_stmt.Close()

	w.upd_stmt, err = w.conn.Prepare(`
UPDATE inventory.oncall_duty_teams
SET    oncall_duty_name = CASE WHEN $1::varchar IS NOT NULL
                          THEN $1::varchar
						  ELSE oncall_duty_name END,
       oncall_duty_phone_number = CASE WHEN $2::numeric IS NOT NULL
	                              THEN $2::numeric
								  ELSE oncall_duty_phone_number END
WHERE  oncall_duty_id = $3;`)
	if err != nil {
		log.Fatal("oncall/update: ", err)
	}
	defer w.upd_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM inventory.oncall_duty_teams
WHERE  oncall_duty_id = $1;`)
	if err != nil {
		log.Fatal("oncall/delete: ", err)
	}
	defer w.del_stmt.Close()

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

func (w *somaOncallWriteHandler) process(q *somaOncallRequest) {
	var (
		name   sql.NullString
		number sql.NullInt64
		res    sql.Result
		n      int // ensure err not redeclared in if block
		err    error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: oncall/add for %s", q.Oncall.Name)
		id := uuid.NewV4()
		res, err = w.add_stmt.Exec(
			id.String(),
			q.Oncall.Name,
			q.Oncall.Number,
		)
		q.Oncall.Id = id.String()
	case "update":
		log.Printf("R: oncall/update for %s", q.Oncall.Id)
		// our update statement uses NULL to check which of the values
		// should be updated
		if q.Oncall.Name == "" {
			name = sql.NullString{String: "", Valid: false}
		} else {
			name = sql.NullString{String: q.Oncall.Name, Valid: true}
		}

		if q.Oncall.Number != "" {
			n, err = strconv.Atoi(q.Oncall.Number)
			if err != nil {
				break
			}
			number = sql.NullInt64{Int64: int64(n), Valid: true}
		} else {
			number = sql.NullInt64{Int64: 0, Valid: false}
		}
		res, err = w.upd_stmt.Exec(
			name,
			number,
			q.Oncall.Id,
		)
	case "delete":
		log.Printf("R: oncall/del for %s", q.Oncall.Id)
		res, err = w.del_stmt.Exec(
			q.Oncall.Id,
		)
	default:
		log.Printf("R: unimplemented oncall/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaOncallResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaOncallResult{})
	default:
		result.Append(nil, &somaOncallResult{
			Oncall: q.Oncall,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
