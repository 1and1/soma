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
	var err error
	switch a.Type {
	case `repository`:
		_, err = stm[`RepositoryPropertyCustomCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Property.Custom.RepositoryId,
			a.Property.View,
			a.Property.Custom.Id,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
			a.Property.Custom.Value,
		)
	case `bucket`:
		_, err = stm[`BucketPropertyCustomCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Bucket.Id,
			a.Property.View,
			a.Property.Custom.Id,
			a.Property.Custom.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
			a.Property.Custom.Value,
		)
	case `group`:
		_, err = stm[`GroupPropertyCustomCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Group.Id,
			a.Property.View,
			a.Property.Custom.Id,
			a.Property.BucketId,
			a.Property.Custom.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
			a.Property.Custom.Value,
		)
	case `cluster`:
		_, err = stm[`ClusterPropertyCustomCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Cluster.Id,
			a.Property.View,
			a.Property.Custom.Id,
			a.Property.BucketId,
			a.Property.Custom.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
			a.Property.Custom.Value,
		)
	case `node`:
		_, err = stm[`NodePropertyCustomCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Node.Id,
			a.Property.View,
			a.Property.Custom.Id,
			a.Property.BucketId,
			a.Property.Custom.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
			a.Property.Custom.Value,
		)
	}
	return err
}

func (tk *treeKeeper) txPropertyNewSystem(a *tree.Action,
	stm *map[string]*sql.Stmt) error {
	var err error
	switch a.Type {
	case `repository`:
		_, err = stm[`RepositoryPropertySystemCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Repository.Id,
			a.Property.View,
			a.Property.System.Name,
			a.Property.SourceType,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
			a.Property.System.Value,
			a.Property.IsInherited,
		)
	case `bucket`:
		_, err = stm[`BucketPropertySystemCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Bucket.Id,
			a.Property.View,
			a.Property.System.Name,
			a.Property.SourceType,
			a.Property.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
			a.Property.System.Value,
			a.Property.IsInherited,
		)
	case `group`:
		_, err = stm[`GroupPropertySystemCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Group.Id,
			a.Property.View,
			a.Property.System.Name,
			a.Property.SourceType,
			a.Property.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
			a.Property.System.Value,
			a.Property.IsInherited,
		)
	case `cluster`:
		_, err = stm[`ClusterPropertySystemCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Cluster.Id,
			a.Property.View,
			a.Property.System.Name,
			a.Property.SourceType,
			a.Property.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
			a.Property.System.Value,
			a.Property.IsInherited,
		)
	case `node`:
		_, err = stm[`NodePropertySystemCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Node.Id,
			a.Property.View,
			a.Property.System.Name,
			a.Property.SourceType,
			a.Property.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
			a.Property.System.Value,
			a.Property.IsInherited,
		)
	}
	return err
}

func (tk *treeKeeper) txPropertyNewService(a *tree.Action,
	stm *map[string]*sql.Stmt) error {
	var err error
	switch a.Type {
	case `repository`:
		_, err = stm[`RepositoryPropertyServiceCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Repository.Id,
			a.Property.View,
			a.Property.Service.Name,
			a.Property.Service.TeamId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
		)
	case `bucket`:
		_, err = stm[`BucketPropertyServiceCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Bucket.Id,
			a.Property.View,
			a.Property.Service.Name,
			a.Property.Service.TeamId,
			a.Property.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
		)
	case `group`:
		_, err = stm[`GroupPropertyServiceCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Group.Id,
			a.Property.View,
			a.Property.Service.Name,
			a.Property.Service.TeamId,
			a.Property.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
		)
	case `cluster`:
		_, err = stm[`ClusterPropertyServiceCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Cluster.Id,
			a.Property.View,
			a.Property.Service.Name,
			a.Property.Service.TeamId,
			a.Property.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
		)
	case `node`:
		_, err = stm[`NodePropertyServiceCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Node.Id,
			a.Property.View,
			a.Property.Service.Name,
			a.Property.Service.TeamId,
			a.Property.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
		)
	}
	return err
}

func (tk *treeKeeper) txPropertyNewOncall(a *tree.Action,
	stm *map[string]*sql.Stmt) error {
	var err error
	switch a.Type {
	case `repository`:
		_, err = stm[`RepositoryPropertyOncallCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Repository.Id,
			a.Property.View,
			a.Property.Oncall.Id,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
		)
	case `bucket`:
		_, err = stm[`BucketPropertyOncallCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Bucket.Id,
			a.Property.View,
			a.Property.Oncall.Id,
			a.Property.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
		)
	case `group`:
		_, err = stm[`GroupPropertyOncallCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Group.Id,
			a.Property.View,
			a.Property.Oncall.Id,
			a.Property.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
		)
	case `cluster`:
		_, err = stm[`ClusterPropertyOncallCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Cluster.Id,
			a.Property.View,
			a.Property.Oncall.Id,
			a.Property.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
		)
	case `node`:
		_, err = stm[`NodePropertyOncallCreate`].Exec(
			a.Property.InstanceId,
			a.Property.SourceInstanceId,
			a.Node.Id,
			a.Property.View,
			a.Property.Oncall.Id,
			a.Property.RepositoryId,
			a.Property.Inheritance,
			a.Property.ChildrenOnly,
		)
	}
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
