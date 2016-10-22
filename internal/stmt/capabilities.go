/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const ListAllCapabilities = `
SELECT smc.capability_id,
       smc.capability_monitoring,
       smc.capability_metric,
       smc.capability_view,
       sms.monitoring_name
FROM   soma.monitoring_capabilities smc
JOIN   soma.monitoring_systems sms
  ON   smc.capability_monitoring = sms.monitoring_id;`

const ListScopedCapabilities = `
WITH sysid AS (
    SELECT sms.monitoring_id
    FROM   inventory.users iu
    JOIN   soma.monitoring_system_users smsu
      ON   iu.organizational_team_id = smsu.organizational_team_id
    JOIN   soma.monitoring_systems sms
      ON   smsu.monitoring_id = sms.monitoring_id
    WHERE  iu.user_id = $1::uuid
      AND  sms.monitoring_system_mode = 'private'
    UNION
    SELECT sms.monitoring_id
    FROM   soma.monitoring_systems sms
    WHERE  sms.monitoring_system_mode = 'public'
)
SELECT smc.capability_id,
       smc.capability_monitoring,
       smc.capability_metric,
       smc.capability_view,
       sms.monitoring_name
FROM   soma.monitoring_capabilities smc
JOIN   soma.monitoring_systems sms
  ON   smc.capability_monitoring = sms.monitoring_id
WHERE  smc.capability_monitoring IN (SELECT monitoring_id FROM sysid);`

const ShowCapability = `
SELECT smc.capability_id,
       smc.capability_monitoring,
       smc.capability_metric,
       smc.capability_view,
       smc.threshold_amount,
       sms.monitoring_name
FROM   soma.monitoring_capabilities smc
JOIN   soma.monitoring_systems sms
ON     smc.capability_monitoring = sms.monitoring_id
WHERE  smc.capability_id = $1::uuid;`

const CapabilityThresholds = `
SELECT threshold_amount
FROM   soma.monitoring_capabilities
WHERE  capability_id = $1::uuid;`

const AddCapability = `
INSERT INTO soma.monitoring_capabilities (
            capability_id,
            capability_monitoring,
            capability_metric,
            capability_view,
            threshold_amount)
SELECT   $1::uuid, $2::uuid, $3::varchar, $4::varchar, $5::integer
WHERE    NOT EXISTS (
    SELECT capability_id
    FROM   soma.monitoring_capabilities
    WHERE  capability_id = $1::uuid
    OR     (    capability_monitoring = $2::uuid
            AND capability_metric     = $3::varchar
            AND capability_view       = $4::varchar));`

const DelCapability = `
DELETE FROM soma.monitoring_capabilities
WHERE  capability_id = $1::uuid;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
