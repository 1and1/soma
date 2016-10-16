package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

type guidePost struct {
	input              chan treeRequest
	system             chan msg.Request
	shutdown           chan bool
	conn               *sql.DB
	jbsv_stmt          *sql.Stmt
	repo_stmt          *sql.Stmt
	name_stmt          *sql.Stmt
	node_stmt          *sql.Stmt
	serv_stmt          *sql.Stmt
	attr_stmt          *sql.Stmt
	cthr_stmt          *sql.Stmt
	cdel_stmt          *sql.Stmt
	bucket_for_node    *sql.Stmt
	bucket_for_cluster *sql.Stmt
	bucket_for_group   *sql.Stmt
	appLog             *log.Logger
	reqLog             *log.Logger
	errLog             *log.Logger
}

func (g *guidePost) run() {
	var err error

	if g.jbsv_stmt, err = g.conn.Prepare(stmt.JobSave); err != nil {
		log.Fatal("guide/job-save: ", err)
	}
	defer g.jbsv_stmt.Close()

	if g.repo_stmt, err = g.conn.Prepare(stmt.RepoByBucketId); err != nil {
		log.Fatal("guide/repo-by-bucket: ", err)
	}
	defer g.repo_stmt.Close()

	if g.node_stmt, err = g.conn.Prepare(stmt.NodeDetails); err != nil {
		log.Fatal("guide/load-node-details: ", err)
	}
	defer g.node_stmt.Close()

	if g.name_stmt, err = g.conn.Prepare(stmt.RepoNameById); err != nil {
		log.Fatal("guide/repo-by-id: ", err)
	}
	defer g.name_stmt.Close()

	if g.serv_stmt, err = g.conn.Prepare(stmt.ServiceLookup); err != nil {
		log.Fatal("guide/service-lookup: ", err)
	}
	defer g.serv_stmt.Close()

	if g.attr_stmt, err = g.conn.Prepare(stmt.ServiceAttributes); err != nil {
		log.Fatal("guide/populate-service-attributes: ", err)
	}
	defer g.attr_stmt.Close()

	if g.cthr_stmt, err = g.conn.Prepare(stmt.CapabilityThresholds); err != nil {
		log.Fatal("guide/capability-threshold-lookup: ", err)
	}
	defer g.cthr_stmt.Close()

	if g.cdel_stmt, err = g.conn.Prepare(stmt.CheckDetailsForDelete); err != nil {
		log.Fatal("guide/get-details-for-delete-check: ", err)
	}
	defer g.cdel_stmt.Close()

	if g.bucket_for_node, err = g.conn.Prepare(stmt.NodeBucketId); err != nil {
		log.Fatal("guide/get-bucketid-for-node: ", err)
	}
	defer g.bucket_for_node.Close()

	if g.bucket_for_cluster, err = g.conn.Prepare(stmt.ClusterBucketId); err != nil {
		log.Fatal("guide/get-bucketid-for-cluster: ", err)
	}
	defer g.bucket_for_cluster.Close()

	if g.bucket_for_group, err = g.conn.Prepare(stmt.GroupBucketId); err != nil {
		log.Fatal("guide/get-bucketid-for-group: ", err)
	}
	defer g.bucket_for_group.Close()

	if SomaCfg.Observer {
		// XXX system/stop_repository should be possible in observer
		// mode
		log.Println(`GuidePost entered observer mode`)
		<-g.shutdown
		goto exit
	}

runloop:
	for {
		select {
		case <-g.shutdown:
			break runloop
		case req := <-g.input:
			g.process(&req)
		case req := <-g.system:
			g.sysprocess(&req)
		}
	}
exit:
}

func (g *guidePost) process(q *treeRequest) {
	var (
		res                      sql.Result
		err                      error
		j                        []byte
		repoId, repoName, keeper string
		nf                       bool
		handler                  *treeKeeper
		rowCnt                   int64
	)
	result := somaResult{}

	// to which tree this request must be forwarded
	if repoId, repoName, err, nf = g.extractRouting(q); err != nil {
		goto bailout
	}

	// verify we can process the request
	if err, nf = g.validateRequest(q); err != nil {
		goto bailout
	}

	// fill in required data for the request
	if err, nf = g.fillReqData(q); err != nil {
		goto bailout
	}

	// check we have a treekeeper for that repository
	if err, nf = g.validateKeeper(repoName); err != nil {
		goto bailout
	}
	keeper = fmt.Sprintf("repository_%s", repoName)
	handler = handlerMap[keeper].(*treeKeeper)

	// store job in database
	log.Printf("R: jobsave/%s", q.Action)
	q.JobId = uuid.NewV4()
	j, _ = json.Marshal(q)
	if res, err = g.jbsv_stmt.Exec(
		q.JobId.String(),
		`queued`,
		`pending`,
		q.Action,
		repoId,
		q.User,
		string(j),
	); err != nil {
		goto bailout
	}
	// insert can have 0 rows affected if the where clause could
	// not find the user
	rowCnt, _ = res.RowsAffected()
	if rowCnt == 0 {
		err = fmt.Errorf("No rows affected while saving job for user %s",
			q.User)
		nf = false
		goto bailout
	}

	handler.input <- *q
	result.JobId = q.JobId.String()
	result.JobType = q.Action

	switch q.RequestType {
	case `repository`:
		result.Append(nil, &somaRepositoryResult{
			Repository: q.Repository.Repository,
		})
	case `bucket`:
		result.Append(nil, &somaBucketResult{
			Bucket: q.Bucket.Bucket,
		})
	case `group`:
		result.Append(nil, &somaGroupResult{
			Group: q.Group.Group,
		})
	case `cluster`:
		result.Append(nil, &somaClusterResult{
			Cluster: q.Cluster.Cluster,
		})
	case `node`:
		result.Append(nil, &somaNodeResult{
			Node: q.Node.Node,
		})
	case `check`:
		result.Append(nil, &somaCheckConfigResult{
			CheckConfig: q.CheckConfig.CheckConfig,
		})
	}

bailout:
	if err != nil {
		if nf {
			result.SetNotFoundErr(err)
		} else {
			result.SetRequestError(err)
		}
	}
	q.reply <- result
}

//
// Process system operation requests
func (g *guidePost) sysprocess(q *msg.Request) {
	var (
		repoName, repoId, keeper string
		err                      error
		handler                  *treeKeeper
	)
	result := msg.Result{Type: `guidepost`, Action: `systemoperation`, System: []proto.SystemOperation{q.System}}

	switch q.System.Request {
	case `stop_repository`:
		repoId = q.System.RepositoryId
	default:
		result.NotImplemented(
			fmt.Errorf("Unknown requested system operation: %s",
				q.System.Request),
		)
		goto exit
	}

	if err = g.name_stmt.QueryRow(repoId).Scan(&repoName); err != nil {
		if err == sql.ErrNoRows {
			result.NotFound(fmt.Errorf(`No such repository`))
		} else {
			result.ServerError(err)
		}
		goto exit
	}

	// check we have a treekeeper for that repository
	keeper = fmt.Sprintf("repository_%s", repoName)
	if _, ok := handlerMap[keeper].(*treeKeeper); !ok {
		// no handler running, nothing to stop
		result.OK()
		goto exit
	}

	// might already be stopped
	handler = handlerMap[keeper].(*treeKeeper)
	if handler.isStopped() {
		result.OK()
		goto exit
	}

	// check the treekeeper is ready for system requests
	if !(handler.isReady() || handler.isBroken()) {
		result.Unavailable(
			fmt.Errorf("Repository %s not fully loaded yet.",
				repoName),
		)
		goto exit
	}

	switch q.System.Request {
	case `stop_repository`:
		if !handler.isStopped() {
			handler.stopchan <- true
		}
		result.OK()
	}

exit:
	q.Reply <- result
}

/* Ops Access
 */
func (g *guidePost) shutdownNow() {
	g.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
