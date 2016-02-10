package main

func schemaInserts(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 7)

	queryMap["insertSystemGroupWheel"] = `
INSERT INTO inventory.organizational_teams (
  organizational_team_id,
  organizational_team_name,
  organizational_team_ldap_id,
  organizational_team_system )
VALUES (
  '00000000-0000-0000-0000-000000000000',
  'wheel',
  0,
  'yes'
);`
	queries[idx] = "insertSystemGroupWheel"
	idx++

	queryMap["insertSystemUserRoot"] = `
INSERT INTO inventory.users (
  user_id,
  user_uid,
  user_first_name,
  user_last_name,
  user_employee_number,
  user_mail_address,
  user_is_active,
  user_is_system,
  user_is_deleted,
  organizational_team_id )
VALUES (
  '00000000-0000-0000-0000-000000000000',
  'root',
  'Charlie',
  'Root',
  0,
  'monitoring@1und1.de',
  'yes',
  'yes',
  'no',
  '00000000-0000-0000-0000-000000000000'
);`
	queries[idx] = "insertSystemUserRoot"
	idx++

	queryMap["insertJobStatus"] = `
INSERT INTO soma.job_status (
  job_status )
VALUES
  ( 'queued' ),
  ( 'in_progress' ),
  ( 'processed' );`
	queries[idx] = "insertJobStatus"
	idx++

	queryMap["insertJobResult"] = `
INSERT INTO soma.job_result (
  job_result )
VALUES
  ( 'pending' ),
  ( 'success' ),
  ( 'failed' );`
	queries[idx] = "insertJobResult"
	idx++

	queryMap["insertJobType"] = `
INSERT INTO soma.job_type (
  job_type )
VALUES
  ( 'create_bucket' ),
  ( 'create_group' ),
  ( 'create_cluster' ),
  ( 'assign_node' );`
	queries[idx] = "insertJobType"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
