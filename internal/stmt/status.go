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
	StatusStatements = ``

	StatusList = `
SELECT status
FROM   soma.check_instance_status;`

	StatusShow = `
SELECT status
FROM   soma.check_instance_status
WHERE  status = $1;`

	StatusAdd = `
INSERT INTO soma.check_instance_status (
            status)
SELECT $1::varchar
WHERE  NOT EXISTS (
   SELECT status
   FROM   soma.check_instance_status
   WHERE  status = $1::varchar);`

	StatusDel = `
DELETE FROM soma.check_instance_status
WHERE  status = $1;`
)

func init() {
	m[StatusAdd] = `StatusAdd`
	m[StatusDel] = `StatusDel`
	m[StatusList] = `StatusList`
	m[StatusShow] = `StatusShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
