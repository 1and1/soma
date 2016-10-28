/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const MetricVerify = `
SELECT metric
FROM   soma.metrics
WHERE  metric = $1::varchar;`

const MetricList = `
SELECT metric
FROM   soma.metrics;`

const MetricShow = `
SELECT metric,
       metric_unit,
       description
FROM   soma.metrics
WHERE  metric = $1::varchar;`

const MetricAdd = `
INSERT INTO soma.metrics (
            metric,
            metric_unit,
            description)
SELECT   $1::varchar, $2::varchar, $3::text
WHERE    NOT EXISTS (
   SELECT metric
   FROM   soma.metrics
   WHERE  metric = $1::varchar);`

const MetricDel = `
DELETE FROM soma.metrics
WHERE       metric = $1::varchar;`

const MetricPkgAdd = `
INSERT INTO soma.metric_packages (
            metric,
            metric_provider,
            package)
SELECT   $1::varchar, $2::varchar, $3::varchar
WHERE    NOT EXISTS (
   SELECT metric
   FROM   soma.metric_packages
   WHERE  metric = $1::varchar
   AND    metric_provider = $2::varchar);`

const MetricPkgDel = `
DELETE FROM soma.metric_packages
WHERE       metric = $1::varchar;`

func init() {
	m[MetricAdd] = `MetricAdd`
	m[MetricDel] = `MetricDel`
	m[MetricList] = `MetricList`
	m[MetricPkgAdd] = `MetricPkgAdd`
	m[MetricPkgDel] = `MetricPkgDel`
	m[MetricShow] = `MetricShow`
	m[MetricVerify] = `MetricVerify`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
