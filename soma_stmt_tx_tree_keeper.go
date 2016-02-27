package main

var tkStmtPropertyInstanceCreate = `
INSERT INTO soma.property_instances (
            instance_id,
            repository_id,
            source_instance_id,
            source_object_type,
            source_object_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid;`

/*
 * Statements for job state updates
 */

var tkStmtFinishJob = `
UPDATE soma.jobs
SET    job_finished = $2::timestamptz,
       job_status = 'processed',
       job_result = $3::varchar
WHERE  job_id = $1::uuid;`

/*
 * Referential integrity hacking
 */

var tkStmtDeferAllConstraints = `
SET CONSTRAINTS ALL DEFERRED;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
