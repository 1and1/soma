package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type somaValidityRequest struct {
	action   string
	Validity proto.Validity
	reply    chan somaResult
}

type somaValidityResult struct {
	ResultError error
	Validity    proto.Validity
}

func (a *somaValidityResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Validity = append(r.Validity, somaValidityResult{ResultError: err})
	}
}

func (a *somaValidityResult) SomaAppendResult(r *somaResult) {
	r.Validity = append(r.Validity, *a)
}

/* Read Access
 */
type somaValidityReadHandler struct {
	input     chan somaValidityRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaValidityReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ValidityList: r.list_stmt,
		stmt.ValidityShow: r.show_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`validity`, err, statement)
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

func (r *somaValidityReadHandler) process(q *somaValidityRequest) {
	var (
		property, object string
		inherited        bool
		rows             *sql.Rows
		err              error
		m                map[string]map[string]map[string]bool
	)
	result := somaResult{}

	switch q.action {
	case "list":
		r.appLog.Printf("R: validity/list")
		rows, err = r.list_stmt.Query()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&property, &object)
			result.Append(err, &somaValidityResult{
				Validity: proto.Validity{
					SystemProperty: property,
					ObjectType:     object,
				},
			})
		}
		if err = rows.Err(); err != nil {
			_ = result.SetRequestError(err)
			q.reply <- result
			return
		}
	case "show":
		r.appLog.Printf("R: status/show for %s", q.Validity.SystemProperty)
		rows, err = r.show_stmt.Query(q.Validity.SystemProperty)
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		m = make(map[string]map[string]map[string]bool)
		for rows.Next() {
			err = rows.Scan(
				&property,
				&object,
				&inherited,
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
			if m[property] == nil {
				m[property] = make(map[string]map[string]bool)
			}
			if m[property][object] == nil {
				m[property][object] = make(map[string]bool)
			}
			if inherited {
				m[property][object]["inherited"] = true
			} else {
				m[property][object]["direct"] = true
			}
		}
		if err = rows.Err(); err != nil {
			_ = result.SetRequestError(err)
			q.reply <- result
			return
		}
		for p_spec, _ := range m {
			for o_spec, _ := range m[p_spec] {
				result.Append(nil, &somaValidityResult{
					Validity: proto.Validity{
						SystemProperty: p_spec,
						ObjectType:     o_spec,
						Direct:         m[p_spec][o_spec]["direct"],
						Inherited:      m[p_spec][o_spec]["inherited"],
					},
				})
			}
		}
	default:
		r.errLog.Printf("R: unimplemented validity/%s", q.action)
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaValidityWriteHandler struct {
	input    chan somaValidityRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaValidityWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ValidityAdd: w.add_stmt,
		stmt.ValidityDel: w.del_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`validity`, err, statement)
		}
		defer prepStmt.Close()
	}

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

func (w *somaValidityWriteHandler) process(q *somaValidityRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		w.appLog.Printf("R: validity/add for %s", q.Validity.SystemProperty)
		if q.Validity.Direct {
			res, err = w.add_stmt.Exec(
				q.Validity.SystemProperty,
				q.Validity.ObjectType,
				false,
			)
		}
		if err != nil {
			goto errorout
		}
		if q.Validity.Inherited {
			res, err = w.add_stmt.Exec(
				q.Validity.SystemProperty,
				q.Validity.ObjectType,
				true,
			)
		}
	case "delete":
		w.appLog.Printf("R: validity/del for %s", q.Validity.SystemProperty)
		res, err = w.del_stmt.Exec(
			q.Validity.SystemProperty,
		)
	default:
		w.errLog.Printf("R: unimplemented validity/%s", q.action)
		result.SetNotImplemented()
		q.reply <- result
		return
	}
errorout:
	if result.SetRequestError(err) {
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		result.Append(errors.New("No rows affected"), &somaValidityResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaValidityResult{})
	default:
		result.Append(nil, &somaValidityResult{
			Validity: q.Validity,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaValidityReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaValidityWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
