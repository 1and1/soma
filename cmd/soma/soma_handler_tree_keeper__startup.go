package main

import (
	"database/sql"
	"encoding/json"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/internal/tree"
)

func (tk *treeKeeper) startupLoad() {

	stMap := tk.prepareStartupStatements()
	for n, _ := range stMap {
		defer stMap[n].Close()
	}

	tk.startupBuckets(stMap)
	tk.startupGroups(stMap)
	tk.startupGroupMemberGroups(stMap)
	tk.startupGroupedClusters(stMap)
	tk.startupClusters(stMap)
	tk.startupNodes(stMap)

	if len(tk.actionChan) > 0 {
		tk.startLog.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	// attach system properties
	tk.startupSystemProperties(stMap)

	if len(tk.actionChan) > 0 {
		tk.startLog.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	// attach service properties
	tk.startupRepositoryServiceProperties(stMap)
	tk.startupBucketServiceProperties(stMap)
	tk.startupGroupServiceProperties(stMap)
	tk.startupClusterServiceProperties(stMap)
	tk.startupNodeServiceProperties(stMap)

	if len(tk.actionChan) > 0 {
		tk.startLog.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	// attach custom properties
	tk.startupRepositoryCustomProperties(stMap)
	tk.startupBucketCustomProperties(stMap)
	tk.startupGroupCustomProperties(stMap)
	tk.startupClusterCustomProperties(stMap)
	tk.startupNodeCustomProperties(stMap)

	if len(tk.actionChan) > 0 {
		tk.startLog.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	// attach oncall properties
	tk.startupRepositoryOncallProperties(stMap)
	tk.startupBucketOncallProperties(stMap)
	tk.startupGroupOncallProperties(stMap)
	tk.startupClusterOncallProperties(stMap)
	tk.startupNodeOncallProperties(stMap)

	if len(tk.actionChan) > 0 {
		tk.startLog.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	// attach checks
	tk.startupChecks(stMap)

	if !tk.rebuild && len(tk.actionChan) > 0 {
		tk.startLog.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	// these run as part of a job, but not inside the job's transaction. If there are leftovers
	// after a crash, fix them up
	if !SomaCfg.Observer {
		tk.buildDeploymentDetails()
		tk.orderDeploymentDetails()
	}

	// preload pending/unfinished jobs if not rebuilding the tree or
	// running in observer mode
	if !tk.rebuild && !SomaCfg.Observer {
		tk.startupJobs(stMap)
	}

	if !tk.rebuild && len(tk.actionChan) > 0 {
		tk.startLog.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}
}

func (tk *treeKeeper) startupBuckets(stMap map[string]*sql.Stmt) {
	if tk.broken {
		return
	}

	var (
		rows                                      *sql.Rows
		bucketId, bucketName, environment, teamId string
		frozen, deleted                           bool
		err                                       error
	)

	tk.startLog.Printf("TK[%s]: loading buckets\n", tk.repoName)
	rows, err = stMap[`LoadBucket`].Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading buckets: %s", tk.repoName, err.Error())
		tk.broken = true
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		tree.NewBucket(tree.BucketSpec{
			Id:          bucketId,
			Name:        bucketName,
			Environment: environment,
			Team:        teamId,
			Deleted:     deleted,
			Frozen:      frozen,
			Repository:  tk.repoId,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "repository",
			ParentId:   tk.repoId,
			ParentName: tk.repoName,
		})
		tk.drain(`action`)
		tk.drain(`error`)
	}
}

func (tk *treeKeeper) startupGroups(stMap map[string]*sql.Stmt) {
	if tk.broken {
		return
	}

	var (
		rows                                 *sql.Rows
		groupId, groupName, bucketId, teamId string
		err                                  error
	)

	tk.startLog.Printf("TK[%s]: loading groups\n", tk.repoName)
	rows, err = stMap[`LoadGroup`].Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading groups: %s", tk.repoName, err.Error())
		tk.broken = true
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
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

func (tk *treeKeeper) startupGroupMemberGroups(stMap map[string]*sql.Stmt) {
	if tk.broken {
		return
	}

	var (
		rows                  *sql.Rows
		groupId, childGroupId string
		err                   error
	)

	tk.startLog.Printf("TK[%s]: loading group-member-groups\n", tk.repoName)
	rows, err = stMap[`LoadGroupMbrGroup`].Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading groups: %s", tk.repoName, err.Error())
		tk.broken = true
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
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

func (tk *treeKeeper) startupGroupedClusters(stMap map[string]*sql.Stmt) {
	if tk.broken {
		return
	}

	var (
		err                                     error
		rows                                    *sql.Rows
		clusterId, clusterName, teamId, groupId string
	)

	tk.startLog.Printf("TK[%s]: loading grouped-clusters\n", tk.repoName)
	rows, err = stMap[`LoadGroupMbrCluster`].Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading clusters: %s", tk.repoName, err.Error())
		tk.broken = true
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
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

func (tk *treeKeeper) startupClusters(stMap map[string]*sql.Stmt) {
	if tk.broken {
		return
	}

	var (
		err                                      error
		rows                                     *sql.Rows
		clusterId, clusterName, bucketId, teamId string
	)

	tk.startLog.Printf("TK[%s]: loading clusters\n", tk.repoName)
	rows, err = stMap[`LoadCluster`].Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading clusters: %s", tk.repoName, err.Error())
		tk.broken = true
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
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

func (tk *treeKeeper) startupNodes(stMap map[string]*sql.Stmt) {
	if tk.broken {
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

	tk.startLog.Printf("TK[%s]: loading nodes\n", tk.repoName)
	rows, err = stMap[`LoadNode`].Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading nodes: %s", tk.repoName, err.Error())
		tk.broken = true
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
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

func (tk *treeKeeper) startupJobs(stMap map[string]*sql.Stmt) {
	if tk.broken {
		return
	}

	var (
		err  error
		rows *sql.Rows
		job  string
	)

	tk.startLog.Printf("TK[%s]: loading pending jobs\n", tk.repoName)
	rows, err = stMap[`LoadJob`].Query(tk.repoId)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading clusters: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

jobloop:
	for rows.Next() {
		err = rows.Scan(
			&job,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break jobloop
			}
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		tr := treeRequest{}
		err = json.Unmarshal([]byte(job), &tr)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		tk.input <- tr
		tk.startLog.Printf("TK[%s] Loaded job %s (%s)\n", tk.repoName, tr.JobId, tr.Action)
	}
}

func (tk *treeKeeper) prepareStartupStatements() map[string]*sql.Stmt {
	var err error
	stMap := map[string]*sql.Stmt{}

	for name, statement := range map[string]string{
		`LoadBucket`:             stmt.TkStartLoadBuckets,
		`LoadGroup`:              stmt.TkStartLoadGroups,
		`LoadGroupMbrGroup`:      stmt.TkStartLoadGroupMemberGroups,
		`LoadGroupMbrCluster`:    stmt.TkStartLoadGroupedClusters,
		`LoadCluster`:            stmt.TkStartLoadCluster,
		`LoadNode`:               stmt.TkStartLoadNode,
		`LoadJob`:                stmt.TkStartLoadJob,
		`LoadChecks`:             stmt.TkStartLoadChecks,
		`LoadItems`:              stmt.TkStartLoadInheritedChecks,
		`LoadConfig`:             stmt.TkStartLoadCheckConfiguration,
		`LoadAllConfigsForType`:  stmt.TkStartLoadAllCheckConfigurationsForType,
		`LoadThreshold`:          stmt.TkStartLoadCheckThresholds,
		`LoadCustomCstr`:         stmt.TkStartLoadCheckConstraintCustom,
		`LoadNativeCstr`:         stmt.TkStartLoadCheckConstraintNative,
		`LoadOncallCstr`:         stmt.TkStartLoadCheckConstraintOncall,
		`LoadAttributeCstr`:      stmt.TkStartLoadCheckConstraintAttribute,
		`LoadServiceCstr`:        stmt.TkStartLoadCheckConstraintService,
		`LoadSystemCstr`:         stmt.TkStartLoadCheckConstraintSystem,
		`LoadChecksForType`:      stmt.TkStartLoadChecksForType,
		`LoadInstances`:          stmt.TkStartLoadCheckInstances,
		`LoadInstanceCfg`:        stmt.TkStartLoadCheckInstanceConfiguration,
		`LoadGroupState`:         stmt.TkStartLoadCheckGroupState,
		`LoadGroupRelations`:     stmt.TkStartLoadCheckGroupRelations,
		`CapabilityView`:         stmt.TreekeeperGetViewFromCapability,
		`LoadPropRepoSystem`:     stmt.TkStartLoadRepoSysProp,
		`LoadPropBuckSystem`:     stmt.TkStartLoadBucketSysProp,
		`LoadPropGrpSystem`:      stmt.TkStartLoadGroupSysProp,
		`LoadPropClrSystem`:      stmt.TkStartLoadClusterSysProp,
		`LoadPropNodeSystem`:     stmt.TkStartLoadNodeSysProp,
		`LoadPropSystemInstance`: stmt.TkStartLoadSystemPropInstances,
		`LoadPropRepoCustom`:     stmt.TkStartLoadRepositoryCstProp,
		`LoadPropBuckCustom`:     stmt.TkStartLoadBucketCstProp,
		`LoadPropGrpCustom`:      stmt.TkStartLoadGroupCstProp,
		`LoadPropClrCustom`:      stmt.TkStartLoadClusterCstProp,
		`LoadPropNodeCustom`:     stmt.TkStartLoadNodeCstProp,
		`LoadPropCustomInstance`: stmt.TkStartLoadCustomPropInstances,
		`LoadPropRepoOncall`:     stmt.TkStartLoadRepoOncProp,
		`LoadPropBuckOncall`:     stmt.TkStartLoadBucketOncProp,
		`LoadPropGrpOncall`:      stmt.TkStartLoadGroupOncProp,
		`LoadPropClrOncall`:      stmt.TkStartLoadClusterOncProp,
		`LoadPropNodeOncall`:     stmt.TkStartLoadNodeOncProp,
		`LoadPropOncallInstance`: stmt.TkStartLoadOncallPropInstances,
		`LoadPropRepoService`:    stmt.TkStartLoadRepoSvcProp,
		`LoadPropBuckService`:    stmt.TkStartLoadBucketSvcProp,
		`LoadPropGrpService`:     stmt.TkStartLoadGroupSvcProp,
		`LoadPropClrService`:     stmt.TkStartLoadClusterSvcProp,
		`LoadPropNodeService`:    stmt.TkStartLoadNodeSvcProp,
		`LoadPropRepoSvcAttr`:    stmt.TkStartLoadRepoSvcAttr,
		`LoadPropBuckSvcAttr`:    stmt.TkStartLoadBucketSvcAttr,
		`LoadPropGrpSvcAttr`:     stmt.TkStartLoadGroupSvcAttr,
		`LoadPropClrSvcAttr`:     stmt.TkStartLoadClusterSvcAttr,
		`LoadPropNodeSvcAttr`:    stmt.TkStartLoadNodeSvcAttr,
		`LoadPropSvcInstance`:    stmt.TkStartLoadServicePropInstances,
	} {
		if stMap[name], err = tk.conn.Prepare(statement); err != nil {
			tk.startLog.Println(`treekeeper startup`, err,
				stmt.Name(statement))
			tk.broken = true
			return map[string]*sql.Stmt{}
		}
	}
	return stMap
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
