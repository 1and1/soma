/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const ObjectStateList = `
SELECT object_state
FROM   soma.object_states;`

const ObjectStateShow = `
SELECT object_state
FROM   soma.object_states
WHERE  object_state = $1::varchar;`

const ObjectStateAdd = `
INSERT INTO soma.object_states (
            object_state)
SELECT $1::varchar
WHERE  NOT EXISTS (
   SELECT object_state
   FROM   soma.object_states
   WHERE  object_state = $1::varchar);`

const ObjectStateDel = `
DELETE FROM soma.object_states
WHERE       object_state = $1::varchar;`

const ObjectStateRename = `
UPDATE soma.object_states
SET    object_state = $1::varchar
WHERE  object_state = $2::varchar;`

const ObjectTypeList = `
SELECT object_type
FROM   soma.object_types;`

const ObjectTypeShow = `
SELECT object_type
FROM   soma.object_types
WHERE  object_type = $1::varchar;`

const ObjectTypeAdd = `
INSERT INTO soma.object_types (
            object_type)
SELECT $1::varchar
WHERE NOT EXISTS (
   SELECT object_type
   FROM   soma.object_types
   WHERE  object_type = $1::varchar);`

const ObjectTypeDel = `
DELETE FROM soma.object_types
WHERE       object_type = $1::varchar;
  `

const ObjectTypeRename = `
UPDATE soma.object_types
SET    object_type = $1::varchar
WHERE  object_type = $2::varchar;`

func init() {
	m[ObjectStateAdd] = `ObjectStateAdd`
	m[ObjectStateDel] = `ObjectStateDel`
	m[ObjectStateList] = `ObjectStateList`
	m[ObjectStateRename] = `ObjectStateRename`
	m[ObjectStateShow] = `ObjectStateShow`
	m[ObjectTypeAdd] = `ObjectTypeAdd`
	m[ObjectTypeDel] = `ObjectTypeDel`
	m[ObjectTypeList] = `ObjectTypeList`
	m[ObjectTypeRename] = `ObjectTypeRename`
	m[ObjectTypeShow] = `ObjectTypeShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
