package main

import (
  "log"
)

func sqlAuthTables01() {
  var err error;

  _, err = db.Exec(`create table if not exists auth.user_authentication (
    user_id                     uuid            NOT NULL REFERENCES inventory.users ( user_id ) ON DELETE CASCADE,
    algorithm                   varchar(16)     NOT NULL,
    rounds                      numeric(10,0)   NOT NULL,
    salt                        text            NOT NULL,
    crypt                       text            NOT NULL,
    reset_pending               boolean         NOT NULL DEFAULT 'no',
    valid_from                  timestamptz(3)  NOT NULL,
    valid_until                 timestamptz(3)  NOT NULL,
    CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
    CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists auth.user_token_authentication (
    user_id                     uuid            NOT NULL REFERENCES inventory.users ( user_id ) ON DELETE CASCADE,
    token                       varchar(256)    UNIQUE NOT NULL,
    valid_from                  timestamptz(3)  NOT NULL,
    valid_until                 timestamptz(3)  NOT NULL,
    CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
    CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists auth.admins (
    admin_id                    uuid            PRIMARY KEY,
    admin_uid                   varchar(256)    UNIQUE NOT NULL,
    admin_user_id               varchar(256)    NOT NULL REFERENCES inventory.users ( user_uid ) ON DELETE CASCADE,
    admin_is_active             boolean         NOT NULL DEFAULT 'yes',
    CHECK( left( admin_uid, 6 ) = 'admin_' ),
    CHECK( position( admin_user_id in admin_uid ) != 0 )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists auth.admin_authentication (
    admin_id                    uuid            NOT NULL REFERENCES auth.admins ( admin_id ) ON DELETE CASCADE,
    algorithm                   varchar(16)     NOT NULL,
    rounds                      numeric(10,0)   NOT NULL,
    salt                        text            NOT NULL,
    crypt                       text            NOT NULL,
    reset_pending               boolean         NOT NULL DEFAULT 'no',
    valid_from                  timestamptz(3)  NOT NULL,
    valid_until                 timestamptz(3)  NOT NULL,
    CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
    CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists auth.tools (
    tool_name                   varchar(256)    PRIMARY KEY,
    tool_owner_id               uuid            NOT NULL REFERENCES inventory.users ( user_id ) ON DELETE RESTRICT,
    created                     timestamptz(3)  NOT NULL DEFAULT NOW(),
    CHECK( EXTRACT( TIMEZONE FROM created ) = '0' )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists auth.tool_authentication (
    tool_name                   varchar(256)    NOT NULL REFERENCES auth.tools ( tool_name ) ON UPDATE CASCADE ON DELETE CASCADE,
    algorithm                   varchar(16)     NOT NULL,
    rounds                      numeric(10,0)   NOT NULL,
    salt                        text            NOT NULL,
    crypt                       text            NOT NULL,
    reset_pending               boolean         NOT NULL DEFAULT 'no',
    valid_from                  timestamptz(3)  NOT NULL,
    valid_until                 timestamptz(3)  NOT NULL,
    CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
    CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists auth.tool_token_authentication (
    tool_name                   varchar(256)    NOT NULL REFERENCES auth.tools ( tool_name ) ON UPDATE CASCADE ON DELETE CASCADE,
    token                       varchar(256)    UNIQUE NOT NULL,
    valid_from                  timestamptz(3)  NOT NULL,
    valid_until                 timestamptz(3)  NOT NULL,
    CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
    CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists auth.password_reset (
    user_id                     uuid            NULL REFERENCES inventory.users ( user_id ) ON DELETE CASCADE,
    tool_name                   varchar(256)    NULL REFERENCES auth.tools ( tool_name ) ON UPDATE CASCADE ON DELETE CASCADE,
    token                       varchar(256)    UNIQUE NOT NULL,
    valid_from                  timestamptz(3)  NOT NULL,
    valid_until                 timestamptz(3)  NOT NULL,
    token_used                  boolean         NOT NULL DEFAULT 'no',
    token_invalidated           boolean         NOT NULL DEFAULT 'no',
    CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
    CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' ),
    CHECK(    ( user_id IS NOT NULL AND tool_name IS     NULL )
           OR ( user_id IS     NULL AND tool_name IS NOT NULL ) )
  );`); if err != nil {
    log.Fatal( err )
  }
}
