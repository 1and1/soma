package soma

import (
	"database/sql"

	"github.com/1and1/soma/internal/tree"
)

func (tk *TreeKeeper) startupBuckets(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		rows                                      *sql.Rows
		bucketId, bucketName, environment, teamId string
		frozen, deleted                           bool
		err                                       error
	)

	tk.startLog.Printf("TK[%s]: loading buckets\n", tk.meta.repoName)
	rows, err = stMap[`LoadBucket`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading buckets: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

bucketloop:
	for rows.Next() {
		err = rows.Scan(
			&bucketId,
			&bucketName,
			&frozen,
			&deleted,
			&environment,
			&teamId,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break bucketloop
			}
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}
		tree.NewBucket(tree.BucketSpec{
			Id:          bucketId,
			Name:        bucketName,
			Environment: environment,
			Team:        teamId,
			Deleted:     deleted,
			Frozen:      frozen,
			Repository:  tk.meta.repoID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "repository",
			ParentId:   tk.meta.repoID,
			ParentName: tk.meta.repoName,
		})
		tk.drain(`action`)
		tk.drain(`error`)
	}
}

func (tk *TreeKeeper) startupGroups(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		rows                                 *sql.Rows
		groupId, groupName, bucketId, teamId string
		err                                  error
	)

	tk.startLog.Printf("TK[%s]: loading groups\n", tk.meta.repoName)
	rows, err = stMap[`LoadGroup`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading groups: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

grouploop:
	for rows.Next() {
		err = rows.Scan(
			&groupId,
			&groupName,
			&bucketId,
			&teamId,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break grouploop
			}
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}
		tree.NewGroup(tree.GroupSpec{
			Id:   groupId,
			Name: groupName,
			Team: teamId,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "bucket",
			ParentId:   bucketId,
		})
		tk.drain(`action`)
		tk.drain(`error`)
	}
}

func (tk *TreeKeeper) startupGroupMemberGroups(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		rows                  *sql.Rows
		groupId, childGroupId string
		err                   error
	)

	tk.startLog.Printf("TK[%s]: loading group-member-groups\n", tk.meta.repoName)
	rows, err = stMap[`LoadGroupMbrGroup`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading groups: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

memberloop:
	for rows.Next() {
		err = rows.Scan(
			&groupId,
			&childGroupId,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break memberloop
			}
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}

		tk.tree.Find(tree.FindRequest{
			ElementType: "group",
			ElementId:   childGroupId,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentId:   groupId,
		})
	}
	tk.drain(`action`)
	tk.drain(`error`)
}

func (tk *TreeKeeper) startupGroupedClusters(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		err                                     error
		rows                                    *sql.Rows
		clusterId, clusterName, teamId, groupId string
	)

	tk.startLog.Printf("TK[%s]: loading grouped-clusters\n", tk.meta.repoName)
	rows, err = stMap[`LoadGroupMbrCluster`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading clusters: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

clusterloop:
	for rows.Next() {
		err = rows.Scan(
			&clusterId,
			&clusterName,
			&teamId,
			&groupId,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break clusterloop
			}
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}

		tree.NewCluster(tree.ClusterSpec{
			Id:   clusterId,
			Name: clusterName,
			Team: teamId,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentId:   groupId,
		})
	}
	tk.drain(`action`)
	tk.drain(`error`)
}

func (tk *TreeKeeper) startupClusters(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		err                                      error
		rows                                     *sql.Rows
		clusterId, clusterName, bucketId, teamId string
	)

	tk.startLog.Printf("TK[%s]: loading clusters\n", tk.meta.repoName)
	rows, err = stMap[`LoadCluster`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading clusters: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

clusterloop:
	for rows.Next() {
		err = rows.Scan(
			&clusterId,
			&clusterName,
			&bucketId,
			&teamId,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break clusterloop
			}
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}

		tree.NewCluster(tree.ClusterSpec{
			Id:   clusterId,
			Name: clusterName,
			Team: teamId,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "bucket",
			ParentId:   bucketId,
		})
	}
	tk.drain(`action`)
	tk.drain(`error`)
}

func (tk *TreeKeeper) startupNodes(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		err                                          error
		rows                                         *sql.Rows
		nodeId, nodeName, teamId, serverId, bucketId string
		assetId                                      int
		nodeOnline, nodeDeleted                      bool
		clusterId, groupId                           sql.NullString
	)

	tk.startLog.Printf("TK[%s]: loading nodes\n", tk.meta.repoName)
	rows, err = stMap[`LoadNode`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading nodes: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

nodeloop:
	for rows.Next() {
		err = rows.Scan(
			&nodeId,
			&assetId,
			&nodeName,
			&teamId,
			&serverId,
			&nodeOnline,
			&nodeDeleted,
			&bucketId,
			&clusterId,
			&groupId,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break nodeloop
			}
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}

		node := tree.NewNode(tree.NodeSpec{
			Id:       nodeId,
			AssetId:  uint64(assetId),
			Name:     nodeName,
			Team:     teamId,
			ServerId: serverId,
			Online:   nodeOnline,
			Deleted:  nodeDeleted,
		})
		if clusterId.Valid {
			node.Attach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: "cluster",
				ParentId:   clusterId.String,
			})
		} else if groupId.Valid {
			node.Attach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: "group",
				ParentId:   groupId.String,
			})
		} else {
			node.Attach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: "bucket",
				ParentId:   bucketId,
			})
		}
	}
	tk.drain(`action`)
	tk.drain(`error`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
