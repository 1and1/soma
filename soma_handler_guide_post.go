package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"


	"github.com/satori/go.uuid"
)

type guidePost struct {
	input              chan treeRequest
	shutdown           chan bool
	conn               *sql.DB
	jbsv_stmt          *sql.Stmt
	repo_stmt          *sql.Stmt
	name_stmt          *sql.Stmt
	node_stmt          *sql.Stmt
	serv_stmt          *sql.Stmt
	attr_stmt          *sql.Stmt
	cthr_stmt          *sql.Stmt
	cdel_stmt          *sql.Stmt
	bucket_for_node    *sql.Stmt
	bucket_for_cluster *sql.Stmt
	bucket_for_group   *sql.Stmt
}

func (g *guidePost) run() {
	var err error

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
		iu.user_id,
		iu.organizational_team_id,
		$7::jsonb
FROM    inventory.users iu
WHERE   iu.user_uid = $6::varchar;`)
	if err != nil {
		log.Fatal("guide/job-save: ", err)
	}
	defer g.jbsv_stmt.Close()

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

	g.name_stmt, err = g.conn.Prepare(`
SELECT repository_name
FROM   soma.repositories
WHERE  repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("guide/repo-by-id: ", err)
	}
	defer g.name_stmt.Close()

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

	g.cthr_stmt, err = g.conn.Prepare(`
SELECT threshold_amount
FROM   soma.monitoring_capabilities
WHERE  capability_id = $1::uuid;`)
	if err != nil {
		log.Fatal("guide/capability-threshold-lookup: ", err)
	}
	defer g.cthr_stmt.Close()

	if g.cdel_stmt, err = g.conn.Prepare(stmt.CheckDetailsForDelete); err != nil {
		log.Fatal("guide/get-details-for-delete-check: ", err)
	}
	defer g.cdel_stmt.Close()

	if g.bucket_for_node, err = g.conn.Prepare(stmt.NodeBucketId); err != nil {
		log.Fatal("guide/get-bucketid-for-node: ", err)
	}
	defer g.bucket_for_node.Close()

	if g.bucket_for_cluster, err = g.conn.Prepare(stmt.ClusterBucketId); err != nil {
		log.Fatal("guide/get-bucketid-for-cluster: ", err)
	}
	defer g.bucket_for_cluster.Close()

	if g.bucket_for_group, err = g.conn.Prepare(stmt.GroupBucketId); err != nil {
		log.Fatal("guide/get-bucketid-for-group: ", err)
	}
	defer g.bucket_for_group.Close()

	if SomaCfg.Observer {
		fmt.Println(`GuidePost entered observer mode`)
		<-g.shutdown
		goto exit
	}

runloop:
	for {
		select {
		case <-g.shutdown:
			break runloop
		case req := <-g.input:
			g.process(&req)
		}
	}
exit:
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
		// vars used for validation
		valNodeBId, valClusterBId, valGroupBId, valChGBId string
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
		fallthrough
	case `delete_system_property_from_repository`:
		fallthrough
	case `delete_custom_property_from_repository`:
		fallthrough
	case `delete_oncall_property_from_repository`:
		fallthrough
	case `delete_service_property_from_repository`:
		repoId = q.Repository.Repository.Id

	case "create_bucket":
		repoId = q.Bucket.Bucket.RepositoryId
	case "add_system_property_to_bucket":
		fallthrough
	case "add_custom_property_to_bucket":
		fallthrough
	case "add_oncall_property_to_bucket":
		fallthrough
	case "add_service_property_to_bucket":
		fallthrough
	case `delete_system_property_from_bucket`:
		fallthrough
	case `delete_custom_property_from_bucket`:
		fallthrough
	case `delete_oncall_property_from_bucket`:
		fallthrough
	case `delete_service_property_from_bucket`:
		bucketId = q.Bucket.Bucket.Id

	case "add_group_to_group":
		// retrieve bucketId for child group
		if err = g.bucket_for_group.QueryRow(
			(*q.Group.Group.MemberGroups)[0].Id,
		).Scan(&valChGBId); err != nil {
			if err == sql.ErrNoRows {
				result.SetRequestError(
					fmt.Errorf(`GuidePost: child group is not assigned to a bucket`),
				)
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}
		fallthrough
	case "add_cluster_to_group":
		// retrieve bucketId for cluster
		if err = g.bucket_for_cluster.QueryRow(
			(*q.Group.Group.MemberClusters)[0].Id,
		).Scan(&valClusterBId); err != nil {
			if err == sql.ErrNoRows {
				result.SetRequestError(
					fmt.Errorf(`GuidePost: cluster is not assigned to a bucket`),
				)
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}
		fallthrough
	case "add_node_to_group":
		// retrieve bucketId for node
		if err = g.bucket_for_node.QueryRow(
			(*q.Group.Group.MemberNodes)[0].Id,
		).Scan(&valNodeBId); err != nil {
			if err == sql.ErrNoRows {
				result.SetRequestError(
					fmt.Errorf(`GuidePost: node is not assigned to a bucket`),
				)
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}
		fallthrough
	case "create_group":
		fallthrough
	case "add_system_property_to_group":
		fallthrough
	case "add_custom_property_to_group":
		fallthrough
	case "add_oncall_property_to_group":
		fallthrough
	case "add_service_property_to_group":
		fallthrough
	case `delete_system_property_from_group`:
		fallthrough
	case `delete_custom_property_from_group`:
		fallthrough
	case `delete_oncall_property_from_group`:
		fallthrough
	case `delete_service_property_from_group`:
		// group bucketId sent by client
		bucketId = q.Group.Group.BucketId

		// retrieve bucketId for group
		if err = g.bucket_for_group.QueryRow(
			q.Group.Group.Id,
		).Scan(&valGroupBId); err != nil {
			if err == sql.ErrNoRows {
				result.SetRequestError(
					fmt.Errorf(`GuidePost: parent group is not assigned to a bucket`),
				)
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}
		// check if client sent correct bucketId
		if bucketId != valGroupBId {
			result.SetRequestError(
				fmt.Errorf(`GuidePost: parent group is not in specified bucket`),
			)
			q.reply <- result
			return
		}
		// check if node and group are in the same bucket
		if valNodeBId != "" && valNodeBId != bucketId {
			result.SetRequestError(
				fmt.Errorf(`GuidePost: parent group and node are not in the same bucket`),
			)
			q.reply <- result
			return
		}
		// check if cluster and group are in the same bucket
		if valClusterBId != "" && valClusterBId != bucketId {
			result.SetRequestError(
				fmt.Errorf(`GuidePost: parent group and cluster are not in the same bucket`),
			)
			q.reply <- result
			return
		}
		// check if group and group are in the same bucket
		if valChGBId != "" && valChGBId != bucketId {
			result.SetRequestError(
				fmt.Errorf(`GuidePost: parent and child group are not in the same bucket`),
			)
			q.reply <- result
			return
		}

	case "add_node_to_cluster":
		// retrieve bucketId for node
		if err = g.bucket_for_node.QueryRow(
			(*q.Cluster.Cluster.Members)[0].Id,
		).Scan(&valNodeBId); err != nil {
			if err == sql.ErrNoRows {
				result.SetRequestError(
					fmt.Errorf(`GuidePost: node is not assigned to a bucket`),
				)
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
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
		fallthrough
	case `delete_system_property_from_cluster`:
		fallthrough
	case `delete_custom_property_from_cluster`:
		fallthrough
	case `delete_oncall_property_from_cluster`:
		fallthrough
	case `delete_service_property_from_cluster`:
		// cluster bucketId sent by client
		bucketId = q.Cluster.Cluster.BucketId

		// retrieve bucketId for cluster
		if err = g.bucket_for_cluster.QueryRow(
			q.Cluster.Cluster.Id,
		).Scan(&valClusterBId); err != nil {
			if err == sql.ErrNoRows {
				result.SetRequestError(
					fmt.Errorf(`GuidePost: cluster is not assigned to a bucket`),
				)
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}
		// check if client sent correct bucketId
		if bucketId != valClusterBId {
			result.SetRequestError(
				fmt.Errorf(`GuidePost: cluster is not in specified bucket`),
			)
			q.reply <- result
			return
		}
		// check if node and cluster are in the same bucket
		if valNodeBId != "" && valNodeBId != bucketId {
			result.SetRequestError(
				fmt.Errorf(`GuidePost: cluster and node are not in the same bucket`),
			)
			q.reply <- result
			return
		}

	case "add_check_to_repository":
		fallthrough
	case "add_check_to_bucket":
		fallthrough
	case "add_check_to_group":
		fallthrough
	case "add_check_to_cluster":
		fallthrough
	case "add_check_to_node":
		fallthrough
	case `remove_check`:
		repoId = q.CheckConfig.CheckConfig.RepositoryId

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
		q.Node.Node.TeamId = ndTeam
		q.Node.Node.ServerId = ndServer
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
		fallthrough
	case `delete_system_property_from_node`:
		fallthrough
	case `delete_custom_property_from_node`:
		fallthrough
	case `delete_oncall_property_from_node`:
		fallthrough
	case `delete_service_property_from_node`:
		if q.Node.Node.Config == nil {
			_ = result.SetRequestError(fmt.Errorf("NodeConfig subobject missing"))
			q.reply <- result
			return
		}
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
		result.SetRequestError(
			fmt.Errorf(`GuidePost: unable find repository for request`),
		)
		q.reply <- result
		return
	}

	if q.Action == `create_bucket` {
		if !strings.HasPrefix(q.Bucket.Bucket.Name, fmt.Sprintf("%s_", repoName)) {
			result.SetRequestError(
				fmt.Errorf(`Illegal bucket name format, requires reponame_ prefix`),
			)
			q.reply <- result
			return
		}
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
		attrs := []proto.ServiceAttribute{}

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

		// validate the tuple (repo, team, service) is valid
		if err = g.serv_stmt.QueryRow(repoId, svName, svTeam).Scan(&service); err != nil {
			goto inputabort
		}

		// load attributes
		if rows, err = g.attr_stmt.Query(repoId, svName, svTeam); err != nil {
			goto inputabort
		}
		defer rows.Close()

	attrloop:
		for rows.Next() {
			if err = rows.Scan(&attr, &val); err != nil {
				break attrloop
			}
			attrs = append(attrs, proto.ServiceAttribute{
				Name:  attr,
				Value: val,
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

	// check the check configuration to contain fewer thresholds than
	// the limit for the capability
	if strings.Contains(q.Action, "add_check_to_") {
		var (
			thrLimit int
			err      error
		)

		if err = g.cthr_stmt.QueryRow(q.CheckConfig.CheckConfig.CapabilityId).Scan(&thrLimit); err != nil {
			goto validateabort
		}
		if len(q.CheckConfig.CheckConfig.Thresholds) > thrLimit {
			err = fmt.Errorf(
				"Specified %d thresholds exceed limit of %d for capability",
				len(q.CheckConfig.CheckConfig.Thresholds),
				thrLimit)
		}

	validateabort:
		if err != nil {
			if err == sql.ErrNoRows {
				_ = result.SetRequestError(fmt.Errorf(
					"Capability %s not found",
					q.CheckConfig.CheckConfig.CapabilityId))
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		// generate configuration_id
		q.CheckConfig.CheckConfig.Id = uuid.NewV4().String()
	}

	// if the request is a check deletion, populate required IDs
	if q.Action == `remove_check` {
		var delObjId, delObjTyp, delSrcChkId string

		if err = g.cdel_stmt.QueryRow(q.CheckConfig.CheckConfig.Id, q.CheckConfig.CheckConfig.RepositoryId).
			Scan(&delObjId, &delObjTyp, &delSrcChkId); err != nil {
			if err == sql.ErrNoRows {
				result.SetRequestError(fmt.Errorf(
					"Failed to find source check for config %s",
					q.CheckConfig.CheckConfig.Id))
			} else {
				result.SetRequestError(err)
			}
			q.reply <- result
			return
		}
		q.CheckConfig.CheckConfig.ObjectId = delObjId
		q.CheckConfig.CheckConfig.ObjectType = delObjTyp
		q.CheckConfig.CheckConfig.ExternalId = delSrcChkId
		q.Action = fmt.Sprintf("remove_check_from_%s", delObjTyp)
	}

	// if the request is a property deletion, populate required IDs
	if strings.HasPrefix(q.Action, `delete_`) &&
		(strings.HasSuffix(q.Action, `property_from_repository`) ||
			strings.HasSuffix(q.Action, `property_from_bucket`) ||
			strings.HasSuffix(q.Action, `property_from_group`) ||
			strings.HasSuffix(q.Action, `property_from_cluster`) ||
			strings.HasSuffix(q.Action, `property_from_node`)) {
		var (
			err                                             error
			row                                             *sql.Row
			queryStmt, view, sysProp, value, cstId, cstProp string
			svcProp, oncId, oncName                         string
			oncNumber                                       int
		)

		// select SQL statement
		switch q.Action {
		case `delete_system_property_from_repository`:
			queryStmt = stmt.RepoSystemPropertyForDelete
		case `delete_custom_property_from_repository`:
			queryStmt = stmt.RepoCustomPropertyForDelete
		case `delete_service_property_from_repository`:
			queryStmt = stmt.RepoServicePropertyForDelete
		case `delete_oncall_property_from_repository`:
			queryStmt = stmt.RepoOncallPropertyForDelete
		case `delete_system_property_from_bucket`:
			queryStmt = stmt.BucketSystemPropertyForDelete
		case `delete_custom_property_from_bucket`:
			queryStmt = stmt.BucketCustomPropertyForDelete
		case `delete_service_property_from_bucket`:
			queryStmt = stmt.BucketServicePropertyForDelete
		case `delete_oncall_property_from_bucket`:
			queryStmt = stmt.BucketOncallPropertyForDelete
		case `delete_system_property_from_group`:
			queryStmt = stmt.GroupSystemPropertyForDelete
		case `delete_custom_property_from_group`:
			queryStmt = stmt.GroupCustomPropertyForDelete
		case `delete_service_property_from_group`:
			queryStmt = stmt.GroupServicePropertyForDelete
		case `delete_oncall_property_from_group`:
			queryStmt = stmt.GroupOncallPropertyForDelete
		case `delete_system_property_from_cluster`:
			queryStmt = stmt.ClusterSystemPropertyForDelete
		case `delete_custom_property_from_cluster`:
			queryStmt = stmt.ClusterCustomPropertyForDelete
		case `delete_service_property_from_cluster`:
			queryStmt = stmt.ClusterServicePropertyForDelete
		case `delete_oncall_property_from_cluster`:
			queryStmt = stmt.ClusterOncallPropertyForDelete
		case `delete_system_property_from_node`:
			queryStmt = stmt.NodeSystemPropertyForDelete
		case `delete_custom_property_from_node`:
			queryStmt = stmt.NodeCustomPropertyForDelete
		case `delete_service_property_from_node`:
			queryStmt = stmt.NodeServicePropertyForDelete
		case `delete_oncall_property_from_node`:
			queryStmt = stmt.NodeOncallPropertyForDelete
		}

		// execute and scan
		switch q.RequestType {
		case `repository`:
			row = g.conn.QueryRow(queryStmt, (*q.Repository.Repository.Properties)[0].SourceInstanceId)
		case `bucket`:
			row = g.conn.QueryRow(queryStmt, (*q.Bucket.Bucket.Properties)[0].SourceInstanceId)
		case `group`:
			row = g.conn.QueryRow(queryStmt, (*q.Group.Group.Properties)[0].SourceInstanceId)
		case `cluster`:
			row = g.conn.QueryRow(queryStmt, (*q.Cluster.Cluster.Properties)[0].SourceInstanceId)
		case `node`:
			row = g.conn.QueryRow(queryStmt, (*q.Node.Node.Properties)[0].SourceInstanceId)
		}
		switch {
		case strings.HasPrefix(q.Action, `delete_system_`):
			err = row.Scan(&view, &sysProp, &value)

		case strings.HasPrefix(q.Action, `delete_custom_`):
			err = row.Scan(&view, &cstId, &value, &cstProp)

		case strings.HasPrefix(q.Action, `delete_service_`):
			err = row.Scan(&view, &svcProp)

		case strings.HasPrefix(q.Action, `delete_oncall_`):
			err = row.Scan(&view, &oncId, &oncName, &oncNumber)
		}
		if err != nil {
			if err == sql.ErrNoRows {
				result.SetRequestError(fmt.Errorf(
					"Failed to find source property for %s",
					(*q.Repository.Repository.Properties)[0].SourceInstanceId,
				))
			} else {
				result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		// assemble and set results: property specification
		var (
			pSys *proto.PropertySystem
			pCst *proto.PropertyCustom
			pSvc *proto.PropertyService
			pOnc *proto.PropertyOncall
		)
		switch {
		case strings.HasPrefix(q.Action, `delete_system_`):
			pSys = &proto.PropertySystem{
				Name:  sysProp,
				Value: value,
			}
		case strings.HasPrefix(q.Action, `delete_custom_`):
			pCst = &proto.PropertyCustom{
				Id:    cstId,
				Name:  cstProp,
				Value: value,
			}
		case strings.HasPrefix(q.Action, `delete_service_`):
			pSvc = &proto.PropertyService{
				Name: svcProp,
			}
		case strings.HasPrefix(q.Action, `delete_oncall_`):
			num := strconv.Itoa(oncNumber)
			pOnc = &proto.PropertyOncall{
				Id:     oncId,
				Name:   oncName,
				Number: num,
			}
		}

		// assemble and set results: view
		switch {
		case strings.HasSuffix(q.Action, `_repository`):
			(*q.Repository.Repository.Properties)[0].View = view
		case strings.HasSuffix(q.Action, `_bucket`):
			(*q.Bucket.Bucket.Properties)[0].View = view
		case strings.HasSuffix(q.Action, `_group`):
			(*q.Group.Group.Properties)[0].View = view
		case strings.HasSuffix(q.Action, `_cluster`):
			(*q.Cluster.Cluster.Properties)[0].View = view
		case strings.HasSuffix(q.Action, `_node`):
			(*q.Node.Node.Properties)[0].View = view
		}

		// final assembly step
		switch q.Action {
		case `delete_system_property_from_repository`:
			(*q.Repository.Repository.Properties)[0].System = pSys
		case `delete_custom_property_from_repository`:
			(*q.Repository.Repository.Properties)[0].Custom = pCst
		case `delete_service_property_from_repository`:
			(*q.Repository.Repository.Properties)[0].Service = pSvc
		case `delete_oncall_property_from_repository`:
			(*q.Repository.Repository.Properties)[0].Oncall = pOnc

		case `delete_system_property_from_bucket`:
			(*q.Bucket.Bucket.Properties)[0].System = pSys
		case `delete_custom_property_from_bucket`:
			(*q.Bucket.Bucket.Properties)[0].Custom = pCst
		case `delete_service_property_from_bucket`:
			(*q.Bucket.Bucket.Properties)[0].Service = pSvc
		case `delete_oncall_property_from_bucket`:
			(*q.Bucket.Bucket.Properties)[0].Oncall = pOnc

		case `delete_system_property_from_group`:
			(*q.Group.Group.Properties)[0].System = pSys
		case `delete_custom_property_from_group`:
			(*q.Group.Group.Properties)[0].Custom = pCst
		case `delete_service_property_from_group`:
			(*q.Group.Group.Properties)[0].Service = pSvc
		case `delete_oncall_property_from_group`:
			(*q.Group.Group.Properties)[0].Oncall = pOnc

		case `delete_system_property_from_cluster`:
			(*q.Cluster.Cluster.Properties)[0].System = pSys
		case `delete_custom_property_from_cluster`:
			(*q.Cluster.Cluster.Properties)[0].Custom = pCst
		case `delete_service_property_from_cluster`:
			(*q.Cluster.Cluster.Properties)[0].Service = pSvc
		case `delete_oncall_property_from_cluster`:
			(*q.Cluster.Cluster.Properties)[0].Oncall = pOnc

		case `delete_system_property_from_node`:
			(*q.Node.Node.Properties)[0].System = pSys
		case `delete_custom_property_from_node`:
			(*q.Node.Node.Properties)[0].Custom = pCst
		case `delete_service_property_from_node`:
			(*q.Node.Node.Properties)[0].Service = pSvc
		case `delete_oncall_property_from_node`:
			(*q.Node.Node.Properties)[0].Oncall = pOnc
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
		q.User,
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
		case "check":
			result.Append(err, &somaCheckConfigResult{})
		}
		q.reply <- result
		return
	}

	handler.input <- *q
	result.JobId = q.JobId.String()
	result.JobType = q.Action

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
	case "check":
		result.Append(nil, &somaCheckConfigResult{
			CheckConfig: q.CheckConfig.CheckConfig,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
