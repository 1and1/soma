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
	SupervisorPermissionStatements = ``

	LoadPermissions = `
SELECT permission_id,
       permission_name
FROM   soma.permissions;`


	PermissionAdd = `
INSERT INTO soma.permissions (
            permission_id,
            permission_name,
            category,
            created_by
)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar,
       ( SELECT user_id
         FROM   inventory.users
         WHERE  user_uid = $4::varchar)
WHERE NOT EXISTS (
      SELECT permission_name
      FROM   soma.permissions
      WHERE  permission_name = $2::varchar
);`

	PermissionLinkGrant = `
INSERT INTO soma.permission_grant_map (
            category,
            permission_id,
            granted_category,
            granted_permission_id)
SELECT $1::varchar,
       $2::uuid,
       $3::varchar,
       $4::uuid
WHERE  NOT EXISTS (
   -- a permission can not have two grant records
   SELECT permission_id
   FROM   soma.permission_grant_map
   WHERE  permission_id = $2::uuid);`

	PermissionLookupGrantId = `
SELECT permission_id
FROM   soma.permission_grant_map
WHERE  granted_permission_id = $1::uuid;`

	PermissionRevokeGlobal = `
DELETE FROM soma.authorizations_global
WHERE       permission_id = $1::uuid;`

	PermissionRevokeRepository = `
DELETE FROM soma.authorizations_repository
WHERE       permission_id = $1::uuid;`

	PermissionRevokeTeam = `
DELETE FROM soma.authorizations_team
WHERE       permission_id = $1::uuid;`

	PermissionRevokeMonitoring = `
DELETE FROM soma.authorizations_monitoring
WHERE       permission_id = $1::uuid;`

	PermissionUnmapAll = `
DELETE FROM soma.permission_map
WHERE       permission_id = $1::uuid;`

	PermissionRemove = `
DELETE FROM soma.permissions
WHERE       permission_id = $1::uuid;`

	PermissionRemoveLink = `
DELETE FROM soma.permission_grant_map
WHERE       granted_permission_id = $1::uuid;`

	PermissionList = `
SELECT permission_id,
       permission_name
FROM   soma.permissions
WHERE  category = $1::varchar;`

	PermissionShow = `
SELECT sp.permission_id,
       sp.permission_name,
       sp.category,
       iu.user_uid,
       sp.created_at
FROM   soma.permissions sp
JOIN   inventory.users iu
ON     sp.created_by = iu.user_id
WHERE  sp.permission_id = $1::uuid
  AND  sp.category = $2::varchar;`

	PermissionSearchByName = `
SELECT permission_id,
       permission_name
FROM   soma.permissions
WHERE  permission_name = $1::varchar
  AND  category= $2::varchar;`

	PermissionMappedActions = `
SELECT sa.action_id,
       sa.action_name,
       sa.section_id,
       ss.section_name,
       sa.category
FROM   soma.permissions sp
JOIN   soma.permission_map spm
  ON   sp.permission_id = spm.permission_id
JOIN   soma.sections ss
  ON   sp.section_id = ss.section_id
JOIN   soma.actions sa
  ON   sp.action_id = sa.action_id
WHERE  sp.permission_id = $1::uuid
  AND  sp.category = $2::uuid
  AND  spm.action_id IS NOT NULL
  AND  sa.section_id = sp.section_id;`

	PermissionMappedSections = `
SELECT ss.section_id,
       ss.section_name,
       ss.category
FROM   soma.permissions sp
JOIN   soma.permission_map spm
  ON   sp.permission_id = spm.permission_id
JOIN   soma.sections ss
  ON   sp.section_id = ss.section_id
WHERE  sp.permission_id = $1::uuid
  AND  sp.category = $2::uuid
  AND  spm.action_id IS NULL;`

	PermissionMapEntry = `
INSERT INTO soma.permission_map (
            mapping_id,
            category,
            permission_id,
            section_id,
            action_id)
VALUES $1::uuid,
       $2::varchar,
       $3::uuid,
       $4::uuid,
       $5::uuid;`

	PermissionUnmapEntry = `
DELETE FROM soma.permission_map
WHERE       permission_id = $1::uuid
  AND       category = $2::varchar
  AND       section_id = $3::uuid
  AND       (action_id = $4::uuid OR ($4::uuid IS NULL AND action_id IS NULL));`
)

func init() {
	m[LoadPermissions] = `LoadPermissions`
	m[PermissionAdd] = `PermissionAdd`
	m[PermissionLinkGrant] = `PermissionLinkGrant`
	m[PermissionList] = `PermissionList`
	m[PermissionLookupGrantId] = `PermissionLookupGrantId`
	m[PermissionMapEntry] = `PermissionMapEntry`
	m[PermissionMappedActions] = `PermissionMappedActions`
	m[PermissionMappedSections] = `PermissionMappedSections`
	m[PermissionRemoveLink] = `PermissionRemoveLink`
	m[PermissionRemove] = `PermissionRemove`
	m[PermissionRevokeGlobal] = `PermissionRevokeGlobal`
	m[PermissionRevokeMonitoring] = `PermissionRevokeMonitoring`
	m[PermissionRevokeRepository] = `PermissionRevokeRepository`
	m[PermissionRevokeTeam] = `PermissionRevokeTeam`
	m[PermissionSearchByName] = `PermissionSearchByName`
	m[PermissionShow] = `PermissionShow`
	m[PermissionUnmapAll] = `PermissionUnmapAll`
	m[PermissionUnmapEntry] = `PermissionUnmapEntry`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
