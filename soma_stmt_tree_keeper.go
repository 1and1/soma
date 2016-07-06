package main

/*
 * Statements for job state updates outside transaction
 */

const tkStmtStartJob = `
UPDATE soma.jobs
SET    job_started = $2::timestamptz,
       job_status = 'in_progress'
WHERE  job_id = $1::uuid
AND    job_started IS NULL;`

const tkStmtGetViewFromCapability = `
SELECT capability_view
FROM   soma.monitoring_capabilities
WHERE  capability_id = $1::uuid;`

const tkStmtGetComputedDeployments = `
SELECT check_instance_id,
       check_instance_config_id,
	   deployment_details
FROM   soma.check_instance_configurations
WHERE  status = 'computed';`

const tkStmtGetPreviousDeployment = `
SELECT check_instance_config_id,
       version,
	   status,
	   deployment_details
FROM   soma.check_instance_configurations
WHERE  status != 'computed'
AND    status != 'awaiting_computation'
AND    check_instance_id = $1::uuid
ORDER  BY version DESC
LIMIT  1;`

const tkStmtUpdateConfigStatus = `
UPDATE soma.check_instance_configurations
SET    status = $1::varchar,
       next_status = $2::varchar
WHERE  check_instance_config_id = $3::uuid;`

const tkStmtUpdateCheckInstance = `
UPDATE soma.check_instances
SET    last_configuration_created = $1::timestamptz,
       update_available = $2::boolean,
	   current_instance_config_id = $3::uuid
WHERE  check_instance_id = $4::uuid;`

const tkStmtUpdateExistingCheckInstance = `
UPDATE soma.check_instances
SET    last_configuration_created = $1::timestamptz,
       update_available = $2::boolean
WHERE  check_instance_id = $3::uuid;`

const tkStmtDeleteDuplicateDetails = `
DELETE FROM soma.check_instance_configurations
WHERE       check_instance_config_id = $1::uuid;`

const tkStmtSetDependency = `
INSERT INTO soma.check_instance_configuration_dependencies (
	blocked_instance_config_id,
	blocking_instance_config_id,
	unblocking_state)
SELECT $1::uuid,
       $2::uuid,
	   $3::varchar;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
