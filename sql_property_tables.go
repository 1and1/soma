package main

import (
  "log"
)

func sqlPropertyTables01() {
  var err error;

  /* Service Property
  */
  _, err = db.Exec(`create table if not exists soma.service_properties (
    service_property            varchar(64)     PRIMARY KEY
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.service_property_attributes (
    service_property_attribute  varchar(64)     PRIMARY KEY
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.service_property_values (
    service_property            varchar(64)     NOT NULL REFERENCES soma.service_properties ( service_property ),
    service_property_attribute  varchar(64)     NOT NULL REFERENCES soma.service_property_attributes ( service_property_attribute ),
    value                       varchar(64)     NOT NULL
  );`); if err != nil {
    log.Fatal( err )
  }

  /* Team Service Property
  */
  _, err = db.Exec(`create table if not exists soma.team_service_properties (
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
    service_property            varchar(64)     NOT NULL,
    UNIQUE( organizational_team_id, service_property )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.team_service_property_values (
    organizational_team_id      uuid            NOT NULL,
    service_property            varchar(64)     NOT NULL,
    service_property_attribute  varchar(64)     NOT NULL REFERENCES soma.service_property_attributes ( service_property_attribute ),
    value                       varchar(64)     NOT NULL,
    FOREIGN KEY( organizational_team_id, service_property ) REFERENCES soma.team_service_properties ( organizational_team_id, service_property )
  );`); if err != nil {
    log.Fatal( err )
  }

  /* System Property
  */
  _, err = db.Exec(`create table if not exists soma.system_properties (
    system_property             varchar(128)    PRIMARY KEY
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.system_property_validity (
    system_property             varchar(128)    NOT NULL REFERENCES soma.system_properties ( system_property ),
    object_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ),
    UNIQUE( system_property, object_type )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.native_properties (
    native_property             varchar(128)    PRIMARY KEY
  );`); if err != nil {
    log.Fatal( err )
  }
}

func sqlPropertyTables02() {
  var err error;

  _, err = db.Exec(`create table if not exists soma.custom_properties (
    custom_property_id          uuid            PRIMARY KEY,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    custom_property             varchar(128)    NOT NULL,
    UNIQUE( repository_id, custom_property ),
    UNIQUE( repository_id, custom_property_id )
  );`); if err != nil {
    log.Fatal( err )
  }
}
