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
    check_id                    uuid            NOT NULL REFERENCES soma.checks ( check_id ) DEFERRABLE,
    check_configuration_id      uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ) DEFERRABLE,
    current_instance_config_id  uuid            NOT NULL,
    last_configuration_created  timestamptz(3)  NOT NULL DEFAULT NOW()::timestamptz(3),
    update_available            boolean         NOT NULL DEFAULT 'no',
    deleted                     boolean         NOT NULL DEFAULT 'no'
);`
	queries[idx] = "createTableCheckInstances"
	idx++

	queryMap["createTableCheckInstanceConfigurations"] = `
create table if not exists soma.check_instance_configurations (
    check_instance_config_id    uuid            PRIMARY KEY,
    version                     integer         NOT NULL,
    check_instance_id           uuid            NOT NULL REFERENCES soma.check_instances ( check_instance_id ) DEFERRABLE,
    monitoring_id               uuid            REFERENCES soma.monitoring_systems ( monitoring_id ) DEFERRABLE,
    constraint_hash             varchar(256)    NOT NULL,
    constraint_val_hash         varchar(256)    NOT NULL,
    instance_service            varchar(64)     NOT NULL,
    instance_service_cfg_hash   varchar(256)    NOT NULL,
    instance_service_cfg        jsonb           NOT NULL,
    created                     timestamptz(3)  NOT NULL DEFAULT NOW()::timestamptz(3),
    activated_at                timestamptz(3),
    status                      varchar(32)     NOT NULL REFERENCES soma.check_instance_status ( status ) DEFERRABLE,
    next_status                 varchar(32)     NOT NULL REFERENCES soma.check_instance_status ( status ) DEFERRABLE,
    awaiting_deletion           boolean         NOT NULL DEFAULT 'no',
    deployment_details          jsonb           NOT NULL,
    CHECK ( status != 'none' ),
    CHECK ( status = 'awaiting_computation' OR monitoring_id IS NOT NULL )
);`
	queries[idx] = "createTableCheckInstanceConfigurations"
	idx++

	queryMap["createUniqueIndexActiveConfigurations"] = `
create unique index _unique_check_instance_configurations_active
    on soma.check_instance_configurations ( check_instance_id, status )
    where status = 'active' or status = 'awaiting_deprovision'
       or status = 'deprovision_in_progress' or status = 'deprovision_failed'
       or status = 'rollout_in_progress';`
	queries[idx] = "createUniqueIndexActiveConfigurations"
	idx++

	queryMap["createTableCheckInstanceConfigurationDependencies"] = `
create table if not exists soma.check_instance_configuration_dependencies (
    blocked_instance_config_id  uuid            NOT NULL REFERENCES soma.check_instance_configurations ( check_instance_config_id ) DEFERRABLE,
    blocking_instance_config_id uuid            NOT NULL REFERENCES soma.check_instance_configurations ( check_instance_config_id ) DEFERRABLE,
    unblocking_state            varchar(32)     NOT NULL REFERENCES soma.check_instance_status ( status ) DEFERRABLE
);`
	queries[idx] = "createTableCheckInstanceConfigurationDependencies"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
