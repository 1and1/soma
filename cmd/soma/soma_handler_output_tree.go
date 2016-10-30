/*-
Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type outputTree struct {
	input    chan msg.Request
	shutdown chan bool
	conn     *sql.DB
	// object details
	stmtRepo    *sql.Stmt
	stmtBucket  *sql.Stmt
	stmtGroup   *sql.Stmt
	stmtCluster *sql.Stmt
	stmtNode    *sql.Stmt
	// object tree
	stmtRepoBuck *sql.Stmt
	stmtBuckGrp  *sql.Stmt
	stmtBuckClr  *sql.Stmt
	stmtBuckNod  *sql.Stmt
	stmtGrpGrp   *sql.Stmt
	stmtGrpClr   *sql.Stmt
	stmtGrpNod   *sql.Stmt
	stmtClrNod   *sql.Stmt
	appLog       *log.Logger
	reqLog       *log.Logger
	errLog       *log.Logger
}

func (o *outputTree) run() {
	var err error

	// single-object return statements
	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.TreeShowRepository:      o.stmtRepo,
		stmt.TreeShowBucket:          o.stmtBucket,
		stmt.TreeShowGroup:           o.stmtGroup,
		stmt.TreeShowCluster:         o.stmtCluster,
		stmt.TreeShowNode:            o.stmtNode,
		stmt.TreeBucketsInRepository: o.stmtRepoBuck,
		stmt.TreeGroupsInBucket:      o.stmtBuckGrp,
		stmt.TreeClustersInBucket:    o.stmtBuckClr,
		stmt.TreeNodesInBucket:       o.stmtBuckNod,
		stmt.TreeGroupsInGroup:       o.stmtGrpGrp,
		stmt.TreeClustersInGroup:     o.stmtGrpClr,
		stmt.TreeNodesInGroup:        o.stmtGrpNod,
		stmt.TreeNodesInCluster:      o.stmtClrNod,
	} {
		if prepStmt, err = o.conn.Prepare(statement); err != nil {
			o.errLog.Fatal(`outputTree`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-o.shutdown:
			break runloop
		case req := <-o.input:
			go func() {
				o.process(&req)
			}()
		}
	}
}

//
func (o *outputTree) process(q *msg.Request) {
	result := msg.Result{Type: `tree`, Action: q.Type}
	tree := proto.Tree{
		Id:   q.Tree.Id,
		Type: q.Tree.Type,
	}
	var err error

	switch tree.Type {
	case `repository`:
		tree.Repository, err = o.repository(tree.Id, true, 0)
	case `bucket`:
		tree.Bucket, err = o.bucket(tree.Id, true, 0)
	case `group`:
		tree.Group, err = o.group(tree.Id, true, 0)
	case `cluster`:
		tree.Cluster, err = o.cluster(tree.Id, true, 0)
	case `node`:
		tree.Node, err = o.node(tree.Id, true, 0)
	}
	if err == sql.ErrNoRows {
		result.NotFound(fmt.Errorf(`Tree starting point not found`))
	} else if err != nil {
		result.ServerError(err)
	} else {
		result.Tree = tree
		result.OK()
	}

	q.Reply <- result
}

//
func (o *outputTree) repository(id string, recurse bool, depth int) (*proto.Repository, error) {
	var (
		repoName, teamId, repoCreatedBy string
		repoActive, repoDeleted         bool
		repoCreatedAt                   time.Time
	)

	if depth >= 64 {
		return nil, fmt.Errorf(`Maximum recursion depth exceeded`)
	}

	//
	if err := o.stmtRepo.QueryRow(id).Scan(
		&repoName,
		&repoActive,
		&teamId,
		&repoDeleted,
		&repoCreatedBy,
		&repoCreatedAt,
	); err != nil {
		o.errLog.Printf("Error in outputTree.repository() for %s: %s", id, err.Error())
		return nil, err
	}

	//
	repo := proto.Repository{
		Id:        id,
		Name:      repoName,
		TeamId:    teamId,
		IsDeleted: repoDeleted,
		IsActive:  repoActive,
		Details: &proto.Details{
			CreatedBy: repoCreatedBy,
			CreatedAt: repoCreatedAt.UTC().Format(time.RFC3339),
		},
	}

	//
	if !recurse {
		return &repo, nil
	}
	depth++
	repo.Members = &[]proto.Bucket{}

	//
	buckets := o.bucketsInRepository(id)
	for i, _ := range buckets {
		b, err := o.bucket(buckets[i], recurse, depth)
		if err != nil {
			return nil, err
		}
		*repo.Members = append(*repo.Members, *b)
	}
	return &repo, nil
}

//
func (o *outputTree) bucket(id string, recurse bool, depth int) (*proto.Bucket, error) {
	var (
		bucketName, bucketRepositoryId, bucketEnvironment string
		bucketTeamId, bucketCreatedBy                     string
		bucketIsFrozen, bucketIsDeleted                   bool
		bucketCreatedAt                                   time.Time
	)

	if depth >= 64 {
		return nil, fmt.Errorf(`Maximum recursion depth exceeded`)
	}

	//
	if err := o.stmtBucket.QueryRow(id).Scan(
		&bucketName,
		&bucketIsFrozen,
		&bucketIsDeleted,
		&bucketRepositoryId,
		&bucketEnvironment,
		&bucketTeamId,
		&bucketCreatedBy,
		&bucketCreatedAt,
	); err != nil {
		o.errLog.Printf("Error in outputTree.bucket() for %s: %s", id, err.Error())
		return nil, err
	}

	//
	bucket := proto.Bucket{
		Id:           id,
		Name:         bucketName,
		RepositoryId: bucketRepositoryId,
		TeamId:       bucketTeamId,
		Environment:  bucketEnvironment,
		IsDeleted:    bucketIsDeleted,
		IsFrozen:     bucketIsFrozen,
		Details: &proto.Details{
			CreatedBy: bucketCreatedBy,
			CreatedAt: bucketCreatedAt.UTC().Format(time.RFC3339),
		},
	}

	//
	if !recurse {
		return &bucket, nil
	}
	depth++
	bucket.MemberGroups = &[]proto.Group{}
	bucket.MemberClusters = &[]proto.Cluster{}
	bucket.MemberNodes = &[]proto.Node{}

	//
	groups := o.groupsInBucket(id)
	for i, _ := range groups {
		g, err := o.group(groups[i], recurse, depth)
		if err != nil {
			return nil, err
		}
		*bucket.MemberGroups = append(*bucket.MemberGroups, *g)
	}

	//
	clusters := o.clustersInBucket(id)
	for i, _ := range clusters {
		c, err := o.cluster(clusters[i], recurse, depth)
		if err != nil {
			return nil, err
		}
		*bucket.MemberClusters = append(*bucket.MemberClusters, *c)
	}

	//
	nodes := o.nodesInBucket(id)
	for i, _ := range nodes {
		n, err := o.node(nodes[i], recurse, depth)
		if err != nil {
			return nil, err
		}
		*bucket.MemberNodes = append(*bucket.MemberNodes, *n)
	}
	return &bucket, nil
}

//
func (o *outputTree) group(id string, recurse bool, depth int) (*proto.Group, error) {
	var (
		groupBucketId, groupName, groupObjectState string
		groupTeamId, groupCreatedBy                string
		groupCreatedAt                             time.Time
	)

	if depth >= 64 {
		return nil, fmt.Errorf(`Maximum recursion depth exceeded`)
	}

	//
	if err := o.stmtGroup.QueryRow(id).Scan(
		&groupBucketId,
		&groupName,
		&groupObjectState,
		&groupTeamId,
		&groupCreatedBy,
		&groupCreatedAt,
	); err != nil {
		o.errLog.Printf("Error in outputTree.group() for %s: %s", id, err.Error())
		return nil, err
	}

	//
	group := proto.Group{
		Id:          id,
		BucketId:    groupBucketId,
		Name:        groupName,
		ObjectState: groupObjectState,
		TeamId:      groupTeamId,
		Details: &proto.Details{
			CreatedBy: groupCreatedBy,
			CreatedAt: groupCreatedAt.UTC().Format(time.RFC3339),
		},
	}

	//
	if !recurse {
		return &group, nil
	}
	depth++
	group.MemberGroups = &[]proto.Group{}
	group.MemberClusters = &[]proto.Cluster{}
	group.MemberNodes = &[]proto.Node{}

	//
	groups := o.groupsInGroup(id)
	for i, _ := range groups {
		g, err := o.group(groups[i], recurse, depth)
		if err != nil {
			return nil, err
		}
		*group.MemberGroups = append(*group.MemberGroups, *g)
	}

	//
	clusters := o.clustersInGroup(id)
	for i, _ := range clusters {
		c, err := o.cluster(clusters[i], recurse, depth)
		if err != nil {
			return nil, err
		}
		*group.MemberClusters = append(*group.MemberClusters, *c)
	}

	//
	nodes := o.nodesInGroup(id)
	for i, _ := range nodes {
		n, err := o.node(nodes[i], recurse, depth)
		if err != nil {
			return nil, err
		}
		*group.MemberNodes = append(*group.MemberNodes, *n)
	}
	return &group, nil
}

//
func (o *outputTree) cluster(id string, recurse bool, depth int) (*proto.Cluster, error) {
	var (
		clusterName, clusterBucketId                        string
		clusterObjectState, clusterTeamId, clusterCreatedBy string
		clusterCreatedAt                                    time.Time
	)

	if depth >= 64 {
		return nil, fmt.Errorf(`Maximum recursion depth exceeded`)
	}

	//
	if err := o.stmtCluster.QueryRow(id).Scan(
		&clusterName,
		&clusterBucketId,
		&clusterObjectState,
		&clusterTeamId,
		&clusterCreatedBy,
		&clusterCreatedAt,
	); err != nil {
		o.errLog.Printf("Error in outputTree.cluster() for %s: %s", id, err.Error())
		return nil, err
	}

	//
	cluster := proto.Cluster{
		Id:          id,
		Name:        clusterName,
		BucketId:    clusterBucketId,
		ObjectState: clusterObjectState,
		TeamId:      clusterTeamId,
		Details: &proto.Details{
			CreatedBy: clusterCreatedBy,
			CreatedAt: clusterCreatedAt.UTC().Format(time.RFC3339),
		},
	}

	//
	if !recurse {
		return &cluster, nil
	}
	depth++
	cluster.Members = &[]proto.Node{}

	//
	nodes := o.nodesInCluster(id)
	for i, _ := range nodes {
		n, err := o.node(nodes[i], recurse, depth)
		if err != nil {
			return nil, err
		}
		*cluster.Members = append(*cluster.Members, *n)
	}
	return &cluster, nil
}

//
func (o *outputTree) node(id string, recurse bool, depth int) (*proto.Node, error) {
	var (
		nodeAssetId                                   int
		nodeName, nodeTeamId, nodeServerId, nodeState string
		nodeIsOnline, nodeIsDeleted                   bool
		nodeCreatedBy, nodeRepositoryId, nodeBucketId string
		nodeCreatedAt                                 time.Time
	)

	if depth >= 64 {
		return nil, fmt.Errorf(`Maximum recursion depth exceeded`)
	}

	//
	if err := o.stmtNode.QueryRow(id).Scan(
		&nodeAssetId,
		&nodeName,
		&nodeTeamId,
		&nodeServerId,
		&nodeState,
		&nodeIsOnline,
		&nodeIsDeleted,
		&nodeCreatedBy,
		&nodeCreatedAt,
		&nodeRepositoryId,
		&nodeBucketId,
	); err != nil {
		o.errLog.Printf("Error in outputTree.node() for %s: %s", id, err.Error())
		return nil, err
	}

	//
	node := proto.Node{
		Id:        id,
		AssetId:   uint64(nodeAssetId),
		Name:      nodeName,
		TeamId:    nodeTeamId,
		ServerId:  nodeServerId,
		State:     nodeState,
		IsOnline:  nodeIsOnline,
		IsDeleted: nodeIsDeleted,
		Details: &proto.Details{
			CreatedBy: nodeCreatedBy,
			CreatedAt: nodeCreatedAt.UTC().Format(time.RFC3339),
		},
		Config: &proto.NodeConfig{
			RepositoryId: nodeRepositoryId,
			BucketId:     nodeBucketId,
		},
	}
	return &node, nil
}

//
func (o *outputTree) bucketsInRepository(id string) []string {
	rows, err := o.stmtRepoBuck.Query(id)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	res := []string{}
	for rows.Next() {
		bID := ``
		if err := rows.Scan(&bID); err != nil {
			return []string{}
		}
		res = append(res, bID)
	}
	if rows.Err() != nil {
		return []string{}
	}
	return res
}

//
func (o *outputTree) groupsInBucket(id string) []string {
	rows, err := o.stmtBuckGrp.Query(id)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	res := []string{}
	for rows.Next() {
		gID := ``
		if err := rows.Scan(&gID); err != nil {
			return []string{}
		}
		res = append(res, gID)
	}
	if rows.Err() != nil {
		return []string{}
	}
	return res
}

//
func (o *outputTree) clustersInBucket(id string) []string {
	rows, err := o.stmtBuckClr.Query(id)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	res := []string{}
	for rows.Next() {
		cID := ``
		if err := rows.Scan(&cID); err != nil {
			return []string{}
		}
		res = append(res, cID)
	}
	if rows.Err() != nil {
		return []string{}
	}
	return res
}

//
func (o *outputTree) nodesInBucket(id string) []string {
	rows, err := o.stmtBuckNod.Query(id)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	res := []string{}
	for rows.Next() {
		nID := ``
		if err := rows.Scan(&nID); err != nil {
			return []string{}
		}
		res = append(res, nID)
	}
	if rows.Err() != nil {
		return []string{}
	}
	return res
}

//
func (o *outputTree) groupsInGroup(id string) []string {
	rows, err := o.stmtGrpGrp.Query(id)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	res := []string{}
	for rows.Next() {
		gID := ``
		if err := rows.Scan(&gID); err != nil {
			return []string{}
		}
		res = append(res, gID)
	}
	if rows.Err() != nil {
		return []string{}
	}
	return res
}

//
func (o *outputTree) clustersInGroup(id string) []string {
	rows, err := o.stmtGrpClr.Query(id)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	res := []string{}
	for rows.Next() {
		cID := ``
		if err := rows.Scan(&cID); err != nil {
			return []string{}
		}
		res = append(res, cID)
	}
	if rows.Err() != nil {
		return []string{}
	}
	return res
}

//
func (o *outputTree) nodesInGroup(id string) []string {
	rows, err := o.stmtGrpNod.Query(id)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	res := []string{}
	for rows.Next() {
		nID := ``
		if err := rows.Scan(&nID); err != nil {
			return []string{}
		}
		res = append(res, nID)
	}
	if rows.Err() != nil {
		return []string{}
	}
	return res
}

//
func (o *outputTree) nodesInCluster(id string) []string {
	rows, err := o.stmtClrNod.Query(id)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	res := []string{}
	for rows.Next() {
		nID := ``
		if err := rows.Scan(&nID); err != nil {
			return []string{}
		}
		res = append(res, nID)
	}
	if rows.Err() != nil {
		return []string{}
	}
	return res
}

/* Ops Access
 */
func (o *outputTree) shutdownNow() {
	o.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
