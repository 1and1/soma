package main

import (
	"database/sql"
	"log"

	"github.com/satori/go.uuid"
)

func (tk *treeKeeper) startupRepositoryServiceProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                    error
		instanceId, srcInstanceId, repositoryId, view, serviceProperty, teamId string
		inInstanceId, inObjectType, inObjId, attrKey, attrValue                string
		inheritance, childrenOnly                                              bool
		rows, attribute_rows, instance_rows                                    *sql.Rows
		load_properties, load_attributes, load_instances                       *sql.Stmt
	)
	log.Println("Prepare: treekeeper/load-repository-service-properties")
	load_properties, err = tk.conn.Prepare(`
SELECT instance_id,
       source_instance_id,
	   repository_id,
	   view,
	   service_property,
	   organizational_team_id,
	   inheritance_enabled,
	   children_only
FROM   soma.repository_service_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-repository-service-properties: ", err)
	}
	defer load_properties.Close()

	log.Println("Prepare: treekeeper/load-repository-service-property-attributes")
	load_attributes, err = tk.conn.Prepare(`
SELECT service_property_attribute,
       value
FROM   soma.team_service_property_values
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`)
	if err != nil {
		log.Fatal("treekeeper/load-repository-service-property-attributes: ", err)
	}
	defer load_attributes.Close()

	log.Println("Prepare: treekeeper/load-repository-service-property-instances")
	load_instances, err = tk.conn.Prepare(tkStmtLoadServicePropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-repository-service-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading repository service properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading repository custom properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

serviceloop:
	// load all service properties defined directly on repository objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&repositoryId,
			&view,
			&serviceProperty,
			&teamId,
			&inheritance,
			&childrenOnly,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break serviceloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := somatree.PropertyService{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Service:      serviceProperty,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Attributes = make([]somaproto.TreeServiceAttribute, 0)
		prop.Instances = make([]somatree.PropertyInstance, 0)

		attribute_rows, err = load_attributes.Query(
			teamId,
			serviceProperty,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading repository service properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer attribute_rows.Close()

	attributeloop:
		// load service attributes
		for attribute_rows.Next() {
			err = attribute_rows.Scan(
				&attrKey,
				&attrValue,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break attributeloop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			pa := somaproto.TreeServiceAttribute{
				Attribute: attrKey,
				Value:     attrValue,
			}
			prop.Attributes = append(prop.Attributes, pa)
		}

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading repository service properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current repository service property so the IDs can be set correctly
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

			pi := somatree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the repository and set the prepared property
		tk.tree.Find(somatree.FindRequest{
			ElementId: repositoryId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		for i := len(tk.actionChan); i > 0; i-- {
			//a := <-tk.actionChan
			//log.Printf("%s -> %s\n", a.Action, a.Type)
			<-tk.actionChan
		}
		for i := len(tk.errChan); i > 0; i-- {
			<-tk.errChan
		}
	}
}

func (tk *treeKeeper) startupBucketServiceProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                error
		instanceId, srcInstanceId, bucketId, view, serviceProperty, teamId string
		inInstanceId, inObjectType, inObjId, attrKey, attrValue            string
		inheritance, childrenOnly                                          bool
		rows, attribute_rows, instance_rows                                *sql.Rows
		load_properties, load_attributes, load_instances                   *sql.Stmt
	)
	log.Println("Prepare: treekeeper/load-bucket-service-properties")
	load_properties, err = tk.conn.Prepare(`
SELECT instance_id,
       source_instance_id,
	   bucket_id,
	   view,
	   service_property,
	   organizational_team_id,
	   inheritance_enabled,
	   children_only
FROM   soma.bucket_service_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-bucket-service-properties: ", err)
	}
	defer load_properties.Close()

	log.Println("Prepare: treekeeper/load-bucket-service-property-attributes")
	load_attributes, err = tk.conn.Prepare(`
SELECT service_property_attribute,
       value
FROM   soma.team_service_property_values
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`)
	if err != nil {
		log.Fatal("treekeeper/load-bucket-service-property-attributes: ", err)
	}
	defer load_attributes.Close()

	log.Println("Prepare: treekeeper/load-bucket-service-property-instances")
	load_instances, err = tk.conn.Prepare(tkStmtLoadServicePropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-bucket-service-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading bucket service properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading bucket custom properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

serviceloop:
	// load all service properties defined directly on bucket objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&bucketId,
			&view,
			&serviceProperty,
			&teamId,
			&inheritance,
			&childrenOnly,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break serviceloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := somatree.PropertyService{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Service:      serviceProperty,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Attributes = make([]somaproto.TreeServiceAttribute, 0)
		prop.Instances = make([]somatree.PropertyInstance, 0)

		attribute_rows, err = load_attributes.Query(
			teamId,
			serviceProperty,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading bucket service properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer attribute_rows.Close()

	attributeloop:
		// load service attributes
		for attribute_rows.Next() {
			err = attribute_rows.Scan(
				&attrKey,
				&attrValue,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break attributeloop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			pa := somaproto.TreeServiceAttribute{
				Attribute: attrKey,
				Value:     attrValue,
			}
			prop.Attributes = append(prop.Attributes, pa)
		}

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading bucket service properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current bucket service property so the IDs can be set correctly
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

			pi := somatree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the bucket and set the prepared property
		tk.tree.Find(somatree.FindRequest{
			ElementId: bucketId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		for i := len(tk.actionChan); i > 0; i-- {
			<-tk.actionChan
			//a := <-tk.actionChan
			//log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := len(tk.errChan); i > 0; i-- {
			<-tk.errChan
		}
	}
}

func (tk *treeKeeper) startupGroupServiceProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                               error
		instanceId, srcInstanceId, groupId, view, serviceProperty, teamId string
		inInstanceId, inObjectType, inObjId, attrKey, attrValue           string
		inheritance, childrenOnly                                         bool
		rows, attribute_rows, instance_rows                               *sql.Rows
		load_properties, load_attributes, load_instances                  *sql.Stmt
	)
	log.Println("Prepare: treekeeper/load-group-service-properties")
	load_properties, err = tk.conn.Prepare(`
SELECT instance_id,
       source_instance_id,
	   group_id,
	   view,
	   service_property,
	   organizational_team_id,
	   inheritance_enabled,
	   children_only
FROM   soma.group_service_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-group-service-properties: ", err)
	}
	defer load_properties.Close()

	log.Println("Prepare: treekeeper/load-group-service-property-attributes")
	load_attributes, err = tk.conn.Prepare(`
SELECT service_property_attribute,
       value
FROM   soma.team_service_property_values
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`)
	if err != nil {
		log.Fatal("treekeeper/load-group-service-property-attributes: ", err)
	}
	defer load_attributes.Close()

	log.Println("Prepare: treekeeper/load-group-service-property-instances")
	load_instances, err = tk.conn.Prepare(tkStmtLoadServicePropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-group-service-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading group service properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading group custom properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

serviceloop:
	// load all service properties defined directly on group objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&groupId,
			&view,
			&serviceProperty,
			&teamId,
			&inheritance,
			&childrenOnly,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break serviceloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := somatree.PropertyService{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Service:      serviceProperty,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Attributes = make([]somaproto.TreeServiceAttribute, 0)
		prop.Instances = make([]somatree.PropertyInstance, 0)

		attribute_rows, err = load_attributes.Query(
			teamId,
			serviceProperty,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading group service properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer attribute_rows.Close()

	attributeloop:
		// load service attributes
		for attribute_rows.Next() {
			err = attribute_rows.Scan(
				&attrKey,
				&attrValue,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break attributeloop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			pa := somaproto.TreeServiceAttribute{
				Attribute: attrKey,
				Value:     attrValue,
			}
			prop.Attributes = append(prop.Attributes, pa)
		}

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading group service properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current group service property so the IDs can be set correctly
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

			pi := somatree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the group and set the prepared property
		tk.tree.Find(somatree.FindRequest{
			ElementId: groupId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		for i := len(tk.actionChan); i > 0; i-- {
			<-tk.actionChan
			//a := <-tk.actionChan
			//log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := len(tk.errChan); i > 0; i-- {
			<-tk.errChan
		}
	}
}

func (tk *treeKeeper) startupClusterServiceProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                 error
		instanceId, srcInstanceId, clusterId, view, serviceProperty, teamId string
		inInstanceId, inObjectType, inObjId, attrKey, attrValue             string
		inheritance, childrenOnly                                           bool
		rows, attribute_rows, instance_rows                                 *sql.Rows
		load_properties, load_attributes, load_instances                    *sql.Stmt
	)
	log.Println("Prepare: treekeeper/load-cluster-service-properties")
	load_properties, err = tk.conn.Prepare(`
SELECT instance_id,
       source_instance_id,
	   cluster_id,
	   view,
	   service_property,
	   organizational_team_id,
	   inheritance_enabled,
	   children_only
FROM   soma.cluster_service_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-cluster-service-properties: ", err)
	}
	defer load_properties.Close()

	log.Println("Prepare: treekeeper/load-cluster-service-property-attributes")
	load_attributes, err = tk.conn.Prepare(`
SELECT service_property_attribute,
       value
FROM   soma.team_service_property_values
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`)
	if err != nil {
		log.Fatal("treekeeper/load-cluster-service-property-attributes: ", err)
	}
	defer load_attributes.Close()

	log.Println("Prepare: treekeeper/load-cluster-service-property-instances")
	load_instances, err = tk.conn.Prepare(tkStmtLoadServicePropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-cluster-service-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading cluster service properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading cluster custom properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

serviceloop:
	// load all service properties defined directly on cluster objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&clusterId,
			&view,
			&serviceProperty,
			&teamId,
			&inheritance,
			&childrenOnly,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break serviceloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := somatree.PropertyService{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Service:      serviceProperty,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Attributes = make([]somaproto.TreeServiceAttribute, 0)
		prop.Instances = make([]somatree.PropertyInstance, 0)

		attribute_rows, err = load_attributes.Query(
			teamId,
			serviceProperty,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading cluster service properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer attribute_rows.Close()

	attributeloop:
		// load service attributes
		for attribute_rows.Next() {
			err = attribute_rows.Scan(
				&attrKey,
				&attrValue,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break attributeloop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			pa := somaproto.TreeServiceAttribute{
				Attribute: attrKey,
				Value:     attrValue,
			}
			prop.Attributes = append(prop.Attributes, pa)
		}

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading cluster service properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current cluster service property so the IDs can be set correctly
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

			pi := somatree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the cluster and set the prepared property
		tk.tree.Find(somatree.FindRequest{
			ElementId: clusterId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		for i := len(tk.actionChan); i > 0; i-- {
			<-tk.actionChan
			//a := <-tk.actionChan
			//log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := len(tk.errChan); i > 0; i-- {
			<-tk.errChan
		}
	}
}

func (tk *treeKeeper) startupNodeServiceProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                              error
		instanceId, srcInstanceId, nodeId, view, serviceProperty, teamId string
		inInstanceId, inObjectType, inObjId, attrKey, attrValue          string
		inheritance, childrenOnly                                        bool
		rows, attribute_rows, instance_rows                              *sql.Rows
		load_properties, load_attributes, load_instances                 *sql.Stmt
	)
	log.Println("Prepare: treekeeper/load-node-service-properties")
	load_properties, err = tk.conn.Prepare(`
SELECT instance_id,
       source_instance_id,
	   node_id,
	   view,
	   service_property,
	   organizational_team_id,
	   inheritance_enabled,
	   children_only
FROM   soma.node_service_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-node-service-properties: ", err)
	}
	defer load_properties.Close()

	log.Println("Prepare: treekeeper/load-node-service-property-attributes")
	load_attributes, err = tk.conn.Prepare(`
SELECT service_property_attribute,
       value
FROM   soma.team_service_property_values
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`)
	if err != nil {
		log.Fatal("treekeeper/load-node-service-property-attributes: ", err)
	}
	defer load_attributes.Close()

	log.Println("Prepare: treekeeper/load-node-service-property-instances")
	load_instances, err = tk.conn.Prepare(tkStmtLoadServicePropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-node-service-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading node service properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading node service properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

serviceloop:
	// load all service properties defined directly on node objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&nodeId,
			&view,
			&serviceProperty,
			&teamId,
			&inheritance,
			&childrenOnly,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break serviceloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := somatree.PropertyService{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Service:      serviceProperty,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Attributes = make([]somaproto.TreeServiceAttribute, 0)
		prop.Instances = make([]somatree.PropertyInstance, 0)

		attribute_rows, err = load_attributes.Query(
			teamId,
			serviceProperty,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading node service properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer attribute_rows.Close()

	attributeloop:
		// load service attributes
		for attribute_rows.Next() {
			err = attribute_rows.Scan(
				&attrKey,
				&attrValue,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break attributeloop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			pa := somaproto.TreeServiceAttribute{
				Attribute: attrKey,
				Value:     attrValue,
			}
			prop.Attributes = append(prop.Attributes, pa)
		}

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading node service properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current node service property so the IDs can be set correctly
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

			pi := somatree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the node and set the prepared property
		tk.tree.Find(somatree.FindRequest{
			ElementId: nodeId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		for i := len(tk.actionChan); i > 0; i-- {
			<-tk.actionChan
			//a := <-tk.actionChan
			//log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := len(tk.errChan); i > 0; i-- {
			<-tk.errChan
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
