package main

func createTablesChecks(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 12)

	queryMap["createTableConfigurationPredicates"] = `
create table if not exists soma.configuration_predicates (
  predicate                   varchar(4)      PRIMARY KEY
);`
	queries[idx] = "createTableConfigurationPredicates"
	idx++

	queryMap["createTableNotificationLevels"] = `
create table if not exists soma.notification_levels (
  level_name                  varchar(16)     PRIMARY KEY,
  level_shortname             varchar(16)     UNIQUE NOT NULL,
  level_numeric               smallint        UNIQUE NOT NULL,
  CHECK ( level_numeric >= 0 )
);`
	queries[idx] = "createTableNotificationLevels"
	idx++

	queryMap["createTableCheckConfigurations"] = `
create table if not exists soma.check_configurations (
  configuration_id            uuid            PRIMARY KEY,
  repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
  bucket_id                   uuid            REFERENCES soma.buckets ( bucket_id ),
  configuration_name          varchar(256)    NOT NULL,
  configuration_object        uuid            NOT NULL,
  configuration_object_type   varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ),
  configuration_active        boolean         NOT NULL DEFAULT 'yes',
  inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
  children_only               boolean         NOT NULL DEFAULT 'no',
  capability_id               uuid            NOT NULL REFERENCES soma.monitoring_capabilities ( capability_id ),
  interval                    integer         NOT NULL,
  enabled                     boolean         NOT NULL DEFAULT 'yes',
  external_id                 varchar(64)     NOT NULL DEFAULT 'none',
  CHECK ( configuration_object_type != 'server' ),
  CHECK ( external_id = 'none' OR configuration_object_type != 'template' )
);`
	queries[idx] = "createTableCheckConfigurations"
	idx++

	queryMap["createTableCheckConfigurationThresholds"] = `
create table if not exists soma.configuration_thresholds (
  configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ),
  predicate                   varchar(4)      NOT NULL REFERENCES soma.configuration_predicates ( predicate ),
  threshold                   varchar(128)    NOT NULL,
  notification_level          varchar(16)     NOT NULL REFERENCES soma.notification_levels ( level_name )
);`
	queries[idx] = "createTableCheckConfigurationThresholds"
	idx++

	queryMap["createTableCheckConstraintsCustomProperty"] = `
create table if not exists soma.constraints_custom_property (
  configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ),
  custom_property_id          uuid            NOT NULL REFERENCES soma.custom_properties ( custom_property_id ),
  property_value              text            NOT NULL
);`
	queries[idx] = "createTableCheckConstraintsCustomProperty"
	idx++

	queryMap["createTableCheckConstraintsSystemProperty"] = `
create table if not exists soma.constraints_system_property (
  configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ),
  system_property             varchar(128)    NOT NULL REFERENCES soma.system_properties ( system_property ),
  property_value              text            NOT NULL
);`
	queries[idx] = "createTableCheckConstraintsSystemProperty"
	idx++

	queryMap["createTableCheckConstraintsNativeProperty"] = `
create table if not exists soma.constraints_native_property (
  configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ),
  native_property             varchar(128)    NOT NULL REFERENCES soma.native_properties ( native_property ),
  property_value              text            NOT NULL
);`
	queries[idx] = "createTableCheckConstraintsNativeProperty"
	idx++

	queryMap["createTableCheckConstraintsServiceProperty"] = `
create table if not exists soma.constraints_service_property (
  configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ),
  organizational_team_id      uuid            NOT NULL,
  service_property            varchar(64)     NOT NULL,
  FOREIGN KEY( organizational_team_id, service_property ) REFERENCES soma.team_service_properties ( organizational_team_id, service_property )
);`
	queries[idx] = "createTableCheckConstraintsServiceProperty"
	idx++

	queryMap["createTableCheckConstraintsServiceAttributes"] = `
create table if not exists soma.constraints_service_attribute (
  configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ),
  service_property_attribute  varchar(64)     NOT NULL REFERENCES soma.service_property_attributes ( service_property_attribute ),
  attribute_value             varchar(64)
);`
	queries[idx] = "createTableCheckConstraintsServiceAttributes"
	idx++

	queryMap["createTableCheckConstraintsOncallProperty"] = `
create table if not exists soma.constraints_oncall_property (
  configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ),
  oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_duty_teams ( oncall_duty_id )
);`
	queries[idx] = "createTableCheckConstraintsOncallProperty"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
