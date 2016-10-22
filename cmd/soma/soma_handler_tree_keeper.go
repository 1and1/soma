package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/internal/tree"
	log "github.com/Sirupsen/logrus"
	metrics "github.com/rcrowley/go-metrics"
	uuid "github.com/satori/go.uuid"
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
	repoId               string
	repoName             string
	team                 string
	rbLevel              string
	broken               bool
	ready                bool
	stopped              bool
	frozen               bool
	rebuild              bool
	input                chan treeRequest
	shutdown             chan bool
	stopchan             chan bool
	errChan              chan *tree.Error
	actionChan           chan *tree.Action
	conn                 *sql.DB
	tree                 *tree.Tree
	get_view             *sql.Stmt
	start_job            *sql.Stmt
	stmt_CapMonMetric    *sql.Stmt
	stmt_Check           *sql.Stmt
	stmt_CheckConfig     *sql.Stmt
	stmt_CheckInstance   *sql.Stmt
	stmt_Cluster         *sql.Stmt
	stmt_ClusterCustProp *sql.Stmt
	stmt_ClusterOncall   *sql.Stmt
	stmt_ClusterService  *sql.Stmt
	stmt_ClusterSysProp  *sql.Stmt
	stmt_DefaultDC       *sql.Stmt
	stmt_DelDuplicate    *sql.Stmt
	stmt_GetComputed     *sql.Stmt
	stmt_GetPrevious     *sql.Stmt
	stmt_Group           *sql.Stmt
	stmt_GroupCustProp   *sql.Stmt
	stmt_GroupOncall     *sql.Stmt
	stmt_GroupService    *sql.Stmt
	stmt_GroupSysProp    *sql.Stmt
	stmt_List            *sql.Stmt
	stmt_Node            *sql.Stmt
	stmt_NodeCustProp    *sql.Stmt
	stmt_NodeOncall      *sql.Stmt
	stmt_NodeService     *sql.Stmt
	stmt_NodeSysProp     *sql.Stmt
	stmt_Pkgs            *sql.Stmt
	stmt_Team            *sql.Stmt
	stmt_Threshold       *sql.Stmt
	stmt_Update          *sql.Stmt
	appLog               *log.Logger
	reqLog               *log.Logger
	errLog               *log.Logger
}

// run() is the method a treeKeeper executes in its background
// go-routine. It checks and handles the input channels and reacts
// appropriately.
func (tk *treeKeeper) run() {
	tk.appLog.Printf("Starting TreeKeeper for Repo %s (%s)", tk.repoName, tk.repoId)
	c := metrics.GetOrRegisterCounter(`.treekeeper.count`, Metrics[`soma`])
	c.Inc(1)
	defer c.Dec(1)

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
broken:
	if tk.broken {
		tickTack := time.NewTicker(time.Second * 10).C
	hoverloop:
		for {
			select {
			case <-tickTack:
				tk.errLog.Printf("TK[%s]: BROKEN REPOSITORY %s flying holding patterns!\n",
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
	for statement, prepStmt := range map[string]*sql.Stmt{
		tkStmtDeleteDuplicateDetails:                  tk.stmt_DelDuplicate,
		tkStmtDeployDetailClusterCustProp:             tk.stmt_ClusterCustProp,
		tkStmtDeployDetailClusterSysProp:              tk.stmt_ClusterSysProp,
		tkStmtDeployDetailDefaultDatacenter:           tk.stmt_DefaultDC,
		tkStmtDeployDetailNodeCustProp:                tk.stmt_NodeCustProp,
		tkStmtDeployDetailNodeSysProp:                 tk.stmt_NodeSysProp,
		tkStmtDeployDetailsCapabilityMonitoringMetric: tk.stmt_CapMonMetric,
		tkStmtDeployDetailsCheck:                      tk.stmt_Check,
		tkStmtDeployDetailsCheckConfig:                tk.stmt_CheckConfig,
		tkStmtDeployDetailsCheckConfigThreshold:       tk.stmt_Threshold,
		tkStmtDeployDetailsCheckInstance:              tk.stmt_CheckInstance,
		tkStmtDeployDetailsCluster:                    tk.stmt_Cluster,
		tkStmtDeployDetailsClusterOncall:              tk.stmt_ClusterOncall,
		tkStmtDeployDetailsClusterService:             tk.stmt_ClusterService,
		tkStmtDeployDetailsComputeList:                tk.stmt_List,
		tkStmtDeployDetailsGroup:                      tk.stmt_Group,
		tkStmtDeployDetailsGroupCustProp:              tk.stmt_GroupCustProp,
		tkStmtDeployDetailsGroupOncall:                tk.stmt_GroupOncall,
		tkStmtDeployDetailsGroupService:               tk.stmt_GroupService,
		tkStmtDeployDetailsGroupSysProp:               tk.stmt_GroupSysProp,
		tkStmtDeployDetailsNode:                       tk.stmt_Node,
		tkStmtDeployDetailsNodeOncall:                 tk.stmt_NodeOncall,
		tkStmtDeployDetailsNodeService:                tk.stmt_NodeService,
		tkStmtDeployDetailsProviders:                  tk.stmt_Pkgs,
		tkStmtDeployDetailsTeam:                       tk.stmt_Team,
		tkStmtDeployDetailsUpdate:                     tk.stmt_Update,
		tkStmtGetComputedDeployments:                  tk.stmt_GetComputed,
		tkStmtGetPreviousDeployment:                   tk.stmt_GetPrevious,
		tkStmtGetViewFromCapability:                   tk.get_view,
		tkStmtStartJob:                                tk.start_job,
	} {
		if prepStmt, err = tk.conn.Prepare(statement); err != nil {
			tk.errLog.Println("Error preparing SQL statement: ", err)
			tk.errLog.Println("Failed statement: ", statement)
			tk.broken = true
			goto broken
		}
		defer prepStmt.Close()
	}

	// TODO per-treekeeper logfiles:
	// ${SomaCfg.LogPath}/repository/${keepername}.log  <- registered rotate
	// ${SomaCfg.LogPath}/repository/${keepername}_startup.${rfc3339Milli}.log

	tk.appLog.Printf("TK[%s]: ready for service!\n", tk.repoName)
	tk.ready = true

	if SomaCfg.Observer {
		// XXX should listen on stopchan
		tk.appLog.Printf("TreeKeeper [%s] entered observer mode\n", tk.repoName)
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

		tk.appLog.Printf("TreeKeeper [%s] has stopped", tk.repoName)
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
				// buildDeploymentDetails and orderDeploymentDetails can
				// both mark the tree as broken if there was an error
				// preparing required SQL statements
				tk.buildDeploymentDetails()
				if tk.broken {
					goto broken
				}
				tk.orderDeploymentDetails()
				if tk.broken {
					goto broken
				}
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
		err                                   error
		hasErrors, hasJobLog, jobNeverStarted bool
		tx                                    *sql.Tx
		stm                                   map[string]*sql.Stmt
		jobLog                                *log.Logger
		lfh                                   *os.File
	)

	if !tk.rebuild {
		_, err = tk.start_job.Exec(q.JobId.String(), time.Now().UTC())
		if err != nil {
			tk.errLog.Println("Failed starting job %s: %s\n",
				q.JobId.String(),
				err)
			jobNeverStarted = true
			goto bailout
		}
		tk.appLog.Printf("Processing job: %s\n", q.JobId.String())
	} else {
		tk.appLog.Printf("Processing rebuild job: %s\n", q.JobId.String())
	}
	if lfh, err = os.Create(filepath.Join(
		SomaCfg.LogPath,
		`job`,
		fmt.Sprintf("%s_%s_%s.log",
			time.Now().UTC().Format(rfc3339Milli),
			tk.repoName,
			q.JobId.String(),
		),
	)); err != nil {
		tk.errLog.Printf("Failed opening joblog %s: %s\n",
			q.JobId.String(),
			err)
	}
	defer lfh.Close()
	defer lfh.Sync()
	jobLog = log.New()
	log.SetOutput(lfh)
	hasJobLog = true

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
	if tx, stm, err = tk.startTx(); err != nil {
		goto bailout
	}
	defer tx.Rollback()

	// defer constraint checks
	if _, err = tx.Exec(tkStmtDeferAllConstraints); err != nil {
		tk.errLog.Println("Failed to exec: tkStmtDeferAllConstraints")
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
		if _, err = tx.Exec(
			stmt.TxMarkCheckConfigDeleted,
			q.CheckConfig.CheckConfig.Id,
		); err != nil {
			goto bailout
		}
	}

	// if the error channel has entries, we can fully ignore the
	// action channel
	for i := len(tk.errChan); i > 0; i-- {
		e := <-tk.errChan
		if hasJobLog {
			b, _ := json.Marshal(e)
			jobLog.Println(string(b))
		}
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

		// log all actions for the job
		if hasJobLog {
			b, _ := json.Marshal(a)
			jobLog.Println(string(b))
		}

		// only check and check_instance actions are relevant during
		// a rebuild, everything else is ignored. Even some deletes are
		// valid, for example when a property overwrites inheritance of
		// another property, the first will generate deletes.
		// Other deletes should not occur, like node/delete, but will be
		// sorted later. TODO
		if tk.rebuild {
			if tk.rbLevel == `instances` {
				switch a.Action {
				case `check_new`, `check_removed`:
					// ignore only in instance-rebuild mode
					continue actionloop
				}
			}
			switch a.Action {
			case `property_new`, `property_delete`,
				`create`, `update`, `delete`,
				`node_assignment`,
				`member_new`, `member_removed`:
				// ignore in all rebuild modes
				continue actionloop
			}
		}

		switch a.Action {
		case `property_new`, `property_delete`:
			if err = tk.txProperty(a, stm); err != nil {
				break actionloop
			}
		case `check_new`, `check_removed`:
			if err = tk.txCheck(a, stm); err != nil {
				break actionloop
			}
		case `check_instance_create`,
			`check_instance_update`,
			`check_instance_delete`:
			if err = tk.txCheckInstance(a, stm); err != nil {
				break actionloop
			}
		case `create`, `update`, `delete`, `node_assignment`,
			`member_new`, `member_removed`:
			if err = tk.txTree(a, stm, q.User); err != nil {
				break actionloop
			}
		default:
			err = fmt.Errorf(
				"Unhandled message in action stream: %s/%s",
				a.Type,
				a.Action,
			)
			break actionloop
		}

		switch a.Type {
		case "errorchannel":
			continue actionloop
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
	tk.appLog.Printf("SUCCESS - Finished job: %s\n", q.JobId.String())

	// accept tree changes
	tk.tree.Commit()
	return

bailout:
	tk.appLog.Printf("FAILED - Finished job: %s\n", q.JobId.String())
	tk.errLog.Printf("Job-Error(%s): %s\n", q.JobId.String(), err)
	if hasJobLog {
		jobLog.Printf("Aborting error: %s\n", err)
	}

	// if this was a rebuild, the tree will not persist and the
	// job is faked. Also if the job never actually started, then it
	// should never be rolled back nor attempted to mark failed.
	if tk.rebuild || jobNeverStarted {
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
		if hasJobLog {
			jobLog.Printf("Cleaned message: %s\n", string(jB))
		}
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
