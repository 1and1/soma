package main

func createTablesProperties(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 10)

	queryMap["createTableServiceProperties"] = `
create table if not exists soma.service_properties (
    service_property            varchar(128)    PRIMARY KEY
);`
	queries[idx] = "createTableServiceProperties"
	idx++

	queryMap["createTableServicePropertyAttributes"] = `
create table if not exists soma.service_property_attributes (
    service_property_attribute  varchar(128)    PRIMARY KEY,
    cardinality                 varchar(8)      NOT NULL DEFAULT 'multi'
);`
	queries[idx] = "createTableServicePropertyAttributes"
	idx++

	queryMap["createTableServicePropertyValues"] = `
create table if not exists soma.service_property_values (
    service_property            varchar(128)    NOT NULL REFERENCES soma.service_properties ( service_property ) DEFERRABLE,
    service_property_attribute  varchar(128)    NOT NULL REFERENCES soma.service_property_attributes ( service_property_attribute ) DEFERRABLE,
    value                       varchar(512)    NOT NULL
);`
	queries[idx] = "createTableServicePropertyValues"
	idx++

	queryMap["createTableTeamServiceProperties"] = `
create table if not exists soma.team_service_properties (
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    service_property            varchar(128)    NOT NULL,
    UNIQUE( organizational_team_id, service_property )
);`
	queries[idx] = "createTableTeamServiceProperties"
	idx++

	queryMap["createTableTeamServicePropertyValues"] = `
create table if not exists soma.team_service_property_values (
    organizational_team_id      uuid            NOT NULL,
    service_property            varchar(128)    NOT NULL,
    service_property_attribute  varchar(128)    NOT NULL REFERENCES soma.service_property_attributes ( service_property_attribute ) DEFERRABLE,
    value                       varchar(512)    NOT NULL,
    FOREIGN KEY( organizational_team_id, service_property ) REFERENCES soma.team_service_properties ( organizational_team_id, service_property ) DEFERRABLE
);`
	queries[idx] = "createTableTeamServicePropertyValues"
	idx++

	queryMap["createTableSystemProperties"] = `
create table if not exists soma.system_properties (
    system_property             varchar(128)    PRIMARY KEY
);`
	queries[idx] = "createTableSystemProperties"
	idx++

	queryMap["createTableSystemPropertyValidity"] = `
create table if not exists soma.system_property_validity (
    system_property             varchar(128)    NOT NULL REFERENCES soma.system_properties ( system_property ) DEFERRABLE,
    object_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    inherited                   boolean         NOT NULL DEFAULT 'yes',
    UNIQUE( system_property, object_type, inherited )
);`
	queries[idx] = "createTableSystemPropertyValidity"
	idx++

	queryMap["createTableNativeProperties"] = `
create table if not exists soma.native_properties (
    native_property             varchar(128)    PRIMARY KEY
);`
	queries[idx] = "createTableNativeProperties"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

func createTableCustomProperties(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 5)

	queryMap["createTableCustomProperties"] = `
create table if not exists soma.custom_properties (
    custom_property_id          uuid            PRIMARY KEY,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    custom_property             varchar(128)    NOT NULL,
    UNIQUE( repository_id, custom_property ),
    UNIQUE( repository_id, custom_property_id )
);`
	queries[idx] = "createTableCustomProperties"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
