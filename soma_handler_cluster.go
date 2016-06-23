package main

import (
	"database/sql"
	"log"

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
}

func (r *somaClusterReadHandler) run() {
	var err error

	r.list_stmt, err = r.conn.Prepare(`
SELECT cluster_id,
       cluster_name
FROM soma.clusters;`)
	if err != nil {
		log.Fatal("cluster/list: ", err)
	}
	defer r.list_stmt.Close()

	r.show_stmt, err = r.conn.Prepare(`
SELECT cluster_id,
       bucket_id,
	   cluster_name,
	   object_state,
	   organizational_team_id
FROM   soma.clusters
WHERE  cluster_id = $1::uuid;`)
	if err != nil {
		log.Fatal("cluster/show: ", err)
	}
	defer r.show_stmt.Close()

	r.mbnl_stmt, err = r.conn.Prepare(`
SELECT sn.node_id,
       sn.node_name,
	   sc.cluster_name
FROM   soma.cluster_membership scm
JOIN   soma.nodes sn
ON     scm.node_id = sn.node_id
JOIN   soma.clusters sc
ON     scm.cluster_id = sc.cluster_id
WHERE  scm.cluster_id = $1::uuid;`)
	if err != nil {
		log.Fatal("cluster/memberlist-node: ", err)
	}
	defer r.mbnl_stmt.Close()

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
		mNodeId, mNodeName                                     string
		rows                                                   *sql.Rows
		err                                                    error
	)
	result := somaResult{}
	resC := proto.Cluster{}

	switch q.action {
	case "list":
		log.Printf("R: cluster/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&clusterId, &clusterName)
			result.Append(err, &somaClusterResult{
				Cluster: proto.Cluster{
					Id:   clusterId,
					Name: clusterName,
				},
			})
		}
	case "show":
		log.Printf("R: cluster/show for %s", q.Cluster.Id)
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
			q.reply <- result
			return
		}

		result.Append(err, &somaClusterResult{
			Cluster: proto.Cluster{
				Id:          clusterId,
				Name:        clusterName,
				BucketId:    bucketId,
				ObjectState: clusterState,
				TeamId:      teamId,
			},
		})
	case "member_list":
		log.Printf("R: cluster/memberlist for %s", q.Cluster.Id)
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
		result.SetNotImplemented()
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
