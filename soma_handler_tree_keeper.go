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
	finish_job        *sql.Stmt
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
	tk.start_job, err = tk.conn.Prepare(`
UPDATE soma.jobs
SET    job_started = $2::timestamptz,
       job_status = 'in_progress'
WHERE  job_id = $1::uuid
AND    job_started IS NULL;`)
	if err != nil {
		log.Fatal("treekeeper/start-job: ", err)
	}
	defer tk.start_job.Close()

	log.Println("Prepare: treekeeper/finish-job")
	tk.finish_job, err = tk.conn.Prepare(`
UPDATE soma.jobs
SET    job_finished = $2::timestamptz,
       job_status = 'processed',
	   job_result = $3::varchar
WHERE  job_id = $1::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/finish-jobs: ", err)
	}
	defer tk.finish_job.Close()

	log.Println("Prepare: treekeeper/create-bucket")
	tk.create_bucket, err = tk.conn.Prepare(`
INSERT INTO soma.buckets (
	bucket_id,
	bucket_name,
	bucket_frozen,
	bucket_deleted,
	repository_id,
	environment,
	organizational_team_id)
SELECT	$1::uuid,
        $2::varchar,
        $3::boolean,
        $4::boolean,
        $5::uuid,
        $6::varchar,
        $7::uuid;`)
	if err != nil {
		log.Fatal("treekeeper/create-bucket: ", err)
	}
	defer tk.create_bucket.Close()

	log.Println("Prepare: treekeeper/defer-constraints")
	tk.defer_constraints, err = tk.conn.Prepare(`
SET CONSTRAINTS ALL DEFERRED;`)
	if err != nil {
		log.Fatal("treekeeper/defer-constraints: ", err)
	}
	defer tk.defer_constraints.Close()

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
		err                    error
		tx                     *sql.Tx
		txStmtCreateGroup      *sql.Stmt
		txStmtCreateCluster    *sql.Stmt
		txStmtBucketAssignNode *sql.Stmt
		txStmtBucketRemoveNode *sql.Stmt
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
	}
	// open multi-statement transaction
	if tx, err = tk.conn.Begin(); err != nil {
		goto bailout
	}
	defer tx.Rollback()

	// prepare statements within tx context
	if txStmtCreateGroup, err = tx.Prepare(tkStmtCreateGroup); err != nil {
		goto bailout
	}
	defer txStmtCreateGroup.Close()
	if txStmtCreateCluster, err = tx.Prepare(tkStmtCreateCluster); err != nil {
		goto bailout
	}
	defer txStmtCreateCluster.Close()
	if txStmtBucketAssignNode, err = tx.Prepare(tkStmtBucketAssignNode); err != nil {
		goto bailout
	}
	defer txStmtBucketAssignNode.Close()
	if txStmtBucketRemoveNode, err = tx.Prepare(tkStmtBucketRemoveNode); err != nil {
		goto bailout
	}
	defer txStmtBucketRemoveNode.Close()

	// defer constraint checks
	if _, err = tx.Exec(tkStmtDeferConstraints); err != nil {
		goto bailout
	}

actionloop:
	for i := 0; i < len(tk.actionChan); i++ {
		a := <-tk.actionChan
		switch a.Type {
		case "bucket":
			switch a.Action {
			case "create":
				if _, err = tx.Stmt(tk.create_bucket).Exec(
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
			case "node_removal":
				if _, err = txStmtBucketRemoveNode.Exec(
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
		case "group":
			switch a.Action {
			case "create":
				if _, err = txStmtCreateGroup.Exec(
					a.Group.Id,
					a.Group.BucketId,
					a.Group.Name,
					a.Group.ObjectState,
					a.Group.TeamId,
				); err != nil {
					break actionloop
				}
			}
		case "cluster":
			switch a.Action {
			case "create":
				if _, err = txStmtCreateCluster.Exec(
					a.Cluster.Id,
					a.Cluster.Name,
					a.Cluster.BucketId,
					a.Cluster.ObjectState,
					a.Cluster.TeamId,
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
	if _, err = tx.Stmt(tk.finish_job).Exec(
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
	tk.finish_job.Exec(q.JobId.String(), time.Now().UTC(), "failed")
	for i := 0; i < len(tk.actionChan); i++ {
		a := <-tk.actionChan
		jB, _ := json.Marshal(a)
		log.Printf("Cleaned message: %s\n", string(jB))
	}
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
