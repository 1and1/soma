/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const ListAllRepositories = `
SELECT repository_id,
       repository_name
FROM   soma.repositories;`

const ListScopedRepositories = `
-- direct user permissions
SELECT sr.repository_id,
       sr.repository_name
FROM   inventory.users iu
JOIN   soma.authorizations_repository sar
  ON   iu.user_id = sar.user_id
JOIN   soma.permissions sp
  ON   sar.permission_id = sp.permission_id
JOIN   soma.repositories sr
  ON   sar.repository_id = sr.repository_id
WHERE  iu.user_id = $1::uuid
  AND  sp.permission_name = $2::varchar
  AND  sr.repository_active
  AND  NOT sr.repository_deleted
UNION
-- team permissions
SELECT sr.repository_id,
       sr.repository_name
FROM   inventory.users iu
JOIN   soma.authorizations_repository sar
  ON   iu.organizational_team_id = sar.organizational_team_id
JOIN   soma.permissions sp
  ON   sar.permission_id = sp.permission_id
JOIN   soma.repositories sr
  ON   sar.repository_id = sr.repository_id
WHERE  iu.user_id = $1::uuid
  AND  sp.permission_name = $2::varchar
  AND  sr.repository_active
  AND  NOT sr.repository_deleted;`

const ShowRepository = `
SELECT repository_id,
       repository_name,
       repository_active,
       organizational_team_id
FROM   soma.repositories
WHERE  repository_id = $1
AND    NOT repository_deleted;`

const RepoOncProps = `
SELECT op.instance_id,
       op.source_instance_id,
       op.view,
       op.oncall_duty_id,
       iodt.oncall_duty_name
FROM   soma.repository_oncall_properties op
JOIN   inventory.oncall_duty_teams iodt
  ON   op.oncall_duty_id = iodt.oncall_duty_id
WHERE  op.repository_id = $1::uuid;`

const RepoSvcProps = `
SELECT sp.instance_id,
       sp.source_instance_id,
       sp.view,
       sp.service_property
FROM   soma.repository_service_properties sp
WHERE  sp.repository_id = $1::uuid;`

const RepoSysProps = `
SELECT sp.instance_id,
       sp.source_instance_id,
       sp.view,
       sp.system_property,
       sp.value
FROM   soma.repository_system_properties sp
WHERE  sp.repository_id = $1::uuid;`

const RepoCstProps = `
SELECT cp.instance_id,
       cp.source_instance_id,
       cp.view,
       cp.custom_property_id,
       cp.value,
       scp.custom_property
FROM   soma.repository_custom_properties cp
JOIN   soma.custom_properties scp
  ON   cp.custom_property_id = scp.custom_property_id
WHERE  cp.repository_id = $1::uuid;`

const RepoSystemPropertyForDelete = `
SELECT view,
       system_property,
       value
FROM   soma.repository_system_properties
WHERE  source_instance_id = $1::uuid
  AND  source_instance_id = instance_id;`

const RepoCustomPropertyForDelete = `
SELECT srcp.view,
       srcp.custom_property_id,
       srcp.value,
       scp.custom_property
FROM   soma.repository_custom_properties srcp
JOIN   soma.custom_properties scp
  ON   srcp.repository_id = scp.repository_id
 AND   srcp.custom_property_id = scp.custom_property_id
WHERE  source_instance_id = $1::uuid
  AND  source_instance_id = instance_id;`

const RepoOncallPropertyForDelete = `
SELECT srop.view,
       srop.oncall_duty_id,
       iodt.oncall_duty_name,
       iodt.oncall_duty_phone_number
FROM   soma.repository_oncall_properties srop
JOIN   inventory.oncall_duty_teams iodt
  ON   srop.oncall_duty_id = iodt.oncall_duty_id
WHERE  source_instance_id = $1::uuid
  AND  source_instance_id = instance_id;`

const RepoServicePropertyForDelete = `
SELECT srsp.view,
       srsp.service_property
FROM   soma.repository_service_properties srsp
JOIN   soma.team_service_properties stsp
  ON   srsp.organizational_team_id = stsp.organizational_team_id
 AND   srsp.service_property = stsp.service_property
WHERE  source_instance_id = $1::uuid
  AND  source_instance_id = instance_id;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
