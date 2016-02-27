package main

/*
 * Statements for NODE actions
 */

var tkStmtUpdateNodeState = `
UPDATE soma.nodes
SET    object_state = $2::varchar
WHERE  node_id = $1::uuid;`

var tkStmtNodeUnassignFromBucket = `
DELETE FROM soma.node_bucket_assignment
WHERE       node_id = $1::uuid
AND         bucket_id = $2::uuid
AND         organizational_team_id = $3::uuid;`

var tkStmtNodePropertyOncallCreate = `
INSERT INTO soma.node_oncall_properties (
            instance_id,
            source_instance_id,
            node_id,
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

var tkStmtNodePropertyServiceCreate = `
INSERT INTO soma.node_service_properties (
            instance_id,
            source_instance_id,
            node_id,
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

var tkStmtNodePropertySystemCreate = `
INSERT INTO soma.node_system_properties (
            instance_id,
            source_instance_id,
            node_id,
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

var tkStmtNodePropertyCustomCreate = `
INSERT INTO soma.node_custom_properties (
            instance_id,
            source_instance_id,
            node_id,
            view,
            custom_property_id,
            bucket_id,
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
       $7::uuid,
       $8::boolean,
       $9::boolean,
       $10::text;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
