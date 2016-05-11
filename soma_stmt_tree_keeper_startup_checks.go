package main

const tkStmtLoadChecks = `
SELECT check_id,
       bucket_id,
	   source_check_id,
	   source_object_type,
	   source_object_id,
	   configuration_id,
	   capability_id,
	   object_id,
	   object_type
FROM   soma.checks
WHERE  repository_id = $1::uuid
AND    check_id = source_check_id
AND    source_object_type = $2::varchar;`

const tkStmtLoadInheritedChecks = `
SELECT check_id,
	   object_id,
	   object_type
FROM   soma.checks
WHERE  repository_id = $1::uuid
AND    source_check_id = $2::uuid
AND    source_check_id != check_id;`

const tkStmtLoadCheckConfiguration = `
SELECT bucket_id,
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
WHERE  configuration_id = $1::uuid
AND    repository_id = $2::uuid;`

const tkStmtLoadCheckThresholds = `
SELECT predicate,
       threshold,
	   notification_level
FROM   soma.configuration_thresholds
WHERE  configuration_id = $1::uuid;`

const tkStmtLoadCheckConstraintCustom = `
SELECT custom_property_id,
       repository_id,
	   property_value
FROM   soma.constraints_custom_property
WHERE  configuration_id = $1::uuid;`

const tkStmtLoadCheckConstraintNative = `
SELECT native_property,
       property_value
FROM   soma.constraints_native_property
WHERE  configuration_id = $1::uuid;`

const tkStmtLoadCheckConstraintOncall = `
SELECT oncall_duty_id
FROM   soma.constraints_oncall_property
WHERE  configuration_id = $1::uuid;`

const tkStmtLoadCheckConstraintAttribute = `
SELECT service_property_attribute,
       attribute_value
FROM   soma.constraints_service_attribute
WHERE  configuration_id = $1::uuid;`

const tkStmtLoadCheckConstraintService = `
SELECT organizational_team_id,
       service_property
FROM   soma.constraints_service_property
WHERE  configuration_id = $1::uuid;`

const tkStmtLoadCheckConstraintSystem = `
SELECT system_property,
       property_value
FROM   soma.constraints_system_property
WHERE  configuration_id = $1::uuid;`

const tkStmtLoadCheckInstances = `
SELECT check_instance_id,
       check_configuration_id,
	   current_instance_config_id
WHERE  check_id = $1::uuid
AND    NOT deleted;`

const tkStmtLoadCheckInstanceConfiguration = `
SELECT check_instance_config_id,
       version,
	   monitoring_id,
	   constraint_hash,
	   constraint_val_hash,
	   instance_service,
	   instance_service_cfg_hash,
	   instance_service_cfg
FROM   soma.check_instance_configurations
WHERE  check_instance_id = $1::uuid;`

const tkStmtLoadCheckGroupState = `
SELECT sg.group_id,
       sg.object_state
FROM   soma.buckets sb
JOIN   soma.groups  sg
ON     sb.bucket_id = sg.bucket_id
WHERE  sb.repository_id = $1::uuid;`

const tkStmtLoadCheckGroupRelations = `
SELECT sgmg.group_id,
       sgmg.child_group_id
FROM   soma.buckets sb
JOIN   soma.group_membership_groups sgmg
ON     sb.bucket_id = sgmg.bucket_id
WHERE  sb.repository_id = $1::uuid;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
