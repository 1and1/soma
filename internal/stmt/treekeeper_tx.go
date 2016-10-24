/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const TxMarkCheckConfigDeleted = `
UPDATE soma.check_configurations
SET    deleted = 'yes'::boolean
WHERE  configuration_id = $1::uuid;`

const TxCreateCheck = `
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

const TxMarkCheckDeleted = `
UPDATE soma.checks
SET    deleted = 'yes'::boolean
WHERE  check_id = $1::uuid;`

const TxCreateCheckInstance = `
INSERT INTO soma.check_instances (
            check_instance_id,
            check_id,
            check_configuration_id,
            current_instance_config_id,
            last_configuration_created)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::uuid,
       $5::timestamptz;`

const TxMarkCheckInstanceDeleted = `
UPDATE soma.check_instances
SET    deleted = 'yes'::boolean
WHERE  check_instance_id = $1::uuid;`

const TxCreateCheckInstanceConfiguration = `
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
SELECT $1::uuid,
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

const TxCreateCheckConfigurationBase = `
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

const TxCreateCheckConfigurationThreshold = `
INSERT INTO soma.configuration_thresholds (
            configuration_id,
            predicate,
            threshold,
            notification_level)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar,
       $4::varchar;`

const TxCreateCheckConfigurationConstraintSystem = `
INSERT INTO soma.constraints_system_property (
            configuration_id,
            system_property,
            property_value)
SELECT $1::uuid,
       $2::varchar,
       $3::text;`

const TxCreateCheckConfigurationConstraintNative = `
INSERT INTO soma.constraints_native_property (
            configuration_id,
            native_property,
            property_value)
SELECT $1::uuid,
       $2::varchar,
       $3::text;`

const TxCreateCheckConfigurationConstraintOncall = `
INSERT INTO soma.constraints_oncall_property (
            configuration_id,
            oncall_duty_id)
SELECT $1::uuid,
       $2::uuid;`

const TxCreateCheckConfigurationConstraintCustom = `
INSERT INTO soma.constraints_custom_property (
            configuration_id,
            custom_property_id,
            repository_id,
            property_value)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::text;`

const TxCreateCheckConfigurationConstraintService = `
INSERT INTO soma.constraints_service_property (
            configuration_id,
            organizational_team_id,
            service_property)
SELECT $1::uuid,
       $2::uuid,
       $3::varchar;`

const TxCreateCheckConfigurationConstraintAttribute = `
INSERT INTO soma.constraints_service_attribute (
            configuration_id,
            service_property_attribute,
            attribute_value)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar;`

const TxPropertyInstanceCreate = `
INSERT INTO soma.property_instances (
            instance_id,
            repository_id,
            source_instance_id,
            source_object_type,
            source_object_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid;`

const TxPropertyInstanceDelete = `
DELETE FROM soma.property_instances
WHERE       instance_id = $1::uuid;`

const TxFinishJob = `
UPDATE soma.jobs
SET    job_finished = $2::timestamptz,
       job_status = 'processed',
       job_result = $3::varchar,
       job_error = $4::text
WHERE  job_id = $1::uuid;`

const TxDeferAllConstraints = `
SET CONSTRAINTS ALL DEFERRED;`

const TxRepositoryPropertyOncallCreate = `
INSERT INTO soma.repository_oncall_properties (
            instance_id,
            source_instance_id,
            repository_id,
            view,
            oncall_duty_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::boolean,
       $7::boolean;`

const TxRepositoryPropertyOncallDelete = `
DELETE FROM soma.repository_oncall_properties
WHERE       instance_id = $1::uuid;`

const TxRepositoryPropertyServiceCreate = `
INSERT INTO soma.repository_service_properties (
            instance_id,
            source_instance_id,
            repository_id,
            view,
            service_property,
            organizational_team_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::uuid,
       $7::boolean,
       $8::boolean;`

const TxRepositoryPropertyServiceDelete = `
DELETE FROM soma.repository_service_properties
WHERE       instance_id = $1::uuid;`

const TxRepositoryPropertySystemCreate = `
INSERT INTO soma.repository_system_properties (
            instance_id,
            source_instance_id,
            repository_id,
            view,
            system_property,
            source_type,
            inheritance_enabled,
            children_only,
            value,
            inherited)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::varchar,
       $7::boolean,
       $8::boolean,
       $9::text,
       $10::boolean;`

const TxRepositoryPropertySystemDelete = `
DELETE FROM soma.repository_system_properties
WHERE       instance_id = $1::uuid;`

const TxRepositoryPropertyCustomCreate = `
INSERT INTO soma.repository_custom_properties (
            instance_id,
            source_instance_id,
            repository_id,
            view,
            custom_property_id,
            inheritance_enabled,
            children_only,
            value)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::boolean,
       $7::boolean,
       $8::text;`

const TxRepositoryPropertyCustomDelete = `
DELETE FROM soma.repository_custom_properties
WHERE       instance_id = $1::uuid;`

const TxUpdateNodeState = `
UPDATE soma.nodes
SET    object_state = $2::varchar
WHERE  node_id = $1::uuid;`

const TxNodeUnassignFromBucket = `
DELETE FROM soma.node_bucket_assignment
WHERE       node_id = $1::uuid
AND         bucket_id = $2::uuid
AND         organizational_team_id = $3::uuid;`

const TxNodePropertyOncallCreate = `
INSERT INTO soma.node_oncall_property (
            instance_id,
            source_instance_id,
            node_id,
            view,
            oncall_duty_id,
            repository_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::uuid,
       $7::boolean,
       $8::boolean;`

const TxNodePropertyOncallDelete = `
DELETE FROM soma.node_oncall_property
WHERE       instance_id = $1::uuid;`

const TxNodePropertyServiceCreate = `
INSERT INTO soma.node_service_properties (
            instance_id,
            source_instance_id,
            node_id,
            view,
            service_property,
            organizational_team_id,
            repository_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::uuid,
       $7::uuid,
       $8::boolean,
       $9::boolean;`

const TxNodePropertyServiceDelete = `
DELETE FROM soma.node_service_properties
WHERE       instance_id = $1::uuid;`

const TxNodePropertySystemCreate = `
INSERT INTO soma.node_system_properties (
            instance_id,
            source_instance_id,
            node_id,
            view,
            system_property,
            source_type,
            repository_id,
            inheritance_enabled,
            children_only,
            value,
            inherited)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::varchar,
       $7::uuid,
       $8::boolean,
       $9::boolean,
       $10::text,
       $11::boolean;`

const TxNodePropertySystemDelete = `
DELETE FROM soma.node_system_properties
WHERE       instance_id = $1::uuid;`

const TxNodePropertyCustomCreate = `
INSERT INTO soma.node_custom_properties (
            instance_id,
            source_instance_id,
            node_id,
            view,
            custom_property_id,
            bucket_id,
            repository_id,
            inheritance_enabled,
            children_only,
            value)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::uuid,
       $7::uuid,
       $8::boolean,
       $9::boolean,
       $10::text;`

const TxNodePropertyCustomDelete = `
DELETE FROM soma.node_custom_properties
WHERE       instance_id = $1::uuid;`

const TxGroupCreate = `
INSERT INTO soma.groups (
            group_id,
            bucket_id,
            group_name,
            object_state,
            organizational_team_id,
            created_by)
SELECT $1::uuid,
       $2::uuid,
       $3::varchar,
       $4::varchar,
       $5::uuid,
       user_id
FROM   inventory.users iu
WHERE  iu.user_uid = $6::varchar;`

const TxGroupUpdate = `
UPDATE soma.groups
SET    object_state = $2::varchar
WHERE  group_id = $1::uuid;`

const TxGroupDelete = `
DELETE FROM soma.groups
WHERE       group_id = $1::uuid;`

const TxGroupMemberNewNode = `
INSERT INTO soma.group_membership_nodes (
            group_id,
            child_node_id,
            bucket_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

const TxGroupMemberNewCluster = `
INSERT INTO soma.group_membership_clusters (
            group_id,
            child_cluster_id,
            bucket_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

const TxGroupMemberNewGroup = `
INSERT INTO soma.group_membership_groups (
            group_id,
            child_group_id,
            bucket_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

const TxGroupMemberRemoveNode = `
DELETE FROM soma.group_membership_nodes
WHERE       group_id = $1::uuid
AND         child_node_id = $2::uuid;`

const TxGroupMemberRemoveCluster = `
DELETE FROM soma.group_membership_clusters
WHERE       group_id = $1::uuid
AND         child_cluster_id = $2::uuid;`

const TxGroupMemberRemoveGroup = `
DELETE FROM soma.group_membership_groups
WHERE       group_id = $1::uuid
AND         child_group_id = $2::uuid;`

const TxGroupPropertyOncallCreate = `
INSERT INTO soma.group_oncall_properties (
            instance_id,
            source_instance_id,
            group_id,
            view,
            oncall_duty_id,
            repository_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::uuid,
       $7::boolean,
       $8::boolean;`

const TxGroupPropertyOncallDelete = `
DELETE FROM soma.group_oncall_properties
WHERE       instance_id = $1::uuid;`

const TxGroupPropertyServiceCreate = `
INSERT INTO soma.group_service_properties (
            instance_id,
            source_instance_id,
            group_id,
            view,
            service_property,
            organizational_team_id,
            repository_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::uuid,
       $7::uuid,
       $8::boolean,
       $9::boolean;`

const TxGroupPropertyServiceDelete = `
DELETE FROM soma.group_service_properties
WHERE       instance_id = $1::uuid;`

const TxGroupPropertySystemCreate = `
INSERT INTO soma.group_system_properties (
            instance_id,
            source_instance_id,
            group_id,
            view,
            system_property,
            source_type,
            repository_id,
            inheritance_enabled,
            children_only,
            value,
            inherited)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::varchar,
       $7::uuid,
       $8::boolean,
       $9::boolean,
       $10::text,
       $11::boolean;`

const TxGroupPropertySystemDelete = `
DELETE FROM soma.group_system_properties
WHERE       instance_id = $1::uuid;`

const TxGroupPropertyCustomCreate = `
INSERT INTO soma.group_custom_properties (
            instance_id,
            source_instance_id,
            group_id,
            view,
            custom_property_id,
            bucket_id,
            repository_id,
            inheritance_enabled,
            children_only,
            value)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::uuid,
       $7::uuid,
       $8::boolean,
       $9::boolean,
       $10::text;`

const TxGroupPropertyCustomDelete = `
DELETE FROM soma.group_custom_properties
WHERE       instance_id = $1::uuid;`

const TxClusterCreate = `
INSERT INTO soma.clusters (
            cluster_id,
            cluster_name,
            bucket_id,
            object_state,
            organizational_team_id,
            created_by)
SELECT $1::uuid,
       $2::varchar,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       user_id
FROM   inventory.users iu
WHERE  iu.user_uid = $6::varchar;`

const TxClusterUpdate = `
UPDATE soma.clusters
SET    object_state = $2::varchar
WHERE  cluster_id = $1::uuid;`

const TxClusterDelete = `
DELETE FROM soma.clusters
WHERE       cluster_id = $1::uuid;`

const TxClusterMemberNew = `
INSERT INTO soma.cluster_membership (
            cluster_id,
            node_id,
            bucket_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

const TxClusterMemberRemove = `
DELETE FROM soma.cluster_membership
WHERE       cluster_id = $1::uuid
AND         node_id = $2::uuid;`

const TxClusterPropertyOncallCreate = `
INSERT INTO soma.cluster_oncall_properties (
            instance_id,
            source_instance_id,
            cluster_id,
            view,
            oncall_duty_id,
            repository_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::uuid,
       $7::boolean,
       $8::boolean;`

const TxClusterPropertyOncallDelete = `
DELETE FROM soma.cluster_oncall_properties
WHERE       instance_id = $1::uuid;`

const TxClusterPropertyServiceCreate = `
INSERT INTO soma.cluster_service_properties (
            instance_id,
            source_instance_id,
            cluster_id,
            view,
            service_property,
            organizational_team_id,
            repository_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::uuid,
       $7::uuid,
       $8::boolean,
       $9::boolean;`

const TxClusterPropertyServiceDelete = `
DELETE FROM soma.cluster_service_properties
WHERE       instance_id = $1::uuid;`

const TxClusterPropertySystemCreate = `
INSERT INTO soma.cluster_system_properties (
            instance_id,
            source_instance_id,
            cluster_id,
            view,
            system_property,
            source_type,
            repository_id,
            inheritance_enabled,
            children_only,
            value,
            inherited)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::varchar,
       $7::uuid,
       $8::boolean,
       $9::boolean,
       $10::text,
       $11::boolean;`

const TxClusterPropertySystemDelete = `
DELETE FROM soma.cluster_system_properties
WHERE       instance_id = $1::uuid;`

const TxClusterPropertyCustomCreate = `
INSERT INTO soma.cluster_custom_properties (
            instance_id,
            source_instance_id,
            cluster_id,
            view,
            custom_property_id,
            bucket_id,
            repository_id,
            inheritance_enabled,
            children_only,
            value)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::uuid,
       $7::uuid,
       $8::boolean,
       $9::boolean,
       $10::text;`

const TxClusterPropertyCustomDelete = `
DELETE FROM soma.cluster_custom_properties
WHERE       instance_id = $1::uuid;`

const TxCreateBucket = `
INSERT INTO soma.buckets (
            bucket_id,
            bucket_name,
            bucket_frozen,
            bucket_deleted,
            repository_id,
            environment,
            organizational_team_id,
            created_by)
SELECT $1::uuid,
       $2::varchar,
       $3::boolean,
       $4::boolean,
       $5::uuid,
       $6::varchar,
       $7::uuid,
       user_id
FROM   inventory.users iu
WHERE  iu.user_uid = $8::varchar;`

const TxBucketAssignNode = `
INSERT INTO soma.node_bucket_assignment (
            node_id,
            bucket_id,
            organizational_team_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

const TxBucketRemoveNode = `
DELETE FROM soma.node_bucket_assignment (
WHERE       node_id = $1::uuid
AND         bucket_id = $2::uuid
AND         organizational_team_id = $3::uuid;`

const TxBucketPropertyOncallCreate = `
INSERT INTO soma.bucket_oncall_properties (
            instance_id,
            source_instance_id,
            bucket_id,
            view,
            oncall_duty_id,
            repository_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::uuid,
       $7::boolean,
       $8::boolean;`

const TxBucketPropertyOncallDelete = `
DELETE FROM soma.bucket_oncall_properties (
WHERE       instance_id = $1::uuid;`

const TxBucketPropertyServiceCreate = `
INSERT INTO soma.bucket_service_properties (
            instance_id,
            source_instance_id,
            bucket_id,
            view,
            service_property,
            organizational_team_id,
            repository_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::uuid,
       $7::uuid,
       $8::boolean,
       $9::boolean;`

const TxBucketPropertyServiceDelete = `
DELETE FROM soma.bucket_service_properties
WHERE       instance_id = $1::uuid;`

const TxBucketPropertySystemCreate = `
INSERT INTO soma.bucket_system_properties (
            instance_id,
            source_instance_id,
            bucket_id,
            view,
            system_property,
            source_type,
            repository_id,
            inheritance_enabled,
            children_only,
            value,
            inherited)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::varchar,
       $7::uuid,
       $8::boolean,
       $9::boolean,
       $10::text,
       $11::boolean;`

const TxBucketPropertySystemDelete = `
DELETE FROM soma.bucket_system_properties
WHERE       instance_id = $1::uuid;`

const TxBucketPropertyCustomCreate = `
INSERT INTO soma.bucket_custom_properties (
            instance_id,
            source_instance_id,
            bucket_id,
            view,
            custom_property_id,
            repository_id,
            inheritance_enabled,
            children_only,
            value)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::uuid,
       $7::boolean,
       $8::boolean,
       $9::text;`

const TxBucketPropertyCustomDelete = `
DELETE FROM soma.bucket_custom_properties
WHERE       instance_id = $1::uuid;`

const TxDeployDetailsComputeList = `
SELECT scic.check_instance_config_id
FROM   soma.checks sc
JOIN   soma.check_instances sci
  ON   sc.check_id = sci.check_id
JOIN   soma.check_instance_configurations scic
  ON   sci.check_instance_id = scic.check_instance_id
WHERE  scic.status = 'awaiting_computation'
  AND  sc.repository_id = $1::uuid;`

const TxDeployDetailsUpdate = `
UPDATE soma.check_instance_configurations
SET    status = 'computed',
       deployment_details = $1::jsonb,
	   monitoring_id = $2::uuid
WHERE  check_instance_config_id = $3::uuid;`

const TxDeployDetailsCheckInstance = `
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

const TxDeployDetailsCheck = `
SELECT sc.repository_id,
	   sc.source_check_id,
	   sc.source_object_type,
	   sc.source_object_id,
	   sc.capability_id,
	   sc.object_id,
	   sc.object_type,
	   scc.inheritance_enabled,
	   scc.children_only
FROM   soma.checks sc
JOIN   soma.check_configurations scc
ON     sc.configuration_id = scc.configuration_id
WHERE  sc.check_id = $1::uuid;`

const TxDeployDetailsCheckConfig = `
SELECT configuration_name,
       interval,
	   configuration_active,
	   enabled,
	   external_id
FROM   soma.check_configurations
WHERE  configuration_id = $1::uuid;`

const TxDeployDetailsCheckConfigThreshold = `
SELECT sct.predicate,
	   sct.threshold,
	   sct.notification_level,
	   snl.level_shortname,
	   snl.level_numeric
FROM   soma.configuration_thresholds sct
JOIN   soma.notification_levels snl
ON     sct.notification_level = snl.level_name
WHERE  sct.configuration_id = $1::uuid;`

const TxDeployDetailsCapabilityMonitoringMetric = `
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

const TxDeployDetailsProviders = `
SELECT metric_provider,
       package
FROM   soma.metric_packages
WHERE  metric = $1::varchar;`

const TxDeployDetailsGroup = `
SELECT sg.bucket_id,
	   sg.group_name,
	   sg.object_state,
	   sg.organizational_team_id,
	   sb.bucket_name,
	   sb.environment,
	   sr.repository_name
FROM   soma.groups sg
JOIN   soma.buckets sb
ON     sg.bucket_id = sb.bucket_id
JOIN   soma.repositories sr
ON     sb.repository_id = sr.repository_id
WHERE  sg.group_id = $1::uuid;`

const TxDeployDetailsCluster = `
SELECT sc.cluster_name,
       sc.bucket_id,
	   sc.object_state,
	   sc.organizational_team_id,
	   sb.bucket_name,
	   sb.environment,
	   sr.repository_name
FROM   soma.clusters sc
JOIN   soma.buckets sb
ON     sc.bucket_id = sb.bucket_id
JOIN   soma.repositories sr
ON     sb.repository_id = sr.repository_id
WHERE  sc.cluster_id = $1::uuid;`

const TxDeployDetailsNode = `
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

const TxDeployDetailsTeam = `
SELECT organizational_team_name,
       organizational_team_ldap_id
FROM   inventory.organizational_teams
WHERE  organizational_team_id = $1::uuid;`

const TxDeployDetailsNodeOncall = `
SELECT iodt.oncall_duty_id,
       iodt.oncall_duty_name,
       iodt.oncall_duty_phone_number
FROM   soma.node_oncall_property snop
JOIN   inventory.oncall_duty_teams iodt
ON     snop.oncall_duty_id = iodt.oncall_duty_id
WHERE  snop.node_id = $1::uuid
AND    snop.view = $2::varchar;`

const TxDeployDetailsClusterOncall = `
SELECT iodt.oncall_duty_id,
       iodt.oncall_duty_name,
	   iodt.oncall_duty_phone_number
FROM   soma.cluster_oncall_properties scop
JOIN   inventory.oncall_duty_teams iodt
ON     scop.oncall_duty_id = iodt.oncall_duty_id
WHERE  scop.cluster_id = $1::uuid
AND    (scop.view = $2::varchar OR scop.view = 'any');`

const TxDeployDetailsGroupOncall = `
SELECT iodt.oncall_duty_id,
       iodt.oncall_duty_name,
	   iodt.oncall_duty_phone_number
FROM   soma.group_oncall_properties sgop
JOIN   inventory.oncall_duty_teams iodt
ON     sgop.oncall_duty_id = iodt.oncall_duty_id
WHERE  sgop.group_id = $1::uuid
AND    (sgop.view = $2::varchar OR sgop.view = 'any');`

const TxDeployDetailsGroupService = `
SELECT service_property,
       organizational_team_id
FROM   soma.group_service_properties
WHERE  instance_id = $1::uuid
AND    (view = $2::varchar OR view = 'any');`

const TxDeployDetailsClusterService = `
SELECT service_property,
       organizational_team_id
FROM   soma.cluster_service_properties
WHERE  instance_id = $1::uuid
AND    (view = $2::varchar OR view = 'any');`

const TxDeployDetailsNodeService = `
SELECT service_property,
       organizational_team_id
FROM   soma.node_service_properties
WHERE  instance_id = $1::uuid
AND    (view = $2::varchar OR view = 'any');`

const TxDeployDetailsGroupSysProp = `
SELECT system_property,
       value
FROM   soma.group_system_properties
WHERE  group_id = $1::uuid
AND    (view = $2::varchar OR view = 'any');`

const TxDeployDetailsGroupCustProp = `
SELECT sgcp.custom_property_id,
       scp.custom_property,
       sgcp.value
FROM   soma.group_custom_properties sgcp
JOIN   soma.custom_properties scp
ON     sgcp.custom_property_id = scp.custom_property_id
AND    sgcp.repository_id = scp.repository_id
WHERE  sgcp.group_id = $1::uuid
AND    (sgcp.view = $2::varchar OR sgcp.view = 'any');`

const TxDeployDetailClusterSysProp = `
SELECT system_property,
       value
FROM   soma.cluster_system_properties
WHERE  cluster_id = $1::uuid
AND    (view = $2::varchar OR view = 'any');`

const TxDeployDetailClusterCustProp = `
SELECT sccp.custom_property_id,
       scp.custom_property,
	   sccp.value
FROM   soma.cluster_custom_properties sccp
JOIN   soma.custom_properties scp
ON     sccp.custom_property_id = scp.custom_property_id
AND    sccp.repository_id = scp.repository_id
WHERE  sccp.cluster_id = $1::uuid
AND    (sccp.view = $2::varchar OR sccp.view = 'any');`

const TxDeployDetailNodeSysProp = `
SELECT system_property,
       value
FROM   soma.node_system_properties
WHERE  node_id = $1::uuid
AND    (view = $2::varchar OR view = 'any');`

const TxDeployDetailNodeCustProp = `
SELECT sncp.custom_property_id,
       scp.custom_property,
	   sncp.value
FROM   soma.node_custom_properties sncp
JOIN   soma.custom_properties scp
ON     sncp.custom_property_id = scp.custom_property_id
AND    sncp.repository_id = scp.repository_id
WHERE  sncp.node_id = $1::uuid
AND    (sncp.view = $2::varchar OR sncp.view = 'any');`

const TxDeployDetailDefaultDatacenter = `
SELECT server_datacenter_name
FROM   inventory.servers
WHERE  server_asset_id = 0
AND    server_datacenter_location = 'none'
AND    server_name = 'soma-null-server';`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
