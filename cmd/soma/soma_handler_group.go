package main

import (
	"database/sql"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type somaGroupRequest struct {
	action string
	Group  proto.Group
	reply  chan somaResult
}

type somaGroupResult struct {
	ResultError error
	Group       proto.Group
}

func (a *somaGroupResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Groups = append(r.Groups, somaGroupResult{ResultError: err})
	}
}

func (a *somaGroupResult) SomaAppendResult(r *somaResult) {
	r.Groups = append(r.Groups, *a)
}

/* Read Access
 */
type somaGroupReadHandler struct {
	input     chan somaGroupRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	mbgl_stmt *sql.Stmt
	mbcl_stmt *sql.Stmt
	mbnl_stmt *sql.Stmt
	ponc_stmt *sql.Stmt
	psvc_stmt *sql.Stmt
	psys_stmt *sql.Stmt
	pcst_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaGroupReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.GroupList:              r.list_stmt,
		stmt.GroupShow:              r.show_stmt,
		stmt.GroupMemberGroupList:   r.mbgl_stmt,
		stmt.GroupMemberClusterList: r.mbcl_stmt,
		stmt.GroupMemberNodeList:    r.mbnl_stmt,
		stmt.GroupOncProps:          r.ponc_stmt,
		stmt.GroupSvcProps:          r.psvc_stmt,
		stmt.GroupSysProps:          r.psys_stmt,
		stmt.GroupCstProps:          r.pcst_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`Group`, err, statement)
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

func (r *somaGroupReadHandler) process(q *somaGroupRequest) {
	var (
		groupId, groupName, bucketId, groupState, teamId  string
		mGroupId, mGroupName, mClusterId, mClusterName    string
		mNodeId, mNodeName, instanceId, sourceInstanceId  string
		view, oncallId, oncallName, serviceName, customId string
		systemProp, value, customProp                     string
		rows                                              *sql.Rows
		err                                               error
	)
	result := somaResult{}
	resG := proto.Group{}

	switch q.action {
	case "list":
		r.reqLog.Printf("R: group/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			goto dispatch
		}

		for rows.Next() {
			err = rows.Scan(&groupId, &groupName, &bucketId)
			result.Append(err, &somaGroupResult{
				Group: proto.Group{
					Id:       groupId,
					Name:     groupName,
					BucketId: bucketId,
				},
			})
		}
	case "show":
		r.reqLog.Printf("R: group/show for %s", q.Group.Id)
		err = r.show_stmt.QueryRow(q.Group.Id).Scan(
			&groupId,
			&bucketId,
			&groupName,
			&groupState,
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
		group := proto.Group{
			Id:          groupId,
			Name:        groupName,
			BucketId:    bucketId,
			ObjectState: groupState,
			TeamId:      teamId,
		}
		group.Properties = &[]proto.Property{}

		// oncall properties
		rows, err = r.ponc_stmt.Query(q.Group.Id)
		if result.SetRequestError(err) {
			goto dispatch
		}
		for rows.Next() {
			if err = rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&oncallId,
				&oncallName,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*group.Properties = append(
				*group.Properties,
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
		rows, err = r.psvc_stmt.Query(q.Group.Id)
		for rows.Next() {
			if err = rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&serviceName,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*group.Properties = append(
				*group.Properties,
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
		rows, err = r.psys_stmt.Query(q.Group.Id)
		for rows.Next() {
			if err = rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&systemProp,
				&value,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*group.Properties = append(
				*group.Properties,
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
		rows, err = r.pcst_stmt.Query(q.Group.Id)
		for rows.Next() {
			if err = rows.Scan(
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
			*group.Properties = append(
				*group.Properties,
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

		result.Append(err, &somaGroupResult{
			Group: group,
		})
	case "member_list":
		r.reqLog.Printf("R: group/memberlist for %s", q.Group.Id)
		rows, err = r.mbgl_stmt.Query(q.Group.Id)
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		resG = proto.Group{
			Id: q.Group.Id,
		}
		for rows.Next() {
			err = rows.Scan(&mGroupId, &mGroupName, &groupName)
			if err == nil {
				resG.Name = groupName
				*resG.MemberGroups = append(*resG.MemberGroups, proto.Group{
					Id:   mGroupId,
					Name: mGroupName,
				})
			}
		}

		rows, err = r.mbcl_stmt.Query(q.Group.Id)
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&mClusterId, &mClusterName, &groupName)
			if err == nil {
				*resG.MemberClusters = append(*resG.MemberClusters, proto.Cluster{
					Id:   mClusterId,
					Name: mClusterName,
				})
			}
		}

		rows, err = r.mbnl_stmt.Query(q.Group.Id)
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&mNodeId, &mNodeName, &groupName)
			if err == nil {
				*resG.MemberNodes = append(*resG.MemberNodes, proto.Node{
					Id:   mNodeId,
					Name: mNodeName,
				})
			}
		}
		result.Append(err, &somaGroupResult{
			Group: resG,
		})
	default:
		result.SetNotImplemented()
	}

dispatch:
	q.reply <- result
}

/* Ops Access
 */
func (r *somaGroupReadHandler) shutdownNow() {
	r.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
