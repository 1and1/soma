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

	AddPermissionCategory = `
INSERT INTO soma.categories (
            category,
            created_by
)
SELECT $1::varchar,
       $2::uuid
WHERE NOT EXISTS (
      SELECT category
      FROM   soma.categories
      WHERE  category = $1::varchar
);`

	DeletePermissionCategory = `
DELETE FROM soma.categories
WHERE category = $1::varchar;`

	ListPermissionCategory = `
SELECT category
FROM   soma.categories;`

	ShowPermissionCategory = `
SELECT sc.category,
       iu.user_uid,
       sc.created_at
FROM   soma.categories sc
JOIN   inventory.users iu
ON     sc.created_by = iu.user_id
WHERE  sc.category = $1::varchar;`

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

	ListPermission = `
SELECT permission_id,
       permission_name
FROM   soma.permissions;`

	ShowPermission = `
SELECT sp.permission_id,
       sp.permission_name,
       sp.category,
       iu.user_uid,
       sp.created_at
FROM   soma.permissions sp
JOIN   inventory.users iu
ON     sp.created_by = iu.user_id
WHERE  sp.permission_name = $1::varchar;`

	SearchPermissionByName = `
SELECT permission_id,
       permission_name
FROM   soma.permissions
WHERE  permission_name = $1::varchar;`
)

func init() {
	m[AddPermissionCategory] = `AddPermissionCategory`
	m[DeletePermissionCategory] = `DeletePermissionCategory`
	m[ListPermissionCategory] = `ListPermissionCategory`
	m[ListPermission] = `ListPermission`
	m[LoadPermissions] = `LoadPermissions`
	m[PermissionAdd] = `PermissionAdd`
	m[PermissionLinkGrant] = `PermissionLinkGrant`
	m[PermissionLookupGrantId] = `PermissionLookupGrantId`
	m[PermissionRemoveLink] = `PermissionRemoveLink`
	m[PermissionRemove] = `PermissionRemove`
	m[PermissionRevokeGlobal] = `PermissionRevokeGlobal`
	m[PermissionRevokeMonitoring] = `PermissionRevokeMonitoring`
	m[PermissionRevokeRepository] = `PermissionRevokeRepository`
	m[PermissionRevokeTeam] = `PermissionRevokeTeam`
	m[PermissionUnmapAll] = `PermissionUnmapAll`
	m[SearchPermissionByName] = `SearchPermissionByName`
	m[ShowPermissionCategory] = `ShowPermissionCategory`
	m[ShowPermission] = `ShowPermission`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
