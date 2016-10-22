package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type somaProviderRequest struct {
	action   string
	Provider proto.Provider
	reply    chan somaResult
}

type somaProviderResult struct {
	ResultError error
	Provider    proto.Provider
}

func (a *somaProviderResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Providers = append(r.Providers, somaProviderResult{ResultError: err})
	}
}

func (a *somaProviderResult) SomaAppendResult(r *somaResult) {
	r.Providers = append(r.Providers, *a)
}

/* Read Access
 */
type somaProviderReadHandler struct {
	input     chan somaProviderRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaProviderReadHandler) run() {
	var err error

	r.list_stmt, err = r.conn.Prepare(`
SELECT metric_provider
FROM   soma.metric_providers;`)
	if err != nil {
		r.errLog.Fatal("provider/list: ", err)
	}
	defer r.list_stmt.Close()

	r.show_stmt, err = r.conn.Prepare(`
SELECT metric_provider
FROM   soma.metric_providers
WHERE  metric_provider = $1::varchar;`)
	if err != nil {
		r.errLog.Fatal("provider/show: ", err)
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

func (r *somaProviderReadHandler) process(q *somaProviderRequest) {
	var (
		provider string
		rows     *sql.Rows
		err      error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		r.reqLog.Printf("R: providers/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&provider)
			result.Append(err, &somaProviderResult{
				Provider: proto.Provider{
					Name: provider,
				},
			})
		}
	case "show":
		r.reqLog.Printf("R: providers/show for %s", q.Provider.Name)
		err = r.show_stmt.QueryRow(q.Provider.Name).Scan(
			&provider,
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

		result.Append(err, &somaProviderResult{
			Provider: proto.Provider{
				Name: provider,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaProviderWriteHandler struct {
	input    chan somaProviderRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaProviderWriteHandler) run() {
	var err error

	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.metric_providers (
	metric_provider)
SELECT $1::varchar WHERE NOT EXISTS (
	SELECT metric_provider
	FROM   soma.metric_providers
	WHERE  metric_provider = $1::varchar);`)
	if err != nil {
		w.errLog.Fatal("provider/add: ", err)
	}
	defer w.add_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.metric_providers
WHERE  metric_provider = $1::varchar;`)
	if err != nil {
		w.errLog.Fatal("provider/delete: ", err)
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

func (w *somaProviderWriteHandler) process(q *somaProviderRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		w.reqLog.Printf("R: providers/add for %s", q.Provider.Name)
		res, err = w.add_stmt.Exec(
			q.Provider.Name,
		)
	case "delete":
		w.reqLog.Printf("R: providers/del for %s", q.Provider.Name)
		res, err = w.del_stmt.Exec(
			q.Provider.Name,
		)
	default:
		w.reqLog.Printf("R: unimplemented providers/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaProviderResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaProviderResult{})
	default:
		result.Append(nil, &somaProviderResult{
			Provider: q.Provider,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaProviderReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaProviderWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
