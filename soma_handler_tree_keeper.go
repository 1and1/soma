package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/satori/go.uuid"

)

type treeRequest struct {
	RequestType string
	Action      string
	JobId       uuid.UUID
	reply       chan somaResult
	Repository  somaRepositoryRequest
	Bucket      somaBucketRequest
	Group       somaGroupRequest
	Cluster     somaClusterRequest
	Node        somaNodeRequest
}

type treeResult struct {
	ResultType  string
	ResultError error
	JobId       uuid.UUID
	Repository  somaRepositoryResult
	Bucket      somaRepositoryRequest
}

type treeKeeper struct {
	repoId           string
	repoName         string
	input            chan treeRequest
	shutdown         chan bool
	conn             *sql.DB
	tree             *somatree.SomaTree
	errChan          chan *somatree.Error
	actionChan       chan *somatree.Action
	load_bucket      *sql.Stmt
	load_group       *sql.Stmt
	load_grp_mbr_grp *sql.Stmt
	load_grp_cluster *sql.Stmt
}

func (tk *treeKeeper) run() {
	log.Printf("Starting TreeKeeper for Repo %s (%s)", tk.repoName, tk.repoId)
	var err error

	log.Println("Prepare: treekeeper/load-buckets")
	tk.load_bucket, err = tk.conn.Prepare(`
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
	defer tk.load_bucket.Close()

	log.Println("Prepare: treekeeper/load-groups")
	tk.load_group, err = tk.conn.Prepare(`
SELECT sg.group_id,
       sg.group_name,
       sg.bucket_id,
       sg.organizational_team_id
FROM   soma.repositories sr
JOIN   soma.buckets sb
ON     sr.repository_id = sb.repository_id
JOIN   soma.groups sg
ON     sg.bucket_id = sg.bucket_id
WHERE  sr.repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-groups: ", err)
	}
	defer tk.load_group.Close()

	log.Println("Prepare: treekeeper/load-group-member-groups")
	tk.load_grp_mbr_grp, err = tk.conn.Prepare(`
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
	defer tk.load_grp_mbr_grp.Close()

	log.Println("Prepare: treekeeper/load-grouped-clusters")
	tk.load_grp_cluster, err = tk.conn.Prepare(`
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
	defer tk.load_grp_cluster.Close()

	tk.startupLoad()

runloop:
	for {
		select {
		case <-tk.shutdown:
			break runloop
		case <-tk.input:
			fmt.Printf("TK %s received input request", tk.repoName)
		}
	}
}

func (tk *treeKeeper) startupLoad() {
	tk.startupBuckets()
	tk.startupGroups()
	tk.startupGroupMemberGroups()
	tk.startupGroupedClusters()
}

func (tk *treeKeeper) startupBuckets() {
	var (
		rows                                      *sql.Rows
		bucketId, bucketName, environment, teamId string
		frozen, deleted                           bool
		err                                       error
	)
	log.Printf("TK[%s]: loading buckets\n", tk.repoName)
	rows, err = tk.load_bucket.Query(tk.repoId)
	if err != nil {
		log.Fatal(fmt.Errorf("TK[%s] Error loading buckets: %s", tk.repoName, err.Error()))
	}
	defer rows.Close()

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
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
		}
		somatree.NewBucket(somatree.BucketSpec{
			Id:          bucketId,
			Name:        bucketName,
			Environment: environment,
			Team:        teamId,
			Deleted:     deleted,
			Frozen:      frozen,
		}).Attach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "repository",
			ParentId:   tk.repoId,
			ParentName: tk.repoName,
		})
		for i := 0; i < len(tk.actionChan); i++ {
			<-tk.actionChan
		}
		for i := 0; i < len(tk.errChan); i++ {
			<-tk.errChan
		}
	}
}

func (tk *treeKeeper) startupGroups() {
	var (
		rows                                 *sql.Rows
		groupId, groupName, bucketId, teamId string
		err                                  error
	)
	log.Printf("TK[%s]: loading groups\n", tk.repoName)
	rows, err = tk.load_group.Query(tk.repoId)
	if err != nil {
		log.Fatal(fmt.Errorf("TK[%s] Error loading groups: %s", tk.repoName, err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&groupId,
			&groupName,
			&bucketId,
			&teamId,
		)
		if err != nil {
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
		}
		somatree.NewGroup(somatree.GroupSpec{
			Id:   groupId,
			Name: groupName,
			Team: teamId,
		}).Attach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "bucket",
			ParentId:   bucketId,
		})
		for i := 0; i < len(tk.actionChan); i++ {
			<-tk.actionChan
		}
		for i := 0; i < len(tk.errChan); i++ {
			<-tk.errChan
		}
	}
}

func (tk *treeKeeper) startupGroupMemberGroups() {
	var (
		rows                  *sql.Rows
		groupId, childGroupId string
		err                   error
	)
	log.Printf("TK[%s]: loading group-member-groups\n", tk.repoName)
	rows, err = tk.load_grp_mbr_grp.Query(tk.repoId)
	if err != nil {
		log.Fatal(fmt.Errorf("TK[%s] Error loading groups: %s", tk.repoName, err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&groupId,
			&childGroupId,
		)
		if err != nil {
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
		}

		tk.tree.Find(somatree.FindRequest{
			ElementType: "group",
			ElementId:   childGroupId,
		}, true).(somatree.SomaTreeBucketAttacher).ReAttach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentId:   groupId,
		})
	}
	for i := 0; i < len(tk.actionChan); i++ {
		<-tk.actionChan
	}
	for i := 0; i < len(tk.errChan); i++ {
		<-tk.errChan
	}
}

func (tk *treeKeeper) startupGroupedClusters() {
	var (
		err                                     error
		rows                                    *sql.Rows
		clusterId, clusterName, teamId, groupId string
	)
	log.Printf("TK[%s]: loading grouped-clusters\n", tk.repoName)
	rows, err = tk.load_grp_cluster.Query(tk.repoId)
	if err != nil {
		log.Fatal(fmt.Errorf("TK[%s] Error loading clusters: %s", tk.repoName, err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&clusterId,
			&clusterName,
			&teamId,
			&groupId,
		)
		if err != nil {
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
		}

		somatree.NewCluster(somatree.ClusterSpec{
			Id:   clusterId,
			Name: clusterName,
			Team: teamId,
		}).Attach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentId:   groupId,
		})
	}
	for i := 0; i < len(tk.actionChan); i++ {
		<-tk.actionChan
	}
	for i := 0; i < len(tk.errChan); i++ {
		<-tk.errChan
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
