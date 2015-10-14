package main

func createTablesMetricsMonitoring(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 10)

	queryMap["createTableMetricUnits"] = `
create table if not exists soma.metric_units (
  metric_unit                 varchar(8)      PRIMARY KEY,
  metric_unit_long_name       varchar(64)     NOT NULL
);`
	queries[idx] = "createTableMetricUnits"
	idx++

	queryMap["createTableMetrics"] = `
create table if not exists soma.metrics (
  metric                      varchar(512)    PRIMARY KEY,
  metric_unit                 varchar(8)      NOT NULL REFERENCES soma.metric_units ( metric_unit ),
  description                 text            NOT NULL
);`
	queries[idx] = "createTableMetrics"
	idx++

	queryMap["createTableMetricProviders"] = `
create table if not exists soma.metric_providers (
  metric_provider             varchar(64)     PRIMARY KEY
);`
	queries[idx] = "createTableMetricProviders"
	idx++

	queryMap["createTableMetricPackages"] = `
create table if not exists soma.metric_packages (
  metric                      varchar(512)    NOT NULL REFERENCES soma.metrics ( metric ),
  metric_provider             varchar(64)     NOT NULL REFERENCES soma.metric_providers ( metric_provider ),
  package                     varchar(128)    NOT NULL,
  UNIQUE ( metric, metric_provider )
);`
	queries[idx] = "createTableMetricPackages"
	idx++

	queryMap["createTableMonitoringSystemModes"] = `
create table if not exists soma.monitoring_system_modes (
  monitoring_system_mode      varchar(32)     PRIMARY KEY
);`
	queries[idx] = "createTableMonitoringSystemModes"
	idx++

	queryMap["createTableMonitoringSystems"] = `
create table if not exists soma.monitoring_systems (
  monitoring_id               uuid            PRIMARY KEY,
  monitoring_name             varchar(96)     UNIQUE NOT NULL,
  monitoring_system_mode      varchar(32)     NOT NULL REFERENCES soma.monitoring_system_modes ( monitoring_system_mode ),
  monitoring_contact          uuid            NOT NULL REFERENCES inventory.users ( user_id ),
  monitoring_owner_team       uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
  monitoring_callback_uri     text,
  UNIQUE ( monitoring_id, monitoring_system_mode )
);`
	queries[idx] = "createTableMonitoringSystems"
	idx++

	queryMap["createTableMonitoringSystemUsers"] = `
create table if not exists soma.monitoring_system_users (
  monitoring_id               uuid            NOT NULL REFERENCES soma.monitoring_systems ( monitoring_id ),
  monitoring_system_mode      varchar(32)     NOT NULL REFERENCES soma.monitoring_system_modes ( monitoring_system_mode ),
  organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
  FOREIGN KEY ( monitoring_id, monitoring_system_mode ) REFERENCES soma.monitoring_systems ( monitoring_id, monitoring_system_mode ),
  CHECK ( monitoring_system_mode = 'private' )
);`
	queries[idx] = "createTableMonitoringSystemUsers"
	idx++

	queryMap["createTableMonitoringCapabilities"] = `
create table if not exists soma.monitoring_capabilities (
  capability_id               uuid            PRIMARY KEY,
  capability_monitoring       uuid            NOT NULL REFERENCES soma.monitoring_systems ( monitoring_id ),
  capability_metric           varchar(512)    NOT NULL REFERENCES soma.metrics ( metric ),
  capability_view             varchar(64)     NOT NULL REFERENCES soma.views ( view ),
  CHECK ( capability_view != 'any' )
);`
	queries[idx] = "createTableMonitoringCapabilities"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
