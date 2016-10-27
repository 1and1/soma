package main

import (
	"database/sql"
	"strconv"

	"github.com/1and1/soma/lib/proto"
)

// expects stmt.CheckConfigShowBase as prepared statement
func exportCheckConfig(prepStmt *sql.Stmt, queryId string) (
	*proto.CheckConfig, error) {

	var (
		checkConfigId, repositoryId, checkConfigName   string
		objectId, objectType, capabilityId, externalId string
		bucketId                                       string
		isActive, hasInheritance, isChildrenOnly       bool
		isEnabled                                      bool
		bucketIdOrNull                                 sql.NullString
		interval                                       int64
	)

	if err := prepStmt.QueryRow(queryId).Scan(
		&checkConfigId,
		&repositoryId,
		&bucketIdOrNull,
		&checkConfigName,
		&objectId,
		&objectType,
		&isActive,
		&hasInheritance,
		&isChildrenOnly,
		&capabilityId,
		&interval,
		&isEnabled,
		&externalId,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if bucketIdOrNull.Valid {
		bucketId = bucketIdOrNull.String
	}

	checkConfig := proto.CheckConfig{
		Id:           checkConfigId,
		Name:         checkConfigName,
		Interval:     uint64(interval),
		RepositoryId: repositoryId,
		BucketId:     bucketId,
		CapabilityId: capabilityId,
		ObjectId:     objectId,
		ObjectType:   objectType,
		IsActive:     isActive,
		IsEnabled:    isEnabled,
		Inheritance:  hasInheritance,
		ChildrenOnly: isChildrenOnly,
		ExternalId:   externalId,
	}
	return &checkConfig, nil
}

// expects stmt.CheckConfigShowThreshold as prepared statement
func exportCheckConfigThresholds(prepStmt *sql.Stmt, queryId string) (
	[]proto.CheckConfigThreshold, error) {

	var (
		err                                       error
		rows                                      *sql.Rows
		checkConfigId, predicateSymbol, threshold string
		levelName, levelShortName                 string
		levelNumeric, thresholdValue              int64
	)
	thresholds := make([]proto.CheckConfigThreshold, 0)

	if rows, err = prepStmt.Query(queryId); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&checkConfigId,
			&predicateSymbol,
			&threshold,
			&levelName,
			&levelShortName,
			&levelNumeric,
		); err != nil {
			rows.Close()
			return nil, err
		}
		thresholdValue, _ = strconv.ParseInt(threshold, 10, 64)

		thr := proto.CheckConfigThreshold{
			Predicate: proto.Predicate{
				Symbol: predicateSymbol,
			},
			Level: proto.Level{
				Name:      levelName,
				ShortName: levelShortName,
				Numeric:   uint16(levelNumeric),
			},
			Value: thresholdValue,
		}
		thresholds = append(thresholds, thr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return thresholds, nil
}

// expects in that order:
// - stmt.CheckConfigShowConstrCustom
// - stmt.CheckConfigShowConstrSystem
// - stmt.CheckConfigShowConstrNative
// - stmt.CheckConfigShowConstrService
// - stmt.CheckConfigShowConstrAttribute
// - stmt.CheckConfigShowConstrOncall
func exportCheckConfigConstraints(stmtCustom, stmtSystem,
	stmtNative, stmtService, stmtAttribute, stmtOncall *sql.Stmt,
	queryId string) ([]proto.CheckConfigConstraint, error) {

	var err error
	constraints := make([]proto.CheckConfigConstraint, 0)
	for _, cstrType := range []string{
		`custom`,
		`system`,
		`native`,
		`service`,
		`attribute`,
		`oncall`,
	} {
		cstr := []proto.CheckConfigConstraint{}
		switch cstrType {
		case `custom`:
			cstr, err = exportCheckConfigCustomConstraints(
				stmtCustom, queryId)
		case `system`:
			cstr, err = exportCheckConfigSystemConstraints(
				stmtSystem, queryId)
		case `native`:
			cstr, err = exportCheckConfigNativeConstraints(
				stmtNative, queryId)
		case `service`:
			cstr, err = exportCheckConfigServiceConstraints(
				stmtService, queryId)
		case `attribute`:
			cstr, err = exportCheckConfigAttributeConstraints(
				stmtAttribute, queryId)
		case `oncall`:
			cstr, err = exportCheckConfigOncallConstraints(
				stmtOncall, queryId)
		}
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, cstr...)
	}
	return constraints, nil
}

// expects stmt.CheckConfigShowConstrCustom as prepared statement
func exportCheckConfigCustomConstraints(prepStmt *sql.Stmt,
	queryId string) ([]proto.CheckConfigConstraint, error) {

	var (
		checkConfigId, propertyId, repositoryId string
		name, value                             string
		rows                                    *sql.Rows
		err                                     error
	)

	constraints := make([]proto.CheckConfigConstraint, 0)

	if rows, err = prepStmt.Query(queryId); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&checkConfigId,
			&propertyId,
			&repositoryId,
			&value,
			&name,
		); err != nil {
			rows.Close()
			return nil, err
		}
		cstr := proto.CheckConfigConstraint{
			ConstraintType: `custom`,
			Custom: &proto.PropertyCustom{
				Id:           propertyId,
				RepositoryId: repositoryId,
				Name:         name,
				Value:        value,
			},
		}
		constraints = append(constraints, cstr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return constraints, nil
}

// expects stmt.CheckConfigShowConstrSystem as prepared statement
func exportCheckConfigSystemConstraints(prepStmt *sql.Stmt,
	queryId string) ([]proto.CheckConfigConstraint, error) {

	var (
		checkConfigId, name, value string
		rows                       *sql.Rows
		err                        error
	)

	constraints := make([]proto.CheckConfigConstraint, 0)

	if rows, err = prepStmt.Query(queryId); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&checkConfigId,
			&name,
			&value,
		); err != nil {
			rows.Close()
			return nil, err
		}
		cstr := proto.CheckConfigConstraint{
			ConstraintType: `system`,
			System: &proto.PropertySystem{
				Name:  name,
				Value: value,
			},
		}
		constraints = append(constraints, cstr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return constraints, nil
}

// expects stmt.CheckConfigShowConstrNative as prepared statement
func exportCheckConfigNativeConstraints(prepStmt *sql.Stmt,
	queryId string) ([]proto.CheckConfigConstraint, error) {

	var (
		checkConfigId, name, value string
		rows                       *sql.Rows
		err                        error
	)

	constraints := make([]proto.CheckConfigConstraint, 0)

	if rows, err = prepStmt.Query(queryId); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&checkConfigId,
			&name,
			&value,
		); err != nil {
			rows.Close()
			return nil, err
		}
		cstr := proto.CheckConfigConstraint{
			ConstraintType: `native`,
			Native: &proto.PropertyNative{
				Name:  name,
				Value: value,
			},
		}
		constraints = append(constraints, cstr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return constraints, nil
}

// expects stmt.CheckConfigShowConstrService as prepared statement
func exportCheckConfigServiceConstraints(prepStmt *sql.Stmt,
	queryId string) ([]proto.CheckConfigConstraint, error) {

	var (
		checkConfigId, name, teamId string
		rows                        *sql.Rows
		err                         error
	)

	constraints := make([]proto.CheckConfigConstraint, 0)

	if rows, err = prepStmt.Query(queryId); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&checkConfigId,
			&teamId,
			&name,
		); err != nil {
			rows.Close()
			return nil, err
		}
		cstr := proto.CheckConfigConstraint{
			ConstraintType: `service`,
			Service: &proto.PropertyService{
				Name:   name,
				TeamId: teamId,
			},
		}
		constraints = append(constraints, cstr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return constraints, nil
}

// expects stmt.CheckConfigShowConstrAttribute as prepared statement
func exportCheckConfigAttributeConstraints(prepStmt *sql.Stmt,
	queryId string) ([]proto.CheckConfigConstraint, error) {

	var (
		checkConfigId, name, value string
		rows                       *sql.Rows
		err                        error
	)

	constraints := make([]proto.CheckConfigConstraint, 0)

	if rows, err = prepStmt.Query(queryId); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&checkConfigId,
			&name,
			&value,
		); err != nil {
			rows.Close()
			return nil, err
		}
		cstr := proto.CheckConfigConstraint{
			ConstraintType: `attribute`,
			Attribute: &proto.ServiceAttribute{
				Name:  name,
				Value: value,
			},
		}
		constraints = append(constraints, cstr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return constraints, nil
}

// expects stmt.CheckConfigShowConstrOncall as prepared statement
func exportCheckConfigOncallConstraints(prepStmt *sql.Stmt,
	queryId string) ([]proto.CheckConfigConstraint, error) {

	var (
		checkConfigId, oncallId, name, number string
		rows                                  *sql.Rows
		err                                   error
	)

	constraints := make([]proto.CheckConfigConstraint, 0)

	if rows, err = prepStmt.Query(queryId); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&checkConfigId,
			&oncallId,
			&name,
			&number,
		); err != nil {
			rows.Close()
			return nil, err
		}
		cstr := proto.CheckConfigConstraint{
			ConstraintType: `oncall`,
			Oncall: &proto.PropertyOncall{
				Id:     oncallId,
				Name:   name,
				Number: number,
			},
		}
		constraints = append(constraints, cstr)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return constraints, nil
}

// expects stmt.CheckConfigInstanceInfo as prepared statement
func exportCheckInstancesForConfig(prepStmt *sql.Stmt,
	queryId string) ([]proto.CheckInstanceInfo, error) {

	var (
		err                              error
		rows                             *sql.Rows
		instanceId, objectId, objectType string
		currentStatus, nextStatus        string
	)

	instances := make([]proto.CheckInstanceInfo, 0)

	if rows, err = prepStmt.Query(queryId); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceId,
			&objectId,
			&objectType,
			&currentStatus,
			&nextStatus,
		); err != nil {
			rows.Close()
			return nil, err
		}
		info := proto.CheckInstanceInfo{
			Id:            instanceId,
			ObjectId:      objectId,
			ObjectType:    objectType,
			CurrentStatus: currentStatus,
			NextStatus:    nextStatus,
		}
		instances = append(instances, info)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return instances, nil
}

// expects stmt.CheckConfigObjectInstanceInfo as prepared statement
func exportCheckInstancesForObject(prepStmt *sql.Stmt,
	configId, objectId string) ([]proto.CheckInstanceInfo, error) {

	var (
		err                                 error
		rows                                *sql.Rows
		instanceId, resObjectId, objectType string
		currentStatus, nextStatus           string
	)

	instances := make([]proto.CheckInstanceInfo, 0)

	if rows, err = prepStmt.Query(configId, objectId); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceId,
			&resObjectId,
			&objectType,
			&currentStatus,
			&nextStatus,
		); err != nil {
			rows.Close()
			return nil, err
		}
		info := proto.CheckInstanceInfo{
			Id:            instanceId,
			ObjectId:      resObjectId,
			ObjectType:    objectType,
			CurrentStatus: currentStatus,
			NextStatus:    nextStatus,
		}
		instances = append(instances, info)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return instances, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
