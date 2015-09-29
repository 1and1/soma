package main

import (
  "log"
)

func sqlMetaTables01() {
  var err error;

  _, err = db.Exec(`create table if not exists soma.views (
    view                        varchar(64)     PRIMARY KEY
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.environments (
    environment                 varchar(32)     PRIMARY KEY
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.object_states (
    object_state                varchar(64)     PRIMARY KEY
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.object_types (
    object_type                 varchar(64)     PRIMARY KEY
  );`); if err != nil {
    log.Fatal( err )
  }
}

func sqlMetaTables02() {
  var err error;

  _, err = db.Exec(`create table if not exists soma.datacenter_groups (
    datacenter_group            varchar(32)     NOT NULL,
    datacenter                  varchar(32)     NOT NULL REFERENCES inventory.datacenters ( datacenter )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create index _datacenter_groups
    on soma.datacenter_groups ( datacenter_group )
  ;`); if err != nil {
    log.Fatal( err )
  }

}
