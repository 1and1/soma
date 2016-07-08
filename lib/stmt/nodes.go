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

const NodeOncProps = `
SELECT op.instance_id,
       op.source_instance_id,
       op.view,
       op.oncall_duty_id,
       iodt.oncall_duty_name
FROM   soma.node_oncall_property op
JOIN   inventory.oncall_duty_teams iodt
  ON   op.oncall_duty_id = iodt.oncall_duty_id
WHERE  op.node_id = $1::uuid;`

const NodeSvcProps = `
SELECT sp.instance_id,
       sp.source_instance_id,
       sp.view,
       sp.service_property
FROM   soma.node_service_properties sp
WHERE  sp.node_id = $1::uuid;`

const NodeSysProps = `
SELECT sp.instance_id,
       sp.source_instance_id,
       sp.view,
       sp.system_property,
       sp.value
FROM   soma.node_system_properties sp
WHERE  sp.node_id = $1::uuid;`

const NodeCstProps = `
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

const NodeSystemPropertyForDelete = `
SELECT snsp.view,
       snsp.system_property,
       snsp.value
FROM   soma.node_system_properties snsp
WHERE  snsp.source_instance_id = $1::uuid
  AND  snsp.source_instance_id = snsp.instance_id;`

const NodeCustomPropertyForDelete = `
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

const NodeOncallPropertyForDelete = `
SELECT snop.view,
       snop.oncall_duty_id,
       iodt.oncall_duty_name,
       iodt.oncall_duty_phone_number
FROM   soma.node_oncall_property snop
JOIN   inventory.oncall_duty_teams iodt
  ON   snop.oncall_duty_id = iodt.oncall_duty_id
WHERE  snop.source_instance_id = $1::uuid
  AND  snop.source_instance_id = snop.instance_id;`

const NodeServicePropertyForDelete = `
SELECT snsp.view,
       snsp.service_property
FROM   soma.node_service_properties snsp
JOIN   soma.team_service_properties stsp
  ON   snsp.organizational_team_id = stsp.organizational_team_id
 AND   snsp.service_property = stsp.service_property
WHERE  snsp.source_instance_id = $1::uuid
  AND  snsp.source_instance_id = snsp.instance_id;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
