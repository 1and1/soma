package main

func createTablesInventoryAssets(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 10)

	queryMap["createTableDatacenters"] = `
create table if not exists inventory.datacenters (
  datacenter                  varchar(32)     PRIMARY KEY
);`
	queries[idx] = "createTableDatacenters"
	idx++

	queryMap["createTableServers"] = `
create table if not exists inventory.servers (
  server_id                   uuid            PRIMARY KEY,
  server_asset_id             numeric(16,0)   UNIQUE NOT NULL,
  server_datacenter_name      varchar(32)     NOT NULL REFERENCES inventory.datacenters ( datacenter ),
  server_datacenter_location  varchar(256)    NOT NULL,
  server_name                 varchar(256)    NOT NULL,
  server_online               boolean         NOT NULL DEFAULT 'yes',
  server_deleted              boolean         NOT NULL DEFAULT 'no'
);`
	queries[idx] = "createTableServers"
	idx++

	queryMap["createIndexUniqueServersOnline"] = `
create unique index _unique_server_online
  on inventory.servers ( server_name )
  where server_online
;`
	queries[idx] = "createIndexUniqueServersOnline"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

func createTablesInventoryAccounts(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 10)

	queryMap["createTableOrganizationalTeams"] = `
create table if not exists inventory.organizational_teams (
  organizational_team_id      uuid            PRIMARY KEY,
  organizational_team_name    varchar(384)    UNIQUE NOT NULL,
  organizational_team_ldap_id numeric(16,0)   UNIQUE NOT NULL,
  organizational_team_system  boolean         NOT NULL DEFAULT 'no'
);`
	queries[idx] = "createTableOrganizationalTeams"
	idx++

	queryMap["createTableOncallDuty"] = `
create table if not exists inventory.oncall_duty_teams (
  oncall_duty_id              uuid            PRIMARY KEY,
  oncall_duty_name            varchar(256)    UNIQUE NOT NULL,
  oncall_duty_phone_number    numeric(4,0)    UNIQUE NOT NULL
);`
	queries[idx] = "createTableOncallDuty"
	idx++

	queryMap["createTableUsers"] = `
create table if not exists inventory.users (
  user_id                     uuid            PRIMARY KEY,
  user_uid                    varchar(256)    UNIQUE NOT NULL,
  user_first_name             varchar(256)    NOT NULL,
  user_last_name              varchar(256)    NOT NULL,
  user_employee_number        numeric(16,0)   UNIQUE NOT NULL,
  user_mail_address           text            NOT NULL,
  user_active                 boolean         NOT NULL DEFAULT 'yes',
  user_is_system              boolean         NOT NULL DEFAULT 'no',
  organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id )
);`
	queries[idx] = "createTableUsers"
	idx++

	queryMap["createTableOncallDutyMembership"] = `
create table if not exists inventory.oncall_duty_membership (
  user_id                     uuid            NOT NULL REFERENCES inventory.users ( user_id ),
  oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_duty_teams ( oncall_duty_id )
);`
	queries[idx] = "createTableOncallDutyMembership"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
