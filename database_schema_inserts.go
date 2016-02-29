package main

func schemaInserts(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 10)

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
    ( 'processed' )
;`
	queries[idx] = "insertJobStatus"
	idx++

	queryMap["insertJobResults"] = `
INSERT INTO soma.job_results (
    job_result )
VALUES
    ( 'pending' ),
    ( 'success' ),
    ( 'failed' )
;`
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

	/*
			queryMap["insertSystemPropertyValidity"] = `
		INSERT INTO soma.system_property_validity (
		    system_property,
		    object_type,
		    inherited )
		VALUES
		    ( 'disable_ip4',        'repository', 'no'  ),
		    ( 'disable_ip4',        'bucket',     'no'  ),
		    ( 'disable_ip4',        'group',      'no'  ),
		    ( 'disable_ip4',        'cluster',    'no'  ),
		    ( 'disable_ip4',        'node',       'no'  ),
		    ( 'disable_ip4',        'repository', 'yes' ),
		    ( 'disable_ip4',        'bucket',     'yes' ),
		    ( 'disable_ip4',        'group',      'yes' ),
		    ( 'disable_ip4',        'cluster',    'yes' ),
		    ( 'disable_ip4',        'node',       'yes' ),
		    ( 'disable_ip6',        'repository', 'no'  ),
		    ( 'disable_ip6',        'bucket',     'no'  ),
		    ( 'disable_ip6',        'group',      'no'  ),
		    ( 'disable_ip6',        'cluster',    'no'  ),
		    ( 'disable_ip6',        'node',       'no'  ),
		    ( 'disable_ip6',        'repository', 'yes' ),
		    ( 'disable_ip6',        'bucket',     'yes' ),
		    ( 'disable_ip6',        'group',      'yes' ),
		    ( 'disable_ip6',        'cluster',    'yes' ),
		    ( 'disable_ip6',        'node',       'yes' ),
		    ( 'dns_zone',           'repository', 'yes' ),
		    ( 'dns_zone',           'bucket',     'yes' ),
		    ( 'dns_zone',           'group',      'yes' ),
		    ( 'dns_zone',           'cluster',    'yes' ),
		    ( 'dns_zone',           'node',       'yes' ),
		    ( 'dns_zone',           'repository', 'no'  ),
		    ( 'dns_zone',           'bucket',     'no'  ),
		    ( 'dns_zone',           'group',      'no'  ),
		    ( 'dns_zone',           'cluster',    'no'  ),
		    ( 'dns_zone',           'node',       'no'  ),
		    ( 'fqdn',               'group',      'yes' ),
		    ( 'fqdn',               'cluster',    'yes' ),
		    ( 'fqdn',               'node',       'yes' ),
		    ( 'fqdn',               'group',      'no'  ),
		    ( 'fqdn',               'cluster',    'no'  ),
		    ( 'fqdn',               'node',       'no'  ),
		    ( 'cluster_state',      'cluster',    'no'  ),
		    ( 'cluster_state',      'node',       'yes' ),
		    ( 'cluster_ha_address', 'cluster',    'no'  ),
		    ( 'cluster_ha_address', 'node',       'yes' ),
		    ( 'cluster_datacenter', 'cluster',    'no'  ),
		    ( 'cluster_datacenter', 'node',       'yes' ),
		    ( 'group_ha_address',   'group',      'no'  ),
		    ( 'group_ha_address',   'group',      'yes' ),
		    ( 'group_ha_address',   'cluster',    'yes' ),
		    ( 'group_ha_address',   'node',       'yes' ),
		    ( 'group_datacenter',   'group',      'no'  ),
		    ( 'group_datacenter',   'group',      'yes' ),
		    ( 'group_datacenter',   'cluster',    'yes' ),
		    ( 'group_datacenter',   'node',       'yes' ),
		    ( 'information_system', 'repository', 'yes' ),
		    ( 'information_system', 'bucket',     'yes' ),
		    ( 'information_system', 'group',      'yes' ),
		    ( 'information_system', 'cluster',    'yes' ),
		    ( 'information_system', 'node',       'yes' ),
		    ( 'information_system', 'repository', 'no'  ),
		    ( 'information_system', 'bucket',     'no'  ),
		    ( 'information_system', 'group',      'no'  ),
		    ( 'information_system', 'cluster',    'no'  ),
		    ( 'information_system', 'node',       'no'  ),
		    ( 'yp_asset',           'repository', 'yes' ),
		    ( 'yp_asset',           'bucket',     'yes' ),
		    ( 'yp_asset',           'group',      'yes' ),
		    ( 'yp_asset',           'cluster',    'yes' ),
		    ( 'yp_asset',           'node',       'yes' ),
		    ( 'yp_asset',           'repository', 'no'  ),
		    ( 'yp_asset',           'bucket',     'no'  ),
		    ( 'yp_asset',           'group',      'no'  ),
		    ( 'yp_asset',           'cluster',    'no'  ),
		    ( 'yp_asset',           'node',       'no'  ),
		    ( 'frozen',             'repository', 'yes' ),
		    ( 'frozen',             'bucket',     'yes' ),
		    ( 'frozen',             'group',      'yes' ),
		    ( 'frozen',             'cluster',    'yes' ),
		    ( 'frozen',             'node',       'yes' ),
		    ( 'frozen',             'repository', 'no'  ),
		    ( 'frozen',             'bucket',     'no'  ),
		    ( 'frozen',             'group',      'no'  ),
		    ( 'frozen',             'cluster',    'no'  ),
		    ( 'frozen',             'node',       'no'  ),
		    ( 'tag',                'repository', 'yes' ),
		    ( 'tag',                'bucket',     'yes' ),
		    ( 'tag',                'group',      'yes' ),
		    ( 'tag',                'cluster',    'yes' ),
		    ( 'tag',                'node',       'yes' ),
		    ( 'tag',                'repository', 'no'  ),
		    ( 'tag',                'bucket',     'no'  ),
		    ( 'tag',                'group',      'no'  ),
		    ( 'tag',                'cluster',    'no'  ),
		    ( 'tag',                'node',       'no'  ),
		    ( 'documentation',      'repository', 'yes' ),
		    ( 'documentation',      'bucket',     'yes' ),
		    ( 'documentation',      'group',      'yes' ),
		    ( 'documentation',      'cluster',    'yes' ),
		    ( 'documentation',      'node',       'yes' ),
		    ( 'documentation',      'repository', 'no'  ),
		    ( 'documentation',      'bucket',     'no'  ),
		    ( 'documentation',      'group',      'no'  ),
		    ( 'documentation',      'cluster',    'no'  ),
		    ( 'documentation',      'node',       'no'  ),
		    ( 'link',               'repository', 'yes' ),
		    ( 'link',               'bucket',     'yes' ),
		    ( 'link',               'group',      'yes' ),
		    ( 'link',               'cluster',    'yes' ),
		    ( 'link',               'node',       'yes' ),
		    ( 'link',               'repository', 'no'  ),
		    ( 'link',               'bucket',     'no'  ),
		    ( 'link',               'group',      'no'  ),
		    ( 'link',               'cluster',    'no'  ),
		    ( 'link',               'node',       'no'  ),
		    ( 'user_management',    'repository', 'yes' ),
		    ( 'user_management',    'bucket',     'yes' ),
		    ( 'user_management',    'group',      'yes' ),
		    ( 'user_management',    'cluster',    'yes' ),
		    ( 'user_management',    'node',       'yes' ),
		    ( 'user_management',    'repository', 'no'  ),
		    ( 'user_management',    'bucket',     'no'  ),
		    ( 'user_management',    'group',      'no'  ),
		    ( 'user_management',    'cluster',    'no'  ),
		    ( 'user_management',    'node',       'no'  ),
		    ( 'wiki',               'repository', 'yes' ),
		    ( 'wiki',               'bucket',     'yes' ),
		    ( 'wiki',               'group',      'yes' ),
		    ( 'wiki',               'cluster',    'yes' ),
		    ( 'wiki',               'node',       'yes' ),
		    ( 'wiki',               'repository', 'no'  ),
		    ( 'wiki',               'bucket',     'no'  ),
		    ( 'wiki',               'group',      'no'  ),
		    ( 'wiki',               'cluster',    'no'  ),
		    ( 'wiki',               'node',       'no'  )
		;`
			queries[idx] = "insertSystemPropertyValidity"
			idx++
	*/

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
