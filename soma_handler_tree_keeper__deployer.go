package main

import (
	"database/sql"
	"fmt"
	"log"

)

func (tk *treeKeeper) buildDeploymentDetails() {
	var (
		stmt_List, stmt_CheckInstance, stmt_Check, stmt_CheckConfig, stmt_CapMonMetric *sql.Stmt
		stmt_Threshold, stmt_Pkgs, stmt_Group, stmt_Cluster, stmt_Node                 *sql.Stmt
		rows, thresh, pkgs                                                             *sql.Rows
		err                                                                            error
		instanceCfgId                                                                  string
		objId, objType                                                                 string
	)

	//
	if stmt_List, err = tk.conn.Prepare(tkStmtDeployDetailsComputeList); err != nil {
		log.Fatal("treekeeper/tkStmtDeployDetailsComputeList: ", err)
	}
	defer stmt_List.Close()
	if stmt_CheckInstance, err = tk.conn.Prepare(tkStmtDeployDetailsCheckInstance); err != nil {
		log.Fatal("treekeeper/tkStmtDeployDetailsCheckInstance: ", err)
	}
	defer stmt_CheckInstance.Close()

	if stmt_Check, err = tk.conn.Prepare(tkStmtDeployDetailsCheck); err != nil {
		log.Fatal(err)
	}
	defer stmt_Check.Close()
	if stmt_CheckConfig, err = tk.conn.Prepare(tkStmtDeployDetailsCheckConfig); err != nil {
		log.Fatal(err)
	}
	defer stmt_CheckConfig.Close()
	if stmt_Threshold, err = tk.conn.Prepare(tkStmtDeployDetailsCheckConfigThreshold); err != nil {
		log.Fatal(err)
	}
	defer stmt_Threshold.Close()
	if stmt_CapMonMetric, err = tk.conn.Prepare(tkStmtDeployDetailsCapabilityMonitoringMetric); err != nil {
		log.Fatal(err)
	}
	defer stmt_CapMonMetric.Close()
	if stmt_Pkgs, err = tk.conn.Prepare(tkStmtDeployDetailsProviders); err != nil {
		log.Fatal(err)
	}
	defer stmt_Pkgs.Close()
	if stmt_Group, err = tk.conn.Prepare(tkStmtDeployDetailsGroup); err != nil {
		log.Fatal(err)
	}
	defer stmt_Group.Close()
	if stmt_Cluster, err = tk.conn.Prepare(tkStmtDeployDetailsCluster); err != nil {
		log.Fatal(err)
	}
	defer stmt_Cluster.Close()
	if stmt_Node, err = tk.conn.Prepare(tkStmtDeployDetailsNode); err != nil {
		log.Fatal(err)
	}
	defer stmt_Node.Close()

	//
	if rows, err = stmt_List.Query(); err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		detail := somaproto.DeploymentDetails{}

		err = rows.Scan(
			&instanceCfgId,
		)

		//
		detail.CheckInstance = &somaproto.TreeCheckInstance{
			InstanceConfigId: instanceCfgId,
		}
		stmt_CheckInstance.QueryRow(instanceCfgId).Scan(
			&detail.CheckInstance.Version,
			&detail.CheckInstance.InstanceId,
			&detail.CheckInstance.ConstraintHash,
			&detail.CheckInstance.ConstraintValHash,
			&detail.CheckInstance.InstanceService,
			&detail.CheckInstance.InstanceSvcCfgHash,
			&detail.CheckInstance.InstanceServiceConfig,
			&detail.CheckInstance.CheckId,
			&detail.CheckInstance.ConfigId,
		)

		//
		detail.Check = &somaproto.TreeCheck{
			CheckId: detail.CheckInstance.CheckId,
		}
		stmt_Check.QueryRow(detail.CheckInstance.CheckId).Scan(
			&detail.Check.RepositoryId,
			&detail.Check.SourceCheckId,
			&detail.Check.SourceType,
			&detail.Check.InheritedFrom,
			&detail.Check.CapabilityId,
			&objId,
			&objType,
			&detail.Check.Inheritance,
			&detail.Check.ChildrenOnly,
		)
		detail.ObjectType = objType
		if detail.Check.InheritedFrom != objId {
			detail.Check.IsInherited = true
		}
		detail.Check.CheckConfigId = detail.CheckInstance.ConfigId

		//
		detail.CheckConfiguration = &somaproto.CheckConfiguration{
			Id:           detail.Check.CheckConfigId,
			RepositoryId: detail.Check.RepositoryId,
			BucketId:     detail.Check.BucketId,
			CapabilityId: detail.Check.CapabilityId,
			ObjectId:     objId,
			ObjectType:   objType,
			Inheritance:  detail.Check.Inheritance,
			ChildrenOnly: detail.Check.ChildrenOnly,
		}
		stmt_CheckConfig.QueryRow(detail.Check.CheckConfigId).Scan(
			&detail.CheckConfiguration.Name,
			&detail.CheckConfiguration.Interval,
			&detail.CheckConfiguration.IsActive,
			&detail.CheckConfiguration.IsEnabled,
			&detail.CheckConfiguration.ExternalId,
		)

		//
		detail.CheckConfiguration.Thresholds = []somaproto.CheckConfigurationThreshold{}
		thresh, _ = stmt_Threshold.Query(detail.CheckConfiguration.Id)
		defer thresh.Close()

		for thresh.Next() {
			thr := somaproto.CheckConfigurationThreshold{
				Predicate: somaproto.ProtoPredicate{},
				Level:     somaproto.ProtoLevel{},
			}

			err = thresh.Scan(
				&thr.Predicate.Predicate,
				&thr.Value,
				&thr.Level.Name,
				&thr.Level.ShortName,
				&thr.Level.Numeric,
			)
			detail.CheckConfiguration.Thresholds = append(detail.CheckConfiguration.Thresholds, thr)
		}

		detail.CheckConfiguration.Constraints = []somaproto.CheckConfigurationConstraint{}
		// XXX TODO

		//
		detail.Capability = &somaproto.ProtoCapability{
			Id: detail.Check.CapabilityId,
		}
		detail.Monitoring = &somaproto.ProtoMonitoring{}
		detail.Metric = &somaproto.ProtoMetric{}
		detail.Unit = &somaproto.ProtoUnit{}
		stmt_CapMonMetric.QueryRow(detail.Capability.Id).Scan(
			&detail.Capability.Metric,
			&detail.Capability.Monitoring,
			&detail.Capability.View,
			&detail.Capability.Thresholds,
			&detail.Monitoring.Name,
			&detail.Monitoring.Mode,
			&detail.Monitoring.Contact,
			&detail.Monitoring.Team,
			&detail.Monitoring.Callback,
			&detail.Metric.Unit,
			&detail.Metric.Description,
			&detail.Unit.Name,
		)
		detail.Unit.Unit = detail.Metric.Unit
		detail.Metric.Metric = detail.Capability.Metric
		detail.Monitoring.Id = detail.Capability.Monitoring
		detail.Capability.Name = fmt.Sprintf("%s.%s.%s",
			detail.Monitoring.Name,
			detail.Capability.View,
			detail.Metric.Metric,
		)
		detail.View = detail.Capability.View

		//
		detail.Metric.Packages = &[]somaproto.ProtoMetricProviderPackage{}
		pkgs, _ = stmt_Pkgs.Query(detail.Metric.Metric)
		defer pkgs.Close()

		for pkgs.Next() {
			pkg := somaproto.ProtoMetricProviderPackage{}

			err = pkgs.Scan(
				&pkg.Provider,
				&pkg.Package,
			)
			*detail.Metric.Packages = append(*detail.Metric.Packages, pkg)
		}

		//
		switch objType {
		case "group":
			detail.Group = &somaproto.ProtoGroup{
				Id: objId,
			}
			stmt_Group.QueryRow(objId).Scan(
				&detail.Group.BucketId,
				&detail.Group.Name,
				&detail.Group.ObjectState,
				&detail.Group.TeamId,
				&detail.Bucket,
				&detail.Environment,
				&detail.Repository,
			)
		case "cluster":
			detail.Cluster = &somaproto.ProtoCluster{
				Id: objId,
			}
			stmt_Cluster.QueryRow(objId).Scan(
				&detail.Cluster.Name,
				&detail.Cluster.BucketId,
				&detail.Cluster.ObjectState,
				&detail.Cluster.TeamId,
				&detail.Bucket,
				&detail.Environment,
				&detail.Repository,
			)
		case "node":
			detail.Node = &somaproto.ProtoNode{
				Id: objId,
			}
			detail.Server = &somaproto.ProtoServer{}
			stmt_Node.QueryRow(objId).Scan(
				&detail.Node.AssetId,
				&detail.Node.Name,
				&detail.Node.Team,
				&detail.Node.Server,
				&detail.Node.State,
				&detail.Node.IsOnline,
				&detail.Node.IsDeleted,
				&detail.Bucket,
				&detail.Environment,
				&detail.Repository,
				&detail.Server.AssetId,
				&detail.Server.Datacenter,
				&detail.Server.Location,
				&detail.Server.Name,
				&detail.Server.IsOnline,
				&detail.Server.IsDeleted,
			)
			detail.Server.Id = detail.Node.Server
			detail.Datacenter = detail.Server.Datacenter
		}
		/* XXX TODO
		* detail.Datacenter for group/cluster
		* detail.Team (-> group/cluster/node teamid)
		* detail.Oncall (-> oncallprop)
		* detail.Service (-> detail.CheckInstance)
		* detail.Properties
		* detail.CustomProperties
		 */
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
