package main

import (
	"database/sql"
	"log"

	"github.com/satori/go.uuid"
)

func (tk *treeKeeper) startupRepositorySystemProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                              error
		instanceId, srcInstanceId, repositoryId, view, systemProperty, sourceType, value string
		inInstanceId, inObjectType, inObjId                                              string
		inheritance, childrenOnly                                                        bool
		rows, instance_rows                                                              *sql.Rows
		load_properties, load_instances                                                  *sql.Stmt
	)
	log.Println("Prepare: treekeeper/load-repository-system-properties")
	load_properties, err = tk.conn.Prepare(`
SELECT instance_id,
       source_instance_id,
	   repository_id,
	   view,
	   system_property,
	   source_type,
	   inheritance_enabled,
	   children_only,
	   value
FROM   soma.repository_system_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-repository-system-properties: ", err)
	}
	defer load_properties.Close()

	log.Println("Prepare: treekeeper/load-repository-system-property-instances")
	load_instances, err = tk.conn.Prepare(tkStmtLoadSystemPropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-repository-system-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading repository system properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading repository system properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

systemloop:
	// load all system properties defined directly on group objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&repositoryId,
			&view,
			&systemProperty,
			&sourceType,
			&inheritance,
			&childrenOnly,
			&value,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break systemloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := somatree.PropertySystem{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          systemProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Instances = make([]somatree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading repository system properties: %s", tk.repoName, err.Error())
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

			pi := somatree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the group and set the prepared property
		tk.tree.Find(somatree.FindRequest{
			ElementId: repositoryId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		for i := len(tk.actionChan); i > 0; i-- {
			a := <-tk.actionChan
			log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := 0; i < len(tk.errChan); i++ {
			<-tk.errChan
		}
	}
}

func (tk *treeKeeper) startupBucketSystemProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                          error
		instanceId, srcInstanceId, bucketId, view, systemProperty, sourceType, value string
		inInstanceId, inObjectType, inObjId                                          string
		inheritance, childrenOnly                                                    bool
		rows, instance_rows                                                          *sql.Rows
		load_properties, load_instances                                              *sql.Stmt
	)
	log.Println("Prepare: treekeeper/load-bucket-system-properties")
	load_properties, err = tk.conn.Prepare(`
SELECT instance_id,
       source_instance_id,
	   bucket_id,
	   view,
	   system_property,
	   source_type,
	   inheritance_enabled,
	   children_only,
	   value
FROM   soma.bucket_system_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-bucket-system-properties: ", err)
	}
	defer load_properties.Close()

	log.Println("Prepare: treekeeper/load-bucket-system-property-instances")
	load_instances, err = tk.conn.Prepare(tkStmtLoadSystemPropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-bucket-system-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading bucket system properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading bucket system properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

systemloop:
	// load all system properties defined directly on group objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&bucketId,
			&view,
			&systemProperty,
			&sourceType,
			&inheritance,
			&childrenOnly,
			&value,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break systemloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := somatree.PropertySystem{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          systemProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Instances = make([]somatree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading bucket system properties: %s", tk.repoName, err.Error())
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

			pi := somatree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the group and set the prepared property
		tk.tree.Find(somatree.FindRequest{
			ElementId: bucketId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		for i := len(tk.actionChan); i > 0; i-- {
			a := <-tk.actionChan
			log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := 0; i < len(tk.errChan); i++ {
			<-tk.errChan
		}
	}
}

func (tk *treeKeeper) startupGroupSystemProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                         error
		instanceId, srcInstanceId, groupId, view, systemProperty, sourceType, value string
		inInstanceId, inObjectType, inObjId                                         string
		inheritance, childrenOnly                                                   bool
		rows, instance_rows                                                         *sql.Rows
		load_properties, load_instances                                             *sql.Stmt
	)
	log.Println("Prepare: treekeeper/load-group-system-properties")
	load_properties, err = tk.conn.Prepare(`
SELECT instance_id,
       source_instance_id,
	   group_id,
	   view,
	   system_property,
	   source_type,
	   inheritance_enabled,
	   children_only,
	   value
FROM   soma.group_system_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-group-system-properties: ", err)
	}
	defer load_properties.Close()

	log.Println("Prepare: treekeeper/load-group-system-property-instances")
	load_instances, err = tk.conn.Prepare(tkStmtLoadSystemPropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-group-system-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading group system properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading group system properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

systemloop:
	// load all system properties defined directly on group objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&groupId,
			&view,
			&systemProperty,
			&sourceType,
			&inheritance,
			&childrenOnly,
			&value,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break systemloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := somatree.PropertySystem{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          systemProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Instances = make([]somatree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading group system properties: %s", tk.repoName, err.Error())
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
			a := <-tk.actionChan
			log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := 0; i < len(tk.errChan); i++ {
			<-tk.errChan
		}
	}
}

func (tk *treeKeeper) startupClusterSystemProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                           error
		instanceId, srcInstanceId, clusterId, view, systemProperty, sourceType, value string
		inInstanceId, inObjectType, inObjId                                           string
		inheritance, childrenOnly                                                     bool
		rows, instance_rows                                                           *sql.Rows
		load_properties, load_instances                                               *sql.Stmt
	)
	log.Println("Prepare: treekeeper/load-cluster-system-properties")
	load_properties, err = tk.conn.Prepare(`
SELECT instance_id,
       source_instance_id,
	   cluster_id,
	   view,
	   system_property,
	   source_type,
	   inheritance_enabled,
	   children_only,
	   value
FROM   soma.cluster_system_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-cluster-system-properties: ", err)
	}
	defer load_properties.Close()

	log.Println("Prepare: treekeeper/load-cluster-system-property-instances")
	load_instances, err = tk.conn.Prepare(tkStmtLoadSystemPropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-cluster-system-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading cluster system properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading cluster system properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

systemloop:
	// load all system properties defined directly on group objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&clusterId,
			&view,
			&systemProperty,
			&sourceType,
			&inheritance,
			&childrenOnly,
			&value,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break systemloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := somatree.PropertySystem{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          systemProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Instances = make([]somatree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading cluster system properties: %s", tk.repoName, err.Error())
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

			pi := somatree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the group and set the prepared property
		tk.tree.Find(somatree.FindRequest{
			ElementId: clusterId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		for i := len(tk.actionChan); i > 0; i-- {
			a := <-tk.actionChan
			log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := 0; i < len(tk.errChan); i++ {
			<-tk.errChan
		}
	}
}

func (tk *treeKeeper) startupNodeSystemProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                        error
		instanceId, srcInstanceId, nodeId, view, systemProperty, sourceType, value string
		inInstanceId, inObjectType, inObjId                                        string
		inheritance, childrenOnly                                                  bool
		rows, instance_rows                                                        *sql.Rows
		load_properties, load_instances                                            *sql.Stmt
	)
	log.Println("Prepare: treekeeper/load-node-system-properties")
	load_properties, err = tk.conn.Prepare(`
SELECT instance_id,
       source_instance_id,
	   node_id,
	   view,
	   system_property,
	   source_type,
	   inheritance_enabled,
	   children_only,
	   value
FROM   soma.node_system_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-node-system-properties: ", err)
	}
	defer load_properties.Close()

	log.Println("Prepare: treekeeper/load-node-system-property-instances")
	load_instances, err = tk.conn.Prepare(tkStmtLoadSystemPropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-node-system-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading node system properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading node system properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

systemloop:
	// load all system properties defined directly on group objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&nodeId,
			&view,
			&systemProperty,
			&sourceType,
			&inheritance,
			&childrenOnly,
			&value,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break systemloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := somatree.PropertySystem{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          systemProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Instances = make([]somatree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading node system properties: %s", tk.repoName, err.Error())
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

			pi := somatree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the group and set the prepared property
		tk.tree.Find(somatree.FindRequest{
			ElementId: nodeId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		for i := len(tk.actionChan); i > 0; i-- {
			a := <-tk.actionChan
			log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := 0; i < len(tk.errChan); i++ {
			<-tk.errChan
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
