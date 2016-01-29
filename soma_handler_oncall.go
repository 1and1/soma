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
	oncall somaproto.ProtoOncall
	reply  chan []somaOncallResult
}

type somaOncallResult struct {
	rErr   error
	lErr   error
	oncall somaproto.ProtoOncall
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
SELECT oncall_duty_id, oncall_duty_name
FROM inventory.oncall_duty_teams;`)
	if err != nil {
		log.Fatal(err)
	}

	r.show_stmt, err = r.conn.Prepare(`
SELECT oncall_duty_id, oncall_duty_name, oncall_duty_phone_number
FROM inventory.oncall_duty_teams
WHERE oncall_duty_id = $1;`)
	if err != nil {
		log.Fatal(err)
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

func (r *somaOncallReadHandler) process(q *somaOncallRequest) {
	var oncallId, oncallName string
	var oncallNumber int
	var rows *sql.Rows
	var err error
	result := make([]somaOncallResult, 0)

	switch q.action {
	case "list":
		log.Printf("R: oncall/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if err != nil {
			result = append(result, somaOncallResult{
				rErr: err,
			})
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&oncallId, &oncallName)
			if err != nil {
				result = append(result, somaOncallResult{
					lErr: err,
				})
				err = nil
				continue
			}
			result = append(result, somaOncallResult{
				oncall: somaproto.ProtoOncall{
					Id:   oncallId,
					Name: oncallName,
				},
			})
		}
	case "show":
		log.Printf("R: oncall/show for %s", q.oncall.Id)
		err = r.show_stmt.QueryRow(q.oncall.Id).Scan(&oncallId, &oncallName, &oncallNumber)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				result = append(result, somaOncallResult{
					rErr: err,
				})
			}
			q.reply <- result
			return
		}

		result = append(result, somaOncallResult{
			oncall: somaproto.ProtoOncall{
				Id:     oncallId,
				Name:   oncallName,
				Number: strconv.Itoa(oncallNumber),
			},
		})
	default:
		result = append(result, somaOncallResult{
			rErr: errors.New("not implemented"),
		})
	}
	q.reply <- result
}

/* Write Access
 *
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

	log.Println("Prepare: oncall/add")
	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO inventory.oncall_duty_teams (
	oncall_duty_id,
	oncall_duty_name,
	oncall_duty_phone_number)
SELECT $1, $2, $3 WHERE NOT EXISTS (
	SELECT oncall_duty_id
	FROM inventory.oncall_duty_teams
	WHERE oncall_duty_id = $4
	OR oncall_duty_name = $5
	OR oncall_duty_phone_number = $6);`)
	if err != nil {
		log.Fatal(err)
	}
	defer w.add_stmt.Close()

	log.Println("Prepare: oncall/upd")
	w.upd_stmt, err = w.conn.Prepare(`
UPDATE inventory.oncall_duty_teams
SET oncall_duty_name = CASE WHEN $1::varchar IS NOT NULL THEN $2::varchar ELSE oncall_duty_name END,
    oncall_duty_phone_number = CASE WHEN $3::numeric IS NOT NULL THEN $4::numeric ELSE oncall_duty_phone_number END
WHERE oncall_duty_id = $5;`)
	if err != nil {
		log.Fatal(err)
	}
	defer w.upd_stmt.Close()

	log.Println("Prepare: oncall/del")
	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM inventory.oncall_duty_teams
WHERE oncall_duty_id = $1;`)
	if err != nil {
		log.Fatal(err)
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
	var res sql.Result
	var err error
	result := make([]somaOncallResult, 0)

	switch q.action {
	case "add":
		log.Printf("R: oncall/add for %s", q.oncall.Name)
		id := uuid.NewV4()
		res, err = w.add_stmt.Exec(
			id.String(),
			q.oncall.Name,
			q.oncall.Number,
			id.String(),
			q.oncall.Name,
			q.oncall.Number,
		)
		q.oncall.Id = id.String()
	case "update":
		log.Printf("R: oncall/update for %s", q.oncall.Id)
		// our update statement uses NULL to check which of the values
		// should be updated
		var name sql.NullString
		if q.oncall.Name == "" {
			name = sql.NullString{String: "", Valid: false}
		} else {
			name = sql.NullString{String: q.oncall.Name, Valid: true}
		}

		var n int // ensure err not redeclared in if block
		var number sql.NullInt64
		if q.oncall.Number != "" {
			n, err = strconv.Atoi(q.oncall.Number)
			if err != nil {
				break
			}
			number = sql.NullInt64{Int64: int64(n), Valid: true}
		} else {
			number = sql.NullInt64{Int64: 0, Valid: false}
		}
		res, err = w.upd_stmt.Exec(
			name,
			name,
			number,
			number,
			q.oncall.Id,
		)
	case "delete":
		log.Printf("R: oncall/del for %s", q.oncall.Id)
		res, err = w.del_stmt.Exec(
			q.oncall.Id,
		)
	default:
		log.Printf("R: unimplemented oncall/%s", q.action)
		result = append(result, somaOncallResult{
			rErr: errors.New("not implemented"),
		})
		q.reply <- result
		return
	}
	if err != nil {
		result = append(result, somaOncallResult{
			rErr: err,
		})
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		result = append(result, somaOncallResult{
			lErr: errors.New("No rows affected"),
		})
	case rowCnt > 1:
		result = append(result, somaOncallResult{
			lErr: fmt.Errorf("Too many rows affected: %d", rowCnt),
		})
	default:
		result = append(result, somaOncallResult{
			oncall: q.oncall,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
