package main

import (
	"database/sql"
	"errors"
	"log"
)

// Message structs
type somaDatacenterRequest struct {
	action     string
	datacenter string
	group      string
	rename     string
	reply      chan []somaDatacenterResult
}

type somaDatacenterResult struct {
	err        error
	datacenter string
}

/*  Read Access
 *
 */
type somaDatacenterReadHandler struct {
	input     chan somaDatacenterRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	grp_list  *sql.Stmt
	grp_show  *sql.Stmt
}

func (r *somaDatacenterReadHandler) run() {
	var err error

	r.list_stmt, err = r.conn.Prepare(`
	SELECT datacenter
	FROM inventory.datacenters;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer r.list_stmt.Close()

	r.show_stmt, err = r.conn.Prepare(`
	SELECT datacenter
	FROM inventory.datacenters
	WHERE datacenter = $1;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer r.show_stmt.Close()

	r.grp_list, err = r.conn.Prepare(`
	SELECT DISTINCT datacenter_group
	FROM soma.datacenter_groups;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer r.grp_list.Close()

	r.grp_show, err = r.conn.Prepare(`
	SELECT DISTINCT datacenter
	FROM soma.datacenter_groups
	WHERE datacenter_group = $1;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer r.grp_show.Close()

	for {
		select {
		case <-r.shutdown:
			break
		case req := <-r.input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

func (r *somaDatacenterReadHandler) process(q *somaDatacenterRequest) {
	var datacenter string
	var rows *sql.Rows
	var err error
	result := make([]somaDatacenterResult, 0)

	switch q.action {
	case "list":
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if err != nil {
			result = append(result, somaDatacenterResult{
				err:        err,
				datacenter: q.datacenter,
			})
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&datacenter)
			if err != nil {
				result = append(result, somaDatacenterResult{
					err:        err,
					datacenter: q.datacenter,
				})
				err = nil
				continue
			}
			result = append(result, somaDatacenterResult{
				err:        nil,
				datacenter: datacenter,
			})
		}
	case "show":
		err = r.show_stmt.QueryRow(q.datacenter).Scan(&datacenter)
		if err != nil {
			result = append(result, somaDatacenterResult{
				err:        err,
				datacenter: q.datacenter,
			})
			q.reply <- result
			return
		}

		result = append(result, somaDatacenterResult{
			err:        nil,
			datacenter: datacenter,
		})
	case "grouplist":
		rows, err = r.grp_list.Query()
		defer rows.Close()
		if err != nil {
			result = append(result, somaDatacenterResult{
				err:        err,
				datacenter: q.datacenter,
			})
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&datacenter)
			if err != nil {
				result = append(result, somaDatacenterResult{
					err:        err,
					datacenter: q.datacenter,
				})
				err = nil
				continue
			}
			result = append(result, somaDatacenterResult{
				err:        nil,
				datacenter: datacenter,
			})
		}
	case "groupshow":
		rows, err = r.grp_show.Query(q.datacenter)
		if err != nil {
			result = append(result, somaDatacenterResult{
				err:        err,
				datacenter: q.datacenter,
			})
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&datacenter)
			if err != nil {
				result = append(result, somaDatacenterResult{
					err:        err,
					datacenter: q.datacenter,
				})
				err = nil
				continue
			}
			result = append(result, somaDatacenterResult{
				err:        nil,
				datacenter: datacenter,
			})
		}
	default:
		result = append(result, somaDatacenterResult{
			err:        errors.New("not implemented"),
			datacenter: "",
		})
	}
	q.reply <- result
}

/*
 * Write Access
 */

type somaDatacenterWriteHandler struct {
	input    chan somaDatacenterRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	ren_stmt *sql.Stmt
	grp_add  *sql.Stmt
	grp_del  *sql.Stmt
}

func (w *somaDatacenterWriteHandler) run() {
	var err error

	w.add_stmt, err = w.conn.Prepare(`
  INSERT INTO inventory.datacenters (datacenter)
  SELECT $1 WHERE NOT EXISTS (
    SELECT datacenter FROM inventory.datacenters WHERE datacenter = $2
  );
  `)
	if err != nil {
		log.Fatal(err)
	}
	defer w.add_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
  DELETE FROM inventory.datacenters
  WHERE datacenter = $1;
  `)
	if err != nil {
		log.Fatal(err)
	}
	defer w.del_stmt.Close()

	w.ren_stmt, err = w.conn.Prepare(`
  UPDATE inventory.datacenters SET datacenter = $1
  WHERE datacenter = $2;
  `)
	if err != nil {
		log.Fatal(err)
	}
	defer w.ren_stmt.Close()

	w.grp_add, err = w.conn.Prepare(`
	INSERT INTO soma.datacenter_groups ( datacenter_group, datacenter )
	SELECT $1, $2 WHERE NOT EXISTS (
		SELECT datacenter FROM soma.datacenter_groups
		WHERE datacenter_group = $3
		AND datacenter = $4
	);
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer w.grp_add.Close()

	w.grp_del, err = w.conn.Prepare(`
	DELETE FROM soma.datacenter_groups
	WHERE datacenter_group = $1
	AND	  datacenter = $2;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer w.grp_del.Close()

	for {
		select {
		case <-w.shutdown:
			break
		case req := <-w.input:
			w.process(&req)
		}
	}
}

func (w *somaDatacenterWriteHandler) process(q *somaDatacenterRequest) {
	var res sql.Result
	var err error

	result := make([]somaDatacenterResult, 0)
	switch q.action {
	case "add":
		res, err = w.add_stmt.Exec(q.datacenter, q.datacenter)
	case "groupadd":
		res, err = w.grp_add.Exec(q.group, q.datacenter, q.group, q.datacenter)
	case "delete":
		res, err = w.del_stmt.Exec(q.datacenter)
	case "groupdel":
		res, err = w.grp_del.Exec(q.group, q.datacenter)
	case "rename":
		res, err = w.ren_stmt.Exec(q.rename, q.datacenter)
	default:
		result = append(result, somaDatacenterResult{
			err:        errors.New("not implemented"),
			datacenter: "",
		})
		q.reply <- result
		return
	}
	if err != nil {
		result = append(result, somaDatacenterResult{
			err:        err,
			datacenter: q.datacenter,
		})
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	if rowCnt == 0 {
		result = append(result, somaDatacenterResult{
			err:        errors.New("No rows affected"),
			datacenter: q.datacenter,
		})
	} else if rowCnt > 1 {
		result = append(result, somaDatacenterResult{
			err:        errors.New("Too many rows affected"),
			datacenter: q.datacenter,
		})
	} else {
		result = append(result, somaDatacenterResult{
			err:        nil,
			datacenter: q.datacenter,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
