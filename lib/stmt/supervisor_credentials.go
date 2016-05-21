/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const LoadAllUserCredentials = `
SELECT aua.user_id,
       aua.crypt,
       aua.reset_pending,
       aua.valid_from,
       aua.valid_until,
	   iu.user_uid,
	   iu.user_is_active
FROM   auth.user_authentication aua
JOIN   inventory.users iu
ON     aua.user_id = iu.user_id
WHERE  iu.user_id != '00000000-0000-0000-0000-000000000000'::uuid;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
