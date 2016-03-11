package main

import (
	"database/sql"
	"fmt"
	"strconv"

)

func CheckUpdateOrInsert(details *somaproto.DeploymentDetails) error {
	var (
		stmt             *sql.Stmt
		err              error
		itemID, lookupID string
		item             ConfigurationItem
	)
	if stmt, err = Eye.conn.Prepare(stmtCheckItemExists); err != nil {
		return err
	}

	lookupID = CalculateLookupId(details.Node.AssetId, details.Metric.Metric)
	fmt.Println(lookupID)

	item = ConfigurationItem{
		Metric:   details.Metric.Metric,
		Interval: details.CheckConfiguration.Interval,
		HostId:   strconv.FormatUint(details.Node.AssetId, 10),
		Metadata: ConfigurationMetaData{
			Monitoring: details.Monitoring.Name,
			Team:       details.Team.Name,
		},
	}
	//Source:
	//Targethost:
	if details.Oncall.Id != "" {
		item.Oncall = fmt.Sprintf("%s (%s)", details.Oncall.Name, details.Oncall.Number)
	}

	err = stmt.QueryRow(details.CheckInstance.InstanceId).Scan(&itemID)
	if err == sql.ErrNoRows {
		// INSERT
	} else if err != nil {
		return err
	}
	// UPDATE
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
