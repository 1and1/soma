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
)

func init() {
	m[WorkflowSummary] = `WorkflowSummary`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
