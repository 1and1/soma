package main

import (
  "log"
)

func sqlBucketsTables01() {
  var err error;

  _, err = db.Exec(`create table if not exists soma.buckets (
    bucket_id                   uuid            PRIMARY KEY,
    bucket_name                 varchar(512)    UNIQUE NOT NULL,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    environment                 varchar(32)     NOT NULL REFERENCES soma.environments ( environment ),
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
    UNIQUE ( bucket_id, repository_id ),
    UNIQUE ( bucket_id, organizational_team_id )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.bucket_oncall_properties (
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_duty_teams ( oncall_duty_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no'
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.bucket_service_properties (
    bucket_id                   uuid            NOT NULL REFERENCES buckets ( bucket_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES views ( view ),
    service_property            varchar(64)     NOT NULL,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    FOREIGN KEY( organizational_team_id, service_property ) REFERENCES soma.team_service_properties ( organizational_team_id, service_property )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.bucket_system_properties (
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    system_property             varchar(64)     NOT NULL REFERENCES soma.system_properties ( system_property ),
    object_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    FOREIGN KEY ( system_property, object_type ) REFERENCES soma.system_property_validity ( system_property, object_type ),
    CHECK( object_type = 'bucket' )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.bucket_custom_properties (
    bucket_id                   uuid            NOT NULL,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    custom_property_id          uuid            NOT NULL,
    repository_id               uuid            NOT NULL,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ),
    FOREIGN KEY ( repository_id, custom_property_id ) REFERENCES soma.custom_properties ( repository_id, custom_property_id )
  );`); if err != nil {
    log.Fatal( err )
  }
}
