package main

import (
	"database/sql"
	"errors"
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
		err = r.show_stmt.QueryRow(q.level.Name).Scan(&level, &short, &numeric)
		if err != nil {
			result = append(result, somaLevelResult{
				rErr: err,
				lErr: nil,
			})
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
