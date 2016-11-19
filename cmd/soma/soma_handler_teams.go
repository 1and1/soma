package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

type somaTeamRequest struct {
	action string
	Team   proto.Team
	reply  chan somaResult
}

type somaTeamResult struct {
	ResultError error
	Team        proto.Team
}

func (a *somaTeamResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Teams = append(r.Teams, somaTeamResult{ResultError: err})
	}
}

func (a *somaTeamResult) SomaAppendResult(r *somaResult) {
	r.Teams = append(r.Teams, *a)
}

/* Read Access
 */
type somaTeamReadHandler struct {
	input     chan somaTeamRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	sync_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaTeamReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ListTeams: r.list_stmt,
		stmt.ShowTeams: r.show_stmt,
		stmt.SyncTeams: r.sync_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`team`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

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

func (r *somaTeamReadHandler) process(q *somaTeamRequest) {
	var (
		teamId, teamName string
		ldapId           int
		systemFlag       bool
		rows             *sql.Rows
		err              error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		r.reqLog.Printf("R: team/list")
		rows, err = r.list_stmt.Query()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&teamId, &teamName)
			result.Append(err, &somaTeamResult{
				Team: proto.Team{
					Id:   teamId,
					Name: teamName,
				},
			})
		}
	case `sync`:
		r.reqLog.Printf("R: team/sync")
		rows, err = r.sync_stmt.Query()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(
				&teamId,
				&teamName,
				&ldapId,
				&systemFlag,
			)

			result.Append(err, &somaTeamResult{
				Team: proto.Team{
					Id:       teamId,
					Name:     teamName,
					LdapId:   strconv.Itoa(ldapId),
					IsSystem: systemFlag,
				},
			})
		}
	case "show":
		r.reqLog.Printf("R: team/show for %s", q.Team.Id)
		err = r.show_stmt.QueryRow(q.Team.Id).Scan(
			&teamId,
			&teamName,
			&ldapId,
			&systemFlag,
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

		result.Append(err, &somaTeamResult{
			Team: proto.Team{
				Id:       teamId,
				Name:     teamName,
				LdapId:   strconv.Itoa(ldapId),
				IsSystem: systemFlag,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaTeamWriteHandler struct {
	input    chan somaTeamRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	upd_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaTeamWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.TeamAdd:    w.add_stmt,
		stmt.TeamUpdate: w.upd_stmt,
		stmt.TeamDel:    w.del_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`team`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-w.shutdown:
			break runloop
		case req := <-w.input:
			go func() {
				w.process(&req)
			}()
		}
	}
}

func (w *somaTeamWriteHandler) process(q *somaTeamRequest) {
	var (
		res    sql.Result
		err    error
		super  *supervisor
		notify msg.Request
	)
	result := somaResult{}
	super = handlerMap[`supervisor`].(*supervisor)
	notify = msg.Request{Section: `map`, Action: `update`,
		Super: &msg.Supervisor{
			Object: `team`,
			Team:   q.Team,
		},
	}

	switch q.action {
	case "add":
		w.reqLog.Printf("R: team/add for %s", q.Team.Name)
		id := uuid.NewV4()
		res, err = w.add_stmt.Exec(
			id.String(),
			q.Team.Name,
			q.Team.LdapId,
			q.Team.IsSystem,
		)
		q.Team.Id = id.String()
		notify.Action = `add`
	case `update`:
		w.reqLog.Printf("R: team/update for %s", q.Team.Name)
		res, err = w.upd_stmt.Exec(
			q.Team.Name,
			q.Team.LdapId,
			q.Team.IsSystem,
			q.Team.Id,
		)
		notify.Action = `update`
	case "delete":
		w.reqLog.Printf("R: team/del for %s", q.Team.Id)
		res, err = w.del_stmt.Exec(
			q.Team.Id,
		)
		notify.Action = `delete`
	default:
		w.reqLog.Printf("R: unimplemented team/%s", q.action)
		result.SetNotImplemented()
		q.reply <- result
		return
	}
	if result.SetRequestError(err) {
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		result.Append(errors.New("No rows affected"), &somaTeamResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaTeamResult{})
	default:
		result.Append(nil, &somaTeamResult{
			Team: q.Team,
		})
		// send update to supervisor
		super.input <- notify
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaTeamReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaTeamWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
