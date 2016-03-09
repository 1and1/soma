package main

const stmtGetDeployment = `
SELECT scic.check_instance_config_id,
       scic.status,
       scic.next_status,
       scic.deployment_details
FROM   soma.check_instances sci
JOIN   soma.check_instance_configurations scic
ON     sci.check_instance_id = scic.check_instance_id
AND    sci.current_instance_config_id = scic.check_instance_config_id
WHERE  sci.check_instance_id = $1::uuid
AND    (  sci.status = 'awaiting_rollout'
       OR sci.status = 'rollout_in_progress'
	   OR sci.status = 'active'
	   OR sci.status = 'rollout_failed'
	   OR sci.status = 'awaiting_deprovision'
	   OR sci.status = 'deprovision_in_progress'
       OR sci.status = 'deprovision_failed' );`

const stmtUpdateDeployment = `
UPDATE soma.check_instance_configurations
SET    status = $1::varchar,
       next_status = $2::varchar
WHERE  check_instance_config_id = $3::uuid;`

const stmtDeploymentStatus = `
SELECT scic.check_instance_config_id,
       scic.status,
	   scic.next_status
FROM   soma.check_instances sci
JOIN   soma.check_instance_configurations scic
ON     sci.check_instance_id = scic.check_instance_id
AND    sci.current_instance_config_id = scic.check_instance_config_id
WHERE  sci.check_instance_id = $1::uuid;`

const stmtActivateDeployment = `
UPDATE soma.check_instance_configurations
SET    status = $1::varchar,
       next_status = $2::varchar,
	   activated_at = $3::timestamptz,
WHERE  check_instance_config_id = $4::uuid;`

const stmtGetDeploymentList = `
SELECT sci.instance_id
FROM   soma.monitoring_systems sms
JOIN   soma.check_instance_configurations scic
ON     sms.monitoring_id = scic.monitoring_id
JOIN   soma.check_instances sci
ON     scic.check_instance_id = sci.check_instance_id
AND    scic.check_instance_config_id = sci.current_instance_config_id
WHERE  sms.monitoring_id = $1::uuid
AND    sci.update_available
AND    (  scic.status = 'awaiting_rollout'
       OR scic.status = 'awaiting_deprovision' );`

const stmtGetAllDeploymentList = `
SELECT sci.instance_id
FROM   soma.monitoring_systems sms
JOIN   soma.check_instance_configurations scic
ON     sms.monitoring_id = scic.monitoring_id
JOIN   soma.check_instances sci
ON     scic.check_instance_id = sci.check_instance_id
AND    scic.check_instance_config_id = sci.current_instance_config_id
WHERE  sms.monitoring_id = $1::uuid
AND    (  scic.status = 'awaiting_rollout'
	   OR scic.status = 'rollout_in_progress'
       OR scic.status = 'awaiting_deprovision'
	   OR scic.status = 'deprovision_in_progress');`

const stmtDeployClearFlag = `
UPDATE soma.check_instances
SET    update_available = 'false'::boolean
WHERE  check_instance_id = $1::uuid;`

const stmtGetInstancesForNode = `
SELECT sci.check_instance_id
FROM   soma.nodes sn
JOIN   soma.checks sc
ON     sn.node_id = sc.object_id
JOIN   soma.monitoring_capabilities smc
ON     sc.capability_id = smc.monitoring_capability
JOIN   soma.check_instances sci
ON     sc.check_id = sci.check_id
WHERE  sn.node_asset_id = $1::numeric
AND    sc.object_type = 'node'
AND    smc.capability_view = 'local'
AND    smc.capability_monitoring = $2::uuid;`

const stmtGetLastInstanceVersion = `
SELECT deployment_details
FROM   soma.check_instance_configurations
WHERE  check_instance_id = $1::uuid
AND    (   status != 'deprovisioned'
       AND status != 'awaiting_deletion')
ORDER  BY version DESC
LIMIT  1;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
