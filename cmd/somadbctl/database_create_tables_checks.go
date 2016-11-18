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

func createTablesChecks(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 27)

	queryMap["createTableConfigurationPredicates"] = `
create table if not exists soma.configuration_predicates (
    predicate                   varchar(4)      PRIMARY KEY
);`
	queries[idx] = "createTableConfigurationPredicates"
	idx++

	queryMap["createTableNotificationLevels"] = `
create table if not exists soma.notification_levels (
    level_name                  varchar(16)     PRIMARY KEY,
    level_shortname             varchar(16)     UNIQUE NOT NULL,
    level_numeric               smallint        UNIQUE NOT NULL,
    CHECK ( level_numeric >= 0 )
);`
	queries[idx] = "createTableNotificationLevels"
	idx++

	queryMap["createTableCheckConfigurations"] = `
create table if not exists soma.check_configurations (
    configuration_id            uuid            PRIMARY KEY,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    bucket_id                   uuid            REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    configuration_name          varchar(256)    NOT NULL,
    configuration_object        uuid            NOT NULL,
    configuration_object_type   varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    configuration_active        boolean         NOT NULL DEFAULT 'no',
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    capability_id               uuid            NOT NULL REFERENCES soma.monitoring_capabilities ( capability_id ) DEFERRABLE,
    interval                    integer         NOT NULL,
    enabled                     boolean         NOT NULL DEFAULT 'yes',
    external_id                 varchar(64)     NOT NULL DEFAULT 'none',
    deleted                     boolean         NOT NULL DEFAULT 'no',
    -- required for custom property constraint foreign key
    UNIQUE ( configuration_id, repository_id ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE,
    CHECK ( configuration_object_type != 'server' ),
    CHECK ( external_id = 'none' OR configuration_object_type != 'template' ),
    CHECK ( interval > 0 )
);`
	queries[idx] = "createTableCheckConfigurations"
	idx++

	queryMap[`createSingletonCheckNameIndex`] = `
create unique index _singleton_undeleted_checkconfig_name
    on soma.check_configurations (
        repository_id,
        configuration_name)
    WHERE NOT deleted;`
	queries[idx] = `createSingletonCheckNameIndex`
	idx++

	queryMap["createTableChecks"] = `
create table if not exists soma.checks (
    check_id                    uuid            PRIMARY KEY,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    bucket_id                   uuid            REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    source_check_id             uuid            NOT NULL,
    source_object_type          varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    source_object_id            uuid            NOT NULL,
    configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ) DEFERRABLE,
    capability_id               uuid            NOT NULL REFERENCES soma.monitoring_capabilities ( capability_id ) DEFERRABLE,
    object_id                   uuid            NOT NULL,
    object_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    deleted                     boolean         NOT NULL DEFAULT 'no',
    UNIQUE( check_id, repository_id ),
    FOREIGN KEY ( source_check_id, repository_id ) REFERENCES soma.checks ( check_id, repository_id ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableChecks"
	idx++

	queryMap[`createIndexRepoToChecks`] = `
create index _repo_to_checks
    on soma.checks (
    repository_id,
    check_id
);`
	queries[idx] = `createIndexRepoToChecks`
	idx++

	queryMap["createTableCheckConfigurationThresholds"] = `
create table if not exists soma.configuration_thresholds (
    configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ) DEFERRABLE,
    predicate                   varchar(4)      NOT NULL REFERENCES soma.configuration_predicates ( predicate ) DEFERRABLE,
    threshold                   varchar(128)    NOT NULL,
    notification_level          varchar(16)     NOT NULL REFERENCES soma.notification_levels ( level_name ) DEFERRABLE
);`
	queries[idx] = "createTableCheckConfigurationThresholds"
	idx++

	queryMap[`createIndexConfigurationLevelsById`] = `
create index _configuration_id_levels
    on soma.configuration_thresholds (
    configuration_id,
    notification_level
);`
	queries[idx] = `createIndexConfigurationLevelsById`
	idx++

	queryMap["createTableCheckConstraintsCustomProperty"] = `
create table if not exists soma.constraints_custom_property (
    configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ) DEFERRABLE,
    custom_property_id          uuid            NOT NULL REFERENCES soma.custom_properties ( custom_property_id ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    property_value              text            NOT NULL,
    -- ensure this custom property is defined for this repository
    FOREIGN KEY ( repository_id, custom_property_id ) REFERENCES soma.custom_properties ( repository_id, custom_property_id ) DEFERRABLE,
    -- ensure the configuration_id is for the repository the custom property is defined in
    FOREIGN KEY ( configuration_id, repository_id ) REFERENCES soma.check_configurations ( configuration_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableCheckConstraintsCustomProperty"
	idx++

	queryMap["createTableCheckConstraintsSystemProperty"] = `
create table if not exists soma.constraints_system_property (
    configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ) DEFERRABLE,
    system_property             varchar(128)    NOT NULL REFERENCES soma.system_properties ( system_property ) DEFERRABLE,
    property_value              text            NOT NULL
);`
	queries[idx] = "createTableCheckConstraintsSystemProperty"
	idx++

	queryMap["createTableCheckConstraintsNativeProperty"] = `
create table if not exists soma.constraints_native_property (
    configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ) DEFERRABLE,
    native_property             varchar(128)    NOT NULL REFERENCES soma.native_properties ( native_property ) DEFERRABLE,
    property_value              text            NOT NULL
);`
	queries[idx] = "createTableCheckConstraintsNativeProperty"
	idx++

	queryMap["createTableCheckConstraintsServiceProperty"] = `
create table if not exists soma.constraints_service_property (
    configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ) DEFERRABLE,
    organizational_team_id      uuid            NOT NULL,
    service_property            varchar(128)    NOT NULL,
    FOREIGN KEY( organizational_team_id, service_property ) REFERENCES soma.team_service_properties ( organizational_team_id, service_property ) DEFERRABLE
);`
	queries[idx] = "createTableCheckConstraintsServiceProperty"
	idx++

	queryMap["createTableCheckConstraintsServiceAttributes"] = `
create table if not exists soma.constraints_service_attribute (
    configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ) DEFERRABLE,
    service_property_attribute  varchar(128)    NOT NULL REFERENCES soma.service_property_attributes ( service_property_attribute ) DEFERRABLE,
    attribute_value             varchar(512)
);`
	queries[idx] = "createTableCheckConstraintsServiceAttributes"
	idx++

	queryMap["createTableCheckConstraintsOncallProperty"] = `
create table if not exists soma.constraints_oncall_property (
    configuration_id            uuid            NOT NULL REFERENCES soma.check_configurations ( configuration_id ) DEFERRABLE,
    oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_duty_teams ( oncall_duty_id ) DEFERRABLE
);`
	queries[idx] = "createTableCheckConstraintsOncallProperty"

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
