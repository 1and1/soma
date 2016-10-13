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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
