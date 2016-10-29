package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

// Message structs
type somaDatacenterRequest struct {
	action     string
	Datacenter proto.Datacenter
	rename     string
	reply      chan somaResult
}

type somaDatacenterResult struct {
	ResultError error
	Datacenter  proto.Datacenter
}

func (a *somaDatacenterResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Datacenters = append(r.Datacenters, somaDatacenterResult{ResultError: err})
	}
}

func (a *somaDatacenterResult) SomaAppendResult(r *somaResult) {
	r.Datacenters = append(r.Datacenters, *a)
}

/*  Read Access
 *
 */
type somaDatacenterReadHandler struct {
	input     chan somaDatacenterRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	grp_list  *sql.Stmt
	grp_show  *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaDatacenterReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.DatacenterList:      r.list_stmt,
		stmt.DatacenterShow:      r.show_stmt,
		stmt.DatacenterGroupList: r.grp_list,
		stmt.DatacenterGroupShow: r.grp_show,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`datacenter`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

	for {
		select {
		case <-r.shutdown:
			break
		case req := <-r.input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

func (r *somaDatacenterReadHandler) process(q *somaDatacenterRequest) {
	var datacenter string
	var rows *sql.Rows
	var err error
	result := somaResult{}

	switch q.action {
	case `sync`:
		r.appLog.Printf("R: datacenter/sync")
		// right now, sync and list are the same. This allows to later
		// change the sync result if required without disturbing list
		fallthrough
	case "list":
		r.appLog.Printf("R: datacenter/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&datacenter)
			result.Append(err, &somaDatacenterResult{
				Datacenter: proto.Datacenter{
					Locode: datacenter,
				},
			})
		}
	case "show":
		r.appLog.Printf("R: datacenter/show for %s", q.Datacenter.Locode)
		err = r.show_stmt.QueryRow(q.Datacenter.Locode).Scan(&datacenter)
		if err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		result.Append(err, &somaDatacenterResult{
			Datacenter: proto.Datacenter{
				Locode: datacenter,
			},
		})
		/*
			case "grouplist":
				rows, err = r.grp_list.Query()
				defer rows.Close()
				if err != nil {
					result = append(result, somaDatacenterResult{
						err:        err,
						datacenter: q.datacenter,
					})
					q.reply <- result
					return
				}

				for rows.Next() {
					err = rows.Scan(&datacenter)
					if err != nil {
						result = append(result, somaDatacenterResult{
							err:        err,
							datacenter: q.datacenter,
						})
						err = nil
						continue
					}
					result = append(result, somaDatacenterResult{
						err:        nil,
						datacenter: datacenter,
					})
				}
			case "groupshow":
				rows, err = r.grp_show.Query(q.datacenter)
				if err != nil {
					result = append(result, somaDatacenterResult{
						err:        err,
						datacenter: q.datacenter,
					})
					q.reply <- result
					return
				}

				for rows.Next() {
					err = rows.Scan(&datacenter)
					if err != nil {
						result = append(result, somaDatacenterResult{
							err:        err,
							datacenter: q.datacenter,
						})
						err = nil
						continue
					}
					result = append(result, somaDatacenterResult{
						err:        nil,
						datacenter: datacenter,
					})
				}
		*/
	default:
		r.errLog.Printf("R: unimplemented datacenter/%s", q.action)
		result.SetNotImplemented()
	}
	q.reply <- result
}

/*
 * Write Access
 */

type somaDatacenterWriteHandler struct {
	input    chan somaDatacenterRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	ren_stmt *sql.Stmt
	grp_add  *sql.Stmt
	grp_del  *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaDatacenterWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.DatacenterAdd:      w.add_stmt,
		stmt.DatacenterDel:      w.del_stmt,
		stmt.DatacenterRename:   w.ren_stmt,
		stmt.DatacenterGroupAdd: w.grp_add,
		stmt.DatacenterGroupDel: w.grp_del,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`datacenter`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

	for {
		select {
		case <-w.shutdown:
			break
		case req := <-w.input:
			w.process(&req)
		}
	}
}

func (w *somaDatacenterWriteHandler) process(q *somaDatacenterRequest) {
	var res sql.Result
	var err error

	result := somaResult{}
	switch q.action {
	case "add":
		res, err = w.add_stmt.Exec(q.Datacenter.Locode)
	//case "groupadd":
	//	res, err = w.grp_add.Exec(q.group, q.datacenter, q.group, q.datacenter)
	case "delete":
		res, err = w.del_stmt.Exec(q.Datacenter.Locode)
	//case "groupdel":
	//	res, err = w.grp_del.Exec(q.group, q.datacenter)
	case "rename":
		res, err = w.ren_stmt.Exec(q.rename, q.Datacenter.Locode)
	default:
		w.errLog.Printf("R: unimplemented datacenter/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaDatacenterResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaDatacenterResult{})
	default:
		result.Append(nil, &somaDatacenterResult{
			Datacenter: q.Datacenter,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaDatacenterReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaDatacenterWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
