package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/satori/go.uuid"

)

type tkLoader interface {
	GroupState() *sql.Stmt
	GroupRelations() *sql.Stmt
}

type tkLoaderChecks struct {
	loadChecks     *sql.Stmt
	loadItems      *sql.Stmt
	loadConfig     *sql.Stmt
	loadThresh     *sql.Stmt
	loadCstrCustom *sql.Stmt
	loadCstrNative *sql.Stmt
	loadCstrOncall *sql.Stmt
	loadCstrAttr   *sql.Stmt
	loadCstrServ   *sql.Stmt
	loadCstrSystem *sql.Stmt
	loadInstances  *sql.Stmt
	loadInstConfig *sql.Stmt
	loadGroupState *sql.Stmt
	loadGroupRel   *sql.Stmt
}

func (ld *tkLoaderChecks) GroupState() *sql.Stmt {
	return ld.loadGroupState
}

func (ld *tkLoaderChecks) GroupRelations() *sql.Stmt {
	return ld.loadGroupRel
}

func (tk *treeKeeper) startupChecks() {
	var (
		err error
		ld  tkLoaderChecks
	)

	// load all required prepared statements into the loader structure
	ld = tkLoaderChecks{}
	if ld.loadChecks, err = tk.conn.Prepare(tkStmtLoadChecks); err != nil {
		log.Fatal("treekeeper/load-checks: ", err)
	}
	defer ld.loadChecks.Close()

	if ld.loadItems, err = tk.conn.Prepare(tkStmtLoadInheritedChecks); err != nil {
		log.Fatal("treekeeper/load-inherited-checks: ", err)
	}
	defer ld.loadItems.Close()

	if ld.loadConfig, err = tk.conn.Prepare(tkStmtLoadCheckConfiguration); err != nil {
		log.Fatal("treekeeper/load-check-configuration: ", err)
	}
	defer ld.loadConfig.Close()

	if ld.loadThresh, err = tk.conn.Prepare(tkStmtLoadCheckThresholds); err != nil {
		log.Fatal("treekeeper/load-check-thresholds: ", err)
	}
	defer ld.loadThresh.Close()

	if ld.loadCstrCustom, err = tk.conn.Prepare(tkStmtLoadCheckConstraintCustom); err != nil {
		log.Fatal("treekeeper/load-check-constraint-custom: ", err)
	}
	defer ld.loadCstrCustom.Close()

	if ld.loadCstrNative, err = tk.conn.Prepare(tkStmtLoadCheckConstraintNative); err != nil {
		log.Fatal("treekeeper/load-check-constraint-native: ", err)
	}
	defer ld.loadCstrNative.Close()

	if ld.loadCstrOncall, err = tk.conn.Prepare(tkStmtLoadCheckConstraintOncall); err != nil {
		log.Fatal("treekeeper/load-check-constraint-oncall: ", err)
	}
	defer ld.loadCstrOncall.Close()

	if ld.loadCstrAttr, err = tk.conn.Prepare(tkStmtLoadCheckConstraintAttribute); err != nil {
		log.Fatal("treekeeper/load-check-constraint-attribute: ", err)
	}
	defer ld.loadCstrAttr.Close()

	if ld.loadCstrServ, err = tk.conn.Prepare(tkStmtLoadCheckConstraintService); err != nil {
		log.Fatal("treekeeper/load-check-constraint-service: ", err)
	}
	defer ld.loadCstrServ.Close()

	if ld.loadCstrSystem, err = tk.conn.Prepare(tkStmtLoadCheckConstraintSystem); err != nil {
		log.Fatal("treekeeper/load-check-constraint-system: ", err)
	}
	defer ld.loadCstrSystem.Close()

	if ld.loadInstances, err = tk.conn.Prepare(tkStmtLoadCheckInstances); err != nil {
		log.Fatal("treekeeper/load-check-instance: ", err)
	}
	defer ld.loadInstances.Close()

	if ld.loadInstConfig, err = tk.conn.Prepare(tkStmtLoadCheckInstanceConfiguration); err != nil {
		log.Fatal("treekeeper/load-check-instance-configuration: ", err)
	}
	defer ld.loadInstConfig.Close()

	// since groups are the only tree elements that can be stacked,
	// additional ordering is required
	if ld.loadGroupState, err = tk.conn.Prepare(tkStmtLoadCheckGroupState); err != nil {
		log.Fatal("treekeeper/load-check-group-state: ", err)
	}
	defer ld.loadGroupState.Close()

	if ld.loadGroupRel, err = tk.conn.Prepare(tkStmtLoadCheckGroupRelations); err != nil {
		log.Fatal("treekeeper/load-check-group-relations: ", err)
	}
	defer ld.loadGroupRel.Close()

	// load checks for the entire tree, in order from root to leaf.
	for _, typ := range []string{`repository`, `bucket`, `group`, `cluster`, `node`} {
		tk.startupScopedChecks(typ, &ld)
	}
}

func (tk *treeKeeper) startupScopedChecks(typ string, ld *tkLoaderChecks) {
	if tk.broken {
		return
	}

	var (
		err                                                           error
		checkId, bucketId, srcCheckId, srcObjType, srcObjId, configId string
		capabilityId, objId, objType, cfgName, cfgObjId, cfgObjType   string
		externalId, predicate, threshold, levelName, levelShort       string
		cstrType, value1, value2, value3, itemId                      string
		levelNumeric, numVal                                          int64
		isActive, hasInheritance, isChildrenOnly, isEnabled           bool
		interval                                                      int64
		grOrder                                                       map[string][]string
		grWeird                                                       map[string]string
		ckRows, thrRows, cstrRows, itRows                             *sql.Rows
		cfgMap                                                        map[string]proto.CheckConfig
		victim                                                        proto.CheckConfig // go/issues/3117 workaround
		ckTree                                                        *somatree.Check
		ckItem                                                        somatree.CheckItem
		ckOrder                                                       map[string]map[string]somatree.Check
	)

	switch typ {
	case "group":
		if err, grOrder, grWeird = tk.orderGroups(ld); err != nil {
			goto fail
		}
	}

	if ckRows, err = ld.loadChecks.Query(tk.repoId, typ); err == sql.ErrNoRows {
		// no checks on this element type
		return
	} else if err != nil {
		goto fail
	}
	defer ckRows.Close()

	// load all checks and start the assembly line
	for ckRows.Next() {
		if err = ckRows.Scan(
			&checkId,
			&bucketId,
			&srcCheckId,
			&srcObjType,
			&srcObjId,
			&configId,
			&capabilityId,
			&objId,
			&objType,
		); err != nil {
			goto fail
		}
		// save CheckConfig
		cfgMap[checkId] = proto.CheckConfig{
			Id:           configId,
			RepositoryId: tk.repoId,
			BucketId:     bucketId,
			CapabilityId: capabilityId,
			ObjectId:     objId,
			ObjectType:   objType,
		}
	}
	if ckRows.Err() != nil {
		goto fail
	}

	// iterate over the loaded checks and continue assembly with values
	// from the stored checkconfiguration
	for checkId, _ = range cfgMap {
		if err = ld.loadConfig.QueryRow(cfgMap[checkId].Id, tk.repoId).Scan(
			&bucketId,
			&cfgName,
			&cfgObjId,
			&cfgObjType,
			&isActive,
			&hasInheritance,
			&isChildrenOnly,
			&capabilityId,
			&interval,
			&isEnabled,
			&externalId,
		); err != nil {
			// sql.ErrNoRows is fatal here, the check exists - there
			// must be a configuration for it
			goto fail
		}

		victim = cfgMap[checkId]
		victim.Name = cfgName
		victim.Interval = uint64(interval)
		victim.IsActive = isActive
		victim.IsEnabled = isEnabled
		victim.Inheritance = hasInheritance
		victim.ChildrenOnly = isChildrenOnly
		victim.ExternalId = externalId
		cfgMap[checkId] = victim
	}

	// iterate over the loaded checks and continue assembly with values
	// from the stored thresholds
	for checkId, _ = range cfgMap {
		if thrRows, err = ld.loadThresh.Query(cfgMap[checkId].Id); err != nil {
			// sql.ErrNoRows is fatal here since a check without
			// thresholds is rather useless
			goto fail
		}

		victim = cfgMap[checkId]
		victim.Thresholds = []proto.CheckConfigThreshold{}

		// iterate over returned thresholds
		for thrRows.Next() {
			if err = thrRows.Scan(
				&predicate,
				&threshold,
				&levelName,
				&levelShort,
				&levelNumeric,
			); err != nil {
				thrRows.Close()
				goto fail
			}
			// ignore error since we converted this into the DB from int64
			numVal, _ = strconv.ParseInt(threshold, 10, 64)

			// save threshold
			victim.Thresholds = append(victim.Thresholds,
				proto.CheckConfigThreshold{
					Predicate: proto.Predicate{
						Symbol: predicate,
					},
					Level: proto.Level{
						Name:      levelName,
						ShortName: levelShort,
						Numeric:   uint16(levelNumeric),
					},
					Value: numVal,
				},
			)
		} // implict close by Next() reaching EOF
		if err = thrRows.Err(); err != nil {
			goto fail
		}
		cfgMap[checkId] = victim
		thrRows.Close()
	}

	// iterate over the loaded checks and continue assembly with values
	// from the stored constraints
	for checkId, _ = range cfgMap {
		victim = cfgMap[checkId]
		victim.Constraints = []proto.CheckConfigConstraint{}
		for _, cstrType = range []string{`custom`, `native`, `oncall`, `attribute`, `service`, `system`} {
			switch cstrType {
			case `custom`:
				cstrRows, err = ld.loadCstrCustom.Query(cfgMap[checkId].Id)
			case `native`:
				cstrRows, err = ld.loadCstrNative.Query(cfgMap[checkId].Id)
			case `oncall`:
				cstrRows, err = ld.loadCstrOncall.Query(cfgMap[checkId].Id)
			case `attribute`:
				cstrRows, err = ld.loadCstrAttr.Query(cfgMap[checkId].Id)
			case `service`:
				cstrRows, err = ld.loadCstrServ.Query(cfgMap[checkId].Id)
			case `system`:
				cstrRows, err = ld.loadCstrSystem.Query(cfgMap[checkId].Id)
			}
			if err == nil {
				goto fail
			}

			// iterate over returned thresholds - no rows is valid, as
			// constraints are not mandatory
			for cstrRows.Next() {
				if err = cstrRows.Scan(&value1, &value2, &value3); err != nil {
					cstrRows.Close()
					goto fail
				}
				switch cstrType {
				case `custom`:
					victim.Constraints = append(victim.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Custom: &proto.PropertyCustom{
								Id:           value1,
								Name:         value2,
								RepositoryId: tk.repoId,
								Value:        value3,
							},
						},
					)
				case `native`:
					victim.Constraints = append(victim.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Native: &proto.PropertyNative{
								Name:  value1,
								Value: value2,
							},
						},
					)
				case `oncall`:
					victim.Constraints = append(victim.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Oncall: &proto.PropertyOncall{
								Id:     value1,
								Name:   value2,
								Number: value3,
							},
						},
					)
				case `attribute`:
					victim.Constraints = append(victim.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Attribute: &proto.ServiceAttribute{
								Name:  value1,
								Value: value2,
							},
						},
					)
				case `service`:
					victim.Constraints = append(victim.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Service: &proto.PropertyService{
								Name:   value2,
								TeamId: value1,
							},
						},
					)
				case `system`:
					victim.Constraints = append(victim.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							System: &proto.PropertySystem{
								Name:  value1,
								Value: value2,
							},
						},
					)
				} // switch cstrType
			} // for cstrRows.Next()
			if cstrRows.Err() != nil {
				goto fail
			}
		}
		cfgMap[checkId] = victim
	}

	// iterate over the checks, convert them to tree.Check. Then load
	// the inherited IDs via loadItems and populate tree.Check.Items.
	// Save in datastructure: map[string]map[string]tree.Check
	//		objId -> checkId -> tree.Check
	// this way it is possible to access the checks by objId, which is
	// required to populate groups in the correct order.
	for checkId, _ = range cfgMap {
		victim = cfgMap[checkId]
		ckOrder[victim.ObjectId] = map[string]somatree.Check{}
		if ckTree, err = tk.convertCheck(&victim); err != nil {
			goto fail
		}
		// add source check as well so it gets recreated with the
		// correct UUID
		ckItem = somatree.CheckItem{ObjectType: victim.ObjectType}
		ckItem.ObjectId, _ = uuid.FromString(victim.ObjectId)
		ckItem.ItemId, _ = uuid.FromString(checkId)
		ckTree.Items = []somatree.CheckItem{ckItem}

		if itRows, err = ld.loadItems.Query(tk.repoId, checkId); err != nil {
			goto fail
		}

		for itRows.Next() {
			if err = itRows.Scan(
				&itemId,
				&objId,
				&objType,
			); err != nil {
				itRows.Close()
				goto fail
			}

			// create new object per iteration
			ckItem := somatree.CheckItem{ObjectType: objType}
			ckItem.ObjectId, _ = uuid.FromString(objId)
			ckItem.ItemId, _ = uuid.FromString(itemId)
			ckTree.Items = append(ckTree.Items, ckItem)
			ckOrder[victim.ObjectId][checkId] = *ckTree
		}
		if err = itRows.Err(); err != nil {
			goto fail
		}
	}

	// apply all somatree.Check object to the tree, special case
	// groups due to their ordering requirements
	//
	// grOrder maps from a standalone groupId to an array of child groupIds
	// ckOrder maps from a groupId to all source checks on it
	// ==> not every group has to have a check
	// ==> groups can have more than one check
	switch typ {
	case "group":
		for objKey, _ := range grOrder {
			// objKey is a standalone groupId. Test if there are
			// checks for it
			if _, ok := ckOrder[objKey]; ok {
				// apply all checks for objKey
				for ck, _ := range ckOrder[objKey] {
					tk.tree.Find(somatree.FindRequest{
						ElementType: cfgMap[ck].ObjectType,
						ElementId:   cfgMap[ck].ObjectId,
					}, true).SetCheck(ckOrder[objKey][ck])
					// drain after each check
					for i := len(tk.actionChan); i > 0; i-- {
						<-tk.actionChan
					}
					for i := len(tk.errChan); i > 0; i-- {
						<-tk.errChan
					}
				}
			}
			// iterate through all childgroups of objKey
			for pos, _ := range grOrder[objKey] {
				// test if there is a check for it
				if _, ok := ckOrder[grOrder[objKey][pos]]; ok {
					// apply all checks for grOrder[objKey][pos]
					for ck, _ := range ckOrder[objKey] {
						tk.tree.Find(somatree.FindRequest{
							ElementType: cfgMap[ck].ObjectType,
							ElementId:   cfgMap[ck].ObjectId,
						}, true).SetCheck(ckOrder[objKey][ck])
						// drain after each check
						for i := len(tk.actionChan); i > 0; i-- {
							<-tk.actionChan
						}
						for i := len(tk.errChan); i > 0; i-- {
							<-tk.errChan
						}
					}
				}
			}
		}
		// iterate through all weird groups as well
		for objKey, _ := range grWeird {
			// Test if there are checks for it
			if _, ok := ckOrder[objKey]; ok {
				// apply all checks for objKey
				for ck, _ := range ckOrder[objKey] {
					tk.tree.Find(somatree.FindRequest{
						ElementType: cfgMap[ck].ObjectType,
						ElementId:   cfgMap[ck].ObjectId,
					}, true).SetCheck(ckOrder[objKey][ck])
					// drain after each check
					for i := len(tk.actionChan); i > 0; i-- {
						<-tk.actionChan
					}
					for i := len(tk.errChan); i > 0; i-- {
						<-tk.errChan
					}
				}
			}
		}
	default:
		for objKey, _ := range ckOrder {
			for ckKey, _ := range ckOrder[objKey] {
				tk.tree.Find(somatree.FindRequest{
					ElementType: cfgMap[ckKey].ObjectType,
					ElementId:   cfgMap[ckKey].ObjectId,
				}, true).SetCheck(ckOrder[objKey][ckKey])
				// drain after each check
				for i := len(tk.actionChan); i > 0; i-- {
					<-tk.actionChan
				}
				for i := len(tk.errChan); i > 0; i-- {
					<-tk.errChan
				}
			}
		}
	}

	/*
		Source-Checks laden: tkStmtLoadChecks, DONE
		Source-Checks.each:
			Config laden: tkStmtLoadCheckConfiguration, DONE
			Thresholds laden: tkStmtLoadCheckThresholds, DONE
			Constraints laden:	tkStmtLoadCheckConstraintCustom, DONE
								tkStmtLoadCheckConstraintNative, DONE
								tkStmtLoadCheckConstraintOncall, DONE
								tkStmtLoadCheckConstraintAttribute, DONE
								tkStmtLoadCheckConstraintService, DONE
								tkStmtLoadCheckConstraintSystem, DONE
			Vererbte Checks laden [CheckItem]: tkStmtLoadInheritedChecks, DONE
			--> Check anlegen, DONE
		!! << NICHT tk.tree.ComputeCheckInstances() AUFRUFEN >> !!
		CheckInstanzen laden: tkStmtLoadCheckInstances
		CheckInstanz.each:
			InstanzConfig laden: tkStmtLoadCheckInstanceConfiguration
			--> Instanz anlegen
			--> TODO: somatree CheckInstanz-Load Interface
	*/
	return

fail:
	tk.broken = true
	return
}

// orderGroups orders the groups in a repository so they can be
// processed from root to leaf
func (tk *treeKeeper) orderGroups(ld tkLoader) (error, map[string][]string, map[string]string) {
	if tk.broken {
		return fmt.Errorf("Broken tree detected"), nil, nil
	}

	var (
		err                                                 error
		groupId, groupState, parentId, childId, candidateId string
		stRows, rlRows                                      *sql.Rows
		ok                                                  bool
		grStateMap, grRelMap, grWeirdMap                    map[string]string
		grOrder                                             map[string][]string
		children                                            []string
		oldLen, sameCount                                   int
	)

	grStateMap = map[string]string{}
	grRelMap = map[string]string{}
	grWeirdMap = map[string]string{}
	grOrder = map[string][]string{}
	children = []string{}

	// load groups in this repository
	if stRows, err = ld.GroupState().Query(tk.repoId); err != nil {
		tk.broken = true
		return err, nil, nil
	}
	defer stRows.Close()

	for stRows.Next() {
		if err = stRows.Scan(&groupId, &groupState); err == sql.ErrNoRows {
			// this repository has no groups
			return nil, grOrder, grWeirdMap
		} else if err != nil {
			// error loading group state
			tk.broken = true
			return err, nil, nil
		}
		grStateMap[groupId] = groupState
	}
	if err = stRows.Err(); err == sql.ErrNoRows {
		// this repository has no groups
		return nil, grOrder, grWeirdMap
	} else if err != nil {
		tk.broken = true
		return err, nil, nil
	}

	// load relations between groups in this repository
	if rlRows, err = ld.GroupRelations().Query(tk.repoId); err != nil {
		tk.broken = true
		return err, nil, nil
	}
	defer rlRows.Close()

relations:
	for rlRows.Next() {
		if err = rlRows.Scan(&groupId, &childId); err == sql.ErrNoRows {
			// no stacked groups
			break relations
		} else if err != nil {
			// error loading relations
			tk.broken = true
			return err, nil, nil
		}
		// every group can only be child to one parent group, but
		// a parent group can have multiple child groups => data
		// needs to be stored this way
		grRelMap[childId] = groupId
	}
	if err = rlRows.Err(); err != nil {
		tk.broken = true
		return err, nil, nil
	}

	// iterate over all groups and identify state standalone,
	// attached to the bucket. once ordered, remove them.
	for groupId, groupState = range grStateMap {
		switch groupState {
		case "grouped":
			continue
		case "standalone":
			grOrder[groupId] = []string{}
		default:
			// groups should really not be in a different state
			grWeirdMap[groupId] = groupState
		}
		delete(grStateMap, groupId)
	}

	// simplified first ordering round, only sort children of
	// standalone groups
	for childId, groupId = range grRelMap {
		if _, ok = grOrder[groupId]; !ok {
			// groupId is not standalone
			continue
		}
		// groupId is standalone
		grOrder[groupId] = append(grOrder[groupId], childId)
		delete(grRelMap, childId)
		delete(grStateMap, childId)
	}

	// breaker switch variables to short-circuit an unbounded loop
	oldLen = 0
	sameCount = 0

	// full ordering, grStateMap contains all grouped groups left
orderloop:
	for len(grStateMap) > 0 {
		// install a breaker switch in case the groups can not be
		// ordered. if no elements were deleted from grStateMap
		// for three full iterations => give up
		if len(grStateMap) == oldLen {
			sameCount++
		} else {
			oldLen = len(grStateMap)
			sameCount = 0
		}
		if sameCount >= 3 {
			break orderloop
		}

		// iterate through all unordered children
	childloop:
		for childId, parentId = range grRelMap {
			// since childId was not ordered during the first
			// pass, its parentId is a grouped group itself
			for groupId, children = range grOrder {
				// iterate through all children
				for _, candidateId = range children {
					if candidateId == parentId {
						// this candidateId is the searched parent.
						// since candidateId is a child of
						// groupId, childId is appended there
						grOrder[groupId] = append(grOrder[groupId], childId)
						delete(grStateMap, childId)
						delete(grRelMap, childId)
						continue childloop
					}
				}
			}
		}
	}
	if sameCount >= 3 || len(grRelMap) != 0 {
		// breaker went off or we have unordered grRelMap left
		tk.broken = true
		return fmt.Errorf("Failed to order groups"), nil, nil
	}

	// return order
	return nil, grOrder, grWeirdMap
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
