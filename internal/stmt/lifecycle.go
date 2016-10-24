package stmt

const LifecycleActiveUnblockCondition = `
SELECT  scicd.blocked_instance_config_id,
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

const LifecycleUpdateInstance = `
UPDATE  soma.check_instances
SET     update_available = $1::boolean,
        current_instance_config_id = $2::uuid
WHERE   check_instance_id = $3::uuid;`

const LifecycleUpdateConfig = `
UPDATE  soma.check_instance_configurations
SET     status = $1::varchar,
        next_status = $2::varchar,
        awaiting_deletion = $3::boolean
WHERE   check_instance_config_id = $4::uuid;`

const LifecycleDeleteDependency = `
DELETE FROM soma.check_instance_configuration_dependencies
WHERE       blocked_instance_config_id = $1::uuid
AND         blocking_instance_config_id = $2::uuid
AND         unblocking_state = $3::varchar;`

const LifecycleReadyDeployments = `
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

const LifecycleClearUpdateFlag = `
UPDATE soma.check_instances
SET    update_available = 'false'::boolean
WHERE  check_instance_id = $1::uuid;`

const LifecycleBlockedConfigsForDeletedInstance = `
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

const LifecycleConfigAwaitingDeletion = `
UPDATE soma.check_instance_configurations
SET    status = 'awaiting_deletion'::varchar,
       next_status = 'none'::varchar,
       awaiting_deletion = 'yes'::boolean
WHERE  check_instance_config_id = $1::uuid;`

const LifecycleDeleteGhosts = `
UPDATE soma.check_instance_configurations scic
SET    status = 'awaiting_deletion'::varchar,
       next_status = 'none'::varchar,
       awaiting_deletion = 'yes'::boolean
FROM   soma.check_instances sci
WHERE  scic.check_instance_id = sci.check_instance_id
AND    scic.status = 'awaiting_rollout'
AND    sci.deleted
AND    sci.update_available;`

const LifecycleDeleteFailedRollouts = `
UPDATE soma.check_instance_configurations scic
SET    status = 'awaiting_deletion'::varchar,
       next_status = 'none'::varchar,
       awaiting_deletion = 'yes'::boolean
FROM   soma.check_instances sci
WHERE  scic.check_instance_id = sci.check_instance_id
AND    sci.deleted
AND    scic.status = 'rollout_failed';`

const LifecycleDeleteDeprovisioned = `
UPDATE soma.check_instance_configurations scic
SET    status = 'awaiting_deletion'::varchar,
       next_status = 'none'::varchar,
       awaiting_deletion = 'yes'::boolean
FROM   soma.check_instances sci
WHERE  scic.check_instance_id = sci.check_instance_id
AND    sci.deleted
AND    scic.status = 'deprovisioned'
AND    scic.next_status = 'none';`

const LifecycleDeprovisionDeletedActive = `
SELECT scic.check_instance_config_id,
       sci.check_instance_id
FROM   soma.check_instance_configurations scic
JOIN   soma.check_instances sci
  ON   scic.check_instance_id = sci.check_instance_id
WHERE  sci.deleted
  AND  scic.status = 'active'
  AND  scic.next_status = 'none';
`

const LifecycleDeprovisionConfiguration = `
UPDATE soma.check_instance_configurations
SET    status = 'awaiting_deprovision'::varchar,
       next_status = 'deprovision_in_progress'::varchar
WHERE  check_instance_config_id = $1::uuid;
`

const LifecycleDeadLockResolver = `
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
