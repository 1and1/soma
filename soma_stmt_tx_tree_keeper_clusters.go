package main

/*
 * Statements for CLUSTER actions
 */

const tkStmtClusterCreate = `
INSERT INTO soma.clusters (
            cluster_id,
            cluster_name,
            bucket_id,
            object_state,
            organizational_team_id)
SELECT $1::uuid,
       $2::varchar,
       $3::uuid,
       $4::varchar,
       $5::uuid;`

const tkStmtClusterUpdate = `
UPDATE soma.clusters
SET    object_state = $2::varchar
WHERE  cluster_id = $1::uuid;`

const tkStmtClusterDelete = `
DELETE FROM soma.clusters
WHERE       cluster_id = $1::uuid;`

const tkStmtClusterMemberNew = `
INSERT INTO soma.cluster_membership (
            cluster_id,
            node_id,
            bucket_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

const tkStmtClusterMemberRemove = `
DELETE FROM soma.cluster_membership
WHERE       cluster_id = $1::uuid
AND         node_id = $2::uuid;`

const tkStmtClusterPropertyOncallCreate = `
INSERT INTO soma.cluster_oncall_properties (
            instance_id,
            source_instance_id,
            cluster_id,
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

const tkStmtClusterPropertyOncallDelete = `
DELETE FROM soma.cluster_oncall_properties
WHERE       instance_id = $1::uuid;`

const tkStmtClusterPropertyServiceCreate = `
INSERT INTO soma.cluster_service_properties (
            instance_id,
            source_instance_id,
            cluster_id,
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

const tkStmtClusterPropertyServiceDelete = `
DELETE FROM soma.cluster_service_properties
WHERE       instance_id = $1::uuid;`

const tkStmtClusterPropertySystemCreate = `
INSERT INTO soma.cluster_system_properties (
            instance_id,
            source_instance_id,
            cluster_id,
            view,
            system_property,
            source_type,
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

const tkStmtClusterPropertySystemDelete = `
DELETE FROM soma.cluster_system_properties
WHERE       instance_id = $1::uuid;`

const tkStmtClusterPropertyCustomCreate = `
INSERT INTO soma.cluster_custom_properties (
            instance_id,
            source_instance_id,
            cluster_id,
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

const tkStmtClusterPropertyCustomDelete = `
DELETE FROM soma.cluster_custom_properties
WHERE       instance_id = $1::uuid;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
