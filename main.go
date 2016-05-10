package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
)

// global variables
var conn *sql.DB
var handlerMap = make(map[string]interface{})
var SomaCfg SomaConfig

func main() {
	var (
		configFlag, configFile string
		err                    error
	)
	flag.StringVar(&configFlag, "config", "/srv/soma/conf/soma.conf", "Configuration file location")
	flag.Parse()

	version := "0.7.9"
	log.Printf("Starting runtime config initialization, SOMA v%s", version)
	/*
	 * Read configuration file
	 */
	if configFile, err = filepath.Abs(configFlag); err != nil {
		log.Fatal(err)
	}
	if configFile, err = filepath.EvalSymlinks(configFile); err != nil {
		log.Fatal(err)
	}
	err = SomaCfg.readConfigFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	/*
	 * Construct listen address
	 */
	SomaCfg.Daemon.url = &url.URL{}
	SomaCfg.Daemon.url.Host = fmt.Sprintf("%s:%s", SomaCfg.Daemon.Listen, SomaCfg.Daemon.Port)
	if SomaCfg.Daemon.Tls {
		SomaCfg.Daemon.url.Scheme = "https"
		if ok, pt := govalidator.IsFilePath(SomaCfg.Daemon.Cert); !ok {
			log.Fatal("Missing required certificate configuration config/daemon/cert-file")
		} else {
			if pt != govalidator.Unix {
				log.Fatal("config/daemon/cert-File: valid Windows paths are not helpful")
			}
		}
		if ok, pt := govalidator.IsFilePath(SomaCfg.Daemon.Key); !ok {
			log.Fatal("Missing required key configuration config/daemon/key-file")
		} else {
			if pt != govalidator.Unix {
				log.Fatal("config/daemon/key-file: valid Windows paths are not helpful")
			}
		}
	} else {
		SomaCfg.Daemon.url.Scheme = "http"
	}

	connectToDatabase()
	go pingDatabase()

	startHandlers()

	router := httprouter.New()

	router.GET("/views", ListView)
	router.GET("/views/:view", ShowView)

	router.GET("/environments", ListEnvironments)
	router.GET("/environments/:environment", ShowEnvironment)

	router.GET("/objstates", ListObjectStates)
	router.GET("/objstates/:state", ShowObjectState)

	router.GET("/objtypes", ListObjectTypes)
	router.GET("/objtypes/:type", ShowObjectType)

	router.GET("/datacenters", ListDatacenters)
	router.GET("/datacenters/:datacenter", ShowDatacenter)

	router.GET("/datacentergroups", ListDatacenterGroups)
	router.GET("/datacentergroups/:datacentergroup", ShowDatacenterGroup)

	router.GET("/levels/", ListLevel)
	router.GET("/levels/:level", ShowLevel)
	router.POST("/filter/levels/", ListLevel)

	router.GET("/predicates/", ListPredicate)
	router.GET("/predicates/:predicate", ShowPredicate)

	router.GET("/status/", ListStatus)
	router.GET("/status/:status", ShowStatus)

	router.GET("/oncall/", ListOncall)
	router.GET("/oncall/:oncall", ShowOncall)
	router.POST("/filter/oncall/", ListOncall)

	router.GET("/teams/", ListTeam)
	router.GET("/teams/:team", ShowTeam)
	router.POST("/filter/teams/", ListTeam)

	router.GET("/nodes/", ListNode)
	router.GET("/nodes/:node", ShowNode)
	router.GET("/nodes/:node/config", ShowNodeConfig)
	router.POST("/filter/nodes/", ListNode)

	router.GET("/servers/", ListServer)
	router.GET("/servers/:server", ShowServer)
	router.POST("/filter/servers/", ListServer)

	router.GET("/units/", ListUnit)
	router.GET("/units/:unit", ShowUnit)

	router.GET("/providers/", ListProvider)
	router.GET("/providers/:provider", ShowProvider)

	router.GET("/metrics/", ListMetric)
	router.GET("/metrics/:metric", ShowMetric)

	router.GET("/modes/", ListMode)
	router.GET("/modes/:mode", ShowMode)

	router.GET("/users/", ListUser)
	router.GET("/users/:user", ShowUser)
	router.POST("/filter/users/", ListUser)

	router.GET("/monitoring/", ListMonitoring)
	router.GET("/monitoring/:monitoring", ShowMonitoring)
	router.POST("/filter/monitoring/", ListMonitoring)

	router.GET("/capability/", ListCapability)
	router.GET("/capability/:capability", ShowCapability)
	router.POST("/filter/capability/", ListCapability)

	router.GET("/property/native/", ListProperty)
	router.GET("/property/native/:native", ShowProperty)

	router.GET("/property/system/", ListProperty)
	router.GET("/property/system/:system", ShowProperty)
	router.POST("/filter/property/system/", ListProperty)

	router.GET("/property/custom/:repository/", ListProperty)
	router.GET("/property/custom/:repository/:custom", ShowProperty)
	router.POST("/filter/property/custom/:repository/", ListProperty)

	router.GET("/property/service/global/", ListProperty)
	router.GET("/property/service/global/:service", ShowProperty)
	router.POST("/filter/property/service/global/", ListProperty)

	router.GET("/property/service/team/:team/", ListProperty)
	router.GET("/property/service/team/:team/:service", ShowProperty)
	router.POST("/filter/property/service/team/:team/", ListProperty)

	router.GET("/validity/", ListValidity)
	router.GET("/validity/:property", ShowValidity)

	router.GET("/attributes/", ListAttribute)
	router.GET("/attributes/:attribute", ShowAttribute)

	router.GET("/repository/", ListRepository)
	router.GET("/repository/:repository", ShowRepository)
	router.POST("/filter/repository/", ListRepository)

	router.GET("/buckets/", ListBucket)
	router.GET("/buckets/:bucket", ShowBucket)
	router.POST("/filter/buckets/", ListBucket)

	router.GET("/groups/", ListGroup)
	router.GET("/groups/:group", ShowGroup)
	router.GET("/groups/:group/members/", ListGroupMembers)
	router.POST("/filter/groups/", ListGroup)

	router.GET("/clusters/", ListCluster)
	router.GET("/clusters/:cluster", ShowCluster)
	router.GET("/clusters/:cluster/members/", ListClusterMembers)
	router.POST("/filter/clusters/", ListCluster)

	router.GET("/checks/:repository/", ListCheckConfiguration)
	router.GET("/checks/:repository/:check", ShowCheckConfiguration)
	router.POST("/filter/checks/:repository/", ListCheckConfiguration)

	router.GET("/hostdeployment/:system/:assetid", GetHostDeployment)
	router.POST("/hostdeployment/:system/:assetid", AssembleHostUpdate)

	if !SomaCfg.ReadOnly {
		router.POST("/views/", AddView)
		router.DELETE("/views/:view", DeleteView)
		router.PATCH("/views/:view", RenameView)

		router.POST("/environments", AddEnvironment)
		router.DELETE("/environments/:environment", DeleteEnvironment)
		router.PUT("/environments/:environment", RenameEnvironment)

		router.POST("/objstates", AddObjectState)
		router.DELETE("/objstates/:state", DeleteObjectState)
		router.PUT("/objstates/:state", RenameObjectState)

		router.POST("/objtypes", AddObjectType)
		router.DELETE("/objtypes/:type", DeleteObjectType)
		router.PUT("/objtypes/:type", RenameObjectType)

		router.POST("/datacenters", AddDatacenter)
		router.DELETE("/datacenters/:datacenter", DeleteDatacenter)
		router.PUT("/datacenters/:datacenter", RenameDatacenter)

		router.PATCH("/datacentergroups/:datacentergroup", AddDatacenterToGroup)
		router.DELETE("/datacentergroups/:datacentergroup", DeleteDatacenterFromGroup)

		router.POST("/levels/", AddLevel)
		router.DELETE("/levels/:level", DeleteLevel)

		router.POST("/predicates/", AddPredicate)
		router.DELETE("/predicates/:predicate", DeletePredicate)

		router.POST("/status/", AddStatus)
		router.DELETE("/status/:status", DeleteStatus)

		router.POST("/oncall/", AddOncall)
		router.PATCH("/oncall/:oncall", UpdateOncall)
		router.DELETE("/oncall/:oncall", DeleteOncall)

		router.POST("/teams/", AddTeam)
		router.DELETE("/teams/:team", DeleteTeam)

		router.POST("/servers/", AddServer)
		router.DELETE("/servers/:server", DeleteServer)
		router.PUT("/servers/:server", InsertNullServer)

		router.POST("/units/", AddUnit)
		router.DELETE("/units/:unit", DeleteUnit)

		router.POST("/providers/", AddProvider)
		router.DELETE("/providers/:provider", DeleteProvider)

		router.POST("/metrics/", AddMetric)
		router.DELETE("/metrics/:metric", DeleteMetric)

		router.POST("/modes/", AddMode)
		router.DELETE("/modes/:mode", DeleteMode)

		router.POST("/users/", AddUser)
		router.DELETE("/users/:user", DeleteUser)

		router.POST("/monitoring/", AddMonitoring)
		router.DELETE("/monitoring/:monitoring", DeleteMonitoring)

		router.POST("/capability/", AddCapability)
		router.DELETE("/capability/:capability", DeleteCapability)

		router.POST("/property/native/", AddProperty)
		router.DELETE("/property/native/:native", DeleteProperty)

		router.POST("/property/system/", AddProperty)
		router.DELETE("/property/system/:system", DeleteProperty)

		router.POST("/property/custom/:repository/", AddProperty)
		router.DELETE("/property/custom/:repository/:custom", DeleteProperty)

		router.POST("/property/service/global/", AddProperty)
		router.DELETE("/property/service/global/:service", DeleteProperty)

		router.POST("/property/service/team/:team/", AddProperty)
		router.DELETE("/property/service/team/:team/:service", DeleteProperty)

		router.POST("/validity/", AddValidity)
		router.DELETE("/validity/:property", DeleteValidity)

		router.POST("/attributes/", AddAttribute)
		router.DELETE("/attributes/:attribute", DeleteAttribute)

		router.POST("/repository/", AddRepository)
		router.POST("/repository/:repository/property/:type/", AddPropertyToRepository)

		router.POST("/buckets/", AddBucket)
		router.POST("/buckets/:bucket/property/:type/", AddPropertyToBucket)

		router.POST("/groups/", AddGroup)
		router.POST("/groups/:group/members/", AddMemberToGroup)
		router.POST("/groups/:group/property/:type/", AddPropertyToGroup)

		router.POST("/clusters/", AddCluster)
		router.POST("/clusters/:cluster/members/", AddMemberToCluster)
		router.POST("/clusters/:cluster/property/:type/", AddPropertyToCluster)

		router.POST("/nodes/", AddNode)
		router.DELETE("/nodes/:node", DeleteNode)
		router.PUT("/nodes/:node/config", AssignNode)
		router.POST("/nodes/:node/property/:type/", AddPropertyToNode)

		router.POST("/checks/:repository/", AddCheckConfiguration)

		router.GET("/deployments/id/:uuid", DeliverDeploymentDetails)
		router.GET("/deployments/monitoring/:uuid", DeliverMonitoringDeployments)
		router.GET("/deployments/monitoring/:uuid/:all", DeliverMonitoringDeployments)
		router.PATCH("/deployments/id/:uuid/:result", UpdateDeploymentDetails)
	}

	if SomaCfg.Daemon.Tls {
		log.Fatal(http.ListenAndServeTLS(
			SomaCfg.Daemon.url.Host,
			SomaCfg.Daemon.Cert,
			SomaCfg.Daemon.Key,
			router))
	} else {
		log.Fatal(http.ListenAndServe(SomaCfg.Daemon.url.Host, router))
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
