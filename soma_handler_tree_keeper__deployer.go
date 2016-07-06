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
		callback                                                                         sql.NullString
	)

	//
	if stmt_List, err = tk.conn.Prepare(tkStmtDeployDetailsComputeList); err != nil {
		log.Fatal("treekeeper/tkStmtDeployDetailsComputeList: ", err)
	}
	defer stmt_List.Close()
	if stmt_Update, err = tk.conn.Prepare(tkStmtDeployDetailsUpdate); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsUpdate")
		log.Fatal(err)
	}
	defer stmt_Update.Close()
	if stmt_CheckInstance, err = tk.conn.Prepare(tkStmtDeployDetailsCheckInstance); err != nil {
		log.Fatal("treekeeper/tkStmtDeployDetailsCheckInstance: ", err)
	}
	defer stmt_CheckInstance.Close()

	if stmt_Check, err = tk.conn.Prepare(tkStmtDeployDetailsCheck); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsCheck")
		log.Fatal(err)
	}
	defer stmt_Check.Close()
	if stmt_CheckConfig, err = tk.conn.Prepare(tkStmtDeployDetailsCheckConfig); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsCheckConfig")
		log.Fatal(err)
	}
	defer stmt_CheckConfig.Close()
	if stmt_Threshold, err = tk.conn.Prepare(tkStmtDeployDetailsCheckConfigThreshold); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsCheckConfigThreshold")
		log.Fatal(err)
	}
	defer stmt_Threshold.Close()
	if stmt_CapMonMetric, err = tk.conn.Prepare(tkStmtDeployDetailsCapabilityMonitoringMetric); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsCapabilityMonitoringMetric")
		log.Fatal(err)
	}
	defer stmt_CapMonMetric.Close()
	if stmt_Pkgs, err = tk.conn.Prepare(tkStmtDeployDetailsProviders); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsProviders")
		log.Fatal(err)
	}
	defer stmt_Pkgs.Close()
	if stmt_Group, err = tk.conn.Prepare(tkStmtDeployDetailsGroup); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsGroup")
		log.Fatal(err)
	}
	defer stmt_Group.Close()
	if stmt_Cluster, err = tk.conn.Prepare(tkStmtDeployDetailsCluster); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsCluster")
		log.Fatal(err)
	}
	defer stmt_Cluster.Close()
	if stmt_Node, err = tk.conn.Prepare(tkStmtDeployDetailsNode); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsNode")
		log.Fatal(err)
	}
	defer stmt_Node.Close()
	if stmt_Team, err = tk.conn.Prepare(tkStmtDeployDetailsTeam); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsTeam")
		log.Fatal(err)
	}
	defer stmt_Team.Close()

	if stmt_GroupOncall, err = tk.conn.Prepare(tkStmtDeployDetailsGroupOncall); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsGroupOncall")
		log.Fatal(err)
	}
	defer stmt_GroupOncall.Close()
	if stmt_GroupService, err = tk.conn.Prepare(tkStmtDeployDetailsGroupService); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsGroupService")
		log.Fatal(err)
	}
	defer stmt_GroupService.Close()
	if stmt_ClusterOncall, err = tk.conn.Prepare(tkStmtDeployDetailsClusterOncall); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsClusterOncall")
		log.Fatal(err)
	}
	defer stmt_ClusterOncall.Close()
	if stmt_ClusterService, err = tk.conn.Prepare(tkStmtDeployDetailsClusterService); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsClusterService")
		log.Fatal(err)
	}
	defer stmt_ClusterService.Close()
	if stmt_NodeOncall, err = tk.conn.Prepare(tkStmtDeployDetailsNodeOncall); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsNodeOncall")
		log.Fatal(err)
	}
	defer stmt_NodeOncall.Close()
	if stmt_NodeService, err = tk.conn.Prepare(tkStmtDeployDetailsNodeService); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsNodeService")
		log.Fatal(err)
	}
	defer stmt_NodeService.Close()
	if stmt_GroupSysProp, err = tk.conn.Prepare(tkStmtDeployDetailsGroupSysProp); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsGroupSysProp")
		log.Fatal(err)
	}
	defer stmt_GroupSysProp.Close()
	if stmt_GroupCustProp, err = tk.conn.Prepare(tkStmtDeployDetailsGroupCustProp); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailsGroupCustProp")
		log.Fatal(err)
	}
	defer stmt_GroupCustProp.Close()
	if stmt_ClusterSysProp, err = tk.conn.Prepare(tkStmtDeployDetailClusterSysProp); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailClusterSysProp")
		log.Fatal(err)
	}
	defer stmt_ClusterSysProp.Close()
	if stmt_ClusterCustProp, err = tk.conn.Prepare(tkStmtDeployDetailClusterCustProp); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailClusterCustProp")
		log.Fatal(err)
	}
	defer stmt_ClusterCustProp.Close()
	if stmt_NodeSysProp, err = tk.conn.Prepare(tkStmtDeployDetailNodeSysProp); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailNodeSysProp")
		log.Fatal(err)
	}
	defer stmt_NodeSysProp.Close()
	if stmt_NodeCustProp, err = tk.conn.Prepare(tkStmtDeployDetailNodeCustProp); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailNodeCustProp")
		log.Fatal(err)
	}
	defer stmt_NodeCustProp.Close()
	if stmt_DefaultDC, err = tk.conn.Prepare(tkStmtDeployDetailDefaultDatacenter); err != nil {
		log.Println("Failed to prepare: tkStmtDeployDetailDefaultDatacenter")
		log.Fatal(err)
	}
	defer stmt_DefaultDC.Close()

	//
	if rows, err = stmt_List.Query(); err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		detail := proto.Deployment{}

		err = rows.Scan(
			&instanceCfgId,
		)

		//
		detail.CheckInstance = &proto.CheckInstance{
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
		detail.Check = &proto.Check{
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
		detail.CheckConfig = &proto.CheckConfig{
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
			&detail.CheckConfig.Name,
			&detail.CheckConfig.Interval,
			&detail.CheckConfig.IsActive,
			&detail.CheckConfig.IsEnabled,
			&detail.CheckConfig.ExternalId,
		)

		//
		detail.CheckConfig.Thresholds = []proto.CheckConfigThreshold{}
		thresh, err = stmt_Threshold.Query(detail.CheckConfig.Id)
		if err != nil {
			log.Println(`DANGER WILL ROBINSON! Failed to get thresholds for:`, detail.CheckConfig.Id)
			continue
		}
		defer thresh.Close()

		for thresh.Next() {
			thr := proto.CheckConfigThreshold{
				Predicate: proto.Predicate{},
				Level:     proto.Level{},
			}

			err = thresh.Scan(
				&thr.Predicate.Symbol,
				&thr.Value,
				&thr.Level.Name,
				&thr.Level.ShortName,
				&thr.Level.Numeric,
			)
			detail.CheckConfig.Thresholds = append(detail.CheckConfig.Thresholds, thr)
		}

		// XXX TODO
		//detail.CheckConfiguration.Constraints = []somaproto.CheckConfigurationConstraint{}
		detail.CheckConfig.Constraints = nil

		//
		detail.Capability = &proto.Capability{
			Id: detail.Check.CapabilityId,
		}
		detail.Monitoring = &proto.Monitoring{}
		detail.Metric = &proto.Metric{}
		detail.Unit = &proto.Unit{}
		stmt_CapMonMetric.QueryRow(detail.Capability.Id).Scan(
			&detail.Capability.Metric,
			&detail.Capability.MonitoringId,
			&detail.Capability.View,
			&detail.Capability.Thresholds,
			&detail.Monitoring.Name,
			&detail.Monitoring.Mode,
			&detail.Monitoring.Contact,
			&detail.Monitoring.TeamId,
			&callback,
			&detail.Metric.Unit,
			&detail.Metric.Description,
			&detail.Unit.Name,
		)
		if callback.Valid {
			detail.Monitoring.Callback = callback.String
		} else {
			detail.Monitoring.Callback = ""
		}
		detail.Unit.Unit = detail.Metric.Unit
		detail.Metric.Path = detail.Capability.Metric
		detail.Monitoring.Id = detail.Capability.MonitoringId
		detail.Capability.Name = fmt.Sprintf("%s.%s.%s",
			detail.Monitoring.Name,
			detail.Capability.View,
			detail.Metric.Path,
		)
		detail.View = detail.Capability.View

		//
		detail.Metric.Packages = &[]proto.MetricPackage{}
		pkgs, _ = stmt_Pkgs.Query(detail.Metric.Path)
		defer pkgs.Close()

		for pkgs.Next() {
			pkg := proto.MetricPackage{}

			err = pkgs.Scan(
				&pkg.Provider,
				&pkg.Name,
			)
			*detail.Metric.Packages = append(*detail.Metric.Packages, pkg)
		}

		//
		detail.Oncall = &proto.Oncall{}
		detail.Service = &proto.PropertyService{}
		switch objType {
		case "group":
			// fetch the group object
			detail.Group = &proto.Group{
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
			detail.Team = &proto.Team{
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
					detail.Service.Attributes = []proto.ServiceAttribute{}
					fm := map[string]string{}
					_ = json.Unmarshal([]byte(detail.CheckInstance.InstanceServiceConfig), &fm)
					for k, v := range fm {
						a := proto.ServiceAttribute{
							Name:  k,
							Value: v,
						}
						detail.Service.Attributes = append(detail.Service.Attributes, a)
					}
				}
			}
			// fetch system properties
			detail.Properties = &[]proto.PropertySystem{}
			gSysProps, _ = stmt_GroupSysProp.Query(detail.Group.Id, detail.View)
			defer gSysProps.Close()

			for gSysProps.Next() {
				prop := proto.PropertySystem{}
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
			detail.CustomProperties = &[]proto.PropertyCustom{}
			gCustProps, _ = stmt_GroupCustProp.Query(detail.Group.Id, detail.View)
			defer gCustProps.Close()

			for gCustProps.Next() {
				prop := proto.PropertyCustom{}
				gCustProps.Scan(
					&prop.Id,
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
			detail.Cluster = &proto.Cluster{
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
			detail.Team = &proto.Team{
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
					detail.Service.Attributes = []proto.ServiceAttribute{}
					fm := map[string]string{}
					_ = json.Unmarshal([]byte(detail.CheckInstance.InstanceServiceConfig), &fm)
					for k, v := range fm {
						a := proto.ServiceAttribute{
							Name:  k,
							Value: v,
						}
						detail.Service.Attributes = append(detail.Service.Attributes, a)
					}
				}
			}
			// fetch system properties
			detail.Properties = &[]proto.PropertySystem{}
			cSysProps, _ = stmt_ClusterSysProp.Query(detail.Cluster.Id, detail.View)
			defer cSysProps.Close()

			for cSysProps.Next() {
				prop := proto.PropertySystem{}
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
			detail.CustomProperties = &[]proto.PropertyCustom{}
			cCustProps, _ = stmt_ClusterCustProp.Query(detail.Cluster.Id, detail.View)
			defer cCustProps.Close()

			for cCustProps.Next() {
				prop := proto.PropertyCustom{}
				cCustProps.Scan(
					&prop.Id,
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
			detail.Server = &proto.Server{}
			detail.Node = &proto.Node{
				Id: objId,
			}
			stmt_Node.QueryRow(objId).Scan(
				&detail.Node.AssetId,
				&detail.Node.Name,
				&detail.Node.TeamId,
				&detail.Node.ServerId,
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
			detail.Server.Id = detail.Node.ServerId
			detail.Datacenter = detail.Server.Datacenter
			// fetch team information
			detail.Team = &proto.Team{
				Id: detail.Node.TeamId,
			}
			// fetch oncall information if the property is set,
			// otherwise cleanup detail.Oncall
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
					detail.Service.Attributes = []proto.ServiceAttribute{}
					fm := map[string]string{}
					_ = json.Unmarshal([]byte(detail.CheckInstance.InstanceServiceConfig), &fm)
					for k, v := range fm {
						a := proto.ServiceAttribute{
							Name:  k,
							Value: v,
						}
						detail.Service.Attributes = append(detail.Service.Attributes, a)
					}
				}
			}
			// fetch system properties
			detail.Properties = &[]proto.PropertySystem{}
			nSysProps, _ = stmt_NodeSysProp.Query(detail.Node.Id, detail.View)
			defer nSysProps.Close()

			for nSysProps.Next() {
				prop := proto.PropertySystem{}
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
			detail.CustomProperties = &[]proto.PropertyCustom{}
			nCustProps, _ = stmt_NodeCustProp.Query(detail.Node.Id, detail.View)
			defer nCustProps.Close()

			for nCustProps.Next() {
				prop := proto.PropertyCustom{}
				gCustProps.Scan(
					&prop.Id,
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
			&detail.Team.LdapId,
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
