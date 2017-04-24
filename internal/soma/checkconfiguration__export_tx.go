package soma

import (
	"database/sql"
	"strconv"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
)

func exportCheckConfigObjectTX(tx *sql.Tx, objectID string) (
	*[]proto.CheckConfig, error) {

	var (
		err           error
		checkconfigs  []proto.CheckConfig
		checkConfigID string
		checkConfig   *proto.CheckConfig
		txMap         map[string]*sql.Stmt
		instances     []proto.CheckInstanceInfo
		rows          *sql.Rows
	)

	// declare this tx as deferrable read-only
	if _, err = tx.Exec(stmt.ReadOnlyTransaction); err != nil {
		return nil, err
	}

	txMap = make(map[string]*sql.Stmt)
	checkconfigs = make([]proto.CheckConfig, 0)

	for name, statement := range map[string]string{
		`configs`:       stmt.CheckConfigForChecksOnObject,
		`base`:          stmt.CheckConfigShowBase,
		`threshold`:     stmt.CheckConfigShowThreshold,
		`cstrCustom`:    stmt.CheckConfigShowConstrCustom,
		`cstrSystem`:    stmt.CheckConfigShowConstrSystem,
		`cstrNative`:    stmt.CheckConfigShowConstrNative,
		`cstrService`:   stmt.CheckConfigShowConstrService,
		`cstrAttribute`: stmt.CheckConfigShowConstrAttribute,
		`cstrOncall`:    stmt.CheckConfigShowConstrOncall,
		`instance`:      stmt.CheckConfigObjectInstanceInfo,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			return nil, err
		}
	}

	if rows, err = txMap[`configs`].Query(objectID); err != nil {
		return nil, err
	}
	for rows.Next() {
		if err = rows.Scan(
			&checkConfigID,
		); err != nil {
			rows.Close()
			return nil, err
		}

		if checkConfig, err = exportCheckConfig(
			txMap[`base`],
			checkConfigID,
		); err != nil {
			rows.Close()
			return nil, err
		}
		if checkConfig.Thresholds, err = exportCheckConfigThresholds(
			txMap[`threshold`],
			checkConfigID,
		); err != nil {
			rows.Close()
			return nil, err
		}
		if checkConfig.Constraints, err = exportCheckConfigConstraints(
			txMap[`cstrCustom`],
			txMap[`cstrSystem`],
			txMap[`cstrNative`],
			txMap[`cstrService`],
			txMap[`cstrAttribute`],
			txMap[`cstrOncall`],
			checkConfigID,
		); err != nil {
			rows.Close()
			return nil, err
		}
		if instances, err = exportCheckInstancesForObject(
			txMap[`instance`],
			checkConfigID,
			objectID,
		); err != nil {
			rows.Close()
			return nil, err
		}
		if len(instances) > 0 {
			checkConfig.Details = &proto.CheckConfigDetails{
				Instances: instances,
			}
		}
		checkconfigs = append(checkconfigs, *checkConfig)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	tx.Commit()

	return &checkconfigs, nil
}

// expects stmt.CheckConfigShowBase as prepared statement
func exportCheckConfig(prepStmt *sql.Stmt, queryID string) (
	*proto.CheckConfig, error) {

	var (
		checkConfigID, repositoryID, checkConfigName   string
		objectID, objectType, capabilityID, externalID string
		bucketID                                       string
		isActive, hasInheritance, isChildrenOnly       bool
		isEnabled                                      bool
		bucketIDOrNull                                 sql.NullString
		interval                                       int64
	)

	if err := prepStmt.QueryRow(queryID).Scan(
		&checkConfigID,
		&repositoryID,
		&bucketIDOrNull,
		&checkConfigName,
		&objectID,
		&objectType,
		&isActive,
		&hasInheritance,
		&isChildrenOnly,
		&capabilityID,
		&interval,
		&isEnabled,
		&externalID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if bucketIDOrNull.Valid {
		bucketID = bucketIDOrNull.String
	}

	checkConfig := proto.CheckConfig{
		Id:           checkConfigID,
		Name:         checkConfigName,
		Interval:     uint64(interval),
		RepositoryId: repositoryID,
		BucketId:     bucketID,
		CapabilityId: capabilityID,
		ObjectId:     objectID,
		ObjectType:   objectType,
		IsActive:     isActive,
		IsEnabled:    isEnabled,
		Inheritance:  hasInheritance,
		ChildrenOnly: isChildrenOnly,
		ExternalId:   externalID,
	}
	return &checkConfig, nil
}

// expects stmt.CheckConfigShowThreshold as prepared statement
func exportCheckConfigThresholds(prepStmt *sql.Stmt, queryID string) (
	[]proto.CheckConfigThreshold, error) {

	var (
		err                                       error
		rows                                      *sql.Rows
		checkConfigID, predicateSymbol, threshold string
		levelName, levelShortName                 string
		levelNumeric, thresholdValue              int64
		thresholds                                []proto.CheckConfigThreshold
	)

	if rows, err = prepStmt.Query(queryID); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&checkConfigID,
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
	queryID string) ([]proto.CheckConfigConstraint, error) {

	var err error
	var constraints []proto.CheckConfigConstraint
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
				stmtCustom, queryID)
		case `system`:
			cstr, err = exportCheckConfigSystemConstraints(
				stmtSystem, queryID)
		case `native`:
			cstr, err = exportCheckConfigNativeConstraints(
				stmtNative, queryID)
		case `service`:
			cstr, err = exportCheckConfigServiceConstraints(
				stmtService, queryID)
		case `attribute`:
			cstr, err = exportCheckConfigAttributeConstraints(
				stmtAttribute, queryID)
		case `oncall`:
			cstr, err = exportCheckConfigOncallConstraints(
				stmtOncall, queryID)
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
	queryID string) ([]proto.CheckConfigConstraint, error) {

	var (
		checkConfigID, propertyID, repositoryID string
		name, value                             string
		rows                                    *sql.Rows
		err                                     error
		constraints                             []proto.CheckConfigConstraint
	)

	if rows, err = prepStmt.Query(queryID); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&checkConfigID,
			&propertyID,
			&repositoryID,
			&value,
			&name,
		); err != nil {
			rows.Close()
			return nil, err
		}
		cstr := proto.CheckConfigConstraint{
			ConstraintType: `custom`,
			Custom: &proto.PropertyCustom{
				Id:           propertyID,
				RepositoryId: repositoryID,
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
	queryID string) ([]proto.CheckConfigConstraint, error) {

	var (
		checkConfigID, name, value string
		rows                       *sql.Rows
		err                        error
		constraints                []proto.CheckConfigConstraint
	)

	if rows, err = prepStmt.Query(queryID); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&checkConfigID,
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
	queryID string) ([]proto.CheckConfigConstraint, error) {

	var (
		checkConfigID, name, value string
		rows                       *sql.Rows
		err                        error
		constraints                []proto.CheckConfigConstraint
	)

	if rows, err = prepStmt.Query(queryID); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&checkConfigID,
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
	queryID string) ([]proto.CheckConfigConstraint, error) {

	var (
		checkConfigID, name, teamID string
		rows                        *sql.Rows
		err                         error
		constraints                 []proto.CheckConfigConstraint
	)

	if rows, err = prepStmt.Query(queryID); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&checkConfigID,
			&teamID,
			&name,
		); err != nil {
			rows.Close()
			return nil, err
		}
		cstr := proto.CheckConfigConstraint{
			ConstraintType: `service`,
			Service: &proto.PropertyService{
				Name:   name,
				TeamId: teamID,
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
	queryID string) ([]proto.CheckConfigConstraint, error) {

	var (
		checkConfigID, name, value string
		rows                       *sql.Rows
		err                        error
		constraints                []proto.CheckConfigConstraint
	)

	if rows, err = prepStmt.Query(queryID); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&checkConfigID,
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
	queryID string) ([]proto.CheckConfigConstraint, error) {

	var (
		checkConfigID, oncallID, name, number string
		rows                                  *sql.Rows
		err                                   error
		constraints                           []proto.CheckConfigConstraint
	)

	if rows, err = prepStmt.Query(queryID); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&checkConfigID,
			&oncallID,
			&name,
			&number,
		); err != nil {
			rows.Close()
			return nil, err
		}
		cstr := proto.CheckConfigConstraint{
			ConstraintType: `oncall`,
			Oncall: &proto.PropertyOncall{
				Id:     oncallID,
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
	queryID string) ([]proto.CheckInstanceInfo, error) {

	var (
		err                              error
		rows                             *sql.Rows
		instanceID, objectID, objectType string
		currentStatus, nextStatus        string
		instances                        []proto.CheckInstanceInfo
	)

	if rows, err = prepStmt.Query(queryID); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&objectID,
			&objectType,
			&currentStatus,
			&nextStatus,
		); err != nil {
			rows.Close()
			return nil, err
		}
		info := proto.CheckInstanceInfo{
			Id:            instanceID,
			ObjectId:      objectID,
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
	configID, objectID string) ([]proto.CheckInstanceInfo, error) {

	var (
		err                                 error
		rows                                *sql.Rows
		instanceID, resObjectID, objectType string
		currentStatus, nextStatus           string
		instances                           []proto.CheckInstanceInfo
	)

	if rows, err = prepStmt.Query(configID, objectID); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&resObjectID,
			&objectType,
			&currentStatus,
			&nextStatus,
		); err != nil {
			rows.Close()
			return nil, err
		}
		info := proto.CheckInstanceInfo{
			Id:            instanceID,
			ObjectId:      resObjectID,
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
