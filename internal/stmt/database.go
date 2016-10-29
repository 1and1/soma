/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const (
	DatabaseStatements = ``

	DatabaseTimezone = `SET TIME ZONE 'UTC';`

	DatabaseIsolationLevel = `SET SESSION CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL SERIALIZABLE;`

	DatabaseSchemaVersion = `
SELECT schema,
       MAX(version) AS version
FROM   public.schema_versions
GROUP  BY schema;`

	ReadOnlyTransaction = `SET TRANSACTION READ ONLY, DEFERRABLE;`
)

func init() {
	m[DatabaseTimezone] = `DatabaseTimezone`
	m[DatabaseIsolationLevel] = `DatabaseIsolationLevel`
	m[DatabaseSchemaVersion] = `DatabaseSchemaVersion`
	m[ReadOnlyTransaction] = `ReadOnlyTransaction`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
