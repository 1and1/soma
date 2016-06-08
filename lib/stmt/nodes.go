/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const ListNodes = `
SELECT node_id,
       node_name
FROM   soma.nodes
WHERE  node_online;`

const ShowNodes = `
SELECT node_id,
       node_asset_id,
       node_name,
       organizational_team_id,
       server_id,
       object_state,
       node_online,
       node_deleted
FROM   soma.nodes
WHERE  node_id = $1;`

const ShowConfigNodes = `
SELECT nodes.node_id,
       nodes.node_name,
       buckets.bucket_id,
       buckets.repository_id
FROM   soma.nodes
JOIN   soma.node_bucket_assignment
  ON   nodes.node_id = node_bucket_assignment.node_id
JOIN   soma.buckets
  ON   node_bucket_assignment.bucket_id = buckets.bucket_id
WHERE  nodes.node_id = $1;`

const SyncNodes = `
SELECT node_id,
       node_asset_id,
       node_name,
       organizational_team_id,
       server_id,
       node_online,
       node_deleted
FROM   soma.nodes;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
