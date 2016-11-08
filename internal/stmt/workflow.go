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
	WorkflowStatements = ``

	// WorkflowSummary returns a summary of the current workflow
	// status distribution in the system
	WorkflowSummary = `
SELECT scic.status,
       count(1)
FROM   soma.check_instances sci
JOIN   soma.check_instance_configurations scic
  ON   sci.check_instance_id = scic.check_instance_id
WHERE  NOT sci.deleted
GROUP  BY scic.status;`

	// WorkflowList returns information about check instance
	// configurations in a specific workflow state
	WorkflowList = `
SELECT sci.check_instance_id,
       sc.check_id,
       sc.repository_id,
       sc.configuration_id,
       scic.check_instance_config_id,
       scic.version,
       scic.status,
       scic.created,
       scic.activated_at,
       scic.deprovisioned_at,
       scic.status_last_updated_at,
       scic.notified_at,
	   (sc.object_id = sc.source_object_id)::boolean
FROM   soma.check_instances sci
JOIN   soma.checks sc
  ON   sci.check_id = sc.check_id
JOIN   soma.check_instance_configurations scic
  ON   sci.check_instance_id = scic.check_instance_id
WHERE  NOT sci.deleted
  AND  scic.status = $1::varchar;`

	// Reset rollout status for retry of check instance
	// configuration rollout
	WorkflowRetry = `
UPDATE soma.check_instance_configurations scic
SET    status = (CASE status
                 WHEN 'rollout_failed'     THEN 'awaiting_rollout'
                 WHEN 'deprovision_failed' THEN 'awaiting_deprovision'
                 END),
       next_status = (CASE status
                      WHEN 'rollout_failed'     THEN 'rollout_in_progress'
                      WHEN 'deprovision_failed' THEN 'deprovision_in_progress'
                      END)
FROM   soma.check_instances sci
WHERE  sci.current_instance_config_id = scic.check_instance_config_id
  AND  scic.status IN ( 'rollout_failed', 'deprovision_failed' )
  AND  sci.check_instance_id = $1::uuid;`

	// Set update available flag for check instance
	WorkflowUpdateAvailable = `
UPDATE soma.check_instances
SET    update_available = 'true'::boolean
WHERE  check_instance_id = $1::uuid;`
)

func init() {
	m[WorkflowList] = `WorkflowList`
	m[WorkflowRetry] = `WorkflowRetry`
	m[WorkflowSummary] = `WorkflowSummary`
	m[WorkflowUpdateAvailable] = `WorkflowUpdateAvailable`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
