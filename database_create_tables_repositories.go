package main

func createTableRepositories(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 5)

	queryMap["createTableRepositories"] = `
create table if not exists soma.repositories (
    repository_id               uuid            PRIMARY KEY,
    repository_name             varchar(128)    UNIQUE NOT NULL,
    repository_deleted          boolean         NOT NULL DEFAULT 'no',
    repository_active           boolean         NOT NULL DEFAULT 'yes',
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    UNIQUE( repository_id, organizational_team_id )
);`
	queries[idx] = "createTableRepositories"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

func createTablesRepositoryProperties(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 5)

	queryMap["createTableRepositoryOncallProperty"] = `
create table if not exists soma.repository_oncall_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_duty_teams ( oncall_duty_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableRepositoryOncallProperty"
	idx++

	queryMap["createTableRepositoryServiceProperty"] = `
create table if not exists soma.repository_service_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    service_property            varchar(64)     NOT NULL,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    FOREIGN KEY( organizational_team_id, service_property ) REFERENCES soma.team_service_properties ( organizational_team_id, service_property ) DEFERRABLE,
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableRepositoryServiceProperty"
	idx++

	queryMap["createTableRepositorySystemProperties"] = `
create table if not exists soma.repository_system_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    system_property             varchar(64)     NOT NULL REFERENCES soma.system_properties ( system_property ) DEFERRABLE,
    object_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    inherited                   boolean         NOT NULL DEFAULT 'yes',
    value                       text            NOT NULL,
    FOREIGN KEY ( system_property, object_type, inherited ) REFERENCES soma.system_property_validity ( system_property, object_type, inherited ) DEFERRABLE,
    CHECK( object_type = 'repository' ),
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableRepositorySystemProperties"
	idx++

	queryMap["createTableRepositoryCustomProperty"] = `
create table if not exists soma.repository_custom_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    custom_property_id          uuid            NOT NULL REFERENCES soma.custom_properties ( custom_property_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    FOREIGN KEY ( repository_id, custom_property_id ) REFERENCES soma.custom_properties ( repository_id, custom_property_id ) DEFERRABLE,
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableRepositoryCustomProperty"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
