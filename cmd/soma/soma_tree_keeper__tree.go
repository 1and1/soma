/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"github.com/1and1/soma/internal/tree"
	"github.com/satori/go.uuid"
)

func (tk *treeKeeper) treeBucket(q *treeRequest) {
	switch q.Action {
	case `create_bucket`:
		tree.NewBucket(tree.BucketSpec{
			Id:          uuid.NewV4().String(),
			Name:        q.Bucket.Bucket.Name,
			Environment: q.Bucket.Bucket.Environment,
			Team:        tk.team,
			Deleted:     q.Bucket.Bucket.IsDeleted,
			Frozen:      q.Bucket.Bucket.IsFrozen,
			Repository:  q.Bucket.Bucket.RepositoryId,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `repository`,
			ParentId:   tk.repoId,
			ParentName: tk.repoName,
		})
	}
}

func (tk *treeKeeper) treeGroup(q *treeRequest) {
	switch q.Action {
	case `create_group`:
		tree.NewGroup(tree.GroupSpec{
			Id:   uuid.NewV4().String(),
			Name: q.Group.Group.Name,
			Team: tk.team,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `bucket`,
			ParentId:   q.Group.Group.BucketId,
		})
	case `delete_group`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `group`,
			ElementId:   q.Group.Group.Id,
		}, true).(tree.BucketAttacher).Destroy()
	case `reset_group_to_bucket`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `group`,
			ElementId:   q.Group.Group.Id,
		}, true).(tree.BucketAttacher).Detach()
	case `add_group_to_group`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `group`,
			ElementId:   (*q.Group.Group.MemberGroups)[0].Id,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `group`,
			ParentId:   q.Group.Group.Id,
		})
	}
}

func (tk *treeKeeper) treeCluster(q *treeRequest) {
	switch q.Action {
	case `create_cluster`:
		tree.NewCluster(tree.ClusterSpec{
			Id:   uuid.NewV4().String(),
			Name: q.Cluster.Cluster.Name,
			Team: tk.team,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `bucket`,
			ParentId:   q.Cluster.Cluster.BucketId,
		})
	case `delete_cluster`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `cluster`,
			ElementId:   q.Cluster.Cluster.Id,
		}, true).(tree.BucketAttacher).Destroy()
	case `reset_cluster_to_bucket`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `cluster`,
			ElementId:   q.Cluster.Cluster.Id,
		}, true).(tree.BucketAttacher).Detach()
	case `add_cluster_to_group`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `cluster`,
			ElementId:   (*q.Group.Group.MemberClusters)[0].Id,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `group`,
			ParentId:   q.Group.Group.Id,
		})
	}
}

func (tk *treeKeeper) treeNode(q *treeRequest) {
	switch q.Action {
	case `assign_node`:
		tree.NewNode(tree.NodeSpec{
			Id:       q.Node.Node.Id,
			AssetId:  q.Node.Node.AssetId,
			Name:     q.Node.Node.Name,
			Team:     q.Node.Node.TeamId,
			ServerId: q.Node.Node.ServerId,
			Online:   q.Node.Node.IsOnline,
			Deleted:  q.Node.Node.IsDeleted,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `bucket`,
			ParentId:   q.Node.Node.Config.BucketId,
		})
	case `delete_node`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementId:   q.Node.Node.Id,
		}, true).(tree.BucketAttacher).Destroy()
	case `reset_node_to_bucket`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementId:   q.Node.Node.Id,
		}, true).(tree.BucketAttacher).Detach()
	case `add_node_to_group`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementId:   (*q.Group.Group.MemberNodes)[0].Id,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `group`,
			ParentId:   q.Group.Group.Id,
		})
	case `add_node_to_cluster`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementId:   (*q.Cluster.Cluster.Members)[0].Id,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `cluster`,
			ParentId:   q.Cluster.Cluster.Id,
		})
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
