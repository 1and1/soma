/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const (
	OncallStatements = ``

	OncallList = `
SELECT oncall_duty_id,
       oncall_duty_name
FROM   inventory.oncall_duty_teams;`

	OncallShow = `
SELECT oncall_duty_id,
       oncall_duty_name,
       oncall_duty_phone_number
FROM   inventory.oncall_duty_teams
WHERE  oncall_duty_id = $1::varchar;`

	OncallAdd = `
INSERT INTO inventory.oncall_duty_teams (
            oncall_duty_id,
            oncall_duty_name,
            oncall_duty_phone_number)
SELECT $1::uuid, $2::varchar, $3::numeric
WHERE  NOT EXISTS (
   SELECT oncall_duty_id
   FROM   inventory.oncall_duty_teams
   WHERE  oncall_duty_id = $1::uuid
      OR  oncall_duty_name = $2::varchar
      OR  oncall_duty_phone_number = $3::numeric);`

	OncallUpdate = `
UPDATE inventory.oncall_duty_teams
SET    oncall_duty_name = CASE WHEN $1::varchar IS NOT NULL
                          THEN      $1::varchar
                          ELSE      oncall_duty_name
                          END,
       oncall_duty_phone_number = CASE WHEN $2::numeric IS NOT NULL
                                  THEN      $2::numeric
                                  ELSE      oncall_duty_phone_number
                                  END
WHERE  oncall_duty_id = $3::uuid;`

	OncallDel = `
DELETE FROM inventory.oncall_duty_teams
WHERE  oncall_duty_id = $1::varchar;`
)

func init() {
	m[OncallAdd] = `OncallAdd`
	m[OncallDel] = `OncallDel`
	m[OncallList] = `OncallList`
	m[OncallShow] = `OncallShow`
	m[OncallUpdate] = `OncallUpdate`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
