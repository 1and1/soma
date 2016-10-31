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
	InstanceStatements = ``

	// InstanceScopedList can return all instances (both parameters
	// null), all instances of a specific repository (first parameter
	// not null) or a specific bucket (second parameter not null).
	// If both parameters are specified with an invalid repositoryId+
	// bucketId combination, there resultset is empty.
	// Result columns are sufficient to fill proto.Instance.
	InstanceScopedList = `
SELECT sci.check_instance_id,
       scic.version,
       sc.check_id,
       sc.configuration_id,
       scic.current_instance_config_id,
       sc.repository_id,
       sc.bucket_id,
       sc.object_id,
       sc.object_type,
       scic.status,
       scic.next_status,
       (sc.object_id = sc.source_object_id)::boolean
FROM   soma.checks sc
JOIN   soma.check_instances sci
  ON   sc.check_id = sci.check_id
JOIN   soma.check_instance_configurations scic
  ON   sci.current_instance_config_id = scic.check_instance_config_id
WHERE  (sc.repository_id = $1::uuid OR $1::uuid IS NULL)
  AND  (sc.bucket_id = $2::uuid OR $2::uuid IS NULL)
  AND  NOT sc.deleted
  AND  NOT sci.deleted;`

	// InstanceShow returns information about a single check instance.
	// Result columns are sufficient to fill proto.Instance.
	InstanceShow = `
SELECT sci.check_instance_id,
       scic.version,
       sc.check_id,
       sc.configuration_id,
       scic.current_instance_config_id,
       sc.repository_id,
       sc.bucket_id,
       sc.object_id,
       sc.object_type,
       scic.status,
       scic.next_status,
       (sc.object_id = sc.source_object_id)::boolean
FROM   soma.checks sc
JOIN   soma.check_instances sci
  ON   sc.check_id = sci.check_id
JOIN   soma.check_instance_configurations scic
  ON   sci.current_instance_config_id = scic.check_instance_config_id
WHERE  sci.check_instance_id = $1::uuid
  AND  NOT sc.deleted
  AND  NOT sci.deleted;`
)

func init() {
	m[InstanceScopedList] = `InstanceScopedList`
	m[InstanceShow] = `InstanceShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
