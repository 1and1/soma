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
	SupervisorRootStatements = ``

	// the bootstrap token to initialize the system
	SelectRootToken = `
SELECT token
FROM   root.token;`

	// 'restricted' => true|false
	// 'disabled' => true|false
	LoadRootFlags = `
SELECT flag,
       status
FROM   root.flags;`

	LoadRootPassword = `
SELECT aua.crypt,
       aua.valid_from,
       aua.valid_until
FROM   inventory.users ui
JOIN   auth.user_authentication aua
ON     ui.user_id = aua.user_id
WHERE  ui.user_id = '00000000-0000-0000-0000-000000000000'::uuid
AND    ui.user_uid = 'root'
AND    ui.user_is_system
AND    ui.user_is_active
AND    aua.valid_from < NOW()
AND    aua.valid_until > NOW();`

	SetRootCredentials = `
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
            'infinity'::timestamptz);`
)

func init() {
	m[LoadRootFlags] = `LoadRootFlags`
	m[LoadRootPassword] = `LoadRootPassword`
	m[SelectRootToken] = `SelectRootToken`
	m[SetRootCredentials] = `SetRootCredentials`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
