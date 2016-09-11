package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/satori/go.uuid"

)

type somaPropertyRequest struct {
	action  string
	prType  string
	System  proto.PropertySystem
	Native  proto.PropertyNative
	Service proto.PropertyService
	Custom  proto.PropertyCustom
	reply   chan somaResult
}

type somaPropertyResult struct {
	ResultError error
	prType      string
	System      proto.PropertySystem
	Native      proto.PropertyNative
	Service     proto.PropertyService
	Custom      proto.PropertyCustom
}

func (a *somaPropertyResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Properties = append(r.Properties, somaPropertyResult{ResultError: err})
	}
}

func (a *somaPropertyResult) SomaAppendResult(r *somaResult) {
	r.Properties = append(r.Properties, *a)
}

/* Read Access
 */
type somaPropertyReadHandler struct {
	input         chan somaPropertyRequest
	shutdown      chan bool
	conn          *sql.DB
	list_sys_stmt *sql.Stmt
	list_srv_stmt *sql.Stmt
	list_nat_stmt *sql.Stmt
	list_tpl_stmt *sql.Stmt
	list_cst_stmt *sql.Stmt
	show_sys_stmt *sql.Stmt
	show_srv_stmt *sql.Stmt
	show_nat_stmt *sql.Stmt
	show_tpl_stmt *sql.Stmt
	show_cst_stmt *sql.Stmt
}

func (r *somaPropertyReadHandler) run() {
	var err error

	r.list_sys_stmt, err = r.conn.Prepare(`
SELECT system_property
FROM   soma.system_properties;`)
	if err != nil {
		log.Fatal("property/list-system: ", err)
	}
	defer r.list_sys_stmt.Close()

	r.list_srv_stmt, err = r.conn.Prepare(`
SELECT service_property,
       organizational_team_id
FROM   soma.team_service_properties
WHERE  organizational_team_id = $1::uuid;`)
	if err != nil {
		log.Fatal("property/list-service: ", err)
	}
	defer r.list_srv_stmt.Close()

	r.list_nat_stmt, err = r.conn.Prepare(`
SELECT native_property
FROM   soma.native_properties;`)
	if err != nil {
		log.Fatal("property/list-native: ", err)
	}
	defer r.list_nat_stmt.Close()

	r.list_tpl_stmt, err = r.conn.Prepare(`
SELECT service_property
FROM   soma.service_properties;`)
	if err != nil {
		log.Fatal("property/list-service-template: ", err)
	}
	defer r.list_tpl_stmt.Close()

	r.list_cst_stmt, err = r.conn.Prepare(`
SELECT custom_property_id,
       repository_id,
	   custom_property
FROM   soma.custom_properties
WHERE  repository_id = $1::uuid;`)
	if err != nil {
		log.Fatal("property/list-custom: ", err)
	}
	defer r.list_cst_stmt.Close()

	r.show_sys_stmt, err = r.conn.Prepare(`
SELECT system_property
FROM   soma.system_properties
WHERE  system_property = $1::varchar;`)
	if err != nil {
		log.Fatal("property/show-system: ", err)
	}
	defer r.show_sys_stmt.Close()

	r.show_nat_stmt, err = r.conn.Prepare(`
SELECT native_property
FROM   soma.native_properties
WHERE  native_property = $1::varchar;`)
	if err != nil {
		log.Fatal("property/show-native: ", err)
	}
	defer r.show_nat_stmt.Close()

	r.show_cst_stmt, err = r.conn.Prepare(`
SELECT custom_property_id,
       repository_id,
	   custom_property
FROM   soma.custom_properties
WHERE  custom_property_id = $1::uuid
AND    repository_id = $2::uuid;`)
	if err != nil {
		log.Fatal("property/show-custom: ", err)
	}
	defer r.show_cst_stmt.Close()

	r.show_srv_stmt, err = r.conn.Prepare(`
SELECT tsp.service_property,
       tsp.organizational_team_id,
	   tspv.service_property_attribute,
	   tspv.value
FROM   soma.team_service_properties tsp
JOIN   soma.team_service_property_values tspv
ON     tsp.service_property = tspv.service_property
WHERE  tsp.service_property = $1::varchar
AND    tsp.organizational_team_id = $2::uuid;`)
	if err != nil {
		log.Fatal("property/show-service")
	}
	defer r.show_srv_stmt.Close()

	r.show_tpl_stmt, err = r.conn.Prepare(`
SELECT sp.service_property,
	   spv.service_property_attribute,
	   spv.value
FROM   soma.service_properties sp
JOIN   soma.service_property_values spv
ON     sp.service_property = spv.service_property
WHERE  sp.service_property = $1::varchar;`)
	if err != nil {
		log.Fatal("property/show-service-template")
	}
	defer r.show_tpl_stmt.Close()

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

func (r *somaPropertyReadHandler) process(q *somaPropertyRequest) {
	var (
		property, team, repository, id, attribute, value string
		rows                                             *sql.Rows
		err                                              error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		switch q.prType {
		case "system":
			log.Printf("R: property/list-system")
			rows, err = r.list_sys_stmt.Query()
		case "native":
			log.Printf("R: property/list-native")
			rows, err = r.list_nat_stmt.Query()
		case "custom":
			log.Printf("R: property/list-custom")
			rows, err = r.list_cst_stmt.Query(q.Custom.RepositoryId)
		case "service":
			log.Printf("R: property/list-service")
			rows, err = r.list_srv_stmt.Query(q.Service.TeamId)
		case "template":
			log.Printf("R: property/list-service-template")
			rows, err = r.list_tpl_stmt.Query()
		}
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			switch q.prType {
			case "system":
				err := rows.Scan(&property)
				result.Append(err, &somaPropertyResult{
					prType: q.prType,
					System: proto.PropertySystem{
						Name: property,
					},
				})
			case "service":
				err := rows.Scan(&property, &team)
				result.Append(err, &somaPropertyResult{
					prType: q.prType,
					Service: proto.PropertyService{
						Name:   property,
						TeamId: team,
					},
				})
			case "native":
				err := rows.Scan(&property)
				result.Append(err, &somaPropertyResult{
					prType: q.prType,
					Native: proto.PropertyNative{
						Name: property,
					},
				})
			case "template":
				err := rows.Scan(&property)
				result.Append(err, &somaPropertyResult{
					prType: q.prType,
					Service: proto.PropertyService{
						Name: property,
					},
				})
			case "custom":
				err := rows.Scan(&id, &repository, &property)
				result.Append(err, &somaPropertyResult{
					prType: q.prType,
					Custom: proto.PropertyCustom{
						Id:           id,
						RepositoryId: repository,
						Name:         property,
					},
				})
			}
		}
	case "show":
		switch q.prType {
		case "system":
			log.Printf("R: property/show-system for %s", q.System.Name)
			err = r.show_sys_stmt.QueryRow(q.System.Name).Scan(
				&property,
			)
		case "native":
			log.Printf("R: property/show-native for %s", q.Native.Name)
			err = r.show_nat_stmt.QueryRow(q.Native.Name).Scan(
				&property,
			)
		case "custom":
			log.Printf("R: property/show-custom for %s", q.Custom.Id)
			err = r.show_cst_stmt.QueryRow(
				q.Custom.Id,
				q.Custom.RepositoryId,
			).Scan(
				&id,
				&repository,
				&property,
			)
		case "service":
			log.Printf("R: property/show-service for %s/%s", q.Service.TeamId, q.Service.Name)
			rows, err = r.show_srv_stmt.Query(
				q.Service.Name,
				q.Service.TeamId,
			)
			defer rows.Close()
		case "template":
			log.Printf("R: property/show-service-template for %s", q.Service.Name)
			rows, err = r.show_tpl_stmt.Query(
				q.Service.Name,
			)
			defer rows.Close()
		}
		if err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		switch q.prType {
		case "system":
			result.Append(err, &somaPropertyResult{
				prType: q.prType,
				System: proto.PropertySystem{
					Name: property,
				},
			})
		case "native":
			result.Append(err, &somaPropertyResult{
				prType: q.prType,
				Native: proto.PropertyNative{
					Name: property,
				},
			})
		case "custom":
			result.Append(err, &somaPropertyResult{
				prType: q.prType,
				Custom: proto.PropertyCustom{
					Id:           id,
					RepositoryId: repository,
					Name:         property,
				},
			})
		case "service":
			propTempl := proto.PropertyService{}
			var fErr error
			for rows.Next() {
				err := rows.Scan(
					&property,
					&team,
					&attribute,
					&value,
				)

				if err == nil {
					propTempl.Name = property
					propTempl.TeamId = team
					propTempl.Attributes = append(propTempl.Attributes, proto.ServiceAttribute{
						Name:  attribute,
						Value: value,
					})
				} else {
					fErr = err
				}
			}
			result.Append(fErr, &somaPropertyResult{
				prType:  q.prType,
				Service: propTempl,
			})
		case "template":
			propTempl := proto.PropertyService{}
			propTempl.Attributes = make([]proto.ServiceAttribute, 0)
			var fErr error
			for rows.Next() {
				err := rows.Scan(
					&property,
					&attribute,
					&value,
				)

				if err == nil {
					propTempl.Name = property
					propTempl.Attributes = append(propTempl.Attributes, proto.ServiceAttribute{
						Name:  attribute,
						Value: value,
					})
				} else {
					fErr = err
				}
			}
			result.Append(fErr, &somaPropertyResult{
				prType:  q.prType,
				Service: propTempl,
			})
		}
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaPropertyWriteHandler struct {
	input             chan somaPropertyRequest
	shutdown          chan bool
	conn              *sql.DB
	add_sys_stmt      *sql.Stmt
	add_nat_stmt      *sql.Stmt
	add_cst_stmt      *sql.Stmt
	add_srv_stmt      *sql.Stmt
	add_tpl_stmt      *sql.Stmt
	add_srv_attr_stmt *sql.Stmt
	add_tpl_attr_stmt *sql.Stmt
	del_sys_stmt      *sql.Stmt
	del_nat_stmt      *sql.Stmt
	del_cst_stmt      *sql.Stmt
	del_srv_stmt      *sql.Stmt
	del_tpl_stmt      *sql.Stmt
	del_srv_attr_stmt *sql.Stmt
	del_tpl_attr_stmt *sql.Stmt
}

func (w *somaPropertyWriteHandler) run() {
	var err error

	w.add_sys_stmt, err = w.conn.Prepare(`
INSERT INTO soma.system_properties (
	system_property)
SELECT $1::varchar WHERE NOT EXISTS (
	SELECT system_property
	FROM   soma.system_properties
	WHERE  system_property = $1::varchar);`)
	if err != nil {
		log.Fatal("property/add-system: ", err)
	}
	defer w.add_sys_stmt.Close()

	w.add_nat_stmt, err = w.conn.Prepare(`
INSERT INTO soma.native_properties (
	native_property)
SELECT $1::varchar WHERE NOT EXISTS (
	SELECT native_property
	FROM   soma.native_properties
	WHERE  native_property = $1::varchar);`)
	if err != nil {
		log.Fatal("property/add-native: ", err)
	}
	defer w.add_nat_stmt.Close()

	w.add_cst_stmt, err = w.conn.Prepare(`
INSERT INTO soma.custom_properties (
	custom_property_id,
	repository_id,
	custom_property)
SELECT $1::uuid, $2::uuid, $3::varchar WHERE NOT EXISTS (
	SELECT custom_property
	FROM   soma.custom_properties
	WHERE  custom_property = $3::varchar
    AND    repository_id = $2::uuid);`)
	if err != nil {
		log.Fatal("property/add-system: ", err)
	}
	defer w.add_cst_stmt.Close()

	w.add_srv_stmt, err = w.conn.Prepare(`
INSERT INTO soma.team_service_properties (
	organizational_team_id,
	service_property)
SELECT $1::uuid, $2::varchar WHERE NOT EXISTS (
	SELECT service_property
	FROM   soma.team_service_properties
	WHERE  organizational_team_id = $1::uuid
	AND    service_property = $2::varchar);`)
	if err != nil {
		log.Fatal("property/add-service: ", err)
	}
	defer w.add_srv_stmt.Close()

	w.add_srv_attr_stmt, err = w.conn.Prepare(`
INSERT INTO soma.team_service_property_values (
	organizational_team_id,
	service_property,
	service_property_attribute,
	value)
SELECT $1::uuid, $2::varchar, $3::varchar, $4::varchar;`)
	if err != nil {
		log.Fatal("property/add-service-attribute: ", err)
	}
	defer w.add_srv_attr_stmt.Close()

	w.add_tpl_stmt, err = w.conn.Prepare(`
INSERT INTO soma.service_properties (
	service_property)
SELECT $1::varchar WHERE NOT EXISTS (
	SELECT service_property
	FROM   soma.service_properties
	WHERE  service_property = $1::varchar);`)
	if err != nil {
		log.Fatal("property/add-service-template: ", err)
	}
	defer w.add_tpl_stmt.Close()

	w.add_tpl_attr_stmt, err = w.conn.Prepare(`
INSERT INTO soma.service_property_values (
	service_property,
	service_property_attribute,
	value)
SELECT $1::varchar, $2::varchar, $3::varchar;`)
	if err != nil {
		log.Fatal("property/add-service-template-attribute: ", err)
	}
	defer w.add_tpl_attr_stmt.Close()

	w.del_sys_stmt, err = w.conn.Prepare(`
DELETE FROM soma.system_properties
WHERE  system_property = $1::varchar;`)
	if err != nil {
		log.Fatal("property/delete-system: ", err)
	}
	defer w.del_sys_stmt.Close()

	w.del_nat_stmt, err = w.conn.Prepare(`
DELETE FROM soma.native_properties
WHERE  native_property = $1::varchar;`)
	if err != nil {
		log.Fatal("property/delete-native: ", err)
	}
	defer w.del_nat_stmt.Close()

	w.del_cst_stmt, err = w.conn.Prepare(`
DELETE FROM soma.custom_properties
WHERE  repository_id = $1::uuid
AND    custom_property_id = $2::uuid;`)
	if err != nil {
		log.Fatal("property/delete-custom: ", err)
	}
	defer w.del_cst_stmt.Close()

	w.del_srv_stmt, err = w.conn.Prepare(`
DELETE FROM soma.team_service_properties
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`)
	if err != nil {
		log.Fatal("property/delete-service: ", err)
	}
	defer w.del_srv_stmt.Close()

	w.del_srv_attr_stmt, err = w.conn.Prepare(`
DELETE FROM soma.team_service_property_values
WHERE  organizational_team_id = $1::uuid
AND    service_property = $2::varchar;`)
	if err != nil {
		log.Fatal("property/delete-service-attributes: ", err)
	}
	defer w.del_srv_attr_stmt.Close()

	w.del_tpl_stmt, err = w.conn.Prepare(`
DELETE FROM soma.service_properties
WHERE  service_property = $1::varchar;`)
	if err != nil {
		log.Fatal("property/delete-service-template: ", err)
	}
	defer w.del_tpl_stmt.Close()

	w.del_tpl_attr_stmt, err = w.conn.Prepare(`
DELETE FROM soma.service_property_values
WHERE  service_property = $1::varchar;`)
	if err != nil {
		log.Fatal("property/delete-service-template-attributes: ", err)
	}
	defer w.del_tpl_attr_stmt.Close()

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

func (w *somaPropertyWriteHandler) process(q *somaPropertyRequest) {
	var (
		res    sql.Result
		err    error
		tx     *sql.Tx
		attr   proto.ServiceAttribute
		rowCnt int64
	)
	result := somaResult{}

	switch q.action {
	case "add":
		switch q.prType {
		case "system":
			log.Printf("R: property/add-system for %s", q.System.Name)
			res, err = w.add_sys_stmt.Exec(
				q.System.Name,
			)
			rowCnt, _ = res.RowsAffected()
		case "native":
			log.Printf("R: property/add-native for %s", q.Native.Name)
			res, err = w.add_nat_stmt.Exec(
				q.Native.Name,
			)
			rowCnt, _ = res.RowsAffected()
		case "custom":
			q.Custom.Id = uuid.NewV4().String()
			log.Printf("R: property/add-custom for %s", q.Custom.Name)
			res, err = w.add_cst_stmt.Exec(
				q.Custom.Id,
				q.Custom.RepositoryId,
				q.Custom.Name,
			)
			rowCnt, _ = res.RowsAffected()
		case "service":
			log.Printf("R: property/add-service for %s/%s", q.Service.TeamId, q.Service.Name)
			tx, err = w.conn.Begin()
			if err != nil {
				goto bailout
			}
			defer tx.Rollback()

			res, err = tx.Stmt(w.add_srv_stmt).Exec(
				q.Service.TeamId,
				q.Service.Name,
			)
			if err != nil {
				goto bailout
			}
			rowCnt, _ = res.RowsAffected()
			if rowCnt == 0 {
				goto bailout
			}

			for _, attr = range q.Service.Attributes {
				res, err = tx.Stmt(w.add_srv_attr_stmt).Exec(
					q.Service.TeamId,
					q.Service.Name,
					attr.Name,
					attr.Value,
				)
				if err != nil {
					break
				}
			}
			if err != nil {
				goto bailout
			}

			err = tx.Commit()
		case "template":
			log.Printf("R: property/add-service-template for %s", q.Service.Name)
			tx, err = w.conn.Begin()
			if err != nil {
				goto bailout
			}
			defer tx.Rollback()

			res, err = tx.Stmt(w.add_tpl_stmt).Exec(
				q.Service.Name,
			)
			if err != nil {
				goto bailout
			}
			rowCnt, _ = res.RowsAffected()
			if rowCnt == 0 {
				goto bailout
			}

			for _, attr = range q.Service.Attributes {
				res, err = tx.Stmt(w.add_tpl_attr_stmt).Exec(
					q.Service.Name,
					attr.Name,
					attr.Value,
				)
				if err != nil {
					break
				}
			}
			if err != nil {
				goto bailout
			}

			err = tx.Commit()
		}
	case "delete":
		switch q.prType {
		case "system":
			log.Printf("R: property/delete-system for %s", q.System.Name)
			res, err = w.del_sys_stmt.Exec(
				q.System.Name,
			)
			rowCnt, _ = res.RowsAffected()
		case "native":
			log.Printf("R: property/delete-native for %s", q.Native.Name)
			res, err = w.del_nat_stmt.Exec(
				q.Native.Name,
			)
			rowCnt, _ = res.RowsAffected()
		case "custom":
			log.Printf("R: property/delete-custom for %s", q.Custom.Id)
			res, err = w.del_cst_stmt.Exec(
				q.Custom.RepositoryId,
				q.Custom.Id,
			)
			rowCnt, _ = res.RowsAffected()
		case "service":
			log.Printf("R: property/delete-service for %s/%s", q.Service.TeamId, q.Service.Name)
			tx, err = w.conn.Begin()
			if err != nil {
				goto bailout
			}
			defer tx.Rollback()

			res, err = tx.Stmt(w.del_srv_attr_stmt).Exec(
				q.Service.TeamId,
				q.Service.Name,
			)
			if err != nil {
				goto bailout
			}

			res, err = tx.Stmt(w.del_srv_stmt).Exec(
				q.Service.TeamId,
				q.Service.Name,
			)
			if err != nil {
				goto bailout
			}

			rowCnt, _ = res.RowsAffected()
			err = tx.Commit()
		case "template":
			log.Printf("R: property/delete-service-template for %s", q.Service.Name)
			tx, err = w.conn.Begin()
			if err != nil {
				goto bailout
			}
			defer tx.Rollback()

			res, err = tx.Stmt(w.del_tpl_attr_stmt).Exec(
				q.Service.Name,
			)
			if err != nil {
				goto bailout
			}

			res, err = tx.Stmt(w.del_tpl_stmt).Exec(
				q.Service.Name,
			)
			if err != nil {
				goto bailout
			}

			rowCnt, _ = res.RowsAffected()
			err = tx.Commit()
		} // switch q.prType
	default:
		log.Printf("R: unimplemented property/%s", q.action)
		result.SetNotImplemented()
		q.reply <- result
		return
	}

bailout:
	if result.SetRequestError(err) {
		q.reply <- result
		return
	}

	switch {
	case rowCnt == 0:
		result.Append(errors.New("No rows affected"), &somaPropertyResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaPropertyResult{})
	default:
		switch q.prType {
		case "system":
			result.Append(nil, &somaPropertyResult{
				prType: q.prType,
				System: q.System,
			})
		case "native":
			result.Append(nil, &somaPropertyResult{
				prType: q.prType,
				Native: q.Native,
			})
		case "custom":
			result.Append(nil, &somaPropertyResult{
				prType: q.prType,
				Custom: q.Custom,
			})
		case "service":
			result.Append(nil, &somaPropertyResult{
				prType:  q.prType,
				Service: q.Service,
			})
		case "template":
			result.Append(nil, &somaPropertyResult{
				prType:  q.prType,
				Service: q.Service,
			})
		}
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
