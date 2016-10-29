/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const (
	TreekeeperStartupStatements = ``

	TkStartLoadChecks = `
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
AND    source_object_type = $2::varchar
AND    NOT deleted;`

	TkStartLoadInheritedChecks = `
SELECT check_id,
       object_id,
       object_type
FROM   soma.checks
WHERE  repository_id = $1::uuid
AND    source_check_id = $2::uuid
AND    source_check_id != check_id
AND    NOT deleted;`

	TkStartLoadChecksForType = `
SELECT check_id,
       object_id
FROM   soma.checks
WHERE  repository_id = $1::uuid
AND    object_type = $2::varchar
AND    NOT deleted;`

	TkStartLoadCheckConfiguration = `
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
AND    repository_id = $2::uuid
AND    NOT deleted;`

	TkStartLoadAllCheckConfigurationsForType = `
SELECT configuration_id,
       bucket_id,
       configuration_name,
       configuration_object,
       inheritance_enabled,
       children_only,
       capability_id,
       interval,
       enabled,
       external_id
FROM   soma.check_configurations
WHERE  configuration_object_type = $1::varchar
AND    repository_id = $2::uuid
AND    NOT deleted;`

	TkStartLoadCheckThresholds = `
SELECT sct.predicate,
       sct.threshold,
       snl.level_name,
       snl.level_shortname,
       snl.level_numeric
FROM   soma.configuration_thresholds sct
JOIN   soma.notification_levels snl
ON     sct.notification_level = snl.level_name
WHERE  configuration_id = $1::uuid;`

	TkStartLoadCheckConstraintCustom = `
SELECT sccp.custom_property_id,
       scp.custom_property,
       sccp.property_value
FROM   soma.constraints_custom_property sccp
JOIN   soma.custom_properties scp
ON     sccp.custom_property_id = scp.custom_property_id
AND    sccp.repository_id = scp.repository_id
WHERE  configuration_id = $1::uuid;`

	// do not get distracted by the squirrels! All constraint
	// statements are constructed to use three result variables,
	// so they can be loaded in one unified loop.
	TkStartLoadCheckConstraintNative = `
SELECT native_property,
       property_value,
       'squirrel'
FROM   soma.constraints_native_property
WHERE  configuration_id = $1::uuid;`

	// return configuration id: every constraint query has 2 columns
	TkStartLoadCheckConstraintOncall = `
SELECT scop.oncall_duty_id,
       oncall_duty_name,
       oncall_duty_phone_number
FROM   soma.constraints_oncall_property scop
JOIN   inventory.oncall_duty_teams iodt
ON     scop.oncall_duty_id = iodt.oncall_duty_id
WHERE  scop.configuration_id = $1::uuid;`

	TkStartLoadCheckConstraintAttribute = `
SELECT service_property_attribute,
       attribute_value,
       'squirrel'
FROM   soma.constraints_service_attribute
WHERE  configuration_id = $1::uuid;`

	TkStartLoadCheckConstraintService = `
SELECT organizational_team_id,
       service_property,
       'squirrel'
FROM   soma.constraints_service_property
WHERE  configuration_id = $1::uuid;`

	TkStartLoadCheckConstraintSystem = `
SELECT system_property,
       property_value,
       'squirrel'
FROM   soma.constraints_system_property
WHERE  configuration_id = $1::uuid;`

	TkStartLoadCheckInstances = `
SELECT check_instance_id,
       check_configuration_id
FROM   soma.check_instances
WHERE  check_id = $1::uuid
AND    NOT deleted;`

	// load the most recent configuration for this instance, which is
	// not always the current one, since a newer version could be blocked
	// by the current versions rollout
	TkStartLoadCheckInstanceConfiguration = `
SELECT check_instance_config_id,
       version,
       monitoring_id,
       constraint_hash,
       constraint_val_hash,
       instance_service,
       instance_service_cfg_hash,
       instance_service_cfg
FROM   soma.check_instance_configurations
WHERE  check_instance_id = $1::uuid
ORDER  BY created DESC
LIMIT  1;`

	TkStartLoadCheckGroupState = `
SELECT sg.group_id,
       sg.object_state
FROM   soma.buckets sb
JOIN   soma.groups  sg
ON     sb.bucket_id = sg.bucket_id
WHERE  sb.repository_id = $1::uuid;`

	TkStartLoadCheckGroupRelations = `
SELECT sgmg.group_id,
       sgmg.child_group_id
FROM   soma.buckets sb
JOIN   soma.group_membership_groups sgmg
ON     sb.bucket_id = sgmg.bucket_id
WHERE  sb.repository_id = $1::uuid;`

	TkStartLoadBuckets = `
SELECT sb.bucket_id,
       sb.bucket_name,
       sb.bucket_frozen,
       sb.bucket_deleted,
       sb.environment,
       sb.organizational_team_id
FROM   soma.repositories sr
JOIN   soma.buckets sb
ON     sr.repository_id = sb.repository_id
WHERE  sr.repository_id = $1::uuid;`

	TkStartLoadGroups = `
SELECT sg.group_id,
       sg.group_name,
       sg.bucket_id,
       sg.organizational_team_id
FROM   soma.repositories sr
JOIN   soma.buckets sb
ON     sr.repository_id = sb.repository_id
JOIN   soma.groups sg
ON     sb.bucket_id = sg.bucket_id
WHERE  sr.repository_id = $1::uuid;`

	TkStartLoadGroupMemberGroups = `
SELECT sgmg.group_id,
       sgmg.child_group_id
FROM   soma.repositories sr
JOIN   soma.buckets sb
ON     sr.repository_id = sb.repository_id
JOIN   soma.group_membership_groups sgmg
ON     sb.bucket_id = sgmg.bucket_id
WHERE  sr.repository_id = $1::uuid;`

	TkStartLoadGroupedClusters = `
SELECT sc.cluster_id,
       sc.cluster_name,
       sc.organizational_team_id,
       sgmc.group_id
FROM   soma.repositories sr
JOIN   soma.buckets sb
ON     sr.repository_id = sb.repository_id
JOIN   soma.clusters sc
ON     sb.bucket_id = sc.bucket_id
JOIN   soma.group_membership_clusters sgmc
ON     sc.bucket_id = sgmc.bucket_id
AND    sc.cluster_id = sgmc.child_cluster_id
WHERE  sr.repository_id = $1::uuid;`

	TkStartLoadCluster = `
SELECT sc.cluster_id,
       sc.cluster_name,
       sc.bucket_id,
       sc.organizational_team_id
FROM   soma.repositories sr
JOIN   soma.buckets sb
ON     sr.repository_id = sb.repository_id
JOIN   soma.clusters sc
ON     sb.bucket_id = sc.bucket_id
WHERE  sr.repository_id = $1::uuid
AND    sc.object_state != 'grouped';`

	TkStartLoadNode = `
SELECT    sn.node_id,
          sn.node_asset_id,
          sn.node_name,
          sn.organizational_team_id,
          sn.server_id,
          sn.node_online,
          sn.node_deleted,
          snba.bucket_id,
          scm.cluster_id,
          sgmn.group_id
FROM      soma.repositories sr
JOIN      soma.buckets sb
ON        sr.repository_id = sb.repository_id
JOIN      soma.node_bucket_assignment snba
ON        sb.bucket_id = snba.bucket_id
JOIN      soma.nodes sn
ON        snba.node_id = sn.node_id
LEFT JOIN soma.cluster_membership scm
ON        sn.node_id = scm.node_id
LEFT JOIN soma.group_membership_nodes sgmn
ON        sn.node_id = sgmn.child_node_id
WHERE     sr.repository_id = $1::uuid`

	TkStartLoadJob = `
SELECT   job
FROM     soma.jobs
WHERE    repository_id = $1::uuid
AND      job_status != 'processed'
ORDER BY job_serial ASC;`

	TkStartLoadSystemPropInstances = `
SELECT      CASE WHEN srsp.instance_id IS NOT NULL THEN srsp.instance_id
            ELSE CASE WHEN sbsp.instance_id IS NOT NULL THEN sbsp.instance_id
                 ELSE CASE WHEN sgsp.instance_id IS NOT NULL THEN sgsp.instance_id
                      ELSE CASE WHEN scsp.instance_id IS NOT NULL THEN scsp.instance_id
                           ELSE CASE WHEN snsp.instance_id IS NOT NULL THEN snsp.instance_id
                                ELSE '00000000-0000-0000-0000-000000000000'
                                END
                           END
                      END
                 END
            END AS "instance_id",
            CASE WHEN srsp.object_type IS NOT NULL THEN srsp.object_type
            ELSE CASE WHEN sbsp.object_type IS NOT NULL THEN sbsp.object_type
                 ELSE CASE WHEN sgsp.object_type IS NOT NULL THEN sgsp.object_type
                      ELSE CASE WHEN scsp.object_type IS NOT NULL THEN scsp.object_type
                           ELSE CASE WHEN snsp.object_type IS NOT NULL THEN snsp.object_type
                                ELSE 'MAGIC_NO_RESULT_VALUE'
                                END
                           END
                      END
                 END
            END AS "object_type",
            CASE WHEN srsp.repository_id IS NOT NULL THEN srsp.repository_id
            ELSE CASE WHEN sbsp.bucket_id IS NOT NULL THEN sbsp.bucket_id
                 ELSE CASE WHEN sgsp.group_id IS NOT NULL THEN sgsp.group_id
                      ELSE CASE WHEN scsp.cluster_id IS NOT NULL THEN scsp.cluster_id
                           ELSE CASE WHEN snsp.node_id IS NOT NULL THEN snsp.node_id
                                ELSE '00000000-0000-0000-0000-000000000000'
                                END
                           END
                      END
                 END
            END AS "object_id"
FROM        soma.property_instances spi
LEFT JOIN   soma.repository_system_properties srsp
  ON        spi.instance_id = srsp.instance_id
  AND       spi.source_instance_id = srsp.source_instance_id
LEFT JOIN   soma.bucket_system_properties sbsp
  ON        spi.instance_id = sbsp.instance_id
  AND       spi.source_instance_id = sbsp.source_instance_id
LEFT JOIN   soma.group_system_properties sgsp
  ON        spi.instance_id = sgsp.instance_id
  AND       spi.source_instance_id = sgsp.source_instance_id
LEFT JOIN   soma.cluster_system_properties scsp
  ON        spi.instance_id = scsp.instance_id
  AND       spi.source_instance_id = scsp.source_instance_id
LEFT JOIN   soma.node_system_properties snsp
  ON        spi.instance_id = snsp.instance_id
  AND       spi.source_instance_id = snsp.source_instance_id
WHERE       spi.instance_id != spi.source_instance_id
  AND       spi.repository_id = $1::uuid
  AND       spi.source_instance_id = $2::uuid;`

	TkStartLoadCustomPropInstances = `
SELECT      CASE WHEN srsp.instance_id IS NOT NULL THEN srsp.instance_id
            ELSE CASE WHEN sbsp.instance_id IS NOT NULL THEN sbsp.instance_id
                 ELSE CASE WHEN sgsp.instance_id IS NOT NULL THEN sgsp.instance_id
                      ELSE CASE WHEN scsp.instance_id IS NOT NULL THEN scsp.instance_id
                           ELSE CASE WHEN snsp.instance_id IS NOT NULL THEN snsp.instance_id
                                ELSE '00000000-0000-0000-0000-000000000000'
                                END
                           END
                      END
                 END
            END AS "instance_id",
            CASE WHEN srsp.instance_id IS NOT NULL THEN 'repository'
            ELSE CASE WHEN sbsp.instance_id IS NOT NULL THEN 'bucket'
                 ELSE CASE WHEN sgsp.instance_id IS NOT NULL THEN 'group'
                      ELSE CASE WHEN scsp.instance_id IS NOT NULL THEN 'cluster'
                           ELSE CASE WHEN snsp.instance_id IS NOT NULL THEN 'node'
                                ELSE 'MAGIC_NO_RESULT_VALUE'
                                END
                           END
                      END
                 END
            END AS "object_type",
            CASE WHEN srsp.repository_id IS NOT NULL THEN srsp.repository_id
            ELSE CASE WHEN sbsp.bucket_id IS NOT NULL THEN sbsp.bucket_id
                 ELSE CASE WHEN sgsp.group_id IS NOT NULL THEN sgsp.group_id
                      ELSE CASE WHEN scsp.cluster_id IS NOT NULL THEN scsp.cluster_id
                           ELSE CASE WHEN snsp.node_id IS NOT NULL THEN snsp.node_id
                                ELSE '00000000-0000-0000-0000-000000000000'
                                END
                           END
                      END
                 END
            END AS "object_id"
FROM        soma.property_instances spi
LEFT JOIN   soma.repository_custom_properties srsp
  ON        spi.instance_id = srsp.instance_id
  AND       spi.source_instance_id = srsp.source_instance_id
LEFT JOIN   soma.bucket_custom_properties sbsp
  ON        spi.instance_id = sbsp.instance_id
  AND       spi.source_instance_id = sbsp.source_instance_id
LEFT JOIN   soma.group_custom_properties sgsp
  ON        spi.instance_id = sgsp.instance_id
  AND       spi.source_instance_id = sgsp.source_instance_id
LEFT JOIN   soma.cluster_custom_properties scsp
  ON        spi.instance_id = scsp.instance_id
  AND       spi.source_instance_id = scsp.source_instance_id
LEFT JOIN   soma.node_custom_properties snsp
  ON        spi.instance_id = snsp.instance_id
  AND       spi.source_instance_id = snsp.source_instance_id
WHERE       spi.instance_id != spi.source_instance_id
  AND       spi.repository_id = $1::uuid
  AND       spi.source_instance_id = $2::uuid;`

	TkStartLoadServicePropInstances = `
SELECT      CASE WHEN srsp.instance_id IS NOT NULL THEN srsp.instance_id
            ELSE CASE WHEN sbsp.instance_id IS NOT NULL THEN sbsp.instance_id
                 ELSE CASE WHEN sgsp.instance_id IS NOT NULL THEN sgsp.instance_id
                      ELSE CASE WHEN scsp.instance_id IS NOT NULL THEN scsp.instance_id
                           ELSE CASE WHEN snsp.instance_id IS NOT NULL THEN snsp.instance_id
                                ELSE '00000000-0000-0000-0000-000000000000'
                                END
                           END
                      END
                 END
            END AS "instance_id",
            CASE WHEN srsp.instance_id IS NOT NULL THEN 'repository'
            ELSE CASE WHEN sbsp.instance_id IS NOT NULL THEN 'bucket'
                 ELSE CASE WHEN sgsp.instance_id IS NOT NULL THEN 'group'
                      ELSE CASE WHEN scsp.instance_id IS NOT NULL THEN 'cluster'
                           ELSE CASE WHEN snsp.instance_id IS NOT NULL THEN 'node'
                                ELSE 'MAGIC_NO_RESULT_VALUE'
                                END
                           END
                      END
                 END
            END AS "object_type",
            CASE WHEN srsp.repository_id IS NOT NULL THEN srsp.repository_id
            ELSE CASE WHEN sbsp.bucket_id IS NOT NULL THEN sbsp.bucket_id
                 ELSE CASE WHEN sgsp.group_id IS NOT NULL THEN sgsp.group_id
                      ELSE CASE WHEN scsp.cluster_id IS NOT NULL THEN scsp.cluster_id
                           ELSE CASE WHEN snsp.node_id IS NOT NULL THEN snsp.node_id
                                ELSE '00000000-0000-0000-0000-000000000000'
                                END
                           END
                      END
                 END
            END AS "object_id"
FROM        soma.property_instances spi
LEFT JOIN   soma.repository_service_properties srsp
  ON        spi.instance_id = srsp.instance_id
  AND       spi.source_instance_id = srsp.source_instance_id
LEFT JOIN   soma.bucket_service_properties sbsp
  ON        spi.instance_id = sbsp.instance_id
  AND       spi.source_instance_id = sbsp.source_instance_id
LEFT JOIN   soma.group_service_properties sgsp
  ON        spi.instance_id = sgsp.instance_id
  AND       spi.source_instance_id = sgsp.source_instance_id
LEFT JOIN   soma.cluster_service_properties scsp
  ON        spi.instance_id = scsp.instance_id
  AND       spi.source_instance_id = scsp.source_instance_id
LEFT JOIN   soma.node_service_properties snsp
  ON        spi.instance_id = snsp.instance_id
  AND       spi.source_instance_id = snsp.source_instance_id
WHERE       spi.instance_id != spi.source_instance_id
  AND       spi.repository_id = $1::uuid
  AND       spi.source_instance_id = $2::uuid;`

	TkStartLoadOncallPropInstances = `
SELECT      CASE WHEN srsp.instance_id IS NOT NULL THEN srsp.instance_id
            ELSE CASE WHEN sbsp.instance_id IS NOT NULL THEN sbsp.instance_id
                 ELSE CASE WHEN sgsp.instance_id IS NOT NULL THEN sgsp.instance_id
                      ELSE CASE WHEN scsp.instance_id IS NOT NULL THEN scsp.instance_id
                           ELSE CASE WHEN snsp.instance_id IS NOT NULL THEN snsp.instance_id
                                ELSE '00000000-0000-0000-0000-000000000000'
                                END
                           END
                      END
                 END
            END AS "instance_id",
            CASE WHEN srsp.instance_id IS NOT NULL THEN 'repository'
            ELSE CASE WHEN sbsp.instance_id IS NOT NULL THEN 'bucket'
                 ELSE CASE WHEN sgsp.instance_id IS NOT NULL THEN 'group'
                      ELSE CASE WHEN scsp.instance_id IS NOT NULL THEN 'cluster'
                           ELSE CASE WHEN snsp.instance_id IS NOT NULL THEN 'node'
                                ELSE 'MAGIC_NO_RESULT_VALUE'
                                END
                           END
                      END
                 END
            END AS "object_type",
            CASE WHEN srsp.repository_id IS NOT NULL THEN srsp.repository_id
            ELSE CASE WHEN sbsp.bucket_id IS NOT NULL THEN sbsp.bucket_id
                 ELSE CASE WHEN sgsp.group_id IS NOT NULL THEN sgsp.group_id
                      ELSE CASE WHEN scsp.cluster_id IS NOT NULL THEN scsp.cluster_id
                           ELSE CASE WHEN snsp.node_id IS NOT NULL THEN snsp.node_id
                                ELSE '00000000-0000-0000-0000-000000000000'
                                END
                           END
                      END
                 END
            END AS "object_id"
FROM        soma.property_instances spi
LEFT JOIN   soma.repository_oncall_properties srsp
  ON        spi.instance_id = srsp.instance_id
  AND       spi.source_instance_id = srsp.source_instance_id
LEFT JOIN   soma.bucket_oncall_properties sbsp
  ON        spi.instance_id = sbsp.instance_id
  AND       spi.source_instance_id = sbsp.source_instance_id
LEFT JOIN   soma.group_oncall_properties sgsp
  ON        spi.instance_id = sgsp.instance_id
  AND       spi.source_instance_id = sgsp.source_instance_id
LEFT JOIN   soma.cluster_oncall_properties scsp
  ON        spi.instance_id = scsp.instance_id
  AND       spi.source_instance_id = scsp.source_instance_id
LEFT JOIN   soma.node_oncall_property snsp
  ON        spi.instance_id = snsp.instance_id
  AND       spi.source_instance_id = snsp.source_instance_id
WHERE       spi.instance_id != spi.source_instance_id
  AND       spi.repository_id = $1::uuid
  AND       spi.source_instance_id = $2::uuid;`

	TkStartLoadRepositoryCstProp = `
SELECT srcp.instance_id,
       srcp.source_instance_id,
       srcp.repository_id,
       srcp.view,
       srcp.custom_property_id,
       srcp.inheritance_enabled,
       srcp.children_only,
       srcp.value,
       scp.custom_property
FROM   soma.repository_custom_properties srcp
JOIN   soma.custom_properties scp
ON     srcp.custom_property_id = scp.custom_property_id
WHERE  srcp.instance_id = srcp.source_instance_id
AND    srcp.repository_id = $1::uuid;`

	TkStartLoadBucketCstProp = `
SELECT sbcp.instance_id,
       sbcp.source_instance_id,
       sbcp.bucket_id,
       sbcp.view,
       sbcp.custom_property_id,
       sbcp.inheritance_enabled,
       sbcp.children_only,
       sbcp.value,
       scp.custom_property
FROM   soma.bucket_custom_properties sbcp
JOIN   soma.custom_properties scp
ON     sbcp.custom_property_id = scp.custom_property_id
WHERE  sbcp.instance_id = sbcp.source_instance_id
AND    sbcp.repository_id = $1::uuid;`

	TkStartLoadGroupCstProp = `
SELECT sgcp.instance_id,
       sgcp.source_instance_id,
       sgcp.group_id,
       sgcp.view,
       sgcp.custom_property_id,
       sgcp.inheritance_enabled,
       sgcp.children_only,
       sgcp.value,
       scp.custom_property
FROM   soma.group_custom_properties sgcp
JOIN   soma.custom_properties scp
ON     sgcp.custom_property_id = scp.custom_property_id
WHERE  sgcp.instance_id = sgcp.source_instance_id
AND    sgcp.repository_id = $1::uuid;`

	TkStartLoadClusterCstProp = `
SELECT sccp.instance_id,
       sccp.source_instance_id,
       sccp.cluster_id,
       sccp.view,
       sccp.custom_property_id,
       sccp.inheritance_enabled,
       sccp.children_only,
       sccp.value,
       scp.custom_property
FROM   soma.cluster_custom_properties sccp
JOIN   soma.custom_properties scp
ON     sccp.custom_property_id = scp.custom_property_id
WHERE  sccp.instance_id = sccp.source_instance_id
AND    sccp.repository_id = $1::uuid;`

	TkStartLoadNodeCstProp = `
SELECT sncp.instance_id,
       sncp.source_instance_id,
       sncp.node_id,
       sncp.view,
       sncp.custom_property_id,
       sncp.inheritance_enabled,
       sncp.children_only,
       sncp.value,
       scp.custom_property
FROM   soma.node_custom_properties sncp
JOIN   soma.custom_properties scp
ON     sncp.custom_property_id = scp.custom_property_id
WHERE  sncp.instance_id = sncp.source_instance_id
AND    sncp.repository_id = $1::uuid;`

	TkStartLoadRepoOncProp = `
SELECT  srop.instance_id,
        srop.source_instance_id,
        srop.repository_id,
        srop.view,
        srop.oncall_duty_id,
        srop.inheritance_enabled,
        srop.children_only,
        iodt.oncall_duty_name,
        iodt.oncall_duty_phone_number
FROM    soma.repository_oncall_properties srop
JOIN    inventory.oncall_duty_teams iodt
  ON    srop.oncall_duty_id = iodt.oncall_duty_id
WHERE   srop.instance_id = srop.source_instance_id
  AND   srop.repository_id = $1::uuid;`

	TkStartLoadBucketOncProp = `
SELECT  sgop.instance_id,
        sgop.source_instance_id,
        sgop.bucket_id,
        sgop.view,
        sgop.oncall_duty_id,
        sgop.inheritance_enabled,
        sgop.children_only,
        iodt.oncall_duty_name,
        iodt.oncall_duty_phone_number
FROM    soma.bucket_oncall_properties sgop
JOIN    inventory.oncall_duty_teams iodt
  ON    sgop.oncall_duty_id = iodt.oncall_duty_id
WHERE   sgop.instance_id = sgop.source_instance_id
  AND   sgop.repository_id = $1::uuid;`

	TkStartLoadGroupOncProp = `
SELECT  sgop.instance_id,
        sgop.source_instance_id,
        sgop.group_id,
        sgop.view,
        sgop.oncall_duty_id,
        sgop.inheritance_enabled,
        sgop.children_only,
        iodt.oncall_duty_name,
        iodt.oncall_duty_phone_number
FROM    soma.group_oncall_properties sgop
JOIN    inventory.oncall_duty_teams iodt
  ON    sgop.oncall_duty_id = iodt.oncall_duty_id
WHERE   sgop.instance_id = sgop.source_instance_id
  AND   sgop.repository_id = $1::uuid;`

	TkStartLoadClusterOncProp = `
SELECT  scop.instance_id,
        scop.source_instance_id,
        scop.cluster_id,
        scop.view,
        scop.oncall_duty_id,
        scop.inheritance_enabled,
        scop.children_only,
        iodt.oncall_duty_name,
        iodt.oncall_duty_phone_number
FROM    soma.cluster_oncall_properties scop
JOIN    inventory.oncall_duty_teams iodt
  ON    scop.oncall_duty_id = iodt.oncall_duty_id
WHERE   scop.instance_id = scop.source_instance_id
  AND   scop.repository_id = $1::uuid;`

	TkStartLoadNodeOncProp = `
SELECT  snop.instance_id,
        snop.source_instance_id,
        snop.node_id,
        snop.view,
        snop.oncall_duty_id,
        snop.inheritance_enabled,
        snop.children_only,
        iodt.oncall_duty_name,
        iodt.oncall_duty_phone_number
FROM    soma.node_oncall_property snop
JOIN    inventory.oncall_duty_teams iodt
  ON    snop.oncall_duty_id = iodt.oncall_duty_id
WHERE   snop.instance_id = snop.source_instance_id
  AND   snop.repository_id = $1::uuid;`

	TkStartLoadRepoSvcProp = `
SELECT instance_id,
       source_instance_id,
       repository_id,
       view,
       service_property,
       organizational_team_id,
       inheritance_enabled,
       children_only
FROM   soma.repository_service_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`

	TkStartLoadRepoSvcAttr = `
SELECT service_property_attribute,
       value
FROM   soma.team_service_property_values
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`

	TkStartLoadBucketSvcProp = `
SELECT instance_id,
       source_instance_id,
       bucket_id,
       view,
       service_property,
       organizational_team_id,
       inheritance_enabled,
       children_only
FROM   soma.bucket_service_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`

	TkStartLoadBucketSvcAttr = `
SELECT service_property_attribute,
       value
FROM   soma.team_service_property_values
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`

	TkStartLoadGroupSvcProp = `
SELECT instance_id,
       source_instance_id,
       group_id,
       view,
       service_property,
       organizational_team_id,
       inheritance_enabled,
       children_only
FROM   soma.group_service_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`

	TkStartLoadGroupSvcAttr = `
SELECT service_property_attribute,
       value
FROM   soma.team_service_property_values
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`

	TkStartLoadClusterSvcProp = `
SELECT instance_id,
       source_instance_id,
       cluster_id,
       view,
       service_property,
       organizational_team_id,
       inheritance_enabled,
       children_only
FROM   soma.cluster_service_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`

	TkStartLoadClusterSvcAttr = `
SELECT service_property_attribute,
       value
FROM   soma.team_service_property_values
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`

	TkStartLoadNodeSvcProp = `
SELECT instance_id,
       source_instance_id,
       node_id,
       view,
       service_property,
       organizational_team_id,
       inheritance_enabled,
       children_only
FROM   soma.node_service_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`

	TkStartLoadNodeSvcAttr = `
SELECT service_property_attribute,
       value
FROM   soma.team_service_property_values
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`

	TkStartLoadRepoSysProp = `
SELECT instance_id,
       source_instance_id,
       repository_id,
       view,
       system_property,
       source_type,
       inheritance_enabled,
       children_only,
       value
FROM   soma.repository_system_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`

	TkStartLoadBucketSysProp = `
SELECT instance_id,
       source_instance_id,
       bucket_id,
       view,
       system_property,
       source_type,
       inheritance_enabled,
       children_only,
       value
FROM   soma.bucket_system_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`

	TkStartLoadGroupSysProp = `
SELECT instance_id,
       source_instance_id,
       group_id,
       view,
       system_property,
       source_type,
       inheritance_enabled,
       children_only,
       value
FROM   soma.group_system_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`

	TkStartLoadClusterSysProp = `
SELECT instance_id,
       source_instance_id,
       cluster_id,
       view,
       system_property,
       source_type,
       inheritance_enabled,
       children_only,
       value
FROM   soma.cluster_system_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`

	TkStartLoadNodeSysProp = `
SELECT instance_id,
       source_instance_id,
       node_id,
       view,
       system_property,
       source_type,
       inheritance_enabled,
       children_only,
       value
FROM   soma.node_system_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`
)

func init() {
	m[TkStartLoadAllCheckConfigurationsForType] = `TkStartLoadAllCheckConfigurationsForType`
	m[TkStartLoadBucketCstProp] = `TkStartLoadBucketCstProp`
	m[TkStartLoadBucketOncProp] = `TkStartLoadBucketOncProp`
	m[TkStartLoadBucketSvcAttr] = `TkStartLoadBucketSvcAttr`
	m[TkStartLoadBucketSvcProp] = `TkStartLoadBucketSvcProp`
	m[TkStartLoadBucketSysProp] = `TkStartLoadBucketSysProp`
	m[TkStartLoadBuckets] = `TkStartLoadBuckets`
	m[TkStartLoadCheckConfiguration] = `TkStartLoadCheckConfiguration`
	m[TkStartLoadCheckConstraintAttribute] = `TkStartLoadCheckConstraintAttribute`
	m[TkStartLoadCheckConstraintCustom] = `TkStartLoadCheckConstraintCustom`
	m[TkStartLoadCheckConstraintNative] = `TkStartLoadCheckConstraintNative`
	m[TkStartLoadCheckConstraintOncall] = `TkStartLoadCheckConstraintOncall`
	m[TkStartLoadCheckConstraintService] = `TkStartLoadCheckConstraintService`
	m[TkStartLoadCheckConstraintSystem] = `TkStartLoadCheckConstraintSystem`
	m[TkStartLoadCheckGroupRelations] = `TkStartLoadCheckGroupRelations`
	m[TkStartLoadCheckGroupState] = `TkStartLoadCheckGroupState`
	m[TkStartLoadCheckInstanceConfiguration] = `TkStartLoadCheckInstanceConfiguration`
	m[TkStartLoadCheckInstances] = `TkStartLoadCheckInstances`
	m[TkStartLoadCheckThresholds] = `TkStartLoadCheckThresholds`
	m[TkStartLoadChecksForType] = `TkStartLoadChecksForType`
	m[TkStartLoadChecks] = `TkStartLoadChecks`
	m[TkStartLoadClusterCstProp] = `TkStartLoadClusterCstProp`
	m[TkStartLoadClusterOncProp] = `TkStartLoadClusterOncProp`
	m[TkStartLoadClusterSvcAttr] = `TkStartLoadClusterSvcAttr`
	m[TkStartLoadClusterSvcProp] = `TkStartLoadClusterSvcProp`
	m[TkStartLoadClusterSysProp] = `TkStartLoadClusterSysProp`
	m[TkStartLoadCluster] = `TkStartLoadCluster`
	m[TkStartLoadCustomPropInstances] = `TkStartLoadCustomPropInstances`
	m[TkStartLoadGroupCstProp] = `TkStartLoadGroupCstProp`
	m[TkStartLoadGroupMemberGroups] = `TkStartLoadGroupMemberGroups`
	m[TkStartLoadGroupOncProp] = `TkStartLoadGroupOncProp`
	m[TkStartLoadGroupSvcAttr] = `TkStartLoadGroupSvcAttr`
	m[TkStartLoadGroupSvcProp] = `TkStartLoadGroupSvcProp`
	m[TkStartLoadGroupSysProp] = `TkStartLoadGroupSysProp`
	m[TkStartLoadGroupedClusters] = `TkStartLoadGroupedClusters`
	m[TkStartLoadGroups] = `TkStartLoadGroups`
	m[TkStartLoadInheritedChecks] = `TkStartLoadInheritedChecks`
	m[TkStartLoadJob] = `TkStartLoadJob`
	m[TkStartLoadNodeCstProp] = `TkStartLoadNodeCstProp`
	m[TkStartLoadNodeOncProp] = `TkStartLoadNodeOncProp`
	m[TkStartLoadNodeSvcAttr] = `TkStartLoadNodeSvcAttr`
	m[TkStartLoadNodeSvcProp] = `TkStartLoadNodeSvcProp`
	m[TkStartLoadNodeSysProp] = `TkStartLoadNodeSysProp`
	m[TkStartLoadNode] = `TkStartLoadNode`
	m[TkStartLoadOncallPropInstances] = `TkStartLoadOncallPropInstances`
	m[TkStartLoadRepoOncProp] = `TkStartLoadRepoOncProp`
	m[TkStartLoadRepoSvcAttr] = `TkStartLoadRepoSvcAttr`
	m[TkStartLoadRepoSvcProp] = `TkStartLoadRepoSvcProp`
	m[TkStartLoadRepoSysProp] = `TkStartLoadRepoSysProp`
	m[TkStartLoadRepositoryCstProp] = `TkStartLoadRepositoryCstProp`
	m[TkStartLoadServicePropInstances] = `TkStartLoadServicePropInstances`
	m[TkStartLoadSystemPropInstances] = `TkStartLoadSystemPropInstances`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
