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

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/internal/tree"
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

// run() is the method a treeKeeper executes in its background
// go-routine. It checks and handles the input channels and reacts
// appropriately.
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
	// then exit. We issue a fake job for this.
	if tk.rebuild {
		req := treeRequest{
			RequestType: `rebuild`,
			Action:      `rebuild`,
			JobId:       uuid.NewV4(),
		}
		tk.process(&req)
		tk.buildDeploymentDetails()
		tk.orderDeploymentDetails()
		tk.conn.Close()
		return
	}

	// there was an error during startupLoad(), the repository is
	// considered broken.
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

	// prepare statements
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
		// XXX should listen on stopchan
		log.Printf("TreeKeeper [%s] entered observer mode\n", tk.repoName)
		<-tk.shutdown
		goto exit
	}

stopsign:
	if tk.stopped {
		// drain the input channel, it could be currently full and
		// writers blocked on it. Future writers will check
		// isStopped() before writing (and/or remove this tree from
		// the handlerMap)
	drain:
		for i := len(tk.input); i > 0; i-- {
			<-tk.input
		}
		if len(tk.input) > 0 {
			// there were blocked writers on a full buffered channel
			goto drain
		}

		log.Printf("TreeKeeper [%s] has stopped", tk.repoName)
		for {
			select {
			case <-tk.shutdown:
				goto exit
			case <-tk.stopchan:
			}
		}
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
			handlerMap[`jobDelay`].(*jobDelay).notify <- req.JobId.String()
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
		nullBucket sql.NullString
		stm        map[string]*sql.Stmt
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
	// TREE MANIPULATION REQUESTS
	case
		`create_bucket`:
		tk.treeBucket(q)

	case
		`create_group`,
		`delete_group`,
		`reset_group_to_bucket`,
		`add_group_to_group`:
		tk.treeGroup(q)

	case
		`create_cluster`,
		`delete_cluster`,
		`reset_cluster_to_bucket`,
		`add_cluster_to_group`:
		tk.treeCluster(q)

	case
		"assign_node",
		"delete_node",
		"reset_node_to_bucket",
		"add_node_to_group",
		"add_node_to_cluster":
		tk.treeNode(q)

	//
	// PROPERTY MANIPULATION REQUESTS
	case
		`add_system_property_to_repository`,
		`add_system_property_to_bucket`,
		`add_system_property_to_group`,
		`add_system_property_to_cluster`,
		`add_system_property_to_node`,
		`add_service_property_to_repository`,
		`add_service_property_to_bucket`,
		`add_service_property_to_group`,
		`add_service_property_to_cluster`,
		`add_service_property_to_node`,
		`add_oncall_property_to_repository`,
		`add_oncall_property_to_bucket`,
		`add_oncall_property_to_group`,
		`add_oncall_property_to_cluster`,
		`add_oncall_property_to_node`,
		`add_custom_property_to_repository`,
		`add_custom_property_to_bucket`,
		`add_custom_property_to_group`,
		`add_custom_property_to_cluster`,
		`add_custom_property_to_node`:
		tk.addProperty(q)

	case
		`delete_system_property_from_repository`,
		`delete_system_property_from_bucket`,
		`delete_system_property_from_group`,
		`delete_system_property_from_cluster`,
		`delete_system_property_from_node`,
		`delete_service_property_from_repository`,
		`delete_service_property_from_bucket`,
		`delete_service_property_from_group`,
		`delete_service_property_from_cluster`,
		`delete_service_property_from_node`,
		`delete_oncall_property_from_repository`,
		`delete_oncall_property_from_bucket`,
		`delete_oncall_property_from_group`,
		`delete_oncall_property_from_cluster`,
		`delete_oncall_property_from_node`,
		`delete_custom_property_from_repository`,
		`delete_custom_property_from_bucket`,
		`delete_custom_property_from_group`,
		`delete_custom_property_from_cluster`,
		`delete_custom_property_from_node`:
		tk.rmProperty(q)

	//
	// CHECK MANIPULATION REQUESTS
	case
		`add_check_to_repository`,
		`add_check_to_bucket`,
		`add_check_to_group`,
		`add_check_to_cluster`,
		`add_check_to_node`:
		err = tk.addCheck(&q.CheckConfig.CheckConfig)

	case
		`remove_check_from_repository`,
		`remove_check_from_bucket`,
		`remove_check_from_group`,
		`remove_check_from_cluster`,
		`remove_check_from_node`:
		err = tk.rmCheck(&q.CheckConfig.CheckConfig)
	}

	// check if we accumulated an error in one of the switch cases
	if err != nil {
		goto bailout
	}

	// recalculate check instances
	tk.tree.ComputeCheckInstances()

	// open multi-statement transaction
	if tx, err, stm = tk.startTx(); err != nil {
		goto bailout
	}
	defer tx.Rollback()

	// defer constraint checks
	if _, err = tx.Exec(tkStmtDeferAllConstraints); err != nil {
		log.Println("Failed to exec: tkStmtDeferAllConstraints")
		goto bailout
	}

	// save the check configuration as part of the transaction before
	// processing the action channel
	if strings.Contains(q.Action, "add_check_to_") {
		if err = tk.txCheckConfig(q.CheckConfig.CheckConfig,
			stm); err != nil {
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

/* Ops Access
 */
func (tk *treeKeeper) shutdownNow() {
	tk.shutdown <- true
}

func (tk *treeKeeper) stopNow() {
	tk.stopchan <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
