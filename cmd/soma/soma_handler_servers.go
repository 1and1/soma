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

type somaServerRequest struct {
	action string
	Server proto.Server
	Filter proto.Filter
	reply  chan somaResult
}

type somaServerResult struct {
	ResultError error
	Server      proto.Server
}

func (a *somaServerResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Servers = append(r.Servers, somaServerResult{ResultError: err})
	}
}

func (a *somaServerResult) SomaAppendResult(r *somaResult) {
	r.Servers = append(r.Servers, *a)
}

/* Read Access
 */
type somaServerReadHandler struct {
	input     chan somaServerRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	sync_stmt *sql.Stmt
	snam_stmt *sql.Stmt
	sass_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaServerReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ListServers:           r.list_stmt,
		stmt.ShowServers:           r.show_stmt,
		stmt.SyncServers:           r.sync_stmt,
		stmt.SearchServerByName:    r.snam_stmt,
		stmt.SearchServerByAssetId: r.sass_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`servers`, err, stmt.Name(statement))
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

func (r *somaServerReadHandler) process(q *somaServerRequest) {
	var (
		serverId, serverDc, serverDcLoc, serverName string
		serverAsset                                 int
		serverOnline, serverDeleted                 bool
		rows                                        *sql.Rows
		err                                         error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		r.reqLog.Printf("R: server/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&serverId, &serverName, &serverAsset)
			result.Append(err, &somaServerResult{
				Server: proto.Server{
					Id:      serverId,
					Name:    serverName,
					AssetId: uint64(serverAsset),
				},
			})
		}
	case `sync`:
		r.reqLog.Printf("R: server/sync")
		rows, err = r.sync_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(
				&serverId,
				&serverAsset,
				&serverDc,
				&serverDcLoc,
				&serverName,
				&serverOnline,
				&serverDeleted,
			)

			result.Append(err, &somaServerResult{
				Server: proto.Server{
					Id:         serverId,
					AssetId:    uint64(serverAsset),
					Datacenter: serverDc,
					Location:   serverDcLoc,
					Name:       serverName,
					IsOnline:   serverOnline,
					IsDeleted:  serverDeleted,
				},
			})
		}
	case `search/name`:
		r.reqLog.Printf("R: server/search-name for %s", q.Filter.Server.Name)
		if err = r.snam_stmt.QueryRow(q.Filter.Server.Name).Scan(
			&serverId,
			&serverName,
			&serverAsset,
		); err == sql.ErrNoRows {
			result.SetNotFound()
		} else if err != nil {
			_ = result.SetRequestError(err)
		} else {
			result.Append(nil, &somaServerResult{
				Server: proto.Server{
					Id:      serverId,
					Name:    serverName,
					AssetId: uint64(serverAsset),
				},
			})
		}
	case `search/asset`:
		r.reqLog.Printf("R: server/search-asset for %d", q.Filter.Server.AssetId)
		if err = r.sass_stmt.QueryRow(q.Filter.Server.AssetId).Scan(
			&serverId,
			&serverName,
			&serverAsset,
		); err == sql.ErrNoRows {
			result.SetNotFound()
		} else if err != nil {
			_ = result.SetRequestError(err)
		} else {
			result.Append(nil, &somaServerResult{
				Server: proto.Server{
					Id:      serverId,
					Name:    serverName,
					AssetId: uint64(serverAsset),
				},
			})
		}
	case "show":
		r.reqLog.Printf("R: server/show for %s", q.Server.Id)
		err = r.show_stmt.QueryRow(q.Server.Id).Scan(
			&serverId,
			&serverAsset,
			&serverDc,
			&serverDcLoc,
			&serverName,
			&serverOnline,
			&serverDeleted,
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

		result.Append(err, &somaServerResult{
			Server: proto.Server{
				Id:         serverId,
				AssetId:    uint64(serverAsset),
				Datacenter: serverDc,
				Location:   serverDcLoc,
				Name:       serverName,
				IsOnline:   serverOnline,
				IsDeleted:  serverDeleted,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaServerWriteHandler struct {
	input    chan somaServerRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	prg_stmt *sql.Stmt
	upd_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaServerWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.AddServers:    w.add_stmt,
		stmt.DeleteServers: w.del_stmt,
		stmt.PurgeServers:  w.prg_stmt,
		stmt.UpdateServers: w.upd_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`server`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-w.shutdown:
			break runloop
		case req := <-w.input:
			go func() {
				w.process(&req)
			}()
		}
	}
}

func (w *somaServerWriteHandler) process(q *somaServerRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		w.reqLog.Printf("R: server/add for %s", q.Server.Name)
		id := uuid.NewV4()
		res, err = w.add_stmt.Exec(
			id.String(),
			q.Server.AssetId,
			q.Server.Datacenter,
			q.Server.Location,
			q.Server.Name,
			q.Server.IsOnline,
			false,
		)
		q.Server.Id = id.String()
	case "remove":
		w.reqLog.Printf("R: server/remove for %s", q.Server.Id)
		res, err = w.del_stmt.Exec(
			q.Server.Id,
		)
	case "purge":
		w.reqLog.Printf("R: server/purge for %s", q.Server.Id)
		res, err = w.del_stmt.Exec(
			q.Server.Id,
		)
	case `update`:
		w.reqLog.Printf("R: server/update for %s", q.Server.Id)
		res, err = w.upd_stmt.Exec(
			q.Server.Id,
			q.Server.AssetId,
			q.Server.Datacenter,
			q.Server.Location,
			q.Server.Name,
			q.Server.IsOnline,
			q.Server.IsDeleted,
		)
	case "insert-null":
		w.reqLog.Printf("R: server/insert-null")
		q.Server.Id = "00000000-0000-0000-0000-000000000000"
		q.Server.AssetId = 0
		q.Server.Location = "none"
		q.Server.Name = "soma-null-server"
		q.Server.IsOnline = true
		q.Server.IsDeleted = false
		res, err = w.add_stmt.Exec(
			q.Server.Id,
			q.Server.AssetId,
			q.Server.Datacenter,
			q.Server.Location,
			q.Server.Name,
			q.Server.IsOnline,
			q.Server.IsDeleted,
		)
	default:
		w.reqLog.Printf("R: unimplemented server/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaServerResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaServerResult{})
	default:
		result.Append(nil, &somaServerResult{
			Server: q.Server,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaServerReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaServerWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
