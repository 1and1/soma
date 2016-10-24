package main

import (
	"database/sql"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/internal/tree"
	uuid "github.com/satori/go.uuid"
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

	load_properties, err = tk.conn.Prepare(stmt.TkStartLoadRepositoryCstProp)
	if err != nil {
		tk.startLog.Println("treekeeper/load-repository-custom-properties: ", err)
		tk.broken = true
		return
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(stmt.TkStartLoadCustomPropInstances)
	if err != nil {
		tk.startLog.Println("treekeeper/load-repository-custom-property-instances: ", err)
		tk.broken = true
		return
	}
	defer load_instances.Close()

	tk.startLog.Printf("TK[%s]: loading repository custom properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading repository custom properties: %s", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error loading repository custom properties: %s", tk.repoName, err.Error())
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
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			var propObjectId, propInstanceId uuid.UUID
			if propObjectId, err = uuid.FromString(inObjId); err != nil {
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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

	load_properties, err = tk.conn.Prepare(stmt.TkStartLoadBucketCstProp)
	if err != nil {
		tk.startLog.Println("treekeeper/load-bucket-custom-properties: ", err)
		tk.broken = true
		return
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(stmt.TkStartLoadCustomPropInstances)
	if err != nil {
		tk.startLog.Println("treekeeper/load-bucket-custom-property-instances: ", err)
		tk.broken = true
		return
	}
	defer load_instances.Close()

	tk.startLog.Printf("TK[%s]: loading bucket custom properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading bucket custom properties: %s", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error loading bucket custom properties: %s", tk.repoName, err.Error())
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
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			var propObjectId, propInstanceId uuid.UUID
			if propObjectId, err = uuid.FromString(inObjId); err != nil {
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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
	load_properties, err = tk.conn.Prepare(stmt.TkStartLoadGroupCstProp)
	if err != nil {
		tk.startLog.Println("treekeeper/load-group-custom-properties: ", err)
		tk.broken = true
		return
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(stmt.TkStartLoadCustomPropInstances)
	if err != nil {
		tk.startLog.Println("treekeeper/load-group-custom-property-instances: ", err)
		tk.broken = true
		return
	}
	defer load_instances.Close()

	tk.startLog.Printf("TK[%s]: loading group custom properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading group custom properties: %s", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error loading group custom properties: %s", tk.repoName, err.Error())
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
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			var propObjectId, propInstanceId uuid.UUID
			if propObjectId, err = uuid.FromString(inObjId); err != nil {
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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
	load_properties, err = tk.conn.Prepare(stmt.TkStartLoadClusterCstProp)
	if err != nil {
		tk.startLog.Println("treekeeper/load-cluster-custom-properties: ", err)
		tk.broken = true
		return
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(stmt.TkStartLoadCustomPropInstances)
	if err != nil {
		tk.startLog.Println("treekeeper/load-cluster-custom-property-instances: ", err)
		tk.broken = true
		return
	}
	defer load_instances.Close()

	tk.startLog.Printf("TK[%s]: loading cluster custom properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading cluster custom properties: %s", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error loading cluster custom properties: %s", tk.repoName, err.Error())
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
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			var propObjectId, propInstanceId uuid.UUID
			if propObjectId, err = uuid.FromString(inObjId); err != nil {
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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
	load_properties, err = tk.conn.Prepare(stmt.TkStartLoadNodeCstProp)
	if err != nil {
		tk.startLog.Println("treekeeper/load-node-custom-properties: ", err)
		tk.broken = true
		return
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(stmt.TkStartLoadCustomPropInstances)
	if err != nil {
		tk.startLog.Println("treekeeper/load-node-custom-property-instances: ", err)
		tk.broken = true
		return
	}
	defer load_instances.Close()

	tk.startLog.Printf("TK[%s]: loading node custom properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading node custom properties: %s", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error loading node custom properties: %s", tk.repoName, err.Error())
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
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			var propObjectId, propInstanceId uuid.UUID
			if propObjectId, err = uuid.FromString(inObjId); err != nil {
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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
