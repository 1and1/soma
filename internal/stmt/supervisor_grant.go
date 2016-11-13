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
	SupervisorGrantStatements = ``

	RevokeGlobalAuthorization = `
DELETE FROM soma.authorizations_global
WHERE       grant_id = $1::uuid
  AND       permission_id = $2::uuid
  AND       category = $3::varchar;`

	RevokeRepositoryAuthorization = `
DELETE FROM soma.authorizations_repository
WHERE       grant_id = $1::uuid
  AND       permission_id = $2::uuid
  AND       category = $3::varchar;`

	RevokeTeamAuthorization = `
DELETE FROM soma.authorizations_team
WHERE       grant_id = $1::uuid
  AND       permission_id = $2::uuid
  AND       category = $3::varchar;`

	RevokeMonitoringAuthorization = `
DELETE FROM soma.authorizations_monitoring
WHERE       grant_id = $1::uuid
  AND       permission_id = $2::uuid
  AND       category = $3::varchar;`

	GrantGlobalAuthorization = `
INSERT INTO soma.authorizations_global (
            grant_id,
            admin_id,
            user_id,
            tool_id,
            organizational_team_id,
            permission_id,
            category,
            created_by)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::uuid,
       $5::uuid,
       $6::uuid,
       $7::varchar,
       iu.user_id
FROM   inventory.users iu
WHERE  iu.user_uid = $8::varchar;`

	GrantRepositoryAuthorization = `
INSERT INTO soma.authorizations_repository (
            grant_id,
            user_id,
            tool_id,
            organizational_team_id,
            category,
            permission_id,
            object_type,
            repository_id,
            bucket_id,
            group_id,
            cluster_id,
            node_id,
            created_by)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::uuid,
       $5::varchar,
       $6::uuid,
       $7::varchar,
       $8::uuid,
       $9::uuid,
       $10::uuid,
       $11::uuid,
       $12::uuid,
       iu.user_id
FROM   inventory.users iu
WHERE  iu.user_uid = $13::varchar;`

	GrantTeamAuthorization = `
INSERT INTO soma.authorizations_team (
            grant_id,
            user_id,
            tool_id,
            organizational_team_id,
            category,
            permission_id,
            authorized_team_id,
            created_by)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::uuid,
       $5::varchar,
       $6::uuid,
       $7::uuid,
       iu.user_id
FROM   inventory.users iu
WHERE  iu.user_uid = $8::varchar;`

	GrantMonitoringAuthorization = `
INSERT INTO soma.authorizations_monitoring (
            grant_id,
            user_id,
            tool_id,
            organizational_team_id,
            category,
            permission_id,
            monitoring_id,
            created_by)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::uuid,
       $5::varchar,
       $6::uuid,
       $7::uuid,
       iu.user_id
FROM   inventory.users iu
WHERE  iu.user_uid = $8::varchar;`

	LoadGlobalOrSystemUserGrants = `
SELECT grant_id,
       user_id,
       permission_id
FROM   soma.authorizations_global;`

	SearchGlobalSystemGrant = `
SELECT grant_id
FROM   soma.authorizations_global
WHERE  permission_id = $1::uuid
  AND  permission_type = $2::varchar
  AND  (   admin_id = $3::uuid
        OR user_id  = $3::uuid
		OR tool_id  = $3::uuid);`
)

func init() {
	m[LoadGlobalOrSystemUserGrants] = `LoadGlobalOrSystemUserGrants`
	m[SearchGlobalSystemGrant] = `SearchGlobalSystemGrant`
	m[GrantGlobalAuthorization] = `GrantGlobalAuthorization`
	m[GrantMonitoringAuthorization] = `GrantMonitoringAuthorization`
	m[GrantRepositoryAuthorization] = `GrantRepositoryAuthorization`
	m[GrantTeamAuthorization] = `GrantTeamAuthorization`
	m[RevokeGlobalAuthorization] = `RevokeGlobalAuthorization`
	m[RevokeMonitoringAuthorization] = `RevokeMonitoringAuthorization`
	m[RevokeRepositoryAuthorization] = `RevokeRepositoryAuthorization`
	m[RevokeTeamAuthorization] = `RevokeTeamAuthorization`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
