package soma

import (
	"database/sql"

	"github.com/1and1/soma/internal/tree"
	uuid "github.com/satori/go.uuid"
)

func (tk *TreeKeeper) startupOncallProperties(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		err                                                           error
		instanceId, srcInstanceId, objectId, view, oncallId           string
		inInstanceId, inObjectType, inObjId, oncallName, oncallNumber string
		inheritance, childrenOnly                                     bool
		rows, instance_rows                                           *sql.Rows
	)

	for loopType, loopStmt := range map[string]string{
		`repository`: `LoadPropRepoOncall`,
		`bucket`:     `LoadPropBuckOncall`,
		`group`:      `LoadPropGrpOncall`,
		`cluster`:    `LoadPropClrOncall`,
		`node`:       `LoadPropNodeOncall`,
	} {

		tk.startLog.Printf("TK[%s]: loading %s oncall properties\n", tk.meta.repoName, loopType)
		rows, err = stMap[loopStmt].Query(tk.meta.repoID)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error loading %s oncall properties: %s", tk.meta.repoName, loopType, err.Error())
			tk.status.isBroken = true
			return
		}
		defer rows.Close()

	oncallloop:
		// load all oncall properties defined directly on objects
		for rows.Next() {
			err = rows.Scan(
				&instanceId,
				&srcInstanceId,
				&objectId,
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
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
				tk.status.isBroken = true
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

			instance_rows, err = stMap[`LoadPropOncallInstance`].Query(
				tk.meta.repoID,
				srcInstanceId,
			)
			if err != nil {
				tk.startLog.Printf("TK[%s] Error loading %s oncall properties: %s", tk.meta.repoName, loopType, err.Error())
				tk.status.isBroken = true
				return
			}
			defer instance_rows.Close()

		inproploop:
			// load all all ids for properties that were inherited from the
			// current oncall property so the IDs can be set correctly
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
					tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
					tk.status.isBroken = true
					return
				}

				var propObjectId, propInstanceId uuid.UUID
				if propObjectId, err = uuid.FromString(inObjId); err != nil {
					tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
					tk.status.isBroken = true
					return
				}
				if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
					tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
					tk.status.isBroken = true
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

			// lookup the object and set the prepared property
			tk.tree.Find(tree.FindRequest{
				ElementType: loopType,
				ElementId:   objectId,
			}, true).SetProperty(&prop)

			// throw away all generated actions, we do this for every
			// property since with inheritance this can create a lot of
			// actions
			tk.drain(`action`)
			tk.drain(`error`)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
