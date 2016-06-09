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
	Node   proto.Node
	reply  chan somaResult
}

type somaNodeResult struct {
	ResultError error
	Node        proto.Node
}

func (a *somaNodeResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Nodes = append(r.Nodes, somaNodeResult{ResultError: err})
	}
}

func (a *somaNodeResult) SomaAppendResult(r *somaResult) {
	r.Nodes = append(r.Nodes, *a)
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
	conf_stmt *sql.Stmt
	sync_stmt *sql.Stmt
}

func (r *somaNodeReadHandler) run() {
	var err error

	if r.list_stmt, err = r.conn.Prepare(stmt.ListNodes); err != nil {
		log.Fatal("node/list: ", err)
	}
	defer r.list_stmt.Close()

	if r.show_stmt, err = r.conn.Prepare(stmt.ShowNodes); err != nil {
		log.Fatal("node/show: ", err)
	}
	defer r.show_stmt.Close()

	if r.conf_stmt, err = r.conn.Prepare(stmt.ShowConfigNodes); err != nil {
		log.Fatal("node/get-config: ", err)
	}
	defer r.conf_stmt.Close()

	if r.sync_stmt, err = r.conn.Prepare(stmt.SyncNodes); err != nil {
		log.Fatal("node/sync: ", err)
	}
	defer r.sync_stmt.Close()

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
	var nodeId, nodeName, nodeTeam, nodeServer, nodeState, bucketId, repositoryId string
	var nodeAsset int
	var nodeOnline, nodeDeleted bool
	var rows *sql.Rows
	var err error
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: node/list")
		rows, err = r.list_stmt.Query()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&nodeId, &nodeName)
			result.Append(err, &somaNodeResult{
				Node: proto.Node{
					Id:   nodeId,
					Name: nodeName,
				},
			})
		}
	case `sync`:
		log.Printf(`R: node/sync`)
		rows, err = r.sync_stmt.Query()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(
				&nodeId,
				&nodeAsset,
				&nodeName,
				&nodeTeam,
				&nodeServer,
				&nodeOnline,
				&nodeDeleted,
			)
			result.Append(err, &somaNodeResult{
				Node: proto.Node{
					Id:        nodeId,
					AssetId:   uint64(nodeAsset),
					Name:      nodeName,
					TeamId:    nodeTeam,
					ServerId:  nodeServer,
					IsOnline:  nodeOnline,
					IsDeleted: nodeDeleted,
				},
			})
		}
	case "show":
		log.Printf("R: node/show")
		err = r.show_stmt.QueryRow(q.Node.Id).Scan(
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
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		result.Append(err, &somaNodeResult{
			Node: proto.Node{
				Id:        nodeId,
				AssetId:   uint64(nodeAsset),
				Name:      nodeName,
				TeamId:    nodeTeam,
				ServerId:  nodeServer,
				State:     nodeState,
				IsOnline:  nodeOnline,
				IsDeleted: nodeDeleted,
			},
		})
	case "get_config":
		log.Printf("R: node/get_config")
		err = r.conf_stmt.QueryRow(q.Node.Id).Scan(
			&nodeId,
			&nodeName,
			&bucketId,
			&repositoryId,
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

		result.Append(err, &somaNodeResult{
			Node: proto.Node{
				Id:   nodeId,
				Name: nodeName,
				Config: &proto.NodeConfig{
					RepositoryId: repositoryId,
					BucketId:     bucketId,
				},
			},
		})
	default:
		result.SetNotImplemented()
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
SELECT $1::uuid, $2::numeric, $3::varchar, $4, $5, $6, $7, $8
WHERE NOT EXISTS (
	SELECT node_id
	FROM   soma.nodes
	WHERE  node_id = $1::uuid
	OR     node_asset_id = $2::numeric
	OR     (node_name = $3::varchar AND node_online));`)
	if err != nil {
		log.Fatal("node/add: ", err)
	}
	defer w.add_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
UPDATE soma.nodes
SET    node_deleted = 'yes'
WHERE  node_id = $1
AND    node_deleted = 'no';`)
	if err != nil {
		log.Fatal("node/delete: ", err)
	}
	defer w.del_stmt.Close()

	w.prg_stmt, err = w.conn.Prepare(`
DELETE FROM soma.nodes
WHERE       node_id = $1
AND         node_deleted;`)
	if err != nil {
		log.Fatal("node/purge: ", err)
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

func (w *somaNodeWriteHandler) process(q *somaNodeRequest) {
	var res sql.Result
	var err error
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: node/add for %s", q.Node.Name)
		id := uuid.NewV4()
		if q.Node.ServerId == "" {
			q.Node.ServerId = "00000000-0000-0000-0000-000000000000"
		}
		res, err = w.add_stmt.Exec(
			id.String(),
			q.Node.AssetId,
			q.Node.Name,
			q.Node.TeamId,
			q.Node.ServerId,
			q.Node.State,
			q.Node.IsOnline,
			false,
		)
		q.Node.Id = id.String()
	case "delete":
		log.Printf("R: node/delete for %s", q.Node.Id)
		res, err = w.del_stmt.Exec(
			q.Node.Id,
		)
		// TODO trigger undeployment
	case "purge":
		log.Printf("R: node/purge for %s", q.Node.Id)
		res, err = w.prg_stmt.Exec(
			q.Node.Id,
		)
	default:
		log.Printf("R: unimplemented node/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaNodeResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaNodeResult{})
	default:
		result.Append(nil, &somaNodeResult{
			Node: q.Node,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
