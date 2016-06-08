/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const SyncServers = `
SELECT server_id,
       server_asset_id,
       server_datacenter_name,
       server_datacenter_location,
       server_name,
       server_online,
       server_deleted
FROM   inventory.servers
WHERE  server_id != '00000000-0000-0000-0000-000000000000';`

const ListServers = `
SELECT server_id,
       server_name,
       server_asset_id
FROM   inventory.servers
WHERE  server_online
AND    NOT server_deleted
AND    NOT server_id = '00000000-0000-0000-0000-000000000000';`

const ShowServers = `
SELECT server_id,
       server_asset_id,
       server_datacenter_name,
       server_datacenter_location,
       server_name,
       server_online,
       server_deleted
FROM   inventory.servers
WHERE  server_id = $1;`

const SearchServerByName = `
SELECT server_id,
       server_name,
       server_asset_id
FROM   inventory.servers
WHERE  server_online
AND    NOT server_deleted
AND    NOT server_id = '00000000-0000-0000-0000-000000000000'
AND    server_name = $1::varchar;`

const SearchServerByAssetId = `
SELECT server_id,
       server_name,
       server_asset_id
FROM   inventory.servers
WHERE  server_online
AND    NOT server_deleted
AND    NOT server_id = '00000000-0000-0000-0000-000000000000'
AND    server_asset_id = $1::numeric;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
