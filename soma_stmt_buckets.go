package main

const stmtBucketList = `
SELECT bucket_id,
       bucket_name
FROM   soma.buckets;`

const stmtBucketShow = `
SELECT bucket_id,
       bucket_name,
       bucket_frozen,
       bucket_deleted,
       repository_id,
       environment,
       organizational_team_id
FROM   soma.buckets
WHERE  bucket_id = $1::uuid;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
