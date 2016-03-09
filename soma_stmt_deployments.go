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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
