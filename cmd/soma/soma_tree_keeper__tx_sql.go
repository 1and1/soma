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
	for name, stmt := range map[string]string{
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
		if stMap[name], err = tx.Prepare(stmt); err != nil {
			delete(stMap, name)
			goto bailout
		}
	}

	//
	// CHECK STATEMENTS
	for name, stmt := range map[string]string{
		`CreateCheck`: stmt.TxCreateCheck,
		`DeleteCheck`: stmt.TxMarkCheckDeleted,
	} {
		if stMap[name], err = tx.Prepare(stmt); err != nil {
			delete(stMap, name)
			goto bailout
		}
	}

	//
	// CHECK INSTANCE STATEMENTS
	for name, stmt := range map[string]string{
		`CreateCheckInstance`:              stmt.TxCreateCheckInstance,
		`CreateCheckInstanceConfiguration`: stmt.TxCreateCheckInstanceConfiguration,
		`DeleteCheckInstance`:              stmt.TxMarkCheckInstanceDeleted,
	} {
		if stMap[name], err = tx.Prepare(stmt); err != nil {
			delete(stMap, name)
			goto bailout
		}
	}

	//
	// CHECK CONFIGURATION STATEMENTS
	for name, stmt := range map[string]string{
		`CreateCheckConfigurationBase`:                stmt.TxCreateCheckConfigurationBase,
		`CreateCheckConfigurationThreshold`:           stmt.TxCreateCheckConfigurationThreshold,
		`CreateCheckConfigurationConstraintSystem`:    stmt.TxCreateCheckConfigurationConstraintSystem,
		`CreateCheckConfigurationConstraintNative`:    stmt.TxCreateCheckConfigurationConstraintNative,
		`CreateCheckConfigurationConstraintOncall`:    stmt.TxCreateCheckConfigurationConstraintOncall,
		`CreateCheckConfigurationConstraintCustom`:    stmt.TxCreateCheckConfigurationConstraintCustom,
		`CreateCheckConfigurationConstraintService`:   stmt.TxCreateCheckConfigurationConstraintService,
		`CreateCheckConfigurationConstraintAttribute`: stmt.TxCreateCheckConfigurationConstraintAttribute,
	} {
		if stMap[name], err = tx.Prepare(stmt); err != nil {
			delete(stMap, name)
			goto bailout
		}
	}

	//
	//
	if stMap[`CreateBucket`], err = tx.Prepare(
		tkStmtCreateBucket,
	); err != nil {
		delete(stMap, `CreateBucket`)
		goto bailout
	}

	if stMap[`GroupCreate`], err = tx.Prepare(
		tkStmtGroupCreate,
	); err != nil {
		delete(stMap, `GroupCreate`)
		goto bailout
	}

	if stMap[`GroupUpdate`], err = tx.Prepare(
		tkStmtGroupUpdate,
	); err != nil {
		delete(stMap, `GroupUpdate`)
		goto bailout
	}

	if stMap[`GroupDelete`], err = tx.Prepare(
		tkStmtGroupDelete,
	); err != nil {
		delete(stMap, `GroupDelete`)
		goto bailout
	}

	if stMap[`GroupMemberNewNode`], err = tx.Prepare(
		tkStmtGroupMemberNewNode,
	); err != nil {
		delete(stMap, `GroupMemberNewNode`)
		goto bailout
	}

	if stMap[`GroupMemberNewCluster`], err = tx.Prepare(
		tkStmtGroupMemberNewCluster,
	); err != nil {
		delete(stMap, `GroupMemberNewCluster`)
		goto bailout
	}

	if stMap[`GroupMemberNewGroup`], err = tx.Prepare(
		tkStmtGroupMemberNewGroup,
	); err != nil {
		delete(stMap, `GroupMemberNewGroup`)
		goto bailout
	}

	if stMap[`GroupMemberRemoveNode`], err = tx.Prepare(
		tkStmtGroupMemberRemoveNode,
	); err != nil {
		delete(stMap, `GroupMemberRemoveNode`)
		goto bailout
	}

	if stMap[`GroupMemberRemoveCluster`], err = tx.Prepare(
		tkStmtGroupMemberRemoveCluster,
	); err != nil {
		delete(stMap, `GroupMemberRemoveCluster`)
		goto bailout
	}

	if stMap[`GroupMemberRemoveGroup`], err = tx.Prepare(
		tkStmtGroupMemberRemoveGroup,
	); err != nil {
		delete(stMap, `GroupMemberRemoveGroup`)
		goto bailout
	}

	if stMap[`ClusterCreate`], err = tx.Prepare(
		tkStmtClusterCreate,
	); err != nil {
		delete(stMap, `ClusterCreate`)
		goto bailout
	}

	if stMap[`ClusterUpdate`], err = tx.Prepare(
		tkStmtClusterUpdate,
	); err != nil {
		delete(stMap, `ClusterUpdate`)
		goto bailout
	}

	if stMap[`ClusterDelete`], err = tx.Prepare(
		tkStmtClusterDelete,
	); err != nil {
		delete(stMap, `ClusterDelete`)
		goto bailout
	}

	if stMap[`ClusterMemberNew`], err = tx.Prepare(
		tkStmtClusterMemberNew,
	); err != nil {
		delete(stMap, `ClusterMemberNew`)
		goto bailout
	}

	if stMap[`ClusterMemberRemove`], err = tx.Prepare(
		tkStmtClusterMemberRemove,
	); err != nil {
		delete(stMap, `ClusterMemberRemove`)
		goto bailout
	}

	if stMap[`BucketAssignNode`], err = tx.Prepare(
		tkStmtBucketAssignNode,
	); err != nil {
		delete(stMap, `BucketAssignNode`)
		goto bailout
	}

	if stMap[`UpdateNodeState`], err = tx.Prepare(
		tkStmtUpdateNodeState,
	); err != nil {
		delete(stMap, `UpdateNodeState`)
		goto bailout
	}

	if stMap[`NodeUnassignFromBucket`], err = tx.Prepare(
		tkStmtNodeUnassignFromBucket,
	); err != nil {
		delete(stMap, `NodeUnassignFromBucket`)
		goto bailout
	}

	return tx, stMap, nil

bailout:
	if open {
		defer tx.Rollback()
	}
	for _, stmt := range stMap {
		defer stmt.Close()
	}
	return nil, nil, err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
