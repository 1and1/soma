package main

func createTablesInstances(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 5)

	queryMap["createTableCheckInstanceStatus"] = `
create table if not exists soma.check_instance_status (
  status                      varchar(32)     PRIMARY KEY
);`
	queries[idx] = "createTableCheckInstanceStatus"
	idx++

	queryMap["createTableCheckInstances"] = `
create table if not exists soma.check_instances (
  check_instance_id           uuid            PRIMARY KEY,
  last_configuration_created  timestamptz(3)  NOT NULL DEFAULT NOW(),
  update_available            boolean         NOT NULL DEFAULT 'no',
  deleted                     boolean         NOT NULL DEFAULT 'no'
);`
	queries[idx] = "createTableCheckInstances"
	idx++

	queryMap["createTableCheckInstanceConfigurations"] = `
create table if not exists soma.check_instance_configurations (
  check_instance_config_id    uuid            PRIMARY KEY,
  check_id                    uuid            NOT NULL REFERENCES soma.check_instances ( check_instance_id ),
  monitoring_id               uuid            NOT NULL REFERENCES soma.monitoring_systems ( monitoring_id ),
  created                     timestamptz(3)  NOT NULL DEFAULT NOW(),
  activated_at                timestamptz(3)  NOT NULL DEFAULT NOW(),
  status                      varchar(32)     NOT NULL REFERENCES soma.check_instance_status ( status ),
  next_status                 varchar(32)     NOT NULL REFERENCES soma.check_instance_status ( status ),
  awaiting_deletion           boolean         NOT NULL DEFAULT 'no',
  configuration               jsonb           NOT NULL,
  CHECK ( status != 'none' )
);`
	queries[idx] = "createTableCheckInstanceConfigurations"
	idx++

	queryMap["createUniqueIndexActiveConfigurations"] = `
create unique index _unique_check_instance_configurations_active
  on soma.check_instance_configurations ( check_id, status )
  where status = 'active';`
	queries[idx] = "createUniqueIndexActiveConfigurations"
	idx++

	queryMap["createTableCheckInstanceConfigurationDependencies"] = `
create table if not exists soma.check_instance_configuration_dependencies (
  blocked_instance_config_id  uuid            NOT NULL REFERENCES soma.check_instance_configurations ( check_instance_config_id ),
  blocking_instance_config_id uuid            NOT NULL REFERENCES soma.check_instance_configurations ( check_instance_config_id ),
  unblocking_state            varchar(32)     NOT NULL REFERENCES soma.check_instance_status ( status )
);`
	queries[idx] = "createTableCheckInstanceConfigurationDependencies"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
