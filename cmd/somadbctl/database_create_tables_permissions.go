package main

func createTablesPermissions(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 25)

	queryMap["createTableCategories"] = `
create table if not exists soma.categories (
    category                    varchar(32)     PRIMARY KEY,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW()
);`
	queries[idx] = "createTablePermissionTypes"
	idx++

	queryMap["createTableSections"] = `
create table if not exists soma.sections (
    section_id                  uuid            PRIMARY KEY,
    section_name                varchar(64)     UNIQUE NOT NULL,
    category                    varchar(32)     NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    UNIQUE ( section_id, category )
);`
	queries[idx] = "createTableSections"
	idx++

	queryMap["createTableActions"] = `
create table if not exists soma.actions (
    action_id                   uuid            PRIMARY KEY,
    action_name                 varchar(64)     NOT NULL,
    section_id                  uuid            NOT NULL REFERENCES soma.sections ( section_id ) DEFERRABLE,
    category                    varchar(32)     NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    --
    UNIQUE ( section_id, action_name ),
    UNIQUE ( section_id, action_id ),
    FOREIGN KEY ( section_id, category ) REFERENCES soma.sections ( section_id, category ) DEFERRABLE
);`
	queries[idx] = "createTableActions"
	idx++

	queryMap["createTableSomaPermissions"] = `
create table if not exists soma.permissions (
    permission_id               uuid            PRIMARY KEY,
    permission_name             varchar(128)    NOT NULL,
    category                    varchar(32)     NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    UNIQUE ( permission_name, category ),
    UNIQUE ( permission_id, category ),
    -- only omnipotence is category omnipotence
    CHECK  ( category != 'omnipotence' OR permission_name = 'omnipotence' )
);`
	queries[idx] = "createTableSomaPermissions"
	idx++

	queryMap["createTablePermissionMap"] = `
create table if not exists soma.permission_map (
    mapping_id                  uuid            PRIMARY KEY,
    category                    varchar(32)     NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE,
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE,
    section_id                  uuid            NOT NULL REFERENCES soma.sections ( section_id ) DEFERRABLE,
    action_id                   uuid            REFERENCES soma.actions ( action_id ) DEFERRABLE,
    FOREIGN KEY ( permission_id, category )     REFERENCES soma.permissions ( permission_id, category ),
    FOREIGN KEY ( section_id, category )        REFERENCES soma.sections ( section_id, category ),
    FOREIGN KEY ( section_id, action_id )       REFERENCES soma.actions ( section_id, action_id )
);`
	queries[idx] = "createTableSomaPermissionMap"
	idx++

	queryMap[`createTablePermissionGrantMap`] = `
create table if not exists soma.permission_grant_map (
    category                    varchar(32)     NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE,
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE,
    granted_category            varchar(32)     NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE,
    granted_permission_id       uuid            NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE,
    FOREIGN KEY ( permission_id, category ) REFERENCES soma.permissions ( permission_id, category ),
    FOREIGN KEY ( granted_permission_id, granted_category ) REFERENCES soma.permissions ( permission_id, category ),
    CHECK ( permission_id != granted_permission_id ),
    CHECK ( category ~ ':grant$' ),
    CHECK ( granted_category = substring(category from '^([^:]+):'))
);`
	queries[idx] = `createTablePermissionGrantMap`
	idx++

	queryMap["createTableGlobalAuthorizations"] = `
create table if not exists soma.authorizations_global (
    grant_id                    uuid            PRIMARY KEY,
    admin_id                    uuid            REFERENCES auth.admins ( admin_id ) DEFERRABLE,
    user_id                     uuid            REFERENCES inventory.users ( user_id ) DEFERRABLE,
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ) DEFERRABLE,
    organizational_team_id      uuid            REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE,
    category                    varchar(32)     NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, category ) REFERENCES soma.permissions ( permission_id, category ) DEFERRABLE,
    CHECK (   ( admin_id IS NOT NULL AND user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( admin_id IS     NULL AND user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( admin_id IS     NULL AND user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( admin_id IS     NULL AND user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( category IN ( 'omnipotence','system','global','global:grant','permission','permission:grant','operations','operations:grant' ) ),
    -- only admin accounts can have system permissions
    CHECK ( admin_id IS NULL OR category != 'system' ),
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
    category                    varchar(32)     NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, category ) REFERENCES soma.permissions ( permission_id, category ) DEFERRABLE,
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( category IN ( 'repository', 'repository:grant' ) )
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
    category                    varchar(32)     NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, category ) REFERENCES soma.permissions ( permission_id, category ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE,
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( category = 'repository' )
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
    category                    varchar(32)     NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, category ) REFERENCES soma.permissions ( permission_id, category ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, group_id ) REFERENCES soma.groups ( bucket_id, group_id ) DEFERRABLE,
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( category = 'repository' )
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
    category                    varchar(32)     NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, category ) REFERENCES soma.permissions ( permission_id, category ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, cluster_id ) REFERENCES soma.clusters ( bucket_id, cluster_id ) DEFERRABLE,
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( category = 'repository' )
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
    category                    varchar(32)     NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, category ) REFERENCES soma.permissions ( permission_id, category ) DEFERRABLE,
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( category IN ( 'monitoring', 'monitoring:grant' ) )
);`
	queries[idx] = "createTableMonitoringAuthorizations"
	idx++

	queryMap["createTableTeamAuthorizations"] = `
create table if not exists soma.authorizations_team (
    grant_id                    uuid            PRIMARY KEY,
    user_id                     uuid            REFERENCES inventory.users ( user_id ) DEFERRABLE,
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ) DEFERRABLE,
    organizational_team_id      uuid            REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    authorized_team_id          uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    permission_id               uuid            NOT NULL REFERENCES soma.permissions ( permission_id ) DEFERRABLE,
    category                    varchar(32)     NOT NULL REFERENCES soma.categories ( category ) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, category ) REFERENCES soma.permissions ( permission_id, category ) DEFERRABLE,
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND organizational_team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND organizational_team_id IS NOT NULL ) ),
    CHECK ( category = 'team', 'team:grant' )
);`
	queries[idx] = "createTableTeamAuthorizations"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
