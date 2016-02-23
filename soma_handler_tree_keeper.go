package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

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
	repoId            string
	repoName          string
	team              string
	broken            bool
	ready             bool
	input             chan treeRequest
	shutdown          chan bool
	conn              *sql.DB
	tree              *somatree.SomaTree
	errChan           chan *somatree.Error
	actionChan        chan *somatree.Action
	start_job         *sql.Stmt
	create_bucket     *sql.Stmt
	defer_constraints *sql.Stmt
}

func (tk *treeKeeper) run() {
	log.Printf("Starting TreeKeeper for Repo %s (%s)", tk.repoName, tk.repoId)
	tk.startupLoad()

	if tk.broken {
		tickTack := time.NewTicker(time.Second * 10).C
	hoverloop:
		for {
			select {
			case <-tickTack:
				log.Printf("TK[%s]: BROKEN REPOSITORY %s flying holding patterns!\n",
					tk.repoName, tk.repoId)
			case <-tk.shutdown:
				break hoverloop
			}
		}
		return
	}
	log.Printf("TK[%s]: ready for service!\n", tk.repoName)
	tk.ready = true

	var err error
	log.Println("Prepare: treekeeper/start-job")
	tk.start_job, err = tk.conn.Prepare(tkStmtStartJob)
	if err != nil {
		log.Fatal("treekeeper/start-job: ", err)
	}
	defer tk.start_job.Close()

runloop:
	for {
		select {
		case <-tk.shutdown:
			break runloop
		case req := <-tk.input:
			tk.process(&req)
		}
	}
}

func (tk *treeKeeper) isReady() bool {
	return tk.ready
}

func (tk *treeKeeper) isBroken() bool {
	return tk.broken
}

func (tk *treeKeeper) process(q *treeRequest) {
	var (
		err                          error
		tx                           *sql.Tx
		txStmtPropertyInstanceCreate *sql.Stmt
		txStmtCreateBucket           *sql.Stmt

		txStmtGroupCreate                *sql.Stmt
		txStmtGroupUpdate                *sql.Stmt
		txStmtGroupDelete                *sql.Stmt
		txStmtGroupMemberNewNode         *sql.Stmt
		txStmtGroupMemberNewCluster      *sql.Stmt
		txStmtGroupMemberNewGroup        *sql.Stmt
		txStmtGroupMemberRemoveNode      *sql.Stmt
		txStmtGroupMemberRemoveCluster   *sql.Stmt
		txStmtGroupMemberRemoveGroup     *sql.Stmt
		txStmtGroupPropertyServiceCreate *sql.Stmt
		txStmtGroupPropertySystemCreate  *sql.Stmt
		txStmtGroupPropertyOncallCreate  *sql.Stmt
		txStmtGroupPropertyCustomCreate  *sql.Stmt

		txStmtClusterCreate       *sql.Stmt
		txStmtClusterUpdate       *sql.Stmt
		txStmtClusterDelete       *sql.Stmt
		txStmtClusterMemberNew    *sql.Stmt
		txStmtClusterMemberRemove *sql.Stmt

		txStmtBucketAssignNode       *sql.Stmt
		txStmtUpdateNodeState        *sql.Stmt
		txStmtNodeUnassignFromBucket *sql.Stmt
	)
	_, err = tk.start_job.Exec(q.JobId.String(), time.Now().UTC())
	if err != nil {
		log.Println(err)
	}
	log.Printf("Processing job: %s\n", q.JobId.String())

	tk.tree.Begin()

	switch q.Action {
	case "create_bucket":
		somatree.NewBucket(somatree.BucketSpec{
			Id:          uuid.NewV4().String(),
			Name:        q.Bucket.Bucket.Name,
			Environment: q.Bucket.Bucket.Environment,
			Team:        tk.team,
			Deleted:     q.Bucket.Bucket.IsDeleted,
			Frozen:      q.Bucket.Bucket.IsFrozen,
			Repository:  q.Bucket.Bucket.Repository,
		}).Attach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "repository",
			ParentId:   tk.repoId,
			ParentName: tk.repoName,
		})

	// GROUP MANIPULATION REQUESTS
	case "create_group":
		somatree.NewGroup(somatree.GroupSpec{
			Id:   uuid.NewV4().String(),
			Name: q.Group.Group.Name,
			Team: q.Group.Group.TeamId,
		}).Attach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "bucket",
			ParentId:   q.Group.Group.BucketId,
		})
	case "delete_group":
		tk.tree.Find(somatree.FindRequest{
			ElementType: "group",
			ElementId:   q.Group.Group.Id,
		}, true).(somatree.SomaTreeBucketAttacher).Destroy()
	case "reset_group_to_bucket":
		tk.tree.Find(somatree.FindRequest{
			ElementType: "group",
			ElementId:   q.Group.Group.Id,
		}, true).(somatree.SomaTreeBucketAttacher).Detach()
	case "add_group_to_group":
		tk.tree.Find(somatree.FindRequest{
			ElementType: "group",
			ElementId:   q.Group.Group.MemberGroups[0].Id,
		}, true).(somatree.SomaTreeBucketAttacher).ReAttach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentId:   q.Group.Group.Id,
		})
	case "add_system_property_to_group": // XXX MOCKUP DATA
		tk.tree.Find(somatree.FindRequest{
			ElementType: "group",
			ElementId:   q.Group.Group.Id,
		}, true).(somatree.SomaTreePropertier).SetProperty(&somatree.PropertySystem{
			Id:           uuid.NewV4(),
			Inheritance:  true,
			ChildrenOnly: false,
			View:         "internal",
			Key:          "dns_zone",
			Value:        "mw.server.lan",
		})
	case "add_service_property_to_group": // XXX MOCKUP DATA
		tk.tree.Find(somatree.FindRequest{
			ElementType: "group",
			ElementId:   q.Group.Group.Id,
		}, true).(somatree.SomaTreePropertier).SetProperty(&somatree.PropertyService{
			Id:           uuid.NewV4(),
			Inheritance:  true,
			ChildrenOnly: false,
			View:         "internal",
			Service:      "pmrmw",
			Attributes: []somatree.PropertyServiceAttribute{
				{
					Attribute: "proto_transport",
					Value:     "tcp",
				},
				{
					Attribute: "port",
					Value:     "9192",
				},
				{
					Attribute: "port",
					Value:     "9193",
				},
			},
		})
	case "add_oncall_property_to_group":
	case "add_custom_property_to_group":

	// CLUSTER MANIPULATION REQUESTS
	case "create_cluster":
		somatree.NewCluster(somatree.ClusterSpec{
			Id:   uuid.NewV4().String(),
			Name: q.Cluster.Cluster.Name,
			Team: q.Cluster.Cluster.TeamId,
		}).Attach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "bucket",
			ParentId:   q.Cluster.Cluster.BucketId,
		})
	case "delete_cluster":
		tk.tree.Find(somatree.FindRequest{
			ElementType: "cluster",
			ElementId:   q.Cluster.Cluster.Id,
		}, true).(somatree.SomaTreeBucketAttacher).Destroy()
	case "reset_cluster_to_bucket":
		tk.tree.Find(somatree.FindRequest{
			ElementType: "cluster",
			ElementId:   q.Cluster.Cluster.Id,
		}, true).(somatree.SomaTreeBucketAttacher).Detach()
	case "add_cluster_to_group":
		tk.tree.Find(somatree.FindRequest{
			ElementType: "cluster",
			ElementId:   q.Group.Group.MemberClusters[0].Id,
		}, true).(somatree.SomaTreeBucketAttacher).ReAttach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentId:   q.Group.Group.Id,
		})

	// NODE MANIPULATION REQUESTS
	case "create_node":
		somatree.NewNode(somatree.NodeSpec{
			Id:       q.Node.Node.Id,
			AssetId:  q.Node.Node.AssetId,
			Name:     q.Node.Node.Name,
			Team:     q.Node.Node.Team,
			ServerId: q.Node.Node.Server,
			Online:   q.Node.Node.IsOnline,
			Deleted:  q.Node.Node.IsDeleted,
		}).Attach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "bucket",
			ParentId:   q.Node.Node.Config.BucketId,
		})
	case "delete_node":
		tk.tree.Find(somatree.FindRequest{
			ElementType: "node",
			ElementId:   q.Node.Node.Id,
		}, true).(somatree.SomaTreeBucketAttacher).Destroy()
	case "reset_node_to_bucket":
		tk.tree.Find(somatree.FindRequest{
			ElementType: "node",
			ElementId:   q.Node.Node.Id,
		}, true).(somatree.SomaTreeBucketAttacher).Detach()
	case "add_node_to_group":
		tk.tree.Find(somatree.FindRequest{
			ElementType: "node",
			ElementId:   q.Group.Group.MemberNodes[0].Id,
		}, true).(somatree.SomaTreeBucketAttacher).ReAttach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentId:   q.Group.Group.Id,
		})
	case "add_node_to_cluster":
		tk.tree.Find(somatree.FindRequest{
			ElementType: "node",
			ElementId:   q.Cluster.Cluster.Members[0].Id,
		}, true).(somatree.SomaTreeBucketAttacher).ReAttach(somatree.AttachRequest{
			Root:       tk.tree,
			ParentType: "cluster",
			ParentId:   q.Cluster.Cluster.Id,
		})
	}

	// open multi-statement transaction
	if tx, err = tk.conn.Begin(); err != nil {
		goto bailout
	}
	defer tx.Rollback()

	// prepare statements within tx context
	if txStmtCreateBucket, err = tx.Prepare(tkStmtCreateBucket); err != nil {
		goto bailout
	}
	defer txStmtCreateBucket.Close()

	if txStmtGroupCreate, err = tx.Prepare(tkStmtGroupCreate); err != nil {
		goto bailout
	}
	defer txStmtGroupCreate.Close()

	if txStmtGroupUpdate, err = tx.Prepare(tkStmtGroupUpdate); err != nil {
		goto bailout
	}
	defer txStmtGroupUpdate.Close()

	if txStmtGroupDelete, err = tx.Prepare(tkStmtGroupDelete); err != nil {
		goto bailout
	}
	defer txStmtGroupDelete.Close()

	if txStmtGroupMemberNewNode, err = tx.Prepare(tkStmtGroupMemberNewNode); err != nil {
		goto bailout
	}
	defer txStmtGroupMemberNewNode.Close()

	if txStmtGroupMemberNewCluster, err = tx.Prepare(tkStmtGroupMemberNewCluster); err != nil {
		goto bailout
	}
	defer txStmtGroupMemberNewCluster.Close()

	if txStmtGroupMemberNewGroup, err = tx.Prepare(tkStmtGroupMemberNewGroup); err != nil {
		goto bailout
	}
	defer txStmtGroupMemberNewGroup.Close()

	if txStmtGroupMemberRemoveNode, err = tx.Prepare(tkStmtGroupMemberRemoveNode); err != nil {
		goto bailout
	}
	defer txStmtGroupMemberRemoveNode.Close()

	if txStmtGroupMemberRemoveCluster, err = tx.Prepare(tkStmtGroupMemberRemoveCluster); err != nil {
		goto bailout
	}
	defer txStmtGroupMemberRemoveCluster.Close()

	if txStmtGroupMemberRemoveGroup, err = tx.Prepare(tkStmtGroupMemberRemoveGroup); err != nil {
		goto bailout
	}
	defer txStmtGroupMemberRemoveGroup.Close()

	if txStmtGroupPropertyOncallCreate, err = tx.Prepare(tkStmtGroupPropertyOncallCreate); err != nil {
		goto bailout
	}
	defer txStmtGroupPropertyOncallCreate.Close()

	if txStmtGroupPropertyServiceCreate, err = tx.Prepare(tkStmtGroupPropertyServiceCreate); err != nil {
		goto bailout
	}
	defer txStmtGroupPropertyServiceCreate.Close()

	if txStmtGroupPropertySystemCreate, err = tx.Prepare(tkStmtGroupPropertySystemCreate); err != nil {
		goto bailout
	}
	defer txStmtGroupPropertySystemCreate.Close()

	if txStmtGroupPropertyCustomCreate, err = tx.Prepare(tkStmtGroupPropertyCustomCreate); err != nil {
		goto bailout
	}
	defer txStmtGroupPropertyCustomCreate.Close()

	// CLUSTER
	if txStmtClusterCreate, err = tx.Prepare(tkStmtClusterCreate); err != nil {
		goto bailout
	}
	defer txStmtClusterCreate.Close()

	if txStmtClusterUpdate, err = tx.Prepare(tkStmtClusterUpdate); err != nil {
		goto bailout
	}
	defer txStmtClusterUpdate.Close()

	if txStmtClusterDelete, err = tx.Prepare(tkStmtClusterDelete); err != nil {
		goto bailout
	}
	defer txStmtClusterDelete.Close()

	if txStmtClusterMemberNew, err = tx.Prepare(tkStmtClusterMemberNew); err != nil {
		goto bailout
	}
	defer txStmtClusterMemberNew.Close()

	if txStmtClusterMemberRemove, err = tx.Prepare(tkStmtClusterMemberRemove); err != nil {
		goto bailout
	}
	defer txStmtClusterMemberRemove.Close()

	// NODE?
	if txStmtBucketAssignNode, err = tx.Prepare(tkStmtBucketAssignNode); err != nil {
		goto bailout
	}
	defer txStmtBucketAssignNode.Close()

	if txStmtUpdateNodeState, err = tx.Prepare(tkStmtUpdateNodeState); err != nil {
		goto bailout
	}
	defer txStmtUpdateNodeState.Close()

	if txStmtNodeUnassignFromBucket, err = tx.Prepare(tkStmtNodeUnassignFromBucket); err != nil {
		goto bailout
	}
	defer txStmtNodeUnassignFromBucket.Close()

	// defer constraint checks
	if _, err = tx.Exec(tkStmtDeferAllConstraints); err != nil {
		goto bailout
	}

actionloop:
	for i := 0; i < len(tk.actionChan); i++ {
		a := <-tk.actionChan
		switch a.Type {
		// BUCKET
		case "bucket":
			switch a.Action {
			case "create":
				if _, err = txStmtCreateBucket.Exec(
					a.Bucket.Id,
					a.Bucket.Name,
					a.Bucket.IsFrozen,
					a.Bucket.IsDeleted,
					a.Bucket.Repository,
					a.Bucket.Environment,
					a.Bucket.Team,
				); err != nil {
					break actionloop
				}
			case "node_assignment":
				if _, err = txStmtBucketAssignNode.Exec(
					a.Node.Id,
					a.Bucket.Id,
					a.Bucket.Team,
				); err != nil {
					break actionloop
				}
			default:
				jB, _ := json.Marshal(a)
				log.Printf("Unhandled message: %s\n", string(jB))
			}
		// GROUP
		case "group":
			switch a.Action {
			case "create":
				if _, err = txStmtGroupCreate.Exec(
					a.Group.Id,
					a.Group.BucketId,
					a.Group.Name,
					a.Group.ObjectState,
					a.Group.TeamId,
				); err != nil {
					break actionloop
				}
			case "update":
				if _, err = txStmtGroupUpdate.Exec(
					a.Group.Id,
					a.Group.ObjectState,
				); err != nil {
					break actionloop
				}
			case "delete":
				if _, err = txStmtGroupDelete.Exec(
					a.Group.Id,
				); err != nil {
					break actionloop
				}
			case "member_new":
				switch a.ChildType {
				case "group":
					if _, err = txStmtGroupMemberNewGroup.Exec(
						a.Group.Id,
						a.ChildGroup.Id,
						a.Group.BucketId,
					); err != nil {
						break actionloop
					}
				case "cluster":
					if _, err = txStmtGroupMemberNewCluster.Exec(
						a.Group.Id,
						a.ChildCluster.Id,
						a.Group.BucketId,
					); err != nil {
						break actionloop
					}
				case "node":
					if _, err = txStmtGroupMemberNewNode.Exec(
						a.Group.Id,
						a.ChildNode.Id,
						a.Group.BucketId,
					); err != nil {
						break actionloop
					}
				}
			case "member_removed":
				switch a.ChildType {
				case "group":
					if _, err = txStmtGroupMemberRemoveGroup.Exec(
						a.Group.Id,
						a.ChildGroup.Id,
					); err != nil {
						break actionloop
					}
				case "cluster":
					if _, err = txStmtGroupMemberRemoveCluster.Exec(
						a.Group.Id,
						a.ChildCluster.Id,
					); err != nil {
						break actionloop
					}
				case "node":
					if _, err = txStmtGroupMemberRemoveNode.Exec(
						a.Group.Id,
						a.ChildNode.Id,
					); err != nil {
						break actionloop
					}
				}
			case "property_new":
				if _, err = txStmtPropertyInstanceCreate.Exec(
					a.Property.InstanceId,
					a.Property.RepositoryId,
					a.Property.SourceInstanceId,
					a.Property.SourceType,
					a.Property.InheritedFrom,
				); err != nil {
					break actionloop
				}
				switch a.Property.PropertyType {
				case "custom":
					if _, err = txStmtGroupPropertyCustomCreate.Exec(
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Group.Id,
						a.Property.View,
						a.Property.Custom.CustomId,
						a.Property.BucketId,
						a.Property.Custom.RepositoryId,
						a.Property.Inheritance,
						a.Property.ChildrenOnly,
						a.Property.Custom.Value,
					); err != nil {
						break actionloop
					}
				case "system":
					if _, err = txStmtGroupPropertySystemCreate.Exec(
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
					); err != nil {
						break actionloop
					}
				case "service":
					if _, err = txStmtGroupPropertyServiceCreate.Exec(
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Group.Id,
						a.Property.View,
						a.Property.Service.Name,
						a.Property.Service.TeamId,
						a.Property.RepositoryId,
						a.Property.Inheritance,
						a.Property.ChildrenOnly,
					); err != nil {
						break actionloop
					}
				case "oncall":
					if _, err = txStmtGroupPropertyOncallCreate.Exec(
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Group.Id,
						a.Property.View,
						a.Property.Oncall.OncallId,
						a.Property.RepositoryId,
						a.Property.Inheritance,
						a.Property.ChildrenOnly,
					); err != nil {
						break actionloop
					}
				}
			}
		// CLUSTER
		case "cluster":
			switch a.Action {
			case "create":
				if _, err = txStmtClusterCreate.Exec(
					a.Cluster.Id,
					a.Cluster.Name,
					a.Cluster.BucketId,
					a.Cluster.ObjectState,
					a.Cluster.TeamId,
				); err != nil {
					break actionloop
				}
			case "update":
				if _, err = txStmtClusterUpdate.Exec(
					a.Cluster.Id,
					a.Cluster.ObjectState,
				); err != nil {
					break actionloop
				}
			case "delete":
				if _, err = txStmtClusterDelete.Exec(
					a.Cluster.Id,
				); err != nil {
					break actionloop
				}
			case "member_new":
				if _, err = txStmtClusterMemberNew.Exec(
					a.Cluster.Id,
					a.ChildNode.Id,
					a.Cluster.BucketId,
				); err != nil {
					break actionloop
				}
			case "member_removed":
				if _, err = txStmtClusterMemberRemove.Exec(
					a.Cluster.Id,
					a.ChildNode.Id,
				); err != nil {
					break actionloop
				}
			}
		// NODE
		case "node":
			switch a.Action {
			case "delete":
				if _, err = txStmtNodeUnassignFromBucket.Exec(
					a.Node.Id,
					a.Node.Config.BucketId,
					a.Node.Team,
				); err != nil {
					break actionloop
				}
				fallthrough // need to call txStmtUpdateNodeState for delete as well
			case "update":
				if _, err = txStmtUpdateNodeState.Exec(
					a.Node.Id,
					a.Node.State,
				); err != nil {
					break actionloop
				}
			}
		case "errorchannel":
			continue actionloop
		default:
			jB, _ := json.Marshal(a)
			log.Printf("Unhandled message: %s\n", string(jB))
		}
	}
	if err != nil {
		goto bailout
	}

	// mark job as finished
	if _, err = tx.Exec(
		tkStmtFinishJob,
		q.JobId.String(),
		time.Now().UTC(),
		"success",
	); err != nil {
		goto bailout
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		goto bailout
	}
	log.Printf("SUCCESS - Finished job: %s\n", q.JobId.String())

	// accept tree changes
	tk.tree.Commit()
	return

bailout:
	log.Printf("FAILED - Finished job: %s\n", q.JobId.String())
	log.Println(err)
	tk.tree.Rollback()
	tx.Rollback()
	tk.conn.Exec(
		tkStmtFinishJob,
		q.JobId.String(),
		time.Now().UTC(),
		"failed",
	)
	for i := 0; i < len(tk.actionChan); i++ {
		a := <-tk.actionChan
		jB, _ := json.Marshal(a)
		log.Printf("Cleaned message: %s\n", string(jB))
	}
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
