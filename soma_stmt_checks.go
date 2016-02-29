package main

var stmtCheckConfigList = `
SELECT configuration_id,
       repository_id,
       bucket_id,
	   configuration_name
FROM   soma.check_configurations
WHERE  repository_id = $1::uuid;`

var stmtCheckConfigShowBase = `
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

var stmtCheckConfigShowThreshold = `
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

var stmtCheckConfigShowConstrCustom = `
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

var stmtCheckConfigShowConstrSystem = `
SELECT scc.configuration_id,
       scsp.system_property,
	   scsp.property_value
FROM   soma.check_configurations scc
JOIN   soma.constraints_system_property scsp
ON     scc.configuration_id = scsp.configuration_id
WHERE  scc.configuration_id = $1::uuid;`

var stmtCheckConfigShowConstrNative = `
SELECT scc.configuration_id,
       scnp.native_property,
	   scnp.property_value
FROM   soma.check_configurations scc
JOIN   soma.constraints_native_property scnp
ON     scc.configuration_id = scnp.configuration_id
WHERE  scc.configuration_id = $1::uuid;`

var stmtCheckConfigShowConstrService = `
SELECT scc.configuration_id,
       scsvp.organizational_team_id,
	   scsvp.service_property
FROM   soma.check_configurations scc
JOIN   soma.constraints_service_property scsvp
ON     scc.configuration_id = scsvp.configuration_id
WHERE  scc.configuration_id = $1::uuid;`

var stmtCheckConfigShowConstrAttribute = `
SELECT scc.configuration_id,
       scsa.service_property_attribute,
	   scsa.attribute_value
FROM   soma.check_configurations scc
JOIN   soma.constraints_service_attribute scsa
ON     scc.configuration_id = scsa.configuration_id
WHERE  scc.configuration_id = $1::uuid;`

var stmtCheckConfigShowConstrOncall = `
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
