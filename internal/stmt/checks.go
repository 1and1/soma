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
  AND  sc.check_id          = sc.source_check_id
  AND  NOT sc.deleted;`

const CheckConfigList = `
SELECT configuration_id,
       repository_id,
       bucket_id,
       configuration_name
FROM   soma.check_configurations
WHERE  repository_id = $1::uuid
AND    NOT deleted;`

const CheckConfigShowBase = `
SELECT configuration_id,
       repository_id,
       bucket_id,
       configuration_name,
       configuration_object,
       configuration_object_type,
       configuration_active,
       inheritance_enabled,
       children_only,
       capability_id,
       interval,
       enabled,
       external_id
FROM   soma.check_configurations
WHERE  configuration_id = $1::uuid;`

const CheckConfigShowThreshold = `
SELECT scc.configuration_id,
       sct.predicate,
       sct.threshold,
       sct.notification_level,
       snl.level_shortname,
       snl.level_numeric
FROM   soma.check_configurations scc
JOIN   soma.configuration_thresholds sct
ON     scc.configuration_id = sct.configuration_id
JOIN   soma.notification_levels snl
ON     sct.notification_level = snl.level_name
WHERE  scc.configuration_id = $1::uuid;`

const CheckConfigShowConstrCustom = `
SELECT scc.configuration_id,
       sccp.custom_property_id,
       sccp.repository_id,
       sccp.property_value,
       scp.custom_property
FROM   soma.check_configurations scc
JOIN   soma.constraints_custom_property sccp
ON     scc.configuration_id = sccp.configuration_id
JOIN   soma.custom_properties scp
ON     sccp.custom_property_id = scp.custom_property_id
AND    sccp.repository_id = scp.repository_id
WHERE  scc.configuration_id = $1::uuid;`

const CheckConfigShowConstrSystem = `
SELECT scc.configuration_id,
       scsp.system_property,
       scsp.property_value
FROM   soma.check_configurations scc
JOIN   soma.constraints_system_property scsp
ON     scc.configuration_id = scsp.configuration_id
WHERE  scc.configuration_id = $1::uuid;`

const CheckConfigShowConstrNative = `
SELECT scc.configuration_id,
       scnp.native_property,
       scnp.property_value
FROM   soma.check_configurations scc
JOIN   soma.constraints_native_property scnp
ON     scc.configuration_id = scnp.configuration_id
WHERE  scc.configuration_id = $1::uuid;`

const CheckConfigShowConstrService = `
SELECT scc.configuration_id,
       scsvp.organizational_team_id,
       scsvp.service_property
FROM   soma.check_configurations scc
JOIN   soma.constraints_service_property scsvp
ON     scc.configuration_id = scsvp.configuration_id
WHERE  scc.configuration_id = $1::uuid;`

const CheckConfigShowConstrAttribute = `
SELECT scc.configuration_id,
       scsa.service_property_attribute,
       scsa.attribute_value
FROM   soma.check_configurations scc
JOIN   soma.constraints_service_attribute scsa
ON     scc.configuration_id = scsa.configuration_id
WHERE  scc.configuration_id = $1::uuid;`

const CheckConfigShowConstrOncall = `
SELECT scc.configuration_id,
       scop.oncall_duty_id,
       iodt.oncall_duty_name,
       iodt.oncall_duty_phone_number
FROM   soma.check_configurations scc
JOIN   soma.constraints_oncall_property scop
ON     scc.configuration_id = scop.configuration_id
JOIN   inventory.oncall_duty_teams iodt
ON     scop.oncall_duty_id = iodt.oncall_duty_id
WHERE  scc.configuration_id = $1::uuid;`

const CheckConfigInstanceInfo = `
SELECT sci.check_instance_id,
       sc.object_id,
       sc.object_type,
       scic.status,
       scic.next_status
FROM   soma.check_configurations scc
JOIN   soma.check_instances sci
  ON   scc.configuration_id = sci.check_configuration_id
JOIN   soma.checks sc
  ON   sci.check_id = sc.check_id
JOIN   soma.check_instance_configurations scic
  ON   sci.current_instance_config_id = scic.check_instance_config_id
WHERE  scc.configuration_id = $1::uuid
  AND  scic.status != 'awaiting_deletion';`

const CheckConfigObjectInstanceInfo = `
SELECT sci.check_instance_id,
       sc.object_id,
       sc.object_type,
       scic.status,
       scic.next_status
FROM   soma.check_configurations scc
JOIN   soma.check_instances sci
  ON   scc.configuration_id = sci.check_configuration_id
JOIN   soma.checks sc
  ON   sci.check_id = sc.check_id
JOIN   soma.check_instance_configurations scic
  ON   sci.current_instance_config_id = scic.check_instance_config_id
WHERE  scc.configuration_id = $1::uuid
  AND  sc.object_id = $2::uuid
  AND  scic.status != 'awaiting_deletion';`

const CheckConfigForChecksOnObject = `
SELECT sc.configuration_id
FROM   soma.checks sc
WHERE  sc.object_id = $1::uuid;`

func init() {
	m[CheckConfigForChecksOnObject] = `CheckConfigForChecksOnObject`
	m[CheckConfigInstanceInfo] = `CheckConfigInstanceInfo`
	m[CheckConfigList] = `CheckConfigList`
	m[CheckConfigObjectInstanceInfo] = `CheckConfigObjectInstanceInfo`
	m[CheckConfigShowBase] = `CheckConfigShowBase`
	m[CheckConfigShowConstrAttribute] = `CheckConfigShowConstrAttribute`
	m[CheckConfigShowConstrCustom] = `CheckConfigShowConstrCustom`
	m[CheckConfigShowConstrNative] = `CheckConfigShowConstrNative`
	m[CheckConfigShowConstrOncall] = `CheckConfigShowConstrOncall`
	m[CheckConfigShowConstrService] = `CheckConfigShowConstrService`
	m[CheckConfigShowConstrSystem] = `CheckConfigShowConstrSystem`
	m[CheckConfigShowThreshold] = `CheckConfigShowThreshold`
	m[CheckDetailsForDelete] = `CheckDetailsForDelete`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
