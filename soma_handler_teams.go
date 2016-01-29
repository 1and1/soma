package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/satori/go.uuid"
)

type somaTeamRequest struct {
	action string
	team   somaproto.ProtoTeam
	reply  chan []somaTeamResult
}

type somaTeamResult struct {
	rErr error
	lErr error
	team somaproto.ProtoTeam
}

/* Read Access
 *
 */
type somaTeamReadHandler struct {
	input     chan somaTeamRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
}

func (r *somaTeamReadHandler) run() {
	var err error

	log.Println("Prepare: team/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT organizational_team_id,
       organizational_team_name 
FROM   inventory.organizational_teams;`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Prepare: team/show")
	r.show_stmt, err = r.conn.Prepare(`
SELECT organizational_team_id,
       organizational_team_name,
       organizational_team_ldap_id,
       organizational_team_system
FROM   inventory.organizational_teams
WHERE  organizational_team_id = $1;`)
	if err != nil {
		log.Fatal(err)
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
	var teamId, teamName string
	var ldapId int
	var systemFlag bool
	var rows *sql.Rows
	var err error
	result := make([]somaTeamResult, 0)

	switch q.action {
	case "list":
		log.Printf("R: team/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if err != nil {
			result = append(result, somaTeamResult{
				rErr: err,
			})
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&teamId, &teamName)
			if err != nil {
				result = append(result, somaTeamResult{
					lErr: err,
				})
				err = nil
				continue
			}
			result = append(result, somaTeamResult{
				team: somaproto.ProtoTeam{
					Id:   teamId,
					Name: teamName,
				},
			})
		}
	case "show":
		log.Printf("R: team/show for %s", q.team.Id)
		err = r.show_stmt.QueryRow(q.team.Id).Scan(
			&teamId,
			&teamName,
			&ldapId,
			&systemFlag,
		)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				result = append(result, somaTeamResult{
					rErr: err,
				})
				q.reply <- result
				return
			}
		}

		result = append(result, somaTeamResult{
			team: somaproto.ProtoTeam{
				Id:     teamId,
				Name:   teamName,
				Ldap:   strconv.Itoa(ldapId),
				System: systemFlag,
			},
		})
	default:
		result = append(result, somaTeamResult{
			rErr: errors.New("not implemented"),
		})
	}
	q.reply <- result
}

/* Write Access
 *
 */
type somaTeamWriteHandler struct {
	input    chan somaTeamRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	//  upd_stmt *sql.Stmt
	del_stmt *sql.Stmt
}

func (w *somaTeamWriteHandler) run() {
	var err error

	log.Println("Prepare: team/add")
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

	log.Println("Prepare: team/del")
	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM inventory.organizational_teams
WHERE organizational_team_id = $1;`)
	if err != nil {
		log.Fatal(err)
	}
	defer w.del_stmt.Close()

runloop:
	for {
		select {
		case <-w.shutdown:
			break runloop
		case req := <-w.input:
			w.process(&req)
		}
	}
}

func (w *somaTeamWriteHandler) process(q *somaTeamRequest) {
	var res sql.Result
	var err error
	result := make([]somaTeamResult, 0)

	switch q.action {
	case "add":
		log.Printf("R: team/add for %s", q.team.Name)
		id := uuid.NewV4()
		res, err = w.add_stmt.Exec(
			id.String(),
			q.team.Name,
			q.team.Ldap,
			q.team.System,
		)
		q.team.Id = id.String()
	case "delete":
		log.Printf("R: team/del for %s", q.team.Id)
		res, err = w.del_stmt.Exec(
			q.team.Id,
		)
	default:
		log.Printf("R: unimplemented team/%s", q.action)
		result = append(result, somaTeamResult{
			rErr: errors.New("not implemented"),
		})
		q.reply <- result
		return
	}
	if err != nil {
		result = append(result, somaTeamResult{
			rErr: err,
		})
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		result = append(result, somaTeamResult{
			lErr: errors.New("No rows affected"),
		})
	case rowCnt > 1:
		result = append(result, somaTeamResult{
			lErr: fmt.Errorf("Too many rows affected: %d", rowCnt),
		})
	default:
		result = append(result, somaTeamResult{
			team: q.team,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
