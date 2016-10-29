/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const (
	PredicateStatements = ``

	PredicateList = `
SELECT predicate
FROM   soma.configuration_predicates;`

	PredicateShow = `
SELECT predicate
FROM   soma.configuration_predicates
WHERE  predicate = $1;`

	PredicateAdd = `
INSERT INTO soma.configuration_predicates (
            predicate)
SELECT $1::varchar
WHERE  NOT EXISTS (
   SELECT predicate
   FROM   soma.configuration_predicates
   WHERE  predicate = $1::varchar);`

	PredicateDel = `
DELETE FROM soma.configuration_predicates
WHERE       predicate = $1::varchar;`
)

func init() {
	m[PredicateAdd] = `PredicateAdd`
	m[PredicateDel] = `PredicateDel`
	m[PredicateList] = `PredicateList`
	m[PredicateShow] = `PredicateShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
