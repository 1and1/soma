/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const (
	PropertyStatements = ``

	ServiceLookup = `
SELECT stsp.service_property
FROM   soma.repositories sr
JOIN   soma.team_service_properties stsp
ON     sr.organizational_team_id = stsp.organizational_team_id
WHERE  sr.repository_id = $1::uuid
AND    stsp.service_property = $2::varchar
AND    sr.organizational_team_id = $3::uuid;`

	ServiceAttributes = `
SELECT stspv.service_property_attribute,
       stspv.value
FROM   soma.repositories sr
JOIN   soma.team_service_properties stsp
ON     sr.organizational_team_id = stsp.organizational_team_id
JOIN   soma.team_service_property_values stspv
ON     stsp.organizational_team_id = stspv.organizational_team_id
AND    stsp.service_property = stspv.service_property
WHERE  sr.repository_id = $1::uuid
AND    stsp.service_property = $2::varchar
AND    sr.organizational_team_id = $3::uuid;`

	PropertySystemList = `
SELECT system_property
FROM   soma.system_properties;`

	PropertyServiceList = `
SELECT service_property,
       organizational_team_id
FROM   soma.team_service_properties
WHERE  organizational_team_id = $1::uuid;`

	PropertyNativeList = `
SELECT native_property
FROM   soma.native_properties;`

	PropertyTemplateList = `
SELECT service_property
FROM   soma.service_properties;`

	PropertyCustomList = `
SELECT custom_property_id,
       repository_id,
       custom_property
FROM   soma.custom_properties
WHERE  repository_id = $1::uuid;`

	PropertySystemShow = `
SELECT system_property
FROM   soma.system_properties
WHERE  system_property = $1::varchar;`

	PropertyNativeShow = `
SELECT native_property
FROM   soma.native_properties
WHERE  native_property = $1::varchar;`

	PropertyCustomShow = `
SELECT custom_property_id,
       repository_id,
       custom_property
FROM   soma.custom_properties
WHERE  custom_property_id = $1::uuid
AND    repository_id = $2::uuid;`

	PropertyServiceShow = `
SELECT tsp.service_property,
       tsp.organizational_team_id,
       tspv.service_property_attribute,
       tspv.value
FROM   soma.team_service_properties tsp
JOIN   soma.team_service_property_values tspv
ON     tsp.service_property = tspv.service_property
WHERE  tsp.service_property = $1::varchar
AND    tsp.organizational_team_id = $2::uuid;`

	PropertyTemplateShow = `
SELECT sp.service_property,
       spv.service_property_attribute,
       spv.value
FROM   soma.service_properties sp
JOIN   soma.service_property_values spv
ON     sp.service_property = spv.service_property
WHERE  sp.service_property = $1::varchar;`

	PropertySystemAdd = `
INSERT INTO soma.system_properties (
            system_property)
SELECT $1::varchar
WHERE  NOT EXISTS (
   SELECT system_property
   FROM   soma.system_properties
   WHERE  system_property = $1::varchar);`

	PropertyNativeAdd = `
INSERT INTO soma.native_properties (
            native_property)
SELECT $1::varchar
WHERE  NOT EXISTS (
   SELECT native_property
   FROM   soma.native_properties
   WHERE  native_property = $1::varchar);`

	PropertyCustomAdd = `
INSERT INTO soma.custom_properties (
            custom_property_id,
            repository_id,
            custom_property)
SELECT $1::uuid, $2::uuid, $3::varchar
WHERE  NOT EXISTS (
   SELECT custom_property
   FROM   soma.custom_properties
   WHERE  custom_property = $3::varchar
     AND  repository_id = $2::uuid);`

	PropertyServiceAdd = `
INSERT INTO soma.team_service_properties (
            organizational_team_id,
            service_property)
SELECT $1::uuid, $2::varchar
WHERE  NOT EXISTS (
   SELECT service_property
   FROM   soma.team_service_properties
   WHERE  organizational_team_id = $1::uuid
   AND    service_property = $2::varchar);`

	PropertyServiceAttributeAdd = `
INSERT INTO soma.team_service_property_values (
            organizational_team_id,
            service_property,
            service_property_attribute,
            value)
SELECT $1::uuid, $2::varchar, $3::varchar, $4::varchar;`

	PropertyTemplateAdd = `
INSERT INTO soma.service_properties (
            service_property)
SELECT $1::varchar
WHERE  NOT EXISTS (
   SELECT service_property
   FROM   soma.service_properties
   WHERE  service_property = $1::varchar);`

	PropertyTemplateAttributeAdd = `
INSERT INTO soma.service_property_values (
            service_property,
            service_property_attribute,
            value)
SELECT $1::varchar, $2::varchar, $3::varchar;`

	PropertySystemDel = `
DELETE FROM soma.system_properties
WHERE  system_property = $1::varchar;`

	PropertyNativeDel = `
DELETE FROM soma.native_properties
WHERE  native_property = $1::varchar;`

	PropertyCustomDel = `
DELETE FROM soma.custom_properties
WHERE  repository_id = $1::uuid
AND    custom_property_id = $2::uuid;`

	PropertyServiceDel = `
DELETE FROM soma.team_service_properties
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`

	PropertyServiceAttributeDel = `
DELETE FROM soma.team_service_property_values
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`

	PropertyTemplateDel = `
DELETE FROM soma.service_properties
WHERE  service_property = $1::varchar;`

	PropertyTemplateAttributeDel = `
DELETE FROM soma.service_property_values
WHERE  service_property = $1::varchar;`
)

func init() {
	m[PropertyCustomAdd] = `PropertyCustomAdd`
	m[PropertyCustomDel] = `PropertyCustomDel`
	m[PropertyCustomList] = `PropertyCustomList`
	m[PropertyCustomShow] = `PropertyCustomShow`
	m[PropertyNativeAdd] = `PropertyNativeAdd`
	m[PropertyNativeDel] = `PropertyNativeDel`
	m[PropertyNativeList] = `PropertyNativeList`
	m[PropertyNativeShow] = `PropertyNativeShow`
	m[PropertyServiceAdd] = `PropertyServiceAdd`
	m[PropertyServiceAttributeAdd] = `PropertyServiceAttributeAdd`
	m[PropertyServiceAttributeDel] = `PropertyServiceAttributeDel`
	m[PropertyServiceDel] = `PropertyServiceDel`
	m[PropertyServiceList] = `PropertyServiceList`
	m[PropertyServiceShow] = `PropertyServiceShow`
	m[PropertySystemAdd] = `PropertySystemAdd`
	m[PropertySystemDel] = `PropertySystemDel`
	m[PropertySystemList] = `PropertySystemList`
	m[PropertySystemShow] = `PropertySystemShow`
	m[PropertyTemplateAdd] = `PropertyTemplateAdd`
	m[PropertyTemplateAttributeAdd] = `PropertyTemplateAttributeAdd`
	m[PropertyTemplateAttributeDel] = `PropertyTemplateAttributeDel`
	m[PropertyTemplateDel] = `PropertyTemplateDel`
	m[PropertyTemplateList] = `PropertyTemplateList`
	m[PropertyTemplateShow] = `PropertyTemplateShow`
	m[ServiceAttributes] = `ServiceAttributes`
	m[ServiceLookup] = `ServiceLookup`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
