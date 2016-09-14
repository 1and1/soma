package main

func createTablesPermissions(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 25)

	queryMap["createTablePermissionTypes"] = `
create table if not exists soma.permission_types (
    permission_type             varchar(32)     PRIMARY KEY,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW()
    -- omnipotence
    -- grant_system
    -- system
    -- global
    -- grant_limited
    -- limited
);`
	queries[idx] = "createTablePermissionTypes"
	idx++

	queryMap["createTableSomaPermissions"] = `
create table if not exists soma.permissions (
    permission_id               uuid            PRIMARY KEY,
    permission_name             varchar(128)    NOT NULL,
    permission_type             varchar(32)     NOT NULL REFERENCES soma.permission_types ( permission_type ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    UNIQUE ( permission_name ),
    UNIQUE ( permission_id, permission_type ),
    -- only omnipotence is type omnipotence
    CHECK  ( permission_type != 'omnipotence' OR permission_name = 'omnipotence' )
);`
	queries[idx] = "createTableSomaPermissions"
	idx++

	queryMap["createTableGlobalAuthorizations"] = `
create table if not exists soma.authorizations_global (
    grant_id                    uuid            PRIMARY KEY,
    admin_id                    uuid            REFERENCES auth.admins ( admin_id ) DEFERRABLE,
    user_id                     uuid            REFERENCES inventory.users ( user_id ) DEFERRABLE,
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ) DEFERRABLE,
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE,
    permission_type             varchar(32)     NOT NULL REFERENCES soma.permission_types ( permission_type ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, permission_type ) REFERENCES soma.permissions ( permission_id, permission_type ) DEFERRABLE,
    CHECK (   ( admin_id IS NOT NULL AND user_id IS     NULL AND tool_id IS     NULL )
           OR ( admin_id IS     NULL AND user_id IS NOT NULL AND tool_id IS     NULL )
           OR ( admin_id IS     NULL AND user_id IS     NULL AND tool_id IS NOT NULL ) ),
    CHECK ( permission_type IN ( 'omnipotence', 'grant_system', 'system', 'global' ) ),
    -- only root can have omnipotence
    CHECK ( permission_id != '00000000-0000-0000-0000-000000000000'::uuid OR user_id = '00000000-0000-0000-0000-000000000000'::uuid )
);`
	queries[idx] = "createTableGlobalAuthorizations"
	idx++

	queryMap[`createUniqueIndexAdminGlobalAuthorization`] = `
create unique index _unique_admin_global_authoriz
    on soma.authorizations_global ( admin_id, permission_id )
    where admin_id IS NOT NULL;`
	queries[idx] = `createUniqueIndexAdminGlobalAuthorization`
	idx++

	queryMap[`createUniqueIndexUserGlobalAuthorization`] = `
create unique index _unique_user_global_authoriz
    on soma.authorizations_global ( user_id, permission_id )
    where user_id IS NOT NULL;`
	queries[idx] = `createUniqueIndexUserGlobalAuthorization`
	idx++

	queryMap[`createUniqueIndexToolGlobalAuthorization`] = `
create unique index _unique_tool_global_authoriz
    on soma.authorizations_global ( tool_id, permission_id )
    where tool_id IS NOT NULL;`
	queries[idx] = `createUniqueIndexToolGlobalAuthorization`
	idx++

	queryMap["createTableRepoAuthorizations"] = `
create table if not exists soma.authorizations_repository (
    grant_id                    uuid            PRIMARY KEY,
    user_id                     uuid            REFERENCES inventory.users ( user_id ) DEFERRABLE,
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ) DEFERRABLE,
    organizational_team_id      uuid            REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE,
    permission_type             varchar(32)     NOT NULL REFERENCES soma.permission_types ( permission_type ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, permission_type ) REFERENCES soma.permissions ( permission_id, permission_type ) DEFERRABLE,
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( permission_type IN ( 'grant_limited', 'limited' ) )
);`
	queries[idx] = "createTableRepoAuthorizations"
	idx++

	queryMap["createTableBucketAuthorizations"] = `
create table if not exists soma.authorizations_bucket (
    grant_id                    uuid            PRIMARY KEY,
    user_id                     uuid            REFERENCES inventory.users ( user_id ) DEFERRABLE,
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ) DEFERRABLE,
    organizational_team_id      uuid            REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE,
    permission_type             varchar(32)     NOT NULL REFERENCES soma.permission_types ( permission_type ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, permission_type ) REFERENCES soma.permissions ( permission_id, permission_type ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE,
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( permission_type = 'limited' )
);`
	queries[idx] = "createTableBucketAuthorizations"
	idx++

	queryMap["createTableGroupAuthorizations"] = `
create table if not exists soma.authorizations_group (
    grant_id                    uuid            PRIMARY KEY,
    user_id                     uuid            REFERENCES inventory.users ( user_id ) DEFERRABLE,
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ) DEFERRABLE,
    organizational_team_id      uuid            REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ) DEFERRABLE,
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE,
    permission_type             varchar(32)     NOT NULL REFERENCES soma.permission_types ( permission_type ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, permission_type ) REFERENCES soma.permissions ( permission_id, permission_type ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, group_id ) REFERENCES soma.groups ( bucket_id, group_id ) DEFERRABLE,
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( permission_type = 'limited' )
);`
	queries[idx] = "createTableGroupAuthorizations"
	idx++

	queryMap["createTableClusterAuthorizations"] = `
create table if not exists soma.authorizations_cluster (
    grant_id                    uuid            PRIMARY KEY,
    user_id                     uuid            REFERENCES inventory.users ( user_id ) DEFERRABLE,
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ) DEFERRABLE,
    organizational_team_id      uuid            REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    cluster_id                  uuid            NOT NULL REFERENCES soma.clusters ( cluster_id ) DEFERRABLE,
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE,
    permission_type             varchar(32)     NOT NULL REFERENCES soma.permission_types ( permission_type ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, permission_type ) REFERENCES soma.permissions ( permission_id, permission_type ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, cluster_id ) REFERENCES soma.clusters ( bucket_id, cluster_id ) DEFERRABLE,
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( permission_type = 'limited' )
);`
	queries[idx] = "createTableClusterAuthorizations"
	idx++

	queryMap["createTableMonitoringAuthorizations"] = `
create table if not exists soma.authorizations_monitoring (
    grant_id                    uuid            PRIMARY KEY,
    user_id                     uuid            REFERENCES inventory.users ( user_id ) DEFERRABLE,
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ) DEFERRABLE,
    organizational_team_id      uuid            REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    monitoring_id               uuid            NOT NULL REFERENCES soma.monitoring_systems ( monitoring_id ) DEFERRABLE,
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE,
    permission_type             varchar(32)     NOT NULL REFERENCES soma.permission_types ( permission_type ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, permission_type ) REFERENCES soma.permissions ( permission_id, permission_type ) DEFERRABLE,
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( permission_type = 'limited' )
);`
	queries[idx] = "createTableMonitoringAuthorizations"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
