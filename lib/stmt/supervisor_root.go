/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const SelectRootToken = `
SELECT token
FROM   root.token;`

// 'restricted' => true|false
const DiscoverRootStatus = `
SELECT flag,
       status
FROM   root.flags;`

const DetectRootPasswordSet = `
SELECT aua.crypt,
       aua.valid_from,
       aua.valid_until
FROM   inventory.users ui
JOIN   auth.user_authentication aua
ON     ui.user_id = aua.user_id
WHERE  ui.user_id = '00000000-0000-0000-0000-000000000000'::uuid
AND    ui.user_uid = 'root'
AND    ui.user_is_system
AND    aua.valid_from < NOW()
AND    aua.valid_until > NOW();`

const SetRootCredentials = `
INSERT INTO auth.user_authentication (
    user_id,
    crypt,
    reset_pending,
    valid_from,
    valid_until
) VALUES (
    $1::uuid,
    $2::text,
    $3::timestamptz,
    infinity
);`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
