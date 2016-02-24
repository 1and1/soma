package main

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/satori/go.uuid"

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

	tk.startupGroupSystemProperties()
	tk.startupGroupCustomProperties()
	tk.startupGroupOncallProperties()
	tk.startupGroupServiceProperties()

	if len(tk.actionChan) > 0 {
		log.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	tk.startupJobs()

	if len(tk.actionChan) > 0 {
		log.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.repoName)
		tk.broken = true
		return
	}

	b, _ := json.Marshal(tk.tree)
	log.Println(string(b))
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
	log.Println("Prepare: treekeeper/load-buckets")
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
		somatree.NewBucket(somatree.BucketSpec{
			Id:          bucketId,
			Name:        bucketName,
			Environment: environment,
			Team:        teamId,
			Deleted:     deleted,
			Frozen:      frozen,
			Repository:  tk.repoId,
		}).Attach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "repository",
			ParentId:   tk.repoId,
			ParentName: tk.repoName,
		})
		for i := len(tk.actionChan); i > 0; i-- {
			a := <-tk.actionChan
			log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := 0; i < len(tk.errChan); i++ {
			<-tk.errChan
		}
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
	log.Println("Prepare: treekeeper/load-groups")
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
		somatree.NewGroup(somatree.GroupSpec{
			Id:   groupId,
			Name: groupName,
			Team: teamId,
		}).Attach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "bucket",
			ParentId:   bucketId,
		})
		for i := len(tk.actionChan); i > 0; i-- {
			a := <-tk.actionChan
			log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := 0; i < len(tk.errChan); i++ {
			<-tk.errChan
		}
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
	log.Println("Prepare: treekeeper/load-group-member-groups")
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

		tk.tree.Find(somatree.FindRequest{
			ElementType: "group",
			ElementId:   childGroupId,
		}, true).(somatree.SomaTreeBucketAttacher).ReAttach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentId:   groupId,
		})
	}
	for i := len(tk.actionChan); i > 0; i-- {
		a := <-tk.actionChan
		log.Printf("%s -> %s\n", a.Action, a.Type)
	}
	for i := 0; i < len(tk.errChan); i++ {
		<-tk.errChan
	}
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
	log.Println("Prepare: treekeeper/load-grouped-clusters")
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
	for i := len(tk.actionChan); i > 0; i-- {
		a := <-tk.actionChan
		log.Printf("%s -> %s\n", a.Action, a.Type)
	}
	for i := 0; i < len(tk.errChan); i++ {
		<-tk.errChan
	}
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
	log.Println("Prepare: treekeeper/load-clusters")
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

		somatree.NewCluster(somatree.ClusterSpec{
			Id:   clusterId,
			Name: clusterName,
			Team: teamId,
		}).Attach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "bucket",
			ParentId:   bucketId,
		})
	}
	for i := len(tk.actionChan); i > 0; i-- {
		a := <-tk.actionChan
		log.Printf("%s -> %s [%s(%s)]\n", a.Action, a.Type, a.Cluster.Id, a.Cluster.Name)
	}
	for i := 0; i < len(tk.errChan); i++ {
		<-tk.errChan
	}
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
	log.Println("Prepare: treekeeper/load-nodes")
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

		node := somatree.NewNode(somatree.NodeSpec{
			Id:       nodeId,
			AssetId:  uint64(assetId),
			Name:     nodeName,
			Team:     teamId,
			ServerId: serverId,
			Online:   nodeOnline,
			Deleted:  nodeDeleted,
		})
		if clusterId.Valid {
			node.Attach(somatree.AttachRequest{
				Root:       tk.tree,
				ParentType: "cluster",
				ParentId:   clusterId.String,
			})
		} else if groupId.Valid {
			node.Attach(somatree.AttachRequest{
				Root:       tk.tree,
				ParentType: "group",
				ParentId:   groupId.String,
			})
		} else {
			node.Attach(somatree.AttachRequest{
				Root:       tk.tree,
				ParentType: "bucket",
				ParentId:   bucketId,
			})
		}
	}
	for i := len(tk.actionChan); i > 0; i-- {
		a := <-tk.actionChan
		log.Printf("%s -> %s\n", a.Action, a.Type)
	}
	for i := 0; i < len(tk.errChan); i++ {
		<-tk.errChan
	}
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
	log.Println("Prepare: treekeeper/load-jobs")
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

func (tk *treeKeeper) startupGroupSystemProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                         error
		instanceId, srcInstanceId, groupId, view, systemProperty, objectType, value string
		inInstanceId, inObjectType, inObjId                                         string
		inheritance, childrenOnly                                                   bool
		rows, instance_rows                                                         *sql.Rows
		load_properties, load_instances                                             *sql.Stmt
	)
	log.Println("Prepare: treekeeper/load-group-system-properties")
	load_properties, err = tk.conn.Prepare(`
SELECT instance_id,
       source_instance_id,
	   group_id,
	   view,
	   system_property,
	   object_type,
	   inheritance_enabled,
	   children_only,
	   value
FROM   soma.group_system_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-group-system-properties: ", err)
	}
	defer load_properties.Close()

	log.Println("Prepare: treekeeper/load-group-system-property-instances")
	load_instances, err = tk.conn.Prepare(tkStmtLoadSystemPropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-group-system-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading group system properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading group system properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

systemloop:
	// load all system properties defined directly on group objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&groupId,
			&view,
			&systemProperty,
			&objectType,
			&inheritance,
			&childrenOnly,
			&value,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break systemloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := somatree.PropertySystem{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          systemProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Instances = make([]somatree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading group system properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current group system property so the IDs can be set correctly
		for instance_rows.Next() {
			err = instance_rows.Scan(
				&inInstanceId,
				&inObjectType,
				&inObjId,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break inproploop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			var propObjectId, propInstanceId uuid.UUID
			if propObjectId, err = uuid.FromString(inObjId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if uuid.Equal(uuid.Nil, propObjectId) || uuid.Equal(uuid.Nil, propInstanceId) {
				continue inproploop
			}
			if inObjectType == "MAGIC_NO_RESULT_VALUE" {
				continue inproploop
			}

			pi := somatree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the group and set the prepared property
		tk.tree.Find(somatree.FindRequest{
			ElementId: groupId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		for i := len(tk.actionChan); i > 0; i-- {
			a := <-tk.actionChan
			log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := 0; i < len(tk.errChan); i++ {
			<-tk.errChan
		}
	}
}

func (tk *treeKeeper) startupGroupCustomProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                       error
		instanceId, srcInstanceId, groupId, view, customId, customProperty, value string
		inInstanceId, inObjectType, inObjId                                       string
		inheritance, childrenOnly                                                 bool
		rows, instance_rows                                                       *sql.Rows
		load_properties, load_instances                                           *sql.Stmt
	)
	log.Println("Prepare: treekeeper/load-group-custom-properties")
	load_properties, err = tk.conn.Prepare(`
SELECT sgcp.instance_id,
       sgcp.source_instance_id,
	   sgcp.group_id,
	   sgcp.view,
	   sgcp.custom_property_id,
	   sgcp.inheritance_enabled,
	   sgcp.children_only,
	   sgcp.value,
	   scp.custom_property
FROM   soma.group_custom_properties sgcp
JOIN   soma.custom_properties scp
ON     sgcp.custom_property_id = scp.custom_property_id
WHERE  sgcp.instance_id = sgcp.source_instance_id
AND    sgcp.repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-group-custom-properties: ", err)
	}
	defer load_properties.Close()

	log.Println("Prepare: treekeeper/load-group-custom-property-instances")
	load_instances, err = tk.conn.Prepare(tkStmtLoadCustomPropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-group-custom-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading group custom properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading group custom properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

customloop:
	// load all system properties defined directly on group objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&groupId,
			&view,
			&customId,
			&inheritance,
			&childrenOnly,
			&value,
			&customProperty,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break customloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := somatree.PropertyCustom{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Key:          customProperty,
			Value:        value,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.CustomId, _ = uuid.FromString(customId)
		prop.Instances = make([]somatree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading group system properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current group system property so the IDs can be set correctly
		for instance_rows.Next() {
			err = instance_rows.Scan(
				&inInstanceId,
				&inObjectType,
				&inObjId,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break inproploop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			var propObjectId, propInstanceId uuid.UUID
			if propObjectId, err = uuid.FromString(inObjId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if uuid.Equal(uuid.Nil, propObjectId) || uuid.Equal(uuid.Nil, propInstanceId) {
				continue inproploop
			}
			if inObjectType == "MAGIC_NO_RESULT_VALUE" {
				continue inproploop
			}

			pi := somatree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the group and set the prepared property
		tk.tree.Find(somatree.FindRequest{
			ElementId: groupId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		for i := len(tk.actionChan); i > 0; i-- {
			a := <-tk.actionChan
			log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := 0; i < len(tk.errChan); i++ {
			<-tk.errChan
		}
	}
}

func (tk *treeKeeper) startupGroupServiceProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                               error
		instanceId, srcInstanceId, groupId, view, serviceProperty, teamId string
		inInstanceId, inObjectType, inObjId, attrKey, attrValue           string
		inheritance, childrenOnly                                         bool
		rows, attribute_rows, instance_rows                               *sql.Rows
		load_properties, load_attributes, load_instances                  *sql.Stmt
	)
	log.Println("Prepare: treekeeper/load-group-service-properties")
	load_properties, err = tk.conn.Prepare(`
SELECT instance_id,
       source_instance_id,
	   group_id,
	   view,
	   service_property,
	   organizational_team_id,
	   inheritance_enabled,
	   children_only
FROM   soma.group_service_properties
WHERE  instance_id = source_instance_id
AND    repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-group-service-properties: ", err)
	}
	defer load_properties.Close()

	log.Println("Prepare: treekeeper/load-group-service-property-attributes")
	load_attributes, err = tk.conn.Prepare(`
SELECT service_property_attribute,
       value
FROM   soma.team_service_property_values
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`)
	if err != nil {
		log.Fatal("treekeeper/load-group-service-property-attributes: ", err)
	}
	defer load_attributes.Close()

	log.Println("Prepare: treekeeper/load-group-service-property-instances")
	load_instances, err = tk.conn.Prepare(tkStmtLoadServicePropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-group-service-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading group service properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading group custom properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

serviceloop:
	// load all system properties defined directly on group objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&groupId,
			&view,
			&serviceProperty,
			&teamId,
			&inheritance,
			&childrenOnly,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break serviceloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := somatree.PropertyService{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Service:      serviceProperty,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.Attributes = make([]somatree.PropertyServiceAttribute, 0)
		prop.Instances = make([]somatree.PropertyInstance, 0)

		attribute_rows, err = load_attributes.Query(
			teamId,
			serviceProperty,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading group service properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer attribute_rows.Close()

	attributeloop:
		// load service attributes
		for attribute_rows.Next() {
			err = attribute_rows.Scan(
				&attrKey,
				&attrValue,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break attributeloop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			pa := somatree.PropertyServiceAttribute{
				Attribute: attrKey,
				Value:     attrValue,
			}
			prop.Attributes = append(prop.Attributes, pa)
		}

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading group service properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current group system property so the IDs can be set correctly
		for instance_rows.Next() {
			err = instance_rows.Scan(
				&inInstanceId,
				&inObjectType,
				&inObjId,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break inproploop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			var propObjectId, propInstanceId uuid.UUID
			if propObjectId, err = uuid.FromString(inObjId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if uuid.Equal(uuid.Nil, propObjectId) || uuid.Equal(uuid.Nil, propInstanceId) {
				continue inproploop
			}
			if inObjectType == "MAGIC_NO_RESULT_VALUE" {
				continue inproploop
			}

			pi := somatree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the group and set the prepared property
		tk.tree.Find(somatree.FindRequest{
			ElementId: groupId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		for i := len(tk.actionChan); i > 0; i-- {
			a := <-tk.actionChan
			log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := 0; i < len(tk.errChan); i++ {
			<-tk.errChan
		}
	}
}

func (tk *treeKeeper) startupGroupOncallProperties() {
	if tk.broken {
		return
	}

	var (
		err                                                                          error
		instanceId, srcInstanceId, groupId, view, oncallId, oncallName, oncallNumber string
		inInstanceId, inObjectType, inObjId                                          string
		inheritance, childrenOnly                                                    bool
		rows, instance_rows                                                          *sql.Rows
		load_properties, load_instances                                              *sql.Stmt
	)
	log.Println("Prepare: treekeeper/load-group-oncall-properties")
	load_properties, err = tk.conn.Prepare(`
SELECT  sgop.instance_id,
        sgop.source_instance_id,
        sgop.group_id,
        sgop.view,
        sgop.oncall_duty_id,
        sgop.inheritance_enabled,
        sgop.children_only,
        iodt.oncall_duty_name,
        iodt.oncall_duty_phone_number
FROM    soma.group_oncall_properties sgop
JOIN    inventory.oncall_duty_teams iodt
  ON    sgop.oncall_duty_id = iodt.oncall_duty_id
WHERE   sgop.instance_id = sgop.source_instance_id
  AND   sgop.repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/load-group-oncall-properties: ", err)
	}
	defer load_properties.Close()

	log.Println("Prepare: treekeeper/load-group-oncall-property-instances")
	load_instances, err = tk.conn.Prepare(tkStmtLoadOncallPropInstances)
	if err != nil {
		log.Fatal("treekeeper/load-group-oncall-property-instances: ", err)
	}
	defer load_instances.Close()

	log.Printf("TK[%s]: loading group oncall properties\n", tk.repoName)
	rows, err = load_properties.Query(tk.repoId)
	if err != nil {
		log.Printf("TK[%s] Error loading group oncall properties: %s", tk.repoName, err.Error())
		tk.broken = true
		return
	}
	defer rows.Close()

oncallloop:
	// load all system properties defined directly on group objects
	for rows.Next() {
		err = rows.Scan(
			&instanceId,
			&srcInstanceId,
			&groupId,
			&view,
			&oncallId,
			&inheritance,
			&childrenOnly,
			&oncallName,
			&oncallNumber,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break oncallloop
			}
			log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
			tk.broken = true
			return
		}

		// build the property
		prop := somatree.PropertyOncall{
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			View:         view,
			Name:         oncallName,
			Number:       oncallNumber,
		}
		prop.Id, _ = uuid.FromString(instanceId)
		prop.OncallId, _ = uuid.FromString(oncallId)
		prop.Instances = make([]somatree.PropertyInstance, 0)

		instance_rows, err = load_instances.Query(
			tk.repoId,
			srcInstanceId,
		)
		if err != nil {
			log.Printf("TK[%s] Error loading group custom properties: %s", tk.repoName, err.Error())
			tk.broken = true
			return
		}
		defer instance_rows.Close()

	inproploop:
		// load all all ids for properties that were inherited from the
		// current group system property so the IDs can be set correctly
		for instance_rows.Next() {
			err = instance_rows.Scan(
				&inInstanceId,
				&inObjectType,
				&inObjId,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break inproploop
				}
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}

			var propObjectId, propInstanceId uuid.UUID
			if propObjectId, err = uuid.FromString(inObjId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if propInstanceId, err = uuid.FromString(inInstanceId); err != nil {
				log.Printf("TK[%s] Error: %s\n", tk.repoName, err.Error())
				tk.broken = true
				return
			}
			if uuid.Equal(uuid.Nil, propObjectId) || uuid.Equal(uuid.Nil, propInstanceId) {
				continue inproploop
			}
			if inObjectType == "MAGIC_NO_RESULT_VALUE" {
				continue inproploop
			}

			pi := somatree.PropertyInstance{
				ObjectId:   propObjectId,
				ObjectType: inObjectType,
				InstanceId: propInstanceId,
			}
			prop.Instances = append(prop.Instances, pi)
		}

		// lookup the group and set the prepared property
		tk.tree.Find(somatree.FindRequest{
			ElementId: groupId,
		}, true).SetProperty(&prop)

		// throw away all generated actions, we do this for every
		// property since with inheritance this can create a lot of
		// actions
		for i := len(tk.actionChan); i > 0; i-- {
			a := <-tk.actionChan
			log.Printf("%s -> %s\n", a.Action, a.Type)
		}
		for i := 0; i < len(tk.errChan); i++ {
			<-tk.errChan
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
