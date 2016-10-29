package main

import (
	"database/sql"

	"github.com/1and1/soma/internal/tree"
	uuid "github.com/satori/go.uuid"
)

func (tk *treeKeeper) startupRepositorySystemProperties(stMap map[string]*sql.Stmt) {
	if tk.broken {
		return
	}

	var (
		err                                                                              error
		instanceId, srcInstanceId, repositoryId, view, systemProperty, sourceType, value string
		inInstanceId, inObjectType, inObjId                                              string
		inheritance, childrenOnly                                                        bool
		rows, instance_rows                                                              *sql.Rows
	)

	tk.startLog.Printf("TK[%s]: loading repository system properties\n", tk.repoName)
	rows, err = stMap[`LoadPropRepoSystem`].Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading repository system properties: %s", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertySystem{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          systemProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Instances = make([]tree.PropertyInstance, 0)

		instance_rows, err = stMap[`LoadPropSystemInstance`].Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error loading repository system properties: %s", tk.repoName, err.Error())
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
			ElementId: repositoryId,
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

func (tk *treeKeeper) startupBucketSystemProperties(stMap map[string]*sql.Stmt) {
	if tk.broken {
		return
	}

	var (
		err                                                                          error
		instanceId, srcInstanceId, bucketId, view, systemProperty, sourceType, value string
		inInstanceId, inObjectType, inObjId                                          string
		inheritance, childrenOnly                                                    bool
		rows, instance_rows                                                          *sql.Rows
	)

	tk.startLog.Printf("TK[%s]: loading bucket system properties\n", tk.repoName)
	rows, err = stMap[`LoadPropBuckSystem`].Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading bucket system properties: %s", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertySystem{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          systemProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Instances = make([]tree.PropertyInstance, 0)

		instance_rows, err = stMap[`LoadPropSystemInstance`].Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error loading bucket system properties: %s", tk.repoName, err.Error())
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

func (tk *treeKeeper) startupGroupSystemProperties(stMap map[string]*sql.Stmt) {
	if tk.broken {
		return
	}

	var (
		err                                                                         error
		instanceId, srcInstanceId, groupId, view, systemProperty, sourceType, value string
		inInstanceId, inObjectType, inObjId                                         string
		inheritance, childrenOnly                                                   bool
		rows, instance_rows                                                         *sql.Rows
	)

	tk.startLog.Printf("TK[%s]: loading group system properties\n", tk.repoName)
	rows, err = stMap[`LoadPropGrpSystem`].Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading group system properties: %s", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertySystem{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          systemProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Instances = make([]tree.PropertyInstance, 0)

		instance_rows, err = stMap[`LoadPropSystemInstance`].Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error loading group system properties: %s", tk.repoName, err.Error())
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

func (tk *treeKeeper) startupClusterSystemProperties(stMap map[string]*sql.Stmt) {
	if tk.broken {
		return
	}

	var (
		err                                                                           error
		instanceId, srcInstanceId, clusterId, view, systemProperty, sourceType, value string
		inInstanceId, inObjectType, inObjId                                           string
		inheritance, childrenOnly                                                     bool
		rows, instance_rows                                                           *sql.Rows
	)

	tk.startLog.Printf("TK[%s]: loading cluster system properties\n", tk.repoName)
	rows, err = stMap[`LoadPropClrSystem`].Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading cluster system properties: %s", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertySystem{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          systemProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Instances = make([]tree.PropertyInstance, 0)

		instance_rows, err = stMap[`LoadPropSystemInstance`].Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error loading cluster system properties: %s", tk.repoName, err.Error())
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

func (tk *treeKeeper) startupNodeSystemProperties(stMap map[string]*sql.Stmt) {
	if tk.broken {
		return
	}

	var (
		err                                                                        error
		instanceId, srcInstanceId, nodeId, view, systemProperty, sourceType, value string
		inInstanceId, inObjectType, inObjId                                        string
		inheritance, childrenOnly                                                  bool
		rows, instance_rows                                                        *sql.Rows
	)

	tk.startLog.Printf("TK[%s]: loading node system properties\n", tk.repoName)
	rows, err = stMap[`LoadPropNodeSystem`].Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading node system properties: %s", tk.repoName, err.Error())
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := tree.PropertySystem{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          systemProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Instances = make([]tree.PropertyInstance, 0)

		instance_rows, err = stMap[`LoadPropSystemInstance`].Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error loading node system properties: %s", tk.repoName, err.Error())
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
