/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const ListAllOutstandingJobs = `
SELECT job_id,
       job_type
FROM   soma.jobs
WHERE  job_status != 'processed';`

const ListScopedOutstandingJobs = `
SELECT sj.job_id,
       sj.job_type
FROM   inventory.users iu
JOIN   soma.jobs sj
  ON   iu.user_id = sj.user_id
WHERE  iu.user_uid = $1::varchar
UNION
SELECT sj.job_id,
       sj.job_type
FROM   inventory.users iu
JOIN   soma.jobs sj
  ON   iu.organizational_team_id = sj.organizational_team_id
WHERE  iu.user_uid = $1::varchar
  AND  sj.user_id NOT IN
  (    SELECT user_id FROM inventory.users
       WHERE user_uid = $1::varchar);`

const JobResultForId = `
SELECT job_id,
       job_status,
       job_result,
       job_type,
       job_serial,
       repository_id,
       user_id,
       organizational_team_id,
       job_queued,
       job_started,
       job_finished,
       job_error,
       job
FROM   soma.jobs
WHERE  job_id = $1::uuid;`

const JobResultsForList = `
SELECT job_id,
       job_status,
       job_result,
       job_type,
       job_serial,
       repository_id,
       user_id,
       organizational_team_id,
       job_queued,
       job_started,
       job_finished,
       job_error,
       job
FROM   soma.jobs
WHERE  job_id = any($1::uuid[]);`

const JobSave = `
INSERT INTO soma.jobs (
            job_id,
            job_status,
            job_result,
            job_type,
            repository_id,
            user_id,
            organizational_team_id,
            job)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar,
       $4::varchar,
       $5::uuid,
       iu.user_id,
       iu.organizational_team_id,
       $7::jsonb
FROM   inventory.users iu
WHERE  iu.user_uid = $6::varchar;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
