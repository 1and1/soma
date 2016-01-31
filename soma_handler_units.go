package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

)

type somaUnitRequest struct {
	action string
	Unit   somaproto.ProtoUnit
	reply  chan somaResult
}

type somaUnitResult struct {
	ResultError error
	Unit        somaproto.ProtoUnit
}

func (a *somaUnitResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Units = append(r.Units, somaUnitResult{ResultError: err})
	}
}

func (a *somaUnitResult) SomaAppendResult(r *somaResult) {
	r.Units = append(r.Units, *a)
}

/* Read Access
 */
type somaUnitReadHandler struct {
	input     chan somaUnitRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
}

func (r *somaUnitReadHandler) run() {
	var err error

	log.Println("Prepare: unit/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT metric_unit
FROM   soma.metric_units;`)
	if err != nil {
		log.Fatal("unit/list: ", err)
	}
	defer r.list_stmt.Close()

	log.Println("Prepare: unit/show")
	r.show_stmt, err = r.conn.Prepare(`
SELECT metric_unit,
       metric_unit_long_name
FROM   soma.metric_units
WHERE  metric_unit = $1::varchar;`)
	if err != nil {
		log.Fatal("unit/show: ", err)
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

func (r *somaUnitReadHandler) process(q *somaUnitRequest) {
	var (
		unit, name string
		rows       *sql.Rows
		err        error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: units/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&unit)
			result.Append(err, &somaUnitResult{
				Unit: somaproto.ProtoUnit{
					Name: unit,
				},
			})
		}
	case "show":
		log.Printf("R: units/show for %s", q.Unit.Unit)
		err = r.show_stmt.QueryRow(q.Unit.Unit).Scan(
			&unit,
			&name,
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

		result.Append(err, &somaUnitResult{
			Unit: somaproto.ProtoUnit{
				Unit: unit,
				Name: name,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaUnitWriteHandler struct {
	input    chan somaUnitRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
}

func (w *somaUnitWriteHandler) run() {
	var err error

	log.Println("Prepare: unit/add")
	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.metric_units (
	metric_unit,
	metric_unit_long_name)
SELECT $1::varchar, $2::varchar WHERE NOT EXISTS (
	SELECT metric_unit
	FROM   soma.metric_units
	WHERE  metric_unit = $1::varchar
	OR     metric_unit_long_name = $2::varchar);`)
	if err != nil {
		log.Fatal("unit/add: ", err)
	}
	defer w.add_stmt.Close()

	log.Println("Prepare: unit/delete")
	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.metric_units
WHERE  metric_unit = $1::varchar;`)
	if err != nil {
		log.Fatal("unit/delete: ", err)
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

func (w *somaUnitWriteHandler) process(q *somaUnitRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: units/add for %s", q.Unit.Unit)
		res, err = w.add_stmt.Exec(
			q.Unit.Unit,
			q.Unit.Name,
		)
	case "delete":
		log.Printf("R: units/del for %s", q.Unit.Unit)
		res, err = w.del_stmt.Exec(
			q.Unit.Unit,
		)
	default:
		log.Printf("R: unimplemented units/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaUnitResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaUnitResult{})
	default:
		result.Append(nil, &somaUnitResult{
			Unit: q.Unit,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix