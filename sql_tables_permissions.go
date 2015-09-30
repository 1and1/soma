package main

import (
  "log"
)

func sqlPermissionTables01() {
  var err error;

  _, err = db.Exec(`create table if not exists soma.permission_types (
    permission_type             varchar(32)     PRIMARY KEY
    -- OMNIPOTENCE
    -- SYSTEM_GRANT
    -- SYSTEM
    -- GLOBAL
    -- REPO_GRANT
    -- REPO
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.permissions (
    permission_id               uuid            PRIMARY KEY,
    permission_name             varchar(128)    NOT NULL,
    permission_type             varchar(32)     NOT NULL REFERENCES soma.permission_types ( permission_type ),
    UNIQUE ( permission_name ),
    UNIQUE ( permission_id, permission_type )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.global_authorizations (
    admin_id                    uuid            REFERENCES auth.admins ( admin_id ),
    user_id                     uuid            REFERENCES inventory.users ( user_id ),
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ),
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ),
    permission_type             varchar(32)     NOT NULL REFERENCES soma.permission_types ( permission_type ),
    FOREIGN KEY ( permission_id, permission_type ) REFERENCES soma.permissions ( permission_id, permission_type ),
    CHECK (   ( admin_id IS NOT NULL AND user_id IS     NULL AND tool_id IS     NULL )
           OR ( admin_id IS     NULL AND user_id IS NOT NULL AND tool_id IS     NULL )
           OR ( admin_id IS     NULL AND user_id IS     NULL AND tool_id IS NOT NULL ) ),
    CHECK ( permission_type IN ( 'OMNIPOTENCE', 'SYSTEM_GRANT', 'SYSTEM', 'GLOBAL' ) )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.repo_authorizations (
    user_id                     uuid            REFERENCES inventory.users ( user_id ),
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ),
    organizational_team_id      uuid            REFERENCES inventory.organizational_teams ( organizational_team_id ),
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ),
    permission_type             varchar(32)     NOT NULL REFERENCES soma.permission_types ( permission_type ),
    FOREIGN KEY ( permission_id, permission_type ) REFERENCES soma.permissions ( permission_id, permission_type ),
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( permission_type IN ( 'REPO_GRANT', 'REPO' ) )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.bucket_authorizations (
    user_id                     uuid            REFERENCES inventory.users ( user_id ),
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ),
    organizational_team_id      uuid            REFERENCES inventory.organizational_teams ( organizational_team_id ),
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ),
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ),
    permission_type             varchar(32)     NOT NULL REFERENCES soma.permission_types ( permission_type ),
    FOREIGN KEY ( permission_id, permission_type ) REFERENCES soma.permissions ( permission_id, permission_type ),
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ),
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( permission_type = 'REPO' )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.group_authorizations (
    user_id                     uuid            REFERENCES inventory.users ( user_id ),
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ),
    organizational_team_id      uuid            REFERENCES inventory.organizational_teams ( organizational_team_id ),
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ),
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ),
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ),
    permission_type             varchar(32)     NOT NULL REFERENCES soma.permission_types ( permission_type ),
    FOREIGN KEY ( permission_id, permission_type ) REFERENCES soma.permissions ( permission_id, permission_type ),
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ),
    FOREIGN KEY ( bucket_id, group_id ) REFERENCES soma.groups ( bucket_id, group_id ),
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( permission_type = 'REPO' )
  );`); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec(`create table if not exists soma.cluster_authorizations (
    user_id                     uuid            REFERENCES inventory.users ( user_id ),
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ),
    organizational_team_id      uuid            REFERENCES inventory.organizational_teams ( organizational_team_id ),
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ),
    cluster_id                  uuid            NOT NULL REFERENCES soma.clusters ( cluster_id ),
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ),
    permission_type             varchar(32)     NOT NULL REFERENCES soma.permission_types ( permission_type ),
    FOREIGN KEY ( permission_id, permission_type ) REFERENCES soma.permissions ( permission_id, permission_type ),
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ),
    FOREIGN KEY ( bucket_id, cluster_id ) REFERENCES soma.clusters ( bucket_id, cluster_id ),
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( permission_type = 'REPO' )
  );`); if err != nil {
    log.Fatal( err )
  }
}
