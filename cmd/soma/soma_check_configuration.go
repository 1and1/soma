package main

import (
	"database/sql"
	"strconv"

	"github.com/1and1/soma/lib/proto"
)

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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
