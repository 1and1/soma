/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const DatacenterList = `
SELECT datacenter
FROM   inventory.datacenters;`

const DatacenterShow = `
SELECT datacenter
FROM   inventory.datacenters
WHERE  datacenter = $1::varchar;`

const DatacenterGroupList = `
SELECT DISTINCT datacenter_group
FROM   soma.datacenter_groups;`

const DatacenterGroupShow = `
SELECT DISTINCT datacenter
FROM   soma.datacenter_groups
WHERE  datacenter_group = $1::varchar;`

const DatacenterAdd = `
INSERT INTO inventory.datacenters (
            datacenter)
SELECT $1::varchar
WHERE  NOT EXISTS (
   SELECT datacenter
   FROM inventory.datacenters
   WHERE datacenter = $1::varchar);`

const DatacenterDel = `
DELETE FROM inventory.datacenters
WHERE datacenter = $1::varchar;`

const DatacenterRename = `
UPDATE inventory.datacenters
SET    datacenter = $1::varchar
WHERE  datacenter = $2::varchar;`

const DatacenterGroupAdd = `
INSERT INTO soma.datacenter_groups (
            datacenter_group,
            datacenter)
SELECT $1::varchar, $2::varchar
WHERE  NOT EXISTS (
   SELECT datacenter
   FROM   soma.datacenter_groups
   WHERE  datacenter_group = $3::varchar
     AND  datacenter = $4::varchar);`

const DatacenterGroupDel = `
DELETE FROM soma.datacenter_groups
WHERE       datacenter_group = $1::varchar
  AND       datacenter = $2::varchar;`

func init() {
	m[DatacenterAdd] = `DatacenterAdd`
	m[DatacenterDel] = `DatacenterDel`
	m[DatacenterGroupAdd] = `DatacenterGroupAdd`
	m[DatacenterGroupDel] = `DatacenterGroupDel`
	m[DatacenterGroupList] = `DatacenterGroupList`
	m[DatacenterGroupShow] = `DatacenterGroupShow`
	m[DatacenterList] = `DatacenterList`
	m[DatacenterRename] = `DatacenterRename`
	m[DatacenterShow] = `DatacenterShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
