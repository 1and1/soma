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
       $3::jsonb
WHERE NOT EXISTS (
      SELECT configuration_item_id
      FROM   monsoon.configuration_items
      WHERE  configuration_item_id = $1::uuid;`

const stmtUpdateConfigurationItem = `
UPDATE monsoon.configuration_items
SET    lookup_id = $2::char,
       configuration = $3::jsonb
WHERE  configuration_item_id = $1::uuid;`

const stmtGetLookupIdForItem = `
SELECT lookup_id
FROM   monsoon.configuration_items
WHERE  configuration_item_id = $1::uuid;`

const stmtGetItemCountForLookupId = `
SELECT COUNT(1)::integer
FROM   monsoon.configuration_items
WHERE  lookup_id = $1::char;`

const stmtDeleteConfigurationItem = `
DELETE FROM monsoon.configuration_items
WHERE       configuration_item_id = $1::uuid;`

const stmtDeleteLookupId = `
DELETE FROM monsoon.configuration_lookup
WHERE       lookup_id = $1::char;`

const stmtGetConfigurationItemIds = `
SELECT configuration_item_id
FROM   monsoon.configuration_items;`

const stmtGetSingleConfiguration = `
SELECT configuration
FROM   monsoon.configuration_items
WHERE  configuration_item_id = $1::uuid;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
