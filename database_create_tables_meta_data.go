package main

func createTablesMetaData(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 5)

	queryMap["createTableViews"] = `
create table if not exists soma.views (
  view                        varchar(64)     PRIMARY KEY
);`
	queries[idx] = "createTableViews"
	idx++

	queryMap["createTableEnvironments"] = `
create table if not exists soma.environments (
  environment                 varchar(32)     PRIMARY KEY
);`
	queries[idx] = "createTableEnvironments"
	idx++

	queryMap["createTableObjectStates"] = `
create table if not exists soma.object_states (
  object_state                varchar(64)     PRIMARY KEY
);`
	queries[idx] = "createTableObjectStates"
	idx++

	queryMap["createTableObjectTypes"] = `
create table if not exists soma.object_types (
  object_type                 varchar(64)     PRIMARY KEY
);`
	queries[idx] = "createTableObjectTypes"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

func createTablesDatacenterMetaData(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 5)

	queryMap["createTableDatacenterGroups"] = `
create table if not exists soma.datacenter_groups (
  datacenter_group            varchar(32)     NOT NULL,
  datacenter                  varchar(32)     NOT NULL REFERENCES inventory.datacenters ( datacenter )
);`
	queries[idx] = "createTableDatacenterGroups"
	idx++

	queryMap["createIndexDatacenterGroups"] = `
create index _datacenter_groups
  on soma.datacenter_groups ( datacenter_group );`
	queries[idx] = "createTableDatacenterGroups"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
