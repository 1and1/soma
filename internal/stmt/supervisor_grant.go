/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const (
	SupervisorGrantStatements = ``

	GrantGlobalOrSystemToUser = `
INSERT INTO soma.authorizations_global (
    grant_id,
    user_id,
    permission_id,
    permission_type,
    created_by
)
VALUES (
    $1::uuid,
    $2::uuid,
    $3::uuid,
    $4::varchar,
    $5::uuid
);`

	RevokeGlobalOrSystemFromUser = `
DELETE FROM soma.authorizations_global
WHERE grant_id = $1::uuid;`

	LoadGlobalOrSystemUserGrants = `
SELECT grant_id,
       user_id,
       permission_id
FROM   soma.authorizations_global;`

	GrantLimitedRepoToUser = `
INSERT INTO soma.authorizations_repository (
	grant_id,
	user_id,
	repository_id,
	permission_id,
	permission_type,
	created_by
)
VALUES (
	$1::uuid,
	$2::uuid,
	$3::uuid,
	$4::uuid,
	$5::varchar,
	$6::uuid
);`

	RevokeLimitedRepoFromUser = `
DELETE FROM soma.authorizations_repository
WHERE grant_id = $1::uuid;`

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
	m[GrantGlobalOrSystemToUser] = `GrantGlobalOrSystemToUser`
	m[GrantLimitedRepoToUser] = `GrantLimitedRepoToUser`
	m[LoadGlobalOrSystemUserGrants] = `LoadGlobalOrSystemUserGrants`
	m[RevokeGlobalOrSystemFromUser] = `RevokeGlobalOrSystemFromUser`
	m[RevokeLimitedRepoFromUser] = `RevokeLimitedRepoFromUser`
	m[SearchGlobalSystemGrant] = `SearchGlobalSystemGrant`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
