package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	"github.com/satori/go.uuid"
)

//
func (g *guidePost) fillReqData(q *treeRequest) (error, bool) {
	switch {
	case strings.Contains(q.Action, "add_service_property_to_"):
		return g.fillServiceAttributes(q)
	case q.Action == `assign_node`:
		return g.fillNode(q)
	case q.Action == `remove_check`:
		return g.fillCheckDeleteInfo(q)
	case strings.HasPrefix(q.Action, `delete_`) &&
		strings.Contains(q.Action, `_property_from_`):
		return g.fillPropertyDeleteInfo(q)
	case strings.HasPrefix(q.Action, `add_check_to_`):
		return g.fillCheckConfigId(q)
	default:
		return nil, false
	}
}

// generate CheckConfigId
func (g *guidePost) fillCheckConfigId(q *treeRequest) (error, bool) {
	q.CheckConfig.CheckConfig.Id = uuid.NewV4().String()
	return nil, false
}

// Populate the node structure with data, overwriting the client
// submitted values.
func (g *guidePost) fillNode(q *treeRequest) (error, bool) {
	var (
		err                      error
		ndName, ndTeam, ndServer string
		ndAsset                  int64
		ndOnline, ndDeleted      bool
	)
	if err = g.node_stmt.QueryRow(q.Node.Node.Id).Scan(
		&ndAsset,
		&ndName,
		&ndTeam,
		&ndServer,
		&ndOnline,
		&ndDeleted,
	); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Node not found: %s", q.Node.Node.Id), true
		}
		return err, false
	}
	q.Node.Node.AssetId = uint64(ndAsset)
	q.Node.Node.Name = ndName
	q.Node.Node.TeamId = ndTeam
	q.Node.Node.ServerId = ndServer
	q.Node.Node.IsOnline = ndOnline
	q.Node.Node.IsDeleted = ndDeleted
	return nil, false
}

// load authoritative copy of the service attributes from the
// database. Replaces whatever the client sent in.
func (g *guidePost) fillServiceAttributes(q *treeRequest) (error, bool) {
	var (
		service, attr, val, svName, svTeam, repoId string
		rows                                       *sql.Rows
		err                                        error
		nf                                         bool
	)
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

	// ignore error since it would have been caught by guidePost
	repoId, _, _, _ = g.extractRouting(q)

	// validate the tuple (repo, team, service) is valid
	if err = g.serv_stmt.QueryRow(repoId, svName, svTeam).Scan(&service); err != nil {
		if err == sql.ErrNoRows {
			nf = true
			err = fmt.Errorf("Requested service %s not available for team %s",
				svName, svTeam)
		}
		goto abort
	}

	// load attributes
	if rows, err = g.attr_stmt.Query(repoId, svName, svTeam); err != nil {
		goto abort
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
abort:
	if err != nil {
		return err, nf
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
	return nil, false
}

// if the request is a check deletion, populate required IDs
func (g *guidePost) fillCheckDeleteInfo(q *treeRequest) (error, bool) {
	var delObjId, delObjTyp, delSrcChkId string
	var err error

	if err = g.cdel_stmt.QueryRow(
		q.CheckConfig.CheckConfig.Id,
		q.CheckConfig.CheckConfig.RepositoryId,
	).Scan(
		&delObjId,
		&delObjTyp,
		&delSrcChkId,
	); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf(
				"Failed to find source check for config %s",
				q.CheckConfig.CheckConfig.Id), true
		}
		return err, false
	}
	q.CheckConfig.CheckConfig.ObjectId = delObjId
	q.CheckConfig.CheckConfig.ObjectType = delObjTyp
	q.CheckConfig.CheckConfig.ExternalId = delSrcChkId
	q.Action = fmt.Sprintf("remove_check_from_%s", delObjTyp)
	return nil, false
}

// if the request is a property deletion, populate required IDs
func (g *guidePost) fillPropertyDeleteInfo(q *treeRequest) (error, bool) {
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
		row = g.conn.QueryRow(queryStmt,
			(*q.Repository.Repository.Properties)[0].SourceInstanceId)
	case `bucket`:
		row = g.conn.QueryRow(queryStmt,
			(*q.Bucket.Bucket.Properties)[0].SourceInstanceId)
	case `group`:
		row = g.conn.QueryRow(queryStmt,
			(*q.Group.Group.Properties)[0].SourceInstanceId)
	case `cluster`:
		row = g.conn.QueryRow(queryStmt,
			(*q.Cluster.Cluster.Properties)[0].SourceInstanceId)
	case `node`:
		row = g.conn.QueryRow(queryStmt,
			(*q.Node.Node.Properties)[0].SourceInstanceId)
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
			return fmt.Errorf(
				"Failed to find source property for %s",
				(*q.Repository.Repository.Properties)[0].SourceInstanceId), true
		}
		return err, false
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
	return nil, false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
