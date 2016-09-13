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
       iu.user_uid
FROM   inventory.users iu
JOIN   auth.user_authentication aua
ON     iu.user_id = aua.user_id
WHERE  iu.user_id != '00000000-0000-0000-0000-000000000000'::uuid
AND    NOW() < aua.valid_until
AND    NOT iu.user_is_deleted
AND    iu.user_is_active;`

const FindUserID = `
SELECT user_id
FROM   inventory.users
WHERE  user_uid = $1::varchar
AND    NOT user_is_deleted;`

const CheckUserActive = `
SELECT user_is_active
FROM   inventory.users
WHERE  user_id = $1::uuid
AND    NOT user_is_deleted;`

const SetUserCredential = `
INSERT INTO auth.user_authentication (
            user_id,
            crypt,
            reset_pending,
            valid_from,
            valid_until
) VALUES (
            $1::uuid,
            $2::text,
            'no'::boolean,
            $3::timestamptz,
			$4::timestamptz
);`

const ActivateUser = `
UPDATE inventory.users
SET    user_is_active = 'yes'::boolean
WHERE  user_id = $1::uuid;`

const InvalidateUserCredential = `
UPDATE auth.user_authentication
SET    valid_until = $1::timestamptz
WHERE  user_id = $2::uuid
  AND  user_is_active = 'yes'::boolean
  AND  user_id != '00000000-0000-0000-0000-000000000000';`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
