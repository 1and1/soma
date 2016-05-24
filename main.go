package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
)

// global variables
var (
	// database connection
	conn *sql.DB
	// lookup table for go routine input channels
	handlerMap = make(map[string]interface{})
	// config file runtime configuration
	SomaCfg SomaConfig
	// this offset influences the biggest date representable in
	// the system without overflow
	unixToInternalOffset int64 = 62135596800
	// this will be used as mapping for the PostgreSQL time value
	// -infinity. Dates earlier than this will be truncated to
	// NegTimeInf. RFC3339: -8192-01-01T00:00:00Z
	NegTimeInf = time.Date(-8192, time.January, 1, 0, 0, 0, 0, time.UTC)
	// this will be used as mapping for the PostgreSQL time value
	// +infinity. It is as far as research showed close to the highest
	// time value Go can represent.
	// RFC: 219248499-12-06 15:30:07.999999999 +0000 UTC
	PosTimeInf = time.Unix(1<<63-1-unixToInternalOffset, 999999999)
)

const (
	// Format string for millisecond precision RFC3339
	rfc3339Milli string = "2006-01-02T15:04:05.000Z07:00"
	// SOMA version
	SomaVersion string = `0.7.21`
	// Logging format strings
	LogStrOK  = `Subsystem=%s, Request=%s, InternalCode=%d, ExternalCode=%d`
	LogStrErr = `Subsystem=%s, Action=%s, InternalCode=%d, Error=%s`
)

func main() {
	var (
		configFlag, configFile string
		err                    error
	)
	flag.StringVar(&configFlag, "config", "/srv/soma/conf/soma.conf", "Configuration file location")
	flag.Parse()

	log.Printf("Starting runtime config initialization, SOMA v%s", SomaVersion)
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

	router.HEAD("/", Ping)

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

	router.GET("/authenticate/validate/", BasicAuth(AuthenticationValidate))

	if !SomaCfg.ReadOnly {
		router.POST("/views/", BasicAuth(AddView))
		router.DELETE("/views/:view", BasicAuth(DeleteView))
		router.PATCH("/views/:view", BasicAuth(RenameView))

		router.POST("/environments", BasicAuth(AddEnvironment))
		router.DELETE("/environments/:environment", BasicAuth(DeleteEnvironment))
		router.PUT("/environments/:environment", BasicAuth(RenameEnvironment))

		router.POST("/objstates", BasicAuth(AddObjectState))
		router.DELETE("/objstates/:state", BasicAuth(DeleteObjectState))
		router.PUT("/objstates/:state", BasicAuth(RenameObjectState))

		router.POST("/objtypes", BasicAuth(AddObjectType))
		router.DELETE("/objtypes/:type", BasicAuth(DeleteObjectType))
		router.PUT("/objtypes/:type", BasicAuth(RenameObjectType))

		router.POST("/datacenters", BasicAuth(AddDatacenter))
		router.DELETE("/datacenters/:datacenter", BasicAuth(DeleteDatacenter))
		router.PUT("/datacenters/:datacenter", BasicAuth(RenameDatacenter))

		router.PATCH("/datacentergroups/:datacentergroup", BasicAuth(AddDatacenterToGroup))
		router.DELETE("/datacentergroups/:datacentergroup", BasicAuth(DeleteDatacenterFromGroup))

		router.POST("/levels/", BasicAuth(AddLevel))
		router.DELETE("/levels/:level", BasicAuth(DeleteLevel))

		router.POST("/predicates/", BasicAuth(AddPredicate))
		router.DELETE("/predicates/:predicate", BasicAuth(DeletePredicate))

		router.POST("/status/", BasicAuth(AddStatus))
		router.DELETE("/status/:status", BasicAuth(DeleteStatus))

		router.POST("/oncall/", BasicAuth(AddOncall))
		router.PATCH("/oncall/:oncall", BasicAuth(UpdateOncall))
		router.DELETE("/oncall/:oncall", BasicAuth(DeleteOncall))

		router.POST("/teams/", BasicAuth(AddTeam))
		router.DELETE("/teams/:team", BasicAuth(DeleteTeam))

		router.POST("/servers/", BasicAuth(AddServer))
		router.DELETE("/servers/:server", BasicAuth(DeleteServer))
		router.PUT("/servers/:server", BasicAuth(InsertNullServer))

		router.POST("/units/", BasicAuth(AddUnit))
		router.DELETE("/units/:unit", BasicAuth(DeleteUnit))

		router.POST("/providers/", BasicAuth(AddProvider))
		router.DELETE("/providers/:provider", BasicAuth(DeleteProvider))

		router.POST("/metrics/", BasicAuth(AddMetric))
		router.DELETE("/metrics/:metric", BasicAuth(DeleteMetric))

		router.POST("/modes/", BasicAuth(AddMode))
		router.DELETE("/modes/:mode", BasicAuth(DeleteMode))

		router.POST("/users/", BasicAuth(AddUser))
		router.DELETE("/users/:user", BasicAuth(DeleteUser))

		router.POST("/monitoring/", BasicAuth(AddMonitoring))
		router.DELETE("/monitoring/:monitoring", BasicAuth(DeleteMonitoring))

		router.POST("/capability/", BasicAuth(AddCapability))
		router.DELETE("/capability/:capability", BasicAuth(DeleteCapability))

		router.POST("/property/native/", BasicAuth(AddProperty))
		router.DELETE("/property/native/:native", BasicAuth(DeleteProperty))

		router.POST("/property/system/", BasicAuth(AddProperty))
		router.DELETE("/property/system/:system", BasicAuth(DeleteProperty))

		router.POST("/property/custom/:repository/", BasicAuth(AddProperty))
		router.DELETE("/property/custom/:repository/:custom", BasicAuth(DeleteProperty))

		router.POST("/property/service/global/", BasicAuth(AddProperty))
		router.DELETE("/property/service/global/:service", BasicAuth(DeleteProperty))

		router.POST("/property/service/team/:team/", BasicAuth(AddProperty))
		router.DELETE("/property/service/team/:team/:service", BasicAuth(DeleteProperty))

		router.POST("/validity/", BasicAuth(AddValidity))
		router.DELETE("/validity/:property", BasicAuth(DeleteValidity))

		router.POST("/attributes/", BasicAuth(AddAttribute))
		router.DELETE("/attributes/:attribute", BasicAuth(DeleteAttribute))

		router.POST("/repository/", BasicAuth(AddRepository))
		router.POST("/repository/:repository/property/:type/", BasicAuth(AddPropertyToRepository))

		router.POST("/buckets/", BasicAuth(AddBucket))
		router.POST("/buckets/:bucket/property/:type/", BasicAuth(AddPropertyToBucket))

		router.POST("/groups/", BasicAuth(AddGroup))
		router.POST("/groups/:group/members/", BasicAuth(AddMemberToGroup))
		router.POST("/groups/:group/property/:type/", BasicAuth(AddPropertyToGroup))

		router.POST("/clusters/", BasicAuth(AddCluster))
		router.POST("/clusters/:cluster/members/", BasicAuth(AddMemberToCluster))
		router.POST("/clusters/:cluster/property/:type/", BasicAuth(AddPropertyToCluster))

		router.POST("/nodes/", BasicAuth(AddNode))
		router.DELETE("/nodes/:node", BasicAuth(DeleteNode))
		router.PUT("/nodes/:node/config", BasicAuth(AssignNode))
		router.POST("/nodes/:node/property/:type/", BasicAuth(AddPropertyToNode))

		router.POST("/checks/:repository/", BasicAuth(AddCheckConfiguration))

		router.GET("/deployments/id/:uuid", DeliverDeploymentDetails)
		router.GET("/deployments/monitoring/:uuid", DeliverMonitoringDeployments)
		router.GET("/deployments/monitoring/:uuid/:all", DeliverMonitoringDeployments)
		router.PATCH("/deployments/id/:uuid/:result", UpdateDeploymentDetails)

		router.POST("/authenticate/", AuthenticationKex)
		router.PUT("/authenticate/bootstrap/:uuid", AuthenticationBootstrapRoot)
		//router.PATCH("/authenticate/root/restrict", AuthenticationRestrictRoot) XXX -> move to somadbctl
		//router.GET("/authenticate/token/", AuthenticationListTokens)
		router.PUT("/authenticate/token/:uuid", BasicAuth(AuthenticationIssueToken))
		router.PUT("/authenticate/activate/:uuid", AuthenticationActivateUser)
		//router.DELETE("/authenticate/invalidate/token/", AuthenticationInvalidateToken)
		//router.DELETE("/authenticate/invalidate/all/", AuthenticationInvalidateAllTokens)
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
