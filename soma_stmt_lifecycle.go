package main

const lcStmtActiveUnblockCondition = `
SELECT 	scicd.blocked_instance_config_id,
		scicd.blocking_instance_config_id,
		scicd.unblocking_state,
		p.status,
		p.next_status,
		p.check_instance_id
FROM    soma.check_instance_configuration_dependencies scicd
JOIN    soma.check_instance_configurations scic
ON      scicd.blocking_instance_config_id = scic.check_instance_config_id
AND     scicd.unblocking_state = scic.status
JOIN    soma.check_instance_configurations p
ON      scicd.blocked_instance_config_id = p.check_instance_config_id
JOIN    soma.check_instances sci
ON      p.check_instance_id = sci.check_instance_id
AND     scicd.blocking_instance_config_id = sci.current_instance_config_id;`

const lcStmtUpdateInstance = `
UPDATE	soma.check_instances
SET     update_available = $1::boolean,
        current_instance_config_id = $2::uuid
WHERE   check_instance_id = $3::uuid;`

const lcStmtUpdateConfig = `
UPDATE  soma.check_instance_configurations
SET     status = $1::varchar,
        next_status = $2::varchar,
		awaiting_deletion = $3::boolean
WHERE   check_instance_config_id = $4::uuid;`

const lcStmtDeleteDependency = `
DELETE FROM soma.check_instance_configuration_dependencies
WHERE       blocked_instance_config_id = $1::uuid
AND         blocking_instance_config_id = $2::uuid
AND         unblocking_state = $3::varchar;`

const lcStmtReadyDeployments = `
SELECT scic.check_instance_id,
       scic.monitoring_id,
	   sms.monitoring_callback_uri
FROM   soma.check_instance_configurations scic
JOIN   soma.monitoring_systems sms
ON     scic.monitoring_id = sms.monitoring_id
JOIN   soma.check_instances sci
ON     scic.check_instance_id = sci.check_instance_id
AND    scic.check_instance_config_id = sci.current_instance_config_id
WHERE  (  scic.status = 'awaiting_rollout'
       OR scic.status = 'awaiting_deprovision' )
AND    sms.monitoring_callback_uri IS NOT NULL
AND    sci.update_available;`

const lcStmtClearUpdateFlag = `
UPDATE soma.check_instances
SET    update_available = 'false'::boolean
WHERE  check_instance_id = $1::uuid;`

const lcStmtBlockedConfigsForDeletedInstance = `
SELECT scicd.blocked_instance_config_id,
       scicd.blocking_instance_config_id,
       scicd.unblocking_state
FROM   soma.check_instances sci
JOIN   soma.check_instance_configurations scic
  ON   sci.check_instance_id = scic.check_instance_id
JOIN   soma.check_instance_configuration_dependencies scicd
  ON   scic.check_instance_config_id = scicd.blocked_instance_config_id
WHERE  sci.deleted
  AND  scic.status = 'blocked';`

const lcStmtConfigAwaitingDeletion = `
UPDATE soma.check_instance_configurations
SET    status = 'awaiting_deletion"::varchar,
       next_status = 'none'::varchar
WHERE  check_instance_config_id = $1::uuid;`

const lcStmtDeleteGhosts = `
UPDATE soma.check_instance_configurations scic
SET    status = 'awaiting_deletion'::varchar,
       next_status = 'none'::varchar
FROM   soma.check_instances sci
WHERE  scic.check_instance_id = sci.check_instance_id
AND    scic.status = 'awaiting_rollout'
AND    sci.deleted
AND    sci.update_available;`

const lcStmtDeleteFailedRollouts = `
UPDATE soma.check_instance_configurations scic
SET    status = 'awaiting_deletion'::varchar,
       next_status = 'none'::varchar
FROM   soma.check_instances sci
WHERE  scic.check_instance_id = sci.check_instance_id
AND    sci.deleted
AND    scic.status = 'rollout_failed';`

const lcStmtDeleteDeprovisioned = `
UPDATE soma.check_instance_configurations scic
SET    status = 'awaiting_deletion'::varchar,
       next_status = 'none'::varchar
FROM   soma.check_instances sci
WHERE  scic.check_instance_id = sci.check_instance_id
AND    sci.deleted
AND    scic.status = 'deprovisioned'
AND    scic.next_status = 'none';`

const lcStmtDeprovisionDeletedActive = `
SELECT scic.check_instance_config_id,
       sci.check_instance_id
FROM   soma.check_instance_configurations scic
JOIN   soma.check_instances sci
  ON   scic.check_instance_id = sci.check_instance_id
WHERE  sci.deleted
  AND  scic.status = 'active'
  AND  scic.next_status = 'none';
`

const lcStmtDeprovisionConfiguration = `
UPDATE soma.check_instance_configurations
SET    status = 'awaiting_deprovision'::varchar,
       next_status = 'deprovision_in_progress'::varchar
WHERE  check_instance_config_id = $1::uuid;
`

const lcStmtDeadLockResolver = `
SELECT ci.check_instance_id,
       ci.current_instance_config_id
FROM   check_instances ci
JOIN   check_instance_configurations cic
  ON   ci.check_instance_id = cic.check_instance_id
 AND   ci.current_instance_config_id = cic.check_instance_config_id
JOIN   check_instance_configuration_dependencies cicd
  ON   ci.current_instance_config_id = cicd.blocking_instance_config_id
WHERE  cic.status = 'active';`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
