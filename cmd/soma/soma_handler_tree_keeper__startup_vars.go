package main

var tkStmtLoadSystemPropInstances = `
SELECT      CASE WHEN srsp.instance_id IS NOT NULL THEN srsp.instance_id
            ELSE CASE WHEN sbsp.instance_id IS NOT NULL THEN sbsp.instance_id
                 ELSE CASE WHEN sgsp.instance_id IS NOT NULL THEN sgsp.instance_id
                      ELSE CASE WHEN scsp.instance_id IS NOT NULL THEN scsp.instance_id
                           ELSE CASE WHEN snsp.instance_id IS NOT NULL THEN snsp.instance_id
                                ELSE '00000000-0000-0000-0000-000000000000'
                                END
                           END
                      END
                 END
            END AS "instance_id",
            CASE WHEN srsp.object_type IS NOT NULL THEN srsp.object_type
            ELSE CASE WHEN sbsp.object_type IS NOT NULL THEN sbsp.object_type
                 ELSE CASE WHEN sgsp.object_type IS NOT NULL THEN sgsp.object_type
                      ELSE CASE WHEN scsp.object_type IS NOT NULL THEN scsp.object_type
                           ELSE CASE WHEN snsp.object_type IS NOT NULL THEN snsp.object_type
                                ELSE 'MAGIC_NO_RESULT_VALUE'
                                END
                           END
                      END
                 END
            END AS "object_type",
            CASE WHEN srsp.repository_id IS NOT NULL THEN srsp.repository_id
            ELSE CASE WHEN sbsp.bucket_id IS NOT NULL THEN sbsp.bucket_id
                 ELSE CASE WHEN sgsp.group_id IS NOT NULL THEN sgsp.group_id
                      ELSE CASE WHEN scsp.cluster_id IS NOT NULL THEN scsp.cluster_id
                           ELSE CASE WHEN snsp.node_id IS NOT NULL THEN snsp.node_id
                                ELSE '00000000-0000-0000-0000-000000000000' 
                                END
                           END
                      END
                 END
            END AS "object_id"
FROM        soma.property_instances spi
LEFT JOIN   soma.repository_system_properties srsp
  ON        spi.instance_id = srsp.instance_id
  AND       spi.source_instance_id = srsp.source_instance_id
LEFT JOIN   soma.bucket_system_properties sbsp
  ON        spi.instance_id = sbsp.instance_id
  AND       spi.source_instance_id = sbsp.source_instance_id
LEFT JOIN   soma.group_system_properties sgsp
  ON        spi.instance_id = sgsp.instance_id
  AND       spi.source_instance_id = sgsp.source_instance_id
LEFT JOIN   soma.cluster_system_properties scsp
  ON        spi.instance_id = scsp.instance_id
  AND       spi.source_instance_id = scsp.source_instance_id
LEFT JOIN   soma.node_system_properties snsp
  ON        spi.instance_id = snsp.instance_id
  AND       spi.source_instance_id = snsp.source_instance_id
WHERE       spi.instance_id != spi.source_instance_id
  AND       spi.repository_id = $1::uuid
  AND       spi.source_instance_id = $2::uuid;`

var tkStmtLoadCustomPropInstances = `
SELECT      CASE WHEN srsp.instance_id IS NOT NULL THEN srsp.instance_id
            ELSE CASE WHEN sbsp.instance_id IS NOT NULL THEN sbsp.instance_id
                 ELSE CASE WHEN sgsp.instance_id IS NOT NULL THEN sgsp.instance_id
                      ELSE CASE WHEN scsp.instance_id IS NOT NULL THEN scsp.instance_id
                           ELSE CASE WHEN snsp.instance_id IS NOT NULL THEN snsp.instance_id
                                ELSE '00000000-0000-0000-0000-000000000000'
                                END
                           END
                      END
                 END
            END AS "instance_id",
            CASE WHEN srsp.instance_id IS NOT NULL THEN 'repository'
            ELSE CASE WHEN sbsp.instance_id IS NOT NULL THEN 'bucket'
                 ELSE CASE WHEN sgsp.instance_id IS NOT NULL THEN 'group'
                      ELSE CASE WHEN scsp.instance_id IS NOT NULL THEN 'cluster'
                           ELSE CASE WHEN snsp.instance_id IS NOT NULL THEN 'node'
                                ELSE 'MAGIC_NO_RESULT_VALUE'
                                END
                           END
                      END
                 END
            END AS "object_type",
            CASE WHEN srsp.repository_id IS NOT NULL THEN srsp.repository_id
            ELSE CASE WHEN sbsp.bucket_id IS NOT NULL THEN sbsp.bucket_id
                 ELSE CASE WHEN sgsp.group_id IS NOT NULL THEN sgsp.group_id
                      ELSE CASE WHEN scsp.cluster_id IS NOT NULL THEN scsp.cluster_id
                           ELSE CASE WHEN snsp.node_id IS NOT NULL THEN snsp.node_id
                                ELSE '00000000-0000-0000-0000-000000000000' 
                                END
                           END
                      END
                 END
            END AS "object_id"
FROM        soma.property_instances spi
LEFT JOIN   soma.repository_custom_properties srsp
  ON        spi.instance_id = srsp.instance_id
  AND       spi.source_instance_id = srsp.source_instance_id
LEFT JOIN   soma.bucket_custom_properties sbsp
  ON        spi.instance_id = sbsp.instance_id
  AND       spi.source_instance_id = sbsp.source_instance_id
LEFT JOIN   soma.group_custom_properties sgsp
  ON        spi.instance_id = sgsp.instance_id
  AND       spi.source_instance_id = sgsp.source_instance_id
LEFT JOIN   soma.cluster_custom_properties scsp
  ON        spi.instance_id = scsp.instance_id
  AND       spi.source_instance_id = scsp.source_instance_id
LEFT JOIN   soma.node_custom_properties snsp
  ON        spi.instance_id = snsp.instance_id
  AND       spi.source_instance_id = snsp.source_instance_id
WHERE       spi.instance_id != spi.source_instance_id
  AND       spi.repository_id = $1::uuid
  AND       spi.source_instance_id = $2::uuid;`

var tkStmtLoadServicePropInstances = `
SELECT      CASE WHEN srsp.instance_id IS NOT NULL THEN srsp.instance_id
            ELSE CASE WHEN sbsp.instance_id IS NOT NULL THEN sbsp.instance_id
                 ELSE CASE WHEN sgsp.instance_id IS NOT NULL THEN sgsp.instance_id
                      ELSE CASE WHEN scsp.instance_id IS NOT NULL THEN scsp.instance_id
                           ELSE CASE WHEN snsp.instance_id IS NOT NULL THEN snsp.instance_id
                                ELSE '00000000-0000-0000-0000-000000000000'
                                END
                           END
                      END
                 END
            END AS "instance_id",
            CASE WHEN srsp.instance_id IS NOT NULL THEN 'repository'
            ELSE CASE WHEN sbsp.instance_id IS NOT NULL THEN 'bucket'
                 ELSE CASE WHEN sgsp.instance_id IS NOT NULL THEN 'group'
                      ELSE CASE WHEN scsp.instance_id IS NOT NULL THEN 'cluster'
                           ELSE CASE WHEN snsp.instance_id IS NOT NULL THEN 'node'
                                ELSE 'MAGIC_NO_RESULT_VALUE'
                                END
                           END
                      END
                 END
            END AS "object_type",
            CASE WHEN srsp.repository_id IS NOT NULL THEN srsp.repository_id
            ELSE CASE WHEN sbsp.bucket_id IS NOT NULL THEN sbsp.bucket_id
                 ELSE CASE WHEN sgsp.group_id IS NOT NULL THEN sgsp.group_id
                      ELSE CASE WHEN scsp.cluster_id IS NOT NULL THEN scsp.cluster_id
                           ELSE CASE WHEN snsp.node_id IS NOT NULL THEN snsp.node_id
                                ELSE '00000000-0000-0000-0000-000000000000' 
                                END
                           END
                      END
                 END
            END AS "object_id"
FROM        soma.property_instances spi
LEFT JOIN   soma.repository_service_properties srsp
  ON        spi.instance_id = srsp.instance_id
  AND       spi.source_instance_id = srsp.source_instance_id
LEFT JOIN   soma.bucket_service_properties sbsp
  ON        spi.instance_id = sbsp.instance_id
  AND       spi.source_instance_id = sbsp.source_instance_id
LEFT JOIN   soma.group_service_properties sgsp
  ON        spi.instance_id = sgsp.instance_id
  AND       spi.source_instance_id = sgsp.source_instance_id
LEFT JOIN   soma.cluster_service_properties scsp
  ON        spi.instance_id = scsp.instance_id
  AND       spi.source_instance_id = scsp.source_instance_id
LEFT JOIN   soma.node_service_properties snsp
  ON        spi.instance_id = snsp.instance_id
  AND       spi.source_instance_id = snsp.source_instance_id
WHERE       spi.instance_id != spi.source_instance_id
  AND       spi.repository_id = $1::uuid
  AND       spi.source_instance_id = $2::uuid;`

var tkStmtLoadOncallPropInstances = `
SELECT      CASE WHEN srsp.instance_id IS NOT NULL THEN srsp.instance_id
            ELSE CASE WHEN sbsp.instance_id IS NOT NULL THEN sbsp.instance_id
                 ELSE CASE WHEN sgsp.instance_id IS NOT NULL THEN sgsp.instance_id
                      ELSE CASE WHEN scsp.instance_id IS NOT NULL THEN scsp.instance_id
                           ELSE CASE WHEN snsp.instance_id IS NOT NULL THEN snsp.instance_id
                                ELSE '00000000-0000-0000-0000-000000000000'
                                END
                           END
                      END
                 END
            END AS "instance_id",
            CASE WHEN srsp.instance_id IS NOT NULL THEN 'repository'
            ELSE CASE WHEN sbsp.instance_id IS NOT NULL THEN 'bucket'
                 ELSE CASE WHEN sgsp.instance_id IS NOT NULL THEN 'group'
                      ELSE CASE WHEN scsp.instance_id IS NOT NULL THEN 'cluster'
                           ELSE CASE WHEN snsp.instance_id IS NOT NULL THEN 'node'
                                ELSE 'MAGIC_NO_RESULT_VALUE'
                                END
                           END
                      END
                 END
            END AS "object_type",
            CASE WHEN srsp.repository_id IS NOT NULL THEN srsp.repository_id
            ELSE CASE WHEN sbsp.bucket_id IS NOT NULL THEN sbsp.bucket_id
                 ELSE CASE WHEN sgsp.group_id IS NOT NULL THEN sgsp.group_id
                      ELSE CASE WHEN scsp.cluster_id IS NOT NULL THEN scsp.cluster_id
                           ELSE CASE WHEN snsp.node_id IS NOT NULL THEN snsp.node_id
                                ELSE '00000000-0000-0000-0000-000000000000' 
                                END
                           END
                      END
                 END
            END AS "object_id"
FROM        soma.property_instances spi
LEFT JOIN   soma.repository_oncall_properties srsp
  ON        spi.instance_id = srsp.instance_id
  AND       spi.source_instance_id = srsp.source_instance_id
LEFT JOIN   soma.bucket_oncall_properties sbsp
  ON        spi.instance_id = sbsp.instance_id
  AND       spi.source_instance_id = sbsp.source_instance_id
LEFT JOIN   soma.group_oncall_properties sgsp
  ON        spi.instance_id = sgsp.instance_id
  AND       spi.source_instance_id = sgsp.source_instance_id
LEFT JOIN   soma.cluster_oncall_properties scsp
  ON        spi.instance_id = scsp.instance_id
  AND       spi.source_instance_id = scsp.source_instance_id
LEFT JOIN   soma.node_oncall_property snsp
  ON        spi.instance_id = snsp.instance_id
  AND       spi.source_instance_id = snsp.source_instance_id
WHERE       spi.instance_id != spi.source_instance_id
  AND       spi.repository_id = $1::uuid
  AND       spi.source_instance_id = $2::uuid;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
