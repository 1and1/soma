/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/tree"
	"github.com/satori/go.uuid"
)

func (tk *TreeKeeper) treeBucket(q *msg.Request) {
	switch q.Action {
	case `create_bucket`:
		tree.NewBucket(tree.BucketSpec{
			Id:          uuid.NewV4().String(),
			Name:        q.Bucket.Name,
			Environment: q.Bucket.Environment,
			Team:        tk.meta.teamID,
			Deleted:     q.Bucket.IsDeleted,
			Frozen:      q.Bucket.IsFrozen,
			Repository:  q.Bucket.RepositoryId,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `repository`,
			ParentId:   tk.meta.repoID,
			ParentName: tk.meta.repoName,
		})
	}
}

func (tk *TreeKeeper) treeGroup(q *msg.Request) {
	switch q.Action {
	case `create_group`:
		tree.NewGroup(tree.GroupSpec{
			Id:   uuid.NewV4().String(),
			Name: q.Group.Name,
			Team: tk.meta.teamID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `bucket`,
			ParentId:   q.Group.BucketId,
		})
	case `delete_group`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `group`,
			ElementId:   q.Group.Id,
		}, true).(tree.BucketAttacher).Destroy()
	case `reset_group_to_bucket`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `group`,
			ElementId:   q.Group.Id,
		}, true).(tree.BucketAttacher).Detach()
	case `add_group_to_group`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `group`,
			ElementId:   (*q.Group.MemberGroups)[0].Id,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `group`,
			ParentId:   q.Group.Id,
		})
	}
}

func (tk *TreeKeeper) treeCluster(q *msg.Request) {
	switch q.Action {
	case `create_cluster`:
		tree.NewCluster(tree.ClusterSpec{
			Id:   uuid.NewV4().String(),
			Name: q.Cluster.Name,
			Team: tk.meta.teamID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `bucket`,
			ParentId:   q.Cluster.BucketId,
		})
	case `delete_cluster`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `cluster`,
			ElementId:   q.Cluster.Id,
		}, true).(tree.BucketAttacher).Destroy()
	case `reset_cluster_to_bucket`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `cluster`,
			ElementId:   q.Cluster.Id,
		}, true).(tree.BucketAttacher).Detach()
	case `add_cluster_to_group`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `cluster`,
			ElementId:   (*q.Group.MemberClusters)[0].Id,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `group`,
			ParentId:   q.Group.Id,
		})
	}
}

func (tk *TreeKeeper) treeNode(q *msg.Request) {
	switch q.Action {
	case `assign_node`:
		tree.NewNode(tree.NodeSpec{
			Id:       q.Node.Id,
			AssetId:  q.Node.AssetId,
			Name:     q.Node.Name,
			Team:     q.Node.TeamId,
			ServerId: q.Node.ServerId,
			Online:   q.Node.IsOnline,
			Deleted:  q.Node.IsDeleted,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `bucket`,
			ParentId:   q.Node.Config.BucketId,
		})
	case `delete_node`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementId:   q.Node.Id,
		}, true).(tree.BucketAttacher).Destroy()
	case `reset_node_to_bucket`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementId:   q.Node.Id,
		}, true).(tree.BucketAttacher).Detach()
	case `add_node_to_group`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementId:   (*q.Group.MemberNodes)[0].Id,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `group`,
			ParentId:   q.Group.Id,
		})
	case `add_node_to_cluster`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementId:   (*q.Cluster.Members)[0].Id,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `cluster`,
			ParentId:   q.Cluster.Id,
		})
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
