package main

import (
	"database/sql"
	"fmt"
	"strings"
)

func (g *guidePost) validateRequest(q *treeRequest) (error, bool) {
	switch q.RequestType {
	case `check`:
		if err, nf := g.validateCheckObjectInBucket(q); err != nil {
			return err, nf
		}
	case `node`:
		if err, nf := g.validateNodeConfig(q); err != nil {
			return err, nf
		}
		fallthrough
	case `cluster`, `group`:
		if err, nf := g.validateCorrectBucket(q); err != nil {
			return err, nf
		}
	case `bucket`:
		if err, nf := g.validateBucketInRepository(
			q.Bucket.Bucket.RepositoryId,
			q.Bucket.Bucket.Id,
		); err != nil {
			return err, nf
		}
	case `repository`:
		// since repository ids are the routing information,
		// it is unnecessary to check that the object is where the
		// routing would point to
	default:
		return fmt.Errorf("Invalid request type %s", q.RequestType), false
	}

	switch q.Action {
	case
		`add_node_to_cluster`,
		`add_node_to_group`,
		`add_cluster_to_group`,
		`add_group_to_group`:
		return g.validateObjectMatch(q)
	case
		`add_check_to_bucket`,
		`add_check_to_cluster`,
		`add_check_to_group`,
		`add_check_to_node`,
		`add_check_to_repository`:
		return g.validateCheckThresholds(q)
	case
		`create_bucket`:
		return g.validateBucketName(q)
	case
		`add_custom_property_to_bucket`,
		`add_custom_property_to_cluster`,
		`add_custom_property_to_group`,
		`add_custom_property_to_node`,
		`add_custom_property_to_repository`,
		`add_oncall_property_to_bucket`,
		`add_oncall_property_to_cluster`,
		`add_oncall_property_to_group`,
		`add_oncall_property_to_node`,
		`add_oncall_property_to_repository`,
		`add_service_property_to_bucket`,
		`add_service_property_to_cluster`,
		`add_service_property_to_group`,
		`add_service_property_to_node`,
		`add_service_property_to_repository`,
		`add_system_property_to_bucket`,
		`add_system_property_to_cluster`,
		`add_system_property_to_group`,
		`add_system_property_to_node`,
		`add_system_property_to_repository`,
		`assign_node`,
		`create_cluster`,
		`create_group`,
		`delete_custom_property_from_bucket`,
		`delete_custom_property_from_cluster`,
		`delete_custom_property_from_group`,
		`delete_custom_property_from_node`,
		`delete_custom_property_from_repository`,
		`delete_oncall_property_from_bucket`,
		`delete_oncall_property_from_cluster`,
		`delete_oncall_property_from_group`,
		`delete_oncall_property_from_node`,
		`delete_oncall_property_from_repository`,
		`delete_service_property_from_bucket`,
		`delete_service_property_from_cluster`,
		`delete_service_property_from_group`,
		`delete_service_property_from_node`,
		`delete_service_property_from_repository`,
		`delete_system_property_from_bucket`,
		`delete_system_property_from_cluster`,
		`delete_system_property_from_group`,
		`delete_system_property_from_node`,
		`delete_system_property_from_repository`,
		`remove_check`:
		// actions are accepted, but require no further validation
		return nil, false
	default:
		return fmt.Errorf("Unimplemented GuidePost/%s", q.Action), false
	}
}

func (g *guidePost) validateObjectMatch(q *treeRequest) (error, bool) {
	var (
		nodeId, clusterId, groupId, childGroupId              string
		valNodeBId, valClusterBId, valGroupBId, valChGroupBId string
	)

	switch q.Action {
	case `add_node_to_cluster`:
		nodeId = (*q.Cluster.Cluster.Members)[0].Id
		clusterId = q.Cluster.Cluster.Id
	case `add_node_to_group`:
		nodeId = (*q.Group.Group.MemberNodes)[0].Id
		groupId = q.Group.Group.Id
	case `add_cluster_to_group`:
		clusterId = (*q.Group.Group.MemberClusters)[0].Id
		groupId = q.Group.Group.Id
	case `add_group_to_group`:
		childGroupId = (*q.Group.Group.MemberGroups)[0].Id
		groupId = q.Group.Group.Id
	default:
		return fmt.Errorf("Incorrect validation attempted for %s",
			q.Action), false
	}

	if nodeId != `` {
		if err := g.bucket_for_node.QueryRow(nodeId).Scan(
			&valNodeBId,
		); err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("Unknown node %s", nodeId), true
			}
			return err, false
		}
	}
	if clusterId != `` {
		if err := g.bucket_for_cluster.QueryRow(clusterId).Scan(
			&valClusterBId,
		); err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("Unknown cluster %s", clusterId), true
			}
			return err, false
		}
	}
	if groupId != `` {
		if err := g.bucket_for_group.QueryRow(groupId).Scan(
			&valGroupBId,
		); err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("Unknown group %s", groupId), true
			}
			return err, false
		}
	}
	if childGroupId != `` {
		if err := g.bucket_for_group.QueryRow(childGroupId).Scan(
			&valChGroupBId,
		); err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("Unknown group %s", childGroupId), true
			}
			return err, false
		}
	}

	switch q.Action {
	case `add_node_to_cluster`:
		if valNodeBId != valClusterBId {
			return fmt.Errorf(
				"Node and Cluster are in different buckets (%s/%s)",
				valNodeBId, valClusterBId,
			), false
		}
	case `add_node_to_group`:
		if valNodeBId != valGroupBId {
			return fmt.Errorf(
				"Node and Group are in different buckets (%s/%s)",
				valNodeBId, valGroupBId,
			), false
		}
	case `add_cluster_to_group`:
		if valClusterBId != valGroupBId {
			return fmt.Errorf(
				"Cluster and Group are in different buckets (%s/%s)",
				valClusterBId, valGroupBId,
			), false
		}
	case `add_group_to_group`:
		if valChGroupBId != valGroupBId {
			return fmt.Errorf(
				"Groups are in different buckets (%s/%s)",
				valGroupBId, valChGroupBId,
			), false
		}
	}
	return nil, false
}

// Verify that an object is assigned to the specified bucket.
func (g *guidePost) validateCorrectBucket(q *treeRequest) (error, bool) {
	switch q.Action {
	case `assign_node`:
		return g.validateNodeUnassigned(q)
	case `create_cluster`, `create_group`:
		return nil, false
	}
	var bid string
	var err error
	switch q.RequestType {
	case `node`:
		err = g.bucket_for_node.QueryRow(
			q.Node.Node.Id,
		).Scan(
			&bid,
		)
	case `cluster`:
		err = g.bucket_for_cluster.QueryRow(
			q.Cluster.Cluster.Id,
		).Scan(
			&bid,
		)
	case `group`:
		err = g.bucket_for_group.QueryRow(
			q.Group.Group.Id,
		).Scan(
			&bid,
		)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			// unassigned
			return fmt.Errorf("%s is not assigned to any bucket",
				q.RequestType), true
		}
		return err, false
	}
	switch q.RequestType {
	case `node`:
		if bid != q.Node.Node.Config.BucketId {
			return fmt.Errorf("Node assigned to different bucket %s",
				bid), false
		}
	case `cluster`:
		if bid != q.Cluster.Cluster.BucketId {
			return fmt.Errorf("Cluster in different bucket %s",
				bid), false
		}
	case `group`:
		if bid != q.Group.Group.BucketId {
			return fmt.Errorf("Group in different bucket %s",
				bid), false
		}
	}
	return nil, false
}

// Verify that a node is not yet assigned to a bucket. Returns nil
// on success.
func (g *guidePost) validateNodeUnassigned(q *treeRequest) (error, bool) {
	var bid string
	if err := g.bucket_for_node.QueryRow(q.Node.Node.Id).Scan(
		&bid,
	); err != nil {
		if err == sql.ErrNoRows {
			// unassigned is not an error here
			return nil, false
		}
		return err, false
	}
	return fmt.Errorf("Node already assigned to bucket %s", bid), false
}

// Verify the node has a Config section
func (g *guidePost) validateNodeConfig(q *treeRequest) (error, bool) {
	if q.Node.Node.Config == nil {
		return fmt.Errorf("NodeConfig subobject missing"), false
	}
	return g.validateBucketInRepository(
		q.Node.Node.Config.RepositoryId,
		q.Node.Node.Config.BucketId,
	)
}

// Verify that the ObjectId->BucketId->RepositoryId chain is part of
// the same tree.
func (g *guidePost) validateCheckObjectInBucket(q *treeRequest) (error, bool) {
	var err error
	var bid string
	switch q.CheckConfig.CheckConfig.ObjectType {
	case `repository`:
		if q.CheckConfig.CheckConfig.RepositoryId !=
			q.CheckConfig.CheckConfig.ObjectId {
			return fmt.Errorf("Conflicting repository ids: %s, %s",
				q.CheckConfig.CheckConfig.RepositoryId,
				q.CheckConfig.CheckConfig.ObjectId,
			), false
		}
		return nil, false
	case `bucket`:
		bid = q.CheckConfig.CheckConfig.ObjectId
	case `group`:
		err = g.bucket_for_group.QueryRow(
			q.CheckConfig.CheckConfig.ObjectId,
		).Scan(&bid)
	case `cluster`:
		err = g.bucket_for_cluster.QueryRow(
			q.CheckConfig.CheckConfig.ObjectId,
		).Scan(&bid)
	case `node`:
		err = g.bucket_for_node.QueryRow(
			q.CheckConfig.CheckConfig.ObjectId,
		).Scan(&bid)
	default:
		return fmt.Errorf("Unknown object type: %s",
			q.CheckConfig.CheckConfig.ObjectType,
		), false
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("No bucketId for object found"), true
		}
		return err, false
	}
	if bid != q.CheckConfig.CheckConfig.BucketId {
		return fmt.Errorf("Object is in bucket %s, not %s",
			bid, q.CheckConfig.CheckConfig.BucketId,
		), false
	}
	return g.validateBucketInRepository(
		q.CheckConfig.CheckConfig.RepositoryId,
		q.CheckConfig.CheckConfig.BucketId,
	)
}

// Verify that the bucket is part of the specified repository
func (g *guidePost) validateBucketInRepository(
	repo, bucket string) (error, bool) {
	var repoId, repoName string
	if err := g.repo_stmt.QueryRow(bucket).Scan(
		&repoId,
		&repoName,
	); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("No repository found for bucket %s",
				bucket), true
		}
		return err, false
	}
	if repo != repoId {
		return fmt.Errorf("Bucket is in different repository: %s",
			repoId), false
	}
	return nil, false
}

// check the check configuration to contain fewer thresholds than
// the limit for the capability
func (g *guidePost) validateCheckThresholds(q *treeRequest) (error, bool) {
	var (
		thrLimit int
		err      error
	)

	if err = g.cthr_stmt.QueryRow(
		q.CheckConfig.CheckConfig.CapabilityId,
	).Scan(&thrLimit); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf(
				"Capability %s not found",
				q.CheckConfig.CheckConfig.CapabilityId), true
		}
		return err, false
	}
	if len(q.CheckConfig.CheckConfig.Thresholds) > thrLimit {
		return fmt.Errorf(
			"Specified %d thresholds exceed limit of %d for capability",
			len(q.CheckConfig.CheckConfig.Thresholds),
			thrLimit), false
	}
	return nil, false
}

// check the naming schema for the bucket (global unique object)
func (g *guidePost) validateBucketName(q *treeRequest) (error, bool) {
	_, repoName, _, _ := g.extractRouting(q)

	if !strings.HasPrefix(
		q.Bucket.Bucket.Name,
		fmt.Sprintf("%s_", repoName),
	) {
		return fmt.Errorf("Illegal bucket name format, " +
			"requires reponame_ prefix"), false
	}
	return nil, false
}

// validate current treekeeper state
func (g *guidePost) validateKeeper(repoName string) (error, bool) {
	// check we have a treekeeper for that repository
	keeper := fmt.Sprintf("repository_%s", repoName)
	if _, ok := handlerMap[keeper].(*treeKeeper); !ok {
		return fmt.Errorf(
			"No handler for repository %s currently registered.",
			repoName), true
	}
	handler := handlerMap[keeper].(*treeKeeper)

	// check the treekeeper has not been stopped
	if handler.isStopped() {
		return fmt.Errorf(
			"Repository %s is currently stopped.", repoName), false
	}

	// check the treekeeper has not encountered a broken tree
	if handler.isBroken() {
		return fmt.Errorf(
			"Repository %s is broken.", repoName), false
	}

	// check the treekeeper has finished loading
	if !handler.isReady() {
		return fmt.Errorf(
			"Repository %s not fully loaded yet.", repoName), false
	}
	return nil, false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
