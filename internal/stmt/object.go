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
	ObjectStatements = ``

	ObjectStateList = `
SELECT object_state
FROM   soma.object_states;`

	ObjectStateShow = `
SELECT object_state
FROM   soma.object_states
WHERE  object_state = $1::varchar;`

	ObjectStateAdd = `
INSERT INTO soma.object_states (
            object_state)
SELECT $1::varchar
WHERE  NOT EXISTS (
   SELECT object_state
   FROM   soma.object_states
   WHERE  object_state = $1::varchar);`

	ObjectStateRemove = `
DELETE FROM soma.object_states
WHERE       object_state = $1::varchar;`

	ObjectStateRename = `
UPDATE soma.object_states
SET    object_state = $1::varchar
WHERE  object_state = $2::varchar;`

	ObjectTypeList = `
SELECT object_type
FROM   soma.object_types;`

	ObjectTypeShow = `
SELECT object_type
FROM   soma.object_types
WHERE  object_type = $1::varchar;`

	ObjectTypeAdd = `
INSERT INTO soma.object_types (
            object_type)
SELECT $1::varchar
WHERE NOT EXISTS (
   SELECT object_type
   FROM   soma.object_types
   WHERE  object_type = $1::varchar);`

	ObjectTypeDel = `
DELETE FROM soma.object_types
WHERE       object_type = $1::varchar;`

	ObjectTypeRename = `
UPDATE soma.object_types
SET    object_type = $1::varchar
WHERE  object_type = $2::varchar;`
)

func init() {
	m[ObjectStateAdd] = `ObjectStateAdd`
	m[ObjectStateList] = `ObjectStateList`
	m[ObjectStateRemove] = `ObjectStateRemove`
	m[ObjectStateRename] = `ObjectStateRename`
	m[ObjectStateShow] = `ObjectStateShow`
	m[ObjectTypeAdd] = `ObjectTypeAdd`
	m[ObjectTypeDel] = `ObjectTypeDel`
	m[ObjectTypeList] = `ObjectTypeList`
	m[ObjectTypeRename] = `ObjectTypeRename`
	m[ObjectTypeShow] = `ObjectTypeShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
