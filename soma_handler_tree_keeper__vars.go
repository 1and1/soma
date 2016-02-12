package main

var tkStmtCreateGroup = `
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

var tkStmtCreateCluster = `
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

var tkStmtDeferConstraints = `SET CONSTRAINTS ALL DEFERRED;`

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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
