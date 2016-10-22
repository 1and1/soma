package main

import (
	"database/sql"

	"github.com/1and1/soma/internal/tree"
	uuid "github.com/satori/go.uuid"
)

func (tk *treeKeeper) startupRepositoryOncallProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                               error
		instanceId, srcInstanceId, repositoryId, view, oncallId, oncallName, oncallNumber string
		inInstanceId, inObjectType, inObjId                                               string
		inheritance, childrenOnly                                                         bool
		rows, instance_rows                                                               *sql.Rows
		load_properties, load_instances                                                   *sql.Stmt
	)
	load_properties, err = tk.conn.Prepare(`
SELECT  srop.instance_id,
        srop.source_instance_id,
        srop.repository_id,
        srop.view,
        srop.oncall_duty_id,
        srop.inheritance_enabled,
        srop.children_only,
        iodt.oncall_duty_name,
        iodt.oncall_duty_phone_number
FROM    soma.repository_oncall_properties srop
JOIN    inventory.oncall_duty_teams iodt
  ON    srop.oncall_duty_id = iodt.oncall_duty_id
WHERE   srop.instance_id = srop.source_instance_id
  AND   srop.repository_id = $1::uuid;`)
	if err != nil {
		tk.startLog.Println("treekeeper/load-repository-oncall-properties: ", err)
		tk.broken = true
		return
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(tkStmtLoadOncallPropInstances)
	if err != nil {
		tk.startLog.Println("treekeeper/load-repository-oncall-property-instances: ", err)
		tk.broken = true
		return
	}
	defer load_instances.Close()

	tk.startLog.Printf("TK[%s]: loading repository oncall properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading repository oncall properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

oncallloop:
	// load all oncall properties defined directly on repository objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&repositoryId,
			&view,
			&oncallId,
			&inheritance,
			&childrenOnly,
			&oncallName,
			&oncallNumber,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break oncallloop
			}
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertyOncall{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Name:         oncallName,
			Number:       oncallNumber,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.OncallId, _ = uuid.FromString(oncallId)
		prop.Instances = make([]tree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error loading repository oncall properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current repository oncall property so the IDs can be set correctly
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

func (tk *treeKeeper) startupBucketOncallProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                           error
		instanceId, srcInstanceId, bucketId, view, oncallId, oncallName, oncallNumber string
		inInstanceId, inObjectType, inObjId                                           string
		inheritance, childrenOnly                                                     bool
		rows, instance_rows                                                           *sql.Rows
		load_properties, load_instances                                               *sql.Stmt
	)
	load_properties, err = tk.conn.Prepare(`
SELECT  sgop.instance_id,
        sgop.source_instance_id,
        sgop.bucket_id,
        sgop.view,
        sgop.oncall_duty_id,
        sgop.inheritance_enabled,
        sgop.children_only,
        iodt.oncall_duty_name,
        iodt.oncall_duty_phone_number
FROM    soma.bucket_oncall_properties sgop
JOIN    inventory.oncall_duty_teams iodt
  ON    sgop.oncall_duty_id = iodt.oncall_duty_id
WHERE   sgop.instance_id = sgop.source_instance_id
  AND   sgop.repository_id = $1::uuid;`)
	if err != nil {
		tk.startLog.Println("treekeeper/load-bucket-oncall-properties: ", err)
		tk.broken = true
		return
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(tkStmtLoadOncallPropInstances)
	if err != nil {
		tk.startLog.Println("treekeeper/load-bucket-oncall-property-instances: ", err)
		tk.broken = true
		return
	}
	defer load_instances.Close()

	tk.startLog.Printf("TK[%s]: loading bucket oncall properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading bucket oncall properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

oncallloop:
	// load all oncall properties defined directly on bucket objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&bucketId,
			&view,
			&oncallId,
			&inheritance,
			&childrenOnly,
			&oncallName,
			&oncallNumber,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break oncallloop
			}
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertyOncall{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Name:         oncallName,
			Number:       oncallNumber,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.OncallId, _ = uuid.FromString(oncallId)
		prop.Instances = make([]tree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error loading bucket oncall properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current bucket oncall property so the IDs can be set correctly
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

func (tk *treeKeeper) startupGroupOncallProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                          error
		instanceId, srcInstanceId, groupId, view, oncallId, oncallName, oncallNumber string
		inInstanceId, inObjectType, inObjId                                          string
		inheritance, childrenOnly                                                    bool
		rows, instance_rows                                                          *sql.Rows
		load_properties, load_instances                                              *sql.Stmt
	)
	load_properties, err = tk.conn.Prepare(`
SELECT  sgop.instance_id,
        sgop.source_instance_id,
        sgop.group_id,
        sgop.view,
        sgop.oncall_duty_id,
        sgop.inheritance_enabled,
        sgop.children_only,
        iodt.oncall_duty_name,
        iodt.oncall_duty_phone_number
FROM    soma.group_oncall_properties sgop
JOIN    inventory.oncall_duty_teams iodt
  ON    sgop.oncall_duty_id = iodt.oncall_duty_id
WHERE   sgop.instance_id = sgop.source_instance_id
  AND   sgop.repository_id = $1::uuid;`)
	if err != nil {
		tk.startLog.Println("treekeeper/load-group-oncall-properties: ", err)
		tk.broken = true
		return
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(tkStmtLoadOncallPropInstances)
	if err != nil {
		tk.startLog.Println("treekeeper/load-group-oncall-property-instances: ", err)
		tk.broken = true
		return
	}
	defer load_instances.Close()

	tk.startLog.Printf("TK[%s]: loading group oncall properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading group oncall properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

oncallloop:
	// load all system properties defined directly on group objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&groupId,
			&view,
			&oncallId,
			&inheritance,
			&childrenOnly,
			&oncallName,
			&oncallNumber,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break oncallloop
			}
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertyOncall{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Name:         oncallName,
			Number:       oncallNumber,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.OncallId, _ = uuid.FromString(oncallId)
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
			ElementType: `group`,
			ElementId:   groupId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		tk.drain(`action`)
		tk.drain(`error`)
	}
}

func (tk *treeKeeper) startupClusterOncallProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                            error
		instanceId, srcInstanceId, clusterId, view, oncallId, oncallName, oncallNumber string
		inInstanceId, inObjectType, inObjId                                            string
		inheritance, childrenOnly                                                      bool
		rows, instance_rows                                                            *sql.Rows
		load_properties, load_instances                                                *sql.Stmt
	)
	load_properties, err = tk.conn.Prepare(`
SELECT  scop.instance_id,
        scop.source_instance_id,
        scop.cluster_id,
        scop.view,
        scop.oncall_duty_id,
        scop.inheritance_enabled,
        scop.children_only,
        iodt.oncall_duty_name,
        iodt.oncall_duty_phone_number
FROM    soma.cluster_oncall_properties scop
JOIN    inventory.oncall_duty_teams iodt
  ON    scop.oncall_duty_id = iodt.oncall_duty_id
WHERE   scop.instance_id = scop.source_instance_id
  AND   scop.repository_id = $1::uuid;`)
	if err != nil {
		tk.startLog.Println("treekeeper/load-cluster-oncall-properties: ", err)
		tk.broken = true
		return
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(tkStmtLoadOncallPropInstances)
	if err != nil {
		tk.startLog.Println("treekeeper/load-cluster-oncall-property-instances: ", err)
		tk.broken = true
		return
	}
	defer load_instances.Close()

	tk.startLog.Printf("TK[%s]: loading cluster oncall properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading cluster oncall properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

oncallloop:
	// load all system properties defined directly on group objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&clusterId,
			&view,
			&oncallId,
			&inheritance,
			&childrenOnly,
			&oncallName,
			&oncallNumber,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break oncallloop
			}
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertyOncall{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Name:         oncallName,
			Number:       oncallNumber,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.OncallId, _ = uuid.FromString(oncallId)
		prop.Instances = make([]tree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error loading cluster oncall properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current cluster oncall property so the IDs can be set correctly
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

func (tk *treeKeeper) startupNodeOncallProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                         error
		instanceId, srcInstanceId, nodeId, view, oncallId, oncallName, oncallNumber string
		inInstanceId, inObjectType, inObjId                                         string
		inheritance, childrenOnly                                                   bool
		rows, instance_rows                                                         *sql.Rows
		load_properties, load_instances                                             *sql.Stmt
	)
	load_properties, err = tk.conn.Prepare(`
SELECT  snop.instance_id,
        snop.source_instance_id,
        snop.node_id,
        snop.view,
        snop.oncall_duty_id,
        snop.inheritance_enabled,
        snop.children_only,
        iodt.oncall_duty_name,
        iodt.oncall_duty_phone_number
FROM    soma.node_oncall_property snop
JOIN    inventory.oncall_duty_teams iodt
  ON    snop.oncall_duty_id = iodt.oncall_duty_id
WHERE   snop.instance_id = snop.source_instance_id
  AND   snop.repository_id = $1::uuid;`)
	if err != nil {
		tk.startLog.Println("treekeeper/load-node-oncall-properties: ", err)
		tk.broken = true
		return
	}
	defer load_properties.Close()

	load_instances, err = tk.conn.Prepare(tkStmtLoadOncallPropInstances)
	if err != nil {
		tk.startLog.Println("treekeeper/load-node-oncall-property-instances: ", err)
		tk.broken = true
		return
	}
	defer load_instances.Close()

	tk.startLog.Printf("TK[%s]: loading node oncall properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading group oncall properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

oncallloop:
	// load all system properties defined directly on group objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&nodeId,
			&view,
			&oncallId,
			&inheritance,
			&childrenOnly,
			&oncallName,
			&oncallNumber,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break oncallloop
			}
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertyOncall{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Name:         oncallName,
			Number:       oncallNumber,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.OncallId, _ = uuid.FromString(oncallId)
		prop.Instances = make([]tree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error loading node oncall properties: %s", tk.repoName, err.Error())
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
