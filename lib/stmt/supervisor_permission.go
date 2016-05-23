/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * All rights reserved
 */

package stmt

const LoadPermissions = `
SELECT permission_id,
       permission_name
FROM   soma.permissions;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
