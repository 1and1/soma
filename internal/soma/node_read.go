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
	"log"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
)

// NodeRead handles read requests for nodes
type NodeRead struct {
	Input           chan msg.Request
	Shutdown        chan bool
	conn            *sql.DB
	stmtList        *sql.Stmt
	stmtShow        *sql.Stmt
	stmtShowConfig  *sql.Stmt
	stmtSync        *sql.Stmt
	stmtPropOncall  *sql.Stmt
	stmtPropService *sql.Stmt
	stmtPropSystem  *sql.Stmt
	stmtPropCustom  *sql.Stmt
	appLog          *log.Logger
	reqLog          *log.Logger
	errLog          *log.Logger
}

// run is the event loop for NodeRead
func (r *NodeRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.NodeList:       r.stmtList,
		stmt.NodeShow:       r.stmtShow,
		stmt.NodeShowConfig: r.stmtShowConfig,
		stmt.NodeSync:       r.stmtSync,
		stmt.NodeOncProps:   r.stmtPropOncall,
		stmt.NodeSvcProps:   r.stmtPropService,
		stmt.NodeSysProps:   r.stmtPropSystem,
		stmt.NodeCstProps:   r.stmtPropCustom,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`node`, err, stmt.Name(statement))
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
func (r *NodeRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case `list`:
		r.list(q, &result)
	case `sync`:
		r.sync(q, &result)
	case `show`:
		r.show(q, &result)
	case `show-config`:
		r.showConfig(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all nodes
func (r *NodeRead) list(q *msg.Request, mr *msg.Result) {
	var (
		rows             *sql.Rows
		err              error
		nodeID, nodeName string
	)

	r.reqLog.Printf("R: node/list")
	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&nodeID, &nodeName); err != nil {
			rows.Close()
			mr.ServerError(err)
			mr.Clear(q.Section)
			return
		}
		mr.Node = append(mr.Node, proto.Node{
			Id:   nodeID,
			Name: nodeName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

// sync returns all nodes with all details attached
func (r *NodeRead) sync(q *msg.Request, mr *msg.Result) {
	var (
		rows                                   *sql.Rows
		err                                    error
		nodeID, nodeName, nodeTeam, nodeServer string
		nodeAsset                              int
		nodeOnline, nodeDeleted                bool
	)

	r.reqLog.Printf(`R: node/sync`)
	if rows, err = r.stmtSync.Query(); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&nodeID,
			&nodeAsset,
			&nodeName,
			&nodeTeam,
			&nodeServer,
			&nodeOnline,
			&nodeDeleted,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			mr.Clear(q.Section)
			return
		}
		mr.Node = append(mr.Node, proto.Node{
			Id:        nodeID,
			AssetId:   uint64(nodeAsset),
			Name:      nodeName,
			TeamId:    nodeTeam,
			ServerId:  nodeServer,
			IsOnline:  nodeOnline,
			IsDeleted: nodeDeleted,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		mr.Clear(q.Section)
		return
	}
	mr.OK()
}

// show returns the details for a specific node
func (r *NodeRead) show(q *msg.Request, mr *msg.Result) {
	var (
		rows                                    *sql.Rows
		err                                     error
		nodeID, nodeName, nodeTeam, nodeServer  string
		repositoryID, bucketID, instanceID      string
		view, oncallID, oncallName, serviceName string
		systemProp, sourceInstanceID            string
		nodeState, value                        string
		nodeOnline, nodeDeleted                 bool
		nodeAsset                               int
		node                                    proto.Node
		tx                                      *sql.Tx
		checkConfigs                            *[]proto.CheckConfig
	)

	r.reqLog.Printf("R: node/show")
	if err = r.stmtShow.QueryRow(
		q.Node.Id,
	).Scan(
		&nodeID,
		&nodeAsset,
		&nodeName,
		&nodeTeam,
		&nodeServer,
		&nodeState,
		&nodeOnline,
		&nodeDeleted,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		mr.Clear(q.Section)
		return
	} else if err != nil {
		mr.ServerError(err)
		mr.Clear(q.Section)
		return
	}
	node = proto.Node{
		Id:        nodeID,
		AssetId:   uint64(nodeAsset),
		Name:      nodeName,
		TeamId:    nodeTeam,
		ServerId:  nodeServer,
		State:     nodeState,
		IsOnline:  nodeOnline,
		IsDeleted: nodeDeleted,
	}

	// add configuration data
	if err = r.stmtShowConfig.QueryRow(
		q.Node.Id,
	).Scan(
		&nodeID,
		&nodeName,
		&bucketID,
		&repositoryID,
	); err == sql.ErrNoRows {
		// sql.ErrNoRows means the node is unassigned, which is
		// valid and not an error. But an unconfigured node can
		// not have properties or checks, which means the request
		// is done.
		mr.OK()
		return
	} else if err != nil {
		mr.ServerError(err)
		mr.Clear(q.Section)
		return
	}
	// node is assigned in this codepath
	node.Config = &proto.NodeConfig{
		RepositoryId: repositoryID,
		BucketId:     bucketID,
	}

	// fetch node properties
	node.Properties = &[]proto.Property{}

	// oncall properties
	if rows, err = r.stmtPropOncall.Query(
		q.Node.Id,
	); err != nil {
		mr.ServerError(err)
		mr.Clear(q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&sourceInstanceID,
			&view,
			&oncallID,
			&oncallName,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			mr.Clear(q.Section)
			return
		}
		*node.Properties = append(*node.Properties, proto.Property{
			Type:             `oncall`,
			RepositoryId:     repositoryID,
			BucketId:         bucketID,
			InstanceId:       instanceID,
			SourceInstanceId: sourceInstanceID,
			View:             view,
			Oncall: &proto.PropertyOncall{
				Id:   oncallID,
				Name: oncallName,
			},
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		mr.Clear(q.Section)
		return
	}

	// service properties
	if rows, err = r.stmtPropService.Query(
		q.Node.Id,
	); err != nil {
		mr.ServerError(err)
		mr.Clear(q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&sourceInstanceID,
			&view,
			&serviceName,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			mr.Clear(q.Section)
			return
		}
		*node.Properties = append(*node.Properties, proto.Property{
			Type:             `service`,
			RepositoryId:     repositoryID,
			BucketId:         bucketID,
			InstanceId:       instanceID,
			SourceInstanceId: sourceInstanceID,
			View:             view,
			Service: &proto.PropertyService{
				Name: serviceName,
			},
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		mr.Clear(q.Section)
		return
	}

	// system properties
	if rows, err = r.stmtPropSystem.Query(
		q.Node.Id,
	); err != nil {
		mr.ServerError(err)
		mr.Clear(q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&sourceInstanceID,
			&view,
			&systemProp,
			&value,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			mr.Clear(q.Section)
			return
		}
		*node.Properties = append(*node.Properties, proto.Property{
			Type:             `system`,
			RepositoryId:     repositoryID,
			BucketId:         bucketID,
			InstanceId:       instanceID,
			SourceInstanceId: sourceInstanceID,
			View:             view,
			System: &proto.PropertySystem{
				Name:  systemProp,
				Value: value,
			},
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		mr.Clear(q.Section)
		return
	}

	// custom properties
	if err = r.customProperties(&node); err != nil {
		mr.ServerError(err)
		mr.Clear(q.Section)
		return
	}

	// add check configuration and instance information
	if tx, err = r.conn.Begin(); err != nil {
		mr.ServerError(err)
		mr.Clear(q.Section)
		return
	}
	if checkConfigs, err = exportCheckConfigObjectTX(
		tx,
		q.Node.Id,
	); err != nil {
		tx.Rollback()
		mr.ServerError(err)
		mr.Clear(q.Section)
		return
	}
	if checkConfigs != nil && len(*checkConfigs) > 0 {
		node.Details = &proto.Details{
			CheckConfigs: checkConfigs,
		}
	}

	mr.Node = append(mr.Node, node)
	mr.OK()
}

// showConfig returns the repository configuration of the node
func (r *NodeRead) showConfig(q *msg.Request, mr *msg.Result) {
	var (
		err                                      error
		nodeID, nodeName, repositoryID, bucketID string
	)
	if err = r.stmtShowConfig.QueryRow(
		q.Node.Id,
	).Scan(
		&nodeID,
		&nodeName,
		&bucketID,
		&repositoryID,
	); err == sql.ErrNoRows {
		// TODO need a better way to transport 'unassigned'
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}
	mr.Node = append(mr.Node, proto.Node{
		Id:   nodeID,
		Name: nodeName,
		Config: &proto.NodeConfig{
			RepositoryId: repositoryID,
			BucketId:     bucketID,
		},
	})
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
