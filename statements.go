package main

const stmtCheckItemExists = `
SELECT configuration_item_id
FROM   monsoon.configuration_items
WHERE  configuration_item_id = $1::uuid;`

const stmtCheckLookupExists = `
SELECT lookup_id
FROM   monsoon.configuration_lookup
WHERE  configuration_item_id = $1::char;`

const stmtInsertLookupInformation = `
INSERT INTO monsoon.configuration_lookup (
	lookup_id,
	host_id,
	metric)
SELECT $1::char,
       $2::numeric,
	   $3::text
WHERE NOT EXISTS (
    SELECT lookup_id
	FROM   monsoon.configuration_lookup
	WHERE  lookup_id = $1::char
	OR     ( host_id = $2::numeric AND metric = $3::text);`

const stmtInsertConfigurationItem = `
INSERT INTO monsoon.configuration_items (
	configuration_item_id,
	lookup_id,
	configuration
SELECT $1::uuid,
       $2::char,
	   $3::jsonb;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
