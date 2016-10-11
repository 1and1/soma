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

	"github.com/1and1/soma/internal/tree"
)

func (tk *treeKeeper) txProperty(a *tree.Action,
	stm *map[string]*sql.Stmt) error {
	switch a.Action {
	case `property_new`:
		return tk.txPropertyNew(a, stm)
	}
}

//
// PROPERTY NEW
func (tk *treeKeeper) txPropertyNew(a *tree.Action,
	stm *map[string]*sql.Stmt) error {
	if err := stm[`PropertyInstanceCreate`].Exec(
		a.Property.InstanceId,
		a.Property.RepositoryId,
		a.Property.SourceInstanceId,
		a.Property.SourceType,
		a.Property.InheritedFrom,
	); err != nil {
		return err
	}

	switch a.Property.Type {
	case `custom`:
		return tk.txPropertyNewCustom(a, stm)
	case `system`:
		return tk.txPropertyNewSystem(a, stm)
	case `service`:
		return tk.txPropertyNewService(a, stm)
	case `oncall`:
		return tk.txPropertyNewOncall(a, stm)
	}
	return fmt.Errorf(`Impossible property type`)
}

func (tk *treeKeeper) txPropertyNewCustom(a *tree.Action,
	stm *map[string]*sql.Stmt) error {
	var (
		err       error
		statement *sql.Stmt
		id        string
	)
	switch a.Type {
	case `repository`:
		statement = stm[`RepositoryPropertyCustomCreate`]
		id = a.Property.Custom.RepositoryId
	case `bucket`:
		statement = stm[`BucketPropertyCustomCreate`]
		id = a.Bucket.Id
	case `group`:
		statement = stm[`GroupPropertyCustomCreate`]
		id = a.Group.Id
	case `cluster`:
		statement = stm[`ClusterPropertyCustomCreate`]
		id = a.Cluster.Id
	case `node`:
		statement = stm[`NodePropertyCustomCreate`]
		id = a.Node.Id
	}
	_, err = statement.Exec(
		a.Property.InstanceId,
		a.Property.SourceInstanceId,
		id,
		a.Property.View,
		a.Property.Custom.Id,
		a.Property.Inheritance,
		a.Property.ChildrenOnly,
		a.Property.Custom.Value,
	)
	return err
}

func (tk *treeKeeper) txPropertyNewSystem(a *tree.Action,
	stm *map[string]*sql.Stmt) error {
	var (
		err       error
		statement *sql.Stmt
		id        string
	)
	switch a.Type {
	case `repository`:
		statement = stm[`RepositoryPropertySystemCreate`]
		id = a.Repository.Id
	case `bucket`:
		statement = stm[`BucketPropertySystemCreate`]
		id = a.Bucket.Id
	case `group`:
		statement = stm[`GroupPropertySystemCreate`]
		id = a.Group.Id
	case `cluster`:
		statement = stm[`ClusterPropertySystemCreate`]
		id = a.Cluster.Id
	case `node`:
		statement = stm[`NodePropertySystemCreate`]
		id = a.Node.Id
	}
	_, err = statement.Exec(
		a.Property.InstanceId,
		a.Property.SourceInstanceId,
		id,
		a.Property.View,
		a.Property.System.Name,
		a.Property.SourceType,
		a.Property.RepositoryId,
		a.Property.Inheritance,
		a.Property.ChildrenOnly,
		a.Property.System.Value,
		a.Property.IsInherited,
	)
	return err
}

func (tk *treeKeeper) txPropertyNewService(a *tree.Action,
	stm *map[string]*sql.Stmt) error {
	var (
		err       error
		statement *sql.Stmt
		id        string
	)
	switch a.Type {
	case `repository`:
		statement = stm[`RepositoryPropertyServiceCreate`]
		id = a.Repository.Id
	case `bucket`:
		statement = stm[`BucketPropertyServiceCreate`]
		id = a.Bucket.Id
	case `group`:
		statement = stm[`GroupPropertyServiceCreate`]
		id = a.Group.Id
	case `cluster`:
		statement = stm[`ClusterPropertyServiceCreate`]
		id = a.Cluster.Id
	case `node`:
		statement = stm[`NodePropertyServiceCreate`]
		id = a.Node.Id
	}
	_, err = statement.Exec(
		a.Property.InstanceId,
		a.Property.SourceInstanceId,
		id,
		a.Property.View,
		a.Property.Service.Name,
		a.Property.Service.TeamId,
		a.Property.RepositoryId,
		a.Property.Inheritance,
		a.Property.ChildrenOnly,
	)
	return err
}

func (tk *treeKeeper) txPropertyNewOncall(a *tree.Action,
	stm *map[string]*sql.Stmt) error {
	var (
		err       error
		statement *sql.Stmt
		id        string
	)
	switch a.Type {
	case `repository`:
		statement = stm[`RepositoryPropertyOncallCreate`]
		id = a.Repository.Id
	case `bucket`:
		statement = stm[`BucketPropertyOncallCreate`]
		id = a.Bucket.Id
	case `group`:
		statement = stm[`GroupPropertyOncallCreate`]
		id = a.Group.Id
	case `cluster`:
		statement = stm[`ClusterPropertyOncallCreate`]
		id = a.Cluster.Id
	case `node`:
		statement = stm[`NodePropertyOncallCreate`]
		id = a.Node.Id
	}
	_, err = statement.Exec(
		a.Property.InstanceId,
		a.Property.SourceInstanceId,
		id,
		a.Property.View,
		a.Property.Oncall.Id,
		a.Property.RepositoryId,
		a.Property.Inheritance,
		a.Property.ChildrenOnly,
	)
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
