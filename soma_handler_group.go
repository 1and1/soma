package main

import (
	"database/sql"
	"log"

)

type somaGroupRequest struct {
	action string
	Group  somaproto.ProtoGroup
	reply  chan somaResult
}

type somaGroupResult struct {
	ResultError error
	Group       somaproto.ProtoGroup
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
}

func (r *somaGroupReadHandler) run() {
	var err error

	log.Println("Prepare: group/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT group_id,
       group_name
FROM soma.groups;`)
	if err != nil {
		log.Fatal("group/list: ", err)
	}
	defer r.list_stmt.Close()

	log.Println("Prepare: group/show")
	r.show_stmt, err = r.conn.Prepare(`
SELECT group_id,
       bucket_id,
	   group_name,
	   object_state,
	   organizational_team_id
FROM   soma.groups
WHERE  group_id = $1::uuid;`)
	if err != nil {
		log.Fatal("group/show: ", err)
	}
	defer r.show_stmt.Close()

	log.Println("Prepare: group/memberlist-group")
	r.mbgl_stmt, err = r.conn.Prepare(`
SELECT sg.group_id,
       sg.group_name,
	   osg.group_name
FROM   soma.group_membership_groups sgmg
JOIN   soma.groups sg
ON     sgmg.child_group_id = sg.group_id
JOIN   soma.groups osg
ON     sgmg.group_id = osg.group_id
WHERE  sgmg.group_id = $1::uuid;`)
	if err != nil {
		log.Fatal("group/memberlist-group: ", err)
	}
	defer r.mbgl_stmt.Close()

	log.Println("Prepare: group/memberlist-cluster")
	r.mbcl_stmt, err = r.conn.Prepare(`
SELECT sc.cluster_id,
       sc.cluster_name,
	   sg.group_name
FROM   soma.group_membership_clusters sgmc
JOIN   soma.clusters sc
ON     sgmc.child_cluster_id = sc.cluster_id
JOIN   soma.groups sg
ON     sgmc.group_id = sg.group_id
WHERE  sgmc.group_id = $1::uuid;`)
	if err != nil {
		log.Fatal("group/memberlist-cluster: ", err)
	}
	defer r.mbcl_stmt.Close()

	log.Println("Prepare: group/memberlist-node")
	r.mbnl_stmt, err = r.conn.Prepare(`
SELECT sn.node_id,
       sn.node_name,
	   sg.group_name
FROM   soma.group_membership_nodes sgmn
JOIN   soma.nodes sn
ON     sgmn.child_node_id = sn.node_id
JOIN   soma.groups sg
ON     sgmn.group_id = sg.group_id
WHERE  sgmn.group_id = $1::uuid;`)
	if err != nil {
		log.Fatal("group/memberlist-node: ", err)
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

func (r *somaGroupReadHandler) process(q *somaGroupRequest) {
	var (
		groupId, groupName, bucketId, groupState, teamId string
		mGroupId, mGroupName, mClusterId, mClusterName   string
		mNodeId, mNodeName                               string
		rows                                             *sql.Rows
		err                                              error
	)
	result := somaResult{}
	resG := somaproto.ProtoGroup{}

	switch q.action {
	case "list":
		log.Printf("R: group/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&groupId, &groupName)
			result.Append(err, &somaGroupResult{
				Group: somaproto.ProtoGroup{
					Id:   groupId,
					Name: groupName,
				},
			})
		}
	case "show":
		log.Printf("R: group/show for %s", q.Group.Id)
		err = r.show_stmt.QueryRow(q.Group.Id).Scan(
			&groupId,
			&bucketId,
			&groupName,
			&groupState,
			&teamId,
		)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		result.Append(err, &somaGroupResult{
			Group: somaproto.ProtoGroup{
				Id:          groupId,
				Name:        groupName,
				BucketId:    bucketId,
				ObjectState: groupState,
				TeamId:      teamId,
			},
		})
	case "member_list":
		log.Printf("R: group/memberlist for %s", q.Group.Id)
		rows, err = r.mbgl_stmt.Query(q.Group.Id)
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		resG = somaproto.ProtoGroup{
			Id: q.Group.Id,
		}
		for rows.Next() {
			err := rows.Scan(&mGroupId, &mGroupName, &groupName)
			if err == nil {
				resG.Name = groupName
				resG.MemberGroups = append(resG.MemberGroups, somaproto.ProtoGroup{
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
			err := rows.Scan(&mClusterId, &mClusterName, &groupName)
			if err == nil {
				resG.MemberClusters = append(resG.MemberClusters, somaproto.ProtoCluster{
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
				resG.MemberNodes = append(resG.MemberNodes, somaproto.ProtoNode{
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
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
