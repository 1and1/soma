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
	UserStatements = ``

	ListUsers = `
SELECT user_id,
       user_uid
FROM   inventory.users;`

	ShowUsers = `
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

	SyncUsers = `
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

	UserAdd = `
INSERT INTO inventory.users (
            user_id,
            user_uid,
            user_first_name,
            user_last_name,
            user_employee_number,
            user_mail_address,
            user_is_active,
            user_is_system,
            user_is_deleted,
            organizational_team_id)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar,
       $4::varchar,
       $5::numeric,
       $6::text,
       $7::boolean,
       $8::boolean,
       $9::boolean,
       $10::uuid
WHERE  NOT EXISTS (
  SELECT user_id
  FROM   inventory.users
  WHERE  user_id = $1::uuid
     OR  user_uid = $2::varchar
     OR  user_employee_number = $5::numeric);`

	UserUpdate = `
UPDATE inventory.users
SET    user_uid = $1::varchar,
       user_first_name = $2::varchar,
       user_last_name = $3::varchar,
       user_employee_number = $4::numeric,
       user_mail_address = $5::text,
       user_is_deleted = $6::boolean,
       organizational_team_id = $7::uuid
WHERE  user_id = $8::uuid;`

	UserDel = `
UPDATE inventory.users
SET    user_is_deleted = 'yes',
       user_is_active = 'no'
WHERE  user_id = $1::uuid;`

	UserPurge = `
DELETE FROM inventory.users
WHERE  user_id = $1::uuid
AND    user_is_deleted;`
)

func init() {
	m[ListUsers] = `ListUsers`
	m[ShowUsers] = `ShowUsers`
	m[SyncUsers] = `SyncUsers`
	m[UserAdd] = `UserAdd`
	m[UserDel] = `UserDel`
	m[UserPurge] = `UserPurge`
	m[UserUpdate] = `UserUpdate`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
