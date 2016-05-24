/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * All rights reserved
 */

package stmt

const LoadPermissions = `
SELECT permission_id,
       permission_name
FROM   soma.permissions;`

const AddPermissionCategory = `
INSERT INTO soma.permission_types (
            permission_type,
            created_by
)
SELECT $1::varchar,
       $2::uuid
WHERE NOT EXISTS (
      SELECT permission_type
      FROM   soma.permission_types
      WHERE  permission_type = $1::varchar
);`

const DeletePermissionCategory = `
DELETE FROM soma.permission_types
WHERE permission_type = $1::varchar;`

const ListPermissionCategory = `
SELECT spt.permission_type
FROM   soma.permission_types spt:`

const ShowPermissionCategory = `
SELECT spt.permission_type,
       iu.user_uid,
       spt.created_by
FROM   soma.permission_types spt
JOIN   inventory.users iu
ON     spt.created_by = iu.user_id
WHERE  spt.permission_type = $1::varchar;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
