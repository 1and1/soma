package main

import (
	"database/sql"
	"encoding/json"
	"log"

)

func (tk *treeKeeper) startupLoad() {

	tk.startupBuckets()
	tk.startupGroups()
	tk.startupGroupMemberGroups()
	tk.startupGroupedClusters()
	tk.startupClusters()
	tk.startupNodes()

	if len(tk.actionChan) > 0 {
		log.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	// attach system properties
	tk.startupRepositorySystemProperties()
	tk.startupBucketSystemProperties()
	tk.startupGroupSystemProperties()
	tk.startupClusterSystemProperties()
	tk.startupNodeSystemProperties()

	if len(tk.actionChan) > 0 {
		log.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	// attach service properties
	tk.startupRepositoryServiceProperties()
	tk.startupBucketServiceProperties()
	tk.startupGroupServiceProperties()
	tk.startupClusterServiceProperties()
	tk.startupNodeServiceProperties()

	if len(tk.actionChan) > 0 {
		log.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	// attach custom properties
	tk.startupRepositoryCustomProperties()
	tk.startupBucketCustomProperties()
	tk.startupGroupCustomProperties()
	tk.startupClusterCustomProperties()
	tk.startupNodeCustomProperties()

	if len(tk.actionChan) > 0 {
		log.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	// attach oncall properties
	tk.startupRepositoryOncallProperties()
	tk.startupBucketOncallProperties()
	tk.startupGroupOncallProperties()
	tk.startupClusterOncallProperties()
	tk.startupNodeOncallProperties()

	if len(tk.actionChan) > 0 {
		log.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	// attach checks
	tk.startupChecks()

	if !tk.rebuild && len(tk.actionChan) > 0 {
		log.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	// these run as part of a job, but not inside the job's transaction. If there are leftovers
	// after a crash, fix them up
	tk.buildDeploymentDetails()
	tk.orderDeploymentDetails()

	// preload pending/unfinished jobs if not rebuilding the tree or
	// running in observer mode
	if !tk.rebuild && !SomaCfg.Observer {
		tk.startupJobs()
	}

	if !tk.rebuild && len(tk.actionChan) > 0 {
		log.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	// XXX DEBUG: enable/disable dumping JSON of the entire tree after startup
	//b, _ := json.Marshal(tk.tree)
	//log.Println(string(b))
}

func (tk *treeKeeper) startupBuckets() {
	if tk.broken {
		return
	}

	var (
		rows                                      *sql.Rows
		bucketId, bucketName, environment, teamId string
		frozen, deleted                           bool
		err                                       error
		load_bucket                               *sql.Stmt
	)
	load_bucket, err = tk.conn.Prepare(`
SELECT sb.bucket_id,
       sb.bucket_name,
       sb.bucket_frozen,
       sb.bucket_deleted,
       sb.environment,
       sb.organizational_team_id
FROM   soma.repositories sr
JOIN   soma.buckets sb
ON     sr.repository_id = sb.repository_id
WHERE  sr.repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-buckets: ", err)
	}
	defer load_bucket.Close()

	log.Printf("TK[%s]: loading buckets\n", tk.repoName)
	rows, err = load_bucket.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading buckets: %s", tk.repoName, err.Error())
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
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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

func (tk *treeKeeper) startupGroups() {
	if tk.broken {
		return
	}

	var (
		rows                                 *sql.Rows
		groupId, groupName, bucketId, teamId string
		err                                  error
		load_group                           *sql.Stmt
	)
	load_group, err = tk.conn.Prepare(`
SELECT sg.group_id,
       sg.group_name,
       sg.bucket_id,
       sg.organizational_team_id
FROM   soma.repositories sr
JOIN   soma.buckets sb
ON     sr.repository_id = sb.repository_id
JOIN   soma.groups sg
ON     sb.bucket_id = sg.bucket_id
WHERE  sr.repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-groups: ", err)
	}
	defer load_group.Close()

	log.Printf("TK[%s]: loading groups\n", tk.repoName)
	rows, err = load_group.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading groups: %s", tk.repoName, err.Error())
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
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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

func (tk *treeKeeper) startupGroupMemberGroups() {
	if tk.broken {
		return
	}

	var (
		rows                  *sql.Rows
		groupId, childGroupId string
		err                   error
		load_grp_mbr_grp      *sql.Stmt
	)
	load_grp_mbr_grp, err = tk.conn.Prepare(`
SELECT sgmg.group_id,
       sgmg.child_group_id
FROM   soma.repositories sr
JOIN   soma.buckets sb
ON     sr.repository_id = sb.repository_id
JOIN   soma.group_membership_groups sgmg
ON     sb.bucket_id = sgmg.bucket_id
WHERE  sr.repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-group-member-groups: ", err)
	}
	defer load_grp_mbr_grp.Close()

	log.Printf("TK[%s]: loading group-member-groups\n", tk.repoName)
	rows, err = load_grp_mbr_grp.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading groups: %s", tk.repoName, err.Error())
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
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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

func (tk *treeKeeper) startupGroupedClusters() {
	if tk.broken {
		return
	}

	var (
		err                                     error
		rows                                    *sql.Rows
		clusterId, clusterName, teamId, groupId string
		load_grp_cluster                        *sql.Stmt
	)
	load_grp_cluster, err = tk.conn.Prepare(`
SELECT sc.cluster_id,
       sc.cluster_name,
       sc.organizational_team_id,
       sgmc.group_id
FROM   soma.repositories sr
JOIN   soma.buckets sb
ON     sr.repository_id = sb.repository_id
JOIN   soma.clusters sc
ON     sb.bucket_id = sc.bucket_id
JOIN   soma.group_membership_clusters sgmc
ON     sc.bucket_id = sgmc.bucket_id
AND    sc.cluster_id = sgmc.child_cluster_id
WHERE  sr.repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-grouped-clusters: ", err)
	}
	defer load_grp_cluster.Close()

	log.Printf("TK[%s]: loading grouped-clusters\n", tk.repoName)
	rows, err = load_grp_cluster.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading clusters: %s", tk.repoName, err.Error())
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
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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

func (tk *treeKeeper) startupClusters() {
	if tk.broken {
		return
	}

	var (
		err                                      error
		rows                                     *sql.Rows
		clusterId, clusterName, bucketId, teamId string
		load_cluster                             *sql.Stmt
	)
	load_cluster, err = tk.conn.Prepare(`
SELECT sc.cluster_id,
       sc.cluster_name,
	   sc.bucket_id,
	   sc.organizational_team_id
FROM   soma.repositories sr
JOIN   soma.buckets sb
ON     sr.repository_id = sb.repository_id
JOIN   soma.clusters sc
ON     sb.bucket_id = sc.bucket_id
WHERE  sr.repository_id = $1::uuid
AND    sc.object_state != 'grouped';`)
	if err != nil {
		log.Fatal("treekeeper/load-clusters: ", err)
	}
	defer load_cluster.Close()

	log.Printf("TK[%s]: loading clusters\n", tk.repoName)
	rows, err = load_cluster.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading clusters: %s", tk.repoName, err.Error())
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
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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

func (tk *treeKeeper) startupNodes() {
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
		load_nodes                                   *sql.Stmt
	)
	load_nodes, err = tk.conn.Prepare(`
SELECT    sn.node_id,
          sn.node_asset_id,
		  sn.node_name,
		  sn.organizational_team_id,
		  sn.server_id,
		  sn.node_online,
		  sn.node_deleted,
		  snba.bucket_id,
		  scm.cluster_id,
		  sgmn.group_id
FROM      soma.repositories sr
JOIN      soma.buckets sb
ON        sr.repository_id = sb.repository_id
JOIN      soma.node_bucket_assignment snba
ON        sb.bucket_id = snba.bucket_id
JOIN      soma.nodes sn
ON        snba.node_id = sn.node_id
LEFT JOIN soma.cluster_membership scm
ON        sn.node_id = scm.node_id
LEFT JOIN soma.group_membership_nodes sgmn
ON        sn.node_id = sgmn.child_node_id
WHERE     sr.repository_id = $1::uuid`)
	if err != nil {
		log.Fatal("treekeeper/load-nodes: ", err)
	}
	defer load_nodes.Close()

	log.Printf("TK[%s]: loading nodes\n", tk.repoName)
	rows, err = load_nodes.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading nodes: %s", tk.repoName, err.Error())
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
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
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

func (tk *treeKeeper) startupJobs() {
	if tk.broken {
		return
	}

	var (
		err       error
		rows      *sql.Rows
		job       string
		load_jobs *sql.Stmt
	)
	load_jobs, err = tk.conn.Prepare(`
SELECT   job
FROM     soma.jobs
WHERE    repository_id = $1::uuid
AND      job_status != 'processed'
ORDER BY job_serial ASC;`)
	if err != nil {
		log.Fatal("treekeeper/load-jobs: ", err)
	}
	defer load_jobs.Close()

	log.Printf("TK[%s]: loading pending jobs\n", tk.repoName)
	rows, err = load_jobs.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading clusters: %s", tk.repoName, err.Error())
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
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		tr := treeRequest{}
		err = json.Unmarshal([]byte(job), &tr)
		if err != nil {
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		tk.input <- tr
		log.Printf("TK[%s] Loaded job %s (%s)\n", tk.repoName, tr.JobId, tr.Action)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
