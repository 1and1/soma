package main

/*
 * Statements for CHECK actions
 */

const tkStmtCreateCheckConfigurationBase = `
INSERT INTO soma.check_configurations (
            configuration_id,
            configuration_name,
            interval,
            repository_id,
            bucket_id,
            capability_id,
            configuration_object,
            configuration_object_type,
            configuration_active,
            enabled,
            inheritance_enabled,
            children_only,
            external_id)
SELECT $1::uuid,
       $2::varchar,
       $3::integer,
       $4::uuid,
       $5::uuid,
       $6::uuid,
       $7::uuid,
       $8::varchar,
       $9::boolean,
       $10::boolean,
       $11::boolean,
       $12::boolean,
       $13::varchar;`

const tkStmtCreateCheckConfigurationThreshold = `
INSERT INTO soma.configuration_thresholds (
            configuration_id,
            predicate,
            threshold,
            notification_level)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar,
       $4::varchar;`

const tkStmtCreateCheckConfigurationConstraintSystem = `
INSERT INTO soma.constraints_system_property (
            configuration_id,
            system_property,
            property_value)
SELECT $1::uuid,
       $2::varchar,
       $3::text;`

const tkStmtCreateCheckConfigurationConstraintNative = `
INSERT INTO soma.constraints_native_property (
            configuration_id,
            native_property,
            property_value)
SELECT $1::uuid,
       $2::varchar,
       $3::text;`

const tkStmtCreateCheckConfigurationConstraintOncall = `
INSERT INTO soma.constraints_oncall_property (
            configuration_id,
            oncall_duty_id)
SELECT $1::uuid,
       $2::uuid;`

const tkStmtCreateCheckConfigurationConstraintCustom = `
INSERT INTO soma.constraints_custom_property (
            configuration_id,
            custom_property_id,
            repository_id,
            property_value)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::text;`

const tkStmtCreateCheckConfigurationConstraintService = `
INSERT INTO soma.constraints_service_property (
            configuration_id,
            organizational_team_id,
            service_property)
SELECT $1::uuid,
       $2::uuid,
       $3::varchar;`

const tkStmtCreateCheckConfigurationConstraintAttribute = `
INSERT INTO soma.constraints_service_attribute (
            configuration_id,
            service_property_attribute,
            attribute_value)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
