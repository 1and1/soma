package main

/*
 * Statements for job state updates outside transaction
 */

var tkStmtStartJob = `
UPDATE soma.jobs
SET    job_started = $2::timestamptz,
       job_status = 'in_progress'
WHERE  job_id = $1::uuid
AND    job_started IS NULL;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
