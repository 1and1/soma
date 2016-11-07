package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/client9/reopen"
	"github.com/julienschmidt/httprouter"
	metrics "github.com/rcrowley/go-metrics"
)

// global variables
var (
	// main database connection pool
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
	// Orderly shutdown of the system has been called. GrimReaper is active
	ShutdownInProgress bool = false
	// lookup table of logfile handles for logrotate reopen
	logFileMap = make(map[string]*reopen.FileWriter)
	// Global metrics registry
	Metrics = make(map[string]metrics.Registry)
)

const (
	// Format string for millisecond precision RFC3339
	rfc3339Milli string = "2006-01-02T15:04:05.000Z07:00"
	// SOMA version
	SomaVersion string = `0.8.2`
	// Logging format strings
	LogStrReq = `Subsystem=%s, Request=%s, User=%s, Addr=%s`
	LogStrArg = `Subsystem=%s, Request=%s, User=%s, Addr=%s, Arg=%s`
	LogStrOK  = `Subsystem=%s, Result=%s, InternalCode=%d, ExternalCode=%d`
	LogStrErr = `Subsystem=%s, Action=%s, InternalCode=%d, Error=%s`
)

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	var (
		configFlag, configFile, obsRepoFlag string
		noPokeFlag, forcedCorruption        bool
		err                                 error
		appLog, reqLog, errLog              *log.Logger
		lfhGlobal, lfhApp, lfhReq, lfhErr   *reopen.FileWriter
	)

	// Daemon command line flags
	flag.StringVar(&configFlag, "config", "/srv/soma/huxley/conf/soma.conf", "Configuration file location")
	flag.StringVar(&obsRepoFlag, "repo", "", "Single-repository mode target repository")
	flag.BoolVar(&noPokeFlag, "nopoke", false, "Disable lifecycle pokes")
	flag.BoolVar(&forcedCorruption, `allowdatacorruption`, false, `Allow single-repo mode on production`)
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

	// Open logfiles
	if lfhGlobal, err = reopen.NewFileWriter(
		filepath.Join(SomaCfg.LogPath, `global.log`),
	); err != nil {
		log.Fatalf("Unable to open global output log: %s", err)
	}
	log.SetOutput(lfhGlobal)
	logFileMap[`global`] = lfhGlobal

	appLog = log.New()
	if lfhApp, err = reopen.NewFileWriter(
		filepath.Join(SomaCfg.LogPath, `application.log`),
	); err != nil {
		log.Fatalf("Unable to open application log: %s", err)
	}
	appLog.Out = lfhApp
	logFileMap[`application`] = lfhApp

	reqLog = log.New()
	if lfhReq, err = reopen.NewFileWriter(
		filepath.Join(SomaCfg.LogPath, `request.log`),
	); err != nil {
		log.Fatalf("Unable to open request log: %s", err)
	}
	reqLog.Out = lfhReq
	logFileMap[`request`] = lfhReq

	errLog = log.New()
	if lfhErr, err = reopen.NewFileWriter(
		filepath.Join(SomaCfg.LogPath, `error.log`),
	); err != nil {
		log.Fatalf("Unable to open error log: %s", err)
	}
	errLog.Out = lfhErr
	logFileMap[`error`] = lfhErr

	// signal handler will reopen all logfiles on USR2
	sigChanLogRotate := make(chan os.Signal, 1)
	signal.Notify(sigChanLogRotate, syscall.SIGUSR2)
	go logrotate(sigChanLogRotate)

	// print selected runtime mode
	if SomaCfg.ReadOnly {
		appLog.Println(`Instance has been configured as: read-only mode`)
	} else if SomaCfg.Observer {
		appLog.Println(`Instance has been configured as: observer mode`)
	} else {
		appLog.Println(`Instance has been configured as: normal mode`)
	}

	// single-repo cli argument overwrites config file
	if obsRepoFlag != `` {
		SomaCfg.ObserverRepo = obsRepoFlag
	}
	if SomaCfg.ObserverRepo != `` {
		appLog.Printf("Single-repository mode active for: %s", SomaCfg.ObserverRepo)
	}

	// disallow single-repository mode on production r/w instances
	if !SomaCfg.ReadOnly && !SomaCfg.Observer &&
		SomaCfg.ObserverRepo != `` && SomaCfg.Environment == `production` &&
		!forcedCorruption {
		errLog.Fatal(`Single-repository r/w mode disallowed for production environments. ` +
			`Use the -allowdatacorruption flag if you are sure this will be the only ` +
			`running SOMA instance.`)
	}

	if noPokeFlag {
		SomaCfg.NoPoke = true
		appLog.Println(`Instance has disabled outgoing pokes by lifeCycle manager`)
	}

	/*
	 * Register metrics collections
	 */
	Metrics[`golang`] = metrics.NewPrefixedRegistry(`golang.`)
	metrics.RegisterRuntimeMemStats(Metrics[`golang`])
	go metrics.CaptureRuntimeMemStats(Metrics[`golang`], time.Second*60)

	Metrics[`soma`] = metrics.NewPrefixedRegistry(`soma`)
	Metrics[`soma`].Register(`requests.latency`,
		// TODO NewCustomTimer(Histogram, Meter) so there is access
		// to Histogram.Clear()
		metrics.NewTimer())

	/*
	 * Construct listen address
	 */
	SomaCfg.Daemon.url = &url.URL{}
	SomaCfg.Daemon.url.Host = fmt.Sprintf("%s:%s", SomaCfg.Daemon.Listen, SomaCfg.Daemon.Port)
	if SomaCfg.Daemon.Tls {
		SomaCfg.Daemon.url.Scheme = "https"
		if ok, pt := govalidator.IsFilePath(SomaCfg.Daemon.Cert); !ok {
			errLog.Fatal("Missing required certificate configuration config/daemon/cert-file")
		} else {
			if pt != govalidator.Unix {
				errLog.Fatal("config/daemon/cert-File: valid Windows paths are not helpful")
			}
		}
		if ok, pt := govalidator.IsFilePath(SomaCfg.Daemon.Key); !ok {
			errLog.Fatal("Missing required key configuration config/daemon/key-file")
		} else {
			if pt != govalidator.Unix {
				errLog.Fatal("config/daemon/key-file: valid Windows paths are not helpful")
			}
		}
	} else {
		SomaCfg.Daemon.url.Scheme = "http"
	}

	connectToDatabase(appLog, errLog)
	go pingDatabase(errLog)

	startHandlers(appLog, reqLog, errLog)

	router := httprouter.New()

	router.HEAD(`/`, Check(Ping))

	router.GET(`/attributes/:attribute`, Check(BasicAuth(ShowAttribute)))
	router.GET(`/attributes/`, Check(BasicAuth(ListAttribute)))
	router.GET(`/authenticate/validate/`, Check(BasicAuth(AuthenticationValidate)))
	router.GET(`/buckets/:bucket/instances/`, Check(BasicAuth(InstanceList)))
	router.GET(`/buckets/:bucket/tree/:tree`, Check(BasicAuth(OutputTree)))
	router.GET(`/buckets/:bucket`, Check(BasicAuth(ShowBucket)))
	router.GET(`/buckets/`, Check(BasicAuth(ListBucket)))
	router.GET(`/capability/:capability`, Check(BasicAuth(ShowCapability)))
	router.GET(`/capability/`, Check(BasicAuth(ListCapability)))
	router.GET(`/category/:category`, Check(BasicAuth(ShowCategory)))
	router.GET(`/category/`, Check(BasicAuth(ListCategory)))
	router.GET(`/checks/:repository/:check`, Check(BasicAuth(ShowCheckConfiguration)))
	router.GET(`/checks/:repository/`, Check(BasicAuth(ListCheckConfiguration)))
	router.GET(`/clusters/:cluster/instances/`, Check(BasicAuth(InstanceList)))
	router.GET(`/clusters/:cluster/members/`, Check(BasicAuth(ListClusterMembers)))
	router.GET(`/clusters/:cluster/tree/:tree`, Check(BasicAuth(OutputTree)))
	router.GET(`/clusters/:cluster`, Check(BasicAuth(ShowCluster)))
	router.GET(`/clusters/`, Check(BasicAuth(ListCluster)))
	router.GET(`/datacentergroups/:datacentergroup`, Check(BasicAuth(ShowDatacenterGroup)))
	router.GET(`/datacentergroups/`, Check(BasicAuth(ListDatacenterGroups)))
	router.GET(`/datacenters/:datacenter`, Check(BasicAuth(ShowDatacenter)))
	router.GET(`/datacenters/`, Check(BasicAuth(ListDatacenters)))
	router.GET(`/environments/:environment`, Check(BasicAuth(ShowEnvironment)))
	router.GET(`/environments/`, Check(BasicAuth(ListEnvironments)))
	router.GET(`/groups/:group/instances/`, Check(BasicAuth(InstanceList)))
	router.GET(`/groups/:group/members/`, Check(BasicAuth(ListGroupMembers)))
	router.GET(`/groups/:group/tree/:tree`, Check(BasicAuth(OutputTree)))
	router.GET(`/groups/:group`, Check(BasicAuth(ShowGroup)))
	router.GET(`/groups/`, Check(BasicAuth(ListGroup)))
	router.GET(`/hostdeployment/:system/:assetid`, Check(GetHostDeployment))
	router.GET(`/instances/:instance/versions`, Check(BasicAuth(InstanceVersions)))
	router.GET(`/instances/:instance`, Check(BasicAuth(InstanceShow)))
	router.GET(`/instances/`, Check(BasicAuth(InstanceListAll)))
	router.GET(`/jobs/:jobid`, Check(BasicAuth(ShowJob)))
	router.GET(`/jobs/`, Check(BasicAuth(ListJobs)))
	router.GET(`/levels/:level`, Check(BasicAuth(ShowLevel)))
	router.GET(`/levels/`, Check(BasicAuth(ListLevel)))
	router.GET(`/metrics/:metric`, Check(BasicAuth(ShowMetric)))
	router.GET(`/metrics/`, Check(BasicAuth(ListMetric)))
	router.GET(`/modes/:mode`, Check(BasicAuth(ShowMode)))
	router.GET(`/modes/`, Check(BasicAuth(ListMode)))
	router.GET(`/monitoring/:monitoring`, Check(BasicAuth(ShowMonitoring)))
	router.GET(`/monitoring/`, Check(BasicAuth(ListMonitoring)))
	router.GET(`/nodes/:node/config`, Check(BasicAuth(ShowNodeConfig)))
	router.GET(`/nodes/:node/instances/`, Check(BasicAuth(InstanceList)))
	router.GET(`/nodes/:node/tree/:tree`, Check(BasicAuth(OutputTree)))
	router.GET(`/nodes/:node`, Check(BasicAuth(ShowNode)))
	router.GET(`/nodes/`, Check(BasicAuth(ListNode)))
	router.GET(`/objstates/:state`, Check(BasicAuth(ShowObjectState)))
	router.GET(`/objstates/`, Check(BasicAuth(ListObjectStates)))
	router.GET(`/objtypes/:type`, Check(BasicAuth(ShowObjectType)))
	router.GET(`/objtypes/`, Check(BasicAuth(ListObjectTypes)))
	router.GET(`/oncall/:oncall`, Check(BasicAuth(ShowOncall)))
	router.GET(`/oncall/`, Check(BasicAuth(ListOncall)))
	router.GET(`/permission/:permission`, Check(BasicAuth(ShowPermission)))
	router.GET(`/permission/`, Check(BasicAuth(ListPermission)))
	router.GET(`/predicates/:predicate`, Check(BasicAuth(ShowPredicate)))
	router.GET(`/predicates/`, Check(BasicAuth(ListPredicate)))
	router.GET(`/property/custom/:repository/:custom`, Check(BasicAuth(ShowProperty)))
	router.GET(`/property/custom/:repository/`, Check(BasicAuth(ListProperty)))
	router.GET(`/property/native/:native`, Check(BasicAuth(ShowProperty)))
	router.GET(`/property/native/`, Check(BasicAuth(ListProperty)))
	router.GET(`/property/service/global/:service`, Check(BasicAuth(ShowProperty)))
	router.GET(`/property/service/global/`, Check(BasicAuth(ListProperty)))
	router.GET(`/property/service/team/:team/:service`, Check(BasicAuth(ShowProperty)))
	router.GET(`/property/service/team/:team/`, Check(BasicAuth(ListProperty)))
	router.GET(`/property/system/:system`, Check(BasicAuth(ShowProperty)))
	router.GET(`/property/system/`, Check(BasicAuth(ListProperty)))
	router.GET(`/providers/:provider`, Check(BasicAuth(ShowProvider)))
	router.GET(`/providers/`, Check(BasicAuth(ListProvider)))
	router.GET(`/repository/:repository/instances/`, Check(BasicAuth(InstanceList)))
	router.GET(`/repository/:repository/tree/:tree`, Check(BasicAuth(OutputTree)))
	router.GET(`/repository/:repository`, Check(BasicAuth(ShowRepository)))
	router.GET(`/repository/`, Check(BasicAuth(ListRepository)))
	router.GET(`/servers/:server`, Check(BasicAuth(ShowServer)))
	router.GET(`/servers/`, Check(BasicAuth(ListServer)))
	router.GET(`/status/:status`, Check(BasicAuth(ShowStatus)))
	router.GET(`/status/`, Check(BasicAuth(ListStatus)))
	router.GET(`/sync/datacenters/`, Check(BasicAuth(SyncDatacenters)))
	router.GET(`/sync/nodes/`, Check(BasicAuth(SyncNode)))
	router.GET(`/sync/servers/`, Check(BasicAuth(SyncServer)))
	router.GET(`/sync/teams/`, Check(BasicAuth(SyncTeam)))
	router.GET(`/sync/users/`, Check(BasicAuth(SyncUser)))
	router.GET(`/teams/:team`, Check(BasicAuth(ShowTeam)))
	router.GET(`/teams/`, Check(BasicAuth(ListTeam)))
	router.GET(`/units/:unit`, Check(BasicAuth(ShowUnit)))
	router.GET(`/units/`, Check(BasicAuth(ListUnit)))
	router.GET(`/users/:user`, Check(BasicAuth(ShowUser)))
	router.GET(`/users/`, Check(BasicAuth(ListUser)))
	router.GET(`/validity/:property`, Check(BasicAuth(ShowValidity)))
	router.GET(`/validity/`, Check(BasicAuth(ListValidity)))
	router.GET(`/views/:view`, Check(BasicAuth(ShowView)))
	router.GET(`/views/`, Check(BasicAuth(ListView)))
	router.GET(`/workflow/summary`, Check(BasicAuth(WorkflowSummary)))
	router.POST(`/filter/buckets/`, Check(BasicAuth(ListBucket)))
	router.POST(`/filter/capability/`, Check(BasicAuth(ListCapability)))
	router.POST(`/filter/checks/:repository/`, Check(BasicAuth(ListCheckConfiguration)))
	router.POST(`/filter/clusters/`, Check(BasicAuth(ListCluster)))
	router.POST(`/filter/grant/`, Check(BasicAuth(SearchGrant)))
	router.POST(`/filter/groups/`, Check(BasicAuth(ListGroup)))
	router.POST(`/filter/jobs/`, Check(BasicAuth(SearchJob)))
	router.POST(`/filter/levels/`, Check(BasicAuth(ListLevel)))
	router.POST(`/filter/monitoring/`, Check(BasicAuth(ListMonitoring)))
	router.POST(`/filter/nodes/`, Check(BasicAuth(ListNode)))
	router.POST(`/filter/oncall/`, Check(BasicAuth(ListOncall)))
	router.POST(`/filter/permission/`, Check(BasicAuth(SearchPermission)))
	router.POST(`/filter/property/custom/:repository/`, Check(BasicAuth(ListProperty)))
	router.POST(`/filter/property/service/global/`, Check(BasicAuth(ListProperty)))
	router.POST(`/filter/property/service/team/:team/`, Check(BasicAuth(ListProperty)))
	router.POST(`/filter/property/system/`, Check(BasicAuth(ListProperty)))
	router.POST(`/filter/repository/`, Check(BasicAuth(ListRepository)))
	router.POST(`/filter/servers/`, Check(BasicAuth(SearchServer)))
	router.POST(`/filter/teams/`, Check(BasicAuth(ListTeam)))
	router.POST(`/filter/users/`, Check(BasicAuth(ListUser)))
	router.POST(`/filter/workflow/`, Check(BasicAuth(WorkflowList)))
	router.POST(`/hostdeployment/:system/:assetid`, Check(AssembleHostUpdate))

	if !SomaCfg.ReadOnly {
		if !SomaCfg.Observer {
			router.DELETE(`/attributes/:attribute`, Check(BasicAuth(DeleteAttribute)))
			router.DELETE(`/buckets/:bucket/property/:type/:source`, Check(BasicAuth(DeletePropertyFromBucket)))
			router.DELETE(`/capability/:capability`, Check(BasicAuth(DeleteCapability)))
			router.DELETE(`/category/:category`, Check(BasicAuth(DeleteCategory)))
			router.DELETE(`/checks/:repository/:check`, Check(BasicAuth(DeleteCheckConfiguration)))
			router.DELETE(`/clusters/:cluster/property/:type/:source`, Check(BasicAuth(DeletePropertyFromCluster)))
			router.DELETE(`/datacentergroups/:datacentergroup`, Check(BasicAuth(DeleteDatacenterFromGroup)))
			router.DELETE(`/datacenters/:datacenter`, Check(BasicAuth(DeleteDatacenter)))
			router.DELETE(`/environments/:environment`, Check(BasicAuth(DeleteEnvironment)))
			router.DELETE(`/grant/global/:rtyp/:rid/:grant`, Check(BasicAuth(RevokeGlobalRight)))
			router.DELETE(`/grant/limited/:rtyp/:rid/:scope/:uuid/:grant`, Check(BasicAuth(RevokeLimitedRight)))
			router.DELETE(`/grant/system/:rtyp/:rid/:grant`, Check(BasicAuth(RevokeSystemRight)))
			router.DELETE(`/groups/:group/property/:type/:source`, Check(BasicAuth(DeletePropertyFromGroup)))
			router.DELETE(`/levels/:level`, Check(BasicAuth(DeleteLevel)))
			router.DELETE(`/metrics/:metric`, Check(BasicAuth(DeleteMetric)))
			router.DELETE(`/modes/:mode`, Check(BasicAuth(DeleteMode)))
			router.DELETE(`/monitoring/:monitoring`, Check(BasicAuth(DeleteMonitoring)))
			router.DELETE(`/nodes/:node/property/:type/:source`, Check(BasicAuth(DeletePropertyFromNode)))
			router.DELETE(`/nodes/:node`, Check(BasicAuth(DeleteNode)))
			router.DELETE(`/objstates/:state`, Check(BasicAuth(DeleteObjectState)))
			router.DELETE(`/objtypes/:type`, Check(BasicAuth(DeleteObjectType)))
			router.DELETE(`/oncall/:oncall`, Check(BasicAuth(DeleteOncall)))
			router.DELETE(`/permission/:permission`, Check(BasicAuth(DeletePermission)))
			router.DELETE(`/predicates/:predicate`, Check(BasicAuth(DeletePredicate)))
			router.DELETE(`/property/custom/:repository/:custom`, Check(BasicAuth(DeleteProperty)))
			router.DELETE(`/property/native/:native`, Check(BasicAuth(DeleteProperty)))
			router.DELETE(`/property/service/global/:service`, Check(BasicAuth(DeleteProperty)))
			router.DELETE(`/property/service/team/:team/:service`, Check(BasicAuth(DeleteProperty)))
			router.DELETE(`/property/system/:system`, Check(BasicAuth(DeleteProperty)))
			router.DELETE(`/providers/:provider`, Check(BasicAuth(DeleteProvider)))
			router.DELETE(`/repository/:repository/property/:type/:source`, Check(BasicAuth(DeletePropertyFromRepository)))
			router.DELETE(`/servers/:server`, Check(BasicAuth(DeleteServer)))
			router.DELETE(`/status/:status`, Check(BasicAuth(DeleteStatus)))
			router.DELETE(`/teams/:team`, Check(BasicAuth(DeleteTeam)))
			router.DELETE(`/units/:unit`, Check(BasicAuth(DeleteUnit)))
			router.DELETE(`/users/:user`, Check(BasicAuth(DeleteUser)))
			router.DELETE(`/validity/:property`, Check(BasicAuth(DeleteValidity)))
			router.DELETE(`/views/:view`, Check(BasicAuth(DeleteView)))
			router.GET(`/deployments/id/:uuid`, Check(DeliverDeploymentDetails))
			router.GET(`/deployments/monitoring/:uuid/:all`, Check(DeliverMonitoringDeployments))
			router.GET(`/deployments/monitoring/:uuid`, Check(DeliverMonitoringDeployments))
			router.PATCH(`/authenticate/user/password/:uuid`, Check(AuthenticationChangeUserPassword))
			router.PATCH(`/datacentergroups/:datacentergroup`, Check(BasicAuth(AddDatacenterToGroup)))
			router.PATCH(`/deployments/id/:uuid/:result`, Check(UpdateDeploymentDetails))
			router.PATCH(`/oncall/:oncall`, Check(BasicAuth(UpdateOncall)))
			router.PATCH(`/views/:view`, Check(BasicAuth(RenameView)))
			router.POST(`/attributes/`, Check(BasicAuth(AddAttribute)))
			router.POST(`/buckets/:bucket/property/:type/`, Check(BasicAuth(AddPropertyToBucket)))
			router.POST(`/buckets/`, Check(BasicAuth(AddBucket)))
			router.POST(`/capability/`, Check(BasicAuth(AddCapability)))
			router.POST(`/category/`, Check(BasicAuth(AddCategory)))
			router.POST(`/checks/:repository/`, Check(BasicAuth(AddCheckConfiguration)))
			router.POST(`/clusters/:cluster/members/`, Check(BasicAuth(AddMemberToCluster)))
			router.POST(`/clusters/:cluster/property/:type/`, Check(BasicAuth(AddPropertyToCluster)))
			router.POST(`/clusters/`, Check(BasicAuth(AddCluster)))
			router.POST(`/datacenters/`, Check(BasicAuth(AddDatacenter)))
			router.POST(`/environments/`, Check(BasicAuth(AddEnvironment)))
			router.POST(`/grant/global/:rtyp/:rid/`, Check(BasicAuth(GrantGlobalRight)))
			router.POST(`/grant/limited/:rtyp/:rid/:scope/:uuid/`, Check(BasicAuth(GrantLimitedRight)))
			router.POST(`/grant/system/:rtyp/:rid/`, Check(BasicAuth(GrantSystemRight)))
			router.POST(`/groups/:group/members/`, Check(BasicAuth(AddMemberToGroup)))
			router.POST(`/groups/:group/property/:type/`, Check(BasicAuth(AddPropertyToGroup)))
			router.POST(`/groups/`, Check(BasicAuth(AddGroup)))
			router.POST(`/levels/`, Check(BasicAuth(AddLevel)))
			router.POST(`/metrics/`, Check(BasicAuth(AddMetric)))
			router.POST(`/modes/`, Check(BasicAuth(AddMode)))
			router.POST(`/monitoring/`, Check(BasicAuth(AddMonitoring)))
			router.POST(`/nodes/:node/property/:type/`, Check(BasicAuth(AddPropertyToNode)))
			router.POST(`/nodes/`, Check(BasicAuth(AddNode)))
			router.POST(`/objstates/`, Check(BasicAuth(AddObjectState)))
			router.POST(`/objtypes/`, Check(BasicAuth(AddObjectType)))
			router.POST(`/oncall/`, Check(BasicAuth(AddOncall)))
			router.POST(`/permission/`, Check(BasicAuth(AddPermission)))
			router.POST(`/predicates/`, Check(BasicAuth(AddPredicate)))
			router.POST(`/property/custom/:repository/`, Check(BasicAuth(AddProperty)))
			router.POST(`/property/native/`, Check(BasicAuth(AddProperty)))
			router.POST(`/property/service/global/`, Check(BasicAuth(AddProperty)))
			router.POST(`/property/service/team/:team/`, Check(BasicAuth(AddProperty)))
			router.POST(`/property/system/`, Check(BasicAuth(AddProperty)))
			router.POST(`/providers/`, Check(BasicAuth(AddProvider)))
			router.POST(`/repository/:repository/property/:type/`, Check(BasicAuth(AddPropertyToRepository)))
			router.POST(`/repository/`, Check(BasicAuth(AddRepository)))
			router.POST(`/servers/:server`, Check(BasicAuth(InsertNullServer)))
			router.POST(`/servers/`, Check(BasicAuth(AddServer)))
			router.POST(`/status/`, Check(BasicAuth(AddStatus)))
			router.POST(`/system/`, Check(BasicAuth(SystemOperation)))
			router.POST(`/teams/`, Check(BasicAuth(AddTeam)))
			router.POST(`/units/`, Check(BasicAuth(AddUnit)))
			router.POST(`/users/`, Check(BasicAuth(AddUser)))
			router.POST(`/validity/`, Check(BasicAuth(AddValidity)))
			router.POST(`/views/`, Check(BasicAuth(AddView)))
			router.PUT(`/authenticate/activate/:uuid`, Check(AuthenticationActivateUser))
			router.PUT(`/authenticate/bootstrap/:uuid`, Check(AuthenticationBootstrapRoot))
			router.PUT(`/authenticate/user/password/:uuid`, Check(AuthenticationResetUserPassword))
			router.PUT(`/datacenters/:datacenter`, Check(BasicAuth(RenameDatacenter)))
			router.PUT(`/environments/:environment`, Check(BasicAuth(RenameEnvironment)))
			router.PUT(`/jobs/:jobid`, Check(BasicAuth(JobDelay)))
			router.PUT(`/nodes/:node/config`, Check(BasicAuth(AssignNode)))
			router.PUT(`/nodes/:node`, Check(BasicAuth(UpdateNode)))
			router.PUT(`/objstates/:state`, Check(BasicAuth(RenameObjectState)))
			router.PUT(`/objtypes/:type`, Check(BasicAuth(RenameObjectType)))
			router.PUT(`/servers/:server`, Check(BasicAuth(UpdateServer)))
			router.PUT(`/teams/:team`, Check(BasicAuth(UpdateTeam)))
			router.PUT(`/users/:user`, Check(BasicAuth(UpdateUser)))
		}
		router.POST(`/authenticate/`, Check(AuthenticationKex))
		router.PUT(`/authenticate/token/:uuid`, Check(AuthenticationIssueToken))
	}

	if SomaCfg.Daemon.Tls {
		errLog.Fatal(http.ListenAndServeTLS(
			SomaCfg.Daemon.url.Host,
			SomaCfg.Daemon.Cert,
			SomaCfg.Daemon.Key,
			router))
	} else {
		errLog.Fatal(http.ListenAndServe(SomaCfg.Daemon.url.Host, router))
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
