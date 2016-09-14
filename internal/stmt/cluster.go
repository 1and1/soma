/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const ClusterBucketId = `
SELECT sc.bucket_id
FROM   soma.clusters sc
WHERE  sc.cluster_id = $1;`

const ClusterOncProps = `
SELECT op.instance_id,
       op.source_instance_id,
       op.view,
       op.oncall_duty_id,
       iodt.oncall_duty_name
FROM   soma.cluster_oncall_properties op
JOIN   inventory.oncall_duty_teams iodt
  ON   op.oncall_duty_id = iodt.oncall_duty_id
WHERE  op.cluster_id = $1::uuid;`

const ClusterSvcProps = `
SELECT sp.instance_id,
       sp.source_instance_id,
       sp.view,
       sp.service_property
FROM   soma.cluster_service_properties sp
WHERE  sp.cluster_id = $1::uuid;`

const ClusterSysProps = `
SELECT sp.instance_id,
       sp.source_instance_id,
       sp.view,
       sp.system_property,
       sp.value
FROM   soma.cluster_system_properties sp
WHERE  sp.cluster_id = $1::uuid;`

const ClusterCstProps = `
SELECT cp.instance_id,
       cp.source_instance_id,
       cp.view,
       cp.custom_property_id,
       cp.value,
       scp.custom_property
FROM   soma.cluster_custom_properties cp
JOIN   soma.custom_properties scp
  ON   cp.custom_property_id = scp.custom_property_id
WHERE  cp.cluster_id = $1::uuid;`

const ClusterSystemPropertyForDelete = `
SELECT view,
       system_property,
       value
FROM   soma.cluster_system_properties
WHERE  source_instance_id = $1::uuid
  AND  source_instance_id = instance_id;`

const ClusterCustomPropertyForDelete = `
SELECT sccp.view,
       sccp.custom_property_id,
       sccp.value,
       scp.custom_property
FROM   soma.cluster_custom_properties sccp
JOIN   soma.custom_properties scp
  ON   sccp.repository_id = scp.repository_id
 AND   sccp.custom_property_id = scp.custom_property_id
WHERE  sccp.source_instance_id = $1::uuid
  AND  sccp.source_instance_id = sccp.instance_id;`

const ClusterOncallPropertyForDelete = `
SELECT scop.view,
       scop.oncall_duty_id,
       iodt.oncall_duty_name,
       iodt.oncall_duty_phone_number
FROM   soma.cluster_oncall_properties scop
JOIN   inventory.oncall_duty_teams iodt
  ON   scop.oncall_duty_id = iodt.oncall_duty_id
WHERE  scop.source_instance_id = $1::uuid
  AND  scop.source_instance_id = scop.instance_id;`

const ClusterServicePropertyForDelete = `
SELECT scsp.view,
       scsp.service_property
FROM   soma.cluster_service_properties scsp
JOIN   soma.team_service_properties stsp
  ON   scsp.organizational_team_id = stsp.organizational_team_id
 AND   scsp.service_property = stsp.service_property
WHERE  scsp.source_instance_id = $1::uuid
  AND  scsp.source_instance_id = scsp.instance_id;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
