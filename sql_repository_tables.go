package main

import (
  "log"
)

func sqlRepositoryTables01() {
  var err error;

  _, err = db.Exec(`create table if not exists soma.repositories (
    repository_id               uuid            PRIMARY KEY,
    repository_name             varchar(128)    UNIQUE NOT NULL,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
    UNIQUE( repository_id, organizational_team_id )
  );`); if err != nil {
    log.Fatal( err )
  }
}

func sqlRepositoryTables02() {
  var err error;

  _, err = db.Exec(`create table if not exists soma.repository_oncall_properties (
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_duty_teams ( oncall_duty_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no'
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.repository_service_properties (
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    service_property            varchar(64)     NOT NULL,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    FOREIGN KEY( organizational_team_id, service_property ) REFERENCES soma.team_service_properties ( organizational_team_id, service_property )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.repository_system_properties (
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    system_property             varchar(64)     NOT NULL REFERENCES soma.system_properties ( system_property ),
    object_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    FOREIGN KEY ( system_property, object_type ) REFERENCES soma.system_property_validity ( system_property, object_type ),
    CHECK( object_type = 'repository' )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.repository_custom_properties (
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    custom_property_id          uuid            NOT NULL REFERENCES soma.custom_properties ( custom_property_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    FOREIGN KEY ( repository_id, custom_property_id ) REFERENCES soma.custom_properties ( repository_id, custom_property_id )
  );`); if err != nil {
    log.Fatal( err )
  }
}
