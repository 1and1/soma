package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/1and1/soma/lib/msg"
	"github.com/1and1/soma/lib/stmt"
	"github.com/1and1/soma/lib/proto"
	"github.com/satori/go.uuid"
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
}

func (r *somaTeamReadHandler) run() {
	var err error

	if r.list_stmt, err = r.conn.Prepare(stmt.ListTeams); err != nil {
		log.Fatal("team/list: ", err)
	}
	defer r.list_stmt.Close()

	if r.show_stmt, err = r.conn.Prepare(stmt.ShowTeams); err != nil {
		log.Fatal("team/show: ", err)
	}
	defer r.show_stmt.Close()

	if r.sync_stmt, err = r.conn.Prepare(stmt.SyncTeams); err != nil {
		log.Fatal("team/sync: ", err)
	}
	defer r.sync_stmt.Close()

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
		log.Printf("R: team/list")
		rows, err = r.list_stmt.Query()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&teamId, &teamName)
			result.Append(err, &somaTeamResult{
				Team: proto.Team{
					Id:   teamId,
					Name: teamName,
				},
			})
		}
	case `sync`:
		log.Printf("R: team/sync")
		rows, err = r.sync_stmt.Query()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(
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
		log.Printf("R: team/show for %s", q.Team.Id)
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
}

func (w *somaTeamWriteHandler) run() {
	var err error

	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO inventory.organizational_teams (
	organizational_team_id,
	organizational_team_name,
	organizational_team_ldap_id,
	organizational_team_system)
SELECT $1::uuid, $2::varchar, $3::numeric, $4 WHERE NOT EXISTS (
	SELECT organizational_team_id
	FROM   inventory.organizational_teams
	WHERE  organizational_team_id = $1::uuid
	OR     organizational_team_name = $2::varchar
	OR     organizational_team_ldap_id = $3::numeric);`)
	if err != nil {
		log.Fatal("team/add: ", err)
	}
	defer w.add_stmt.Close()

	if w.upd_stmt, err = w.conn.Prepare(`
UPDATE inventory.organizational_teams
SET    organizational_team_name = $1::varchar,
       organizational_team_ldap_id = $2::numeric,
       organizational_team_system = $3::boolean
WHERE  organizational_team_id = $4::uuid;`); err != nil {
		log.Fatal(`team/update: `, err)
	}
	defer w.upd_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM inventory.organizational_teams
WHERE organizational_team_id = $1;`)
	if err != nil {
		log.Fatal("team/delete: ", err)
	}
	defer w.del_stmt.Close()

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
	notify = msg.Request{Type: `supervisor`, Action: `update_map`,
		Super: &msg.Supervisor{
			Object: `team`,
			Team:   q.Team,
		},
	}

	switch q.action {
	case "add":
		log.Printf("R: team/add for %s", q.Team.Name)
		id := uuid.NewV4()
		res, err = w.add_stmt.Exec(
			id.String(),
			q.Team.Name,
			q.Team.LdapId,
			q.Team.IsSystem,
		)
		q.Team.Id = id.String()
		notify.Super.Action = `add`
	case `update`:
		log.Printf("R: team/update for %s", q.Team.Name)
		res, err = w.upd_stmt.Exec(
			q.Team.Name,
			q.Team.LdapId,
			q.Team.IsSystem,
			q.Team.Id,
		)
		notify.Super.Action = `update`
	case "delete":
		log.Printf("R: team/del for %s", q.Team.Id)
		res, err = w.del_stmt.Exec(
			q.Team.Id,
		)
		notify.Super.Action = `delete`
	default:
		log.Printf("R: unimplemented team/%s", q.action)
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
