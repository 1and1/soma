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

func createTablesClusters(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 10)

	queryMap["createTableClusters"] = `
create table if not exists soma.clusters (
    cluster_id                  uuid            PRIMARY KEY,
    cluster_name                varchar(256)    NOT NULL,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    object_state                varchar(64)     NOT NULL DEFAULT 'standalone' REFERENCES soma.object_states ( object_state ) DEFERRABLE,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    -- enforce unique cluster names per bucket
    UNIQUE( bucket_id, cluster_name ),
    -- required for FK relation from soma.cluster_membership
    UNIQUE( bucket_id, cluster_id ),
    -- required for FK relation from soma.cluster_service_properties
    UNIQUE( cluster_id, organizational_team_id ),
    -- cluster must be configured like bucket it is in
    FOREIGN KEY ( bucket_id, organizational_team_id ) REFERENCES soma.buckets ( bucket_id, organizational_team_id ) DEFERRABLE
);`
	queries[idx] = "createTableClusters"
	idx++

	queryMap["createTableClusterMembership"] = `
create table if not exists soma.cluster_membership (
    cluster_id                  uuid            NOT NULL REFERENCES soma.clusters ( cluster_id ) DEFERRABLE,
    node_id                     uuid            NOT NULL REFERENCES soma.nodes ( node_id ) DEFERRABLE,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    -- every node can only be in one cluster
    UNIQUE( node_id ),
    -- node and cluster must belong to the same bucket
    FOREIGN KEY ( bucket_id, cluster_id ) REFERENCES soma.clusters ( bucket_id, cluster_id ) DEFERRABLE,
    FOREIGN KEY ( node_id, bucket_id ) REFERENCES soma.node_bucket_assignment ( node_id, bucket_id ) DEFERRABLE
);`
	queries[idx] = "createTableClusterMembership"
	idx++

	queryMap["createTableClusterOncallProperty"] = `
create table if not exists soma.cluster_oncall_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    cluster_id                  uuid            NOT NULL REFERENCES soma.clusters ( cluster_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_duty_teams ( oncall_duty_id ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    UNIQUE ( cluster_id, view ),
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableClusterOncallProperty"
	idx++

	queryMap["createTableClusterServiceProperties"] = `
create table if not exists soma.cluster_service_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    cluster_id                  uuid            NOT NULL REFERENCES clusters ( cluster_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES views ( view ) DEFERRABLE,
    service_property            varchar(64)     NOT NULL,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    UNIQUE( cluster_id, service_property, view ),
    FOREIGN KEY ( organizational_team_id, service_property ) REFERENCES soma.team_service_properties ( organizational_team_id, service_property ) DEFERRABLE,
    FOREIGN KEY ( cluster_id, organizational_team_id ) REFERENCES soma.clusters ( cluster_id, organizational_team_id ) DEFERRABLE,
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableClusterServiceProperties"
	idx++

	queryMap["createTableClusterSystemProperties"] = `
create table if not exists soma.cluster_system_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    cluster_id                  uuid            NOT NULL REFERENCES soma.clusters ( cluster_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    system_property             varchar(64)     NOT NULL REFERENCES soma.system_properties ( system_property ) DEFERRABLE,
    source_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    object_type                 varchar(64)     NOT NULL DEFAULT 'cluster' REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    inherited                   boolean         NOT NULL DEFAULT 'yes',
    value                       text            NOT NULL,
    FOREIGN KEY ( system_property, object_type, inherited ) REFERENCES soma.system_property_validity ( system_property, object_type, inherited ) DEFERRABLE,
    CHECK ( inherited OR object_type = 'cluster' ),
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableClusterSystemProperties"
	idx++

	// restrict all system properties to once per cluster+view, except
	// tags which would be silly if limited to once
	queryMap["createIndexUniqueClusterSystemProperties"] = `
create unique index _unique_cluster_system_properties
    on soma.cluster_system_properties ( cluster_id, system_property, view )
    where system_property != 'tag'
;`
	queries[idx] = "createIndexUniqueClusterSystemProperties"
	idx++

	queryMap["createTableClusterCustomProperties"] = `
create table if not exists soma.cluster_custom_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    cluster_id                  uuid            NOT NULL REFERENCES soma.clusters ( cluster_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    custom_property_id          uuid            NOT NULL REFERENCES soma.custom_properties ( custom_property_id ) DEFERRABLE,
    bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    UNIQUE ( cluster_id, custom_property_id, view ),
    -- ensure cluster is in this bucket, bucket is in this repository and custom_property is defined for this repository.
    -- together these three foreign keys link cluster_id with valid custom_property_id target
    FOREIGN KEY ( bucket_id, cluster_id ) REFERENCES soma.clusters ( bucket_id, cluster_id ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE,
    FOREIGN KEY ( repository_id, custom_property_id ) REFERENCES soma.custom_properties ( repository_id, custom_property_id ) DEFERRABLE,
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableClusterCustomProperties"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
