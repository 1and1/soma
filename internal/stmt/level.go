/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const LevelList = `
SELECT level_name,
       level_shortname
FROM   soma.notification_levels;`

const LevelShow = `
SELECT level_name,
       level_shortname,
       level_numeric
FROM   soma.notification_levels
WHERE  level_name = $1;`

const LevelAdd = `
INSERT INTO soma.notification_levels (
            level_name,
            level_shortname,
            level_numeric)
SELECT $1::varchar, $2::varchar, $3::smallint
WHERE  NOT EXISTS (
   SELECT level_name
   FROM   soma.notification_levels
   WHERE  level_name = $1::varchar
      OR  level_shortname = $2::varchar
      OR  level_numeric = $3::smallint);`

const LevelDel = `
DELETE FROM soma.notification_levels
WHERE  level_name = $1;`

func init() {
	m[LevelAdd] = `LevelAdd`
	m[LevelDel] = `LevelDel`
	m[LevelList] = `LevelList`
	m[LevelShow] = `LevelShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
