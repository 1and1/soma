/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const TxDeleteCheckDetails = `
SELECT scc.configuration_object,
       scc.configuration_object_type,
       sc.source_check_id
FROM   soma.check_configurations scc
JOIN   soma.checks sc
  ON   scc.configuration_id = sc.configuration_id
WHERE  scc.configuration_id = $1::uuid
  AND  scc.repository_id    = $2::uuid
  AND  sc.check_id          = sc.source_check_id;`

const TxMarkCheckConfigDeleted = `
UPDATE soma.check_configurations
SET    deleted = 'yes'::boolean
WHERE  configuration_id = $1::uuid;`

const TxCreateCheck = `
INSERT INTO soma.checks (
            check_id,
            repository_id,
            bucket_id,
            source_check_id,
            source_object_type,
            source_object_id,
            configuration_id,
            capability_id,
            object_id,
            object_type)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::uuid,
       $5::varchar,
       $6::uuid,
       $7::uuid,
       $8::uuid,
       $9::uuid,
       $10::varchar;`

const TxMarkCheckDeleted = `
UPDATE soma.checks
SET    deleted = 'yes'::boolean
WHERE  check_id = $1::uuid;`

const TxCreateCheckInstance = `
INSERT INTO soma.check_instances (
            check_instance_id,
            check_id,
            check_configuration_id,
            current_instance_config_id,
            last_configuration_created)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::uuid,
       $5::timestamptz;`

const TxMarkCheckInstanceDeleted = `
UPDATE soma.check_instances
SET    deleted = 'yes'::boolean
WHERE  check_instance_id = $1::uuid;`

const TxCreateCheckInstanceConfiguration = `
INSERT INTO soma.check_instance_configurations (
            check_instance_config_id,
            version,
            check_instance_id,
            constraint_hash,
            constraint_val_hash,
            instance_service,
            instance_service_cfg_hash,
            instance_service_cfg,
            created,
            status,
            next_status,
            awaiting_deletion,
            deployment_details)
SELECT $1::uuid,
       $2::integer,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::varchar,
       $7::varchar,
       $8::jsonb,
       $9::timestamptz,
       $10::varchar,
       $11::varchar,
       $12::boolean,
       $13::jsonb;`

const TxCreateCheckConfigurationBase = `
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

const TxCreateCheckConfigurationThreshold = `
INSERT INTO soma.configuration_thresholds (
            configuration_id,
            predicate,
            threshold,
            notification_level)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar,
       $4::varchar;`

const TxCreateCheckConfigurationConstraintSystem = `
INSERT INTO soma.constraints_system_property (
            configuration_id,
            system_property,
            property_value)
SELECT $1::uuid,
       $2::varchar,
       $3::text;`

const TxCreateCheckConfigurationConstraintNative = `
INSERT INTO soma.constraints_native_property (
            configuration_id,
            native_property,
            property_value)
SELECT $1::uuid,
       $2::varchar,
       $3::text;`

const TxCreateCheckConfigurationConstraintOncall = `
INSERT INTO soma.constraints_oncall_property (
            configuration_id,
            oncall_duty_id)
SELECT $1::uuid,
       $2::uuid;`

const TxCreateCheckConfigurationConstraintCustom = `
INSERT INTO soma.constraints_custom_property (
            configuration_id,
            custom_property_id,
            repository_id,
            property_value)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::text;`

const TxCreateCheckConfigurationConstraintService = `
INSERT INTO soma.constraints_service_property (
            configuration_id,
            organizational_team_id,
            service_property)
SELECT $1::uuid,
       $2::uuid,
       $3::varchar;`

const TxCreateCheckConfigurationConstraintAttribute = `
INSERT INTO soma.constraints_service_attribute (
            configuration_id,
            service_property_attribute,
            attribute_value)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
