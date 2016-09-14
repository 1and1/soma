package main

import "fmt"

func schemaInserts(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 100)

	queryMap["insertSystemGroupWheel"] = `
INSERT INTO inventory.organizational_teams (
            organizational_team_id,
            organizational_team_name,
            organizational_team_ldap_id,
            organizational_team_system
) VALUES (
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
            organizational_team_id
) VALUES (
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

	queryMap["insertCategoryOmnipotence"] = `
INSERT INTO soma.permission_types (
            permission_type,
            created_by
) VALUES (
            'omnipotence',
            '00000000-0000-0000-0000-000000000000'
);`
	queries[idx] = "insertCategoryOmnipotence"
	idx++

	queryMap["insertPermissionOmnipotence"] = `
INSERT INTO soma.permissions (
            permission_id,
            permission_name,
            permission_type,
            created_by
) VALUES (
            '00000000-0000-0000-0000-000000000000',
            'omnipotence',
            'omnipotence',
            '00000000-0000-0000-0000-000000000000'
);`
	queries[idx] = "insertPermissionOmnipotence"
	idx++

	queryMap["grantOmnipotence"] = `
INSERT INTO soma.authorizations_global (
            grant_id,
            user_id,
            permission_id,
            permission_type,
            created_by
) VALUES (
            '00000000-0000-0000-0000-000000000000',
            '00000000-0000-0000-0000-000000000000',
            '00000000-0000-0000-0000-000000000000',
            'omnipotence',
            '00000000-0000-0000-0000-000000000000'
);`
	queries[idx] = "grantOmnipotence"
	idx++

	queryMap["insertJobStatus"] = `
INSERT INTO soma.job_status (
            job_status
) VALUES
            ( 'queued' ),
            ( 'in_progress' ),
            ( 'processed' )
;`
	queries[idx] = "insertJobStatus"
	idx++

	queryMap["insertJobResults"] = `
INSERT INTO soma.job_results (
            job_result
) VALUES
            ( 'pending' ),
            ( 'success' ),
            ( 'failed' )
;`
	queries[idx] = "insertJobResults"
	idx++

	queryMap["insertJobTypes"] = `
INSERT INTO soma.job_types (
            job_type
) VALUES
            ( 'add_check_to_bucket' ),
            ( 'add_check_to_cluster' ),
            ( 'add_check_to_group' ),
            ( 'add_check_to_node' ),
            ( 'add_check_to_repository' ),
            ( 'add_cluster_to_group' ),
            ( 'add_custom_property_to_bucket' ),
            ( 'add_custom_property_to_cluster' ),
            ( 'add_custom_property_to_group' ),
            ( 'add_custom_property_to_node' ),
            ( 'add_custom_property_to_repository' ),
            ( 'add_group_to_group' ),
            ( 'add_node_to_cluster' ),
            ( 'add_node_to_group' ),
            ( 'add_oncall_property_to_bucket' ),
            ( 'add_oncall_property_to_cluster' ),
            ( 'add_oncall_property_to_group' ),
            ( 'add_oncall_property_to_node' ),
            ( 'add_oncall_property_to_repository' ),
            ( 'add_service_property_to_bucket' ),
            ( 'add_service_property_to_cluster' ),
            ( 'add_service_property_to_group' ),
            ( 'add_service_property_to_node' ),
            ( 'add_service_property_to_repository' ),
            ( 'add_system_property_to_bucket' ),
            ( 'add_system_property_to_cluster' ),
            ( 'add_system_property_to_group' ),
            ( 'add_system_property_to_node' ),
            ( 'add_system_property_to_repository' ),
            ( 'assign_node' ),
            ( 'create_bucket' ),
            ( 'create_cluster' ),
            ( 'create_group' ),
            ( 'delete_custom_property_from_bucket' ),
            ( 'delete_custom_property_from_cluster' ),
            ( 'delete_custom_property_from_group' ),
            ( 'delete_custom_property_from_node' ),
            ( 'delete_custom_property_from_repository' ),
            ( 'delete_oncall_property_from_bucket' ),
            ( 'delete_oncall_property_from_cluster' ),
            ( 'delete_oncall_property_from_group' ),
            ( 'delete_oncall_property_from_node' ),
            ( 'delete_oncall_property_from_repository' ),
            ( 'delete_service_property_from_bucket' ),
            ( 'delete_service_property_from_cluster' ),
            ( 'delete_service_property_from_group' ),
            ( 'delete_service_property_from_node' ),
            ( 'delete_service_property_from_repository' ),
            ( 'delete_system_property_from_bucket' ),
            ( 'delete_system_property_from_cluster' ),
            ( 'delete_system_property_from_group' ),
            ( 'delete_system_property_from_node' ),
            ( 'delete_system_property_from_repository' ),
            ( 'remove_check_from_bucket' ),
            ( 'remove_check_from_cluster' ),
            ( 'remove_check_from_group' ),
            ( 'remove_check_from_node' ),
            ( 'remove_check_from_repository' )
;`
	queries[idx] = "insertJobTypes"
	idx++

	queryMap["insertRootRestricted"] = `
INSERT INTO root.flags (
            flag,
            status
) VALUES
            ( 'restricted', false ),
            ( 'disabled', false )
;`
	queries[idx] = "insertRootRestricted"
	idx++

	performDatabaseTask(printOnly, verbose, queries[:idx], queryMap)
}

func schemaVersionInserts(printOnly bool, verbose bool, version string) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 100)

	invString := fmt.Sprintf(`
INSERT INTO public.schema_versions (
            schema,
            version,
            description )
VALUES (
            'inventory',
            201605060001,
            'Initial create - somadbctl %s'
);`, version)
	queryMap["insertInventorySchemaVersion"] = invString
	queries[idx] = "insertInventorySchemaVersion"
	idx++

	authString := fmt.Sprintf(`
INSERT INTO public.schema_versions (
            schema,
            version,
            description
) VALUES (
            'auth',
            201605190001,
            'Initial create - somadbctl %s'
);`, version)
	queryMap["insertAuthSchemaVersion"] = authString
	queries[idx] = "insertAuthSchemaVersion"
	idx++

	somaString := fmt.Sprintf(`
INSERT INTO public.schema_versions (
            schema,
            version,
            description
) VALUES (
            'soma',
            201609120001,
            'Initial create - somadbctl %s'
);`, version)
	queryMap["insertSomaSchemaVersion"] = somaString
	queries[idx] = "insertSomaSchemaVersion"
	idx++

	rootString := fmt.Sprintf(`
INSERT INTO public.schema_versions (
            schema,
            version,
            description
) VALUES (
            'root',
            201605160001,
            'Initial create - somadbctl %s'
);`, version)
	queryMap["insertRootSchemaVersion"] = rootString
	queries[idx] = "insertRootSchemaVersion"
	idx++

	performDatabaseTask(printOnly, verbose, queries[:idx], queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
