/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"fmt"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	"github.com/Sirupsen/logrus"
)

// PropertyRead handles read requests for properties
type PropertyRead struct {
	Input            chan msg.Request
	Shutdown         chan struct{}
	conn             *sql.DB
	stmtListCustom   *sql.Stmt
	stmtListNative   *sql.Stmt
	stmtListService  *sql.Stmt
	stmtListSystem   *sql.Stmt
	stmtListTemplate *sql.Stmt
	stmtShowCustom   *sql.Stmt
	stmtShowNative   *sql.Stmt
	stmtShowService  *sql.Stmt
	stmtShowSystem   *sql.Stmt
	stmtShowTemplate *sql.Stmt
	appLog           *logrus.Logger
	reqLog           *logrus.Logger
	errLog           *logrus.Logger
}

// newPropertyRead return a new PropertyRead handler with input
// buffer of length
func newPropertyRead(length int) (r *PropertyRead) {
	r = &PropertyRead{}
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (r *PropertyRead) register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// run is the event loop for PropertyRead
func (r *PropertyRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.PropertyCustomList:   r.stmtListCustom,
		stmt.PropertyCustomShow:   r.stmtShowCustom,
		stmt.PropertyNativeList:   r.stmtListNative,
		stmt.PropertyNativeShow:   r.stmtShowNative,
		stmt.PropertyServiceList:  r.stmtListService,
		stmt.PropertyServiceShow:  r.stmtShowService,
		stmt.PropertySystemList:   r.stmtListSystem,
		stmt.PropertySystemShow:   r.stmtShowSystem,
		stmt.PropertyTemplateList: r.stmtListTemplate,
		stmt.PropertyTemplateShow: r.stmtShowTemplate,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`property`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-r.Shutdown:
			break runloop
		case req := <-r.Input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

// process is the request dispatcher
func (r *PropertyRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case `list`:
		r.list(q, &result)
	case `show`:
		r.show(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all properties
func (r *PropertyRead) list(q *msg.Request, mr *msg.Result) {
	switch q.Property.Type {
	case `custom`:
		r.listCustom(q, mr)
	case `native`:
		r.listNative(q, mr)
	case `service`:
		r.listService(q, mr)
	case `system`:
		r.listSystem(q, mr)
	case `template`:
		r.listTemplate(q, mr)
	default:
		mr.NotImplemented(fmt.Errorf("Unknown property type: %s",
			q.Property.Type))
	}
}

// listCustom returns all custom properties for a repository
func (r *PropertyRead) listCustom(q *msg.Request, mr *msg.Result) {
	var (
		property, repository, id string
		rows                     *sql.Rows
		err                      error
	)

	if rows, err = r.stmtListCustom.Query(
		q.Property.Custom.RepositoryId,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&id, &repository, &property); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Property = append(mr.Property, proto.Property{
			Type: q.Property.Type,
			Custom: &proto.PropertyCustom{
				Id:           id,
				RepositoryId: repository,
				Name:         property,
			},
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// listNative returns all native properties
func (r *PropertyRead) listNative(q *msg.Request, mr *msg.Result) {
	var (
		property string
		rows     *sql.Rows
		err      error
	)

	if rows, err = r.stmtListNative.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&property); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Property = append(mr.Property, proto.Property{
			Type: q.Property.Type,
			Native: &proto.PropertyNative{
				Name: property,
			},
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// listService returns all service properties for a team
func (r *PropertyRead) listService(q *msg.Request, mr *msg.Result) {
	var (
		property, team string
		rows           *sql.Rows
		err            error
	)

	if rows, err = r.stmtListService.Query(
		q.Property.Service.TeamId,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&property, &team); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Property = append(mr.Property, proto.Property{
			Type: q.Property.Type,
			Service: &proto.PropertyService{
				Name:   property,
				TeamId: team,
			},
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// listSystem returns all system properties
func (r *PropertyRead) listSystem(q *msg.Request, mr *msg.Result) {
	var (
		property string
		rows     *sql.Rows
		err      error
	)

	if rows, err = r.stmtListSystem.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&property); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Property = append(mr.Property, proto.Property{
			Type: q.Property.Type,
			System: &proto.PropertySystem{
				Name: property,
			},
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// listTemplate returns all service templates
func (r *PropertyRead) listTemplate(q *msg.Request, mr *msg.Result) {
	var (
		property string
		rows     *sql.Rows
		err      error
	)

	if rows, err = r.stmtListTemplate.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	for rows.Next() {
		if err = rows.Scan(&property); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Property = append(mr.Property, proto.Property{
			Type: q.Property.Type,
			Service: &proto.PropertyService{
				Name: property,
			},
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns details about a specific property
func (r *PropertyRead) show(q *msg.Request, mr *msg.Result) {
	switch q.Property.Type {
	case `custom`:
		r.showCustom(q, mr)
	case `native`:
		r.showNative(q, mr)
	case `service`:
		r.showService(q, mr)
	case `system`:
		r.showSystem(q, mr)
	case `template`:
		r.showTemplate(q, mr)
	default:
		mr.NotImplemented(fmt.Errorf("Unknown property type: %s",
			q.Property.Type))
	}
}

// showCustom returns the details for a specific custom property
func (r *PropertyRead) showCustom(q *msg.Request, mr *msg.Result) {
	var (
		property, repository, id string
		err                      error
	)

	if err = r.stmtShowCustom.QueryRow(
		q.Property.Custom.Id,
		q.Property.Custom.RepositoryId,
	).Scan(
		&id,
		&repository,
		&property,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Property = append(mr.Property, proto.Property{
		Type: q.Property.Type,
		Custom: &proto.PropertyCustom{
			Id:           id,
			RepositoryId: repository,
			Name:         property,
		},
	})
	mr.OK()
}

// showNative returns the details for a specific native property
func (r *PropertyRead) showNative(q *msg.Request, mr *msg.Result) {
	var (
		property string
		err      error
	)

	if err = r.stmtShowNative.QueryRow(
		q.Property.Native.Name,
	).Scan(
		&property,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Property = append(mr.Property, proto.Property{
		Type: q.Property.Type,
		Native: &proto.PropertyNative{
			Name: property,
		},
	})
	mr.OK()
}

// showService returns the details for a specific service
func (r *PropertyRead) showService(q *msg.Request, mr *msg.Result) {
	var (
		property, team, attribute, value string
		rows                             *sql.Rows
		err                              error
		service                          proto.PropertyService
	)

	if rows, err = r.stmtShowService.Query(
		q.Property.Service.Name,
		q.Property.Service.TeamId,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	service.Attributes = make([]proto.ServiceAttribute, 0)

	for rows.Next() {
		if err = rows.Scan(
			&property,
			&team,
			&attribute,
			&value,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}

		service.Name = property
		service.TeamId = team
		service.Attributes = append(service.Attributes,
			proto.ServiceAttribute{
				Name:  attribute,
				Value: value,
			})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	mr.Property = append(mr.Property, proto.Property{
		Type:    q.Property.Type,
		Service: &service,
	})
	mr.OK()
}

// showSystem returns the details about a specific system property
func (r *PropertyRead) showSystem(q *msg.Request, mr *msg.Result) {
	var (
		property string
		err      error
	)

	if err = r.stmtShowSystem.QueryRow(
		q.Property.System.Name,
	).Scan(
		&property,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Property = append(mr.Property, proto.Property{
		Type: q.Property.Type,
		System: &proto.PropertySystem{
			Name: property,
		},
	})
	mr.OK()
}

// showTemplate returns the details about a specific service
// template
func (r *PropertyRead) showTemplate(q *msg.Request, mr *msg.Result) {
	var (
		property, attribute, value string
		rows                       *sql.Rows
		err                        error
		template                   proto.PropertyService
	)

	if rows, err = r.stmtShowTemplate.Query(
		q.Property.Service.Name,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	template.Attributes = make([]proto.ServiceAttribute, 0)

	for rows.Next() {
		if err = rows.Scan(
			&property,
			&attribute,
			&value,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}

		template.Name = property
		template.Attributes = append(template.Attributes,
			proto.ServiceAttribute{
				Name:  attribute,
				Value: value,
			})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	mr.Property = append(mr.Property, proto.Property{
		Type:    q.Property.Type,
		Service: &template,
	})
	mr.OK()
}

// shutdown signals the handler to shut down
func (r *PropertyRead) shutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
