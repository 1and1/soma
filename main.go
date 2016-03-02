package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// global variables
var conn *sql.DB
var handlerMap = make(map[string]interface{})
var SomaCfg SomaConfig

func main() {
	version := "0.2.7"
	log.Printf("Starting runtime config initialization, SOMA v%s", version)
	err := SomaCfg.readConfigFile("soma.conf")
	if err != nil {
		log.Fatal(err)
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

	router.GET("/property/custom/:repository/", ListProperty)
	router.GET("/property/custom/:repository/:custom", ShowProperty)
	router.POST("/filter/property/custom/", ListProperty)

	router.GET("/property/service/global/", ListProperty)
	router.GET("/property/service/global/:service", ShowProperty)

	router.GET("/property/service/team/:team/", ListProperty)
	router.GET("/property/service/team/:team/:service", ShowProperty)

	router.GET("/attributes/", ListAttribute)
	router.GET("/attributes/:attribute", ShowAttribute)

	router.GET("/repository/", ListRepository)
	router.GET("/repository/:repository", ShowRepository)
	router.GET("/filter/repository/", ListRepository)

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
	}

	log.Fatal(http.ListenAndServe(":8888", router))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
