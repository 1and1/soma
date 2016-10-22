package main

import (
	"database/sql"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type somaClusterRequest struct {
	action  string
	Cluster proto.Cluster
	reply   chan somaResult
}

type somaClusterResult struct {
	ResultError error
	Cluster     proto.Cluster
}

func (a *somaClusterResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Clusters = append(r.Clusters, somaClusterResult{ResultError: err})
	}
}

func (a *somaClusterResult) SomaAppendResult(r *somaResult) {
	r.Clusters = append(r.Clusters, *a)
}

/* Read Access
 */
type somaClusterReadHandler struct {
	input     chan somaClusterRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	mbnl_stmt *sql.Stmt
	ponc_stmt *sql.Stmt
	psvc_stmt *sql.Stmt
	psys_stmt *sql.Stmt
	pcst_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaClusterReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ClusterList:       r.list_stmt,
		stmt.ClusterShow:       r.show_stmt,
		stmt.ClusterMemberList: r.mbnl_stmt,
		stmt.ClusterOncProps:   r.ponc_stmt,
		stmt.ClusterSvcProps:   r.psvc_stmt,
		stmt.ClusterSysProps:   r.psys_stmt,
		stmt.ClusterCstProps:   r.pcst_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`cluster`, err, statement)
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

func (r *somaClusterReadHandler) process(q *somaClusterRequest) {
	var (
		clusterId, clusterName, bucketId, clusterState, teamId string
		mNodeId, mNodeName, instanceId, sourceInstanceId       string
		view, oncallId, oncallName, serviceName, customId      string
		systemProp, value, customProp                          string
		rows                                                   *sql.Rows
		err                                                    error
	)
	result := somaResult{}
	resC := proto.Cluster{}

	switch q.action {
	case "list":
		r.appLog.Printf("R: cluster/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&clusterId, &clusterName, &bucketId)
			result.Append(err, &somaClusterResult{
				Cluster: proto.Cluster{
					Id:       clusterId,
					Name:     clusterName,
					BucketId: bucketId,
				},
			})
		}
	case "show":
		r.appLog.Printf("R: cluster/show for %s", q.Cluster.Id)
		err = r.show_stmt.QueryRow(q.Cluster.Id).Scan(
			&clusterId,
			&bucketId,
			&clusterName,
			&clusterState,
			&teamId,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			goto dispatch
		}
		cluster := proto.Cluster{
			Id:          clusterId,
			Name:        clusterName,
			BucketId:    bucketId,
			ObjectState: clusterState,
			TeamId:      teamId,
		}
		cluster.Properties = &[]proto.Property{}

		// oncall properties
		rows, err = r.ponc_stmt.Query(q.Cluster.Id)
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
			*cluster.Properties = append(
				*cluster.Properties,
				proto.Property{
					Type:             `oncall`,
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
		rows, err = r.psvc_stmt.Query(q.Cluster.Id)
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
			*cluster.Properties = append(
				*cluster.Properties,
				proto.Property{
					Type:             `service`,
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
		rows, err = r.psys_stmt.Query(q.Cluster.Id)
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
			*cluster.Properties = append(
				*cluster.Properties,
				proto.Property{
					Type:             `system`,
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
		rows, err = r.pcst_stmt.Query(q.Cluster.Id)
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
			*cluster.Properties = append(
				*cluster.Properties,
				proto.Property{
					Type:             `custom`,
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

		result.Append(err, &somaClusterResult{
			Cluster: cluster,
		})
	case "member_list":
		r.appLog.Printf("R: cluster/memberlist for %s", q.Cluster.Id)
		rows, err = r.mbnl_stmt.Query(q.Cluster.Id)
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		resC.Id = q.Cluster.Id
		for rows.Next() {
			err := rows.Scan(&mNodeId, &mNodeName, &clusterName)
			if err == nil {
				resC.Name = clusterName
				*resC.Members = append(*resC.Members, proto.Node{
					Id:   mNodeId,
					Name: mNodeName,
				})
			}
		}

		result.Append(err, &somaClusterResult{
			Cluster: resC,
		})
	default:
		r.errLog.Printf("R: unimplemented cluster/%s", q.action)
		result.SetNotImplemented()
	}
dispatch:
	q.reply <- result
}

/* Ops Access
 */
func (r *somaClusterReadHandler) shutdownNow() {
	r.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
