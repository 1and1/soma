/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const ListUsers = `
SELECT user_id,
       user_uid
FROM   inventory.users;`

const ShowUsers = `
SELECT user_id,
       user_uid,
       user_first_name,
       user_last_name,
       user_employee_number,
       user_mail_address,
       user_is_active,
       user_is_system,
       user_is_deleted,
       organizational_team_id
FROM   inventory.users
WHERE  user_id = $1::uuid;`

const SyncUsers = `
SELECT user_id,
       user_uid,
       user_first_name,
       user_last_name,
       user_employee_number,
       user_mail_address,
       user_is_deleted,
       organizational_team_id
FROM   inventory.users
WHERE  NOT user_is_system;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
