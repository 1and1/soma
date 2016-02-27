package main

/*
 * Statements for REPOSITORY actions
 */

var tkStmtRepositoryPropertyOncallCreate = `
INSERT INTO soma.repository_oncall_properties (
            instance_id,
            source_instance_id,
            repository_id,
            view,
            oncall_duty_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::boolean,
       $7::boolean;`

var tkStmtRepositoryPropertyServiceCreate = `
INSERT INTO soma.repository_service_properties (
            instance_id,
            source_instance_id,
            repository_id,
            view,
            service_property,
            organizational_team_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::uuid,
       $7::boolean,
       $8::boolean;`

var tkStmtRepositoryPropertySystemCreate = `
INSERT INTO soma.repository_system_properties (
            instance_id,
            source_instance_id,
            repository_id,
            view,
            system_property,
            object_type,
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
       $7::boolean,
       $8::boolean,
       $9::text,
       $10::boolean;`

var tkStmtRepositoryPropertyCustomCreate = `
INSERT INTO soma.repository_custom_properties (
            instance_id,
            source_instance_id,
            repository_id,
            view,
            custom_property_id,
            inheritance_enabled,
            children_only,
            value)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::boolean,
       $7::boolean,
       $8::text;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
