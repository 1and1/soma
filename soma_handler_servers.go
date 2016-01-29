package main

import (
	"database/sql"
	"log"

)

type somaServerRequest struct {
	action string
	server somaproto.ProtoServer
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
		if err != nil {
			result.SetRequestError(err)
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
		err = r.show_stmt.QueryRow(q.server.Id).Scan(
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
				result.SetRequestError(err)
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
