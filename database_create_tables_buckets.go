package main

func createTablesBuckets(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 10)

	queryMap["createTableBuckets"] = `
create table if not exists soma.buckets (
    bucket_id                   uuid            PRIMARY KEY,
    bucket_name                 varchar(512)    UNIQUE NOT NULL,
    bucket_frozen               boolean         NOT NULL DEFAULT 'no',
    bucket_deleted              boolean         NOT NULL DEFAULT 'no',
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    environment                 varchar(32)     NOT NULL REFERENCES soma.environments ( environment ) DEFERRABLE,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    UNIQUE ( bucket_id, repository_id ),
    UNIQUE ( bucket_id, organizational_team_id )
);`
	queries[idx] = "createTableBuckets"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

func createTablesBucketsProperties(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 5)

	queryMap["createTableBucketOncall"] = `
create table if not exists soma.bucket_oncall_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_duty_teams ( oncall_duty_id ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableBucketOncall"
	idx++

	queryMap["createTableBucketService"] = `
create table if not exists soma.bucket_service_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    bucket_id                   uuid            NOT NULL REFERENCES buckets ( bucket_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES views ( view ) DEFERRABLE,
    service_property            varchar(64)     NOT NULL,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    FOREIGN KEY( organizational_team_id, service_property ) REFERENCES soma.team_service_properties ( organizational_team_id, service_property ),
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id )
);`
	queries[idx] = "createTableBucketService"
	idx++

	queryMap["createTableBucketSystem"] = `
create table if not exists soma.bucket_system_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    system_property             varchar(64)     NOT NULL REFERENCES soma.system_properties ( system_property ) DEFERRABLE,
    object_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    inherited                   boolean         NOT NULL DEFAULT 'yes',
    value                       text            NOT NULL,
    FOREIGN KEY ( system_property, object_type, inherited ) REFERENCES soma.system_property_validity ( system_property, object_type, inherited ) DEFERRABLE,
    CHECK( object_type = 'bucket' ),
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableBucketSystem"
	idx++

	queryMap["createTableBucketCustom"] = `
create table if not exists soma.bucket_custom_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    bucket_id                   uuid            NOT NULL,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    custom_property_id          uuid            NOT NULL,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE,
    FOREIGN KEY ( repository_id, custom_property_id ) REFERENCES soma.custom_properties ( repository_id, custom_property_id ) DEFERRABLE,
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableBucketCustom"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
