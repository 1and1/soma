package main

const stmtCheckItemExists = `
SELECT configuration_item_id
FROM   eye.configuration_items
WHERE  configuration_item_id = $1::uuid;`

const stmtCheckLookupExists = `
SELECT lookup_id
FROM   eye.configuration_lookup
WHERE  lookup_id = $1::char;`

const stmtInsertLookupInformation = `
INSERT INTO eye.configuration_lookup (
            lookup_id,
            host_id,
            metric)
SELECT $1::char,
       $2::numeric,
       $3::text
WHERE NOT EXISTS (
    SELECT lookup_id
    FROM   eye.configuration_lookup
    WHERE  lookup_id = $1::char
    OR     ( host_id = $2::numeric AND metric = $3::text));`

const stmtInsertConfigurationItem = `
INSERT INTO eye.configuration_items (
            configuration_item_id,
            lookup_id,
            configuration)
SELECT $1::uuid,
       $2::char,
       $3::jsonb
WHERE NOT EXISTS (
      SELECT configuration_item_id
      FROM   eye.configuration_items
      WHERE  configuration_item_id = $1::uuid);`

const stmtUpdateConfigurationItem = `
UPDATE eye.configuration_items
SET    lookup_id = $2::char,
       configuration = $3::jsonb
WHERE  configuration_item_id = $1::uuid;`

const stmtGetLookupIdForItem = `
SELECT lookup_id
FROM   eye.configuration_items
WHERE  configuration_item_id = $1::uuid;`

const stmtGetItemCountForLookupId = `
SELECT COUNT(1)::integer
FROM   eye.configuration_items
WHERE  lookup_id = $1::char;`

const stmtDeleteConfigurationItem = `
DELETE FROM eye.configuration_items
WHERE       configuration_item_id = $1::uuid;`

const stmtDeleteLookupId = `
DELETE FROM eye.configuration_lookup
WHERE       lookup_id = $1::char;`

const stmtGetConfigurationItemIds = `
SELECT configuration_item_id
FROM   eye.configuration_items;`

const stmtGetSingleConfiguration = `
SELECT configuration
FROM   eye.configuration_items
WHERE  configuration_item_id = $1::uuid;`

const stmtRetrieveConfigurationsByLookup = `
SELECT configuration
FROM   eye.configuration_items
WHERE  lookup_id = $1::char;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
