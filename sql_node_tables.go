package main

import (
  "log"
)

func sqlNodeTables01() {
  var err error;

  _, err = db.Exec(`create table if not exists soma.nodes (
    node_id                     uuid            PRIMARY KEY,
    node_asset_id               numeric(16,0)   UNIQUE NOT NULL,
    node_name                   varchar(256)    NOT NULL,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
    server_id                   uuid            NOT NULL REFERENCES inventory.servers ( server_id ),
    object_state                varchar(64)     NOT NULL DEFAULT 'standalone' REFERENCES soma.object_states ( object_state ),
    node_online                 boolean         NOT NULL DEFAULT 'yes',
    UNIQUE ( node_id, organizational_team_id )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.node_bucket_assignment (
    node_id                     uuid            NOT NULL,
    bucket_id                   uuid            NOT NULL,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
    UNIQUE ( node_id ),
    UNIQUE ( node_id, bucket_id ),
    FOREIGN KEY ( node_id, organizational_team_id ) REFERENCES soma.nodes ( node_id, organizational_team_id ),
    FOREIGN KEY ( bucket_id, organizational_team_id ) REFERENCES soma.buckets ( bucket_id, organizational_team_id )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create unique index _unique_node_online
    on soma.nodes ( node_name, node_online )
    where node_online
  ;`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.node_oncall_property (
    node_id                     uuid            NOT NULL REFERENCES soma.nodes ( node_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_duty_teams ( oncall_duty_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    UNIQUE ( node_id )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.node_service_properties (
    node_id                     uuid            NOT NULL REFERENCES soma.nodes ( node_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    service_property            varchar(64)     NOT NULL,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    UNIQUE ( node_id, service_property, view ),
    FOREIGN KEY ( organizational_team_id, service_property ) REFERENCES soma.team_service_properties ( organizational_team_id, service_property ),
    FOREIGN KEY ( node_id, organizational_team_id ) REFERENCES soma.nodes ( node_id, organizational_team_id )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.node_system_properties (
    node_id                     uuid            NOT NULL REFERENCES nodes ( node_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES views ( view ),
    system_property             varchar(64)     NOT NULL REFERENCES soma.system_properties ( system_property ),
    object_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    FOREIGN KEY ( system_property, object_type ) REFERENCES soma.system_property_validity ( system_property, object_type ),
    CHECK( object_type = 'node' )
  );`); if err != nil {
    log.Fatal( err )
  }

  // restrict all system properties to once per cluster+view, except
  // tags which would be silly if limited to once
  _, err = db.Exec(`create unique index _unique_node_system_properties
    on soma.node_system_properties ( node_id, system_property, view )
    where system_property != 'tag'
  ;`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.node_custom_properties (
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
  );`); if err != nil {
    log.Fatal( err )
  }
}
