package main

import (
	"database/sql"
	"strconv"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
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
	show_instance_info    *sql.Stmt
	appLog                *log.Logger
	reqLog                *log.Logger
	errLog                *log.Logger
}

func (r *somaCheckConfigurationReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.CheckConfigList:                r.list_stmt,
		stmt.CheckConfigShowBase:            r.show_base,
		stmt.CheckConfigShowThreshold:       r.show_threshold,
		stmt.CheckConfigShowConstrCustom:    r.show_constr_custom,
		stmt.CheckConfigShowConstrSystem:    r.show_constr_system,
		stmt.CheckConfigShowConstrNative:    r.show_constr_native,
		stmt.CheckConfigShowConstrService:   r.show_constr_service,
		stmt.CheckConfigShowConstrAttribute: r.show_constr_attribute,
		stmt.CheckConfigShowConstrOncall:    r.show_constr_oncall,
		stmt.CheckConfigInstanceInfo:        r.show_instance_info,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`checkconfig`, err, stmt.Name(statement))
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
	instanceInfo := true

	switch q.action {
	case "list":
		r.appLog.Printf("R: checkconfig/list for %s", q.CheckConfig.RepositoryId)
		rows, err = r.list_stmt.Query(q.CheckConfig.RepositoryId)
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(
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
		r.appLog.Printf("R: checkconfig/show for %s", q.CheckConfig.Id)
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

		if instanceInfo {
			instances := make([]proto.CheckInstanceInfo, 0)
			var (
				rows                                                     *sql.Rows
				instanceId, instObjId, instObjType, instStatus, instNext string
			)

			rows, err = r.show_instance_info.Query(q.CheckConfig.Id)
			if result.SetRequestError(err) {
				q.reply <- result
				return
			}

			for rows.Next() {
				_ = rows.Scan(
					&instanceId,
					&instObjId,
					&instObjType,
					&instStatus,
					&instNext,
				)
				info := proto.CheckInstanceInfo{
					Id:            instanceId,
					ObjectId:      instObjId,
					ObjectType:    instObjType,
					CurrentStatus: instStatus,
					NextStatus:    instNext,
				}
				instances = append(instances, info)
			}
			rows.Close()
			if len(instances) > 0 {
				chkConfig.Details = &proto.CheckConfigDetails{
					Instances: instances,
				}
			}
		}
		result.Append(err, &somaCheckConfigResult{
			CheckConfig: chkConfig,
		})
	default:
		r.errLog.Printf("R: unimplemented checkconfig/%s", q.action)
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaCheckConfigurationReadHandler) shutdownNow() {
	r.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
