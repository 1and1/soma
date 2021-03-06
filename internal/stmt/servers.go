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
	ServerStatements = ``

	SyncServers = `
SELECT server_id,
       server_asset_id,
       server_datacenter_name,
       server_datacenter_location,
       server_name,
       server_online,
       server_deleted
FROM   inventory.servers
WHERE  server_id != '00000000-0000-0000-0000-000000000000'::uuid;`

	ListServers = `
SELECT server_id,
       server_name,
       server_asset_id
FROM   inventory.servers
WHERE  server_online
AND    NOT server_deleted
AND    NOT server_id = '00000000-0000-0000-0000-000000000000'::uuid;`

	ShowServers = `
SELECT server_id,
       server_asset_id,
       server_datacenter_name,
       server_datacenter_location,
       server_name,
       server_online,
       server_deleted
FROM   inventory.servers
WHERE  server_id = $1::uuid;`

	SearchServerByName = `
SELECT server_id,
       server_name,
       server_asset_id
FROM   inventory.servers
WHERE  server_online
AND    NOT server_deleted
AND    NOT server_id = '00000000-0000-0000-0000-000000000000'
AND    server_name = $1::varchar;`

	SearchServerByAssetId = `
SELECT server_id,
       server_name,
       server_asset_id
FROM   inventory.servers
WHERE  server_online
AND    NOT server_deleted
AND    NOT server_id = '00000000-0000-0000-0000-000000000000'
AND    server_asset_id = $1::numeric;`

	AddServers = `
INSERT INTO inventory.servers (
            server_id,
            server_asset_id,
            server_datacenter_name,
            server_datacenter_location,
            server_name,
            server_online,
            server_deleted)
SELECT      $1::uuid, $2::numeric, $3, $4, $5, $6, $7
WHERE NOT   EXISTS(
    SELECT  server_id
    FROM    inventory.servers
    WHERE   server_id = $1::uuid
       OR   server_asset_id = $2::numeric);`

	UpdateServers = `
UPDATE inventory.servers
SET    server_asset_id = $2::numeric,
       server_datacenter_name = $3::varchar,
       server_datacenter_location = $4::varchar,
       server_name = $5::varchar,
       server_online = $6::boolean,
       server_deleted = $7::boolean
WHERE  server_id = $1::uuid;`

	DeleteServers = `
UPDATE inventory.servers
SET    server_deleted = 'yes'::boolean,
       server_online = 'no'::boolean
WHERE  server_id = $1::uuid
AND    server_id != '00000000-0000-0000-0000-000000000000'::uuid;`

	PurgeServers = `
DELETE FROM inventory.servers
WHERE  server_id = $1::uuid
  AND  server_deleted
  AND  server_id != '00000000-0000-0000-0000-000000000000'::uuid;`
)

func init() {
	m[AddServers] = `AddServers`
	m[DeleteServers] = `DeleteServers`
	m[ListServers] = `ListServers`
	m[PurgeServers] = `PurgeServers`
	m[SearchServerByAssetId] = `SearchServerByAssetId`
	m[SearchServerByName] = `SearchServerByName`
	m[ShowServers] = `ShowServers`
	m[SyncServers] = `SyncServers`
	m[UpdateServers] = `UpdateServers`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
