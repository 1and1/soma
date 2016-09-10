package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/satori/go.uuid"

)

type treeRequest struct {
	RequestType string
	Action      string
	User        string
	JobId       uuid.UUID
	reply       chan somaResult
	Repository  somaRepositoryRequest
	Bucket      somaBucketRequest
	Group       somaGroupRequest
	Cluster     somaClusterRequest
	Node        somaNodeRequest
	CheckConfig somaCheckConfigRequest
}

type treeResult struct {
	ResultType  string
	ResultError error
	JobId       uuid.UUID
	Repository  somaRepositoryResult
	Bucket      somaRepositoryRequest
}

type treeKeeper struct {
	repoId     string
	repoName   string
	team       string
	broken     bool
	ready      bool
	stopped    bool
	frozen     bool
	rebuild    bool
	rbLevel    string
	input      chan treeRequest
	shutdown   chan bool
	stopchan   chan bool
	conn       *sql.DB
	tree       *tree.Tree
	errChan    chan *tree.Error
	actionChan chan *tree.Action
	start_job  *sql.Stmt
	get_view   *sql.Stmt
}

func (tk *treeKeeper) run() {
	log.Printf("Starting TreeKeeper for Repo %s (%s)", tk.repoName, tk.repoId)
	tk.startupLoad()
	var err error

	// treekeepers have a dedicated connection pool
	defer tk.conn.Close()

	// if this was a rebuild, simply return if it failed
	if tk.broken && tk.rebuild {
		return
	}

	// rebuild was successful, process events from initial loading
	// then exit
	if tk.rebuild {
		req := treeRequest{
			RequestType: `rebuild`,
			Action:      `rebuild`,
			JobId:       uuid.NewV4(),
		}
		tk.process(&req)
		return
	}

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
			case <-tk.stopchan:
				tk.stop()
				goto stopsign
			}
		}
		return
	}

	if tk.start_job, err = tk.conn.Prepare(tkStmtStartJob); err != nil {
		log.Fatal("treekeeper/start-job: ", err)
	}
	defer tk.start_job.Close()

	if tk.get_view, err = tk.conn.Prepare(tkStmtGetViewFromCapability); err != nil {
		log.Fatal("treekeeper/get-view-by-capability: ", err)
	}
	defer tk.get_view.Close()

	log.Printf("TK[%s]: ready for service!\n", tk.repoName)
	tk.ready = true

	if SomaCfg.Observer {
		log.Printf("TreeKeeper [%s] entered observer mode\n", tk.repoName)
		<-tk.shutdown
		goto exit
	}

stopsign:
	if tk.stopped {
		<-tk.shutdown
		goto exit
	}
runloop:
	for {
		select {
		case <-tk.shutdown:
			break runloop
		case <-tk.stopchan:
			tk.stop()
			goto stopsign
		case req := <-tk.input:
			tk.process(&req)
			handlerMap[`jobDelay`].(jobDelay).notify <- req.JobId.String()
			if !tk.frozen {
				tk.buildDeploymentDetails()
				tk.orderDeploymentDetails()
			}
		}
	}
exit:
}

func (tk *treeKeeper) isReady() bool {
	return tk.ready
}

func (tk *treeKeeper) isBroken() bool {
	return tk.broken
}

func (tk *treeKeeper) stop() {
	tk.stopped = true
	tk.ready = false
	tk.broken = false
}

func (tk *treeKeeper) isStopped() bool {
	return tk.stopped
}

func (tk *treeKeeper) process(q *treeRequest) {
	var (
		err        error
		hasErrors  bool
		tx         *sql.Tx
		treeCheck  *tree.Check
		nullBucket sql.NullString

		txStmtPropertyInstanceCreate *sql.Stmt

		txStmtRepositoryPropertyServiceCreate *sql.Stmt
		txStmtRepositoryPropertySystemCreate  *sql.Stmt
		txStmtRepositoryPropertyOncallCreate  *sql.Stmt
		txStmtRepositoryPropertyCustomCreate  *sql.Stmt

		txStmtCreateBucket                *sql.Stmt
		txStmtBucketPropertyServiceCreate *sql.Stmt
		txStmtBucketPropertySystemCreate  *sql.Stmt
		txStmtBucketPropertyOncallCreate  *sql.Stmt
		txStmtBucketPropertyCustomCreate  *sql.Stmt

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

		txStmtClusterCreate                *sql.Stmt
		txStmtClusterUpdate                *sql.Stmt
		txStmtClusterDelete                *sql.Stmt
		txStmtClusterMemberNew             *sql.Stmt
		txStmtClusterMemberRemove          *sql.Stmt
		txStmtClusterPropertyServiceCreate *sql.Stmt
		txStmtClusterPropertySystemCreate  *sql.Stmt
		txStmtClusterPropertyOncallCreate  *sql.Stmt
		txStmtClusterPropertyCustomCreate  *sql.Stmt

		txStmtBucketAssignNode          *sql.Stmt
		txStmtUpdateNodeState           *sql.Stmt
		txStmtNodeUnassignFromBucket    *sql.Stmt
		txStmtNodePropertyServiceCreate *sql.Stmt
		txStmtNodePropertySystemCreate  *sql.Stmt
		txStmtNodePropertyOncallCreate  *sql.Stmt
		txStmtNodePropertyCustomCreate  *sql.Stmt

		txStmtCreateCheckConfigurationBase                *sql.Stmt
		txStmtCreateCheckConfigurationThreshold           *sql.Stmt
		txStmtCreateCheckConfigurationConstraintSystem    *sql.Stmt
		txStmtCreateCheckConfigurationConstraintNative    *sql.Stmt
		txStmtCreateCheckConfigurationConstraintOncall    *sql.Stmt
		txStmtCreateCheckConfigurationConstraintCustom    *sql.Stmt
		txStmtCreateCheckConfigurationConstraintService   *sql.Stmt
		txStmtCreateCheckConfigurationConstraintAttribute *sql.Stmt
		txStmtCreateCheck                                 *sql.Stmt
		txStmtCreateCheckInstance                         *sql.Stmt
		txStmtCreateCheckInstanceConfiguration            *sql.Stmt
		txStmtDeleteCheck                                 *sql.Stmt
		txStmtDeleteCheckInstance                         *sql.Stmt
	)
	if !tk.rebuild {
		_, err = tk.start_job.Exec(q.JobId.String(), time.Now().UTC())
		if err != nil {
			log.Println(err)
		}
		log.Printf("Processing job: %s\n", q.JobId.String())
	} else {
		log.Printf("Processing rebuild job: %s\n", q.JobId.String())
	}

	tk.tree.Begin()

	// q.Action == `rebuild` will fall through switch
	switch q.Action {
	//
	// REPOSITORY MANIPULATION REQUESTS
	case "add_system_property_to_repository":
		tk.tree.Find(tree.FindRequest{
			ElementType: "repository",
			ElementId:   q.Repository.Repository.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertySystem{
			Id:           uuid.NewV4(),
			Inheritance:  (*q.Repository.Repository.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Repository.Repository.Properties)[0].ChildrenOnly,
			View:         (*q.Repository.Repository.Properties)[0].View,
			Key:          (*q.Repository.Repository.Properties)[0].System.Name,
			Value:        (*q.Repository.Repository.Properties)[0].System.Value,
		})
	case `delete_system_property_from_repository`:
		srcUUID, _ := uuid.FromString((*q.Repository.Repository.Properties)[0].SourceInstanceId)

		tk.tree.Find(tree.FindRequest{
			ElementType: `repository`,
			ElementId:   q.Repository.Repository.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertySystem{
			SourceId: srcUUID,
			View:     (*q.Repository.Repository.Properties)[0].View,
			Key:      (*q.Repository.Repository.Properties)[0].System.Name,
			Value:    (*q.Repository.Repository.Properties)[0].System.Value,
		})
	case "add_service_property_to_repository":
		tk.tree.Find(tree.FindRequest{
			ElementType: "repository",
			ElementId:   q.Repository.Repository.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyService{
			Id:           uuid.NewV4(),
			Inheritance:  (*q.Repository.Repository.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Repository.Repository.Properties)[0].ChildrenOnly,
			View:         (*q.Repository.Repository.Properties)[0].View,
			Service:      (*q.Repository.Repository.Properties)[0].Service.Name,
			Attributes:   (*q.Repository.Repository.Properties)[0].Service.Attributes,
		})
	case `delete_service_property_from_repository`:
		srcUUID, _ := uuid.FromString((*q.Repository.Repository.Properties)[0].SourceInstanceId)

		tk.tree.Find(tree.FindRequest{
			ElementType: "repository",
			ElementId:   q.Repository.Repository.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyService{
			SourceId: srcUUID,
			View:     (*q.Repository.Repository.Properties)[0].View,
			Service:  (*q.Repository.Repository.Properties)[0].Service.Name,
		})
	case "add_oncall_property_to_repository":
		oncallId, _ := uuid.FromString((*q.Repository.Repository.Properties)[0].Oncall.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: "repository",
			ElementId:   q.Repository.Repository.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyOncall{
			Id:           uuid.NewV4(),
			OncallId:     oncallId,
			Inheritance:  (*q.Repository.Repository.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Repository.Repository.Properties)[0].ChildrenOnly,
			View:         (*q.Repository.Repository.Properties)[0].View,
			Name:         (*q.Repository.Repository.Properties)[0].Oncall.Name,
			Number:       (*q.Repository.Repository.Properties)[0].Oncall.Number,
		})
	case `delete_oncall_property_from_repository`:
		srcUUID, _ := uuid.FromString((*q.Repository.Repository.Properties)[0].SourceInstanceId)
		oncallId, _ := uuid.FromString((*q.Repository.Repository.Properties)[0].Oncall.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: "repository",
			ElementId:   q.Repository.Repository.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyOncall{
			SourceId: srcUUID,
			OncallId: oncallId,
			View:     (*q.Repository.Repository.Properties)[0].View,
			Name:     (*q.Repository.Repository.Properties)[0].Oncall.Name,
			Number:   (*q.Repository.Repository.Properties)[0].Oncall.Number,
		})
	case "add_custom_property_to_repository":
		customId, _ := uuid.FromString((*q.Repository.Repository.Properties)[0].Custom.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: "repository",
			ElementId:   q.Repository.Repository.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyCustom{
			Id:           uuid.NewV4(),
			CustomId:     customId,
			Inheritance:  (*q.Repository.Repository.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Repository.Repository.Properties)[0].ChildrenOnly,
			View:         (*q.Repository.Repository.Properties)[0].View,
			Key:          (*q.Repository.Repository.Properties)[0].Custom.Name,
			Value:        (*q.Repository.Repository.Properties)[0].Custom.Value,
		})
	case `delete_custom_property_from_repository`:
		srcUUID, _ := uuid.FromString((*q.Repository.Repository.Properties)[0].SourceInstanceId)
		customId, _ := uuid.FromString((*q.Repository.Repository.Properties)[0].Custom.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: "repository",
			ElementId:   q.Repository.Repository.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyCustom{
			SourceId: srcUUID,
			CustomId: customId,
			View:     (*q.Repository.Repository.Properties)[0].View,
			Key:      (*q.Repository.Repository.Properties)[0].Custom.Name,
			Value:    (*q.Repository.Repository.Properties)[0].Custom.Value,
		})

	//
	// BUCKET MANIPULATION REQUESTS
	case "create_bucket":
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
			ParentType: "repository",
			ParentId:   tk.repoId,
			ParentName: tk.repoName,
		})
	case "add_system_property_to_bucket":
		tk.tree.Find(tree.FindRequest{
			ElementType: "bucket",
			ElementId:   q.Bucket.Bucket.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertySystem{
			Id:           uuid.NewV4(),
			Inheritance:  (*q.Bucket.Bucket.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Bucket.Bucket.Properties)[0].ChildrenOnly,
			View:         (*q.Bucket.Bucket.Properties)[0].View,
			Key:          (*q.Bucket.Bucket.Properties)[0].System.Name,
			Value:        (*q.Bucket.Bucket.Properties)[0].System.Value,
		})
	case `delete_system_property_from_bucket`:
		srcUUID, _ := uuid.FromString((*q.Bucket.Bucket.Properties)[0].SourceInstanceId)

		tk.tree.Find(tree.FindRequest{
			ElementType: `bucket`,
			ElementId:   q.Bucket.Bucket.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertySystem{
			SourceId: srcUUID,
			View:     (*q.Bucket.Bucket.Properties)[0].View,
			Key:      (*q.Bucket.Bucket.Properties)[0].System.Name,
			Value:    (*q.Bucket.Bucket.Properties)[0].System.Value,
		})
	case "add_service_property_to_bucket":
		tk.tree.Find(tree.FindRequest{
			ElementType: "bucket",
			ElementId:   q.Bucket.Bucket.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyService{
			Id:           uuid.NewV4(),
			Inheritance:  (*q.Bucket.Bucket.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Bucket.Bucket.Properties)[0].ChildrenOnly,
			View:         (*q.Bucket.Bucket.Properties)[0].View,
			Service:      (*q.Bucket.Bucket.Properties)[0].Service.Name,
			Attributes:   (*q.Bucket.Bucket.Properties)[0].Service.Attributes,
		})
	case `delete_service_property_from_bucket`:
		srcUUID, _ := uuid.FromString((*q.Bucket.Bucket.Properties)[0].SourceInstanceId)

		tk.tree.Find(tree.FindRequest{
			ElementType: `bucket`,
			ElementId:   q.Bucket.Bucket.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyService{
			SourceId: srcUUID,
			View:     (*q.Bucket.Bucket.Properties)[0].View,
			Service:  (*q.Bucket.Bucket.Properties)[0].Service.Name,
		})
	case "add_oncall_property_to_bucket":
		oncallId, _ := uuid.FromString((*q.Bucket.Bucket.Properties)[0].Oncall.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: "bucket",
			ElementId:   q.Bucket.Bucket.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyOncall{
			Id:           uuid.NewV4(),
			OncallId:     oncallId,
			Inheritance:  (*q.Bucket.Bucket.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Bucket.Bucket.Properties)[0].ChildrenOnly,
			View:         (*q.Bucket.Bucket.Properties)[0].View,
			Name:         (*q.Bucket.Bucket.Properties)[0].Oncall.Name,
			Number:       (*q.Bucket.Bucket.Properties)[0].Oncall.Number,
		})
	case `delete_oncall_property_from_bucket`:
		srcUUID, _ := uuid.FromString((*q.Bucket.Bucket.Properties)[0].SourceInstanceId)
		oncallId, _ := uuid.FromString((*q.Bucket.Bucket.Properties)[0].Oncall.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: `bucket`,
			ElementId:   q.Bucket.Bucket.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyOncall{
			SourceId: srcUUID,
			OncallId: oncallId,
			View:     (*q.Bucket.Bucket.Properties)[0].View,
			Name:     (*q.Bucket.Bucket.Properties)[0].Oncall.Name,
			Number:   (*q.Bucket.Bucket.Properties)[0].Oncall.Number,
		})
	case "add_custom_property_to_bucket":
		customId, _ := uuid.FromString((*q.Bucket.Bucket.Properties)[0].Custom.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: "bucket",
			ElementId:   q.Bucket.Bucket.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyCustom{
			Id:           uuid.NewV4(),
			CustomId:     customId,
			Inheritance:  (*q.Bucket.Bucket.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Bucket.Bucket.Properties)[0].ChildrenOnly,
			View:         (*q.Bucket.Bucket.Properties)[0].View,
			Key:          (*q.Bucket.Bucket.Properties)[0].Custom.Name,
			Value:        (*q.Bucket.Bucket.Properties)[0].Custom.Value,
		})
	case `delete_custom_property_from_bucket`:
		srcUUID, _ := uuid.FromString((*q.Bucket.Bucket.Properties)[0].SourceInstanceId)
		customId, _ := uuid.FromString((*q.Bucket.Bucket.Properties)[0].Custom.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: `bucket`,
			ElementId:   q.Bucket.Bucket.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyCustom{
			SourceId: srcUUID,
			CustomId: customId,
			View:     (*q.Bucket.Bucket.Properties)[0].View,
			Key:      (*q.Bucket.Bucket.Properties)[0].Custom.Name,
			Value:    (*q.Bucket.Bucket.Properties)[0].Custom.Value,
		})

	//
	// GROUP MANIPULATION REQUESTS
	case "create_group":
		tree.NewGroup(tree.GroupSpec{
			Id:   uuid.NewV4().String(),
			Name: q.Group.Group.Name,
			Team: tk.team,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "bucket",
			ParentId:   q.Group.Group.BucketId,
		})
	case "delete_group":
		tk.tree.Find(tree.FindRequest{
			ElementType: "group",
			ElementId:   q.Group.Group.Id,
		}, true).(tree.BucketAttacher).Destroy()
	case "reset_group_to_bucket":
		tk.tree.Find(tree.FindRequest{
			ElementType: "group",
			ElementId:   q.Group.Group.Id,
		}, true).(tree.BucketAttacher).Detach()
	case "add_group_to_group":
		tk.tree.Find(tree.FindRequest{
			ElementType: "group",
			ElementId:   (*q.Group.Group.MemberGroups)[0].Id,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentId:   q.Group.Group.Id,
		})
	case "add_system_property_to_group":
		tk.tree.Find(tree.FindRequest{
			ElementType: "group",
			ElementId:   q.Group.Group.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertySystem{
			Id:           uuid.NewV4(),
			Inheritance:  (*q.Group.Group.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Group.Group.Properties)[0].ChildrenOnly,
			View:         (*q.Group.Group.Properties)[0].View,
			Key:          (*q.Group.Group.Properties)[0].System.Name,
			Value:        (*q.Group.Group.Properties)[0].System.Value,
		})
	case `delete_system_property_from_group`:
		srcUUID, _ := uuid.FromString((*q.Group.Group.Properties)[0].SourceInstanceId)

		tk.tree.Find(tree.FindRequest{
			ElementType: `group`,
			ElementId:   q.Group.Group.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertySystem{
			SourceId: srcUUID,
			View:     (*q.Group.Group.Properties)[0].View,
			Key:      (*q.Group.Group.Properties)[0].System.Name,
			Value:    (*q.Group.Group.Properties)[0].System.Value,
		})
	case "add_service_property_to_group":
		tk.tree.Find(tree.FindRequest{
			ElementType: "group",
			ElementId:   q.Group.Group.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyService{
			Id:           uuid.NewV4(),
			Inheritance:  (*q.Group.Group.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Group.Group.Properties)[0].ChildrenOnly,
			View:         (*q.Group.Group.Properties)[0].View,
			Service:      (*q.Group.Group.Properties)[0].Service.Name,
			Attributes:   (*q.Group.Group.Properties)[0].Service.Attributes,
		})
	case `delete_service_property_from_group`:
		srcUUID, _ := uuid.FromString((*q.Group.Group.Properties)[0].SourceInstanceId)

		tk.tree.Find(tree.FindRequest{
			ElementType: `group`,
			ElementId:   q.Group.Group.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyService{
			SourceId: srcUUID,
			View:     (*q.Group.Group.Properties)[0].View,
			Service:  (*q.Group.Group.Properties)[0].Service.Name,
		})
	case "add_oncall_property_to_group":
		oncallId, _ := uuid.FromString((*q.Group.Group.Properties)[0].Oncall.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: "group",
			ElementId:   q.Group.Group.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyOncall{
			Id:           uuid.NewV4(),
			OncallId:     oncallId,
			Inheritance:  (*q.Group.Group.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Group.Group.Properties)[0].ChildrenOnly,
			View:         (*q.Group.Group.Properties)[0].View,
			Name:         (*q.Group.Group.Properties)[0].Oncall.Name,
			Number:       (*q.Group.Group.Properties)[0].Oncall.Number,
		})
	case `delete_oncall_property_from_group`:
		srcUUID, _ := uuid.FromString((*q.Group.Group.Properties)[0].SourceInstanceId)
		oncallId, _ := uuid.FromString((*q.Group.Group.Properties)[0].Oncall.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: `group`,
			ElementId:   q.Group.Group.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyOncall{
			SourceId: srcUUID,
			OncallId: oncallId,
			View:     (*q.Group.Group.Properties)[0].View,
			Name:     (*q.Group.Group.Properties)[0].Oncall.Name,
			Number:   (*q.Group.Group.Properties)[0].Oncall.Number,
		})
	case "add_custom_property_to_group":
		customId, _ := uuid.FromString((*q.Group.Group.Properties)[0].Custom.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: "group",
			ElementId:   q.Group.Group.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyCustom{
			Id:           uuid.NewV4(),
			CustomId:     customId,
			Inheritance:  (*q.Group.Group.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Group.Group.Properties)[0].ChildrenOnly,
			View:         (*q.Group.Group.Properties)[0].View,
			Key:          (*q.Group.Group.Properties)[0].Custom.Name,
			Value:        (*q.Group.Group.Properties)[0].Custom.Value,
		})
	case `delete_custom_property_from_group`:
		srcUUID, _ := uuid.FromString((*q.Group.Group.Properties)[0].SourceInstanceId)
		customId, _ := uuid.FromString((*q.Group.Group.Properties)[0].Custom.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: `group`,
			ElementId:   q.Group.Group.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyCustom{
			SourceId: srcUUID,
			CustomId: customId,
			View:     (*q.Group.Group.Properties)[0].View,
			Key:      (*q.Group.Group.Properties)[0].Custom.Name,
			Value:    (*q.Group.Group.Properties)[0].Custom.Value,
		})

	//
	// CLUSTER MANIPULATION REQUESTS
	case "create_cluster":
		tree.NewCluster(tree.ClusterSpec{
			Id:   uuid.NewV4().String(),
			Name: q.Cluster.Cluster.Name,
			Team: tk.team,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "bucket",
			ParentId:   q.Cluster.Cluster.BucketId,
		})
	case "delete_cluster":
		tk.tree.Find(tree.FindRequest{
			ElementType: "cluster",
			ElementId:   q.Cluster.Cluster.Id,
		}, true).(tree.BucketAttacher).Destroy()
	case "reset_cluster_to_bucket":
		tk.tree.Find(tree.FindRequest{
			ElementType: "cluster",
			ElementId:   q.Cluster.Cluster.Id,
		}, true).(tree.BucketAttacher).Detach()
	case "add_cluster_to_group":
		tk.tree.Find(tree.FindRequest{
			ElementType: "cluster",
			ElementId:   (*q.Group.Group.MemberClusters)[0].Id,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentId:   q.Group.Group.Id,
		})
	case "add_system_property_to_cluster":
		tk.tree.Find(tree.FindRequest{
			ElementType: "cluster",
			ElementId:   q.Cluster.Cluster.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertySystem{
			Id:           uuid.NewV4(),
			Inheritance:  (*q.Cluster.Cluster.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Cluster.Cluster.Properties)[0].ChildrenOnly,
			View:         (*q.Cluster.Cluster.Properties)[0].View,
			Key:          (*q.Cluster.Cluster.Properties)[0].System.Name,
			Value:        (*q.Cluster.Cluster.Properties)[0].System.Value,
		})
	case `delete_system_property_from_cluster`:
		srcUUID, _ := uuid.FromString((*q.Cluster.Cluster.Properties)[0].SourceInstanceId)

		tk.tree.Find(tree.FindRequest{
			ElementType: `cluster`,
			ElementId:   q.Cluster.Cluster.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertySystem{
			SourceId: srcUUID,
			View:     (*q.Cluster.Cluster.Properties)[0].View,
			Key:      (*q.Cluster.Cluster.Properties)[0].System.Name,
			Value:    (*q.Cluster.Cluster.Properties)[0].System.Value,
		})
	case "add_service_property_to_cluster":
		tk.tree.Find(tree.FindRequest{
			ElementType: "cluster",
			ElementId:   q.Cluster.Cluster.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyService{
			Id:           uuid.NewV4(),
			Inheritance:  (*q.Cluster.Cluster.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Cluster.Cluster.Properties)[0].ChildrenOnly,
			View:         (*q.Cluster.Cluster.Properties)[0].View,
			Service:      (*q.Cluster.Cluster.Properties)[0].Service.Name,
			Attributes:   (*q.Cluster.Cluster.Properties)[0].Service.Attributes,
		})
	case `delete_service_property_from_cluster`:
		srcUUID, _ := uuid.FromString((*q.Cluster.Cluster.Properties)[0].SourceInstanceId)

		tk.tree.Find(tree.FindRequest{
			ElementType: `cluster`,
			ElementId:   q.Cluster.Cluster.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyService{
			SourceId: srcUUID,
			View:     (*q.Cluster.Cluster.Properties)[0].View,
			Service:  (*q.Cluster.Cluster.Properties)[0].Service.Name,
		})
	case "add_oncall_property_to_cluster":
		oncallId, _ := uuid.FromString((*q.Cluster.Cluster.Properties)[0].Oncall.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: "cluster",
			ElementId:   q.Cluster.Cluster.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyOncall{
			Id:           uuid.NewV4(),
			OncallId:     oncallId,
			Inheritance:  (*q.Cluster.Cluster.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Cluster.Cluster.Properties)[0].ChildrenOnly,
			View:         (*q.Cluster.Cluster.Properties)[0].View,
			Name:         (*q.Cluster.Cluster.Properties)[0].Oncall.Name,
			Number:       (*q.Cluster.Cluster.Properties)[0].Oncall.Number,
		})
	case `delete_oncall_property_from_cluster`:
		srcUUID, _ := uuid.FromString((*q.Cluster.Cluster.Properties)[0].SourceInstanceId)
		oncallId, _ := uuid.FromString((*q.Cluster.Cluster.Properties)[0].Oncall.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: `cluster`,
			ElementId:   q.Cluster.Cluster.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyOncall{
			SourceId: srcUUID,
			OncallId: oncallId,
			View:     (*q.Cluster.Cluster.Properties)[0].View,
			Name:     (*q.Cluster.Cluster.Properties)[0].Oncall.Name,
			Number:   (*q.Cluster.Cluster.Properties)[0].Oncall.Number,
		})
	case "add_custom_property_to_cluster":
		customId, _ := uuid.FromString((*q.Cluster.Cluster.Properties)[0].Custom.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: "cluster",
			ElementId:   q.Cluster.Cluster.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyCustom{
			Id:           uuid.NewV4(),
			CustomId:     customId,
			Inheritance:  (*q.Cluster.Cluster.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Cluster.Cluster.Properties)[0].ChildrenOnly,
			View:         (*q.Cluster.Cluster.Properties)[0].View,
			Key:          (*q.Cluster.Cluster.Properties)[0].Custom.Name,
			Value:        (*q.Cluster.Cluster.Properties)[0].Custom.Value,
		})
	case `delete_custom_property_from_cluster`:
		srcUUID, _ := uuid.FromString((*q.Cluster.Cluster.Properties)[0].SourceInstanceId)
		customId, _ := uuid.FromString((*q.Cluster.Cluster.Properties)[0].Custom.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: `cluster`,
			ElementId:   q.Cluster.Cluster.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyCustom{
			SourceId: srcUUID,
			CustomId: customId,
			View:     (*q.Cluster.Cluster.Properties)[0].View,
			Key:      (*q.Cluster.Cluster.Properties)[0].Custom.Name,
			Value:    (*q.Cluster.Cluster.Properties)[0].Custom.Value,
		})

	//
	// NODE MANIPULATION REQUESTS
	case "assign_node":
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
			ParentType: "bucket",
			ParentId:   q.Node.Node.Config.BucketId,
		})
	case "delete_node":
		tk.tree.Find(tree.FindRequest{
			ElementType: "node",
			ElementId:   q.Node.Node.Id,
		}, true).(tree.BucketAttacher).Destroy()
	case "reset_node_to_bucket":
		tk.tree.Find(tree.FindRequest{
			ElementType: "node",
			ElementId:   q.Node.Node.Id,
		}, true).(tree.BucketAttacher).Detach()
	case "add_node_to_group":
		tk.tree.Find(tree.FindRequest{
			ElementType: "node",
			ElementId:   (*q.Group.Group.MemberNodes)[0].Id,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentId:   q.Group.Group.Id,
		})
	case "add_node_to_cluster":
		tk.tree.Find(tree.FindRequest{
			ElementType: "node",
			ElementId:   (*q.Cluster.Cluster.Members)[0].Id,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "cluster",
			ParentId:   q.Cluster.Cluster.Id,
		})
	case "add_system_property_to_node":
		tk.tree.Find(tree.FindRequest{
			ElementType: "node",
			ElementId:   q.Node.Node.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertySystem{
			Id:           uuid.NewV4(),
			Inheritance:  (*q.Node.Node.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Node.Node.Properties)[0].ChildrenOnly,
			View:         (*q.Node.Node.Properties)[0].View,
			Key:          (*q.Node.Node.Properties)[0].System.Name,
			Value:        (*q.Node.Node.Properties)[0].System.Value,
		})
	case `delete_system_property_from_node`:
		srcUUID, _ := uuid.FromString((*q.Node.Node.Properties)[0].SourceInstanceId)

		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementId:   q.Node.Node.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertySystem{
			SourceId: srcUUID,
			View:     (*q.Node.Node.Properties)[0].View,
			Key:      (*q.Node.Node.Properties)[0].System.Name,
			Value:    (*q.Node.Node.Properties)[0].System.Value,
		})
	case "add_service_property_to_node":
		tk.tree.Find(tree.FindRequest{
			ElementType: "node",
			ElementId:   q.Node.Node.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyService{
			Id:           uuid.NewV4(),
			Inheritance:  (*q.Node.Node.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Node.Node.Properties)[0].ChildrenOnly,
			View:         (*q.Node.Node.Properties)[0].View,
			Service:      (*q.Node.Node.Properties)[0].Service.Name,
			Attributes:   (*q.Node.Node.Properties)[0].Service.Attributes,
		})
	case `delete_service_property_from_node`:
		srcUUID, _ := uuid.FromString((*q.Node.Node.Properties)[0].SourceInstanceId)

		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementId:   q.Node.Node.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyService{
			SourceId: srcUUID,
			View:     (*q.Node.Node.Properties)[0].View,
			Service:  (*q.Node.Node.Properties)[0].Service.Name,
		})
	case "add_oncall_property_to_node":
		oncallId, _ := uuid.FromString((*q.Node.Node.Properties)[0].Oncall.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: "node",
			ElementId:   q.Node.Node.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyOncall{
			Id:           uuid.NewV4(),
			OncallId:     oncallId,
			Inheritance:  (*q.Node.Node.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Node.Node.Properties)[0].ChildrenOnly,
			View:         (*q.Node.Node.Properties)[0].View,
			Name:         (*q.Node.Node.Properties)[0].Oncall.Name,
			Number:       (*q.Node.Node.Properties)[0].Oncall.Number,
		})
	case `delete_oncall_property_from_node`:
		srcUUID, _ := uuid.FromString((*q.Node.Node.Properties)[0].SourceInstanceId)
		oncallId, _ := uuid.FromString((*q.Node.Node.Properties)[0].Oncall.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementId:   q.Node.Node.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyOncall{
			SourceId: srcUUID,
			OncallId: oncallId,
			View:     (*q.Node.Node.Properties)[0].View,
			Name:     (*q.Node.Node.Properties)[0].Oncall.Name,
			Number:   (*q.Node.Node.Properties)[0].Oncall.Number,
		})
	case "add_custom_property_to_node":
		customId, _ := uuid.FromString((*q.Node.Node.Properties)[0].Custom.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: "node",
			ElementId:   q.Node.Node.Id,
		}, true).(tree.Propertier).SetProperty(&tree.PropertyCustom{
			Id:           uuid.NewV4(),
			CustomId:     customId,
			Inheritance:  (*q.Node.Node.Properties)[0].Inheritance,
			ChildrenOnly: (*q.Node.Node.Properties)[0].ChildrenOnly,
			View:         (*q.Node.Node.Properties)[0].View,
			Key:          (*q.Node.Node.Properties)[0].Custom.Name,
			Value:        (*q.Node.Node.Properties)[0].Custom.Value,
		})
	case `delete_custom_property_from_node`:
		srcUUID, _ := uuid.FromString((*q.Node.Node.Properties)[0].SourceInstanceId)
		customId, _ := uuid.FromString((*q.Node.Node.Properties)[0].Custom.Id)

		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementId:   q.Node.Node.Id,
		}, true).(tree.Propertier).DeleteProperty(&tree.PropertyCustom{
			SourceId: srcUUID,
			CustomId: customId,
			View:     (*q.Node.Node.Properties)[0].View,
			Key:      (*q.Node.Node.Properties)[0].Custom.Name,
			Value:    (*q.Node.Node.Properties)[0].Custom.Value,
		})

	//
	// CHECK MANIPULATION REQUESTS
	case `add_check_to_repository`:
		fallthrough
	case `add_check_to_bucket`:
		fallthrough
	case `add_check_to_group`:
		fallthrough
	case `add_check_to_cluster`:
		fallthrough
	case `add_check_to_node`:
		if treeCheck, err = tk.convertCheck(&q.CheckConfig.CheckConfig); err == nil {
			tk.tree.Find(tree.FindRequest{
				ElementType: q.CheckConfig.CheckConfig.ObjectType,
				ElementId:   q.CheckConfig.CheckConfig.ObjectId,
			}, true).SetCheck(*treeCheck)
		}
	case `remove_check_from_repository`:
		fallthrough
	case `remove_check_from_bucket`:
		fallthrough
	case `remove_check_from_group`:
		fallthrough
	case `remove_check_from_cluster`:
		fallthrough
	case `remove_check_from_node`:
		if treeCheck, err = tk.convertCheckForDelete(&q.CheckConfig.CheckConfig); err == nil {
			tk.tree.Find(tree.FindRequest{
				ElementType: q.CheckConfig.CheckConfig.ObjectType,
				ElementId:   q.CheckConfig.CheckConfig.ObjectId,
			}, true).DeleteCheck(*treeCheck)
		}
	}

	// check if we accumulated an error in one of the switch cases
	if err != nil {
		goto bailout
	}

	// recalculate check instances
	tk.tree.ComputeCheckInstances()

	// open multi-statement transaction
	if tx, err = tk.conn.Begin(); err != nil {
		goto bailout
	}
	defer tx.Rollback()

	// prepare statements within tx context
	if txStmtPropertyInstanceCreate, err = tx.Prepare(tkStmtPropertyInstanceCreate); err != nil {
		log.Println("Failed to prepare: tkStmtPropertyInstanceCreate")
		goto bailout
	}
	defer txStmtPropertyInstanceCreate.Close()

	if txStmtCreateCheckConfigurationBase, err = tx.Prepare(tkStmtCreateCheckConfigurationBase); err != nil {
		log.Println("Failed to prepare: tkStmtCreateCheckConfigurationBase")
		goto bailout
	}
	defer txStmtCreateCheckConfigurationBase.Close()

	if txStmtCreateCheckConfigurationThreshold, err = tx.Prepare(tkStmtCreateCheckConfigurationThreshold); err != nil {
		log.Println("Failed to prepare: tkStmtCreateCheckConfigurationThreshold")
		goto bailout
	}
	defer txStmtCreateCheckConfigurationThreshold.Close()

	if txStmtCreateCheckConfigurationConstraintSystem, err = tx.Prepare(tkStmtCreateCheckConfigurationConstraintSystem); err != nil {
		log.Println("Failed to prepare: tkStmtCreateCheckConfigurationConstraintSystem")
		goto bailout
	}
	defer txStmtCreateCheckConfigurationConstraintSystem.Close()

	if txStmtCreateCheckConfigurationConstraintNative, err = tx.Prepare(tkStmtCreateCheckConfigurationConstraintNative); err != nil {
		log.Println("Failed to prepare: tkStmtCreateCheckConfigurationConstraintNative")
		goto bailout
	}
	defer txStmtCreateCheckConfigurationConstraintNative.Close()

	if txStmtCreateCheckConfigurationConstraintOncall, err = tx.Prepare(tkStmtCreateCheckConfigurationConstraintOncall); err != nil {
		log.Println("Failed to prepare: tkStmtCreateCheckConfigurationConstraintOncall")
		goto bailout
	}
	defer txStmtCreateCheckConfigurationConstraintOncall.Close()

	if txStmtCreateCheckConfigurationConstraintCustom, err = tx.Prepare(tkStmtCreateCheckConfigurationConstraintCustom); err != nil {
		log.Println("Failed to prepare: tkStmtCreateCheckConfigurationConstraintCustom")
		goto bailout
	}
	defer txStmtCreateCheckConfigurationConstraintCustom.Close()

	if txStmtCreateCheckConfigurationConstraintService, err = tx.Prepare(tkStmtCreateCheckConfigurationConstraintService); err != nil {
		log.Println("Failed to prepare: tkStmtCreateCheckConfigurationConstraintService")
		goto bailout
	}
	defer txStmtCreateCheckConfigurationConstraintService.Close()

	if txStmtCreateCheckConfigurationConstraintAttribute, err = tx.Prepare(tkStmtCreateCheckConfigurationConstraintAttribute); err != nil {
		log.Println("Failed to prepare: tkStmtCreateCheckConfigurationConstraintAttribute")
		goto bailout
	}
	defer txStmtCreateCheckConfigurationConstraintAttribute.Close()

	if txStmtCreateCheck, err = tx.Prepare(tkStmtCreateCheck); err != nil {
		log.Println("Failed to prepare: tkStmtCreateCheck")
		goto bailout
	}
	defer txStmtCreateCheck.Close()

	if txStmtCreateCheckInstance, err = tx.Prepare(tkStmtCreateCheckInstance); err != nil {
		log.Println("Failed to prepare: tkStmtCreateCheckInstance")
		goto bailout
	}
	defer txStmtCreateCheckInstance.Close()

	if txStmtCreateCheckInstanceConfiguration, err = tx.Prepare(tkStmtCreateCheckInstanceConfiguration); err != nil {
		log.Println("Failed to prepare: tkStmtCreateCheckInstanceConfiguration")
		goto bailout
	}
	defer txStmtCreateCheckInstanceConfiguration.Close()

	if txStmtDeleteCheck, err = tx.Prepare(stmt.TxMarkCheckDeleted); err != nil {
		log.Println("Failed to prepare: txStmtDeleteCheck")
		goto bailout
	}

	if txStmtDeleteCheckInstance, err = tx.Prepare(stmt.TxMarkCheckInstanceDeleted); err != nil {
		log.Println("Failed to prepare: txStmtDeleteCheckInstance")
		goto bailout
	}

	//
	// REPOSITORY
	if txStmtRepositoryPropertyOncallCreate, err = tx.Prepare(tkStmtRepositoryPropertyOncallCreate); err != nil {
		log.Println("Failed to prepare: tkStmtRepositoryPropertyOncallCreate")
		goto bailout
	}
	defer txStmtRepositoryPropertyOncallCreate.Close()

	if txStmtRepositoryPropertyServiceCreate, err = tx.Prepare(tkStmtRepositoryPropertyServiceCreate); err != nil {
		log.Println("Failed to prepare: tkStmtRepositoryPropertyServiceCreate")
		goto bailout
	}
	defer txStmtRepositoryPropertyServiceCreate.Close()

	if txStmtRepositoryPropertySystemCreate, err = tx.Prepare(tkStmtRepositoryPropertySystemCreate); err != nil {
		log.Println("Failed to prepare: tkStmtRepositoryPropertySystemCreate")
		goto bailout
	}
	defer txStmtRepositoryPropertySystemCreate.Close()

	if txStmtRepositoryPropertyCustomCreate, err = tx.Prepare(tkStmtRepositoryPropertyCustomCreate); err != nil {
		log.Println("Failed to prepare: tkStmtRepositoryPropertyCustomCreate")
		goto bailout
	}
	defer txStmtRepositoryPropertyCustomCreate.Close()

	//
	// BUCKET
	if txStmtCreateBucket, err = tx.Prepare(tkStmtCreateBucket); err != nil {
		log.Println("Failed to prepare: tkStmtCreateBucket")
		goto bailout
	}
	defer txStmtCreateBucket.Close()

	if txStmtBucketPropertyOncallCreate, err = tx.Prepare(tkStmtBucketPropertyOncallCreate); err != nil {
		log.Println("Failed to prepare: tkStmtBucketPropertyOncallCreate")
		goto bailout
	}
	defer txStmtBucketPropertyOncallCreate.Close()

	if txStmtBucketPropertyServiceCreate, err = tx.Prepare(tkStmtBucketPropertyServiceCreate); err != nil {
		log.Println("Failed to prepare: tkStmtBucketPropertyServiceCreate")
		goto bailout
	}
	defer txStmtBucketPropertyServiceCreate.Close()

	if txStmtBucketPropertySystemCreate, err = tx.Prepare(tkStmtBucketPropertySystemCreate); err != nil {
		log.Println("Failed to prepare: tkStmtBucketPropertySystemCreate")
		goto bailout
	}
	defer txStmtBucketPropertySystemCreate.Close()

	if txStmtBucketPropertyCustomCreate, err = tx.Prepare(tkStmtBucketPropertyCustomCreate); err != nil {
		log.Println("Failed to prepare: tkStmtBucketPropertyCustomCreate")
		goto bailout
	}
	defer txStmtBucketPropertyCustomCreate.Close()

	//
	// GROUP
	if txStmtGroupCreate, err = tx.Prepare(tkStmtGroupCreate); err != nil {
		log.Println("Failed to prepare: tkStmtGroupCreate")
		goto bailout
	}
	defer txStmtGroupCreate.Close()

	if txStmtGroupUpdate, err = tx.Prepare(tkStmtGroupUpdate); err != nil {
		log.Println("Failed to prepare: tkStmtGroupUpdate")
		goto bailout
	}
	defer txStmtGroupUpdate.Close()

	if txStmtGroupDelete, err = tx.Prepare(tkStmtGroupDelete); err != nil {
		log.Println("Failed to prepare: tkStmtGroupDelete")
		goto bailout
	}
	defer txStmtGroupDelete.Close()

	if txStmtGroupMemberNewNode, err = tx.Prepare(tkStmtGroupMemberNewNode); err != nil {
		log.Println("Failed to prepare: tkStmtGroupMemberNewNode")
		goto bailout
	}
	defer txStmtGroupMemberNewNode.Close()

	if txStmtGroupMemberNewCluster, err = tx.Prepare(tkStmtGroupMemberNewCluster); err != nil {
		log.Println("Failed to prepare: tkStmtGroupMemberNewCluster")
		goto bailout
	}
	defer txStmtGroupMemberNewCluster.Close()

	if txStmtGroupMemberNewGroup, err = tx.Prepare(tkStmtGroupMemberNewGroup); err != nil {
		log.Println("Failed to prepare: tkStmtGroupMemberNewGroup")
		goto bailout
	}
	defer txStmtGroupMemberNewGroup.Close()

	if txStmtGroupMemberRemoveNode, err = tx.Prepare(tkStmtGroupMemberRemoveNode); err != nil {
		log.Println("Failed to prepare: tkStmtGroupMemberRemoveNode")
		goto bailout
	}
	defer txStmtGroupMemberRemoveNode.Close()

	if txStmtGroupMemberRemoveCluster, err = tx.Prepare(tkStmtGroupMemberRemoveCluster); err != nil {
		log.Println("Failed to prepare: tkStmtGroupMemberRemoveCluster")
		goto bailout
	}
	defer txStmtGroupMemberRemoveCluster.Close()

	if txStmtGroupMemberRemoveGroup, err = tx.Prepare(tkStmtGroupMemberRemoveGroup); err != nil {
		log.Println("Failed to prepare: tkStmtGroupMemberRemoveGroup")
		goto bailout
	}
	defer txStmtGroupMemberRemoveGroup.Close()

	if txStmtGroupPropertyOncallCreate, err = tx.Prepare(tkStmtGroupPropertyOncallCreate); err != nil {
		log.Println("Failed to prepare: tkStmtGroupPropertyOncallCreate")
		goto bailout
	}
	defer txStmtGroupPropertyOncallCreate.Close()

	if txStmtGroupPropertyServiceCreate, err = tx.Prepare(tkStmtGroupPropertyServiceCreate); err != nil {
		log.Println("Failed to prepare: tkStmtGroupPropertyServiceCreate")
		goto bailout
	}
	defer txStmtGroupPropertyServiceCreate.Close()

	if txStmtGroupPropertySystemCreate, err = tx.Prepare(tkStmtGroupPropertySystemCreate); err != nil {
		log.Println("Failed to prepare: tkStmtGroupPropertySystemCreate")
		goto bailout
	}
	defer txStmtGroupPropertySystemCreate.Close()

	if txStmtGroupPropertyCustomCreate, err = tx.Prepare(tkStmtGroupPropertyCustomCreate); err != nil {
		log.Println("Failed to prepare: tkStmtGroupPropertyCustomCreate")
		goto bailout
	}
	defer txStmtGroupPropertyCustomCreate.Close()

	//
	// CLUSTER
	if txStmtClusterCreate, err = tx.Prepare(tkStmtClusterCreate); err != nil {
		log.Println("Failed to prepare: tkStmtClusterCreate")
		goto bailout
	}
	defer txStmtClusterCreate.Close()

	if txStmtClusterUpdate, err = tx.Prepare(tkStmtClusterUpdate); err != nil {
		log.Println("Failed to prepare: tkStmtClusterUpdate")
		goto bailout
	}
	defer txStmtClusterUpdate.Close()

	if txStmtClusterDelete, err = tx.Prepare(tkStmtClusterDelete); err != nil {
		log.Println("Failed to prepare: tkStmtClusterDelete")
		goto bailout
	}
	defer txStmtClusterDelete.Close()

	if txStmtClusterMemberNew, err = tx.Prepare(tkStmtClusterMemberNew); err != nil {
		log.Println("Failed to prepare: tkStmtClusterMemberNew")
		goto bailout
	}
	defer txStmtClusterMemberNew.Close()

	if txStmtClusterMemberRemove, err = tx.Prepare(tkStmtClusterMemberRemove); err != nil {
		log.Println("Failed to prepare: tkStmtClusterMemberRemove")
		goto bailout
	}
	defer txStmtClusterMemberRemove.Close()

	if txStmtClusterPropertyOncallCreate, err = tx.Prepare(tkStmtClusterPropertyOncallCreate); err != nil {
		log.Println("Failed to prepare: tkStmtClusterPropertyOncallCreate")
		goto bailout
	}
	defer txStmtClusterPropertyOncallCreate.Close()

	if txStmtClusterPropertyServiceCreate, err = tx.Prepare(tkStmtClusterPropertyServiceCreate); err != nil {
		log.Println("Failed to prepare: tkStmtClusterPropertyServiceCreate")
		goto bailout
	}
	defer txStmtClusterPropertyServiceCreate.Close()

	if txStmtClusterPropertySystemCreate, err = tx.Prepare(tkStmtClusterPropertySystemCreate); err != nil {
		log.Println("Failed to prepare: tkStmtClusterPropertySystemCreate")
		goto bailout
	}
	defer txStmtClusterPropertySystemCreate.Close()

	if txStmtClusterPropertyCustomCreate, err = tx.Prepare(tkStmtClusterPropertyCustomCreate); err != nil {
		log.Println("Failed to prepare: tkStmtClusterPropertyCustomCreate")
		goto bailout
	}
	defer txStmtClusterPropertyCustomCreate.Close()

	//
	// NODE
	if txStmtBucketAssignNode, err = tx.Prepare(tkStmtBucketAssignNode); err != nil {
		log.Println("Failed to prepare: tkStmtBucketAssignNode")
		goto bailout
	}
	defer txStmtBucketAssignNode.Close()

	if txStmtUpdateNodeState, err = tx.Prepare(tkStmtUpdateNodeState); err != nil {
		log.Println("Failed to prepare: tkStmtUpdateNodeState")
		goto bailout
	}
	defer txStmtUpdateNodeState.Close()

	if txStmtNodeUnassignFromBucket, err = tx.Prepare(tkStmtNodeUnassignFromBucket); err != nil {
		log.Println("Failed to prepare: tkStmtNodeUnassignFromBucket")
		goto bailout
	}
	defer txStmtNodeUnassignFromBucket.Close()

	if txStmtNodePropertyOncallCreate, err = tx.Prepare(tkStmtNodePropertyOncallCreate); err != nil {
		log.Println("Failed to prepare: tkStmtNodePropertyOncallCreate")
		goto bailout
	}
	defer txStmtNodePropertyOncallCreate.Close()

	if txStmtNodePropertyServiceCreate, err = tx.Prepare(tkStmtNodePropertyServiceCreate); err != nil {
		log.Println("Failed to prepare: tkStmtNodePropertyServiceCreate")
		goto bailout
	}
	defer txStmtNodePropertyServiceCreate.Close()

	if txStmtNodePropertySystemCreate, err = tx.Prepare(tkStmtNodePropertySystemCreate); err != nil {
		log.Println("Failed to prepare: tkStmtNodePropertySystemCreate")
		goto bailout
	}
	defer txStmtNodePropertySystemCreate.Close()

	if txStmtNodePropertyCustomCreate, err = tx.Prepare(tkStmtNodePropertyCustomCreate); err != nil {
		log.Println("Failed to prepare: tkStmtNodePropertyCustomCreate")
		goto bailout
	}
	defer txStmtNodePropertyCustomCreate.Close()

	// defer constraint checks
	if _, err = tx.Exec(tkStmtDeferAllConstraints); err != nil {
		log.Println("Failed to exec: tkStmtDeferAllConstraints")
		goto bailout
	}

	// save the check configuration as part of the transaction before
	// processing the action channel
	if strings.Contains(q.Action, "add_check_to_") {
		if q.CheckConfig.CheckConfig.BucketId != "" {
			nullBucket = sql.NullString{
				String: q.CheckConfig.CheckConfig.BucketId,
				Valid:  true,
			}
		} else {
			nullBucket = sql.NullString{String: "", Valid: false}
		}

		if _, err = txStmtCreateCheckConfigurationBase.Exec(
			q.CheckConfig.CheckConfig.Id,
			q.CheckConfig.CheckConfig.Name,
			int64(q.CheckConfig.CheckConfig.Interval),
			q.CheckConfig.CheckConfig.RepositoryId,
			nullBucket,
			q.CheckConfig.CheckConfig.CapabilityId,
			q.CheckConfig.CheckConfig.ObjectId,
			q.CheckConfig.CheckConfig.ObjectType,
			q.CheckConfig.CheckConfig.IsActive,
			q.CheckConfig.CheckConfig.IsEnabled,
			q.CheckConfig.CheckConfig.Inheritance,
			q.CheckConfig.CheckConfig.ChildrenOnly,
			q.CheckConfig.CheckConfig.ExternalId,
		); err != nil {
			goto bailout
		}

	threshloop:
		for _, thr := range q.CheckConfig.CheckConfig.Thresholds {
			if _, err = txStmtCreateCheckConfigurationThreshold.Exec(
				q.CheckConfig.CheckConfig.Id,
				thr.Predicate.Symbol,
				strconv.FormatInt(thr.Value, 10),
				thr.Level.Name,
			); err != nil {
				break threshloop
			}
		}
		if err != nil {
			goto bailout
		}

	constrloop:
		for _, constr := range q.CheckConfig.CheckConfig.Constraints {
			switch constr.ConstraintType {
			case "native":
				if _, err = txStmtCreateCheckConfigurationConstraintNative.Exec(
					q.CheckConfig.CheckConfig.Id,
					constr.Native.Name,
					constr.Native.Value,
				); err != nil {
					break constrloop
				}
			case "oncall":
				if _, err = txStmtCreateCheckConfigurationConstraintOncall.Exec(
					q.CheckConfig.CheckConfig.Id,
					constr.Oncall.Id,
				); err != nil {
					break constrloop
				}
			case "custom":
				if _, err = txStmtCreateCheckConfigurationConstraintCustom.Exec(
					q.CheckConfig.CheckConfig.Id,
					constr.Custom.Id,
					constr.Custom.RepositoryId,
					constr.Custom.Value,
				); err != nil {
					break constrloop
				}
			case "system":
				if _, err = txStmtCreateCheckConfigurationConstraintSystem.Exec(
					q.CheckConfig.CheckConfig.Id,
					constr.System.Name,
					constr.System.Value,
				); err != nil {
					break constrloop
				}
			case "service":
				if constr.Service.TeamId != tk.team {
					err = fmt.Errorf("Service constraint has mismatched TeamID values: %s/%s",
						tk.team, constr.Service.TeamId)
					fmt.Println(err)
					break constrloop
				}
				log.Printf(`SQL: tkStmtCreateCheckConfigurationConstraintService:
CheckConfig ID: %s
Team ID:        %s
Service Name:   %s%s`,
					q.CheckConfig.CheckConfig.Id,
					tk.team,
					constr.Service.Name, "\n")
				if _, err = txStmtCreateCheckConfigurationConstraintService.Exec(
					q.CheckConfig.CheckConfig.Id,
					tk.team,
					constr.Service.Name,
				); err != nil {
					break constrloop
				}
			case "attribute":
				if _, err = txStmtCreateCheckConfigurationConstraintAttribute.Exec(
					q.CheckConfig.CheckConfig.Id,
					constr.Attribute.Name,
					constr.Attribute.Value,
				); err != nil {
					break constrloop
				}
			}
		}
		if err != nil {
			goto bailout
		}
	}

	// mark the check configuration as deleted
	if strings.HasPrefix(q.Action, `remove_check_from_`) {
		if _, err = tx.Exec(stmt.TxMarkCheckConfigDeleted, q.CheckConfig.CheckConfig.Id); err != nil {
			goto bailout
		}
	}

	// if the error channel has entries, we can fully ignore the
	// action channel
	for i := len(tk.errChan); i > 0; i-- {
		e := <-tk.errChan
		b, _ := json.Marshal(e)
		log.Println(string(b))
		hasErrors = true
		if err == nil {
			err = fmt.Errorf(e.Action)
		}
	}
	if hasErrors {
		goto bailout
	}

actionloop:
	for i := len(tk.actionChan); i > 0; i-- {
		a := <-tk.actionChan

		// we need all messages to figure out why for example a deferred
		// constraint later failed
		//jBxX, _ := json.Marshal(a)
		//log.Printf("%s - Processing: %s\n", q.JobId.String(), string(jBxX))

		switch a.Type {
		// REPOSITORY
		case "repository":
			switch a.Action {
			case "property_new":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtPropertyInstanceCreate.Exec(
					a.Property.InstanceId,
					a.Property.RepositoryId,
					a.Property.SourceInstanceId,
					a.Property.SourceType,
					a.Property.InheritedFrom,
				); err != nil {
					break actionloop
				}
				switch a.Property.Type {
				case "custom":
					if _, err = txStmtRepositoryPropertyCustomCreate.Exec(
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Property.Custom.RepositoryId,
						a.Property.View,
						a.Property.Custom.Id,
						a.Property.Inheritance,
						a.Property.ChildrenOnly,
						a.Property.Custom.Value,
					); err != nil {
						break actionloop
					}
				case "system":
					if _, err = txStmtRepositoryPropertySystemCreate.Exec(
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
					); err != nil {
						break actionloop
					}
				case "service":
					if _, err = txStmtRepositoryPropertyServiceCreate.Exec(
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Repository.Id,
						a.Property.View,
						a.Property.Service.Name,
						a.Property.Service.TeamId,
						a.Property.Inheritance,
						a.Property.ChildrenOnly,
					); err != nil {
						break actionloop
					}
				case "oncall":
					if _, err = txStmtRepositoryPropertyOncallCreate.Exec(
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Repository.Id,
						a.Property.View,
						a.Property.Oncall.Id,
						a.Property.Inheritance,
						a.Property.ChildrenOnly,
					); err != nil {
						break actionloop
					}
				}
			case `property_delete`:
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = tx.Exec(tkStmtPropertyInstanceDelete,
					a.Property.InstanceId,
				); err != nil {
					break actionloop
				}
				switch a.Property.Type {
				case `custom`:
					if _, err = tx.Exec(tkStmtRepositoryPropertyCustomDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `system`:
					if _, err = tx.Exec(tkStmtRepositoryPropertySystemDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `service`:
					if _, err = tx.Exec(tkStmtRepositoryPropertyServiceDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `oncall`:
					if _, err = tx.Exec(tkStmtRepositoryPropertyOncallDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				}
			case "check_new":
				if tk.rebuild && tk.rbLevel == `instances` {
					// ignore in instance-rebuild mode
					continue actionloop
				}
				if _, err = txStmtCreateCheck.Exec(
					a.Check.CheckId,
					a.Check.RepositoryId,
					sql.NullString{String: "", Valid: false},
					a.Check.SourceCheckId,
					a.Check.SourceType,
					a.Check.InheritedFrom,
					a.Check.CheckConfigId,
					a.Check.CapabilityId,
					a.Repository.Id,
					"repository",
				); err != nil {
					break actionloop
				}
			case `check_removed`:
				if tk.rebuild && tk.rbLevel == `instances` {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtDeleteCheck.Exec(
					a.Check.CheckId,
				); err != nil {
					break actionloop
				}
			default:
				jB, _ := json.Marshal(a)
				log.Printf("Unhandled message: %s\n", string(jB))
			}

		// BUCKET
		case "bucket":
			switch a.Action {
			case "create":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtCreateBucket.Exec(
					a.Bucket.Id,
					a.Bucket.Name,
					a.Bucket.IsFrozen,
					a.Bucket.IsDeleted,
					a.Bucket.RepositoryId,
					a.Bucket.Environment,
					a.Bucket.TeamId,
					q.User,
				); err != nil {
					break actionloop
				}
			case "node_assignment":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtBucketAssignNode.Exec(
					a.ChildNode.Id,
					a.Bucket.Id,
					a.Bucket.TeamId,
				); err != nil {
					break actionloop
				}
			case "property_new":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtPropertyInstanceCreate.Exec(
					a.Property.InstanceId,
					a.Property.RepositoryId,
					a.Property.SourceInstanceId,
					a.Property.SourceType,
					a.Property.InheritedFrom,
				); err != nil {
					break actionloop
				}
				switch a.Property.Type {
				case "custom":
					if _, err = txStmtBucketPropertyCustomCreate.Exec(
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Bucket.Id,
						a.Property.View,
						a.Property.Custom.Id,
						a.Property.Custom.RepositoryId,
						a.Property.Inheritance,
						a.Property.ChildrenOnly,
						a.Property.Custom.Value,
					); err != nil {
						break actionloop
					}
				case "system":
					if _, err = txStmtBucketPropertySystemCreate.Exec(
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
					); err != nil {
						break actionloop
					}
				case "service":
					log.Printf(`SQL: tkStmtBucketPropertyServiceCreate:
Instance ID:           %s
Source Instance ID:    %s
Bucket ID:             %s
View:                  %s
Service Name:          %s
Service TeamId:        %s
Repository ID:         %s
Inheritance Enabled:   %t
Children Only:         %t%s`,
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Bucket.Id,
						a.Property.View,
						a.Property.Service.Name,
						a.Property.Service.TeamId,
						a.Property.RepositoryId,
						a.Property.Inheritance,
						a.Property.ChildrenOnly, "\n")
					if _, err = txStmtBucketPropertyServiceCreate.Exec(
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Bucket.Id,
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
					if _, err = txStmtBucketPropertyOncallCreate.Exec(
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Bucket.Id,
						a.Property.View,
						a.Property.Oncall.Id,
						a.Property.RepositoryId,
						a.Property.Inheritance,
						a.Property.ChildrenOnly,
					); err != nil {
						break actionloop
					}
				}
			case `property_delete`:
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = tx.Exec(tkStmtPropertyInstanceDelete,
					a.Property.InstanceId,
				); err != nil {
					break actionloop
				}
				switch a.Property.Type {
				case `custom`:
					if _, err = tx.Exec(tkStmtBucketPropertyCustomDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `system`:
					if _, err = tx.Exec(tkStmtBucketPropertySystemDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `service`:
					if _, err = tx.Exec(tkStmtBucketPropertyServiceDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `oncall`:
					if _, err = tx.Exec(tkStmtBucketPropertyOncallDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				}
			case "check_new":
				if tk.rebuild && tk.rbLevel == `instances` {
					// ignore in instance rebuild mode
					continue actionloop
				}
				if _, err = txStmtCreateCheck.Exec(
					a.Check.CheckId,
					a.Check.RepositoryId,
					a.Check.BucketId,
					a.Check.SourceCheckId,
					a.Check.SourceType,
					a.Check.InheritedFrom,
					a.Check.CheckConfigId,
					a.Check.CapabilityId,
					a.Bucket.Id,
					"bucket",
				); err != nil {
					break actionloop
				}
			case `check_removed`:
				if tk.rebuild && tk.rbLevel == `instances` {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtDeleteCheck.Exec(
					a.Check.CheckId,
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
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtGroupCreate.Exec(
					a.Group.Id,
					a.Group.BucketId,
					a.Group.Name,
					a.Group.ObjectState,
					a.Group.TeamId,
					q.User,
				); err != nil {
					break actionloop
				}
			case "update":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtGroupUpdate.Exec(
					a.Group.Id,
					a.Group.ObjectState,
				); err != nil {
					break actionloop
				}
			case "delete":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtGroupDelete.Exec(
					a.Group.Id,
				); err != nil {
					break actionloop
				}
			case "member_new":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				switch a.ChildType {
				case "group":
					log.Println("==> group/new membergroup")
					if _, err = txStmtGroupMemberNewGroup.Exec(
						a.Group.Id,
						a.ChildGroup.Id,
						a.Group.BucketId,
					); err != nil {
						break actionloop
					}
				case "cluster":
					log.Println("==> group/new membercluster")
					if _, err = txStmtGroupMemberNewCluster.Exec(
						a.Group.Id,
						a.ChildCluster.Id,
						a.Group.BucketId,
					); err != nil {
						break actionloop
					}
				case "node":
					log.Println("==> group/new membernode")
					if _, err = txStmtGroupMemberNewNode.Exec(
						a.Group.Id,
						a.ChildNode.Id,
						a.Group.BucketId,
					); err != nil {
						break actionloop
					}
				}
			case "member_removed":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
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
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtPropertyInstanceCreate.Exec(
					a.Property.InstanceId,
					a.Property.RepositoryId,
					a.Property.SourceInstanceId,
					a.Property.SourceType,
					a.Property.InheritedFrom,
				); err != nil {
					break actionloop
				}
				switch a.Property.Type {
				case "custom":
					if _, err = txStmtGroupPropertyCustomCreate.Exec(
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
						a.Property.Oncall.Id,
						a.Property.RepositoryId,
						a.Property.Inheritance,
						a.Property.ChildrenOnly,
					); err != nil {
						break actionloop
					}
				}
			case `property_delete`:
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = tx.Exec(tkStmtPropertyInstanceDelete,
					a.Property.InstanceId,
				); err != nil {
					break actionloop
				}
				switch a.Property.Type {
				case `custom`:
					if _, err = tx.Exec(tkStmtGroupPropertyCustomDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `system`:
					if _, err = tx.Exec(tkStmtGroupPropertySystemDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `service`:
					if _, err = tx.Exec(tkStmtGroupPropertyServiceDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `oncall`:
					if _, err = tx.Exec(tkStmtGroupPropertyOncallDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				}
			case "check_new":
				if tk.rebuild && tk.rbLevel == `instances` {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtCreateCheck.Exec(
					a.Check.CheckId,
					a.Check.RepositoryId,
					a.Check.BucketId,
					a.Check.SourceCheckId,
					a.Check.SourceType,
					a.Check.InheritedFrom,
					a.Check.CheckConfigId,
					a.Check.CapabilityId,
					a.Group.Id,
					"group",
				); err != nil {
					break actionloop
				}
			case `check_removed`:
				if tk.rebuild && tk.rbLevel == `instances` {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtDeleteCheck.Exec(
					a.Check.CheckId,
				); err != nil {
					break actionloop
				}
			case "check_instance_create":
				if _, err = txStmtCreateCheckInstance.Exec(
					a.CheckInstance.InstanceId,
					a.CheckInstance.CheckId,
					a.CheckInstance.ConfigId,
					"00000000-0000-0000-0000-000000000000",
					time.Now().UTC(),
				); err != nil {
					break actionloop
				}
				fallthrough
			case "check_instance_update":
				if _, err = txStmtCreateCheckInstanceConfiguration.Exec(
					a.CheckInstance.InstanceConfigId,
					a.CheckInstance.Version,
					a.CheckInstance.InstanceId,
					a.CheckInstance.ConstraintHash,
					a.CheckInstance.ConstraintValHash,
					a.CheckInstance.InstanceService,
					a.CheckInstance.InstanceSvcCfgHash,
					a.CheckInstance.InstanceServiceConfig,
					time.Now().UTC(),
					"awaiting_computation",
					"none",
					false,
					"{}",
				); err != nil {
					fmt.Println(`Failed CreateCheckInstanceConfiguration`, a.CheckInstance.InstanceConfigId)
					break actionloop
				}
			case "check_instance_delete":
				if _, err = txStmtDeleteCheckInstance.Exec(
					a.CheckInstance.InstanceId,
				); err != nil {
					break actionloop
				}
			default:
				jB, _ := json.Marshal(a)
				log.Printf("Unhandled message: %s\n", string(jB))
			}
		// CLUSTER
		case "cluster":
			switch a.Action {
			case "create":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtClusterCreate.Exec(
					a.Cluster.Id,
					a.Cluster.Name,
					a.Cluster.BucketId,
					a.Cluster.ObjectState,
					a.Cluster.TeamId,
					q.User,
				); err != nil {
					break actionloop
				}
			case "update":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtClusterUpdate.Exec(
					a.Cluster.Id,
					a.Cluster.ObjectState,
				); err != nil {
					break actionloop
				}
			case "delete":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtClusterDelete.Exec(
					a.Cluster.Id,
				); err != nil {
					break actionloop
				}
			case "member_new":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				log.Println("==> cluster/new membernode")
				if _, err = txStmtClusterMemberNew.Exec(
					a.Cluster.Id,
					a.ChildNode.Id,
					a.Cluster.BucketId,
				); err != nil {
					break actionloop
				}
			case "member_removed":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				log.Println("==> cluster/new membernode")
				if _, err = txStmtClusterMemberRemove.Exec(
					a.Cluster.Id,
					a.ChildNode.Id,
				); err != nil {
					break actionloop
				}
			case "property_new":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtPropertyInstanceCreate.Exec(
					a.Property.InstanceId,
					a.Property.RepositoryId,
					a.Property.SourceInstanceId,
					a.Property.SourceType,
					a.Property.InheritedFrom,
				); err != nil {
					break actionloop
				}
				switch a.Property.Type {
				case "custom":
					if _, err = txStmtClusterPropertyCustomCreate.Exec(
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
					); err != nil {
						break actionloop
					}
				case "system":
					if _, err = txStmtClusterPropertySystemCreate.Exec(
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
					); err != nil {
						break actionloop
					}
				case "service":
					log.Printf(`SQL: tkStmtClusterPropertyServiceCreate:
Instance ID:           %s
Source Instance ID:    %s
Cluster ID:            %s
View:                  %s
Service Name:          %s
Service TeamId:        %s
Repository ID:         %s
Inheritance Enabled:   %t
Children Only:         %t%s`,
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Cluster.Id,
						a.Property.View,
						a.Property.Service.Name,
						a.Property.Service.TeamId,
						a.Property.RepositoryId,
						a.Property.Inheritance,
						a.Property.ChildrenOnly, "\n")
					if _, err = txStmtClusterPropertyServiceCreate.Exec(
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Cluster.Id,
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
					if _, err = txStmtClusterPropertyOncallCreate.Exec(
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Cluster.Id,
						a.Property.View,
						a.Property.Oncall.Id,
						a.Property.RepositoryId,
						a.Property.Inheritance,
						a.Property.ChildrenOnly,
					); err != nil {
						break actionloop
					}
				}
			case `property_delete`:
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = tx.Exec(tkStmtPropertyInstanceDelete,
					a.Property.InstanceId,
				); err != nil {
					break actionloop
				}
				switch a.Property.Type {
				case `custom`:
					if _, err = tx.Exec(tkStmtClusterPropertyCustomDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `system`:
					if _, err = tx.Exec(tkStmtClusterPropertySystemDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `service`:
					if _, err = tx.Exec(tkStmtClusterPropertyServiceDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `oncall`:
					if _, err = tx.Exec(tkStmtClusterPropertyOncallDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				}
			case "check_new":
				if tk.rebuild && tk.rbLevel == `instances` {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtCreateCheck.Exec(
					a.Check.CheckId,
					a.Check.RepositoryId,
					a.Check.BucketId,
					a.Check.SourceCheckId,
					a.Check.SourceType,
					a.Check.InheritedFrom,
					a.Check.CheckConfigId,
					a.Check.CapabilityId,
					a.Cluster.Id,
					"cluster",
				); err != nil {
					break actionloop
				}
			case `check_removed`:
				if tk.rebuild && tk.rbLevel == `instances` {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtDeleteCheck.Exec(
					a.Check.CheckId,
				); err != nil {
					break actionloop
				}
			case "check_instance_create":
				if _, err = txStmtCreateCheckInstance.Exec(
					a.CheckInstance.InstanceId,
					a.CheckInstance.CheckId,
					a.CheckInstance.ConfigId,
					"00000000-0000-0000-0000-000000000000",
					time.Now().UTC(),
				); err != nil {
					break actionloop
				}
				fallthrough
			case "check_instance_update":
				if _, err = txStmtCreateCheckInstanceConfiguration.Exec(
					a.CheckInstance.InstanceConfigId,
					a.CheckInstance.Version,
					a.CheckInstance.InstanceId,
					a.CheckInstance.ConstraintHash,
					a.CheckInstance.ConstraintValHash,
					a.CheckInstance.InstanceService,
					a.CheckInstance.InstanceSvcCfgHash,
					a.CheckInstance.InstanceServiceConfig,
					time.Now().UTC(),
					"awaiting_computation",
					"none",
					false,
					"{}",
				); err != nil {
					fmt.Println(`Failed CreateCheckInstanceConfiguration`, a.CheckInstance.InstanceConfigId)
					break actionloop
				}
			case "check_instance_delete":
				if _, err = txStmtDeleteCheckInstance.Exec(
					a.CheckInstance.InstanceId,
				); err != nil {
					break actionloop
				}
			default:
				jB, _ := json.Marshal(a)
				log.Printf("Unhandled message: %s\n", string(jB))
			}
		// NODE
		case "node":
			switch a.Action {
			case "delete":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtNodeUnassignFromBucket.Exec(
					a.Node.Id,
					a.Node.Config.BucketId,
					a.Node.TeamId,
				); err != nil {
					break actionloop
				}
				fallthrough // need to call txStmtUpdateNodeState for delete as well
			case "update":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				log.Println("==> node/update")
				if _, err = txStmtUpdateNodeState.Exec(
					a.Node.Id,
					a.Node.State,
				); err != nil {
					break actionloop
				}
			case "property_new":
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = txStmtPropertyInstanceCreate.Exec(
					a.Property.InstanceId,
					a.Property.RepositoryId,
					a.Property.SourceInstanceId,
					a.Property.SourceType,
					a.Property.InheritedFrom,
				); err != nil {
					break actionloop
				}
				switch a.Property.Type {
				case "custom":
					if _, err = txStmtNodePropertyCustomCreate.Exec(
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
					); err != nil {
						break actionloop
					}
				case "system":
					log.Printf(`SQL: tkStmtNodePropertySystemCreate:
Instance ID:           %s
Source Instance ID:    %s
Node ID:               %s
View:                  %s
SystemProperty:        %s
Object Type:           %s
Repository ID:         %s
Inheritance Enabled:   %t
Children Only:         %t
System Property Value: %s
Is Inherited:          %t%s`,
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
						a.Property.IsInherited, "\n")
					if _, err = txStmtNodePropertySystemCreate.Exec(
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
					); err != nil {
						break actionloop
					}
				case "service":
					if _, err = txStmtNodePropertyServiceCreate.Exec(
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Node.Id,
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
					if _, err = txStmtNodePropertyOncallCreate.Exec(
						a.Property.InstanceId,
						a.Property.SourceInstanceId,
						a.Node.Id,
						a.Property.View,
						a.Property.Oncall.Id,
						a.Property.RepositoryId,
						a.Property.Inheritance,
						a.Property.ChildrenOnly,
					); err != nil {
						break actionloop
					}
				}
			case `property_delete`:
				if tk.rebuild {
					// ignore in rebuild mode
					continue actionloop
				}
				if _, err = tx.Exec(tkStmtPropertyInstanceDelete,
					a.Property.InstanceId,
				); err != nil {
					break actionloop
				}
				switch a.Property.Type {
				case `custom`:
					if _, err = tx.Exec(tkStmtNodePropertyCustomDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `system`:
					if _, err = tx.Exec(tkStmtNodePropertySystemDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `service`:
					if _, err = tx.Exec(tkStmtNodePropertyServiceDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				case `oncall`:
					if _, err = tx.Exec(tkStmtNodePropertyOncallDelete,
						a.Property.InstanceId,
					); err != nil {
						break actionloop
					}
				}
			case "check_new":
				if tk.rebuild && tk.rbLevel == `instances` {
					// ignore in instance-rebuild mode
					continue actionloop
				}
				log.Printf(`SQL: tkStmtCreateCheck:
Check ID:            %s
Repository ID:       %s
Bucket ID:           %s
Source Check ID:     %s
Source Type:         %s
Inherited From:      %s
Check Config ID:     %s
Check Capability ID: %s
Node ID:             %s%s`,
					a.Check.CheckId,
					a.Check.RepositoryId,
					a.Check.BucketId,
					a.Check.SourceCheckId,
					a.Check.SourceType,
					a.Check.InheritedFrom,
					a.Check.CheckConfigId,
					a.Check.CapabilityId,
					a.Node.Id, "\n")
				if _, err = txStmtCreateCheck.Exec(
					a.Check.CheckId,
					a.Check.RepositoryId,
					a.Check.BucketId,
					a.Check.SourceCheckId,
					a.Check.SourceType,
					a.Check.InheritedFrom,
					a.Check.CheckConfigId,
					a.Check.CapabilityId,
					a.Node.Id,
					"node",
				); err != nil {
					break actionloop
				}
			case `check_removed`:
				if tk.rebuild && tk.rbLevel == `instances` {
					// ignore in instance-rebuild mode
					continue actionloop
				}
				if _, err = tx.Exec(stmt.TxMarkCheckDeleted,
					a.Check.CheckId,
				); err != nil {
					break actionloop
				}
			case "check_instance_create":
				if _, err = txStmtCreateCheckInstance.Exec(
					a.CheckInstance.InstanceId,
					a.CheckInstance.CheckId,
					a.CheckInstance.ConfigId,
					"00000000-0000-0000-0000-000000000000",
					time.Now().UTC(),
				); err != nil {
					break actionloop
				}
				fallthrough
			case "check_instance_update":
				if _, err = txStmtCreateCheckInstanceConfiguration.Exec(
					a.CheckInstance.InstanceConfigId,
					a.CheckInstance.Version,
					a.CheckInstance.InstanceId,
					a.CheckInstance.ConstraintHash,
					a.CheckInstance.ConstraintValHash,
					a.CheckInstance.InstanceService,
					a.CheckInstance.InstanceSvcCfgHash,
					a.CheckInstance.InstanceServiceConfig,
					time.Now().UTC(),
					"awaiting_computation",
					"none",
					false,
					"{}",
				); err != nil {
					fmt.Println(`Failed CreateCheckInstanceConfiguration`, a.CheckInstance.InstanceConfigId)
					break actionloop
				}
			case "check_instance_delete":
				if _, err = txStmtDeleteCheckInstance.Exec(
					a.CheckInstance.InstanceId,
				); err != nil {
					break actionloop
				}
			default:
				jB, _ := json.Marshal(a)
				log.Printf("Unhandled message: %s\n", string(jB))
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

	if !tk.rebuild {
		// mark job as finished
		if _, err = tx.Exec(
			tkStmtFinishJob,
			q.JobId.String(),
			time.Now().UTC(),
			"success",
			``, // empty error field
		); err != nil {
			goto bailout
		}
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

	// if this was a rebuild, the tree will not persist and the
	// job is faked
	if tk.rebuild {
		return
	}

	tk.tree.Rollback()
	tx.Rollback()
	tk.conn.Exec(
		tkStmtFinishJob,
		q.JobId.String(),
		time.Now().UTC(),
		"failed",
		err.Error(),
	)
	for i := len(tk.actionChan); i > 0; i-- {
		a := <-tk.actionChan
		jB, _ := json.Marshal(a)
		log.Printf("Cleaned message: %s\n", string(jB))
	}
	return
}

func (tk *treeKeeper) convertCheckForDelete(conf *proto.CheckConfig) (*tree.Check, error) {
	var err error
	treechk := &tree.Check{
		Id:            uuid.Nil,
		InheritedFrom: uuid.Nil,
	}
	if treechk.SourceId, err = uuid.FromString(conf.ExternalId); err != nil {
		return nil, err
	}
	if treechk.ConfigId, err = uuid.FromString(conf.Id); err != nil {
		return nil, err
	}
	return treechk, nil
}

func (tk *treeKeeper) convertCheck(conf *proto.CheckConfig) (*tree.Check, error) {
	treechk := &tree.Check{
		Id:            uuid.Nil,
		SourceId:      uuid.Nil,
		InheritedFrom: uuid.Nil,
		Inheritance:   conf.Inheritance,
		ChildrenOnly:  conf.ChildrenOnly,
		Interval:      conf.Interval,
	}
	treechk.CapabilityId, _ = uuid.FromString(conf.CapabilityId)
	treechk.ConfigId, _ = uuid.FromString(conf.Id)
	if err := tk.get_view.QueryRow(conf.CapabilityId).Scan(&treechk.View); err != nil {
		return &tree.Check{}, err
	}

	treechk.Thresholds = make([]tree.CheckThreshold, len(conf.Thresholds))
	for i, thr := range conf.Thresholds {
		nthr := tree.CheckThreshold{
			Predicate: thr.Predicate.Symbol,
			Level:     uint8(thr.Level.Numeric),
			Value:     thr.Value,
		}
		treechk.Thresholds[i] = nthr
	}

	treechk.Constraints = make([]tree.CheckConstraint, len(conf.Constraints))
	for i, constr := range conf.Constraints {
		ncon := tree.CheckConstraint{
			Type: constr.ConstraintType,
		}
		switch constr.ConstraintType {
		case "native":
			ncon.Key = constr.Native.Name
			ncon.Value = constr.Native.Value
		case "oncall":
			ncon.Key = "OncallId"
			ncon.Value = constr.Oncall.Id
		case "custom":
			ncon.Key = constr.Custom.Id
			ncon.Value = constr.Custom.Value
		case "system":
			ncon.Key = constr.System.Name
			ncon.Value = constr.System.Value
		case "service":
			ncon.Key = "name"
			ncon.Value = constr.Service.Name
		case "attribute":
			ncon.Key = constr.Attribute.Name
			ncon.Value = constr.Attribute.Value
		}
		treechk.Constraints[i] = ncon
	}
	return treechk, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
