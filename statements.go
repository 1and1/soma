/*
 * Copyright (c) 2016, 1&1 Internet SE
 * Written by Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved.
 */

package main

const stmtCheckItemExists = `
SELECT configuration_item_id
FROM   eye.configuration_items
WHERE  configuration_item_id = $1::uuid;`

const stmtCheckLookupExists = `
SELECT lookup_id
FROM   eye.configuration_lookup
WHERE  lookup_id = $1::varchar;`

const stmtInsertLookupInformation = `
INSERT INTO eye.configuration_lookup (
            lookup_id,
            host_id,
            metric)
SELECT $1::varchar,
       $2::numeric,
       $3::text
WHERE NOT EXISTS (
    SELECT lookup_id
    FROM   eye.configuration_lookup
    WHERE  lookup_id = $1::varchar
    OR     ( host_id = $2::numeric AND metric = $3::text));`

const stmtInsertConfigurationItem = `
INSERT INTO eye.configuration_items (
            configuration_item_id,
            lookup_id,
            configuration)
SELECT $1::uuid,
       $2::varchar,
       $3::jsonb
WHERE NOT EXISTS (
      SELECT configuration_item_id
      FROM   eye.configuration_items
      WHERE  configuration_item_id = $1::uuid);`

const stmtUpdateConfigurationItem = `
UPDATE eye.configuration_items
SET    lookup_id = $2::varchar,
       configuration = $3::jsonb
WHERE  configuration_item_id = $1::uuid;`

const stmtGetLookupIdForItem = `
SELECT lookup_id
FROM   eye.configuration_items
WHERE  configuration_item_id = $1::uuid;`

const stmtGetItemCountForLookupId = `
SELECT COUNT(1)::integer
FROM   eye.configuration_items
WHERE  lookup_id = $1::varchar;`

const stmtDeleteConfigurationItem = `
DELETE FROM eye.configuration_items
WHERE       configuration_item_id = $1::uuid;`

const stmtDeleteLookupId = `
DELETE FROM eye.configuration_lookup
WHERE       lookup_id = $1::varchar;`

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
WHERE  lookup_id = $1::varchar;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
