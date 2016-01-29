package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/satori/go.uuid"

)

type somaServerRequest struct {
	action string
	Server somaproto.ProtoServer
	reply  chan somaResult
}

type somaServerResult struct {
	ResultError error
	Server      somaproto.ProtoServer
}

/* Read
 */
type somaServerReadHandler struct {
	input     chan somaServerRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
}

func (r *somaServerReadHandler) run() {
	var err error

	log.Println("Prepare: server/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT server_id,
       server_name
FROM   inventory.servers
WHERE  server_online
AND    NOT server_deleted
AND    NOT server_id = '00000000-0000-0000-0000-000000000000';`)
	if err != nil {
		log.Fatal("server/list: ", err)
	}
	defer r.list_stmt.Close()

	log.Println("Prepare: server/show")
	r.show_stmt, err = r.conn.Prepare(`
SELECT server_id,
       server_asset_id,
	   server_datacenter_name,
	   server_datacenter_location,
	   server_name,
	   server_online,
	   server_deleted
FROM   inventory.servers
WHERE  server_id = $1;`)
	if err != nil {
		log.Fatal("server/show: ", err)
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

func (r *somaServerReadHandler) process(q *somaServerRequest) {
	var serverId, serverDc, serverDcLoc, serverName string
	var serverAsset int
	var serverOnline, serverDeleted bool
	var rows *sql.Rows
	var err error
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: server/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&serverId, &serverName)
			result.Append(err, somaServerResult{
				Server: somaproto.ProtoServer{
					Id:   serverId,
					Name: serverName,
				},
			})
		}
	case "show":
		log.Printf("R: server/show")
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
			if err.Error() == "sql: no rows in result set" {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		result.Append(err, somaServerResult{
			Server: somaproto.ProtoServer{
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

/* Write
 */
type somaServerWriteHandler struct {
	input    chan somaServerRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	prg_stmt *sql.Stmt
}

func (w *somaServerWriteHandler) run() {
	var err error

	log.Println("Prepare: server/add")
	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO inventory.servers (
	server_id,
	server_asset_id,
	server_datacenter_name,
	server_datacenter_location,
	server_name,
	server_online,
	server_deleted)
SELECT	$1::uuid, $2::numeric, $3, $4, $5, $6, $7
WHERE	NOT EXISTS(
	SELECT server_id
	FROM   inventory.servers
	WHERE  server_id = $1::uuid
	OR     server_asset_id = $2::numeric);`)
	if err != nil {
		log.Fatal("server/add: ", err)
	}
	defer w.add_stmt.Close()

	log.Println("Prepare: server/del")
	w.del_stmt, err = w.conn.Prepare(`
UPDATE inventory.servers
SET    server_online = 'no'
WHERE  server_id = $1::uuid
AND    server_online
AND    server_id != '00000000-0000-0000-0000-000000000000';`)
	if err != nil {
		log.Fatal("server/del: ", err)
	}
	defer w.del_stmt.Close()

	log.Println("Prepare: server/prg")
	w.prg_stmt, err = w.conn.Prepare(`
DELETE FROM inventory.servers
WHERE  server_id = $1::uuid
AND    server_deleted
AND    server_id != '00000000-0000-0000-0000-000000000000';`)
	if err != nil {
		log.Fatal("server/prg: ", err)
	}
	defer w.prg_stmt.Close()

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
	var res sql.Result
	var err error
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: server/add for %s", q.Server.Name)
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
	case "delete":
		log.Printf("R: server/delete for %s", q.Server.Id)
		res, err = w.del_stmt.Exec(
			q.Server.Id,
		)
	case "purge":
		log.Printf("R: server/purge for %s", q.Server.Id)
		res, err = w.del_stmt.Exec(
			q.Server.Id,
		)
	case "insert-null":
		log.Printf("R: server/insert-null")
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
		log.Printf("R: unimplemented server/%s", q.action)
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
		result.Append(errors.New("No rows affected"), somaServerResult{})
	case rowCnt >= 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			somaServerResult{})
	default:
		result.Append(nil, somaServerResult{
			Server: q.Server,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
