package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/satori/go.uuid"
)

type somaNodeRequest struct {
	action string
	node   somaproto.ProtoNode
	reply  chan []somaNodeResult
}

type somaNodeResult struct {
	rErr error
	lErr error
	node somaproto.ProtoNode
}

/* Read Access
 *
 */
type somaNodeReadHandler struct {
	input     chan somaNodeRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
}

func (r *somaNodeReadHandler) run() {
	var err error

	log.Println("Prepare: node/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT node_id,
       node_name
FROM   soma.nodes
WHERE  node_online;`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Prepare: node/show")
	r.show_stmt, err = r.conn.Prepare(`
SELECT node_id,
       node_asset_id,
       node_name,
       organizational_team_id,
       server_id,
       object_state,
       node_online,
       node_deleted,
FROM   soma.nodes
WHERE  node_id = $1;`)
	if err != nil {
		log.Fatal(err)
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

func (r *somaNodeReadHandler) process(q *somaNodeRequest) {
	var nodeId, nodeName, nodeTeam, nodeServer, nodeState string
	var nodeAsset int
	var nodeOnline, nodeDeleted bool
	var rows *sql.Rows
	var err error
	result := make([]somaNodeResult, 0)

	switch q.action {
	case "list":
		log.Printf("R: node/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		result = append(result, somaNodeResult{
			rErr: err,
		})
		q.reply <- result
		return

		for rows.Next() {
			err := rows.Scan(&nodeId, &nodeName)
			if err != nil {
				result = append(result, somaNodeResult{
					lErr: err,
				})
				err = nil
				continue
			}
			result = append(result, somaNodeResult{
				node: somaproto.ProtoNode{
					Id:   nodeId,
					Name: nodeName,
				},
			})
		}
	case "show":
		log.Printf("R: node/show")
		err = r.show_stmt.QueryRow(q.node.Id).Scan(
			&nodeId,
			&nodeAsset,
			&nodeName,
			&nodeTeam,
			&nodeServer,
			&nodeState,
			&nodeOnline,
			&nodeDeleted,
		)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				result = append(result, somaNodeResult{
					rErr: err,
				})
			}
			q.reply <- result
			return
		}

		result = append(result, somaNodeResult{
			node: somaproto.ProtoNode{
				Id:        nodeId,
				AssetId:   uint64(nodeAsset),
				Name:      nodeName,
				Team:      nodeTeam,
				Server:    nodeServer,
				State:     nodeState,
				IsOnline:  nodeOnline,
				IsDeleted: nodeDeleted,
			},
		})
	default:
		result = append(result, somaNodeResult{
			rErr: errors.New("not implemented"),
		})
	}
	q.reply <- result
}

/* Write Access
 *
 */
type somaNodeWriteHandler struct {
	input    chan somaNodeRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	prg_stmt *sql.Stmt
}

func (w *somaNodeWriteHandler) run() {
	var err error

	log.Println("Prepare: node/add")
	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.nodes (
	node_id,
	node_asset_id,
	node_name,
	organizational_team_id,
	server_id,
	object_state,
	node_online,
    node_deleted)
SELECT $1, $2, $3, $4, $5, $6, $7, $8
WHERE NOT EXISTS (
	SELECT node_id
	FROM   soma.nodes
	WHERE  node_id = $9
	OR     node_asset_id = $10
    OR     (node_name = $11 AND node_online));`)
	if err != nil {
		log.Fatal("node/add: ", err)
	}
	defer w.add_stmt.Close()

	log.Println("Prepare: node/del")
	w.del_stmt, err = w.conn.Prepare(`
UPDATE soma.nodes
SET    node_deleted = yes
WHERE  node_id = $1
AND    node_deleted = no;`)
	if err != nil {
		log.Fatal("node/del: ", err)
	}
	defer w.del_stmt.Close()

	log.Println("Prepare: node/prg")
	w.prg_stmt, err = w.conn.Prepare(`
DELETE FROM soma.nodes
WHERE       node_id = $1
AND         node_deleted;`)
	if err != nil {
		log.Fatal("node/prg: ", err)
	}
	defer w.prg_stmt.Close()

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

func (w *somaNodeWriteHandler) process(q *somaNodeRequest) {
	var res sql.Result
	var err error
	result := make([]somaNodeResult, 0)

	switch q.action {
	case "add":
		log.Printf("R: node/add for %s", q.node.Name)
		id := uuid.NewV4()
		res, err = w.add_stmt.Exec(
			id.String(),
			q.node.AssetId,
			q.node.Name,
			q.node.Team,
			q.node.Server,
			q.node.State,
			q.node.IsOnline,
			q.node.IsDeleted,
			id.String(),
			q.node.AssetId,
			q.node.Name,
		)
		q.node.Id = id.String()
	case "delete":
		log.Printf("R: node/delete for %s", q.node.Id)
		res, err = w.del_stmt.Exec(
			q.node.Id,
		)
		// TODO trigger undeployment
	case "purge":
		log.Printf("R: node/purge for %s", q.node.Id)
		res, err = w.prg_stmt.Exec(
			q.node.Id,
		)
	default:
		log.Printf("R: unimplemented node/%s", q.action)
		result = append(result, somaNodeResult{
			rErr: errors.New("not implemented"),
		})
		q.reply <- result
		return
	}
	if err != nil {
		result = append(result, somaNodeResult{
			rErr: err,
		})
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		result = append(result, somaNodeResult{
			lErr: errors.New("No rows affected"),
		})
	case rowCnt > 1:
		result = append(result, somaNodeResult{
			lErr: fmt.Errorf("Too many rows affected: %d", rowCnt),
		})
	default:
		result = append(result, somaNodeResult{
			node: q.node,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
