package soma

import (
	"database/sql"

	"github.com/1and1/soma/internal/tree"
	"github.com/1and1/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

func (tk *TreeKeeper) startupServiceProperties(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		err                                          error
		instanceId, srcInstanceId, objectId, view    string
		inInstanceId, inObjectType, inObjId, attrKey string
		serviceProperty, teamId, attrValue           string
		inheritance, childrenOnly                    bool
		rows, attribute_rows, instance_rows          *sql.Rows
	)

	for loopType, loopStmt := range map[string][2]string{
		`repository`: [2]string{
			`LoadPropRepoService`,
			`LoadPropRepoSvcAttr`},
		`bucket`: [2]string{
			`LoadPropBuckService`,
			`LoadPropBuckSvcAttr`},
		`group`: [2]string{
			`LoadPropGrpService`,
			`LoadPropGrpSvcAttr`},
		`cluster`: [2]string{
			`LoadPropClrService`,
			`LoadPropClrSvcAttr`},
		`node`: [2]string{
			`LoadPropNodeService`,
			`LoadPropNodeSvcAttr`},
	} {

		tk.startLog.Printf("TK[%s]: loading %s service properties\n", tk.meta.repoName, loopType)
		rows, err = stMap[loopStmt[0]].Query(tk.meta.repoID)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error loading %s service properties: %s", tk.meta.repoName, loopType, err.Error())
			tk.status.isBroken = true
			return
		}
		defer rows.Close()

	serviceloop:
		// load all service properties defined directly on the object
		for rows.Next() {
			err = rows.Scan(
				&instanceId,
				&srcInstanceId,
				&objectId,
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
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
				tk.status.isBroken = true
				return
			}

			// build the property
			prop := tree.PropertyService{
				Inheritance:  inheritance,
				ChildrenOnly: childrenOnly,
				View:         view,
				Service:      serviceProperty,
			}
			prop.Id, _ = uuid.FromString(instanceId)
			prop.Attributes = make([]proto.ServiceAttribute, 0)
			prop.Instances = make([]tree.PropertyInstance, 0)

			attribute_rows, err = stMap[loopStmt[1]].Query(
				teamId,
				serviceProperty,
			)
			if err != nil {
				tk.startLog.Printf("TK[%s] Error loading %s service properties: %s", tk.meta.repoName, loopType, err.Error())
				tk.status.isBroken = true
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
					tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
					tk.status.isBroken = true
					return
				}

				pa := proto.ServiceAttribute{
					Name:  attrKey,
					Value: attrValue,
				}
				prop.Attributes = append(prop.Attributes, pa)
			}

			instance_rows, err = stMap[`LoadPropSvcInstance`].Query(
				tk.meta.repoID,
				srcInstanceId,
			)
			if err != nil {
				tk.startLog.Printf("TK[%s] Error loading %s service properties: %s", tk.meta.repoName, loopType, err.Error())
				tk.status.isBroken = true
				return
			}
			defer instance_rows.Close()

		inproploop:
			// load all all ids for properties that were inherited from the
			// current service property so the IDs can be set correctly
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
