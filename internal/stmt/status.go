/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const StatusList = `
SELECT status
FROM   soma.check_instance_status;`

const StatusShow = `
SELECT status
FROM   soma.check_instance_status
WHERE  status = $1;`

const StatusAdd = `
INSERT INTO soma.check_instance_status (
            status)
SELECT $1::varchar
WHERE  NOT EXISTS (
   SELECT status
   FROM   soma.check_instance_status
   WHERE  status = $1::varchar);`

const StatusDel = `
DELETE FROM soma.check_instance_status
WHERE  status = $1;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
