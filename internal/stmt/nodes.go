/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const (
	NodeStatements = ``

	ListNodes = `
SELECT node_id,
       node_name
FROM   soma.nodes
WHERE  node_online;`

	ShowNodes = `
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

	ShowConfigNodes = `
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

	NodeBucketId = `
SELECT snba.bucket_id
FROM   soma.node_bucket_assignment snba
WHERE  snba.node_id = $1;`

	SyncNodes = `
SELECT node_id,
       node_asset_id,
       node_name,
       organizational_team_id,
       server_id,
       node_online,
       node_deleted
FROM   soma.nodes;`

	NodeOncProps = `
SELECT op.instance_id,
       op.source_instance_id,
       op.view,
       op.oncall_duty_id,
       iodt.oncall_duty_name
FROM   soma.node_oncall_property op
JOIN   inventory.oncall_duty_teams iodt
  ON   op.oncall_duty_id = iodt.oncall_duty_id
WHERE  op.node_id = $1::uuid;`

	NodeSvcProps = `
SELECT sp.instance_id,
       sp.source_instance_id,
       sp.view,
       sp.service_property
FROM   soma.node_service_properties sp
WHERE  sp.node_id = $1::uuid;`

	NodeSysProps = `
SELECT sp.instance_id,
       sp.source_instance_id,
       sp.view,
       sp.system_property,
       sp.value
FROM   soma.node_system_properties sp
WHERE  sp.node_id = $1::uuid;`

	NodeCstProps = `
SELECT cp.instance_id,
       cp.source_instance_id,
       cp.view,
       cp.custom_property_id,
       cp.value,
       scp.custom_property
FROM   soma.node_custom_properties cp
JOIN   soma.custom_properties scp
  ON   cp.custom_property_id = scp.custom_property_id
WHERE  cp.node_id = $1::uuid;`

	NodeSystemPropertyForDelete = `
SELECT snsp.view,
       snsp.system_property,
       snsp.value
FROM   soma.node_system_properties snsp
WHERE  snsp.source_instance_id = $1::uuid
  AND  snsp.source_instance_id = snsp.instance_id;`

	NodeCustomPropertyForDelete = `
SELECT sncp.view,
       sncp.custom_property_id,
       sncp.value,
       scp.custom_property
FROM   soma.node_custom_properties sncp
JOIN   soma.custom_properties scp
  ON   sncp.repository_id = scp.repository_id
 AND   sncp.custom_property_id = scp.custom_property_id
WHERE  sncp.source_instance_id = $1::uuid
  AND  sncp.source_instance_id = sncp.instance_id;`

	NodeOncallPropertyForDelete = `
SELECT snop.view,
       snop.oncall_duty_id,
       iodt.oncall_duty_name,
       iodt.oncall_duty_phone_number
FROM   soma.node_oncall_property snop
JOIN   inventory.oncall_duty_teams iodt
  ON   snop.oncall_duty_id = iodt.oncall_duty_id
WHERE  snop.source_instance_id = $1::uuid
  AND  snop.source_instance_id = snop.instance_id;`

	NodeServicePropertyForDelete = `
SELECT snsp.view,
       snsp.service_property
FROM   soma.node_service_properties snsp
JOIN   soma.team_service_properties stsp
  ON   snsp.organizational_team_id = stsp.organizational_team_id
 AND   snsp.service_property = stsp.service_property
WHERE  snsp.source_instance_id = $1::uuid
  AND  snsp.source_instance_id = snsp.instance_id;`

	NodeDetails = `
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
AND       sn.node_id = $1::uuid;`

	NodeAdd = `
INSERT INTO soma.nodes (
            node_id,
            node_asset_id,
            node_name,
            organizational_team_id,
            server_id,
            object_state,
            node_online,
            node_deleted,
            created_by)
SELECT $1::uuid,
       $2::numeric,
       $3::varchar,
       $4,
       $5,
       $6,
       $7,
       $8,
       user_id
FROM   inventory.users iu
WHERE  iu.user_uid = $9::varchar
AND    NOT EXISTS (
         SELECT node_id
         FROM   soma.nodes
         WHERE  node_id = $1::uuid
         OR     node_asset_id = $2::numeric
         OR     (node_name = $3::varchar AND node_online)
       );`

	NodeUpdate = `
UPDATE soma.nodes
SET    node_asset_id = $1::numeric,
       node_name = $2::varchar,
       organizational_team_id = $3::uuid,
       server_id = $4::uuid,
       node_online = $5::boolean,
       node_deleted = $6::boolean
WHERE  node_id = $7::uuid;`

	NodeDel = `
UPDATE soma.nodes
SET    node_deleted = 'yes'
WHERE  node_id = $1
AND    node_deleted = 'no';`

	NodePurge = `
DELETE FROM soma.nodes
WHERE       node_id = $1
AND         node_deleted;`
)

func init() {
	m[ListNodes] = `ListNodes`
	m[NodeAdd] = `NodeAdd`
	m[NodeBucketId] = `NodeBucketId`
	m[NodeCstProps] = `NodeCstProps`
	m[NodeCustomPropertyForDelete] = `NodeCustomPropertyForDelete`
	m[NodeDel] = `NodeDel`
	m[NodeDetails] = `NodeDetails`
	m[NodeOncProps] = `NodeOncProps`
	m[NodeOncallPropertyForDelete] = `NodeOncallPropertyForDelete`
	m[NodePurge] = `NodePurge`
	m[NodeServicePropertyForDelete] = `NodeServicePropertyForDelete`
	m[NodeSvcProps] = `NodeSvcProps`
	m[NodeSysProps] = `NodeSysProps`
	m[NodeSystemPropertyForDelete] = `NodeSystemPropertyForDelete`
	m[NodeUpdate] = `NodeUpdate`
	m[ShowConfigNodes] = `ShowConfigNodes`
	m[ShowNodes] = `ShowNodes`
	m[SyncNodes] = `SyncNodes`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
