package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

)

type somaProviderRequest struct {
	action   string
	Provider somaproto.ProtoProvider
	reply    chan somaResult
}

type somaProviderResult struct {
	ResultError error
	Provider    somaproto.ProtoProvider
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
}

func (r *somaProviderReadHandler) run() {
	var err error

	log.Println("Prepare: provider/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT metric_provider
FROM   soma.metric_providers;`)
	if err != nil {
		log.Fatal("provider/list: ", err)
	}
	defer r.list_stmt.Close()

	log.Println("Prepare: provider/show")
	r.show_stmt, err = r.conn.Prepare(`
SELECT metric_provider
FROM   soma.metric_providers
WHERE  metric_provider = $1::varchar;`)
	if err != nil {
		log.Fatal("provider/show: ", err)
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
		log.Printf("R: providers/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&provider)
			result.Append(err, &somaProviderResult{
				Provider: somaproto.ProtoProvider{
					Provider: provider,
				},
			})
		}
	case "show":
		log.Printf("R: providers/show for %s", q.Provider.Provider)
		err = r.show_stmt.QueryRow(q.Provider.Provider).Scan(
			&provider,
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

		result.Append(err, &somaProviderResult{
			Provider: somaproto.ProtoProvider{
				Provider: provider,
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
}

func (w *somaProviderWriteHandler) run() {
	var err error

	log.Println("Prepare: provider/add")
	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.metric_providers (
	metric_provider)
SELECT $1::varchar WHERE NOT EXISTS (
	SELECT metric_provider
	FROM   soma.metric_providers
	WHERE  metric_provider = $1::varchar);`)
	if err != nil {
		log.Fatal("provider/add: ", err)
	}
	defer w.add_stmt.Close()

	log.Println("Prepare: provider/delete")
	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.metric_providers
WHERE  metric_provider = $1::varchar;`)
	if err != nil {
		log.Fatal("provider/delete: ", err)
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
		log.Printf("R: providers/add for %s", q.Provider.Provider)
		res, err = w.add_stmt.Exec(
			q.Provider.Provider,
		)
	case "delete":
		log.Printf("R: providers/del for %s", q.Provider.Provider)
		res, err = w.del_stmt.Exec(
			q.Provider.Provider,
		)
	default:
		log.Printf("R: unimplemented providers/%s", q.action)
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
