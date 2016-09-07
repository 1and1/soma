/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const ForestRebuildDeleteChecks = `
UPDATE soma.checks sc
SET    deleted = 'yes'::boolean
WHERE  sc.repository_id = $1::uuid;`

const ForestRebuildDeleteInstances = `
UPDATE soma.check_instances sci
SET    deleted = 'yes'::boolean
FROM   soma.checks sc
WHERE  sci.check_id = sc.check_id
AND    sc.repository_id = $1::uuid;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
