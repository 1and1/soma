package main

import (
	"database/sql"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
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
	ponc_stmt *sql.Stmt
	psvc_stmt *sql.Stmt
	psys_stmt *sql.Stmt
	pcst_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaBucketReadHandler) run() {
	var err error

	if r.list_stmt, err = r.conn.Prepare(stmtBucketList); err != nil {
		log.Fatal("bucket/list: ", err)
	}
	defer r.list_stmt.Close()

	if r.show_stmt, err = r.conn.Prepare(stmtBucketShow); err != nil {
		log.Fatal("bucket/show: ", err)
	}
	defer r.show_stmt.Close()

	if r.ponc_stmt, err = r.conn.Prepare(stmt.BucketOncProps); err != nil {
		log.Fatal(`bucket/property-oncall: `, err)
	}
	defer r.ponc_stmt.Close()

	if r.psvc_stmt, err = r.conn.Prepare(stmt.BucketSvcProps); err != nil {
		log.Fatal(`bucket/property-service: `, err)
	}
	defer r.psvc_stmt.Close()

	if r.psys_stmt, err = r.conn.Prepare(stmt.BucketSysProps); err != nil {
		log.Fatal(`bucket/property-system: `, err)
	}
	defer r.psys_stmt.Close()

	if r.pcst_stmt, err = r.conn.Prepare(stmt.BucketCstProps); err != nil {
		log.Fatal(`bucket/property-custom: `, err)
	}
	defer r.pcst_stmt.Close()

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
		instanceId, sourceInstanceId, view, oncallId    string
		oncallName, serviceName, customId, systemProp   string
		value, customProp                               string
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
			goto dispatch
		}

		for rows.Next() {
			err := rows.Scan(&bucketId, &bucketName)
			result.Append(err, &somaBucketResult{
				Bucket: proto.Bucket{
					Id:   bucketId,
					Name: bucketName,
				},
			})
		}
	case `show`:
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
			goto dispatch
		}
		bucket := proto.Bucket{
			Id:           bucketId,
			Name:         bucketName,
			RepositoryId: repoId,
			TeamId:       teamId,
			Environment:  bucketEnv,
			IsDeleted:    bucketDeleted,
			IsFrozen:     bucketFrozen,
		}
		bucket.Properties = &[]proto.Property{}

		// oncall properties
		rows, err = r.ponc_stmt.Query(q.Bucket.Id)
		if result.SetRequestError(err) {
			goto dispatch
		}
		for rows.Next() {
			if err := rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&oncallId,
				&oncallName,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*bucket.Properties = append(
				*bucket.Properties,
				proto.Property{
					Type:             `oncall`,
					RepositoryId:     repoId,
					BucketId:         bucketId,
					InstanceId:       instanceId,
					SourceInstanceId: sourceInstanceId,
					View:             view,
					Oncall: &proto.PropertyOncall{
						Id:   oncallId,
						Name: oncallName,
					},
				},
			)
		}
		if err = rows.Err(); err != nil {
			result.SetRequestError(err)
			goto dispatch
		}

		// service properties
		rows, err = r.psvc_stmt.Query(q.Bucket.Id)
		for rows.Next() {
			if err := rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&serviceName,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*bucket.Properties = append(
				*bucket.Properties,
				proto.Property{
					Type:             `service`,
					RepositoryId:     repoId,
					BucketId:         bucketId,
					InstanceId:       instanceId,
					SourceInstanceId: sourceInstanceId,
					View:             view,
					Service: &proto.PropertyService{
						Name: serviceName,
					},
				},
			)
		}
		if err = rows.Err(); err != nil {
			result.SetRequestError(err)
			goto dispatch
		}

		// system properties
		rows, err = r.psys_stmt.Query(q.Bucket.Id)
		for rows.Next() {
			if err := rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&systemProp,
				&value,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*bucket.Properties = append(
				*bucket.Properties,
				proto.Property{
					Type:             `system`,
					RepositoryId:     repoId,
					BucketId:         bucketId,
					InstanceId:       instanceId,
					SourceInstanceId: sourceInstanceId,
					View:             view,
					System: &proto.PropertySystem{
						Name:  systemProp,
						Value: value,
					},
				},
			)
		}
		if err = rows.Err(); err != nil {
			result.SetRequestError(err)
			goto dispatch
		}

		// custom properties
		rows, err = r.pcst_stmt.Query(q.Bucket.Id)
		for rows.Next() {
			if err := rows.Scan(
				&instanceId,
				&sourceInstanceId,
				&view,
				&customId,
				&value,
				&customProp,
			); result.SetRequestError(err) {
				rows.Close()
				goto dispatch
			}
			*bucket.Properties = append(
				*bucket.Properties,
				proto.Property{
					Type:             `custom`,
					RepositoryId:     repoId,
					BucketId:         bucketId,
					InstanceId:       instanceId,
					SourceInstanceId: sourceInstanceId,
					View:             view,
					Custom: &proto.PropertyCustom{
						Id:    customId,
						Name:  customProp,
						Value: value,
					},
				},
			)
		}
		if err = rows.Err(); err != nil {
			result.SetRequestError(err)
			goto dispatch
		}

		result.Append(err, &somaBucketResult{
			Bucket: bucket,
		})
	default:
		result.SetNotImplemented()
	}

dispatch:
	q.reply <- result
}

/* Ops Access
 */
func (r *somaBucketReadHandler) shutdownNow() {
	r.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
