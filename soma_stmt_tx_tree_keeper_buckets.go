package main

/*
 * Statements for BUCKET actions
 */

const tkStmtCreateBucket = `
INSERT INTO soma.buckets (
            bucket_id,
            bucket_name,
            bucket_frozen,
            bucket_deleted,
            repository_id,
            environment,
            organizational_team_id)
SELECT $1::uuid,
       $2::varchar,
       $3::boolean,
       $4::boolean,
       $5::uuid,
       $6::varchar,
       $7::uuid;`

const tkStmtBucketAssignNode = `
INSERT INTO soma.node_bucket_assignment (
            node_id,
            bucket_id,
            organizational_team_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

const tkStmtBucketRemoveNode = `
DELETE FROM soma.node_bucket_assignment (
WHERE       node_id = $1::uuid
AND         bucket_id = $2::uuid
AND         organizational_team_id = $3::uuid;`

const tkStmtBucketPropertyOncallCreate = `
INSERT INTO soma.bucket_oncall_properties (
            instance_id,
            source_instance_id,
            bucket_id,
            view,
            oncall_duty_id,
            repository_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::uuid,
       $7::boolean,
       $8::boolean;`

const tkStmtBucketPropertyServiceCreate = `
INSERT INTO soma.bucket_service_properties (
            instance_id,
            source_instance_id,
            bucket_id,
            view,
            service_property,
            organizational_team_id,
            repository_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::uuid,
       $7::uuid,
       $8::boolean,
       $9::boolean;`

const tkStmtBucketPropertySystemCreate = `
INSERT INTO soma.bucket_system_properties (
            instance_id,
            source_instance_id,
            bucket_id,
            view,
            system_property,
            object_type,
            repository_id,
            inheritance_enabled,
            children_only,
            value,
            inherited)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::varchar,
       $7::uuid,
       $8::boolean,
       $9::boolean,
       $10::text,
       $11::boolean;`

const tkStmtBucketPropertyCustomCreate = `
INSERT INTO soma.bucket_custom_properties (
            instance_id,
            source_instance_id,
            bucket_id,
            view,
            custom_property_id,
            repository_id,
            inheritance_enabled,
            children_only,
            value)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::uuid,
       $7::boolean,
       $8::boolean,
       $9::text;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
