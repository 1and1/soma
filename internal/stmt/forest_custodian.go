/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const ForestRebuildDeleteChecks = `
UPDATE soma.checks sc
SET    deleted = 'yes'::boolean
WHERE  sc.repository_id = $1::uuid;`

const ForestRebuildDeleteInstances = `
UPDATE soma.check_instances sci
SET    deleted = 'yes'::boolean
FROM   soma.checks sc
WHERE  sci.check_id = sc.check_id
AND    sc.repository_id = $1::uuid;`

const ForestRepoNameById = `
SELECT repository_name,
       organizational_team_id
FROM   soma.repositories
WHERE  repository_id = $1::uuid;`

const ForestLoadRepository = `
SELECT repository_id,
       repository_name,
       repository_deleted,
       repository_active,
       organizational_team_id
FROM   soma.repositories;`

const ForestAddRepository = `
INSERT INTO soma.repositories (
            repository_id,
            repository_name,
            repository_active,
            repository_deleted,
            organizational_team_id,
            created_by)
SELECT      $1::uuid,
            $2::varchar,
            $3::boolean,
            $4::boolean,
            $5::uuid,
            user_id
FROM        inventory.users iu
WHERE       iu.user_uid = $6::varchar
AND NOT EXISTS (
	SELECT  repository_id
	FROM    soma.repositories
	WHERE   repository_id   = $1::uuid
	OR      repository_name = $2::varchar);`

func init() {
	m[ForestAddRepository] = `ForestAddRepository`
	m[ForestLoadRepository] = `ForestLoadRepository`
	m[ForestRebuildDeleteChecks] = `ForestRebuildDeleteChecks`
	m[ForestRebuildDeleteInstances] = `ForestRebuildDeleteInstances`
	m[ForestRepoNameById] = `ForestRepoNameById`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
