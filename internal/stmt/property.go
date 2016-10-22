/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const ServiceLookup = `
SELECT stsp.service_property
FROM   soma.repositories sr
JOIN   soma.team_service_properties stsp
ON     sr.organizational_team_id = stsp.organizational_team_id
WHERE  sr.repository_id = $1::uuid
AND    stsp.service_property = $2::varchar
AND    sr.organizational_team_id = $3::uuid;`

const ServiceAttributes = `
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

const PropertySystemList = `
SELECT system_property
FROM   soma.system_properties;`

const PropertyServiceList = `
SELECT service_property,
       organizational_team_id
FROM   soma.team_service_properties
WHERE  organizational_team_id = $1::uuid;`

const PropertyNativeList = `
SELECT native_property
FROM   soma.native_properties;`

const PropertyTemplateList = `
SELECT service_property
FROM   soma.service_properties;`

const PropertyCustomList = `
SELECT custom_property_id,
       repository_id,
       custom_property
FROM   soma.custom_properties
WHERE  repository_id = $1::uuid;`

const PropertySystemShow = `
SELECT system_property
FROM   soma.system_properties
WHERE  system_property = $1::varchar;`

const PropertyNativeShow = `
SELECT native_property
FROM   soma.native_properties
WHERE  native_property = $1::varchar;`

const PropertyCustomShow = `
SELECT custom_property_id,
       repository_id,
       custom_property
FROM   soma.custom_properties
WHERE  custom_property_id = $1::uuid
AND    repository_id = $2::uuid;`

const PropertyServiceShow = `
SELECT tsp.service_property,
       tsp.organizational_team_id,
       tspv.service_property_attribute,
       tspv.value
FROM   soma.team_service_properties tsp
JOIN   soma.team_service_property_values tspv
ON     tsp.service_property = tspv.service_property
WHERE  tsp.service_property = $1::varchar
AND    tsp.organizational_team_id = $2::uuid;`

const PropertyTemplateShow = `
SELECT sp.service_property,
       spv.service_property_attribute,
       spv.value
FROM   soma.service_properties sp
JOIN   soma.service_property_values spv
ON     sp.service_property = spv.service_property
WHERE  sp.service_property = $1::varchar;`

const PropertySystemAdd = `
INSERT INTO soma.system_properties (
            system_property)
SELECT $1::varchar
WHERE  NOT EXISTS (
   SELECT system_property
   FROM   soma.system_properties
   WHERE  system_property = $1::varchar);`

const PropertyNativeAdd = `
INSERT INTO soma.native_properties (
            native_property)
SELECT $1::varchar
WHERE  NOT EXISTS (
   SELECT native_property
   FROM   soma.native_properties
   WHERE  native_property = $1::varchar);`

const PropertyCustomAdd = `
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

const PropertyServiceAdd = `
INSERT INTO soma.team_service_properties (
            organizational_team_id,
            service_property)
SELECT $1::uuid, $2::varchar
WHERE  NOT EXISTS (
   SELECT service_property
   FROM   soma.team_service_properties
   WHERE  organizational_team_id = $1::uuid
   AND    service_property = $2::varchar);`

const PropertyServiceAttributeAdd = `
INSERT INTO soma.team_service_property_values (
            organizational_team_id,
            service_property,
            service_property_attribute,
            value)
SELECT $1::uuid, $2::varchar, $3::varchar, $4::varchar;`

const PropertyTemplateAdd = `
INSERT INTO soma.service_properties (
            service_property)
SELECT $1::varchar
WHERE  NOT EXISTS (
   SELECT service_property
   FROM   soma.service_properties
   WHERE  service_property = $1::varchar);`

const PropertyTemplateAttributeAdd = `
INSERT INTO soma.service_property_values (
            service_property,
            service_property_attribute,
            value)
SELECT $1::varchar, $2::varchar, $3::varchar;`

const PropertySystemDel = `
DELETE FROM soma.system_properties
WHERE  system_property = $1::varchar;`

const PropertyNativeDel = `
DELETE FROM soma.native_properties
WHERE  native_property = $1::varchar;`

const PropertyCustomDel = `
DELETE FROM soma.custom_properties
WHERE  repository_id = $1::uuid
AND    custom_property_id = $2::uuid;`

const PropertyServiceDel = `
DELETE FROM soma.team_service_properties
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`

const PropertyServiceAttributeDel = `
DELETE FROM soma.team_service_property_values
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`

const PropertyTemplateDel = `
DELETE FROM soma.service_properties
WHERE  service_property = $1::varchar;`

const PropertyTemplateAttributeDel = `
DELETE FROM soma.service_property_values
WHERE  service_property = $1::varchar;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
