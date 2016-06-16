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
		fmt.Fprintln(os.Stderr, "NO-EXECUTE: generated token was not inserted!\n")
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
