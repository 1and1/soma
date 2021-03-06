/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const (
	TreeStatements = ``

	TreeShowRepository = `
SELECT repository_name,
       repository_active,
       organizational_team_id,
	   repository_deleted,
	   created_by,
	   created_at
FROM   soma.repositories
WHERE  repository_id = $1::uuid;`

	TreeShowBucket = `
SELECT bucket_name,
	   bucket_frozen,
	   bucket_deleted,
	   repository_id,
	   environment,
	   organizational_team_id,
	   created_by,
	   created_at
FROM   soma.buckets
WHERE  bucket_id = $1::uuid;`

	TreeShowGroup = `
SELECT sg.bucket_id,
       sg.group_name,
	   sg.object_state,
	   sg.organizational_team_id,
	   sg.created_by,
	   sg.created_at
FROM   soma.groups sg
WHERE  sg.group_id = $1::uuid;`

	TreeShowCluster = `
SELECT sc.cluster_name,
       sc.bucket_id,
	   sc.object_state,
	   sc.organizational_team_id,
	   sc.created_by,
	   sc.created_at
FROM   soma.clusters sc
WHERE  sc.cluster_id = $1::uuid;`

	TreeShowNode = `
SELECT sn.node_asset_id,
       sn.node_name,
       sn.organizational_team_id,
	   sn.server_id,
	   sn.object_state,
	   sn.node_online,
	   sn.node_deleted,
	   sn.created_by,
	   sn.created_at,
	   sb.repository_id,
	   snba.bucket_id
FROM   soma.nodes sn
JOIN   soma.node_bucket_assignment snba
  ON   sn.node_id = snba.node_id
  AND  sn.organizational_team_id = snba.organizational_team_id
JOIN   soma.buckets sb
  ON   snba.bucket_id = sb.bucket_id
WHERE  sn.node_id = $1::uuid;`

	//
	//
	TreeBucketsInRepository = `
SELECT bucket_id
FROM   soma.buckets
WHERE  repository_id = $1::uuid;`

	TreeGroupsInBucket = `
SELECT group_id
FROM   soma.groups
WHERE  bucket_id = $1::uuid
AND    object_state = 'standalone';`

	TreeClustersInBucket = `
SELECT cluster_id
FROM   soma.clusters
WHERE  bucket_id = $1::uuid
AND    object_state = 'standalone';`

	TreeNodesInBucket = `
SELECT snba.node_id
FROM   soma.node_bucket_assignment snba
JOIN   soma.nodes sn
  ON   snba.node_id = sn.node_id
WHERE  snba.bucket_id = $1::uuid
  AND  sn.object_state = 'standalone';`

	// groupsingroup
	TreeGroupsInGroup = `
SELECT sgmg.child_group_id
FROM   soma.group_membership_groups sgmg
WHERE  sgmg.group_id = $1::uuid;`

	// clustersingroup
	TreeClustersInGroup = `
SELECT sgmc.child_cluster_id
FROM   soma.group_membership_clusters sgmc
WHERE  sgmc.group_id = $1::uuid;`

	// nodesingroup
	TreeNodesInGroup = `
SELECT sgmn.child_node_id
FROM   soma.group_membership_nodes sgmn
WHERE  sgmn.group_id = $1::uuid;`

	// nodesincluster
	TreeNodesInCluster = `
SELECT scm.node_id
FROM   soma.cluster_membership scm
WHERE  scm.cluster_id = $1::uuid;`
)

func init() {
	m[TreeBucketsInRepository] = `TreeBucketsInRepository`
	m[TreeClustersInBucket] = `TreeClustersInBucket`
	m[TreeClustersInGroup] = `TreeClustersInGroup`
	m[TreeGroupsInBucket] = `TreeGroupsInBucket`
	m[TreeGroupsInGroup] = `TreeGroupsInGroup`
	m[TreeNodesInBucket] = `TreeNodesInBucket`
	m[TreeNodesInCluster] = `TreeNodesInCluster`
	m[TreeNodesInGroup] = `TreeNodesInGroup`
	m[TreeShowBucket] = `TreeShowBucket`
	m[TreeShowCluster] = `TreeShowCluster`
	m[TreeShowGroup] = `TreeShowGroup`
	m[TreeShowNode] = `TreeShowNode`
	m[TreeShowRepository] = `TreeShowRepository`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
