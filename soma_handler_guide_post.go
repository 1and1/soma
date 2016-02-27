package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"


	"github.com/satori/go.uuid"
)

type guidePost struct {
	input     chan treeRequest
	shutdown  chan bool
	conn      *sql.DB
	jbsv_stmt *sql.Stmt
	repo_stmt *sql.Stmt
	name_stmt *sql.Stmt
	node_stmt *sql.Stmt
	serv_stmt *sql.Stmt
	attr_stmt *sql.Stmt
}

func (g *guidePost) run() {
	var err error

	log.Println("Prepare: guide/job-save")
	g.jbsv_stmt, err = g.conn.Prepare(`
INSERT INTO soma.jobs (
	job_id,
	job_status,
	job_result,
	job_type,
	repository_id,
	user_id,
	organizational_team_id,
	job)
SELECT	$1::uuid,
		$2::varchar,
		$3::varchar,
		$4::varchar,
		$5::uuid,
		$6::uuid,
		$7::uuid,
		$8::jsonb;`)
	if err != nil {
		log.Fatal("guide/job-save: ", err)
	}
	defer g.jbsv_stmt.Close()

	log.Println("Prepare: guide/repo-by-bucket")
	g.repo_stmt, err = g.conn.Prepare(`
SELECT	sb.repository_id,
		sr.repository_name
FROM	soma.buckets sb
JOIN    soma.repositories sr
ON		sb.repository_id = sr.repository_id
WHERE	sb.bucket_id = $1::uuid;`)
	if err != nil {
		log.Fatal("guide/repo-by-bucket: ", err)
	}
	defer g.repo_stmt.Close()

	log.Println("Prepare: guide/load-node-details")
	g.node_stmt, err = g.conn.Prepare(`
SELECT    sn.node_asset_id,
	      sn.node_name,
	      sn.organizational_team_id,
	      sn.server_id,
	      sn.node_online,
	      sn.node_deleted
FROM      soma.nodes sn
LEFT JOIN soma.node_bucket_assignment snba
ON        sn.node_id = snba.node_id
WHERE     sn.node_online = 'yes'
AND       sn.node_deleted = 'false'
AND       snba.node_id IS NULL
AND       sn.node_id = $1::uuid;`)
	if err != nil {
		log.Fatal("guide/load-node-details: ", err)
	}
	defer g.node_stmt.Close()

	log.Println("Prepare: guide/repo-by-id")
	g.name_stmt, err = g.conn.Prepare(`
SELECT repository_name
FROM   soma.repositories
WHERE  repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("guide/repo-by-id: ", err)
	}
	defer g.name_stmt.Close()

	log.Println("Prepare: guide/service-lookup")
	g.serv_stmt, err = g.conn.Prepare(`
SELECT stsp.service_property
FROM   soma.repositories sr
JOIN   soma.team_service_properties stsp
ON     sr.organizational_team_id = stsp.organizational_team_id
WHERE  sr.repository_id = $1::uuid
AND    stsp.service_property = $2::varchar
AND    sr.organizational_team_id = $3::uuid;`)
	if err != nil {
		log.Fatal("guide/service-lookup: ", err)
	}
	defer g.serv_stmt.Close()

	log.Println("Prepare: guide/populate-service-attributes")
	g.attr_stmt, err = g.conn.Prepare(`
SELECT stspv.service_property_attribute,
       stspv.value
FROM   soma.repositories sr
JOIN   soma.team_service_properties stsp
ON     sr.organizational_team_id = stsp.organizational_team_id
JOIN   soma.team_service_property_values stspv
ON     stsp.organizational_team_id = stspv.organizational_team_id
AND    stsp.service_property = stspv.service_property
WHERE  sr.repository_id = $1::uuid
AND    stsp.service_property = $2::varchar
AND    sr.organizational_team_id = $3::uuid;`)
	if err != nil {
		log.Fatal("guide/populate-service-attributes: ", err)
	}
	defer g.attr_stmt.Close()

runloop:
	for {
		select {
		case <-g.shutdown:
			break runloop
		case req := <-g.input:
			g.process(&req)
		}
	}
}

func (g *guidePost) process(q *treeRequest) {
	var (
		res                                sql.Result
		err                                error
		j                                  []byte
		repoId, repoName, keeper, bucketId string
		ndName, ndTeam, ndServer           string
		ndAsset                            int64
		ndOnline, ndDeleted                bool
	)
	result := somaResult{}

	switch q.Action {
	case "add_system_property_to_repository":
		fallthrough
	case "add_custom_property_to_repository":
		fallthrough
	case "add_oncall_property_to_repository":
		fallthrough
	case "add_service_property_to_repository":
		repoId = q.Repository.Repository.Id

	case "create_bucket":
		repoId = q.Bucket.Bucket.Repository
	case "add_system_property_to_bucket":
		fallthrough
	case "add_custom_property_to_bucket":
		fallthrough
	case "add_oncall_property_to_bucket":
		fallthrough
	case "add_service_property_to_bucket":
		bucketId = q.Bucket.Bucket.Id

	case "create_group":
		fallthrough
	case "add_group_to_group":
		fallthrough
	case "add_cluster_to_group":
		fallthrough
	case "add_node_to_group":
		fallthrough
	case "add_system_property_to_group":
		fallthrough
	case "add_custom_property_to_group":
		fallthrough
	case "add_oncall_property_to_group":
		fallthrough
	case "add_service_property_to_group":
		bucketId = q.Group.Group.BucketId

	case "add_node_to_cluster":
		if q.Cluster.Cluster.BucketId != q.Cluster.Cluster.Members[0].Config.BucketId {
			panic("This should not happen.")
		}
		fallthrough
	case "create_cluster":
		fallthrough
	case "add_system_property_to_cluster":
		fallthrough
	case "add_custom_property_to_cluster":
		fallthrough
	case "add_oncall_property_to_cluster":
		fallthrough
	case "add_service_property_to_cluster":
		bucketId = q.Cluster.Cluster.BucketId

	case "assign_node":
		if err = g.node_stmt.QueryRow(q.Node.Node.Id).Scan(
			&ndAsset,
			&ndName,
			&ndTeam,
			&ndServer,
			&ndOnline,
			&ndDeleted,
		); err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		q.Node.Node.AssetId = uint64(ndAsset)
		q.Node.Node.Name = ndName
		q.Node.Node.Team = ndTeam
		q.Node.Node.Server = ndServer
		q.Node.Node.IsOnline = ndOnline
		q.Node.Node.IsDeleted = ndDeleted
		fallthrough
	case "add_system_property_to_node":
		fallthrough
	case "add_custom_property_to_node":
		fallthrough
	case "add_oncall_property_to_node":
		fallthrough
	case "add_service_property_to_node":
		repoId = q.Node.Node.Config.RepositoryId
		bucketId = q.Node.Node.Config.BucketId

	default:
		log.Printf("R: unimplemented guidepost/%s", q.Action)
		result.SetNotImplemented()
		q.reply <- result
		return
	}

	// lookup repository by bucket
	if bucketId != "" {
		if err = g.repo_stmt.QueryRow(bucketId).Scan(&repoId, &repoName); err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}
	}

	// lookup repository name
	if repoName == "" && repoId != "" {
		if err = g.name_stmt.QueryRow(repoId).Scan(&repoName); err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}
	}

	// XXX
	if repoName == "" {
		panic(`Have no repository name`)
	}

	// check we have a treekeeper for that repository
	keeper = fmt.Sprintf("repository_%s", repoName)
	if _, ok := handlerMap[keeper].(*treeKeeper); !ok {
		_ = result.SetRequestError(
			fmt.Errorf("No handler for repository %s registered.\n", repoName),
		)
		q.reply <- result
		return
	}

	// check the treekeeper has finished loading
	handler := handlerMap[keeper].(*treeKeeper)
	if !handler.isReady() {
		_ = result.SetRequestError( // TODO should be 503/ServiceUnavailable
			fmt.Errorf("Repository %s not fully loaded yet.\n", repoName),
		)
		q.reply <- result
		return
	}

	// check the treekeeper has not encountered a broken tree
	if handler.isBroken() {
		_ = result.SetRequestError(
			fmt.Errorf("Repository %s is broken.\n", repoName),
		)
		q.reply <- result
		return
	}

	// load authoritative copy of the service attributes from the
	// database. Replaces whatever the client sent in.
	if strings.Contains(q.Action, "add_service_property_to_") {
		var service, attr, val, svName, svTeam string
		var rows *sql.Rows
		attrs := []somaproto.TreeServiceAttribute{}

		switch q.RequestType {
		case "repository":
			svName = (*q.Repository.Repository.Properties)[0].Service.Name
			svTeam = (*q.Repository.Repository.Properties)[0].Service.TeamId
		case "bucket":
			svName = (*q.Bucket.Bucket.Properties)[0].Service.Name
			svTeam = (*q.Bucket.Bucket.Properties)[0].Service.TeamId
		case "group":
			svName = (*q.Group.Group.Properties)[0].Service.Name
			svTeam = (*q.Group.Group.Properties)[0].Service.TeamId
		case "cluster":
			svName = (*q.Cluster.Cluster.Properties)[0].Service.Name
			svTeam = (*q.Cluster.Cluster.Properties)[0].Service.TeamId
		case "node":
			svName = (*q.Node.Node.Properties)[0].Service.Name
			svTeam = (*q.Node.Node.Properties)[0].Service.TeamId
		}

		if err = g.serv_stmt.QueryRow(repoId, svName, svTeam).Scan(&service); err != nil {
			goto inputabort
		}

		if rows, err = g.attr_stmt.Query(repoId, svName, svTeam); err != nil {
			goto inputabort
		}
		defer rows.Close()

	attrloop:
		for rows.Next() {
			if err = rows.Scan(&attr, &val); err != nil {
				break attrloop
			}
			attrs = append(attrs, somaproto.TreeServiceAttribute{
				Attribute: attr,
				Value:     val,
			})
		}

	inputabort:
		if err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}
		// not aborted: set the loaded attributes
		switch q.RequestType {
		case "repository":
			(*q.Repository.Repository.Properties)[0].Service.Attributes = attrs
		case "bucket":
			(*q.Bucket.Bucket.Properties)[0].Service.Attributes = attrs
		case "group":
			(*q.Group.Group.Properties)[0].Service.Attributes = attrs
		case "cluster":
			(*q.Cluster.Cluster.Properties)[0].Service.Attributes = attrs
		case "node":
			(*q.Node.Node.Properties)[0].Service.Attributes = attrs
		}
	}

	// store job in database
	log.Printf("R: jobsave/%s", q.Action)
	q.JobId = uuid.NewV4()
	j, _ = json.Marshal(q)
	res, err = g.jbsv_stmt.Exec(
		q.JobId.String(),
		"queued",
		"pending",
		q.Action,
		repoId,
		"00000000-0000-0000-0000-000000000000", // XXX user uuid
		"00000000-0000-0000-0000-000000000000", // XXX team uuid
		string(j),
	)
	if result.SetRequestError(err) {
		q.reply <- result
		return
	}
	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		err = errors.New("No rows affected")
	case rowCnt > 1:
		err = fmt.Errorf("Too many rows affected: %d", rowCnt)
	case rowCnt < 0:
		err = fmt.Errorf("Space/Time Continuum broke, rows affected: %d", rowCnt)
	}
	if err != nil {
		switch q.RequestType {
		case "repository":
			result.Append(err, &somaRepositoryResult{})
		case "bucket":
			result.Append(err, &somaBucketResult{})
		case "group":
			result.Append(err, &somaGroupResult{})
		case "cluster":
			result.Append(err, &somaClusterResult{})
		case "node":
			result.Append(err, &somaNodeResult{})
		}
		q.reply <- result
		return
	}

	handler.input <- *q
	result.JobId = q.JobId.String()

	switch q.RequestType {
	case "repository":
		result.Append(nil, &somaRepositoryResult{
			Repository: q.Repository.Repository,
		})
	case "bucket":
		result.Append(nil, &somaBucketResult{
			Bucket: q.Bucket.Bucket,
		})
	case "group":
		result.Append(nil, &somaGroupResult{
			Group: q.Group.Group,
		})
	case "cluster":
		result.Append(nil, &somaClusterResult{
			Cluster: q.Cluster.Cluster,
		})
	case "node":
		result.Append(nil, &somaNodeResult{
			Node: q.Node.Node,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
