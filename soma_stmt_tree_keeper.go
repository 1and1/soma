package main

/*
 * Statements for job state updates outside transaction
 */

const tkStmtStartJob = `
UPDATE soma.jobs
SET    job_started = $2::timestamptz,
       job_status = 'in_progress'
WHERE  job_id = $1::uuid
AND    job_started IS NULL;`

const tkStmtGetViewFromCapability = `
SELECT capability_view
FROM   soma.monitoring_capabilities
WHERE  capability_id = $1::uuid;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
