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

type somaNodeRequest struct {
	action string
	user   string
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
	ponc_stmt *sql.Stmt
	psvc_stmt *sql.Stmt
	psys_stmt *sql.Stmt
	pcst_stmt *sql.Stmt
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

	if r.ponc_stmt, err = r.conn.Prepare(stmt.NodeOncProps); err != nil {
		log.Fatal(`node/property-oncall: `, err)
	}
	defer r.ponc_stmt.Close()

	if r.psvc_stmt, err = r.conn.Prepare(stmt.NodeSvcProps); err != nil {
		log.Fatal(`node/property-service: `, err)
	}
	defer r.psvc_stmt.Close()

	if r.psys_stmt, err = r.conn.Prepare(stmt.NodeSysProps); err != nil {
		log.Fatal(`node/property-system: `, err)
	}
	defer r.psys_stmt.Close()

	if r.pcst_stmt, err = r.conn.Prepare(stmt.NodeCstProps); err != nil {
		log.Fatal(`node/property-custom: `, err)
	}
	defer r.pcst_stmt.Close()

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
	var (
		nodeId, nodeName, nodeTeam, nodeServer, nodeState    string
		bucketId, repositoryId, instanceId, sourceInstanceId string
		view, oncallId, oncallName, serviceName, customId    string
		systemProp, value, customProp                        string
		nodeAsset                                            int
		nodeOnline, nodeDeleted                              bool
		rows                                                 *sql.Rows
		err                                                  error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: node/list")
		rows, err = r.list_stmt.Query()
		if result.SetRequestError(err) {
			goto dispatch
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
			goto dispatch
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
			goto dispatch
		}
		node := proto.Node{
			Id:        nodeId,
			AssetId:   uint64(nodeAsset),
			Name:      nodeName,
			TeamId:    nodeTeam,
			ServerId:  nodeServer,
			State:     nodeState,
			IsOnline:  nodeOnline,
			IsDeleted: nodeDeleted,
		}

		// add configuration data
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
			goto dispatch
		}
		node.Config = &proto.NodeConfig{
			RepositoryId: repositoryId,
			BucketId:     bucketId,
		}
		node.Properties = &[]proto.Property{}

		// oncall properties
		rows, err = r.ponc_stmt.Query(q.Node.Id)
		if result.SetRequestError(err) {
			goto dispatch
		}
		for rows.Next() {
			if err := rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&oncallId,
				&oncallName,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*node.Properties = append(
				*node.Properties,
				proto.Property{
					Type:             `oncall`,
					RepositoryId:     repositoryId,
					BucketId:         bucketId,
					InstanceId:       instanceId,
					SourceInstanceId: sourceInstanceId,
					View:             view,
					Oncall: &proto.PropertyOncall{
						Id:   oncallId,
						Name: oncallName,
					},
				},
			)
		}
		if err = rows.Err(); err != nil {
			result.SetRequestError(err)
			goto dispatch
		}

		// service properties
		rows, err = r.psvc_stmt.Query(q.Node.Id)
		for rows.Next() {
			if err := rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&serviceName,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*node.Properties = append(
				*node.Properties,
				proto.Property{
					Type:             `service`,
					RepositoryId:     repositoryId,
					BucketId:         bucketId,
					InstanceId:       instanceId,
					SourceInstanceId: sourceInstanceId,
					View:             view,
					Service: &proto.PropertyService{
						Name: serviceName,
					},
				},
			)
		}
		if err = rows.Err(); err != nil {
			result.SetRequestError(err)
			goto dispatch
		}

		// system properties
		rows, err = r.psys_stmt.Query(q.Node.Id)
		for rows.Next() {
			if err := rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&systemProp,
				&value,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*node.Properties = append(
				*node.Properties,
				proto.Property{
					Type:             `system`,
					RepositoryId:     repositoryId,
					BucketId:         bucketId,
					InstanceId:       instanceId,
					SourceInstanceId: sourceInstanceId,
					View:             view,
					System: &proto.PropertySystem{
						Name:  systemProp,
						Value: value,
					},
				},
			)
		}
		if err = rows.Err(); err != nil {
			result.SetRequestError(err)
			goto dispatch
		}

		// custom properties
		rows, err = r.pcst_stmt.Query(q.Node.Id)
		for rows.Next() {
			if err := rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&customId,
				&value,
				&customProp,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*node.Properties = append(
				*node.Properties,
				proto.Property{
					Type:             `custom`,
					RepositoryId:     repositoryId,
					BucketId:         bucketId,
					InstanceId:       instanceId,
					SourceInstanceId: sourceInstanceId,
					View:             view,
					Custom: &proto.PropertyCustom{
						Id:    customId,
						Name:  customProp,
						Value: value,
					},
				},
			)
		}
		if err = rows.Err(); err != nil {
			result.SetRequestError(err)
			goto dispatch
		}

		result.Append(err, &somaNodeResult{
			Node: node,
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
			goto dispatch
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
dispatch:
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
	upd_stmt *sql.Stmt
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
            node_deleted,
            created_by)
SELECT $1::uuid,
       $2::numeric,
       $3::varchar,
       $4,
       $5,
       $6,
       $7,
       $8,
       user_id
FROM   inventory.users iu
WHERE  iu.user_uid = $9::varchar
AND    NOT EXISTS (
         SELECT node_id
         FROM   soma.nodes
         WHERE  node_id = $1::uuid
         OR     node_asset_id = $2::numeric
         OR     (node_name = $3::varchar AND node_online)
       );`)
	if err != nil {
		log.Fatal("node/add: ", err)
	}
	defer w.add_stmt.Close()

	if w.upd_stmt, err = w.conn.Prepare(`
UPDATE soma.nodes
SET    node_asset_id = $1::numeric,
       node_name = $2::varchar,
       organizational_team_id = $3::uuid,
       server_id = $4::uuid,
       node_online = $5::boolean,
       node_deleted = $6::boolean
WHERE  node_id = $7::uuid;`); err != nil {
		log.Fatal(`node/update: `, err)
	}
	defer w.upd_stmt.Close()

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
			q.user,
		)
		q.Node.Id = id.String()
	case `update`:
		log.Printf("R: node/update for %s", q.Node.Id)
		res, err = w.upd_stmt.Exec(
			q.Node.AssetId,
			q.Node.Name,
			q.Node.TeamId,
			q.Node.ServerId,
			q.Node.IsOnline,
			q.Node.IsDeleted,
			q.Node.Id,
		)
		// TODO what has to be done for this undeployment?
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

/* Ops Access
 */
func (r *somaNodeReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaNodeWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
