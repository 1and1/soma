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
	SomaVersion string = `0.7.46`
	// Logging format strings
	LogStrReq = `Subsystem=%s, Request=%s, User=%s, Addr=%s`
	LogStrArg = `Subsystem=%s, Request=%s, User=%s, Addr=%s, Arg=%s`
	LogStrOK  = `Subsystem=%s, Result=%s, InternalCode=%d, ExternalCode=%d`
	LogStrErr = `Subsystem=%s, Action=%s, InternalCode=%d, Error=%s`
)

func main() {
	var (
		configFlag, configFile, obsRepoFlag string
		noPokeFlag                          bool
		err                                 error
	)
	flag.StringVar(&configFlag, "config", "/srv/soma/conf/soma.conf", "Configuration file location")
	flag.StringVar(&obsRepoFlag, "repo", "", "Observer target repository")
	flag.BoolVar(&noPokeFlag, "nopoke", false, "Disable lifecycle pokes")
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

	// force observer mode if the cli argument was present
	if obsRepoFlag != `` {
		SomaCfg.Observer = true
		SomaCfg.ObserverRepo = obsRepoFlag
	}

	if noPokeFlag {
		SomaCfg.NoPoke = true
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

	router.HEAD(`/`, Ping)

	router.GET(`/attributes/:attribute`, BasicAuth(ShowAttribute))
	router.GET(`/attributes/`, BasicAuth(ListAttribute))
	router.GET(`/authenticate/validate/`, BasicAuth(AuthenticationValidate))
	router.GET(`/buckets/:bucket/:tree`, BasicAuth(OutputTree))
	router.GET(`/buckets/:bucket`, BasicAuth(ShowBucket))
	router.GET(`/buckets/`, BasicAuth(ListBucket))
	router.GET(`/capability/:capability`, BasicAuth(ShowCapability))
	router.GET(`/capability/`, BasicAuth(ListCapability))
	router.GET(`/category/:category`, BasicAuth(ShowCategory))
	router.GET(`/category/`, BasicAuth(ListCategory))
	router.GET(`/checks/:repository/:check`, BasicAuth(ShowCheckConfiguration))
	router.GET(`/checks/:repository/`, BasicAuth(ListCheckConfiguration))
	router.GET(`/clusters/:cluster/members/`, BasicAuth(ListClusterMembers))
	router.GET(`/clusters/:cluster/tree/:tree`, BasicAuth(OutputTree))
	router.GET(`/clusters/:cluster`, BasicAuth(ShowCluster))
	router.GET(`/clusters/`, BasicAuth(ListCluster))
	router.GET(`/datacentergroups/:datacentergroup`, BasicAuth(ShowDatacenterGroup))
	router.GET(`/datacentergroups/`, BasicAuth(ListDatacenterGroups))
	router.GET(`/datacenters/:datacenter`, BasicAuth(ShowDatacenter))
	router.GET(`/datacenters/`, BasicAuth(ListDatacenters))
	router.GET(`/environments/:environment`, BasicAuth(ShowEnvironment))
	router.GET(`/environments/`, BasicAuth(ListEnvironments))
	router.GET(`/groups/:group/members/`, BasicAuth(ListGroupMembers))
	router.GET(`/groups/:group/tree/:tree`, BasicAuth(OutputTree))
	router.GET(`/groups/:group`, BasicAuth(ShowGroup))
	router.GET(`/groups/`, BasicAuth(ListGroup))
	router.GET(`/hostdeployment/:system/:assetid`, GetHostDeployment)
	router.GET(`/jobs/:jobid`, BasicAuth(ShowJob))
	router.GET(`/jobs/`, BasicAuth(ListJobs))
	router.GET(`/levels/:level`, BasicAuth(ShowLevel))
	router.GET(`/levels/`, BasicAuth(ListLevel))
	router.GET(`/metrics/:metric`, BasicAuth(ShowMetric))
	router.GET(`/metrics/`, BasicAuth(ListMetric))
	router.GET(`/modes/:mode`, BasicAuth(ShowMode))
	router.GET(`/modes/`, BasicAuth(ListMode))
	router.GET(`/monitoring/:monitoring`, BasicAuth(ShowMonitoring))
	router.GET(`/monitoring/`, BasicAuth(ListMonitoring))
	router.GET(`/nodes/:node/config`, BasicAuth(ShowNodeConfig))
	router.GET(`/nodes/:node/tree/:tree`, BasicAuth(OutputTree))
	router.GET(`/nodes/:node`, BasicAuth(ShowNode))
	router.GET(`/nodes/`, BasicAuth(ListNode))
	router.GET(`/objstates/:state`, BasicAuth(ShowObjectState))
	router.GET(`/objstates/`, BasicAuth(ListObjectStates))
	router.GET(`/objtypes/:type`, BasicAuth(ShowObjectType))
	router.GET(`/objtypes/`, BasicAuth(ListObjectTypes))
	router.GET(`/oncall/:oncall`, BasicAuth(ShowOncall))
	router.GET(`/oncall/`, BasicAuth(ListOncall))
	router.GET(`/permission/:permission`, BasicAuth(ShowPermission))
	router.GET(`/permission/`, BasicAuth(ListPermission))
	router.GET(`/predicates/:predicate`, BasicAuth(ShowPredicate))
	router.GET(`/predicates/`, BasicAuth(ListPredicate))
	router.GET(`/property/custom/:repository/:custom`, BasicAuth(ShowProperty))
	router.GET(`/property/custom/:repository/`, BasicAuth(ListProperty))
	router.GET(`/property/native/:native`, BasicAuth(ShowProperty))
	router.GET(`/property/native/`, BasicAuth(ListProperty))
	router.GET(`/property/service/global/:service`, BasicAuth(ShowProperty))
	router.GET(`/property/service/global/`, BasicAuth(ListProperty))
	router.GET(`/property/service/team/:team/:service`, BasicAuth(ShowProperty))
	router.GET(`/property/service/team/:team/`, BasicAuth(ListProperty))
	router.GET(`/property/system/:system`, BasicAuth(ShowProperty))
	router.GET(`/property/system/`, BasicAuth(ListProperty))
	router.GET(`/providers/:provider`, BasicAuth(ShowProvider))
	router.GET(`/providers/`, BasicAuth(ListProvider))
	router.GET(`/repository/:repository/:tree`, BasicAuth(OutputTree))
	router.GET(`/repository/:repository`, BasicAuth(ShowRepository))
	router.GET(`/repository/`, BasicAuth(ListRepository))
	router.GET(`/servers/:server`, BasicAuth(ShowServer))
	router.GET(`/servers/`, BasicAuth(ListServer))
	router.GET(`/status/:status`, BasicAuth(ShowStatus))
	router.GET(`/status/`, BasicAuth(ListStatus))
	router.GET(`/sync/datacenters/`, BasicAuth(SyncDatacenters))
	router.GET(`/sync/nodes/`, BasicAuth(SyncNode))
	router.GET(`/sync/servers/`, BasicAuth(SyncServer))
	router.GET(`/sync/teams/`, BasicAuth(SyncTeam))
	router.GET(`/sync/users/`, BasicAuth(SyncUser))
	router.GET(`/teams/:team`, BasicAuth(ShowTeam))
	router.GET(`/teams/`, BasicAuth(ListTeam))
	router.GET(`/units/:unit`, BasicAuth(ShowUnit))
	router.GET(`/units/`, BasicAuth(ListUnit))
	router.GET(`/users/:user`, BasicAuth(ShowUser))
	router.GET(`/users/`, BasicAuth(ListUser))
	router.GET(`/validity/:property`, BasicAuth(ShowValidity))
	router.GET(`/validity/`, BasicAuth(ListValidity))
	router.GET(`/views/:view`, BasicAuth(ShowView))
	router.GET(`/views/`, BasicAuth(ListView))
	router.POST(`/filter/buckets/`, BasicAuth(ListBucket))
	router.POST(`/filter/capability/`, BasicAuth(ListCapability))
	router.POST(`/filter/checks/:repository/`, BasicAuth(ListCheckConfiguration))
	router.POST(`/filter/clusters/`, BasicAuth(ListCluster))
	router.POST(`/filter/grant/`, BasicAuth(SearchGrant))
	router.POST(`/filter/groups/`, BasicAuth(ListGroup))
	router.POST(`/filter/jobs/`, BasicAuth(SearchJob))
	router.POST(`/filter/levels/`, BasicAuth(ListLevel))
	router.POST(`/filter/monitoring/`, BasicAuth(ListMonitoring))
	router.POST(`/filter/nodes/`, BasicAuth(ListNode))
	router.POST(`/filter/oncall/`, BasicAuth(ListOncall))
	router.POST(`/filter/permission/`, BasicAuth(SearchPermission))
	router.POST(`/filter/property/custom/:repository/`, BasicAuth(ListProperty))
	router.POST(`/filter/property/service/global/`, BasicAuth(ListProperty))
	router.POST(`/filter/property/service/team/:team/`, BasicAuth(ListProperty))
	router.POST(`/filter/property/system/`, BasicAuth(ListProperty))
	router.POST(`/filter/repository/`, BasicAuth(ListRepository))
	router.POST(`/filter/servers/`, BasicAuth(SearchServer))
	router.POST(`/filter/teams/`, BasicAuth(ListTeam))
	router.POST(`/filter/users/`, BasicAuth(ListUser))
	router.POST(`/hostdeployment/:system/:assetid`, AssembleHostUpdate)

	if !SomaCfg.ReadOnly {
		router.DELETE(`/attributes/:attribute`, BasicAuth(DeleteAttribute))
		router.DELETE(`/buckets/:bucket/property/:type/:source`, BasicAuth(DeletePropertyFromBucket))
		router.DELETE(`/capability/:capability`, BasicAuth(DeleteCapability))
		router.DELETE(`/category/:category`, BasicAuth(DeleteCategory))
		router.DELETE(`/checks/:repository/:check`, BasicAuth(DeleteCheckConfiguration))
		router.DELETE(`/clusters/:cluster/property/:type/:source`, BasicAuth(DeletePropertyFromCluster))
		router.DELETE(`/datacentergroups/:datacentergroup`, BasicAuth(DeleteDatacenterFromGroup))
		router.DELETE(`/datacenters/:datacenter`, BasicAuth(DeleteDatacenter))
		router.DELETE(`/environments/:environment`, BasicAuth(DeleteEnvironment))
		router.DELETE(`/grant/global/:rtyp/:rid/:grant`, BasicAuth(RevokeGlobalRight))
		router.DELETE(`/grant/limited/:rtyp/:rid/:scope/:uuid/:grant`, BasicAuth(RevokeLimitedRight))
		router.DELETE(`/grant/system/:rtyp/:rid/:grant`, BasicAuth(RevokeSystemRight))
		router.DELETE(`/groups/:group/property/:type/:source`, BasicAuth(DeletePropertyFromGroup))
		router.DELETE(`/levels/:level`, BasicAuth(DeleteLevel))
		router.DELETE(`/metrics/:metric`, BasicAuth(DeleteMetric))
		router.DELETE(`/modes/:mode`, BasicAuth(DeleteMode))
		router.DELETE(`/monitoring/:monitoring`, BasicAuth(DeleteMonitoring))
		router.DELETE(`/nodes/:node/property/:type/:source`, BasicAuth(DeletePropertyFromNode))
		router.DELETE(`/nodes/:node`, BasicAuth(DeleteNode))
		router.DELETE(`/objstates/:state`, BasicAuth(DeleteObjectState))
		router.DELETE(`/objtypes/:type`, BasicAuth(DeleteObjectType))
		router.DELETE(`/oncall/:oncall`, BasicAuth(DeleteOncall))
		router.DELETE(`/permission/:permission`, BasicAuth(DeletePermission))
		router.DELETE(`/predicates/:predicate`, BasicAuth(DeletePredicate))
		router.DELETE(`/property/custom/:repository/:custom`, BasicAuth(DeleteProperty))
		router.DELETE(`/property/native/:native`, BasicAuth(DeleteProperty))
		router.DELETE(`/property/service/global/:service`, BasicAuth(DeleteProperty))
		router.DELETE(`/property/service/team/:team/:service`, BasicAuth(DeleteProperty))
		router.DELETE(`/property/system/:system`, BasicAuth(DeleteProperty))
		router.DELETE(`/providers/:provider`, BasicAuth(DeleteProvider))
		router.DELETE(`/repository/:repository/property/:type/:source`, BasicAuth(DeletePropertyFromRepository))
		router.DELETE(`/servers/:server`, BasicAuth(DeleteServer))
		router.DELETE(`/status/:status`, BasicAuth(DeleteStatus))
		router.DELETE(`/teams/:team`, BasicAuth(DeleteTeam))
		router.DELETE(`/units/:unit`, BasicAuth(DeleteUnit))
		router.DELETE(`/users/:user`, BasicAuth(DeleteUser))
		router.DELETE(`/validity/:property`, BasicAuth(DeleteValidity))
		router.DELETE(`/views/:view`, BasicAuth(DeleteView))
		router.GET(`/deployments/id/:uuid`, DeliverDeploymentDetails)
		router.GET(`/deployments/monitoring/:uuid/:all`, DeliverMonitoringDeployments)
		router.GET(`/deployments/monitoring/:uuid`, DeliverMonitoringDeployments)
		router.PATCH(`/datacentergroups/:datacentergroup`, BasicAuth(AddDatacenterToGroup))
		router.PATCH(`/deployments/id/:uuid/:result`, UpdateDeploymentDetails)
		router.PATCH(`/oncall/:oncall`, BasicAuth(UpdateOncall))
		router.PATCH(`/views/:view`, BasicAuth(RenameView))
		router.POST(`/attributes/`, BasicAuth(AddAttribute))
		router.POST(`/authenticate/`, AuthenticationKex)
		router.POST(`/buckets/:bucket/property/:type/`, BasicAuth(AddPropertyToBucket))
		router.POST(`/buckets/`, BasicAuth(AddBucket))
		router.POST(`/capability/`, BasicAuth(AddCapability))
		router.POST(`/category/`, BasicAuth(AddCategory))
		router.POST(`/checks/:repository/`, BasicAuth(AddCheckConfiguration))
		router.POST(`/clusters/:cluster/members/`, BasicAuth(AddMemberToCluster))
		router.POST(`/clusters/:cluster/property/:type/`, BasicAuth(AddPropertyToCluster))
		router.POST(`/clusters/`, BasicAuth(AddCluster))
		router.POST(`/datacenters/`, BasicAuth(AddDatacenter))
		router.POST(`/environments/`, BasicAuth(AddEnvironment))
		router.POST(`/grant/global/:rtyp/:rid/`, BasicAuth(GrantGlobalRight))
		router.POST(`/grant/limited/:rtyp/:rid/:scope/:uuid/`, BasicAuth(GrantLimitedRight))
		router.POST(`/grant/system/:rtyp/:rid/`, BasicAuth(GrantSystemRight))
		router.POST(`/groups/:group/members/`, BasicAuth(AddMemberToGroup))
		router.POST(`/groups/:group/property/:type/`, BasicAuth(AddPropertyToGroup))
		router.POST(`/groups/`, BasicAuth(AddGroup))
		router.POST(`/levels/`, BasicAuth(AddLevel))
		router.POST(`/metrics/`, BasicAuth(AddMetric))
		router.POST(`/modes/`, BasicAuth(AddMode))
		router.POST(`/monitoring/`, BasicAuth(AddMonitoring))
		router.POST(`/nodes/:node/property/:type/`, BasicAuth(AddPropertyToNode))
		router.POST(`/nodes/`, BasicAuth(AddNode))
		router.POST(`/objstates/`, BasicAuth(AddObjectState))
		router.POST(`/objtypes/`, BasicAuth(AddObjectType))
		router.POST(`/oncall/`, BasicAuth(AddOncall))
		router.POST(`/permission/`, BasicAuth(AddPermission))
		router.POST(`/predicates/`, BasicAuth(AddPredicate))
		router.POST(`/property/custom/:repository/`, BasicAuth(AddProperty))
		router.POST(`/property/native/`, BasicAuth(AddProperty))
		router.POST(`/property/service/global/`, BasicAuth(AddProperty))
		router.POST(`/property/service/team/:team/`, BasicAuth(AddProperty))
		router.POST(`/property/system/`, BasicAuth(AddProperty))
		router.POST(`/providers/`, BasicAuth(AddProvider))
		router.POST(`/repository/:repository/property/:type/`, BasicAuth(AddPropertyToRepository))
		router.POST(`/repository/`, BasicAuth(AddRepository))
		router.POST(`/servers/:server`, BasicAuth(InsertNullServer))
		router.POST(`/servers/`, BasicAuth(AddServer))
		router.POST(`/status/`, BasicAuth(AddStatus))
		router.POST(`/system/`, BasicAuth(SystemOperation))
		router.POST(`/teams/`, BasicAuth(AddTeam))
		router.POST(`/units/`, BasicAuth(AddUnit))
		router.POST(`/users/`, BasicAuth(AddUser))
		router.POST(`/validity/`, BasicAuth(AddValidity))
		router.POST(`/views/`, BasicAuth(AddView))
		router.PUT(`/authenticate/activate/:uuid`, AuthenticationActivateUser)
		router.PUT(`/authenticate/bootstrap/:uuid`, AuthenticationBootstrapRoot)
		router.PUT(`/authenticate/token/:uuid`, AuthenticationIssueToken)
		router.PUT(`/datacenters/:datacenter`, BasicAuth(RenameDatacenter))
		router.PUT(`/environments/:environment`, BasicAuth(RenameEnvironment))
		router.PUT(`/jobs/:jobid`, BasicAuth(JobDelay))
		router.PUT(`/nodes/:node/config`, BasicAuth(AssignNode))
		router.PUT(`/nodes/:node`, BasicAuth(UpdateNode))
		router.PUT(`/objstates/:state`, BasicAuth(RenameObjectState))
		router.PUT(`/objtypes/:type`, BasicAuth(RenameObjectType))
		router.PUT(`/servers/:server`, BasicAuth(UpdateServer))
		router.PUT(`/teams/:team`, BasicAuth(UpdateTeam))
		router.PUT(`/users/:user`, BasicAuth(UpdateUser))
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
