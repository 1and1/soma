/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */
package main

import (
	"database/sql"
	"fmt"

	"github.com/1and1/soma/internal/stmt"
)

func (tk *treeKeeper) startTx() (
	*sql.Tx, map[string]*sql.Stmt, error) {

	var err error
	var tx *sql.Tx
	open := false
	stMap := map[string]*sql.Stmt{}

	if tx, err = tk.conn.Begin(); err != nil {
		goto bailout
	}
	open = true

	//
	// PROPERTY STATEMENTS
	for name, statement := range map[string]string{
		`PropertyInstanceCreate`:          tkStmtPropertyInstanceCreate,
		`PropertyInstanceDelete`:          tkStmtPropertyInstanceDelete,
		`RepositoryPropertyOncallCreate`:  tkStmtRepositoryPropertyOncallCreate,
		`RepositoryPropertyOncallDelete`:  tkStmtRepositoryPropertyOncallDelete,
		`RepositoryPropertyServiceCreate`: tkStmtRepositoryPropertyServiceCreate,
		`RepositoryPropertyServiceDelete`: tkStmtRepositoryPropertyServiceDelete,
		`RepositoryPropertySystemCreate`:  tkStmtRepositoryPropertySystemCreate,
		`RepositoryPropertySystemDelete`:  tkStmtRepositoryPropertySystemDelete,
		`RepositoryPropertyCustomCreate`:  tkStmtRepositoryPropertyCustomCreate,
		`RepositoryPropertyCustomDelete`:  tkStmtRepositoryPropertyCustomDelete,
		`BucketPropertyOncallCreate`:      tkStmtBucketPropertyOncallCreate,
		`BucketPropertyOncallDelete`:      tkStmtBucketPropertyOncallDelete,
		`BucketPropertyServiceCreate`:     tkStmtBucketPropertyServiceCreate,
		`BucketPropertyServiceDelete`:     tkStmtBucketPropertyServiceDelete,
		`BucketPropertySystemCreate`:      tkStmtBucketPropertySystemCreate,
		`BucketPropertySystemDelete`:      tkStmtBucketPropertySystemDelete,
		`BucketPropertyCustomCreate`:      tkStmtBucketPropertyCustomCreate,
		`BucketPropertyCustomDelete`:      tkStmtBucketPropertyCustomDelete,
		`GroupPropertyOncallCreate`:       tkStmtGroupPropertyOncallCreate,
		`GroupPropertyOncallDelete`:       tkStmtGroupPropertyOncallDelete,
		`GroupPropertyServiceCreate`:      tkStmtGroupPropertyServiceCreate,
		`GroupPropertyServiceDelete`:      tkStmtGroupPropertyServiceDelete,
		`GroupPropertySystemCreate`:       tkStmtGroupPropertySystemCreate,
		`GroupPropertySystemDelete`:       tkStmtGroupPropertySystemDelete,
		`GroupPropertyCustomCreate`:       tkStmtGroupPropertyCustomCreate,
		`GroupPropertyCustomDelete`:       tkStmtGroupPropertyCustomDelete,
		`ClusterPropertyOncallCreate`:     tkStmtClusterPropertyOncallCreate,
		`ClusterPropertyOncallDelete`:     tkStmtClusterPropertyOncallDelete,
		`ClusterPropertyServiceCreate`:    tkStmtClusterPropertyServiceCreate,
		`ClusterPropertyServiceDelete`:    tkStmtClusterPropertyServiceDelete,
		`ClusterPropertySystemCreate`:     tkStmtClusterPropertySystemCreate,
		`ClusterPropertySystemDelete`:     tkStmtClusterPropertySystemDelete,
		`ClusterPropertyCustomCreate`:     tkStmtClusterPropertyCustomCreate,
		`ClusterPropertyCustomDelete`:     tkStmtClusterPropertyCustomDelete,
		`NodePropertyOncallCreate`:        tkStmtNodePropertyOncallCreate,
		`NodePropertyOncallDelete`:        tkStmtNodePropertyOncallDelete,
		`NodePropertyServiceCreate`:       tkStmtNodePropertyServiceCreate,
		`NodePropertyServiceDelete`:       tkStmtNodePropertyServiceDelete,
		`NodePropertySystemCreate`:        tkStmtNodePropertySystemCreate,
		`NodePropertySystemDelete`:        tkStmtNodePropertySystemDelete,
		`NodePropertyCustomCreate`:        tkStmtNodePropertyCustomCreate,
		`NodePropertyCustomDelete`:        tkStmtNodePropertyCustomDelete,
	} {
		if stMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("tk.Prepare(%s) error: %s",
				name, err.Error())
			delete(stMap, name)
			goto bailout
		}
	}

	//
	// CHECK STATEMENTS
	for name, statement := range map[string]string{
		`CreateCheck`: stmt.TxCreateCheck,
		`DeleteCheck`: stmt.TxMarkCheckDeleted,
	} {
		if stMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("tk.Prepare(%s) error: %s",
				name, err.Error())
			delete(stMap, name)
			goto bailout
		}
	}

	//
	// CHECK INSTANCE STATEMENTS
	for name, statement := range map[string]string{
		`CreateCheckInstance`:              stmt.TxCreateCheckInstance,
		`CreateCheckInstanceConfiguration`: stmt.TxCreateCheckInstanceConfiguration,
		`DeleteCheckInstance`:              stmt.TxMarkCheckInstanceDeleted,
	} {
		if stMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("tk.Prepare(%s) error: %s",
				name, err.Error())
			delete(stMap, name)
			goto bailout
		}
	}

	//
	// CHECK CONFIGURATION STATEMENTS
	for name, statement := range map[string]string{
		`CreateCheckConfigurationBase`:                stmt.TxCreateCheckConfigurationBase,
		`CreateCheckConfigurationThreshold`:           stmt.TxCreateCheckConfigurationThreshold,
		`CreateCheckConfigurationConstraintSystem`:    stmt.TxCreateCheckConfigurationConstraintSystem,
		`CreateCheckConfigurationConstraintNative`:    stmt.TxCreateCheckConfigurationConstraintNative,
		`CreateCheckConfigurationConstraintOncall`:    stmt.TxCreateCheckConfigurationConstraintOncall,
		`CreateCheckConfigurationConstraintCustom`:    stmt.TxCreateCheckConfigurationConstraintCustom,
		`CreateCheckConfigurationConstraintService`:   stmt.TxCreateCheckConfigurationConstraintService,
		`CreateCheckConfigurationConstraintAttribute`: stmt.TxCreateCheckConfigurationConstraintAttribute,
	} {
		if stMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("tk.Prepare(%s) error: %s",
				name, err.Error())
			delete(stMap, name)
			goto bailout
		}
	}

	//
	// TREE MANIPULATION STATEMENTS
	for name, statement := range map[string]string{
		`BucketAssignNode`:         tkStmtBucketAssignNode,
		`ClusterCreate`:            tkStmtClusterCreate,
		`ClusterDelete`:            tkStmtClusterDelete,
		`ClusterMemberNew`:         tkStmtClusterMemberNew,
		`ClusterMemberRemove`:      tkStmtClusterMemberRemove,
		`ClusterUpdate`:            tkStmtClusterUpdate,
		`CreateBucket`:             tkStmtCreateBucket,
		`GroupCreate`:              tkStmtGroupCreate,
		`GroupDelete`:              tkStmtGroupDelete,
		`GroupMemberNewCluster`:    tkStmtGroupMemberNewCluster,
		`GroupMemberNewGroup`:      tkStmtGroupMemberNewGroup,
		`GroupMemberNewNode`:       tkStmtGroupMemberNewNode,
		`GroupMemberRemoveCluster`: tkStmtGroupMemberRemoveCluster,
		`GroupMemberRemoveGroup`:   tkStmtGroupMemberRemoveGroup,
		`GroupMemberRemoveNode`:    tkStmtGroupMemberRemoveNode,
		`GroupUpdate`:              tkStmtGroupUpdate,
		`NodeUnassignFromBucket`:   tkStmtNodeUnassignFromBucket,
		`UpdateNodeState`:          tkStmtUpdateNodeState,
	} {
		if stMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("tk.Prepare(%s) error: %s",
				name, err.Error())
			delete(stMap, name)
			goto bailout
		}
	}

	return tx, stMap, nil

bailout:
	if open {
		defer tx.Rollback()
	}
	for _, statement := range stMap {
		defer statement.Close()
	}
	return nil, nil, err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
