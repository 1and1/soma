package main

func createTablesGroups(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 15)

	queryMap["createTableGroups"] = `create table if not exists soma.groups (
    group_id                    uuid            PRIMARY KEY,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ),
    group_name                  varchar(256)    NOT NULL,
    object_state                varchar(64)     NOT NULL DEFAULT 'standalone' REFERENCES soma.object_states ( object_state ),
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
    -- enforce unique group names per bucket
    UNIQUE ( bucket_id, group_name ),
    -- group must be configured like bucket it is in
    FOREIGN KEY ( bucket_id, organizational_team_id ) REFERENCES soma.buckets ( bucket_id, organizational_team_id ),
    -- required for FK relations
    UNIQUE ( bucket_id, group_id ),
    UNIQUE ( group_id, organizational_team_id )
  );`
	queries[idx] = "createTableGroups"
	idx++

	queryMap["createTableGroupMembershipNodes"] = `create table if not exists soma.group_membership_nodes (
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ),
    child_node_id               uuid            NOT NULL REFERENCES soma.nodes ( node_id ),
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ),
    UNIQUE ( child_node_id ),
    -- node and group must belong to the same bucket
    FOREIGN KEY ( bucket_id, group_id ) REFERENCES soma.groups ( bucket_id, group_id ),
    FOREIGN KEY ( child_node_id, bucket_id ) REFERENCES soma.node_bucket_assignment ( node_id, bucket_id )
  );`
	queries[idx] = "createTableGroupMembershipNodes"
	idx++

	queryMap["createTableGroupMembershipClusters"] = `create table if not exists soma.group_membership_clusters (
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ),
    child_cluster_id            uuid            NOT NULL REFERENCES clusters ( cluster_id ),
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ),
    UNIQUE ( child_cluster_id ),
    -- cluster and group must belong to the same bucket
    FOREIGN KEY ( bucket_id, group_id ) REFERENCES soma.groups ( bucket_id, group_id ),
    FOREIGN KEY ( bucket_id, child_cluster_id ) REFERENCES soma.clusters ( bucket_id, cluster_id )
  );`
	queries[idx] = "createTableGroupMembershipClusters"
	idx++

	queryMap["createTableGroupMembershipGroups"] = `create table if not exists soma.group_membership_groups (
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ),
    child_group_id              uuid            NOT NULL REFERENCES soma.groups ( group_id ),
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ),
    UNIQUE ( child_group_id ),
    -- no fun for you, sir!
    CHECK ( group_id != child_group_id ),
    -- group and child_group must belong to the same bucket
    FOREIGN KEY ( bucket_id, group_id ) REFERENCES soma.groups ( bucket_id, group_id ),
    FOREIGN KEY ( bucket_id, child_group_id ) REFERENCES soma.groups ( bucket_id, group_id )
  );`
	queries[idx] = "createTableGroupMembershipGroups"
	idx++

	queryMap["createTableGroupOncallProperty"] = `create table if not exists soma.group_oncall_properties (
	instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ),
	source_instance_id          uuid            NOT NULL,
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_duty_teams ( oncall_duty_id ),
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    UNIQUE ( group_id, view ),
	FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id )
  );`
	queries[idx] = "createTableGroupOncallProperty"
	idx++

	queryMap["createTableGroupServiceProperties"] = `create table if not exists soma.group_service_properties (
	instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ),
	source_instance_id          uuid            NOT NULL,
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    service_property            varchar(64)     NOT NULL,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ),
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    UNIQUE( group_id, service_property, view ),
    FOREIGN KEY ( organizational_team_id, service_property ) REFERENCES soma.team_service_properties ( organizational_team_id, service_property ),
    FOREIGN KEY ( group_id, organizational_team_id ) REFERENCES soma.groups ( group_id, organizational_team_id ),
	FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id )
  );`
	queries[idx] = "createTableGroupServiceProperties"
	idx++

	queryMap["createTableGroupSystemProperties"] = `create table if not exists soma.group_system_properties (
	instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ),
	source_instance_id          uuid            NOT NULL,
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    system_property             varchar(64)     NOT NULL REFERENCES soma.system_properties ( system_property ),
    object_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ),
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    FOREIGN KEY ( system_property, object_type ) REFERENCES soma.system_property_validity ( system_property, object_type ),
    CHECK ( object_type = 'group' ),
	FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id )
  );`
	queries[idx] = "createTableGroupSystemProperties"
	idx++

	// restrict all system properties to once per group+view, except
	// tags which would be silly if limited to once
	queryMap["createIndexUniqueGroupSystemProperties"] = `create unique index _unique_group_system_properties
    on soma.group_system_properties ( group_id, system_property, view )
    where system_property != 'tag'
  ;`
	queries[idx] = "createIndexUniqueGroupSystemProperties"
	idx++

	queryMap["createTableGroupsCustomProperties"] = `create table if not exists soma.group_custom_properties (
	instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ),
	source_instance_id          uuid            NOT NULL,
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ),
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ),
    custom_property_id          uuid            NOT NULL REFERENCES soma.custom_properties ( custom_property_id ),
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ),
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    UNIQUE ( group_id, custom_property_id, view ),
    -- ensure group is in this bucket, bucket is in this repository and custom_property is defined for this repository.
    -- together these three foreign keys link group_id with valid custom_property_id target
    FOREIGN KEY ( bucket_id, group_id ) REFERENCES soma.groups ( bucket_id, group_id ),
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ),
    FOREIGN KEY ( repository_id, custom_property_id ) REFERENCES soma.custom_properties ( repository_id, custom_property_id ),
	FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id )
  );`
	queries[idx] = "createTableGroupsCustomProperties"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
