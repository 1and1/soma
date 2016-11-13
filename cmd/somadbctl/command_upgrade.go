package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

const MaxInt = int(^uint(0) >> 1)

var UpgradeVersions = map[string]map[int]func(int, string, bool) int{
	//	"inventory": map[int]func(int, string) int{
	//		201605060001: mock_upgrade_inventory_201605060001,
	//	},
	"auth": map[int]func(int, string, bool) int{
		201605060001: upgrade_auth_to_201605150002,
		201605150002: upgrade_auth_to_201605190001,
	},
	"soma": map[int]func(int, string, bool) int{
		201605060001: upgrade_soma_to_201605210001,
		201605210001: upgrade_soma_to_201605240001,
		201605240001: upgrade_soma_to_201605240002,
		201605240002: upgrade_soma_to_201605270001,
		201605270001: upgrade_soma_to_201605310001,
		201605310001: upgrade_soma_to_201606150001,
		201606150001: upgrade_soma_to_201606160001,
		201606160001: upgrade_soma_to_201606210001,
		201606210001: upgrade_soma_to_201607070001,
		201607070001: upgrade_soma_to_201609080001,
		201609080001: upgrade_soma_to_201609120001,
		201609120001: upgrade_soma_to_201610290001,
		201610290001: upgrade_soma_to_201611060001,
		201611060001: upgrade_soma_to_201611100001,
		201611100001: upgrade_soma_to_201611130001,
	},
	"root": map[int]func(int, string, bool) int{
		000000000001: install_root_201605150001,
		201605150001: upgrade_root_to_201605160001,
	},
}

func commandUpgradeSchema(done chan bool, target int, tool string, printOnly bool) {
	// no specific target specified => upgrade all the way
	if target == 0 {
		target = MaxInt
	}
	dbOpen()
	if printOnly {
		// in printOnly we have to close ourselve
		defer db.Close()
	}

loop:
	for schema, _ := range UpgradeVersions {
		// fetch current version from database
		version := getCurrentSchemaVersion(schema)

		if version >= target {
			// schema is already as updated as we need
			continue loop
		}

		for f, ok := UpgradeVersions[schema][version]; ok; f, ok = UpgradeVersions[schema][version] {
			version = f(version, tool, printOnly)
			if version == 0 {
				// something broke
				// TODO: set error
				break loop
			} else if version >= target {
				// job done, continue with next schema
				continue loop
			}
		}
	}
	done <- true
}

func upgrade_auth_to_201605150002(curr int, tool string, printOnly bool) int {
	if curr != 201605060001 {
		return 0
	}

	stmts := []string{
		`DELETE FROM auth.user_token_authentication;`,
		`ALTER TABLE auth.user_token_authentication ADD COLUMN salt varchar(256) NOT NULL;`,
		`ALTER TABLE auth.user_token_authentication RENAME TO tokens;`,
		`DROP TABLE auth.admin_token_authentication;`,
		`ALTER TABLE auth.tools ADD CHECK( left( tool_name, 5 ) = 'tool_' );`,
		`ALTER TABLE auth.user_authentication DROP COLUMN algorithm;`,
		`ALTER TABLE auth.user_authentication DROP COLUMN rounds;`,
		`ALTER TABLE auth.user_authentication DROP COLUMN salt;`,
		`ALTER TABLE auth.admin_authentication DROP COLUMN algorithm;`,
		`ALTER TABLE auth.admin_authentication DROP COLUMN rounds;`,
		`ALTER TABLE auth.admin_authentication DROP COLUMN salt;`,
		`ALTER TABLE auth.tool_authentication DROP COLUMN algorithm;`,
		`ALTER TABLE auth.tool_authentication DROP COLUMN rounds;`,
		`ALTER TABLE auth.tool_authentication DROP COLUMN salt;`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('auth', 201605150002, 'Upgrade - somadbctl %s');", tool),
	)

	executeUpgrades(stmts, printOnly)

	return 201605150002
}

func upgrade_auth_to_201605190001(curr int, tool string, printOnly bool) int {
	if curr != 201605150002 {
		return 0
	}

	stmts := []string{
		`ALTER TABLE auth.tokens DROP COLUMN IF EXISTS user_id;`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('auth', 201605190001, 'Upgrade - somadbctl %s');", tool),
	)
	executeUpgrades(stmts, printOnly)

	return 201605190001
}

func upgrade_soma_to_201605210001(curr int, tool string, printOnly bool) int {
	if curr != 201605060001 {
		return 0
	}
	stmts := []string{
		`ALTER TABLE soma.permissions ADD CHECK  ( permission_type != 'omnipotence' OR permission_name = 'omnipotence' );`,
		`ALTER TABLE soma.global_authorizations DROP CONSTRAINT "global_authorizations_permission_type_check";`,
		`ALTER TABLE soma.repo_authorizations DROP CONSTRAINT "repo_authorizations_permission_type_check";`,
		`ALTER TABLE soma.bucket_authorizations DROP CONSTRAINT "bucket_authorizations_permission_type_check";`,
		`ALTER TABLE soma.group_authorizations DROP CONSTRAINT "group_authorizations_permission_type_check";`,
		`ALTER TABLE soma.cluster_authorizations DROP CONSTRAINT "cluster_authorizations_permission_type_check";`,
		`ALTER TABLE soma.global_authorizations ADD CHECK ( permission_type IN ( 'omnipotence', 'grant_system', 'system', 'global' ) );`,
		`ALTER TABLE soma.global_authorizations ADD CHECK ( permission_id != '00000000-0000-0000-0000-000000000000'::uuid OR user_id = '00000000-0000-0000-0000-000000000000'::uuid );`,
		`ALTER TABLE soma.global_authorizations ADD CHECK ( permission_type IN ( 'omnipotence', 'grant_system', 'system', 'global' ) );`,
		`ALTER TABLE soma.repo_authorizations ADD CHECK ( permission_type IN ( 'grant_limited', 'limited' ) );`,
		`ALTER TABLE soma.bucket_authorizations ADD CHECK ( permission_type = 'limited' );`,
		`ALTER TABLE soma.group_authorizations ADD CHECK ( permission_type = 'limited' );`,
		`ALTER TABLE soma.cluster_authorizations ADD CHECK ( permission_type = 'limited' );`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201605210001, 'Upgrade - somadbctl %s');", tool),
	)

	executeUpgrades(stmts, printOnly)

	return 201605210001
}

func upgrade_soma_to_201605240001(curr int, tool string, printOnly bool) int {
	if curr != 201605210001 {
		return 0
	}
	stmts := []string{
		`ALTER TABLE soma.permission_types ADD COLUMN created_by uuid NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE;`,
		`ALTER TABLE soma.permission_types ADD COLUMN created_at timestamptz(3) NOT NULL DEFAULT NOW();`,
		`INSERT INTO soma.permission_types ( permission_type, created_by ) VALUES ( 'omnipotence', '00000000-0000-0000-0000-000000000000' );`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201605240001, 'Upgrade - somadbctl %s');", tool),
	)

	executeUpgrades(stmts, printOnly)

	return 201605240001
}

func upgrade_soma_to_201605240002(curr int, tool string, printOnly bool) int {
	if curr != 201605240001 {
		return 0
	}
	stmts := []string{
		`ALTER TABLE soma.permissions ADD COLUMN created_by uuid NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE;`,
		`ALTER TABLE soma.permissions ADD COLUMN created_at timestamptz(3) NOT NULL DEFAULT NOW();`,
		`INSERT INTO soma.permissions (permission_id, permission_name, permission_type, created_by )
		 VALUES ( '00000000-0000-0000-0000-000000000000','omnipotence', 'omnipotence', '00000000-0000-0000-0000-000000000000' );`,
		`INSERT INTO soma.global_authorizations ( user_id, permission_id, permission_type )
		 VALUES ( '00000000-0000-0000-0000-000000000000', '00000000-0000-0000-0000-000000000000', 'omnipotence' );`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201605240002, 'Upgrade - somadbctl %s');", tool),
	)

	executeUpgrades(stmts, printOnly)

	return 201605240002
}

func upgrade_soma_to_201605270001(curr int, tool string, printOnly bool) int {
	if curr != 201605240002 {
		return 0
	}
	stmts := []string{
		`ALTER TABLE soma.service_properties ALTER COLUMN service_property TYPE varchar(128);`,
		`ALTER TABLE soma.service_property_attributes ALTER COLUMN service_property_attribute TYPE varchar(128);`,
		`ALTER TABLE soma.service_property_values ALTER COLUMN service_property TYPE varchar(128);`,
		`ALTER TABLE soma.service_property_values ALTER COLUMN service_property_attribute TYPE varchar(128);`,
		`ALTER TABLE soma.service_property_values ALTER COLUMN value TYPE varchar(512);`,
		`ALTER TABLE soma.team_service_properties ALTER COLUMN service_property TYPE varchar(128);`,
		`ALTER TABLE soma.team_service_property_values ALTER COLUMN service_property TYPE varchar(128);`,
		`ALTER TABLE soma.team_service_property_values ALTER COLUMN service_property_attribute TYPE varchar(128);`,
		`ALTER TABLE soma.team_service_property_values ALTER COLUMN value TYPE varchar(512);`,
		`ALTER TABLE soma.constraints_service_property ALTER COLUMN service_property TYPE varchar(128);`,
		`ALTER TABLE soma.constraints_service_attribute ALTER COLUMN service_property_attribute TYPE varchar(128);`,
		`ALTER TABLE soma.constraints_service_attribute ALTER COLUMN attribute_value TYPE varchar(512);`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201605270001, 'Upgrade - somadbctl %s');", tool),
	)

	executeUpgrades(stmts, printOnly)

	return 201605270001
}

func upgrade_soma_to_201605310001(curr int, tool string, printOnly bool) int {
	if curr != 201605270001 {
		return 0
	}
	stmts := []string{
		`DELETE FROM soma.global_authorizations;`,
		`ALTER TABLE soma.global_authorizations ADD COLUMN grant_id uuid PRIMARY KEY;`,
		`ALTER TABLE soma.global_authorizations ADD COLUMN created_by uuid NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE;`,
		`ALTER TABLE soma.global_authorizations ADD COLUMN created_at timestamptz(3) NOT NULL DEFAULT NOW();`,
		`INSERT INTO soma.global_authorizations ( grant_id, user_id, permission_id, permission_type, created_by )
		 VALUES ( '00000000-0000-0000-0000-000000000000', '00000000-0000-0000-0000-000000000000', '00000000-0000-0000-0000-000000000000', 'omnipotence', '00000000-0000-0000-0000-000000000000' );`,
		`ALTER TABLE soma.global_authorizations RENAME TO authorizations_global;`,
		`DELETE FROM soma.repo_authorizations;`,
		`ALTER TABLE soma.repo_authorizations ADD COLUMN grant_id uuid PRIMARY KEY;`,
		`ALTER TABLE soma.repo_authorizations ADD COLUMN created_by uuid NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE;`,
		`ALTER TABLE soma.repo_authorizations ADD COLUMN created_at timestamptz(3) NOT NULL DEFAULT NOW();`,
		`ALTER TABLE soma.repo_authorizations RENAME TO authorizations_repository;`,
		`DELETE FROM soma.bucket_authorizations;`,
		`ALTER TABLE soma.bucket_authorizations ADD COLUMN grant_id uuid PRIMARY KEY;`,
		`ALTER TABLE soma.bucket_authorizations ADD COLUMN created_by uuid NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE;`,
		`ALTER TABLE soma.bucket_authorizations ADD COLUMN created_at timestamptz(3) NOT NULL DEFAULT NOW();`,
		`ALTER TABLE soma.bucket_authorizations RENAME TO authorizations_bucket;`,
		`DELETE FROM soma.group_authorizations;`,
		`ALTER TABLE soma.group_authorizations ADD COLUMN grant_id uuid PRIMARY KEY;`,
		`ALTER TABLE soma.group_authorizations ADD COLUMN created_by uuid NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE;`,
		`ALTER TABLE soma.group_authorizations ADD COLUMN created_at timestamptz(3) NOT NULL DEFAULT NOW();`,
		`ALTER TABLE soma.group_authorizations RENAME TO authorizations_group;`,
		`DELETE FROM soma.cluster_authorizations;`,
		`ALTER TABLE soma.cluster_authorizations ADD COLUMN grant_id uuid PRIMARY KEY;`,
		`ALTER TABLE soma.cluster_authorizations ADD COLUMN created_by uuid NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE;`,
		`ALTER TABLE soma.cluster_authorizations ADD COLUMN created_at timestamptz(3) NOT NULL DEFAULT NOW();`,
		`ALTER TABLE soma.cluster_authorizations RENAME TO authorizations_cluster;`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201605310001, 'Upgrade - somadbctl %s');", tool),
	)

	executeUpgrades(stmts, printOnly)

	return 201605310001
}

func upgrade_soma_to_201606150001(curr int, tool string, printOnly bool) int {
	if curr != 201605310001 {
		return 0
	}
	stmts := []string{
		`ALTER TABLE soma.repositories ADD COLUMN created_by uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000' REFERENCES inventory.users ( user_id ) DEFERRABLE;`,
		`ALTER TABLE soma.repositories ADD COLUMN created_at timestamptz(3) NOT NULL DEFAULT NOW();`,
		`CREATE UNIQUE INDEX _singleton_default_bucket ON soma.buckets ( organizational_team_id, environment ) WHERE environment = 'default';`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201606150001, 'Upgrade - somadbctl %s');", tool),
	)

	executeUpgrades(stmts, printOnly)

	return 201606150001
}

func upgrade_soma_to_201606160001(curr int, tool string, printOnly bool) int {
	if curr != 201606150001 {
		return 0
	}
	stmts := []string{
		`ALTER TABLE soma.jobs ADD COLUMN job_error text NOT NULL DEFAULT '';`,
		`INSERT INTO soma.job_types ( job_type ) VALUES ('remove_check_from_repository'), ('remove_check_from_bucket'), ('remove_check_from_group'), ('remove_check_from_cluster'), ('remove_check_from_node');`,
		`ALTER TABLE soma.check_configurations ADD COLUMN deleted boolean NOT NULL DEFAULT 'no'::boolean;`,
		`ALTER TABLE soma.checks ADD COLUMN deleted boolean NOT NULL DEFAULT 'no'::boolean;`,
		`ALTER TABLE soma.check_configurations ADD UNIQUE ( repository_id, configuration_name ) DEFERRABLE;`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201606160001, 'Upgrade - somadbctl %s');", tool),
	)

	executeUpgrades(stmts, printOnly)

	return 201606160001
}

func upgrade_soma_to_201606210001(curr int, tool string, printOnly bool) int {
	if curr != 201606160001 {
		return 0
	}
	stmts := []string{
		`ALTER TABLE soma.check_configurations DROP CONSTRAINT check_configurations_repository_id_configuration_name_key;`,
		`CREATE UNIQUE INDEX _singleton_undeleted_checkconfig_name ON soma.check_configurations ( repository_id, configuration_name ) WHERE NOT deleted;`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201606210001, 'Upgrade - somadbctl %s');", tool),
	)

	executeUpgrades(stmts, printOnly)

	return 201606210001
}

func upgrade_soma_to_201607070001(curr int, tool string, printOnly bool) int {
	if curr != 201606210001 {
		return 0
	}
	stmts := []string{
		`CREATE INDEX CONCURRENTLY _checks_to_instances ON soma.check_instances ( check_id, check_instance_id );`,
		`CREATE INDEX CONCURRENTLY _repo_to_checks ON checks ( repository_id, check_id );`,
		`CREATE INDEX CONCURRENTLY _instance_to_config ON soma.check_instance_configurations ( check_instance_id, check_instance_config_id );`,
		`CREATE INDEX CONCURRENTLY _config_dependencies ON soma.check_instance_configuration_dependencies ( blocked_instance_config_id, blocking_instance_config_id );`,
		`CREATE INDEX CONCURRENTLY _instance_config_status ON soma.check_instance_configurations ( status, check_instance_id );`,
		`CREATE UNIQUE INDEX CONCURRENTLY _instance_config_version ON check_instance_configurations ( check_instance_id, version );`,
		`CREATE INDEX CONCURRENTLY _configuration_id_levels ON configuration_thresholds ( configuration_id, notification_level );`,
		`CREATE TABLE IF NOT EXISTS soma.authorizations_monitoring ( grant_id uuid PRIMARY KEY, user_id uuid REFERENCES inventory.users ( user_id ) DEFERRABLE, tool_id uuid REFERENCES auth.tools ( tool_id ) DEFERRABLE, organizational_team_id uuid REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE, monitoring_id uuid NOT NULL REFERENCES soma.monitoring_systems ( monitoring_id ) DEFERRABLE, permission_id uuid NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE, permission_type varchar(32) NOT NULL REFERENCES soma.permission_types ( permission_type ) DEFERRABLE, created_by uuid NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE, created_at timestamptz(3) NOT NULL DEFAULT NOW(), FOREIGN KEY ( permission_id, permission_type ) REFERENCES soma.permissions ( permission_id, permission_type ) DEFERRABLE, CHECK (( user_id IS NOT NULL AND tool_id IS NULL AND organizational_team_id IS NULL ) OR ( user_id IS NULL AND tool_id IS NOT NULL AND organizational_team_id IS NULL ) OR ( user_id IS NULL AND tool_id IS NULL AND organizational_team_id IS NOT NULL )), CHECK ( permission_type = 'limited' ));`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201607070001, 'Upgrade - somadbctl %s');", tool),
	)

	executeUpgrades(stmts, printOnly)

	return 201607070001
}

func upgrade_soma_to_201609080001(curr int, tool string, printOnly bool) int {
	if curr != 201607070001 {
		return 0
	}
	stmts := []string{
		`ALTER TABLE soma.buckets ADD COLUMN created_by uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000' REFERENCES inventory.users ( user_id ) DEFERRABLE;`,
		`ALTER TABLE soma.buckets ADD COLUMN created_at timestamptz(3) NOT NULL DEFAULT NOW();`,
		`ALTER TABLE soma.groups ADD COLUMN created_by uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000' REFERENCES inventory.users ( user_id ) DEFERRABLE;`,
		`ALTER TABLE soma.groups ADD COLUMN created_at timestamptz(3) NOT NULL DEFAULT NOW();`,
		`ALTER TABLE soma.clusters ADD COLUMN created_by uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000' REFERENCES inventory.users ( user_id ) DEFERRABLE;`,
		`ALTER TABLE soma.clusters ADD COLUMN created_at timestamptz(3) NOT NULL DEFAULT NOW();`,
		`ALTER TABLE soma.nodes ADD COLUMN created_by uuid NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000' REFERENCES inventory.users ( user_id ) DEFERRABLE;`,
		`ALTER TABLE soma.nodes ADD COLUMN created_at timestamptz(3) NOT NULL DEFAULT NOW();`,
		`INSERT INTO soma.job_types ( job_type ) VALUES ( 'delete_system_property_from_repository' ), ( 'delete_custom_property_from_repository' ), ( 'delete_oncall_property_from_repository' ), ( 'delete_service_property_from_repository' ), ( 'delete_system_property_from_bucket' ), ( 'delete_custom_property_from_bucket' ), ( 'delete_oncall_property_from_bucket' ), ( 'delete_service_property_from_bucket' ), ( 'delete_system_property_from_group' ), ( 'delete_custom_property_from_group' ), ( 'delete_oncall_property_from_group' ), ( 'delete_service_property_from_group' ), ( 'delete_system_property_from_cluster' ), ( 'delete_custom_property_from_cluster' ), ( 'delete_oncall_property_from_cluster' ), ( 'delete_service_property_from_cluster' ), ( 'delete_system_property_from_node' ), ( 'delete_custom_property_from_node' ), ( 'delete_oncall_property_from_node' ), ( 'delete_service_property_from_node' );`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201609080001, 'Upgrade - somadbctl %s');", tool),
	)
	executeUpgrades(stmts, printOnly)

	return 201609080001
}

func upgrade_soma_to_201609120001(curr int, tool string, printOnly bool) int {
	if curr != 201609080001 {
		return 0
	}
	stmts := []string{
		`create unique index _unique_admin_global_authoriz on soma.authorizations_global ( admin_id, permission_id ) where admin_id IS NOT NULL;`,
		`create unique index _unique_user_global_authoriz on soma.authorizations_global ( user_id, permission_id ) where user_id IS NOT NULL;`,
		`create unique index _unique_tool_global_authoriz on soma.authorizations_global ( tool_id, permission_id ) where tool_id IS NOT NULL;`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201609120001, 'Upgrade - somadbctl %s');", tool),
	)
	executeUpgrades(stmts, printOnly)

	return 201609120001
}

func upgrade_soma_to_201610290001(curr int, tool string, printOnly bool) int {
	if curr != 201609120001 {
		return 0
	}
	stmts := []string{
		`ALTER TABLE soma.check_instance_configurations ADD COLUMN deprovisioned_at timestamptz(3) NULL;`,
		`ALTER TABLE soma.check_instance_configurations ADD COLUMN status_last_updated_at timestamptz(3) NULL;`,
		`ALTER TABLE soma.check_instance_configurations ADD COLUMN notified_at timestamptz(3) NULL;`,
		`SET TIME ZONE 'UTC';`,
		`UPDATE soma.check_instance_configurations SET deprovisioned_at = NOW()::timestamptz(3), status_last_updated_at = NOW()::timestamptz(3) WHERE status IN ('deprovisioned', 'awaiting_deletion');`,
		`UPDATE soma.check_instance_configurations SET status_last_updated_at = NOW()::timestamptz(3) WHERE status IN ('awaiting_rollout', 'rollout_in_progress', 'rollout_failed', 'active', 'awaiting_deprovision', 'deprovision_in_progress','deprovision_failed');`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201610290001, 'Upgrade - somadbctl %s');", tool),
	)
	executeUpgrades(stmts, printOnly)

	return 201610290001
}

func upgrade_soma_to_201611060001(curr int, tool string, printOnly bool) int {
	if curr != 201610290001 {
		return 0
	}
	stmts := []string{}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201611060001, 'Upgrade - somadbctl %s');", tool),
	)
	executeUpgrades(stmts, printOnly)

	return 201611060001
}

func upgrade_soma_to_201611100001(curr int, tool string, printOnly bool) int {
	if curr != 201611060001 {
		return 0
	}
	stmts := []string{
		`ALTER TABLE soma.permission_types RENAME TO categories;`,
		`ALTER TABLE soma.categories RENAME permission_type TO category;`,
		`create table if not exists soma.sections ( section_id uuid PRIMARY KEY, section_name varchar(64) UNIQUE NOT NULL, category varchar(32) NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE, created_by uuid NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE, created_at timestamptz(3) NOT NULL DEFAULT NOW(), UNIQUE ( section_id, category ), UNIQUE( section_name ));`,
		`create table if not exists soma.actions ( action_id uuid PRIMARY KEY, action_name varchar(64) NOT NULL, section_id uuid NOT NULL REFERENCES soma.sections ( section_id ) DEFERRABLE, category varchar(32) NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE, created_by uuid NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE, created_at timestamptz(3) NOT NULL DEFAULT NOW(), UNIQUE ( section_id, action_name ), UNIQUE ( section_id, action_id ), FOREIGN KEY ( section_id, category ) REFERENCES soma.sections ( section_id, category ) DEFERRABLE );`,
		`ALTER TABLE soma.permissions RENAME permission_type TO category;`,
		`ALTER TABLE permissions DROP CONSTRAINT permissions_permission_name_key;`,
		`ALTER TABLE soma.permissions ADD CONSTRAINT permissions_permission_name_category_key UNIQUE (permission_name, category );`,
		`create table if not exists soma.permission_map ( mapping_id uuid PRIMARY KEY, category varchar(32) NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE, permission_id uuid NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE, section_id uuid NOT NULL REFERENCES soma.sections ( section_id ) DEFERRABLE, action_id uuid REFERENCES soma.actions ( action_id ) DEFERRABLE, FOREIGN KEY ( permission_id, category ) REFERENCES soma.permissions ( permission_id, category ), FOREIGN KEY ( section_id, category ) REFERENCES soma.sections ( section_id, category ), FOREIGN KEY ( section_id, action_id ) REFERENCES soma.actions ( section_id, action_id ));`,
		`create table if not exists soma.permission_grant_map ( category varchar(32) NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE, permission_id uuid NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE, granted_category varchar(32) NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE, granted_permission_id uuid NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE, FOREIGN KEY ( permission_id, category ) REFERENCES soma.permissions ( permission_id, category ), FOREIGN KEY ( granted_permission_id, granted_category ) REFERENCES soma.permissions ( permission_id, category ), CHECK ( permission_id != granted_permission_id ), CHECK ( category ~ ':grant$' ), CHECK ( granted_category = substring(category from '^([^:]+):')));`,
		`ALTER TABLE soma.authorizations_global RENAME permission_type TO category;`,
		`ALTER TABLE soma.authorizations_global ADD COLUMN organizational_team_id uuid REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE;`,
		`ALTER TABLE soma.authorizations_global DROP CONSTRAINT authorizations_global_check;`,
		`ALTER TABLE soma.authorizations_global DROP CONSTRAINT authorizations_global_check1;`,
		`ALTER TABLE soma.authorizations_global DROP CONSTRAINT authorizations_global_permission_type_check;`,
		`ALTER TABLE soma.authorizations_global ADD CONSTRAINT authorizations_global_admin_id_check CHECK ( admin_id IS NULL OR category != 'system' );`,
		`ALTER TABLE soma.authorizations_global ADD CONSTRAINT authorizations_global_category_check CHECK ( category IN ( 'omnipotence','system','global','global:grant','permission','permission:grant','operations','operations:grant' ));`,
		`ALTER TABLE soma.authorizations_global ADD CONSTRAINT authorizations_global_check CHECK ( ( admin_id IS NOT NULL AND user_id IS NULL AND tool_id IS NULL AND organizational_team_id IS NULL ) OR ( admin_id IS NULL AND user_id IS NOT NULL AND tool_id IS NULL AND organizational_team_id IS NULL ) OR ( admin_id IS NULL AND user_id IS NULL AND tool_id IS NOT NULL AND organizational_team_id IS NULL ) OR ( admin_id IS NULL AND user_id IS NULL AND tool_id IS NULL AND organizational_team_id IS NOT NULL ));`,
		`ALTER TABLE soma.authorizations_global ADD CONSTRAINT authorizations_global_check1 CHECK ( permission_id != '00000000-0000-0000-0000-000000000000'::uuid OR user_id = '00000000-0000-0000-0000-000000000000'::uuid );`,
		`ALTER TABLE soma.authorizations_repository RENAME permission_type TO category;`,
		`ALTER TABLE soma.authorizations_repository DROP CONSTRAINT authorizations_repository_permission_type_check;`,
		`ALTER TABLE soma.authorizations_repository ADD CONSTRAINT authorizations_repository_category_check CHECK ( category IN ( 'repository', 'repository:grant' ));`,
		`ALTER TABLE soma.authorizations_bucket RENAME permission_type TO category;`,
		`ALTER TABLE soma.authorizations_bucket DROP CONSTRAINT authorizations_bucket_permission_type_check;`,
		`ALTER TABLE soma.authorizations_bucket ADD CONSTRAINT authorizations_bucket_category_check CHECK ( category = 'repository' );`,
		`ALTER TABLE soma.authorizations_group RENAME permission_type TO category;`,
		`ALTER TABLE soma.authorizations_group DROP CONSTRAINT authorizations_group_permission_type_check;`,
		`ALTER TABLE soma.authorizations_group ADD CONSTRAINT authorizations_group_category_check CHECK ( category = 'repository' );`,
		`ALTER TABLE soma.authorizations_cluster RENAME permission_type TO category;`,
		`ALTER TABLE soma.authorizations_cluster DROP CONSTRAINT authorizations_cluster_permission_type_check;`,
		`ALTER TABLE soma.authorizations_cluster ADD CONSTRAINT authorizations_cluster_category_check CHECK ( category = 'repository' );`,
		`ALTER TABLE soma.authorizations_monitoring RENAME permission_type TO category;`,
		`ALTER TABLE soma.authorizations_monitoring DROP CONSTRAINT authorizations_monitoring_permission_type_check;`,
		`ALTER TABLE soma.authorizations_monitoring ADD CONSTRAINT authorizations_monitoring_category_check CHECK ( category IN ( 'monitoring','monitoring:grant' ));`,
		`create table if not exists soma.authorizations_team ( grant_id uuid PRIMARY KEY, user_id uuid REFERENCES inventory.users ( user_id ) DEFERRABLE, tool_id uuid REFERENCES auth.tools ( tool_id ) DEFERRABLE, organizational_team_id uuid REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE, authorized_team_id uuid NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE, permission_id uuid NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE, category varchar(32) NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE, created_by uuid NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE, created_at timestamptz(3) NOT NULL DEFAULT NOW(), FOREIGN KEY ( permission_id, category ) REFERENCES soma.permissions ( permission_id, category ) DEFERRABLE, CHECK (( user_id IS NOT NULL AND tool_id IS NULL AND organizational_team_id IS NULL ) OR ( user_id IS NULL AND tool_id IS NOT NULL AND organizational_team_id IS NULL ) OR ( user_id IS NULL AND tool_id IS NULL AND organizational_team_id IS NOT NULL )), CHECK ( category IN ( 'team', 'team:grant' )));`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201611100001, 'Upgrade - somadbctl %s');", tool),
	)
	executeUpgrades(stmts, printOnly)

	return 201611100001
}

func upgrade_soma_to_201611130001(curr int, tool string, printOnly bool) int {
	if curr != 201611100001 {
		return 0
	}
	stmts := []string{
		`DELETE FROM soma.authorizations_global WHERE category = 'system';`,
		`ALTER TABLE soma.authorizations_global DROP CONSTRAINT authorizations_global_admin_id_check;`,
		`ALTER TABLE soma.authorizations_global ADD CONSTRAINT authorizations_global_admin_id_check CHECK ( admin_id IS NULL OR category = 'system' );`,
		`ALTER TABLE soma.authorizations_global ADD CONSTRAINT authorizations_global_admin_id_check1 CHECK ( category != 'system' OR admin_id IS NOT NULL );`,
		`ALTER TABLE soma.authorizations_repository ALTER COLUMN repository_id DROP NOT NULL;`,
		`DELETE FROM soma.authorizations_repository;`,
		`ALTER TABLE soma.authorizations_repository ADD COLUMN object_type varchar(64) NOT NULL REFERENCES soma.object_types ( object_type ) DEFERRABLE;`,
		//XXX
		`ALTER TABLE soma.authorizations_repository ADD COLUMN bucket_id uuid REFERENCES soma.buckets ( bucket_id ) DEFERRABLE;`,
		`ALTER TABLE soma.authorizations_repository ADD COLUMN group_id uuid REFERENCES soma.groups ( group_id ) DEFERRABLE;`,
		`ALTER TABLE soma.authorizations_repository ADD COLUMN cluster_id uuid REFERENCES soma.clusters ( cluster_id ) DEFERRABLE;`,
		`ALTER TABLE soma.authorizations_repository ADD COLUMN node_id uuid REFERENCES soma.nodes ( node_id ) DEFERRABLE;`,
		`ALTER TABLE soma.authorizations_repository ADD FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE;`,
		`ALTER TABLE soma.authorizations_repository ADD FOREIGN KEY ( bucket_id, group_id ) REFERENCES soma.groups ( bucket_id, group_id ) DEFERRABLE;`,
		`ALTER TABLE soma.authorizations_repository ADD FOREIGN KEY ( bucket_id, cluster_id ) REFERENCES soma.clusters ( bucket_id, cluster_id ) DEFERRABLE;`,
		`ALTER TABLE soma.authorizations_repository ADD FOREIGN KEY ( node_id, bucket_id ) REFERENCES soma.node_bucket_assignment ( node_id, bucket_id ) DEFERRABLE;`,
		`ALTER TABLE soma.authorizations_repository ADD CHECK ( object_type IN ( 'repository', 'bucket', 'group', 'cluster', 'node' ));`,
		`ALTER TABLE soma.authorizations_repository ADD CHECK ( object_type != 'repository' OR repository_id IS NOT NULL );`,
		`ALTER TABLE soma.authorizations_repository ADD CHECK ( object_type != 'bucket' OR bucket_id IS NOT NULL );`,
		`ALTER TABLE soma.authorizations_repository ADD CHECK ( object_type != 'group' OR group_id IS NOT NULL );`,
		`ALTER TABLE soma.authorizations_repository ADD CHECK ( object_type != 'cluster' OR cluster_id IS NOT NULL );`,
		`ALTER TABLE soma.authorizations_repository ADD CHECK ( object_type != 'node' OR node_id IS NOT NULL );`,
		`ALTER TABLE soma.authorizations_repository ADD CHECK ( ( repository_id IS NOT NULL AND bucket_id IS NULL AND group_id IS NULL AND cluster_id IS NULL AND node_id IS NULL ) OR ( repository_id IS NOT NULL AND bucket_id IS NOT NULL AND group_id IS NULL AND cluster_id IS NULL AND node_id IS NULL ) OR ( repository_id IS NOT NULL AND bucket_id IS NOT NULL AND group_id IS NOT NULL AND cluster_id IS NULL AND node_id IS NULL ) OR ( repository_id IS NOT NULL AND bucket_id IS NOT NULL AND group_id IS NULL AND cluster_id IS NOT NULL AND node_id IS NULL ) OR ( repository_id IS NOT NULL AND bucket_id IS NOT NULL AND group_id IS NULL AND cluster_id IS NULL AND node_id IS NOT NULL ));`,
		`DROP TABLE soma.authorizations_bucket;`,
		`DROP TABLE soma.authorizations_cluster;`,
		`DROP TABLE soma.authorizations_group;`,
		`ALTER TABLE soma.authorizations_monitoring ADD UNIQUE ( user_id, tool_id, organizational_team_id, category, permission_id, monitoring_id );`,
		`ALTER TABLE soma.authorizations_team ADD UNIQUE ( user_id, tool_id, organizational_team_id, category, permission_id, authorized_team_id );`,
		`ALTER TABLE soma.authorizations_repository ADD UNIQUE ( user_id, tool_id, organizational_team_id, category, permission_id, object_type, repository_id, bucket_id, group_id, cluster_id, node_id );`,
		`ALTER TABLE soma.authorizations_global ADD UNIQUE( admin_id, user_id, tool_id, organizational_team_id, category, permission_id );`,
		`ALTER TABLE soma.permission_grant_map ADD UNIQUE ( permission_id );`,
		`ALTER TABLE soma.permission_grant_map ADD UNIQUE ( granted_permission_id );`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('soma', 201611130001, 'Upgrade - somadbctl %s');", tool),
	)
	executeUpgrades(stmts, printOnly)

	return 201611130001
}

func install_root_201605150001(curr int, tool string, printOnly bool) int {
	if curr != 000000000001 {
		return 0
	}

	stmts := []string{
		`CREATE SCHEMA IF NOT EXISTS root;`,
		`GRANT USAGE ON SCHEMA root TO soma_svc;`,
		`CREATE TABLE IF NOT EXISTS root.token (token varchar(256) NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS root.flags (flag varchar(256) NOT NULL, status boolean NOT NULL DEFAULT 'no');`,
		`GRANT SELECT ON ALL TABLES IN SCHEMA root TO soma_svc;`,
		`INSERT INTO root.flags (flag, status) VALUES ('restricted', false), ('disabled', false);`,
	}
	stmts = append(stmts,
		fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('root', 201605150001, 'Upgrade create - somadbctl %s');", tool),
	)

	executeUpgrades(stmts, printOnly)

	return 201605150001
}

func upgrade_root_to_201605160001(curr int, tool string, printOnly bool) int {
	if curr != 201605150001 {
		return 0
	}

	token := generateToken()
	if token == "" {
		return 0
	}
	istmt := `INSERT INTO root.token ( token ) VALUES ( $1::varchar );`
	vstmt := fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('root', 201605160001, 'Upgrade - somadbctl %s');", tool)
	if !printOnly {
		db.Exec(istmt, token)
		db.Exec(vstmt)
	} else {
		fmt.Println(vstmt)
	}
	fmt.Fprintf(os.Stderr, "The generated boostrap token was: %s\n", token)
	if printOnly {
		fmt.Fprintln(os.Stderr, "NO-EXECUTE: generated token was not inserted!")
	}
	return 201605160001
}

func executeUpgrades(stmts []string, printOnly bool) {
	var tx *sql.Tx

	if !printOnly {
		tx, _ = db.Begin()
		defer tx.Rollback()
		tx.Exec(`SET CONSTRAINTS ALL DEFERRED;`)
	}

	for _, stmt := range stmts {
		if printOnly {
			fmt.Println(stmt)
			continue
		}
		tx.Exec(stmt)
	}

	if !printOnly {
		tx.Commit()
	}
}

func getCurrentSchemaVersion(schema string) int {
	var (
		version int
		err     error
	)
	stmt := `SELECT MAX(version) AS version FROM public.schema_versions WHERE schema = $1::varchar GROUP BY schema;`
	if err = db.QueryRow(stmt, schema).Scan(&version); err == sql.ErrNoRows {
		return 000000000001
	} else if err != nil {
		log.Fatal(err)
	}
	return version
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
