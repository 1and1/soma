/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const ListTeams = `
SELECT organizational_team_id,
       organizational_team_name
FROM   inventory.organizational_teams;`

const ShowTeams = `
SELECT organizational_team_id,
       organizational_team_name,
       organizational_team_ldap_id,
       organizational_team_system
FROM   inventory.organizational_teams
WHERE  organizational_team_id = $1;`

const SyncTeams = `
SELECT organizational_team_id,
       organizational_team_name,
       organizational_team_ldap_id,
       organizational_team_system
FROM   inventory.organizational_teams
WHERE  NOT organizational_team_system;`

const TeamAdd = `
INSERT INTO inventory.organizational_teams (
            organizational_team_id,
            organizational_team_name,
            organizational_team_ldap_id,
            organizational_team_system)
SELECT $1::uuid, $2::varchar, $3::numeric, $4
WHERE  NOT EXISTS (
   SELECT organizational_team_id
   FROM   inventory.organizational_teams
   WHERE  organizational_team_id = $1::uuid
      OR  organizational_team_name = $2::varchar
      OR  organizational_team_ldap_id = $3::numeric);`

const TeamUpdate = `
UPDATE inventory.organizational_teams
SET    organizational_team_name = $1::varchar,
       organizational_team_ldap_id = $2::numeric,
       organizational_team_system = $3::boolean
WHERE  organizational_team_id = $4::uuid;`

const TeamDel = `
DELETE FROM inventory.organizational_teams
WHERE       organizational_team_id = $1;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
