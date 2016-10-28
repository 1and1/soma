/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const LoadUserTeamMapping = `
SELECT iu.user_id,
       iu.user_uid,
       iu.organizational_team_id,
       iot.organizational_team_name
FROM   inventory.users iu
JOIN   inventory.organizational_teams iot
ON     iu.organizational_team_id = iot.organizational_team_id;`

func init() {
	m[LoadUserTeamMapping] = `LoadUserTeamMapping`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
