/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const ValidityList = `
SELECT system_property,
       object_type
FROM   soma.system_property_validity;`

const ValidityShow = `
SELECT system_property,
       object_type,
       inherited
FROM   soma.system_property_validity
WHERE  system_property = $1;`

const ValidityAdd = `
INSERT INTO soma.system_property_validity (
            system_property,
            object_type,
            inherited)
SELECT $1::varchar,
       $2::varchar,
       $3::boolean
WHERE  NOT EXISTS (
    SELECT system_property,
           object_type
    FROM   soma.system_property_validity
    WHERE  system_property = $1::varchar
    AND    object_type = $2::varchar
    AND    inherited = $3::boolean);`

const ValidityDel = `
DELETE FROM soma.system_property_validity
WHERE       system_property = $1::varchar;`

func init() {
	m[ValidityAdd] = `ValidityAdd`
	m[ValidityDel] = `ValidityDel`
	m[ValidityList] = `ValidityList`
	m[ValidityShow] = `ValidityShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
