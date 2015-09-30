package main

import (
  "log"
)

func sqlClusterTables01() {
  var err error;

  _, err = db.Exec(`create table if not exists soma.clusters (
    cluster_id                  uuid            PRIMARY KEY,
    cluster_name                varchar(256)    NOT NULL,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ),
    object_state                varchar(64)     NOT NULL DEFAULT 'standalone' REFERENCES soma.object_states ( object_state ),
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
    -- enforce unique cluster names per bucket
    UNIQUE( bucket_id, cluster_name ),
    -- required for FK relation from soma.cluster_membership
    UNIQUE( bucket_id, cluster_id ),
    -- required for FK relation from soma.cluster_service_properties
    UNIQUE( cluster_id, organizational_team_id ),
    -- cluster must be configured like bucket it is in
    FOREIGN KEY ( bucket_id, organizational_team_id ) REFERENCES soma.buckets ( bucket_id, organizational_team_id )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.cluster_membership (
    cluster_id                  uuid            NOT NULL REFERENCES soma.clusters ( cluster_id ),
    node_id                     uuid            NOT NULL REFERENCES soma.nodes ( node_id ),
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ),
    -- every node can only be in one cluster
    UNIQUE( node_id ),
    -- node and cluster must belong to the same bucket
    FOREIGN KEY ( bucket_id, cluster_id ) REFERENCES soma.clusters ( bucket_id, cluster_id ),
    FOREIGN KEY ( node_id, bucket_id ) REFERENCES soma.node_bucket_assignment ( node_id, bucket_id )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.cluster_oncall_properties (
    cluster_id                  uuid            NOT NULL REFERENCES soma.clusters ( cluster_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_duty_teams ( oncall_duty_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    UNIQUE ( cluster_id, view )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.cluster_service_properties (
    cluster_id                  uuid            NOT NULL REFERENCES clusters ( cluster_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES views ( view ),
    service_property            varchar(64)     NOT NULL,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    UNIQUE( cluster_id, service_property, view ),
    FOREIGN KEY ( organizational_team_id, service_property ) REFERENCES soma.team_service_properties ( organizational_team_id, service_property ),
    FOREIGN KEY ( cluster_id, organizational_team_id ) REFERENCES soma.clusters ( cluster_id, organizational_team_id )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.cluster_system_properties (
    cluster_id                  uuid            NOT NULL REFERENCES soma.clusters ( cluster_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    system_property             varchar(64)     NOT NULL REFERENCES soma.system_properties ( system_property ),
    object_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    FOREIGN KEY ( system_property, object_type ) REFERENCES soma.system_property_validity ( system_property, object_type ),
    CHECK ( object_type = 'cluster' )
  );`); if err != nil {
    log.Fatal( err )
  }

  // restrict all system properties to once per cluster+view, except
  // tags which would be silly if limited to once
  _, err = db.Exec(`create unique index _unique_cluster_system_properties
    on soma.cluster_system_properties ( cluster_id, system_property, view )
    where system_property != 'tag'
  ;`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.cluster_custom_properties (
    cluster_id                  uuid            NOT NULL REFERENCES soma.clusters ( cluster_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    custom_property_id          uuid            NOT NULL REFERENCES soma.custom_properties ( custom_property_id ),
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ),
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    UNIQUE ( cluster_id, custom_property_id, view ),
    -- ensure cluster is in this bucket, bucket is in this repository and custom_property is defined for this repository.
    -- together these three foreign keys link cluster_id with valid custom_property_id target
    FOREIGN KEY ( bucket_id, cluster_id ) REFERENCES soma.clusters ( bucket_id, cluster_id ),
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ),
    FOREIGN KEY ( repository_id, custom_property_id ) REFERENCES soma.custom_properties ( repository_id, custom_property_id )
  );`); if err != nil {
    log.Fatal( err )
  }
}
