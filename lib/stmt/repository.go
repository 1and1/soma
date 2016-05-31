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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
