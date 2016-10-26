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
		`PropertyInstanceCreate`:          stmt.TxPropertyInstanceCreate,
		`PropertyInstanceDelete`:          stmt.TxPropertyInstanceDelete,
		`RepositoryPropertyOncallCreate`:  stmt.TxRepositoryPropertyOncallCreate,
		`RepositoryPropertyOncallDelete`:  stmt.TxRepositoryPropertyOncallDelete,
		`RepositoryPropertyServiceCreate`: stmt.TxRepositoryPropertyServiceCreate,
		`RepositoryPropertyServiceDelete`: stmt.TxRepositoryPropertyServiceDelete,
		`RepositoryPropertySystemCreate`:  stmt.TxRepositoryPropertySystemCreate,
		`RepositoryPropertySystemDelete`:  stmt.TxRepositoryPropertySystemDelete,
		`RepositoryPropertyCustomCreate`:  stmt.TxRepositoryPropertyCustomCreate,
		`RepositoryPropertyCustomDelete`:  stmt.TxRepositoryPropertyCustomDelete,
		`BucketPropertyOncallCreate`:      stmt.TxBucketPropertyOncallCreate,
		`BucketPropertyOncallDelete`:      stmt.TxBucketPropertyOncallDelete,
		`BucketPropertyServiceCreate`:     stmt.TxBucketPropertyServiceCreate,
		`BucketPropertyServiceDelete`:     stmt.TxBucketPropertyServiceDelete,
		`BucketPropertySystemCreate`:      stmt.TxBucketPropertySystemCreate,
		`BucketPropertySystemDelete`:      stmt.TxBucketPropertySystemDelete,
		`BucketPropertyCustomCreate`:      stmt.TxBucketPropertyCustomCreate,
		`BucketPropertyCustomDelete`:      stmt.TxBucketPropertyCustomDelete,
		`GroupPropertyOncallCreate`:       stmt.TxGroupPropertyOncallCreate,
		`GroupPropertyOncallDelete`:       stmt.TxGroupPropertyOncallDelete,
		`GroupPropertyServiceCreate`:      stmt.TxGroupPropertyServiceCreate,
		`GroupPropertyServiceDelete`:      stmt.TxGroupPropertyServiceDelete,
		`GroupPropertySystemCreate`:       stmt.TxGroupPropertySystemCreate,
		`GroupPropertySystemDelete`:       stmt.TxGroupPropertySystemDelete,
		`GroupPropertyCustomCreate`:       stmt.TxGroupPropertyCustomCreate,
		`GroupPropertyCustomDelete`:       stmt.TxGroupPropertyCustomDelete,
		`ClusterPropertyOncallCreate`:     stmt.TxClusterPropertyOncallCreate,
		`ClusterPropertyOncallDelete`:     stmt.TxClusterPropertyOncallDelete,
		`ClusterPropertyServiceCreate`:    stmt.TxClusterPropertyServiceCreate,
		`ClusterPropertyServiceDelete`:    stmt.TxClusterPropertyServiceDelete,
		`ClusterPropertySystemCreate`:     stmt.TxClusterPropertySystemCreate,
		`ClusterPropertySystemDelete`:     stmt.TxClusterPropertySystemDelete,
		`ClusterPropertyCustomCreate`:     stmt.TxClusterPropertyCustomCreate,
		`ClusterPropertyCustomDelete`:     stmt.TxClusterPropertyCustomDelete,
		`NodePropertyOncallCreate`:        stmt.TxNodePropertyOncallCreate,
		`NodePropertyOncallDelete`:        stmt.TxNodePropertyOncallDelete,
		`NodePropertyServiceCreate`:       stmt.TxNodePropertyServiceCreate,
		`NodePropertyServiceDelete`:       stmt.TxNodePropertyServiceDelete,
		`NodePropertySystemCreate`:        stmt.TxNodePropertySystemCreate,
		`NodePropertySystemDelete`:        stmt.TxNodePropertySystemDelete,
		`NodePropertyCustomCreate`:        stmt.TxNodePropertyCustomCreate,
		`NodePropertyCustomDelete`:        stmt.TxNodePropertyCustomDelete,
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
		`BucketAssignNode`:         stmt.TxBucketAssignNode,
		`ClusterCreate`:            stmt.TxClusterCreate,
		`ClusterDelete`:            stmt.TxClusterDelete,
		`ClusterMemberNew`:         stmt.TxClusterMemberNew,
		`ClusterMemberRemove`:      stmt.TxClusterMemberRemove,
		`ClusterUpdate`:            stmt.TxClusterUpdate,
		`CreateBucket`:             stmt.TxCreateBucket,
		`GroupCreate`:              stmt.TxGroupCreate,
		`GroupDelete`:              stmt.TxGroupDelete,
		`GroupMemberNewCluster`:    stmt.TxGroupMemberNewCluster,
		`GroupMemberNewGroup`:      stmt.TxGroupMemberNewGroup,
		`GroupMemberNewNode`:       stmt.TxGroupMemberNewNode,
		`GroupMemberRemoveCluster`: stmt.TxGroupMemberRemoveCluster,
		`GroupMemberRemoveGroup`:   stmt.TxGroupMemberRemoveGroup,
		`GroupMemberRemoveNode`:    stmt.TxGroupMemberRemoveNode,
		`GroupUpdate`:              stmt.TxGroupUpdate,
		`NodeUnassignFromBucket`:   stmt.TxNodeUnassignFromBucket,
		`UpdateNodeState`:          stmt.TxUpdateNodeState,
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
		// if the transaction was opened, then tx.Rollback() will close all
		// prepared statements. If the transaction was not opened yet, then
		// no statements have been prepared inside it - there is nothing to
		// close
		defer tx.Rollback()
	}
	return nil, nil, err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
