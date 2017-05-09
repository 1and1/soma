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

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	"github.com/Sirupsen/logrus"
)

// ServerRead handles read requests for server
type ServerRead struct {
	Input               chan msg.Request
	Shutdown            chan struct{}
	conn                *sql.DB
	stmtList            *sql.Stmt
	stmtShow            *sql.Stmt
	stmtSync            *sql.Stmt
	stmtSearchByName    *sql.Stmt
	stmtSearchByAssetID *sql.Stmt
	appLog              *logrus.Logger
	reqLog              *logrus.Logger
	errLog              *logrus.Logger
}

// newServerRead return a new ServerRead handler with input buffer of length
func newServerRead(length int) (r *ServerRead) {
	r = &ServerRead{}
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (r *ServerRead) register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// run is the event loop for ServerRead
func (r *ServerRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ListServers:           r.stmtList,
		stmt.ShowServers:           r.stmtShow,
		stmt.SyncServers:           r.stmtSync,
		stmt.SearchServerByName:    r.stmtSearchByName,
		stmt.SearchServerByAssetId: r.stmtSearchByAssetID,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`servers`, err, stmt.Name(statement))
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
func (r *ServerRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case `list`:
		r.list(q, &result)
	case `show`:
		r.show(q, &result)
	case `sync`:
		r.sync(q, &result)
	case `search/name`:
		r.searchByName(q, &result)
	case `search/asset`:
		r.searchByAssetID(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all servers
func (r *ServerRead) list(q *msg.Request, mr *msg.Result) {
	var (
		serverID, serverName string
		serverAssetID        int
		rows                 *sql.Rows
		err                  error
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&serverID,
			&serverName,
			&serverAssetID,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Server = append(mr.Server, proto.Server{
			Id:      serverID,
			Name:    serverName,
			AssetId: uint64(serverAssetID),
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns the details of a specific server
func (r *ServerRead) show(q *msg.Request, mr *msg.Result) {
	var (
		err                         error
		serverID, serverDc          string
		serverDcLoc, serverName     string
		serverAssetID               int
		serverOnline, serverDeleted bool
	)

	if err = r.stmtShow.QueryRow(
		q.Server.Id,
	).Scan(
		&serverID,
		&serverAssetID,
		&serverDc,
		&serverDcLoc,
		&serverName,
		&serverOnline,
		&serverDeleted,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Server = append(mr.Server, proto.Server{
		Id:         serverID,
		AssetId:    uint64(serverAssetID),
		Datacenter: serverDc,
		Location:   serverDcLoc,
		Name:       serverName,
		IsOnline:   serverOnline,
		IsDeleted:  serverDeleted,
	})
	mr.OK()
}

// sync returns details for all servers suitable for sync processing
func (r *ServerRead) sync(q *msg.Request, mr *msg.Result) {
	var (
		err                         error
		serverID, serverDc          string
		serverDcLoc, serverName     string
		serverAssetID               int
		serverOnline, serverDeleted bool
		rows                        *sql.Rows
	)

	if rows, err = r.stmtSync.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&serverID,
			&serverAssetID,
			&serverDc,
			&serverDcLoc,
			&serverName,
			&serverOnline,
			&serverDeleted,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}

		mr.Server = append(mr.Server, proto.Server{
			Id:         serverID,
			AssetId:    uint64(serverAssetID),
			Datacenter: serverDc,
			Location:   serverDcLoc,
			Name:       serverName,
			IsOnline:   serverOnline,
			IsDeleted:  serverDeleted,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// searchByName looks up a server's ID by name
func (r *ServerRead) searchByName(q *msg.Request, mr *msg.Result) {
	var (
		err                  error
		serverID, serverName string
		serverAssetID        int
	)

	if err = r.stmtSearchByName.QueryRow(
		q.Search.Server.Name,
	).Scan(
		&serverID,
		&serverName,
		&serverAssetID,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Server = append(mr.Server, proto.Server{
		Id:      serverID,
		Name:    serverName,
		AssetId: uint64(serverAssetID),
	})
	mr.OK()
}

// searchByAssetID looks up a server's ID by assetID
func (r *ServerRead) searchByAssetID(q *msg.Request, mr *msg.Result) {
	var (
		err                  error
		serverID, serverName string
		serverAssetID        int
	)

	if err = r.stmtSearchByAssetID.QueryRow(
		q.Search.Server.AssetId,
	).Scan(
		&serverID,
		&serverName,
		&serverAssetID,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Server = append(mr.Server, proto.Server{
		Id:      serverID,
		Name:    serverName,
		AssetId: uint64(serverAssetID),
	})
	mr.OK()
}

// shutdown signals the handler to shut down
func (r *ServerRead) shutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
