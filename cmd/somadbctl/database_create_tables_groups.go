/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright notice,
 *   this list of conditions and the following disclaimer in the documentation
 *   and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package main

func createTablesGroups(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 15)

	queryMap["createTableGroups"] = `
create table if not exists soma.groups (
    group_id                    uuid            PRIMARY KEY,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    group_name                  varchar(256)    NOT NULL,
    object_state                varchar(64)     NOT NULL DEFAULT 'standalone' REFERENCES soma.object_states ( object_state ) DEFERRABLE,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    created_by                  uuid            NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000' REFERENCES inventory.users ( user_id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    -- enforce unique group names per bucket
    UNIQUE ( bucket_id, group_name ),
    -- group must be configured like bucket it is in
    FOREIGN KEY ( bucket_id, organizational_team_id ) REFERENCES soma.buckets ( bucket_id, organizational_team_id ) DEFERRABLE,
    -- required for FK relations
    UNIQUE ( bucket_id, group_id ),
    UNIQUE ( group_id, organizational_team_id )
);`
	queries[idx] = "createTableGroups"
	idx++

	queryMap["createTableGroupMembershipNodes"] = `
create table if not exists soma.group_membership_nodes (
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ) DEFERRABLE,
    child_node_id               uuid            NOT NULL REFERENCES soma.nodes ( node_id ) DEFERRABLE,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    UNIQUE ( child_node_id ),
    -- node and group must belong to the same bucket
    FOREIGN KEY ( bucket_id, group_id ) REFERENCES soma.groups ( bucket_id, group_id ) DEFERRABLE,
    FOREIGN KEY ( child_node_id, bucket_id ) REFERENCES soma.node_bucket_assignment ( node_id, bucket_id ) DEFERRABLE
);`
	queries[idx] = "createTableGroupMembershipNodes"
	idx++

	queryMap["createTableGroupMembershipClusters"] = `
create table if not exists soma.group_membership_clusters (
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ) DEFERRABLE,
    child_cluster_id            uuid            NOT NULL REFERENCES clusters ( cluster_id ) DEFERRABLE,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    UNIQUE ( child_cluster_id ),
    -- cluster and group must belong to the same bucket
    FOREIGN KEY ( bucket_id, group_id ) REFERENCES soma.groups ( bucket_id, group_id ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, child_cluster_id ) REFERENCES soma.clusters ( bucket_id, cluster_id ) DEFERRABLE
);`
	queries[idx] = "createTableGroupMembershipClusters"
	idx++

	queryMap["createTableGroupMembershipGroups"] = `
create table if not exists soma.group_membership_groups (
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ) DEFERRABLE,
    child_group_id              uuid            NOT NULL REFERENCES soma.groups ( group_id ) DEFERRABLE,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    UNIQUE ( child_group_id ),
    -- no fun for you, sir!
    CHECK ( group_id != child_group_id ),
    -- group and child_group must belong to the same bucket
    FOREIGN KEY ( bucket_id, group_id ) REFERENCES soma.groups ( bucket_id, group_id ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, child_group_id ) REFERENCES soma.groups ( bucket_id, group_id ) DEFERRABLE
);`
	queries[idx] = "createTableGroupMembershipGroups"
	idx++

	queryMap["createTableGroupOncallProperty"] = `
create table if not exists soma.group_oncall_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_duty_teams ( oncall_duty_id ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    UNIQUE ( group_id, view ),
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableGroupOncallProperty"
	idx++

	queryMap["createTableGroupServiceProperties"] = `
create table if not exists soma.group_service_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    service_property            varchar(64)     NOT NULL,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    UNIQUE( group_id, service_property, view ),
    FOREIGN KEY ( organizational_team_id, service_property ) REFERENCES soma.team_service_properties ( organizational_team_id, service_property ) DEFERRABLE,
    FOREIGN KEY ( group_id, organizational_team_id ) REFERENCES soma.groups ( group_id, organizational_team_id ) DEFERRABLE,
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableGroupServiceProperties"
	idx++

	queryMap["createTableGroupSystemProperties"] = `
create table if not exists soma.group_system_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    system_property             varchar(64)     NOT NULL REFERENCES soma.system_properties ( system_property ) DEFERRABLE,
    source_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    object_type                 varchar(64)     NOT NULL DEFAULT 'group' REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    inherited                   boolean         NOT NULL DEFAULT 'yes',
    value                       text            NOT NULL,
    FOREIGN KEY ( system_property, object_type, inherited ) REFERENCES soma.system_property_validity ( system_property, object_type, inherited ) DEFERRABLE,
    CHECK ( inherited OR object_type = 'group' ),
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableGroupSystemProperties"
	idx++

	// restrict all system properties to once per group+view, except
	// tags which would be silly if limited to once
	queryMap["createIndexUniqueGroupSystemProperties"] = `
create unique index _unique_group_system_properties
    on soma.group_system_properties ( group_id, system_property, view )
    where system_property != 'tag'
;`
	queries[idx] = "createIndexUniqueGroupSystemProperties"
	idx++

	queryMap["createTableGroupsCustomProperties"] = `
create table if not exists soma.group_custom_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    group_id                    uuid            NOT NULL REFERENCES soma.groups ( group_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    custom_property_id          uuid            NOT NULL REFERENCES soma.custom_properties ( custom_property_id ) DEFERRABLE,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    UNIQUE ( group_id, custom_property_id, view ),
    -- ensure group is in this bucket, bucket is in this repository and custom_property is defined for this repository.
    -- together these three foreign keys link group_id with valid custom_property_id target
    FOREIGN KEY ( bucket_id, group_id ) REFERENCES soma.groups ( bucket_id, group_id ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE,
    FOREIGN KEY ( repository_id, custom_property_id ) REFERENCES soma.custom_properties ( repository_id, custom_property_id ) DEFERRABLE,
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableGroupsCustomProperties"

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
