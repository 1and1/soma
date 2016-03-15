package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

)

type somaValidityRequest struct {
	action   string
	Validity somaproto.Validity
	reply    chan somaResult
}

type somaValidityResult struct {
	ResultError error
	Validity    somaproto.Validity
}

func (a *somaValidityResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Validity = append(r.Validity, somaValidityResult{ResultError: err})
	}
}

func (a *somaValidityResult) SomaAppendResult(r *somaResult) {
	r.Validity = append(r.Validity, *a)
}

/* Read Access
 */
type somaValidityReadHandler struct {
	input     chan somaValidityRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
}

func (r *somaValidityReadHandler) run() {
	var err error

	log.Println("Prepare: validity/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT system_property,
       object_type
FROM   soma.system_property_validity;`)
	if err != nil {
		log.Fatal("validity/list: ", err)
	}
	defer r.list_stmt.Close()

	log.Println("Prepare: validity/show")
	r.show_stmt, err = r.conn.Prepare(`
SELECT system_property,
       object_type,
       inherited
FROM   soma.system_property_validity
WHERE  system_property = $1;`)
	if err != nil {
		log.Fatal("validity/show: ", err)
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

func (r *somaValidityReadHandler) process(q *somaValidityRequest) {
	var (
		property, object string
		inherited        bool
		rows             *sql.Rows
		err              error
		m                map[string]map[string]map[string]bool
	)
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: validity/list")
		rows, err = r.list_stmt.Query()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&property, &object)
			result.Append(err, &somaValidityResult{
				Validity: somaproto.Validity{
					SystemProperty: property,
					ObjectType:     object,
				},
			})
		}
		if err = rows.Err(); err != nil {
			_ = result.SetRequestError(err)
			q.reply <- result
			return
		}
	case "show":
		log.Printf("R: status/show for %s", q.Validity.SystemProperty)
		rows, err = r.show_stmt.Query(q.Validity.SystemProperty)
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		m = make(map[string]map[string]map[string]bool)
		for rows.Next() {
			err = rows.Scan(
				&property,
				&object,
				&inherited,
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
			if m[property] == nil {
				m[property] = make(map[string]map[string]bool)
			}
			if m[property][object] == nil {
				m[property][object] = make(map[string]bool)
			}
			if inherited {
				m[property][object]["inherited"] = true
			} else {
				m[property][object]["direct"] = true
			}
		}
		if err = rows.Err(); err != nil {
			_ = result.SetRequestError(err)
			q.reply <- result
			return
		}
		for p_spec, _ := range m {
			for o_spec, _ := range m[p_spec] {
				result.Append(nil, &somaValidityResult{
					Validity: somaproto.Validity{
						SystemProperty: p_spec,
						ObjectType:     o_spec,
						Direct:         m[p_spec][o_spec]["direct"],
						Inherited:      m[p_spec][o_spec]["inherited"],
					},
				})
			}
		}
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaValidityWriteHandler struct {
	input    chan somaValidityRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
}

func (w *somaValidityWriteHandler) run() {
	var err error

	log.Println("Prepare: validity/add")
	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.system_property_validity (
	system_property,
	object_type,
	inherited)
SELECT $1::varchar,
       $2::varchar,
       $3::boolean
WHERE NOT EXISTS (
	SELECT system_property,
           object_type
	FROM   soma.system_property_validity
	WHERE  system_property = $1::varchar
    AND    object_type = $2::varchar
    AND    inherited = $3::boolean);`)
	if err != nil {
		log.Fatal("validity/add: ", err)
	}
	defer w.add_stmt.Close()

	log.Println("Prepare: validity/delete")
	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.system_property_validity
WHERE       system_property = $1::varchar;`)
	if err != nil {
		log.Fatal("validity/delete: ", err)
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

func (w *somaValidityWriteHandler) process(q *somaValidityRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: validity/add for %s", q.Validity.SystemProperty)
		if q.Validity.Direct {
			res, err = w.add_stmt.Exec(
				q.Validity.SystemProperty,
				q.Validity.ObjectType,
				false,
			)
		}
		if err != nil {
			goto errorout
		}
		if q.Validity.Inherited {
			res, err = w.add_stmt.Exec(
				q.Validity.SystemProperty,
				q.Validity.ObjectType,
				true,
			)
		}
	case "delete":
		log.Printf("R: validity/del for %s", q.Validity.SystemProperty)
		res, err = w.del_stmt.Exec(
			q.Validity.SystemProperty,
		)
	default:
		log.Printf("R: unimplemented validity/%s", q.action)
		result.SetNotImplemented()
		q.reply <- result
		return
	}
errorout:
	if result.SetRequestError(err) {
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		result.Append(errors.New("No rows affected"), &somaValidityResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaValidityResult{})
	default:
		result.Append(nil, &somaValidityResult{
			Validity: q.Validity,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
