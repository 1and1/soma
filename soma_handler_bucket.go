package main

import (
	"database/sql"
	"log"

)

type somaBucketRequest struct {
	action string
	Bucket proto.Bucket
	reply  chan somaResult
}

type somaBucketResult struct {
	ResultError error
	Bucket      proto.Bucket
}

func (a *somaBucketResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Buckets = append(r.Buckets, somaBucketResult{ResultError: err})
	}
}

func (a *somaBucketResult) SomaAppendResult(r *somaResult) {
	r.Buckets = append(r.Buckets, *a)
}

/* Read Access
 */
type somaBucketReadHandler struct {
	input     chan somaBucketRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
}

func (r *somaBucketReadHandler) run() {
	var err error

	log.Println("Prepare: bucket/list")
	if r.list_stmt, err = r.conn.Prepare(stmtBucketList); err != nil {
		log.Fatal("bucket/list: ", err)
	}
	defer r.list_stmt.Close()

	log.Println("Prepare: bucket/show")
	if r.show_stmt, err = r.conn.Prepare(stmtBucketShow); err != nil {
		log.Fatal("bucket/show: ", err)
	}
	defer r.show_stmt.Close()

runloop:
	for {
		select {
		case <-r.shutdown:
			break runloop
		case req := <-r.input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

func (r *somaBucketReadHandler) process(q *somaBucketRequest) {
	var (
		bucketId, bucketName, bucketEnv, repoId, teamId string
		rows                                            *sql.Rows
		bucketDeleted, bucketFrozen                     bool
		err                                             error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: bucket/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&bucketId, &bucketName)
			result.Append(err, &somaBucketResult{
				Bucket: somaproto.ProtoBucket{
					Id:   bucketId,
					Name: bucketName,
				},
			})
		}
	case "show":
		log.Printf("R: bucket/show for %s", q.Bucket.Id)
		err = r.show_stmt.QueryRow(q.Bucket.Id).Scan(
			&bucketId,
			&bucketName,
			&bucketFrozen,
			&bucketDeleted,
			&repoId,
			&bucketEnv,
			&teamId,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		result.Append(err, &somaBucketResult{
			Bucket: somaproto.ProtoBucket{
				Id:          bucketId,
				Name:        bucketName,
				Repository:  repoId,
				Team:        teamId,
				Environment: bucketEnv,
				IsDeleted:   bucketDeleted,
				IsFrozen:    bucketFrozen,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
