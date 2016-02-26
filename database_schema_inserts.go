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

	queryMap["insertJobResults"] = `
INSERT INTO soma.job_results (
  job_result )
VALUES
  ( 'pending' ),
  ( 'success' ),
  ( 'failed' );`
	queries[idx] = "insertJobResults"
	idx++

	queryMap["insertJobTypes"] = `
INSERT INTO soma.job_types (
  job_type )
VALUES
  ( 'create_bucket' ),
  ( 'create_group' ),
  ( 'create_cluster' ),
  ( 'assign_node' ),
  ( 'add_group_to_group' ),
  ( 'add_cluster_to_group' ),
  ( 'add_node_to_group' ),
  ( 'add_node_to_cluster' ),
  ( 'add_system_property_to_repository' ),
  ( 'add_custom_property_to_repository' ),
  ( 'add_oncall_property_to_repository' ),
  ( 'add_service_property_to_repository' ),
  ( 'add_system_property_to_bucket' ),
  ( 'add_custom_property_to_bucket' ),
  ( 'add_oncall_property_to_bucket' ),
  ( 'add_service_property_to_bucket' ),
  ( 'add_system_property_to_group' ),
  ( 'add_custom_property_to_group' ),
  ( 'add_oncall_property_to_group' ),
  ( 'add_service_property_to_group' ),
  ( 'add_system_property_to_cluster' ),
  ( 'add_custom_property_to_cluster' ),
  ( 'add_oncall_property_to_cluster' ),
  ( 'add_service_property_to_cluster' ),
  ( 'add_system_property_to_node' ),
  ( 'add_custom_property_to_node' ),
  ( 'add_oncall_property_to_node' ),
  ( 'add_service_property_to_node' )
  ;`
	queries[idx] = "insertJobTypes"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
