package main

func createTablesNodes(printOnly bool, verbose bool) {
  idx := 0
  // map for storing the SQL statements by name
  queryMap := make( map[string]string )
  // slice storing the required statement order so foreign keys can
  // resolve successfully
  queries := make( []string, 10 )


  queryMap["createTableNodes"] = `
create table if not exists soma.nodes (
  node_id                     uuid            PRIMARY KEY,
  node_asset_id               numeric(16,0)   UNIQUE NOT NULL,
  node_name                   varchar(256)    NOT NULL,
  organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
  server_id                   uuid            NOT NULL REFERENCES inventory.servers ( server_id ),
  object_state                varchar(64)     NOT NULL DEFAULT 'standalone' REFERENCES soma.object_states ( object_state ),
  node_online                 boolean         NOT NULL DEFAULT 'yes',
  UNIQUE ( node_id, organizational_team_id )
);`
  queries[idx] = "createTableNodes"; idx++


  queryMap["createTableNodeBucketAssignment"] = `
create table if not exists soma.node_bucket_assignment (
  node_id                     uuid            NOT NULL,
  bucket_id                   uuid            NOT NULL,
  organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
  UNIQUE ( node_id ),
  UNIQUE ( node_id, bucket_id ),
  FOREIGN KEY ( node_id, organizational_team_id ) REFERENCES soma.nodes ( node_id, organizational_team_id ),
  FOREIGN KEY ( bucket_id, organizational_team_id ) REFERENCES soma.buckets ( bucket_id, organizational_team_id )
);`
  queries[idx] = "createTableNodeBucketAssignment"; idx++


  queryMap["createUniqueIndexNodeOnline"] = `
create unique index _unique_node_online
  on soma.nodes ( node_name, node_online )
  where node_online
;`
  queries[idx] = "createUniqueIndexNodeOnline"; idx++


  queryMap["createTableNodeOncallProperty"] = `
create table if not exists soma.node_oncall_property (
  node_id                     uuid            NOT NULL REFERENCES soma.nodes ( node_id ),
  view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
  oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_duty_teams ( oncall_duty_id ),
  inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
  children_only               boolean         NOT NULL DEFAULT 'no',
  UNIQUE ( node_id )
);`
  queries[idx] = "createTableNodeOncallProperty"; idx++


  queryMap["createTableNodeServiceProperties"] = `
create table if not exists soma.node_service_properties (
  node_id                     uuid            NOT NULL REFERENCES soma.nodes ( node_id ),
  view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
  service_property            varchar(64)     NOT NULL,
  organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
  inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
  children_only               boolean         NOT NULL DEFAULT 'no',
  UNIQUE ( node_id, service_property, view ),
  FOREIGN KEY ( organizational_team_id, service_property ) REFERENCES soma.team_service_properties ( organizational_team_id, service_property ),
  FOREIGN KEY ( node_id, organizational_team_id ) REFERENCES soma.nodes ( node_id, organizational_team_id )
);`
  queries[idx] = "createTableNodeServiceProperties"; idx++


  queryMap["createTableNodeSystemProperties"] = `
create table if not exists soma.node_system_properties (
  node_id                     uuid            NOT NULL REFERENCES nodes ( node_id ),
  view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES views ( view ),
  system_property             varchar(64)     NOT NULL REFERENCES soma.system_properties ( system_property ),
  object_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ),
  inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
  children_only               boolean         NOT NULL DEFAULT 'no',
  value                       text            NOT NULL,
  FOREIGN KEY ( system_property, object_type ) REFERENCES soma.system_property_validity ( system_property, object_type ),
  CHECK( object_type = 'node' )
);`
  queries[idx] = "createTableNodeSystemProperties"; idx++


  // restrict all system properties to once per cluster+view, except
  // tags which would be silly if limited to once
  queryMap["createUniqueIndexNodeSystemProperties"] = `
create unique index _unique_node_system_properties
  on soma.node_system_properties ( node_id, system_property, view )
  where system_property != 'tag'
;`
  queries[idx] = "createUniqueIndexNodeSystemProperties"; idx++


  queryMap["createTableNodeCustomProperties"] = `
create table if not exists soma.node_custom_properties (
  node_id                     uuid            NOT NULL REFERENCES soma.nodes ( node_id ),
  view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
  custom_property_id          uuid            NOT NULL,
  bucket_id                   uuid            NOT NULL,
  repository_id               uuid            NOT NULL,
  inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
  children_only               boolean         NOT NULL DEFAULT 'no',
  value                       text            NOT NULL,
  UNIQUE ( node_id, custom_property_id, view ),
  -- ensure node is in this bucket
  -- ensure bucket is in this repository
  -- ensure custom_property is defined for this repository
  FOREIGN KEY ( node_id, bucket_id ) REFERENCES soma.node_bucket_assignment ( node_id, bucket_id ),
  FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ),
  FOREIGN KEY ( repository_id, custom_property_id ) REFERENCES soma.custom_properties ( repository_id, custom_property_id )
);`
  queries[idx] = "createTableNodeCustomProperties"; idx++


  performDatabaseTask( printOnly, verbose, queries, queryMap )
}
