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
	SectionStatements = ``

	SectionList = `
SELECT section_id,
       section_name
FROM   soma.sections;`

	SectionSearch = `
SELECT section_id,
       section_name
FROM   soma.sections
WHERE  section_name = $1::varchar;`

	SectionShow = `
SELECT ss.section_id,
       ss.section_name,
       ss.category,
       iu.user_uid,
       ss.created_at
FROM   soma.sections ss
JOIN   inventory.users iu
  ON   ss.created_by = iu.user_id
WHERE  ss.section_id = $1::uuid;`

	SectionRemoveFromMap = `
DELETE FROM soma.permission_map
WHERE       section_id = $1::uuid
  AND       action_id IS NULL;`

	SectionRemove = `
DELETE FROM soma.sections
WHERE       section_id = $1::uuid;`

	SectionListActions = `
SELECT action_id
FROM   soma.actions
WHERE  section_id = $1::uuid;`

	SectionAdd = `
INSERT INTO soma.sections (
            section_id,
            section_name,
            category,
            created_by)
SELECT      $1::uuid,
            $2::varchar,
            $3::varchar,
            ( SELECT user_id
              FROM   inventory.users
              WHERE  user_uid = $4::varchar)
WHERE       NOT EXISTS (
     SELECT section_id
     FROM   soma.sections
     WHERE  section_name = $2::varchar);`
)

func init() {
	m[SectionAdd] = `SectionAdd`
	m[SectionListActions] = `SectionListActions`
	m[SectionList] = `SectionList`
	m[SectionRemoveFromMap] = `SectionRemoveFromMap`
	m[SectionRemove] = `SectionRemove`
	m[SectionSearch] = `SectionSearch`
	m[SectionShow] = `SectionShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
