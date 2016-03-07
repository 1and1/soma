package main

/*
 * Statements for CHECK actions
 */

const tkStmtCreateCheckConfigurationBase = `
INSERT INTO soma.check_configurations (
            configuration_id,
            configuration_name,
            interval,
            repository_id,
            bucket_id,
            capability_id,
            configuration_object,
            configuration_object_type,
            configuration_active,
            enabled,
            inheritance_enabled,
            children_only,
            external_id)
SELECT $1::uuid,
       $2::varchar,
       $3::integer,
       $4::uuid,
       $5::uuid,
       $6::uuid,
       $7::uuid,
       $8::varchar,
       $9::boolean,
       $10::boolean,
       $11::boolean,
       $12::boolean,
       $13::varchar;`

const tkStmtCreateCheckConfigurationThreshold = `
INSERT INTO soma.configuration_thresholds (
            configuration_id,
            predicate,
            threshold,
            notification_level)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar,
       $4::varchar;`

const tkStmtCreateCheckConfigurationConstraintSystem = `
INSERT INTO soma.constraints_system_property (
            configuration_id,
            system_property,
            property_value)
SELECT $1::uuid,
       $2::varchar,
       $3::text;`

const tkStmtCreateCheckConfigurationConstraintNative = `
INSERT INTO soma.constraints_native_property (
            configuration_id,
            native_property,
            property_value)
SELECT $1::uuid,
       $2::varchar,
       $3::text;`

const tkStmtCreateCheckConfigurationConstraintOncall = `
INSERT INTO soma.constraints_oncall_property (
            configuration_id,
            oncall_duty_id)
SELECT $1::uuid,
       $2::uuid;`

const tkStmtCreateCheckConfigurationConstraintCustom = `
INSERT INTO soma.constraints_custom_property (
            configuration_id,
            custom_property_id,
            repository_id,
            property_value)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::text;`

const tkStmtCreateCheckConfigurationConstraintService = `
INSERT INTO soma.constraints_service_property (
            configuration_id,
            organizational_team_id,
            service_property)
SELECT $1::uuid,
       $2::uuid,
       $3::varchar;`

const tkStmtCreateCheckConfigurationConstraintAttribute = `
INSERT INTO soma.constraints_service_attribute (
            configuration_id,
            service_property_attribute,
            attribute_value)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar;`

const tkStmtCreateCheck = `
INSERT INTO soma.checks (
            check_id,
            repository_id,
            bucket_id,
            source_check_id,
            source_object_type,
            source_object_id,
            configuration_id,
            capability_id,
            object_id,
            object_type)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::uuid,
       $5::varchar,
       $6::uuid,
       $7::uuid,
       $8::uuid,
       $9::uuid,
       $10::varchar;`

const tkStmtCreateCheckInstance = `
INSERT INTO soma.check_instances (
	check_instance_id,
	check_id,
	check_configuration_id,
	last_configuration_created)
VALUES $1::uuid,
       $2::uuid,
	   $3::uuid,
	   $4::timestamptz;`

const tkStmtCreateCheckInstanceConfiguration = `
INSERT INTO soma.check_instance_configurations (
	check_instance_config_id,
	version,
	check_instance_id,
	constraint_hash,
	constraint_val_hash,
	instance_service,
	instance_service_cfg_hash,
	instance_service_cfg,
	created,
	status,
	next_status,
	awaiting_deletion,
	deployment_details)
VALUES $1::uuid,
       $2::integer,
	   $3::uuid,
	   $4::varchar,
	   $5::varchar,
	   $6::varchar,
	   $7::varchar,
	   $8::jsonb,
	   $9::timestamptz,
	   $10::varchar,
	   $11::varchar,
	   $12::boolean,
	   $13::jsonb;`

////////////////////////////////////////////////
const tkStmtDeployDetailsComputeList = `
SELECT check_instance_config_id
FROM   soma.check_instance_configurations
WHERE  status = 'awaiting_computation';`

const tkStmtDeployDetailsCheckInstance = `
SELECT scic.version,
	   scic.check_instance_id,
	   scic.constraint_hash,
	   scic.constraint_val_hash,
	   scic.instance_service,
	   scic.instance_service_cfg_hash,
	   scic.instance_service_cfg,
	   sci.check_id,
	   sci.check_configuration_id
FROM   soma.check_instance_configurations scic
JOIN   soma.check_instances sci
ON     scic.check_instance_id = sci.check_instance_id
WHERE  scic.check_instance_config_id = $1::uuid;`

const tkStmtDeployDetailsCheck = `
SELECT sc.repository_id,
	   sc.source_check_id,
	   sc.source_object_type,
	   sc.source_object_id,
	   sc.capability_id,
	   sc.object_id,
	   sc.object_type,
	   scc.inheritance_enabled,
	   scc.children_only,
FROM   soma.checks sc
JOIN   soma.check_configurations scc
ON     sc.configuration_id = scc.configuration_id
WHERE  sc.check_id = $1::uuid;`

const tkStmtDeployDetailsCheckConfig = `
SELECT configuration_name,
       interval,
	   configuration_active,
	   enabled,
	   external_id
FROM   soma.check_configurations
WHERE  configuration_id = $1::uuid;`

const tkStmtDeployDetailsCheckConfigThreshold = `
SELECT sct.predicate,
	   sct.threshold,
	   sct.notification_level,
	   snl.level_shortname,
	   snl.level_numeric
FROM   soma.configuration_thresholds sct
JOIN   soma.notification_levels snl
WHERE  sct.configuration_id = $1::uuid;`

const tkStmtDeployDetailsCapabilityMonitoringMetric = `
SELECT smc.capability_metric,
       smc.capability_monitoring,
	   smc.capability_view,
	   smc.threshold_amount,
	   sms.monitoring_name,
	   sms.monitoring_system_mode,
	   sms.monitoring_contact,
	   sms.monitoring_owner_team,
	   sms.monitoring_callback_uri,
	   sm.metric_unit,
	   sm.description,
	   smu.metric_unit_long_name
FROM   soma.monitoring_capabilities smc
JOIN   soma.monitoring_systems sms
ON     smc.capability_monitoring = sms.monitoring_id
JOIN   soma.metrics sm
ON     smc.capability_metric = sm.metric
JOIN   soma.metric_units smu
ON     sm.metric_unit = smu.metric_unit
WHERE  smc.capability_id = $1::uuid;`

const tkStmtDeployDetailsProviders = `
SELECT metric_provider,
       package
FROM   soma.metric_packages
WHERE  metric = $1::varchar;`

const tkStmtDeployDetailsGroup = `
SELECT sg.bucket_id,
	   sg.group_name,
	   sg.object_state,
	   sg.organizational_team_id,
	   sb.bucket_name,
	   sb.environment,
	   sr.repository_name,
FROM  soma.groups sg
JOIN  soma.buckets sb
ON    sg.bucket_id = sb.bucket_id
JOIN  soma.repositories sr
ON    sb.repository_id = sr.repository_id
WHERE sg.group_id = $1::uuid;`

const tkStmtDeployDetailsCluster = `
SELECT sc.cluster_name,
       sc.bucket_id,
	   sc.object_state,
	   sc.organizational_team_id,
	   sb.bucket_name,
	   sb.environment,
	   sr.repository_name,
FROM   soma.clusters sc
JOIN   soma.buckets sb
ON     sc.bucket_id = sb.bucket_id
JOIN  soma.repositories sr
ON    sb.repository_id = sr.repository_id
WHERE sr.cluster_id = $1::uuid;`

const tkStmtDeployDetailsNode = `
SELECT sn.node_asset_id,
       sn.node_name,
	   sn.organizational_team_id,
	   sn.server_id,
	   sn.object_state,
	   sn.node_online,
	   sn.node_deleted,
	   sb.bucket_name,
	   sb.environment,
	   sr.repository_name,
	   ins.server_asset_id,
	   ins.server_datacenter_name,
	   ins.server_datacenter_location,
	   ins.server_name,
	   ins.server_online,
	   ins.server_deleted
FROM  soma.nodes sn
JOIN  soma.node_bucket_assignment snba
ON    sn.node_id = snba.node_id
JOIN  soma.buckets sb
ON    snba.bucket_id = sb.bucket_id
JOIN  soma.repositories sr
ON    sb.repository_id = sr.repository_id
JOIN  inventory.servers ins
ON    sn.server_id = ins.server_id
WHERE sn.node_id = $1::uuid;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
