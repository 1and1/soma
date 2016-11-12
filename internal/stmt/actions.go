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
	ActionStatements = ``

	ActionList = `
SELECT action_id,
       action_name,
       section_id
FROM   soma.actions;`

	ActionSearch = `
SELECT action_id,
       action_name,
       section_id
FROM   soma.actions
WHERE  action_name = $1::varchar
  AND  section_id = $2::uuid;`

	ActionShow = `
SELECT sa.action_id,
       sa.action_name,
       sa.section_id,
       sa.category,
       iu.user_uid,
       sa.created_at
FROM   soma.actions sa
JOIN   inventory.users iu
  ON   sa.created_by = iu.user_id
WHERE  sa.action_id = $1::uuid;`

	ActionRemoveFromMap = `
DELETE FROM soma.permission_map
WHERE       action_id = $1::uuid;`

	ActionRemove = `
DELETE FROM soma.actions
WHERE       action_id = $1::uuid;`

	ActionAdd = `
INSERT INTO soma.actions (
            action_id,
            action_name,
            section_id,
            category,
            created_by)
SELECT      $1::uuid,
            $2::varchar,
            $3::uuid,
            ( SELECT category
              FROM   soma.sections
              WHERE  section_id = $3::uuid),
            ( SELECT user_id
              FROM   inventory.users
              WHERE  user_uid = $4::varchar)
WHERE       NOT EXISTS (
     SELECT action_id
     FROM   soma.actions
     WHERE  action_name = $2::varchar
     AND    section_id = $3::uuid);`
)

func init() {
	m[ActionAdd] = `ActionAdd`
	m[ActionList] = `ActionList`
	m[ActionRemoveFromMap] = `ActionRemoveFromMap`
	m[ActionRemove] = `ActionRemove`
	m[ActionSearch] = `ActionSearch`
	m[ActionShow] = `ActionShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
