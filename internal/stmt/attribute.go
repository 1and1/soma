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
	AttributeStatements = ``

	AttributeList = `
SELECT service_property_attribute,
       cardinality
FROM   soma.service_property_attributes;`

	AttributeShow = `
SELECT service_property_attribute,
       cardinality
FROM   soma.service_property_attributes
WHERE  service_property_attribute = $1::varchar;`

	AttributeAdd = `
INSERT INTO soma.service_property_attributes (
            service_property_attribute,
            cardinality)
SELECT $1::varchar, $2::varchar WHERE NOT EXISTS (
       SELECT service_property_attribute
       FROM   soma.service_property_attributes
       WHERE  service_property_attribute = $1::varchar);`

	AttributeDelete = `
DELETE FROM soma.service_property_attributes
WHERE       service_property_attribute = $1::varchar;`
)

func init() {
	m[AttributeAdd] = `AttributeAdd`
	m[AttributeDelete] = `AttributeDelete`
	m[AttributeList] = `AttributeList`
	m[AttributeShow] = `AttributeShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
