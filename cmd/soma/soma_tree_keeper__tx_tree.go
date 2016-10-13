/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */
package main

func (tk *treeKeeper) txTree(a *tree.Action,
	stm *map[string]*sql.Stmt, user string) error {
	switch a.Action {
	case `create`:
		return tk.txTreeCreate(a, stm, user)
	case `update`:
		return tk.txTreeUpdate(a, stm)
	case `delete`:
		return tk.txTreeDelete(a, stm)
	}
}

func (tk *treeKeeper) txTreeCreate(a *tree.Action,
	stm *map[string]*sql.Stmt, user string) error {
	var err error
	switch a.Type {
	case `bucket`:
		_, err = stm[`CreateBucket`].Exec(
			a.Bucket.Id,
			a.Bucket.Name,
			a.Bucket.IsFrozen,
			a.Bucket.IsDeleted,
			a.Bucket.RepositoryId,
			a.Bucket.Environment,
			a.Bucket.TeamId,
			user,
		)
	case `group`:
		_, err = stm[`GroupCreate`].Exec(
			a.Group.Id,
			a.Group.BucketId,
			a.Group.Name,
			a.Group.ObjectState,
			a.Group.TeamId,
			user,
		)
	case `cluster`:
		_, err = stm[`ClusterCreate`].Exec(
			a.Cluster.Id,
			a.Cluster.Name,
			a.Cluster.BucketId,
			a.Cluster.ObjectState,
			a.Cluster.TeamId,
			user,
		)
	}
	return err
}

func (tk *treeKeeper) txTreeUpdate(a *tree.Action,
	stm *map[string]*sql.Stmt) error {
	var (
		err          error
		statement    *sql.Stmt
		id, newState string
	)
	switch a.Type {
	case `group`:
		statement = stm[`GroupUpdate`]
		id = a.Group.Id
		newState = a.Group.ObjectState
	case `cluster`:
		statement = stm[`ClusterUpdate`]
		id = a.Cluster.Id
		newState = a.Cluster.ObjectState
	case `node`:
		statement = stm[`UpdateNodeState`]
		id = a.Node.Id
		newState = a.Node.State
	}
	_, err = statement.Exec(
		id,
		newState,
	)
	return err
}

func (tk *treeKeeper) txTreeDelete(a *tree.Action,
	stm *map[string]*sql.Stmt) error {
	var err error
	switch a.Type {
	case `group`:
		_, err = stm[`GroupDelete`].Exec(
			a.Group.Id,
		)
	case `cluster`:
		_, err = stm[`ClusterDelete`].Exec(
			a.Cluster.Id,
		)
	case `node`:
		if _, err = stm[`NodeUnassignFromBucket`].Exec(
			a.Node.Id,
			a.Node.Config.BucketId,
			a.Node.TeamId,
		); err != nil {
			return err
		}
		// node unassign requires state update
		err = tk.txTreeUpdate(a, stm)
	}
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
