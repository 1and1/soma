package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/satori/go.uuid"
)

type guidePost struct {
	input     chan treeRequest
	shutdown  chan bool
	conn      *sql.DB
	jbsv_stmt *sql.Stmt
	repo_stmt *sql.Stmt
	name_stmt *sql.Stmt
}

func (g *guidePost) run() {
	var err error

	log.Println("Prepare: guide/job-save")
	g.jbsv_stmt, err = g.conn.Prepare(`
INSERT INTO soma.jobs (
	job_id,
	job_status,
	job_result,
	job_type,
	repository_id,
	user_id,
	organizational_team_id,
	job)
SELECT	$1::uuid,
		$2::varchar,
		$3::varchar,
		$4::varchar,
		$5::uuid,
		$6::uuid,
		$7::uuid,
		$8::jsonb;`)
	if err != nil {
		log.Fatal("guide/job-save: ", err)
	}
	defer g.jbsv_stmt.Close()

	log.Println("Prepare: guide/repo-by-bucket")
	g.repo_stmt, err = g.conn.Prepare(`
SELECT	sb.bucket_id,
		sb.repository_id,
		sr.repository_name
FROM	soma.buckets sb
JOIN    soma.repositories sr
ON		sb.repository_id = sr.repository_id
WHERE	sb.bucket_id = $1::uuid;`)
	if err != nil {
		log.Fatal("guide/repo-by-bucket: ", err)
	}
	defer g.repo_stmt.Close()

	log.Println("Prepare: guide/repo-by-id")
	g.name_stmt, err = g.conn.Prepare(`
SELECT repository_name
FROM   soma.repositories
WHERE  repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("guide/repo-by-id: ", err)
	}
	defer g.name_stmt.Close()

runloop:
	for {
		select {
		case <-g.shutdown:
			break runloop
		case req := <-g.input:
			g.process(&req)
		}
	}
}

func (g *guidePost) process(q *treeRequest) {
	var (
		res              sql.Result
		err              error
		j                []byte
		repoId, repoName string
	)
	result := somaResult{}

	log.Printf("R: jobsave/%s", q.Action)
	switch q.Action {
	case "create_bucket":
		repoId = q.Bucket.Bucket.Repository
	case "create_group":
	default:
		log.Printf("R: unimplemented server/%s", q.Action)
		result.SetNotImplemented()
		q.reply <- result
		return
	}
	q.JobId = uuid.NewV4()
	j, _ = json.Marshal(q)
	res, err = g.jbsv_stmt.Exec(
		q.JobId.String(),
		"queued",
		"pending",
		q.Action,
		repoId,
		"00000000-0000-0000-0000-000000000000",
		"00000000-0000-0000-0000-000000000000",
		string(j),
	)
	if result.SetRequestError(err) {
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		result.Append(errors.New("No rows affected"), &somaBucketResult{})
		q.reply <- result
		return
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaBucketResult{})
		q.reply <- result
		return
	}

	switch q.Action {
	case "create_bucket":
		err = g.name_stmt.QueryRow(repoId).Scan(&repoName)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		keeper := fmt.Sprintf("repository_%s", repoName)
		handler := handlerMap[keeper].(treeKeeper)
		handler.input <- *q
	}
	result.Append(nil, &somaBucketResult{
		Bucket: q.Bucket.Bucket,
	})
	result.JobId = q.JobId.String()
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
