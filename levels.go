package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

)

type somaLevelRequest struct {
	action string
	level  somaproto.ProtoLevel
	reply  chan []somaLevelResult
}

type somaLevelResult struct {
	rErr  error
	lErr  error
	level somaproto.ProtoLevel
}

/* Read Access
 *
 */
type somaLevelReadHandler struct {
	input     chan somaLevelRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
}

func (r *somaLevelReadHandler) run() {
	var err error

	r.list_stmt, err = r.conn.Prepare("SELECT level_name FROM soma.notification_levels;")
	if err != nil {
		log.Fatal(err)
	}
	r.show_stmt, err = r.conn.Prepare(`
		SELECT level_name, level_shortname, level_numeric
		FROM soma.notification_levels
		WHERE level_name = $1;`)
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

func (r *somaLevelReadHandler) process(q *somaLevelRequest) {
	var level string
	var short string
	var numeric uint16
	var rows *sql.Rows
	var err error
	result := make([]somaLevelResult, 0)

	switch q.action {
	case "list":
		log.Printf("R: levels/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if err != nil {
			result = append(result, somaLevelResult{
				rErr: err,
				lErr: nil,
			})
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&level)
			if err != nil {
				result = append(result, somaLevelResult{
					rErr: nil,
					lErr: err,
				})
				err = nil
				continue
			}
			result = append(result, somaLevelResult{
				rErr: nil,
				lErr: nil,
				level: somaproto.ProtoLevel{
					Name: level,
				},
			})
		}
	case "show":
		log.Printf("R: levels/show for %s", q.level.Name)
		err = r.show_stmt.QueryRow(q.level.Name).Scan(&level, &short, &numeric)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				result = append(result, somaLevelResult{
					rErr: err,
					lErr: nil,
				})
			}
			q.reply <- result
			return
		}

		result = append(result, somaLevelResult{
			rErr: nil,
			lErr: nil,
			level: somaproto.ProtoLevel{
				Name:      level,
				ShortName: short,
				Numeric:   numeric,
			},
		})
	default:
		result = append(result, somaLevelResult{
			rErr: errors.New("not implemented"),
			lErr: nil,
		})
	}
	q.reply <- result
}

/* Write Access
 *
 */
type somaLevelWriteHandler struct {
	input    chan somaLevelRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
}

func (w *somaLevelWriteHandler) run() {
	var err error

	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.notification_levels (
	level_name,
	level_shortname,
	level_numeric)
SELECT $1, $2, $3 WHERE NOT EXISTS (
	SELECT level_name
	FROM soma.notification_levels
	WHERE level_name = $4
	OR level_shortname = $5
	OR level_numeric = $6);`)
	if err != nil {
		log.Fatal(err)
	}
	defer w.add_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.notification_levels
WHERE level_name = $1;`)
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

func (w *somaLevelWriteHandler) process(q *somaLevelRequest) {
	var res sql.Result
	var err error
	result := make([]somaLevelResult, 0)

	switch q.action {
	case "add":
		log.Printf("R: levels/add for %s", q.level.Name)
		res, err = w.add_stmt.Exec(
			q.level.Name,
			q.level.ShortName,
			q.level.Numeric,
			q.level.Name,
			q.level.ShortName,
			q.level.Numeric,
		)
	case "delete":
		log.Printf("R: levels/del for %s", q.level.Name)
		res, err = w.del_stmt.Exec(
			q.level.Name,
		)
	default:
		log.Printf("R: unimplemented levels/%s", q.action)
		result = append(result, somaLevelResult{
			rErr: errors.New("not implemented"),
		})
		q.reply <- result
		return
	}
	if err != nil {
		result = append(result, somaLevelResult{
			rErr: err,
		})
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		result = append(result, somaLevelResult{
			lErr: errors.New("No rows affected"),
		})
	case rowCnt > 1:
		result = append(result, somaLevelResult{
			lErr: fmt.Errorf("Too many rows affected: %d", rowCnt),
		})
	default:
		result = append(result, somaLevelResult{
			level: q.level,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
