package main

var stmtAttributeList = `
SELECT service_property_attribute
FROM   soma.service_property_attributes;`

var stmtAttributeShow = `
SELECT service_property_attribute
FROM   soma.service_property_attributes
WHERE  service_property_attribute = $1::varchar;`

var stmtAttributeAdd = `
INSERT INTO soma.service_property_attributes (
            service_property_attribute)
SELECT $1::varchar WHERE NOT EXISTS (
       SELECT service_property_attribute
       FROM   soma.service_property_attributes
       WHERE  service_property_attribute = $1::varchar);`

var stmtAttributeDelete = `
DELETE FROM soma.service_property_attributes
WHERE       service_property_attribute = $1::varchar;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
