package main

import (
	"database/sql"
	"log"

)

type somaGroupRequest struct {
	action string
	Group  somaproto.ProtoGroup
	reply  chan somaResult
}

type somaGroupResult struct {
	ResultError error
	Group       somaproto.ProtoGroup
}

func (a *somaGroupResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Groups = append(r.Groups, somaGroupResult{ResultError: err})
	}
}

func (a *somaGroupResult) SomaAppendResult(r *somaResult) {
	r.Groups = append(r.Groups, *a)
}

/* Read Access
 */
type somaGroupReadHandler struct {
	input     chan somaGroupRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
}

func (r *somaGroupReadHandler) run() {
	var err error

	log.Println("Prepare: group/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT group_id,
       group_name
FROM soma.groups;`)
	if err != nil {
		log.Fatal("group/list: ", err)
	}
	defer r.list_stmt.Close()

	log.Println("Prepare: group/show")
	r.show_stmt, err = r.conn.Prepare(`
SELECT group_id,
       bucket_id,
	   group_name,
	   object_state,
	   organizational_team_id
FROM   soma.groups
WHERE  group_id = $1::uuid;`)
	if err != nil {
		log.Fatal("group/show: ", err)
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

func (r *somaGroupReadHandler) process(q *somaGroupRequest) {
	var (
		groupId, groupName, bucketId, groupState, teamId string
		rows                                             *sql.Rows
		err                                              error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: group/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&groupId, &groupName)
			result.Append(err, &somaGroupResult{
				Group: somaproto.ProtoGroup{
					Id:   groupId,
					Name: groupName,
				},
			})
		}
	case "show":
		log.Printf("R: group/show for %s", q.Group.Id)
		err = r.show_stmt.QueryRow(q.Group.Id).Scan(
			&groupId,
			&bucketId,
			&groupName,
			&groupState,
			&teamId,
		)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		result.Append(err, &somaGroupResult{
			Group: somaproto.ProtoGroup{
				Id:          groupId,
				Name:        groupName,
				BucketId:    bucketId,
				ObjectState: groupState,
				TeamId:      teamId,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
