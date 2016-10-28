/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const EnvironmentList = `
SELECT environment
FROM   soma.environments;`

const EnvironmentShow = `
SELECT environment
FROM   soma.environments
WHERE  environment = $1::varchar;`

const EnvironmentAdd = `
INSERT INTO soma.environments (
            environment
)
SELECT $1::varchar
WHERE NOT EXISTS (
    SELECT environment
    FROM   soma.environments
    WHERE  environment = $2::varchar
);`

const EnvironmentDel = `
  DELETE FROM soma.environments
  WHERE environment = $1;
  `

const EnvironmentRename = `
  UPDATE soma.environments SET environment = $1
  WHERE environment = $2;
  `

func init() {
	m[EnvironmentAdd] = `EnvironmentAdd`
	m[EnvironmentDel] = `EnvironmentDel`
	m[EnvironmentList] = `EnvironmentList`
	m[EnvironmentRename] = `EnvironmentRename`
	m[EnvironmentShow] = `EnvironmentShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
