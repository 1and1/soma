package main

import "fmt"

const MaxInt = int(^uint(0) >> 1)

var UpgradeVersions = map[string]map[int]func(int, string) int{
	"inventory": map[int]func(int, string) int{
		201605060001: mock_upgrade_inventory_201605060001,
	},
	"auth": map[int]func(int, string) int{
		201605060001: upgrade_auth_to_201605150002,
	},
	"soma": map[int]func(int, string) int{
		201605060001: mock_upgrade_soma_201605060001,
		201605060002: mock_upgrade_soma_201605060002,
	},
	"root": map[int]func(int, string) int{
		000000000001: install_root_201605150001,
	},
}

func UpgradeSchema(target int, tool string) error {
	// no specific target specified => upgrade all the way
	if target == 0 {
		target = MaxInt
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
			version = f(version, tool)
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
	return nil
}

func mock_upgrade_inventory_201605060001(curr int, tool string) int {
	if curr != 201605060001 {
		return 0
	}
	return 201605060002
}

func upgrade_auth_to_201605150002(curr int, tool string) int {
	if curr != 201605060001 {
		return 0
	}
	dbOpen()

	db.Exec(`DELETE FROM auth.user_token_authentication;`)
	db.Exec(`ALTER TABLE auth.user_token_authentication ADD COLUMN salt varchar(256) NOT NULL;`)
	db.Exec(`ALTER TABLE auth.user_token_authentication RENAME TO tokens;`)
	db.Exec(`DROP TABLE auth.admin_token_authentication;`)
	db.Exec(`ALTER TABLE auth.tools ADD CHECK( left( tool_name, 5 ) = 'tool_' );`)
	db.Exec(`ALTER TABLE auth.user_authentication DROP COLUMN algorithm;`)
	db.Exec(`ALTER TABLE auth.user_authentication DROP COLUMN rounds;`)
	db.Exec(`ALTER TABLE auth.user_authentication DROP COLUMN salt;`)
	db.Exec(`ALTER TABLE auth.admin_authentication DROP COLUMN algorithm;`)
	db.Exec(`ALTER TABLE auth.admin_authentication DROP COLUMN rounds;`)
	db.Exec(`ALTER TABLE auth.admin_authentication DROP COLUMN salt;`)
	db.Exec(`ALTER TABLE auth.tool_authentication DROP COLUMN algorithm;`)
	db.Exec(`ALTER TABLE auth.tool_authentication DROP COLUMN rounds;`)
	db.Exec(`ALTER TABLE auth.tool_authentication DROP COLUMN salt;`)
	insert := fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('auth', 201605150002, 'Upgrade - somadbctl %s');", tool)
	db.Exec(insert)
	return 201605150002
}

func install_root_201605150001(curr int, tool string) int {
	if curr != 000000000001 {
		return 0
	}
	dbOpen()

	db.Exec(`CREATE SCHEMA IF NOT EXISTS root;`)
	db.Exec(`GRANT USAGE ON SCHEMA root TO soma_svc;`)
	db.Exec(`CREATE TABLE IF NOT EXISTS root.token (token varchar(256) NOT NULL);`)
	db.Exec(`CREATE TABLE IF NOT EXISTS root.flags (flag varchar(256) NOT NULL, status boolean NOT NULL DEFAULT 'no');`)
	db.Exec(`GRANT SELECT ON ALL TABLES IN SCHEMA root TO soma_svc;`)
	db.Exec(`INSERT INTO root.flags (flag, status) VALUES ('restricted', false), ('disabled', false);`)
	insert := fmt.Sprintf("INSERT INTO public.schema_versions (schema, version, description) VALUES ('root', 201605150001, 'Upgrade create - somadbctl %s');", tool)
	db.Exec(insert)

	return 201605150001
}

func mock_upgrade_soma_201605060001(curr int, tool string) int {
	if curr != 201605060001 {
		return 0
	}
	return 201605060002
}

func mock_upgrade_soma_201605060002(curr int, tool string) int {
	if curr != 201605060002 {
		return 0
	}
	return 201605060003
}

func getCurrentSchemaVersion(schema string) int {
	// TODO: needs hook to mock current version to report for no-execute
	//       case
	return 201605060001
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
