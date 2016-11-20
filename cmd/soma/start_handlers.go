package main

import (
	"encoding/hex"

	"github.com/1and1/soma/internal/msg"
	log "github.com/Sirupsen/logrus"
)

func startHandlers(appLog, reqLog, errLog *log.Logger) {
	spawnGrimReaperHandler(appLog, reqLog, errLog)
	spawnSupervisorHandler(appLog, reqLog, errLog)

	spawnAttributeRead(appLog, reqLog, errLog)
	spawnBucketReadHandler(appLog, reqLog, errLog)
	spawnCapabilityReadHandler(appLog, reqLog, errLog)
	spawnCheckConfigurationReadHandler(appLog, reqLog, errLog)
	spawnClusterReadHandler(appLog, reqLog, errLog)
	spawnDatacenterReadHandler(appLog, reqLog, errLog)
	spawnEntityRead(appLog, reqLog, errLog)
	spawnEnvironmentReadHandler(appLog, reqLog, errLog)
	spawnGroupReadHandler(appLog, reqLog, errLog)
	spawnHostDeploymentHandler(appLog, reqLog, errLog)
	spawnInstanceReadHandler(appLog, reqLog, errLog)
	spawnJobReadHandler(appLog, reqLog, errLog)
	spawnLevelReadHandler(appLog, reqLog, errLog)
	spawnMetricReadHandler(appLog, reqLog, errLog)
	spawnModeReadHandler(appLog, reqLog, errLog)
	spawnMonitoringRead(appLog, reqLog, errLog)
	spawnNodeReadHandler(appLog, reqLog, errLog)
	spawnObjectStateReadHandler(appLog, reqLog, errLog)
	spawnOncallReadHandler(appLog, reqLog, errLog)
	spawnOutputTreeHandler(appLog, reqLog, errLog)
	spawnPredicateReadHandler(appLog, reqLog, errLog)
	spawnPropertyReadHandler(appLog, reqLog, errLog)
	spawnProviderReadHandler(appLog, reqLog, errLog)
	spawnRepositoryReadHandler(appLog, reqLog, errLog)
	spawnServerReadHandler(appLog, reqLog, errLog)
	spawnStatusReadHandler(appLog, reqLog, errLog)
	spawnTeamReadHandler(appLog, reqLog, errLog)
	spawnUnitReadHandler(appLog, reqLog, errLog)
	spawnUserReadHandler(appLog, reqLog, errLog)
	spawnValidityReadHandler(appLog, reqLog, errLog)
	spawnViewReadHandler(appLog, reqLog, errLog)
	spawnWorkflowReadHandler(appLog, reqLog, errLog)

	if !SomaCfg.ReadOnly {
		spawnForestCustodian(appLog, reqLog, errLog)
		spawnGuidePost(appLog, reqLog, errLog)
		spawnLifeCycle(appLog, reqLog, errLog)

		if !SomaCfg.Observer {
			spawnAttributeWrite(appLog, reqLog, errLog)
			spawnCapabilityWriteHandler(appLog, reqLog, errLog)
			spawnDatacenterWriteHandler(appLog, reqLog, errLog)
			spawnDeploymentHandler(appLog, reqLog, errLog)
			spawnEntityWrite(appLog, reqLog, errLog)
			spawnEnvironmentWriteHandler(appLog, reqLog, errLog)
			spawnJobDelay(appLog, reqLog, errLog)
			spawnLevelWriteHandler(appLog, reqLog, errLog)
			spawnMetricWriteHandler(appLog, reqLog, errLog)
			spawnModeWriteHandler(appLog, reqLog, errLog)
			spawnMonitoringWrite(appLog, reqLog, errLog)
			spawnNodeWriteHandler(appLog, reqLog, errLog)
			spawnObjectStateWriteHandler(appLog, reqLog, errLog)
			spawnOncallWriteHandler(appLog, reqLog, errLog)
			spawnPredicateWriteHandler(appLog, reqLog, errLog)
			spawnPropertyWriteHandler(appLog, reqLog, errLog)
			spawnProviderWriteHandler(appLog, reqLog, errLog)
			spawnServerWriteHandler(appLog, reqLog, errLog)
			spawnStatusWriteHandler(appLog, reqLog, errLog)
			spawnTeamWriteHandler(appLog, reqLog, errLog)
			spawnUnitWriteHandler(appLog, reqLog, errLog)
			spawnUserWriteHandler(appLog, reqLog, errLog)
			spawnValidityWriteHandler(appLog, reqLog, errLog)
			spawnViewWriteHandler(appLog, reqLog, errLog)
			spawnWorkflowWriteHandler(appLog, reqLog, errLog)
		}
	}
}

func spawnViewReadHandler(appLog, reqLog, errLog *log.Logger) {
	var viewReadHandler somaViewReadHandler
	viewReadHandler.input = make(chan somaViewRequest)
	viewReadHandler.shutdown = make(chan bool)
	viewReadHandler.conn = conn
	viewReadHandler.appLog = appLog
	viewReadHandler.reqLog = reqLog
	viewReadHandler.errLog = errLog
	handlerMap["viewReadHandler"] = &viewReadHandler
	go viewReadHandler.run()
}

func spawnViewWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var viewWriteHandler somaViewWriteHandler
	viewWriteHandler.input = make(chan somaViewRequest, 64)
	viewWriteHandler.shutdown = make(chan bool)
	viewWriteHandler.conn = conn
	viewWriteHandler.appLog = appLog
	viewWriteHandler.reqLog = reqLog
	viewWriteHandler.errLog = errLog
	handlerMap["viewWriteHandler"] = &viewWriteHandler
	go viewWriteHandler.run()
}

func spawnEnvironmentReadHandler(appLog, reqLog, errLog *log.Logger) {
	var environmentReadHandler environmentRead
	environmentReadHandler.input = make(chan msg.Request)
	environmentReadHandler.shutdown = make(chan bool)
	environmentReadHandler.conn = conn
	environmentReadHandler.appLog = appLog
	environmentReadHandler.reqLog = reqLog
	environmentReadHandler.errLog = errLog
	handlerMap[`environment_r`] = &environmentReadHandler
	go environmentReadHandler.run()
}

func spawnEnvironmentWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var environmentWriteHandler environmentWrite
	environmentWriteHandler.input = make(chan msg.Request, 64)
	environmentWriteHandler.shutdown = make(chan bool)
	environmentWriteHandler.conn = conn
	environmentWriteHandler.appLog = appLog
	environmentWriteHandler.reqLog = reqLog
	environmentWriteHandler.errLog = errLog
	handlerMap[`environment_w`] = &environmentWriteHandler
	go environmentWriteHandler.run()
}

func spawnObjectStateReadHandler(appLog, reqLog, errLog *log.Logger) {
	var handler stateRead
	handler.input = make(chan msg.Request)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`state_r`] = &handler
	go handler.run()
}

func spawnObjectStateWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var handler stateWrite
	handler.input = make(chan msg.Request, 64)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`state_w`] = &handler
	go handler.run()
}

func spawnEntityRead(appLog, reqLog, errLog *log.Logger) {
	var handler entityRead
	handler.input = make(chan msg.Request)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`entity_r`] = &handler
	go handler.run()
}

func spawnEntityWrite(appLog, reqLog, errLog *log.Logger) {
	var handler entityWrite
	handler.input = make(chan msg.Request, 64)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`entity_w`] = &handler
	go handler.run()
}

func spawnDatacenterReadHandler(appLog, reqLog, errLog *log.Logger) {
	var datacenterReadHandler somaDatacenterReadHandler
	datacenterReadHandler.input = make(chan somaDatacenterRequest)
	datacenterReadHandler.shutdown = make(chan bool)
	datacenterReadHandler.conn = conn
	datacenterReadHandler.appLog = appLog
	datacenterReadHandler.reqLog = reqLog
	datacenterReadHandler.errLog = errLog
	handlerMap["datacenterReadHandler"] = &datacenterReadHandler
	go datacenterReadHandler.run()
}

func spawnDatacenterWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var datacenterWriteHandler somaDatacenterWriteHandler
	datacenterWriteHandler.input = make(chan somaDatacenterRequest, 64)
	datacenterWriteHandler.shutdown = make(chan bool)
	datacenterWriteHandler.conn = conn
	datacenterWriteHandler.appLog = appLog
	datacenterWriteHandler.reqLog = reqLog
	datacenterWriteHandler.errLog = errLog
	handlerMap["datacenterWriteHandler"] = &datacenterWriteHandler
	go datacenterWriteHandler.run()
}

func spawnLevelReadHandler(appLog, reqLog, errLog *log.Logger) {
	var levelReadHandler somaLevelReadHandler
	levelReadHandler.input = make(chan somaLevelRequest, 64)
	levelReadHandler.shutdown = make(chan bool)
	levelReadHandler.conn = conn
	levelReadHandler.appLog = appLog
	levelReadHandler.reqLog = reqLog
	levelReadHandler.errLog = errLog
	handlerMap["levelReadHandler"] = &levelReadHandler
	go levelReadHandler.run()
}

func spawnLevelWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var levelWriteHandler somaLevelWriteHandler
	levelWriteHandler.input = make(chan somaLevelRequest, 64)
	levelWriteHandler.shutdown = make(chan bool)
	levelWriteHandler.conn = conn
	levelWriteHandler.appLog = appLog
	levelWriteHandler.reqLog = reqLog
	levelWriteHandler.errLog = errLog
	handlerMap["levelWriteHandler"] = &levelWriteHandler
	go levelWriteHandler.run()
}

func spawnPredicateReadHandler(appLog, reqLog, errLog *log.Logger) {
	var predicateReadHandler somaPredicateReadHandler
	predicateReadHandler.input = make(chan somaPredicateRequest, 64)
	predicateReadHandler.shutdown = make(chan bool)
	predicateReadHandler.conn = conn
	predicateReadHandler.appLog = appLog
	predicateReadHandler.reqLog = reqLog
	predicateReadHandler.errLog = errLog
	handlerMap["predicateReadHandler"] = &predicateReadHandler
	go predicateReadHandler.run()
}

func spawnPredicateWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var predicateWriteHandler somaPredicateWriteHandler
	predicateWriteHandler.input = make(chan somaPredicateRequest, 64)
	predicateWriteHandler.shutdown = make(chan bool)
	predicateWriteHandler.conn = conn
	predicateWriteHandler.appLog = appLog
	predicateWriteHandler.reqLog = reqLog
	predicateWriteHandler.errLog = errLog
	handlerMap["predicateWriteHandler"] = &predicateWriteHandler
	go predicateWriteHandler.run()
}

func spawnStatusReadHandler(appLog, reqLog, errLog *log.Logger) {
	var statusReadHandler somaStatusReadHandler
	statusReadHandler.input = make(chan somaStatusRequest, 64)
	statusReadHandler.shutdown = make(chan bool)
	statusReadHandler.conn = conn
	statusReadHandler.appLog = appLog
	statusReadHandler.reqLog = reqLog
	statusReadHandler.errLog = errLog
	handlerMap["statusReadHandler"] = &statusReadHandler
	go statusReadHandler.run()
}

func spawnStatusWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var statusWriteHandler somaStatusWriteHandler
	statusWriteHandler.input = make(chan somaStatusRequest, 64)
	statusWriteHandler.shutdown = make(chan bool)
	statusWriteHandler.conn = conn
	statusWriteHandler.appLog = appLog
	statusWriteHandler.reqLog = reqLog
	statusWriteHandler.errLog = errLog
	handlerMap["statusWriteHandler"] = &statusWriteHandler
	go statusWriteHandler.run()
}

func spawnOncallReadHandler(appLog, reqLog, errLog *log.Logger) {
	var oncallReadHandler somaOncallReadHandler
	oncallReadHandler.input = make(chan somaOncallRequest, 64)
	oncallReadHandler.shutdown = make(chan bool)
	oncallReadHandler.conn = conn
	oncallReadHandler.appLog = appLog
	oncallReadHandler.reqLog = reqLog
	oncallReadHandler.errLog = errLog
	handlerMap["oncallReadHandler"] = &oncallReadHandler
	go oncallReadHandler.run()
}

func spawnOncallWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var oncallWriteHandler somaOncallWriteHandler
	oncallWriteHandler.input = make(chan somaOncallRequest, 64)
	oncallWriteHandler.shutdown = make(chan bool)
	oncallWriteHandler.conn = conn
	oncallWriteHandler.appLog = appLog
	oncallWriteHandler.reqLog = reqLog
	oncallWriteHandler.errLog = errLog
	handlerMap["oncallWriteHandler"] = &oncallWriteHandler
	go oncallWriteHandler.run()
}

func spawnTeamReadHandler(appLog, reqLog, errLog *log.Logger) {
	var teamReadHandler somaTeamReadHandler
	teamReadHandler.input = make(chan somaTeamRequest, 64)
	teamReadHandler.shutdown = make(chan bool)
	teamReadHandler.conn = conn
	teamReadHandler.appLog = appLog
	teamReadHandler.reqLog = reqLog
	teamReadHandler.errLog = errLog
	handlerMap["teamReadHandler"] = &teamReadHandler
	go teamReadHandler.run()
}

func spawnTeamWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var teamWriteHandler somaTeamWriteHandler
	teamWriteHandler.input = make(chan somaTeamRequest, 64)
	teamWriteHandler.shutdown = make(chan bool)
	teamWriteHandler.conn = conn
	teamWriteHandler.appLog = appLog
	teamWriteHandler.reqLog = reqLog
	teamWriteHandler.errLog = errLog
	handlerMap["teamWriteHandler"] = &teamWriteHandler
	go teamWriteHandler.run()
}

func spawnNodeReadHandler(appLog, reqLog, errLog *log.Logger) {
	var nodeReadHandler somaNodeReadHandler
	nodeReadHandler.input = make(chan somaNodeRequest, 64)
	nodeReadHandler.shutdown = make(chan bool)
	nodeReadHandler.conn = conn
	nodeReadHandler.appLog = appLog
	nodeReadHandler.reqLog = reqLog
	nodeReadHandler.errLog = errLog
	handlerMap["nodeReadHandler"] = &nodeReadHandler
	go nodeReadHandler.run()
}

func spawnNodeWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var nodeWriteHandler somaNodeWriteHandler
	nodeWriteHandler.input = make(chan somaNodeRequest, 64)
	nodeWriteHandler.shutdown = make(chan bool)
	nodeWriteHandler.conn = conn
	nodeWriteHandler.appLog = appLog
	nodeWriteHandler.reqLog = reqLog
	nodeWriteHandler.errLog = errLog
	handlerMap["nodeWriteHandler"] = &nodeWriteHandler
	go nodeWriteHandler.run()
}

func spawnServerReadHandler(appLog, reqLog, errLog *log.Logger) {
	var serverReadHandler somaServerReadHandler
	serverReadHandler.input = make(chan somaServerRequest, 64)
	serverReadHandler.shutdown = make(chan bool)
	serverReadHandler.conn = conn
	serverReadHandler.appLog = appLog
	serverReadHandler.reqLog = reqLog
	serverReadHandler.errLog = errLog
	handlerMap["serverReadHandler"] = &serverReadHandler
	go serverReadHandler.run()
}

func spawnServerWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var serverWriteHandler somaServerWriteHandler
	serverWriteHandler.input = make(chan somaServerRequest, 64)
	serverWriteHandler.shutdown = make(chan bool)
	serverWriteHandler.conn = conn
	serverWriteHandler.appLog = appLog
	serverWriteHandler.reqLog = reqLog
	serverWriteHandler.errLog = errLog
	handlerMap["serverWriteHandler"] = &serverWriteHandler
	go serverWriteHandler.run()
}

func spawnUnitReadHandler(appLog, reqLog, errLog *log.Logger) {
	var unitReadHandler somaUnitReadHandler
	unitReadHandler.input = make(chan somaUnitRequest, 64)
	unitReadHandler.shutdown = make(chan bool)
	unitReadHandler.conn = conn
	unitReadHandler.appLog = appLog
	unitReadHandler.reqLog = reqLog
	unitReadHandler.errLog = errLog
	handlerMap["unitReadHandler"] = &unitReadHandler
	go unitReadHandler.run()
}

func spawnUnitWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var unitWriteHandler somaUnitWriteHandler
	unitWriteHandler.input = make(chan somaUnitRequest, 64)
	unitWriteHandler.shutdown = make(chan bool)
	unitWriteHandler.conn = conn
	unitWriteHandler.appLog = appLog
	unitWriteHandler.reqLog = reqLog
	unitWriteHandler.errLog = errLog
	handlerMap["unitWriteHandler"] = &unitWriteHandler
	go unitWriteHandler.run()
}

func spawnProviderReadHandler(appLog, reqLog, errLog *log.Logger) {
	var providerReadHandler somaProviderReadHandler
	providerReadHandler.input = make(chan somaProviderRequest, 64)
	providerReadHandler.shutdown = make(chan bool)
	providerReadHandler.conn = conn
	providerReadHandler.appLog = appLog
	providerReadHandler.reqLog = reqLog
	providerReadHandler.errLog = errLog
	handlerMap["providerReadHandler"] = &providerReadHandler
	go providerReadHandler.run()
}

func spawnProviderWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var providerWriteHandler somaProviderWriteHandler
	providerWriteHandler.input = make(chan somaProviderRequest, 64)
	providerWriteHandler.shutdown = make(chan bool)
	providerWriteHandler.conn = conn
	providerWriteHandler.appLog = appLog
	providerWriteHandler.reqLog = reqLog
	providerWriteHandler.errLog = errLog
	handlerMap["providerWriteHandler"] = &providerWriteHandler
	go providerWriteHandler.run()
}

func spawnMetricReadHandler(appLog, reqLog, errLog *log.Logger) {
	var metricReadHandler somaMetricReadHandler
	metricReadHandler.input = make(chan somaMetricRequest, 64)
	metricReadHandler.shutdown = make(chan bool)
	metricReadHandler.conn = conn
	metricReadHandler.appLog = appLog
	metricReadHandler.reqLog = reqLog
	metricReadHandler.errLog = errLog
	handlerMap["metricReadHandler"] = &metricReadHandler
	go metricReadHandler.run()
}

func spawnMetricWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var metricWriteHandler somaMetricWriteHandler
	metricWriteHandler.input = make(chan somaMetricRequest, 64)
	metricWriteHandler.shutdown = make(chan bool)
	metricWriteHandler.conn = conn
	metricWriteHandler.appLog = appLog
	metricWriteHandler.reqLog = reqLog
	metricWriteHandler.errLog = errLog
	handlerMap["metricWriteHandler"] = &metricWriteHandler
	go metricWriteHandler.run()
}

func spawnModeReadHandler(appLog, reqLog, errLog *log.Logger) {
	var modeReadHandler somaModeReadHandler
	modeReadHandler.input = make(chan somaModeRequest, 64)
	modeReadHandler.shutdown = make(chan bool)
	modeReadHandler.conn = conn
	modeReadHandler.appLog = appLog
	modeReadHandler.reqLog = reqLog
	modeReadHandler.errLog = errLog
	handlerMap["modeReadHandler"] = &modeReadHandler
	go modeReadHandler.run()
}

func spawnModeWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var modeWriteHandler somaModeWriteHandler
	modeWriteHandler.input = make(chan somaModeRequest, 64)
	modeWriteHandler.shutdown = make(chan bool)
	modeWriteHandler.conn = conn
	modeWriteHandler.appLog = appLog
	modeWriteHandler.reqLog = reqLog
	modeWriteHandler.errLog = errLog
	handlerMap["modeWriteHandler"] = &modeWriteHandler
	go modeWriteHandler.run()
}

func spawnUserReadHandler(appLog, reqLog, errLog *log.Logger) {
	var userReadHandler somaUserReadHandler
	userReadHandler.input = make(chan somaUserRequest, 64)
	userReadHandler.shutdown = make(chan bool)
	userReadHandler.conn = conn
	userReadHandler.appLog = appLog
	userReadHandler.reqLog = reqLog
	userReadHandler.errLog = errLog
	handlerMap["userReadHandler"] = &userReadHandler
	go userReadHandler.run()
}

func spawnUserWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var userWriteHandler somaUserWriteHandler
	userWriteHandler.input = make(chan somaUserRequest, 64)
	userWriteHandler.shutdown = make(chan bool)
	userWriteHandler.conn = conn
	userWriteHandler.appLog = appLog
	userWriteHandler.reqLog = reqLog
	userWriteHandler.errLog = errLog
	handlerMap["userWriteHandler"] = &userWriteHandler
	go userWriteHandler.run()
}

func spawnMonitoringRead(appLog, reqLog, errLog *log.Logger) {
	var handler monitoringRead
	handler.input = make(chan msg.Request, 64)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`monitoring_r`] = &handler
	go handler.run()
}

func spawnMonitoringWrite(appLog, reqLog, errLog *log.Logger) {
	var handler monitoringWrite
	handler.input = make(chan msg.Request, 64)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`monitoring_w`] = &handler
	go handler.run()
}

func spawnCapabilityReadHandler(appLog, reqLog, errLog *log.Logger) {
	var capabilityReadHandler somaCapabilityReadHandler
	capabilityReadHandler.input = make(chan somaCapabilityRequest, 64)
	capabilityReadHandler.shutdown = make(chan bool)
	capabilityReadHandler.conn = conn
	capabilityReadHandler.appLog = appLog
	capabilityReadHandler.reqLog = reqLog
	capabilityReadHandler.errLog = errLog
	handlerMap["capabilityReadHandler"] = &capabilityReadHandler
	go capabilityReadHandler.run()
}

func spawnCapabilityWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var capabilityWriteHandler somaCapabilityWriteHandler
	capabilityWriteHandler.input = make(chan somaCapabilityRequest, 64)
	capabilityWriteHandler.shutdown = make(chan bool)
	capabilityWriteHandler.conn = conn
	capabilityWriteHandler.appLog = appLog
	capabilityWriteHandler.reqLog = reqLog
	capabilityWriteHandler.errLog = errLog
	handlerMap["capabilityWriteHandler"] = &capabilityWriteHandler
	go capabilityWriteHandler.run()
}

func spawnPropertyReadHandler(appLog, reqLog, errLog *log.Logger) {
	var propertyReadHandler somaPropertyReadHandler
	propertyReadHandler.input = make(chan somaPropertyRequest, 64)
	propertyReadHandler.shutdown = make(chan bool)
	propertyReadHandler.conn = conn
	propertyReadHandler.appLog = appLog
	propertyReadHandler.reqLog = reqLog
	propertyReadHandler.errLog = errLog
	handlerMap["propertyReadHandler"] = &propertyReadHandler
	go propertyReadHandler.run()
}

func spawnPropertyWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var propertyWriteHandler somaPropertyWriteHandler
	propertyWriteHandler.input = make(chan somaPropertyRequest, 64)
	propertyWriteHandler.shutdown = make(chan bool)
	propertyWriteHandler.conn = conn
	propertyWriteHandler.appLog = appLog
	propertyWriteHandler.reqLog = reqLog
	propertyWriteHandler.errLog = errLog
	handlerMap["propertyWriteHandler"] = &propertyWriteHandler
	go propertyWriteHandler.run()
}

func spawnAttributeRead(appLog, reqLog, errLog *log.Logger) {
	var handler attributeRead
	handler.input = make(chan msg.Request, 64)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`attribute_r`] = &handler
	go handler.run()
}

func spawnAttributeWrite(appLog, reqLog, errLog *log.Logger) {
	var handler attributeWrite
	handler.input = make(chan msg.Request, 64)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`attribute_w`] = &handler
	go handler.run()
}

func spawnRepositoryReadHandler(appLog, reqLog, errLog *log.Logger) {
	var repositoryReadHandler somaRepositoryReadHandler
	repositoryReadHandler.input = make(chan somaRepositoryRequest, 64)
	repositoryReadHandler.shutdown = make(chan bool)
	repositoryReadHandler.conn = conn
	repositoryReadHandler.appLog = appLog
	repositoryReadHandler.reqLog = reqLog
	repositoryReadHandler.errLog = errLog
	handlerMap["repositoryReadHandler"] = &repositoryReadHandler
	go repositoryReadHandler.run()
}

func spawnBucketReadHandler(appLog, reqLog, errLog *log.Logger) {
	var bucketReadHandler somaBucketReadHandler
	bucketReadHandler.input = make(chan somaBucketRequest, 64)
	bucketReadHandler.shutdown = make(chan bool)
	bucketReadHandler.conn = conn
	bucketReadHandler.appLog = appLog
	bucketReadHandler.reqLog = reqLog
	bucketReadHandler.errLog = errLog
	handlerMap["bucketReadHandler"] = &bucketReadHandler
	go bucketReadHandler.run()
}

func spawnGroupReadHandler(appLog, reqLog, errLog *log.Logger) {
	var groupReadHandler somaGroupReadHandler
	groupReadHandler.input = make(chan somaGroupRequest, 64)
	groupReadHandler.shutdown = make(chan bool)
	groupReadHandler.conn = conn
	groupReadHandler.appLog = appLog
	groupReadHandler.reqLog = reqLog
	groupReadHandler.errLog = errLog
	handlerMap["groupReadHandler"] = &groupReadHandler
	go groupReadHandler.run()
}

func spawnClusterReadHandler(appLog, reqLog, errLog *log.Logger) {
	var clusterReadHandler somaClusterReadHandler
	clusterReadHandler.input = make(chan somaClusterRequest, 64)
	clusterReadHandler.shutdown = make(chan bool)
	clusterReadHandler.conn = conn
	clusterReadHandler.appLog = appLog
	clusterReadHandler.reqLog = reqLog
	clusterReadHandler.errLog = errLog
	handlerMap["clusterReadHandler"] = &clusterReadHandler
	go clusterReadHandler.run()
}

func spawnForestCustodian(appLog, reqLog, errLog *log.Logger) {
	var fC forestCustodian
	fC.input = make(chan somaRepositoryRequest, 64)
	fC.system = make(chan msg.Request, 32)
	fC.shutdown = make(chan bool)
	fC.conn = conn
	fC.appLog = appLog
	fC.reqLog = reqLog
	fC.errLog = errLog
	handlerMap["forestCustodian"] = &fC
	go fC.run()
}

func spawnGuidePost(appLog, reqLog, errLog *log.Logger) {
	var gP guidePost
	gP.input = make(chan treeRequest, 4096)
	gP.system = make(chan msg.Request, 32)
	gP.shutdown = make(chan bool)
	gP.conn = conn
	gP.appLog = appLog
	gP.reqLog = reqLog
	gP.errLog = errLog
	handlerMap["guidePost"] = &gP
	go gP.run()
}

func spawnCheckConfigurationReadHandler(appLog, reqLog, errLog *log.Logger) {
	var checkConfigurationReadHandler somaCheckConfigurationReadHandler
	checkConfigurationReadHandler.input = make(chan somaCheckConfigRequest, 64)
	checkConfigurationReadHandler.shutdown = make(chan bool)
	checkConfigurationReadHandler.conn = conn
	checkConfigurationReadHandler.appLog = appLog
	checkConfigurationReadHandler.reqLog = reqLog
	checkConfigurationReadHandler.errLog = errLog
	handlerMap["checkConfigurationReadHandler"] = &checkConfigurationReadHandler
	go checkConfigurationReadHandler.run()
}

func spawnLifeCycle(appLog, reqLog, errLog *log.Logger) {
	var lifeCycleHandler lifeCycle
	lifeCycleHandler.shutdown = make(chan bool)
	lifeCycleHandler.conn = conn
	lifeCycleHandler.appLog = appLog
	lifeCycleHandler.reqLog = reqLog
	lifeCycleHandler.errLog = errLog
	handlerMap["lifeCycle"] = &lifeCycleHandler
	go lifeCycleHandler.run()
}

func spawnDeploymentHandler(appLog, reqLog, errLog *log.Logger) {
	var deploymentHandler somaDeploymentHandler
	deploymentHandler.input = make(chan somaDeploymentRequest, 64)
	deploymentHandler.shutdown = make(chan bool)
	deploymentHandler.conn = conn
	deploymentHandler.appLog = appLog
	deploymentHandler.reqLog = reqLog
	deploymentHandler.errLog = errLog
	handlerMap["deploymentHandler"] = &deploymentHandler
	go deploymentHandler.run()
}

func spawnHostDeploymentHandler(appLog, reqLog, errLog *log.Logger) {
	var hostDeploymentHandler somaHostDeploymentHandler
	hostDeploymentHandler.input = make(chan somaHostDeploymentRequest, 64)
	hostDeploymentHandler.shutdown = make(chan bool)
	hostDeploymentHandler.conn = conn
	hostDeploymentHandler.appLog = appLog
	hostDeploymentHandler.reqLog = reqLog
	hostDeploymentHandler.errLog = errLog
	handlerMap["hostDeploymentHandler"] = &hostDeploymentHandler
	go hostDeploymentHandler.run()
}

func spawnValidityReadHandler(appLog, reqLog, errLog *log.Logger) {
	var validityReadHandler somaValidityReadHandler
	validityReadHandler.input = make(chan somaValidityRequest, 64)
	validityReadHandler.shutdown = make(chan bool)
	validityReadHandler.conn = conn
	validityReadHandler.appLog = appLog
	validityReadHandler.reqLog = reqLog
	validityReadHandler.errLog = errLog
	handlerMap["validityReadHandler"] = &validityReadHandler
	go validityReadHandler.run()
}

func spawnValidityWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var validityWriteHandler somaValidityWriteHandler
	validityWriteHandler.input = make(chan somaValidityRequest, 64)
	validityWriteHandler.shutdown = make(chan bool)
	validityWriteHandler.conn = conn
	validityWriteHandler.appLog = appLog
	validityWriteHandler.reqLog = reqLog
	validityWriteHandler.errLog = errLog
	handlerMap["validityWriteHandler"] = &validityWriteHandler
	go validityWriteHandler.run()
}

func spawnSupervisorHandler(appLog, reqLog, errLog *log.Logger) {
	var supervisorHandler supervisor
	var err error
	supervisorHandler.input = make(chan msg.Request, 1024)
	supervisorHandler.shutdown = make(chan bool)
	supervisorHandler.conn = conn
	supervisorHandler.appLog = appLog
	supervisorHandler.reqLog = reqLog
	supervisorHandler.errLog = errLog
	supervisorHandler.readonly = SomaCfg.ReadOnly
	if supervisorHandler.seed, err = hex.DecodeString(SomaCfg.Auth.TokenSeed); err != nil {
		panic(err)
	}
	if len(supervisorHandler.seed) == 0 {
		panic(`token.seed has length 0`)
	}
	if supervisorHandler.key, err = hex.DecodeString(SomaCfg.Auth.TokenKey); err != nil {
		panic(err)
	}
	if len(supervisorHandler.key) == 0 {
		panic(`token.key has length 0`)
	}
	supervisorHandler.tokenExpiry = SomaCfg.Auth.TokenExpirySeconds
	supervisorHandler.kexExpiry = SomaCfg.Auth.KexExpirySeconds
	supervisorHandler.credExpiry = SomaCfg.Auth.CredentialExpiryDays
	supervisorHandler.activation = SomaCfg.Auth.Activation
	handlerMap[`supervisor`] = &supervisorHandler
	go supervisorHandler.run()
}

func spawnJobDelay(appLog, reqLog, errLog *log.Logger) {
	var handler jobDelay
	handler.input = make(chan waitSpec, 128)
	handler.shutdown = make(chan bool)
	handler.notify = make(chan string, 256)
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`jobDelay`] = &handler
	go handler.run()
}

func spawnJobReadHandler(appLog, reqLog, errLog *log.Logger) {
	var handler jobsRead
	handler.input = make(chan msg.Request, 256)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`jobs_r`] = &handler
	go handler.run()
}

func spawnOutputTreeHandler(appLog, reqLog, errLog *log.Logger) {
	var handler outputTree
	handler.input = make(chan msg.Request, 128)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`tree_r`] = &handler
	go handler.run()
}

func spawnGrimReaperHandler(appLog, reqLog, errLog *log.Logger) {
	var reaper grimReaper
	reaper.system = make(chan msg.Request, 1)
	reaper.conn = conn
	reaper.appLog = appLog
	reaper.reqLog = reqLog
	reaper.errLog = errLog
	handlerMap[`grimReaper`] = &reaper
	go reaper.run()
}

func spawnInstanceReadHandler(appLog, reqLog, errLog *log.Logger) {
	var handler instance
	handler.input = make(chan msg.Request, 128)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`instance_r`] = &handler
	go handler.run()
}

func spawnWorkflowReadHandler(appLog, reqLog, errLog *log.Logger) {
	var handler workflowRead
	handler.input = make(chan msg.Request, 128)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`workflow_r`] = &handler
	go handler.run()
}

func spawnWorkflowWriteHandler(appLog, reqLog, errLog *log.Logger) {
	var handler workflowWrite
	handler.input = make(chan msg.Request, 128)
	handler.shutdown = make(chan bool)
	handler.conn = conn
	handler.appLog = appLog
	handler.reqLog = reqLog
	handler.errLog = errLog
	handlerMap[`workflow_w`] = &handler
	go handler.run()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
