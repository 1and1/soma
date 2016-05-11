package main

import (
	"database/sql"
	"fmt"
	"log"
)

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
		err     error
		grOrder map[string][]string
		grWeird map[string]string
	)

	switch typ {
	case "group":
		if err, grOrder, grWeird = tk.orderGroups(ld); err != nil {
			tk.broken = true
			return
		}
	}
	//XXX: silence declared and not used for now
	fmt.Printf("%+v\n", grOrder)
	fmt.Printf("%+v\n", grWeird)

	/*
		Source-Checks laden: tkStmtLoadChecks
		Source-Checks.each:
			Config laden: tkStmtLoadCheckConfiguration
			Thresholds laden: tkStmtLoadCheckThresholds
			Constraints laden:	tkStmtLoadCheckConstraintCustom
								tkStmtLoadCheckConstraintNative
								tkStmtLoadCheckConstraintOncall
								tkStmtLoadCheckConstraintAttribute
								tkStmtLoadCheckConstraintService
								tkStmtLoadCheckConstraintSystem
			Vererbte Checks laden [CheckItem]: tkStmtLoadInheritedChecks
			--> Check anlegen
		!! << NICHT tk.tree.ComputeCheckInstances() AUFRUFEN >> !!
		CheckInstanzen laden: tkStmtLoadCheckInstances
		CheckInstanz.each:
			InstanzConfig laden: tkStmtLoadCheckInstanceConfiguration
			--> Instanz anlegen
			--> TODO: somatree CheckInstanz-Load Interface
	*/
}

// orderGroups orders the groups in a repository so they can be
// processed from root to leaf
func (tk *treeKeeper) orderGroups(ld *tkLoaderChecks) (error, map[string][]string, map[string]string) {
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
	if stRows, err = ld.loadGroupState.Query(tk.repoId); err != nil {
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
	if rlRows, err = ld.loadGroupRel.Query(tk.repoId); err != nil {
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
