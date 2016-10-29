package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

type somaPropertyRequest struct {
	action  string
	prType  string
	System  proto.PropertySystem
	Native  proto.PropertyNative
	Service proto.PropertyService
	Custom  proto.PropertyCustom
	reply   chan somaResult
}

type somaPropertyResult struct {
	ResultError error
	prType      string
	System      proto.PropertySystem
	Native      proto.PropertyNative
	Service     proto.PropertyService
	Custom      proto.PropertyCustom
}

func (a *somaPropertyResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Properties = append(r.Properties, somaPropertyResult{ResultError: err})
	}
}

func (a *somaPropertyResult) SomaAppendResult(r *somaResult) {
	r.Properties = append(r.Properties, *a)
}

/* Read Access
 */
type somaPropertyReadHandler struct {
	input         chan somaPropertyRequest
	shutdown      chan bool
	conn          *sql.DB
	list_sys_stmt *sql.Stmt
	list_srv_stmt *sql.Stmt
	list_nat_stmt *sql.Stmt
	list_tpl_stmt *sql.Stmt
	list_cst_stmt *sql.Stmt
	show_sys_stmt *sql.Stmt
	show_srv_stmt *sql.Stmt
	show_nat_stmt *sql.Stmt
	show_tpl_stmt *sql.Stmt
	show_cst_stmt *sql.Stmt
	appLog        *log.Logger
	reqLog        *log.Logger
	errLog        *log.Logger
}

func (r *somaPropertyReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.PropertyCustomList:   r.list_cst_stmt,
		stmt.PropertyCustomShow:   r.show_cst_stmt,
		stmt.PropertyNativeList:   r.list_nat_stmt,
		stmt.PropertyNativeShow:   r.show_nat_stmt,
		stmt.PropertyServiceList:  r.list_srv_stmt,
		stmt.PropertyServiceShow:  r.show_srv_stmt,
		stmt.PropertySystemList:   r.list_sys_stmt,
		stmt.PropertySystemShow:   r.show_sys_stmt,
		stmt.PropertyTemplateList: r.list_tpl_stmt,
		stmt.PropertyTemplateShow: r.show_tpl_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`property`, err, stmt.Name(statement))
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

func (r *somaPropertyReadHandler) process(q *somaPropertyRequest) {
	var (
		property, team, repository, id, attribute, value string
		rows                                             *sql.Rows
		err                                              error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		switch q.prType {
		case "system":
			r.reqLog.Printf("R: property/list-system")
			rows, err = r.list_sys_stmt.Query()
		case "native":
			r.reqLog.Printf("R: property/list-native")
			rows, err = r.list_nat_stmt.Query()
		case "custom":
			r.reqLog.Printf("R: property/list-custom")
			rows, err = r.list_cst_stmt.Query(q.Custom.RepositoryId)
		case "service":
			r.reqLog.Printf("R: property/list-service")
			rows, err = r.list_srv_stmt.Query(q.Service.TeamId)
		case "template":
			r.reqLog.Printf("R: property/list-service-template")
			rows, err = r.list_tpl_stmt.Query()
		}
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			switch q.prType {
			case "system":
				err = rows.Scan(&property)
				result.Append(err, &somaPropertyResult{
					prType: q.prType,
					System: proto.PropertySystem{
						Name: property,
					},
				})
			case "service":
				err = rows.Scan(&property, &team)
				result.Append(err, &somaPropertyResult{
					prType: q.prType,
					Service: proto.PropertyService{
						Name:   property,
						TeamId: team,
					},
				})
			case "native":
				err = rows.Scan(&property)
				result.Append(err, &somaPropertyResult{
					prType: q.prType,
					Native: proto.PropertyNative{
						Name: property,
					},
				})
			case "template":
				err = rows.Scan(&property)
				result.Append(err, &somaPropertyResult{
					prType: q.prType,
					Service: proto.PropertyService{
						Name: property,
					},
				})
			case "custom":
				err = rows.Scan(&id, &repository, &property)
				result.Append(err, &somaPropertyResult{
					prType: q.prType,
					Custom: proto.PropertyCustom{
						Id:           id,
						RepositoryId: repository,
						Name:         property,
					},
				})
			}
		}
	case "show":
		switch q.prType {
		case "system":
			r.reqLog.Printf("R: property/show-system for %s", q.System.Name)
			err = r.show_sys_stmt.QueryRow(q.System.Name).Scan(
				&property,
			)
		case "native":
			r.reqLog.Printf("R: property/show-native for %s", q.Native.Name)
			err = r.show_nat_stmt.QueryRow(q.Native.Name).Scan(
				&property,
			)
		case "custom":
			r.reqLog.Printf("R: property/show-custom for %s", q.Custom.Id)
			err = r.show_cst_stmt.QueryRow(
				q.Custom.Id,
				q.Custom.RepositoryId,
			).Scan(
				&id,
				&repository,
				&property,
			)
		case "service":
			r.reqLog.Printf("R: property/show-service for %s/%s", q.Service.TeamId, q.Service.Name)
			rows, err = r.show_srv_stmt.Query(
				q.Service.Name,
				q.Service.TeamId,
			)
			defer rows.Close()
		case "template":
			r.reqLog.Printf("R: property/show-service-template for %s", q.Service.Name)
			rows, err = r.show_tpl_stmt.Query(
				q.Service.Name,
			)
			defer rows.Close()
		}
		if err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		switch q.prType {
		case "system":
			result.Append(err, &somaPropertyResult{
				prType: q.prType,
				System: proto.PropertySystem{
					Name: property,
				},
			})
		case "native":
			result.Append(err, &somaPropertyResult{
				prType: q.prType,
				Native: proto.PropertyNative{
					Name: property,
				},
			})
		case "custom":
			result.Append(err, &somaPropertyResult{
				prType: q.prType,
				Custom: proto.PropertyCustom{
					Id:           id,
					RepositoryId: repository,
					Name:         property,
				},
			})
		case "service":
			propTempl := proto.PropertyService{}
			var fErr error
			for rows.Next() {
				err := rows.Scan(
					&property,
					&team,
					&attribute,
					&value,
				)

				if err == nil {
					propTempl.Name = property
					propTempl.TeamId = team
					propTempl.Attributes = append(propTempl.Attributes, proto.ServiceAttribute{
						Name:  attribute,
						Value: value,
					})
				} else {
					fErr = err
				}
			}
			result.Append(fErr, &somaPropertyResult{
				prType:  q.prType,
				Service: propTempl,
			})
		case "template":
			propTempl := proto.PropertyService{}
			propTempl.Attributes = make([]proto.ServiceAttribute, 0)
			var fErr error
			for rows.Next() {
				err := rows.Scan(
					&property,
					&attribute,
					&value,
				)

				if err == nil {
					propTempl.Name = property
					propTempl.Attributes = append(propTempl.Attributes, proto.ServiceAttribute{
						Name:  attribute,
						Value: value,
					})
				} else {
					fErr = err
				}
			}
			result.Append(fErr, &somaPropertyResult{
				prType:  q.prType,
				Service: propTempl,
			})
		}
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaPropertyWriteHandler struct {
	input             chan somaPropertyRequest
	shutdown          chan bool
	conn              *sql.DB
	add_sys_stmt      *sql.Stmt
	add_nat_stmt      *sql.Stmt
	add_cst_stmt      *sql.Stmt
	add_srv_stmt      *sql.Stmt
	add_tpl_stmt      *sql.Stmt
	add_srv_attr_stmt *sql.Stmt
	add_tpl_attr_stmt *sql.Stmt
	del_sys_stmt      *sql.Stmt
	del_nat_stmt      *sql.Stmt
	del_cst_stmt      *sql.Stmt
	del_srv_stmt      *sql.Stmt
	del_tpl_stmt      *sql.Stmt
	del_srv_attr_stmt *sql.Stmt
	del_tpl_attr_stmt *sql.Stmt
	appLog            *log.Logger
	reqLog            *log.Logger
	errLog            *log.Logger
}

func (w *somaPropertyWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.PropertyCustomAdd:            w.add_cst_stmt,
		stmt.PropertyCustomDel:            w.del_cst_stmt,
		stmt.PropertyNativeAdd:            w.add_nat_stmt,
		stmt.PropertyNativeDel:            w.del_nat_stmt,
		stmt.PropertyServiceAdd:           w.add_srv_stmt,
		stmt.PropertyServiceAttributeAdd:  w.add_srv_attr_stmt,
		stmt.PropertyServiceAttributeDel:  w.del_srv_attr_stmt,
		stmt.PropertyServiceDel:           w.del_srv_stmt,
		stmt.PropertySystemAdd:            w.add_sys_stmt,
		stmt.PropertySystemDel:            w.del_sys_stmt,
		stmt.PropertyTemplateAdd:          w.add_tpl_stmt,
		stmt.PropertyTemplateAttributeAdd: w.add_tpl_attr_stmt,
		stmt.PropertyTemplateAttributeDel: w.del_tpl_attr_stmt,
		stmt.PropertyTemplateDel:          w.del_tpl_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`property`, err, stmt.Name(statement))
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

func (w *somaPropertyWriteHandler) process(q *somaPropertyRequest) {
	var (
		res    sql.Result
		err    error
		tx     *sql.Tx
		attr   proto.ServiceAttribute
		rowCnt int64
	)
	result := somaResult{}

	switch q.action {
	case "add":
		switch q.prType {
		case "system":
			w.reqLog.Printf("R: property/add-system for %s", q.System.Name)
			res, err = w.add_sys_stmt.Exec(
				q.System.Name,
			)
			rowCnt, _ = res.RowsAffected()
		case "native":
			w.reqLog.Printf("R: property/add-native for %s", q.Native.Name)
			res, err = w.add_nat_stmt.Exec(
				q.Native.Name,
			)
			rowCnt, _ = res.RowsAffected()
		case "custom":
			q.Custom.Id = uuid.NewV4().String()
			w.reqLog.Printf("R: property/add-custom for %s", q.Custom.Name)
			res, err = w.add_cst_stmt.Exec(
				q.Custom.Id,
				q.Custom.RepositoryId,
				q.Custom.Name,
			)
			rowCnt, _ = res.RowsAffected()
		case "service":
			w.reqLog.Printf("R: property/add-service for %s/%s", q.Service.TeamId, q.Service.Name)
			tx, err = w.conn.Begin()
			if err != nil {
				goto bailout
			}
			defer tx.Rollback()

			res, err = tx.Stmt(w.add_srv_stmt).Exec(
				q.Service.TeamId,
				q.Service.Name,
			)
			if err != nil {
				goto bailout
			}
			rowCnt, _ = res.RowsAffected()
			if rowCnt == 0 {
				goto bailout
			}

			for _, attr = range q.Service.Attributes {
				res, err = tx.Stmt(w.add_srv_attr_stmt).Exec(
					q.Service.TeamId,
					q.Service.Name,
					attr.Name,
					attr.Value,
				)
				if err != nil {
					break
				}
			}
			if err != nil {
				goto bailout
			}

			err = tx.Commit()
		case "template":
			w.reqLog.Printf("R: property/add-service-template for %s", q.Service.Name)
			tx, err = w.conn.Begin()
			if err != nil {
				goto bailout
			}
			defer tx.Rollback()

			res, err = tx.Stmt(w.add_tpl_stmt).Exec(
				q.Service.Name,
			)
			if err != nil {
				goto bailout
			}
			rowCnt, _ = res.RowsAffected()
			if rowCnt == 0 {
				goto bailout
			}

			for _, attr = range q.Service.Attributes {
				res, err = tx.Stmt(w.add_tpl_attr_stmt).Exec(
					q.Service.Name,
					attr.Name,
					attr.Value,
				)
				if err != nil {
					break
				}
			}
			if err != nil {
				goto bailout
			}

			err = tx.Commit()
		}
	case "delete":
		switch q.prType {
		case "system":
			w.reqLog.Printf("R: property/delete-system for %s", q.System.Name)
			res, err = w.del_sys_stmt.Exec(
				q.System.Name,
			)
			rowCnt, _ = res.RowsAffected()
		case "native":
			w.reqLog.Printf("R: property/delete-native for %s", q.Native.Name)
			res, err = w.del_nat_stmt.Exec(
				q.Native.Name,
			)
			rowCnt, _ = res.RowsAffected()
		case "custom":
			w.reqLog.Printf("R: property/delete-custom for %s", q.Custom.Id)
			res, err = w.del_cst_stmt.Exec(
				q.Custom.RepositoryId,
				q.Custom.Id,
			)
			rowCnt, _ = res.RowsAffected()
		case "service":
			w.reqLog.Printf("R: property/delete-service for %s/%s", q.Service.TeamId, q.Service.Name)
			tx, err = w.conn.Begin()
			if err != nil {
				goto bailout
			}
			defer tx.Rollback()

			res, err = tx.Stmt(w.del_srv_attr_stmt).Exec(
				q.Service.TeamId,
				q.Service.Name,
			)
			if err != nil {
				goto bailout
			}

			res, err = tx.Stmt(w.del_srv_stmt).Exec(
				q.Service.TeamId,
				q.Service.Name,
			)
			if err != nil {
				goto bailout
			}

			rowCnt, _ = res.RowsAffected()
			err = tx.Commit()
		case "template":
			w.reqLog.Printf("R: property/delete-service-template for %s", q.Service.Name)
			tx, err = w.conn.Begin()
			if err != nil {
				goto bailout
			}
			defer tx.Rollback()

			res, err = tx.Stmt(w.del_tpl_attr_stmt).Exec(
				q.Service.Name,
			)
			if err != nil {
				goto bailout
			}

			res, err = tx.Stmt(w.del_tpl_stmt).Exec(
				q.Service.Name,
			)
			if err != nil {
				goto bailout
			}

			rowCnt, _ = res.RowsAffected()
			err = tx.Commit()
		} // switch q.prType
	default:
		w.reqLog.Printf("R: unimplemented property/%s", q.action)
		result.SetNotImplemented()
		q.reply <- result
		return
	}

bailout:
	if result.SetRequestError(err) {
		q.reply <- result
		return
	}

	switch {
	case rowCnt == 0:
		result.Append(errors.New("No rows affected"), &somaPropertyResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaPropertyResult{})
	default:
		switch q.prType {
		case "system":
			result.Append(nil, &somaPropertyResult{
				prType: q.prType,
				System: q.System,
			})
		case "native":
			result.Append(nil, &somaPropertyResult{
				prType: q.prType,
				Native: q.Native,
			})
		case "custom":
			result.Append(nil, &somaPropertyResult{
				prType: q.prType,
				Custom: q.Custom,
			})
		case "service":
			result.Append(nil, &somaPropertyResult{
				prType:  q.prType,
				Service: q.Service,
			})
		case "template":
			result.Append(nil, &somaPropertyResult{
				prType:  q.prType,
				Service: q.Service,
			})
		}
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaPropertyReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaPropertyWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
