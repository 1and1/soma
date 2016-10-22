/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const ModeList = `
SELECT monitoring_system_mode
FROM   soma.monitoring_system_modes; `

const ModeShow = `
SELECT monitoring_system_mode
FROM   soma.monitoring_system_modes
WHERE  monitoring_system_mode = $1::varchar;`

const ModeAdd = `
INSERT INTO soma.monitoring_system_modes (
            monitoring_system_mode)
SELECT  $1::varchar
WHERE   NOT EXISTS (
   SELECT monitoring_system_mode
   FROM   soma.monitoring_system_modes
   WHERE  monitoring_system_mode = $1::varchar);`

const ModeDel = `
DELETE FROM soma.monitoring_system_modes
WHERE  monitoring_system_mode = $1::varchar;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
