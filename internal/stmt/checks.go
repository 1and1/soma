/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const CheckDetailsForDelete = `
SELECT scc.configuration_object,
       scc.configuration_object_type,
       sc.source_check_id
FROM   soma.check_configurations scc
JOIN   soma.checks sc
  ON   scc.configuration_id = sc.configuration_id
WHERE  scc.configuration_id = $1::uuid
  AND  scc.repository_id    = $2::uuid
  AND  sc.check_id          = sc.source_check_id;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
