package main

import (
	"database/sql"
	"log"

	"github.com/satori/go.uuid"
)

func (tk *treeKeeper) startupRepositoryCustomProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                            error
		instanceId, srcInstanceId, repositoryId, view, customId, customProperty, value string
		inInstanceId, inObjectType, inObjId                                            string
		inheritance, childrenOnly                                                      bool
		rows, instance_rows                                                            *sql.Rows
		load_properties, load_instances                                                *sql.Stmt
	)
	load_properties, err = tk.conn.Prepare(`
SELECT srcp.instance_id,
       srcp.source_instance_id,
	   srcp.repository_id,
	   srcp.view,
	   srcp.custom_property_id,
	   srcp.inheritance_enabled,
	   srcp.children_only,
	   srcp.value,
	   scp.custom_property
FROM   soma.repository_custom_properties srcp
JOIN   soma.custom_properties scp
ON     srcp.custom_property_id = scp.custom_property_id
WHERE  srcp.instance_id = srcp.source_instance_id
AND    srcp.repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-repository-custom-properties: ", err)
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(tkStmtLoadCustomPropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-repository-custom-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading repository custom properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading repository custom properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

customloop:
	// load all custom properties defined directly on repository objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&repositoryId,
			&view,
			&customId,
			&inheritance,
			&childrenOnly,
			&value,
			&customProperty,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break customloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertyCustom{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          customProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.CustomId, _ = uuid.FromString(customId)
		prop.Instances = make([]tree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading repository custom properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current repository custom property so the IDs can be set correctly
		for instance_rows.Next() {
			err = instance_rows.Scan(
				&inInstanceId,
				&inObjectType,
				&inObjId,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break inproploop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			var propObjectId, propInstanceId uuid.UUID
			if propObjectId, err = uuid.FromString(inObjId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if uuid.Equal(uuid.Nil, propObjectId) || uuid.Equal(uuid.Nil, propInstanceId) {
				continue inproploop
			}
			if inObjectType == "MAGIC_NO_RESULT_VALUE" {
				continue inproploop
			}

			pi := tree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the repository and set the prepared property
		tk.tree.Find(tree.FindRequest{
			ElementType: `repository`,
			ElementId:   repositoryId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		tk.drain(`action`)
		tk.drain(`error`)
	}
}

func (tk *treeKeeper) startupBucketCustomProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                        error
		instanceId, srcInstanceId, bucketId, view, customId, customProperty, value string
		inInstanceId, inObjectType, inObjId                                        string
		inheritance, childrenOnly                                                  bool
		rows, instance_rows                                                        *sql.Rows
		load_properties, load_instances                                            *sql.Stmt
	)
	load_properties, err = tk.conn.Prepare(`
SELECT sbcp.instance_id,
       sbcp.source_instance_id,
	   sbcp.bucket_id,
	   sbcp.view,
	   sbcp.custom_property_id,
	   sbcp.inheritance_enabled,
	   sbcp.children_only,
	   sbcp.value,
	   scp.custom_property
FROM   soma.bucket_custom_properties sbcp
JOIN   soma.custom_properties scp
ON     sbcp.custom_property_id = scp.custom_property_id
WHERE  sbcp.instance_id = sbcp.source_instance_id
AND    sbcp.repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-bucket-custom-properties: ", err)
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(tkStmtLoadCustomPropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-bucket-custom-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading bucket custom properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading bucket custom properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

customloop:
	// load all custom properties defined directly on bucket objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&bucketId,
			&view,
			&customId,
			&inheritance,
			&childrenOnly,
			&value,
			&customProperty,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break customloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertyCustom{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          customProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.CustomId, _ = uuid.FromString(customId)
		prop.Instances = make([]tree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading bucket custom properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current bucket custom property so the IDs can be set correctly
		for instance_rows.Next() {
			err = instance_rows.Scan(
				&inInstanceId,
				&inObjectType,
				&inObjId,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break inproploop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			var propObjectId, propInstanceId uuid.UUID
			if propObjectId, err = uuid.FromString(inObjId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if uuid.Equal(uuid.Nil, propObjectId) || uuid.Equal(uuid.Nil, propInstanceId) {
				continue inproploop
			}
			if inObjectType == "MAGIC_NO_RESULT_VALUE" {
				continue inproploop
			}

			pi := tree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the bucket and set the prepared property
		tk.tree.Find(tree.FindRequest{
			ElementType: `bucket`,
			ElementId:   bucketId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		tk.drain(`action`)
		tk.drain(`error`)
	}
}

func (tk *treeKeeper) startupGroupCustomProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                       error
		instanceId, srcInstanceId, groupId, view, customId, customProperty, value string
		inInstanceId, inObjectType, inObjId                                       string
		inheritance, childrenOnly                                                 bool
		rows, instance_rows                                                       *sql.Rows
		load_properties, load_instances                                           *sql.Stmt
	)
	load_properties, err = tk.conn.Prepare(`
SELECT sgcp.instance_id,
       sgcp.source_instance_id,
	   sgcp.group_id,
	   sgcp.view,
	   sgcp.custom_property_id,
	   sgcp.inheritance_enabled,
	   sgcp.children_only,
	   sgcp.value,
	   scp.custom_property
FROM   soma.group_custom_properties sgcp
JOIN   soma.custom_properties scp
ON     sgcp.custom_property_id = scp.custom_property_id
WHERE  sgcp.instance_id = sgcp.source_instance_id
AND    sgcp.repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-group-custom-properties: ", err)
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(tkStmtLoadCustomPropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-group-custom-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading group custom properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading group custom properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

customloop:
	// load all system properties defined directly on group objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&groupId,
			&view,
			&customId,
			&inheritance,
			&childrenOnly,
			&value,
			&customProperty,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break customloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertyCustom{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          customProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.CustomId, _ = uuid.FromString(customId)
		prop.Instances = make([]tree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading group custom properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current group system property so the IDs can be set correctly
		for instance_rows.Next() {
			err = instance_rows.Scan(
				&inInstanceId,
				&inObjectType,
				&inObjId,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break inproploop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			var propObjectId, propInstanceId uuid.UUID
			if propObjectId, err = uuid.FromString(inObjId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if uuid.Equal(uuid.Nil, propObjectId) || uuid.Equal(uuid.Nil, propInstanceId) {
				continue inproploop
			}
			if inObjectType == "MAGIC_NO_RESULT_VALUE" {
				continue inproploop
			}

			pi := tree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the group and set the prepared property
		tk.tree.Find(tree.FindRequest{
			ElementId: groupId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		tk.drain(`action`)
		tk.drain(`error`)
	}
}

func (tk *treeKeeper) startupClusterCustomProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                         error
		instanceId, srcInstanceId, clusterId, view, customId, customProperty, value string
		inInstanceId, inObjectType, inObjId                                         string
		inheritance, childrenOnly                                                   bool
		rows, instance_rows                                                         *sql.Rows
		load_properties, load_instances                                             *sql.Stmt
	)
	load_properties, err = tk.conn.Prepare(`
SELECT sccp.instance_id,
       sccp.source_instance_id,
	   sccp.cluster_id,
	   sccp.view,
	   sccp.custom_property_id,
	   sccp.inheritance_enabled,
	   sccp.children_only,
	   sccp.value,
	   scp.custom_property
FROM   soma.cluster_custom_properties sccp
JOIN   soma.custom_properties scp
ON     sccp.custom_property_id = scp.custom_property_id
WHERE  sccp.instance_id = sccp.source_instance_id
AND    sccp.repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-cluster-custom-properties: ", err)
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(tkStmtLoadCustomPropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-cluster-custom-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading cluster custom properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading cluster custom properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

customloop:
	// load all custom properties defined directly on cluster objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&clusterId,
			&view,
			&customId,
			&inheritance,
			&childrenOnly,
			&value,
			&customProperty,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break customloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertyCustom{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          customProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.CustomId, _ = uuid.FromString(customId)
		prop.Instances = make([]tree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading cluster custom properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current cluster custom property so the IDs can be set correctly
		for instance_rows.Next() {
			err = instance_rows.Scan(
				&inInstanceId,
				&inObjectType,
				&inObjId,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break inproploop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			var propObjectId, propInstanceId uuid.UUID
			if propObjectId, err = uuid.FromString(inObjId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if uuid.Equal(uuid.Nil, propObjectId) || uuid.Equal(uuid.Nil, propInstanceId) {
				continue inproploop
			}
			if inObjectType == "MAGIC_NO_RESULT_VALUE" {
				continue inproploop
			}

			pi := tree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the cluster and set the prepared property
		tk.tree.Find(tree.FindRequest{
			ElementType: `cluster`,
			ElementId:   clusterId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		tk.drain(`action`)
		tk.drain(`error`)
	}
}

func (tk *treeKeeper) startupNodeCustomProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                      error
		instanceId, srcInstanceId, nodeId, view, customId, customProperty, value string
		inInstanceId, inObjectType, inObjId                                      string
		inheritance, childrenOnly                                                bool
		rows, instance_rows                                                      *sql.Rows
		load_properties, load_instances                                          *sql.Stmt
	)
	load_properties, err = tk.conn.Prepare(`
SELECT sncp.instance_id,
       sncp.source_instance_id,
	   sncp.node_id,
	   sncp.view,
	   sncp.custom_property_id,
	   sncp.inheritance_enabled,
	   sncp.children_only,
	   sncp.value,
	   scp.custom_property
FROM   soma.node_custom_properties sncp
JOIN   soma.custom_properties scp
ON     sncp.custom_property_id = scp.custom_property_id
WHERE  sncp.instance_id = sncp.source_instance_id
AND    sncp.repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-node-custom-properties: ", err)
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(tkStmtLoadCustomPropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-node-custom-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading node custom properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading node custom properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

customloop:
	// load all custom properties defined directly on node objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&nodeId,
			&view,
			&customId,
			&inheritance,
			&childrenOnly,
			&value,
			&customProperty,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break customloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertyCustom{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          customProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.CustomId, _ = uuid.FromString(customId)
		prop.Instances = make([]tree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading node custom properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current node custom property so the IDs can be set correctly
		for instance_rows.Next() {
			err = instance_rows.Scan(
				&inInstanceId,
				&inObjectType,
				&inObjId,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break inproploop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			var propObjectId, propInstanceId uuid.UUID
			if propObjectId, err = uuid.FromString(inObjId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if uuid.Equal(uuid.Nil, propObjectId) || uuid.Equal(uuid.Nil, propInstanceId) {
				continue inproploop
			}
			if inObjectType == "MAGIC_NO_RESULT_VALUE" {
				continue inproploop
			}

			pi := tree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the node and set the prepared property
		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementId:   nodeId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		tk.drain(`action`)
		tk.drain(`error`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
