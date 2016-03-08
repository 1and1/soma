package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

)

func (tk *treeKeeper) buildDeploymentDetails() {
	var (
		stmt_List, stmt_CheckInstance, stmt_Check, stmt_CheckConfig, stmt_CapMonMetric   *sql.Stmt
		stmt_Threshold, stmt_Pkgs, stmt_Group, stmt_Cluster, stmt_Node, stmt_Team        *sql.Stmt
		stmt_GroupOncall, stmt_ClusterOncall, stmt_NodeOncall                            *sql.Stmt
		stmt_GroupService, stmt_ClusterService, stmt_NodeService                         *sql.Stmt
		stmt_GroupSysProp, stmt_GroupCustProp, stmt_ClusterSysProp, stmt_ClusterCustProp *sql.Stmt
		stmt_NodeSysProp, stmt_NodeCustProp, stmt_DefaultDC, stmt_Update                 *sql.Stmt
		err                                                                              error
		instanceCfgId                                                                    string
		objId, objType                                                                   string
		rows, thresh, pkgs, gSysProps, cSysProps, nSysProps                              *sql.Rows
		gCustProps, cCustProps, nCustProps                                               *sql.Rows
	)

	//
	if stmt_List, err = tk.conn.Prepare(tkStmtDeployDetailsComputeList); err != nil {
		log.Fatal("treekeeper/tkStmtDeployDetailsComputeList: ", err)
	}
	defer stmt_List.Close()
	if stmt_Update, err = tk.conn.Prepare(tkStmtDeployDetailsUpdate); err != nil {
		log.Fatal(err)
	}
	defer stmt_Update.Close()
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
	if stmt_Team, err = tk.conn.Prepare(tkStmtDeployDetailsTeam); err != nil {
		log.Fatal(err)
	}
	defer stmt_Team.Close()

	if stmt_GroupOncall, err = tk.conn.Prepare(tkStmtDeployDetailsGroupOncall); err != nil {
		log.Fatal(err)
	}
	defer stmt_GroupOncall.Close()
	if stmt_GroupService, err = tk.conn.Prepare(tkStmtDeployDetailsGroupService); err != nil {
		log.Fatal(err)
	}
	defer stmt_GroupService.Close()
	if stmt_ClusterOncall, err = tk.conn.Prepare(tkStmtDeployDetailsClusterOncall); err != nil {
		log.Fatal(err)
	}
	defer stmt_ClusterOncall.Close()
	if stmt_ClusterService, err = tk.conn.Prepare(tkStmtDeployDetailsClusterService); err != nil {
		log.Fatal(err)
	}
	defer stmt_ClusterService.Close()
	if stmt_NodeOncall, err = tk.conn.Prepare(tkStmtDeployDetailsNodeOncall); err != nil {
		log.Fatal(err)
	}
	defer stmt_NodeOncall.Close()
	if stmt_NodeService, err = tk.conn.Prepare(tkStmtDeployDetailsNodeService); err != nil {
		log.Fatal(err)
	}
	defer stmt_NodeService.Close()
	if stmt_GroupSysProp, err = tk.conn.Prepare(tkStmtDeployDetailsGroupSysProp); err != nil {
		log.Fatal(err)
	}
	defer stmt_GroupSysProp.Close()
	if stmt_GroupCustProp, err = tk.conn.Prepare(tkStmtDeployDetailsGroupCustProp); err != nil {
		log.Fatal(err)
	}
	defer stmt_GroupCustProp.Close()
	if stmt_ClusterSysProp, err = tk.conn.Prepare(tkStmtDeployDetailClusterSysProp); err != nil {
		log.Fatal(err)
	}
	defer stmt_ClusterSysProp.Close()
	if stmt_ClusterCustProp, err = tk.conn.Prepare(tkStmtDeployDetailClusterCustProp); err != nil {
		log.Fatal(err)
	}
	defer stmt_ClusterCustProp.Close()
	if stmt_NodeSysProp, err = tk.conn.Prepare(tkStmtDeployDetailNodeSysProp); err != nil {
		log.Fatal(err)
	}
	defer stmt_NodeSysProp.Close()
	if stmt_NodeCustProp, err = tk.conn.Prepare(tkStmtDeployDetailNodeCustProp); err != nil {
		log.Fatal(err)
	}
	defer stmt_NodeCustProp.Close()
	if stmt_DefaultDC, err = tk.conn.Prepare(tkStmtDeployDetailDefaultDatacenter); err != nil {
		log.Fatal(err)
	}
	defer stmt_DefaultDC.Close()

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

		// XXX TODO
		//detail.CheckConfiguration.Constraints = []somaproto.CheckConfigurationConstraint{}
		detail.CheckConfiguration.Constraints = nil

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
		detail.Oncall = &somaproto.ProtoOncall{}
		detail.Service = &somaproto.TreePropertyService{}
		switch objType {
		case "group":
			// fetch the group object
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
			// fetch team information
			detail.Team = &somaproto.ProtoTeam{
				Id: detail.Group.TeamId,
			}
			// fetch oncall information if the property is set,
			// otherwise cleanup detail.Oncall
			err = stmt_GroupOncall.QueryRow(detail.Group.Id, detail.View).Scan(
				&detail.Oncall.Id,
				&detail.Oncall.Name,
				&detail.Oncall.Number,
			)
			if err != nil {
				detail.Oncall = nil
			}
			// fetch service name, and attributes if applicable
			if detail.CheckInstance.InstanceService != "" {
				err = stmt_GroupService.QueryRow(
					detail.CheckInstance.InstanceService,
					detail.View,
				).Scan(
					&detail.Service.Name,
					&detail.Service.TeamId,
				)
				if err != nil {
					detail.Service = nil
				} else {
					detail.Service.Attributes = []somaproto.TreeServiceAttribute{}
					fm := map[string]string{}
					_ = json.Unmarshal([]byte(detail.CheckInstance.InstanceServiceConfig), &fm)
					for k, v := range fm {
						a := somaproto.TreeServiceAttribute{
							Attribute: k,
							Value:     v,
						}
						detail.Service.Attributes = append(detail.Service.Attributes, a)
					}
				}
			}
			// fetch system properties
			detail.Properties = &[]somaproto.TreePropertySystem{}
			gSysProps, _ = stmt_GroupSysProp.Query(detail.Group.Id, detail.View)
			defer gSysProps.Close()

			for gSysProps.Next() {
				prop := somaproto.TreePropertySystem{}
				err = gSysProps.Scan(
					&prop.Name,
					&prop.Value,
				)
				*detail.Properties = append(*detail.Properties, prop)
				if prop.Name == "group_datacenter" {
					detail.Datacenter = prop.Value
				}
			}
			if len(*detail.Properties) == 0 {
				detail.Properties = nil
			}
			// fetch custom properties
			detail.CustomProperties = &[]somaproto.TreePropertyCustom{}
			gCustProps, _ = stmt_GroupCustProp.Query(detail.Group.Id, detail.View)
			defer gCustProps.Close()

			for gCustProps.Next() {
				prop := somaproto.TreePropertyCustom{}
				gCustProps.Scan(
					&prop.CustomId,
					&prop.Name,
					&prop.Value,
				)
				*detail.CustomProperties = append(*detail.CustomProperties, prop)
			}
			if len(*detail.CustomProperties) == 0 {
				detail.CustomProperties = nil
			}
		case "cluster":
			// fetch the cluster object
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
			// fetch team information
			detail.Team = &somaproto.ProtoTeam{
				Id: detail.Cluster.TeamId,
			}
			// fetch oncall information if the property is set,
			// otherwise cleanup detail.Oncall
			err = stmt_ClusterOncall.QueryRow(detail.Cluster.Id, detail.View).Scan(
				&detail.Oncall.Id,
				&detail.Oncall.Name,
				&detail.Oncall.Number,
			)
			if err != nil {
				detail.Oncall = nil
			}
			// fetch the service name, and attributes if applicable
			if detail.CheckInstance.InstanceService != "" {
				err = stmt_ClusterService.QueryRow(
					detail.CheckInstance.InstanceService,
					detail.View,
				).Scan(
					&detail.Service.Name,
					&detail.Service.TeamId,
				)
				if err != nil {
					detail.Service = nil
				} else {
					detail.Service.Attributes = []somaproto.TreeServiceAttribute{}
					fm := map[string]string{}
					_ = json.Unmarshal([]byte(detail.CheckInstance.InstanceServiceConfig), &fm)
					for k, v := range fm {
						a := somaproto.TreeServiceAttribute{
							Attribute: k,
							Value:     v,
						}
						detail.Service.Attributes = append(detail.Service.Attributes, a)
					}
				}
			}
			// fetch system properties
			detail.Properties = &[]somaproto.TreePropertySystem{}
			cSysProps, _ = stmt_ClusterSysProp.Query(detail.Cluster.Id, detail.View)
			defer cSysProps.Close()

			for cSysProps.Next() {
				prop := somaproto.TreePropertySystem{}
				err = cSysProps.Scan(
					&prop.Name,
					&prop.Value,
				)
				*detail.Properties = append(*detail.Properties, prop)
				if prop.Name == "cluster_datacenter" {
					detail.Datacenter = prop.Value
				}
			}
			if len(*detail.Properties) == 0 {
				detail.Properties = nil
			}
			// fetch custom properties
			detail.CustomProperties = &[]somaproto.TreePropertyCustom{}
			cCustProps, _ = stmt_ClusterCustProp.Query(detail.Cluster.Id, detail.View)
			defer cCustProps.Close()

			for cCustProps.Next() {
				prop := somaproto.TreePropertyCustom{}
				cCustProps.Scan(
					&prop.CustomId,
					&prop.Name,
					&prop.Value,
				)
				*detail.CustomProperties = append(*detail.CustomProperties, prop)
			}
			if len(*detail.CustomProperties) == 0 {
				detail.CustomProperties = nil
			}
		case "node":
			// fetch the node object
			detail.Node = &somaproto.ProtoNode{
				Id: objId,
			}
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
			// fetch team information
			detail.Team = &somaproto.ProtoTeam{
				Id: detail.Node.Team,
			}
			// fetch oncall information if the property is set,
			// otherwise cleanup detail.Oncall
			detail.Oncall = &somaproto.ProtoOncall{}
			err = stmt_NodeOncall.QueryRow(detail.Node.Id, detail.View).Scan(
				&detail.Oncall.Id,
				&detail.Oncall.Name,
				&detail.Oncall.Number,
			)
			if err != nil {
				detail.Oncall = nil
			}
			// fetch the service name, and attributes if applicable
			if detail.CheckInstance.InstanceService != "" {
				err = stmt_NodeService.QueryRow(
					detail.CheckInstance.InstanceService,
					detail.View,
				).Scan(
					&detail.Service.Name,
					&detail.Service.TeamId,
				)
				if err != nil {
					detail.Service = nil
				} else {
					detail.Service.Attributes = []somaproto.TreeServiceAttribute{}
					fm := map[string]string{}
					_ = json.Unmarshal([]byte(detail.CheckInstance.InstanceServiceConfig), &fm)
					for k, v := range fm {
						a := somaproto.TreeServiceAttribute{
							Attribute: k,
							Value:     v,
						}
						detail.Service.Attributes = append(detail.Service.Attributes, a)
					}
				}
			}
			// fetch system properties
			detail.Properties = &[]somaproto.TreePropertySystem{}
			nSysProps, _ = stmt_NodeSysProp.Query(detail.Node.Id, detail.View)
			defer nSysProps.Close()

			for nSysProps.Next() {
				prop := somaproto.TreePropertySystem{}
				err = nSysProps.Scan(
					&prop.Name,
					&prop.Value,
				)
				*detail.Properties = append(*detail.Properties, prop)
			}
			if len(*detail.Properties) == 0 {
				detail.Properties = nil
			}
			// fetch custom properties
			detail.CustomProperties = &[]somaproto.TreePropertyCustom{}
			nCustProps, _ = stmt_NodeCustProp.Query(detail.Node.Id, detail.View)
			defer nCustProps.Close()

			for nCustProps.Next() {
				prop := somaproto.TreePropertyCustom{}
				gCustProps.Scan(
					&prop.CustomId,
					&prop.Name,
					&prop.Value,
				)
				*detail.CustomProperties = append(*detail.CustomProperties, prop)
			}
			if len(*detail.CustomProperties) == 0 {
				detail.CustomProperties = nil
			}
		}

		stmt_Team.QueryRow(detail.Team.Id).Scan(
			&detail.Team.Name,
			&detail.Team.Ldap,
		)

		// if no datacenter information was gathered, use the default DC
		if detail.Datacenter == "" {
			stmt_DefaultDC.QueryRow().Scan(&detail.Datacenter)
		}

		// build JSON of DeploymentDetails
		var detailJSON []byte
		if detailJSON, err = json.Marshal(&detail); err != nil {
			log.Fatal("Failed to JSON marshal deployment details: ", err)
		}
		if _, err = stmt_Update.Exec(
			detailJSON,
			detail.Monitoring.Id,
			detail.CheckInstance.InstanceConfigId,
		); err != nil {
			log.Fatal("Failed to save DeploymentDetails.JSON: ", err)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
