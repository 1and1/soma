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

var tkStmtStartJob = `
UPDATE soma.jobs
SET    job_started = $2::timestamptz,
       job_status = 'in_progress'
WHERE  job_id = $1::uuid
AND    job_started IS NULL;`

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

/*
 * Statements for BUCKET actions
 */

var tkStmtCreateBucket = `
INSERT INTO soma.buckets (
	bucket_id,
	bucket_name,
	bucket_frozen,
	bucket_deleted,
	repository_id,
	environment,
	organizational_team_id)
SELECT	$1::uuid,
        $2::varchar,
        $3::boolean,
        $4::boolean,
        $5::uuid,
        $6::varchar,
        $7::uuid;`

var tkStmtBucketAssignNode = `
INSERT INTO soma.node_bucket_assignment (
            node_id,
            bucket_id,
            organizational_team_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

var tkStmtBucketRemoveNode = `
DELETE FROM soma.node_bucket_assignment (
WHERE       node_id = $1::uuid
AND         bucket_id = $2::uuid
AND         organizational_team_id = $3::uuid;`

var tkStmtBucketPropertyOncallCreate = `
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

var tkStmtBucketPropertyServiceCreate = `
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

var tkStmtBucketPropertySystemCreate = `
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

var tkStmtBucketPropertyCustomCreate = `
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

/*
 * Statements for GROUP actions
 */

var tkStmtGroupCreate = `
INSERT INTO soma.groups (
            group_id,
            bucket_id,
            group_name,
            object_state,
            organizational_team_id)
SELECT $1::uuid,
       $2::uuid,
       $3::varchar,
       $4::varchar,
       $5::uuid;`

var tkStmtGroupUpdate = `
UPDATE soma.groups
SET    object_state = $2::varchar
WHERE  group_id = $1::uuid;`

var tkStmtGroupDelete = `
DELETE FROM soma.groups
WHERE       group_id = $1::uuid;`

var tkStmtGroupMemberNewNode = `
INSERT INTO soma.group_membership_nodes (
            group_id,
            child_node_id,
            bucket_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

var tkStmtGroupMemberNewCluster = `
INSERT INTO soma.group_membership_clusters (
            group_id,
            child_cluster_id,
            bucket_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

var tkStmtGroupMemberNewGroup = `
INSERT INTO soma.group_membership_groups (
            group_id,
            child_group_id,
            bucket_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

var tkStmtGroupMemberRemoveNode = `
DELETE FROM soma.group_membership_nodes
WHERE       group_id = $1::uuid
AND         child_node_id = $2::uuid;`

var tkStmtGroupMemberRemoveCluster = `
DELETE FROM soma.group_membership_clusters
WHERE       group_id = $1::uuid
AND         child_cluster_id = $2::uuid;`

var tkStmtGroupMemberRemoveGroup = `
DELETE FROM soma.group_membership_groups
WHERE       group_id = $1::uuid
AND         child_group_id = $2::uuid;`

var tkStmtGroupPropertyOncallCreate = `
INSERT INTO soma.group_oncall_properties (
            instance_id,
            source_instance_id,
            group_id,
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

var tkStmtGroupPropertyServiceCreate = `
INSERT INTO soma.group_service_properties (
	instance_id,
	source_instance_id,
	group_id,
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

var tkStmtGroupPropertySystemCreate = `
INSERT INTO soma.group_system_properties (
	instance_id,
	source_instance_id,
	group_id,
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

var tkStmtGroupPropertyCustomCreate = `
INSERT INTO soma.group_custom_properties (
	instance_id,
	source_instance_id,
	group_id,
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

/*
 * Statements for CLUSTER actions
 */

var tkStmtClusterCreate = `
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

var tkStmtClusterUpdate = `
UPDATE soma.clusters
SET    object_state = $2::varchar
WHERE  cluster_id = $1::uuid;`

var tkStmtClusterDelete = `
DELETE FROM soma.clusters
WHERE       cluster_id = $1::uuid;`

var tkStmtClusterMemberNew = `
INSERT INTO soma.cluster_membership (
            cluster_id,
            node_id,
            bucket_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

var tkStmtClusterMemberRemove = `
DELETE FROM soma.cluster_membership
WHERE       cluster_id = $1::uuid
AND         node_id = $2::uuid;`

var tkStmtClusterPropertyOncallCreate = `
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

var tkStmtClusterPropertyServiceCreate = `
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

var tkStmtClusterPropertySystemCreate = `
INSERT INTO soma.cluster_system_properties (
	instance_id,
	source_instance_id,
	cluster_id,
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

var tkStmtClusterPropertyCustomCreate = `
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
