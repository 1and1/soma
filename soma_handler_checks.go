package main

import (
	"database/sql"
	"log"
	"strconv"

)

type somaCheckConfigRequest struct {
	action      string
	CheckConfig proto.CheckConfig
	reply       chan somaResult
}

type somaCheckConfigResult struct {
	ResultError error
	CheckConfig proto.CheckConfig
}

func (a *somaCheckConfigResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.CheckConfigs = append(r.CheckConfigs, somaCheckConfigResult{ResultError: err})
	}
}

func (a *somaCheckConfigResult) SomaAppendResult(r *somaResult) {
	r.CheckConfigs = append(r.CheckConfigs, *a)
}

/* Read Access
 */
type somaCheckConfigurationReadHandler struct {
	input                 chan somaCheckConfigRequest
	shutdown              chan bool
	conn                  *sql.DB
	list_stmt             *sql.Stmt
	show_base             *sql.Stmt
	show_threshold        *sql.Stmt
	show_constr_custom    *sql.Stmt
	show_constr_system    *sql.Stmt
	show_constr_native    *sql.Stmt
	show_constr_service   *sql.Stmt
	show_constr_attribute *sql.Stmt
	show_constr_oncall    *sql.Stmt
}

func (r *somaCheckConfigurationReadHandler) run() {
	var err error

	if r.list_stmt, err = r.conn.Prepare(stmtCheckConfigList); err != nil {
		log.Fatal("checkconfig/list: ", err)
	}
	defer r.list_stmt.Close()

	if r.show_base, err = r.conn.Prepare(stmtCheckConfigShowBase); err != nil {
		log.Fatal("checkconfig/show-base: ", err)
	}
	defer r.show_base.Close()

	if r.show_threshold, err = r.conn.Prepare(stmtCheckConfigShowThreshold); err != nil {
		log.Fatal("checkconfig/show-threshold: ", err)
	}
	defer r.show_threshold.Close()

	if r.show_constr_custom, err = r.conn.Prepare(stmtCheckConfigShowConstrCustom); err != nil {
		log.Fatal("checkconfig/show-constraint-custom: ", err)
	}
	defer r.show_constr_custom.Close()

	if r.show_constr_system, err = r.conn.Prepare(stmtCheckConfigShowConstrSystem); err != nil {
		log.Fatal("checkconfig/show-constraint-system: ", err)
	}
	defer r.show_constr_system.Close()

	if r.show_constr_native, err = r.conn.Prepare(stmtCheckConfigShowConstrNative); err != nil {
		log.Fatal("checkconfig/show-constraint-native: ", err)
	}
	defer r.show_constr_native.Close()

	if r.show_constr_service, err = r.conn.Prepare(stmtCheckConfigShowConstrService); err != nil {
		log.Fatal("checkconfig/show-constraint-service: ", err)
	}
	defer r.show_constr_service.Close()

	if r.show_constr_attribute, err = r.conn.Prepare(stmtCheckConfigShowConstrAttribute); err != nil {
		log.Fatal("checkconfig/show-constraint-attribute: ", err)
	}
	defer r.show_constr_attribute.Close()

	if r.show_constr_oncall, err = r.conn.Prepare(stmtCheckConfigShowConstrOncall); err != nil {
		log.Fatal("checkconfig/show-constraint-oncall: ", err)
	}
	defer r.show_constr_oncall.Close()

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

func (r *somaCheckConfigurationReadHandler) process(q *somaCheckConfigRequest) {
	var (
		configId, repoId, configName, configObjId, configObjType, capabilityId, extId, buck string
		configActive, inheritance, childrenOnly, enabled                                    bool
		interval                                                                            int64
		rows                                                                                *sql.Rows
		err                                                                                 error
		bucketId                                                                            sql.NullString
	)
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: checkconfig/list for %s", q.CheckConfig.RepositoryId)
		rows, err = r.list_stmt.Query(q.CheckConfig.RepositoryId)
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(
				&configId,
				&repoId,
				&bucketId,
				&configName,
			)
			if bucketId.Valid {
				buck = bucketId.String
			}
			result.Append(err, &somaCheckConfigResult{
				CheckConfig: proto.CheckConfig{
					Id:           configId,
					RepositoryId: repoId,
					BucketId:     buck,
					Name:         configName,
				},
			})
		}
	case "show":
		log.Printf("R: checkconfig/list for %s", q.CheckConfig.Id)
		if err = r.show_base.QueryRow(q.CheckConfig.Id).Scan(
			&configId,
			&repoId,
			&bucketId,
			&configName,
			&configObjId,
			&configObjType,
			&configActive,
			&inheritance,
			&childrenOnly,
			&capabilityId,
			&interval,
			&enabled,
			&extId,
		); err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		if bucketId.Valid {
			buck = bucketId.String
		}
		chkConfig := proto.CheckConfig{
			Id:           configId,
			Name:         configName,
			Interval:     uint64(interval),
			RepositoryId: repoId,
			BucketId:     buck,
			CapabilityId: capabilityId,
			ObjectId:     configObjId,
			ObjectType:   configObjType,
			IsActive:     configActive,
			IsEnabled:    enabled,
			Inheritance:  inheritance,
			ChildrenOnly: childrenOnly,
			ExternalId:   extId,
		}
		chkConfig.Thresholds = make([]proto.CheckConfigThreshold, 0)
		if rows, err = r.show_threshold.Query(q.CheckConfig.Id); err != nil {
			if err != sql.ErrNoRows {
				if result.SetRequestError(err) {
					q.reply <- result
					return
				}
			}
		}
		// err can be nil or sql.ErrNoRows, everything else already
		// error'ed out
		if err != sql.ErrNoRows {
			// rows is *sql.Rows, so if err != nil then rows == nilptr.
			// this makes rows.Close() a nilptr dereference if we
			// ignored sql.ErrNoRows
			defer rows.Close()

			var (
				predicate, threshold, levelName, levelShort string
				numeric, treshVal                           int64
			)

			for rows.Next() {
				err = rows.Scan(
					&configId,
					&predicate,
					&threshold,
					&levelName,
					&levelShort,
					&numeric,
				)
				treshVal, _ = strconv.ParseInt(threshold, 10, 64)
				thr := proto.CheckConfigThreshold{
					Predicate: proto.Predicate{
						Symbol: predicate,
					},
					Level: proto.Level{
						Name:      levelName,
						ShortName: levelShort,
						Numeric:   uint16(numeric),
					},
					Value: treshVal,
				}

				chkConfig.Thresholds = append(chkConfig.Thresholds, thr)
			}
		}

		chkConfig.Constraints = make([]proto.CheckConfigConstraint, 0)
		for _, tp := range []string{"custom", "system", "native", "service", "attribute", "oncall"} {
			var (
				err                                           error
				configId, propertyId, repoId, property, value string
				rows                                          *sql.Rows
			)

			switch tp {
			case "custom":
				rows, err = r.show_constr_custom.Query(q.CheckConfig.Id)
				if result.SetRequestError(err) {
					q.reply <- result
					return
				}
				defer rows.Close()

				for rows.Next() {
					_ = rows.Scan(
						&configId,
						&propertyId,
						&repoId,
						&value,
						&property,
					)
					constr := proto.CheckConfigConstraint{
						ConstraintType: tp,
						Custom: &proto.PropertyCustom{
							Id:           propertyId,
							RepositoryId: repoId,
							Name:         property,
							Value:        value,
						},
					}
					chkConfig.Constraints = append(chkConfig.Constraints, constr)
				}
			case "system":
				rows, err = r.show_constr_system.Query(q.CheckConfig.Id)
				if result.SetRequestError(err) {
					q.reply <- result
					return
				}
				defer rows.Close()

				for rows.Next() {
					_ = rows.Scan(
						&configId,
						&property,
						&value,
					)
					constr := proto.CheckConfigConstraint{
						ConstraintType: tp,
						System: &proto.PropertySystem{
							Name:  property,
							Value: value,
						},
					}
					chkConfig.Constraints = append(chkConfig.Constraints, constr)
				}
			case "native":
				rows, err = r.show_constr_native.Query(q.CheckConfig.Id)
				if result.SetRequestError(err) {
					q.reply <- result
					return
				}
				defer rows.Close()

				for rows.Next() {
					_ = rows.Scan(
						&configId,
						&property,
						&value,
					)
					constr := proto.CheckConfigConstraint{
						ConstraintType: tp,
						Native: &proto.PropertyNative{
							Name:  property,
							Value: value,
						},
					}
					chkConfig.Constraints = append(chkConfig.Constraints, constr)
				}
			case "service":
				rows, err = r.show_constr_service.Query(q.CheckConfig.Id)
				if result.SetRequestError(err) {
					q.reply <- result
					return
				}
				defer rows.Close()

				for rows.Next() {
					_ = rows.Scan(
						&configId,
						&propertyId,
						&property,
					)
					constr := proto.CheckConfigConstraint{
						ConstraintType: tp,
						Service: &proto.PropertyService{
							Name:   property,
							TeamId: propertyId,
						},
					}
					chkConfig.Constraints = append(chkConfig.Constraints, constr)
				}
			case "attribute":
				rows, err = r.show_constr_attribute.Query(q.CheckConfig.Id)
				if result.SetRequestError(err) {
					q.reply <- result
					return
				}
				defer rows.Close()

				for rows.Next() {
					_ = rows.Scan(
						&configId,
						&property,
						&value,
					)
					constr := proto.CheckConfigConstraint{
						ConstraintType: tp,
						Attribute: &proto.ServiceAttribute{
							Name:  property,
							Value: value,
						},
					}
					chkConfig.Constraints = append(chkConfig.Constraints, constr)
				}
			case "oncall":
				rows, err = r.show_constr_oncall.Query(q.CheckConfig.Id)
				if result.SetRequestError(err) {
					q.reply <- result
					return
				}
				defer rows.Close()

				for rows.Next() {
					_ = rows.Scan(
						&configId,
						&propertyId,
						&property,
						&value,
					)
					constr := proto.CheckConfigConstraint{
						ConstraintType: tp,
						Oncall: &proto.PropertyOncall{
							Id:     propertyId,
							Name:   property,
							Number: value,
						},
					}
					chkConfig.Constraints = append(chkConfig.Constraints, constr)
				}
			} // switch tp
		}
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
