package main

func startHandlers() {
	spawnViewReadHandler()
	spawnEnvironmentReadHandler()
	spawnObjectStateReadHandler()
	spawnObjectTypeReadHandler()
	spawnDatacenterReadHandler()
	spawnLevelReadHandler()
	spawnPredicateReadHandler()
	spawnStatusReadHandler()
	spawnOncallReadHandler()
	spawnTeamReadHandler()
	spawnNodeReadHandler()
	spawnServerReadHandler()
	spawnUnitReadHandler()
	spawnProviderReadHandler()
	spawnMetricReadHandler()
	spawnModeReadHandler()
	spawnUserReadHandler()
	spawnMonitoringReadHandler()
	spawnCapabilityReadHandler()
	spawnPropertyReadHandler()
	spawnAttributeReadHandler()
	spawnRepositoryReadHandler()
	spawnBucketReadHandler()
	spawnGroupReadHandler()
	spawnClusterReadHandler()
	spawnCheckConfigurationReadHandler()
	spawnHostDeploymentHandler()

	if !SomaCfg.ReadOnly {
		spawnViewWriteHandler()
		spawnEnvironmentWriteHandler()
		spawnObjectStateWriteHandler()
		spawnObjectTypeWriteHandler()
		spawnDatacenterWriteHandler()
		spawnLevelWriteHandler()
		spawnPredicateWriteHandler()
		spawnStatusWriteHandler()
		spawnOncallWriteHandler()
		spawnTeamWriteHandler()
		spawnNodeWriteHandler()
		spawnServerWriteHandler()
		spawnUnitWriteHandler()
		spawnProviderWriteHandler()
		spawnMetricWriteHandler()
		spawnModeWriteHandler()
		spawnUserWriteHandler()
		spawnMonitoringWriteHandler()
		spawnCapabilityWriteHandler()
		spawnPropertyWriteHandler()
		spawnAttributeWriteHandler()
		spawnForestCustodian()
		spawnGuidePost()
		spawnLifeCycle()
		spawnDeploymentHandler()
	}
}

func spawnViewReadHandler() {
	var viewReadHandler somaViewReadHandler
	viewReadHandler.input = make(chan somaViewRequest)
	viewReadHandler.shutdown = make(chan bool)
	viewReadHandler.conn = conn
	handlerMap["viewReadHandler"] = viewReadHandler
	go viewReadHandler.run()
}

func spawnViewWriteHandler() {
	var viewWriteHandler somaViewWriteHandler
	viewWriteHandler.input = make(chan somaViewRequest, 64)
	viewWriteHandler.shutdown = make(chan bool)
	viewWriteHandler.conn = conn
	handlerMap["viewWriteHandler"] = viewWriteHandler
	go viewWriteHandler.run()
}

func spawnEnvironmentReadHandler() {
	var environmentReadHandler somaEnvironmentReadHandler
	environmentReadHandler.input = make(chan somaEnvironmentRequest)
	environmentReadHandler.shutdown = make(chan bool)
	environmentReadHandler.conn = conn
	handlerMap["environmentReadHandler"] = environmentReadHandler
	go environmentReadHandler.run()
}

func spawnEnvironmentWriteHandler() {
	var environmentWriteHandler somaEnvironmentWriteHandler
	environmentWriteHandler.input = make(chan somaEnvironmentRequest, 64)
	environmentWriteHandler.shutdown = make(chan bool)
	environmentWriteHandler.conn = conn
	handlerMap["environmentWriteHandler"] = environmentWriteHandler
	go environmentWriteHandler.run()
}

func spawnObjectStateReadHandler() {
	var objectStateReadHandler somaObjectStateReadHandler
	objectStateReadHandler.input = make(chan somaObjectStateRequest)
	objectStateReadHandler.shutdown = make(chan bool)
	objectStateReadHandler.conn = conn
	handlerMap["objectStateReadHandler"] = objectStateReadHandler
	go objectStateReadHandler.run()
}

func spawnObjectStateWriteHandler() {
	var objectStateWriteHandler somaObjectStateWriteHandler
	objectStateWriteHandler.input = make(chan somaObjectStateRequest, 64)
	objectStateWriteHandler.shutdown = make(chan bool)
	objectStateWriteHandler.conn = conn
	handlerMap["objectStateWriteHandler"] = objectStateWriteHandler
	go objectStateWriteHandler.run()
}

func spawnObjectTypeReadHandler() {
	var objectTypeReadHandler somaObjectTypeReadHandler
	objectTypeReadHandler.input = make(chan somaObjectTypeRequest)
	objectTypeReadHandler.shutdown = make(chan bool)
	objectTypeReadHandler.conn = conn
	handlerMap["objectTypeReadHandler"] = objectTypeReadHandler
	go objectTypeReadHandler.run()
}

func spawnObjectTypeWriteHandler() {
	var objectTypeWriteHandler somaObjectTypeWriteHandler
	objectTypeWriteHandler.input = make(chan somaObjectTypeRequest, 64)
	objectTypeWriteHandler.shutdown = make(chan bool)
	objectTypeWriteHandler.conn = conn
	handlerMap["objectTypeWriteHandler"] = objectTypeWriteHandler
	go objectTypeWriteHandler.run()
}

func spawnDatacenterReadHandler() {
	var datacenterReadHandler somaDatacenterReadHandler
	datacenterReadHandler.input = make(chan somaDatacenterRequest)
	datacenterReadHandler.shutdown = make(chan bool)
	datacenterReadHandler.conn = conn
	handlerMap["datacenterReadHandler"] = datacenterReadHandler
	go datacenterReadHandler.run()
}

func spawnDatacenterWriteHandler() {
	var datacenterWriteHandler somaDatacenterWriteHandler
	datacenterWriteHandler.input = make(chan somaDatacenterRequest, 64)
	datacenterWriteHandler.shutdown = make(chan bool)
	datacenterWriteHandler.conn = conn
	handlerMap["datacenterWriteHandler"] = datacenterWriteHandler
	go datacenterWriteHandler.run()
}

func spawnLevelReadHandler() {
	var levelReadHandler somaLevelReadHandler
	levelReadHandler.input = make(chan somaLevelRequest, 64)
	levelReadHandler.shutdown = make(chan bool)
	levelReadHandler.conn = conn
	handlerMap["levelReadHandler"] = levelReadHandler
	go levelReadHandler.run()
}

func spawnLevelWriteHandler() {
	var levelWriteHandler somaLevelWriteHandler
	levelWriteHandler.input = make(chan somaLevelRequest, 64)
	levelWriteHandler.shutdown = make(chan bool)
	levelWriteHandler.conn = conn
	handlerMap["levelWriteHandler"] = levelWriteHandler
	go levelWriteHandler.run()
}

func spawnPredicateReadHandler() {
	var predicateReadHandler somaPredicateReadHandler
	predicateReadHandler.input = make(chan somaPredicateRequest, 64)
	predicateReadHandler.shutdown = make(chan bool)
	predicateReadHandler.conn = conn
	handlerMap["predicateReadHandler"] = predicateReadHandler
	go predicateReadHandler.run()
}

func spawnPredicateWriteHandler() {
	var predicateWriteHandler somaPredicateWriteHandler
	predicateWriteHandler.input = make(chan somaPredicateRequest, 64)
	predicateWriteHandler.shutdown = make(chan bool)
	predicateWriteHandler.conn = conn
	handlerMap["predicateWriteHandler"] = predicateWriteHandler
	go predicateWriteHandler.run()
}

func spawnStatusReadHandler() {
	var statusReadHandler somaStatusReadHandler
	statusReadHandler.input = make(chan somaStatusRequest, 64)
	statusReadHandler.shutdown = make(chan bool)
	statusReadHandler.conn = conn
	handlerMap["statusReadHandler"] = statusReadHandler
	go statusReadHandler.run()
}

func spawnStatusWriteHandler() {
	var statusWriteHandler somaStatusWriteHandler
	statusWriteHandler.input = make(chan somaStatusRequest, 64)
	statusWriteHandler.shutdown = make(chan bool)
	statusWriteHandler.conn = conn
	handlerMap["statusWriteHandler"] = statusWriteHandler
	go statusWriteHandler.run()
}

func spawnOncallReadHandler() {
	var oncallReadHandler somaOncallReadHandler
	oncallReadHandler.input = make(chan somaOncallRequest, 64)
	oncallReadHandler.shutdown = make(chan bool)
	oncallReadHandler.conn = conn
	handlerMap["oncallReadHandler"] = oncallReadHandler
	go oncallReadHandler.run()
}

func spawnOncallWriteHandler() {
	var oncallWriteHandler somaOncallWriteHandler
	oncallWriteHandler.input = make(chan somaOncallRequest, 64)
	oncallWriteHandler.shutdown = make(chan bool)
	oncallWriteHandler.conn = conn
	handlerMap["oncallWriteHandler"] = oncallWriteHandler
	go oncallWriteHandler.run()
}

func spawnTeamReadHandler() {
	var teamReadHandler somaTeamReadHandler
	teamReadHandler.input = make(chan somaTeamRequest, 64)
	teamReadHandler.shutdown = make(chan bool)
	teamReadHandler.conn = conn
	handlerMap["teamReadHandler"] = teamReadHandler
	go teamReadHandler.run()
}

func spawnTeamWriteHandler() {
	var teamWriteHandler somaTeamWriteHandler
	teamWriteHandler.input = make(chan somaTeamRequest, 64)
	teamWriteHandler.shutdown = make(chan bool)
	teamWriteHandler.conn = conn
	handlerMap["teamWriteHandler"] = teamWriteHandler
	go teamWriteHandler.run()
}

func spawnNodeReadHandler() {
	var nodeReadHandler somaNodeReadHandler
	nodeReadHandler.input = make(chan somaNodeRequest, 64)
	nodeReadHandler.shutdown = make(chan bool)
	nodeReadHandler.conn = conn
	handlerMap["nodeReadHandler"] = nodeReadHandler
	go nodeReadHandler.run()
}

func spawnNodeWriteHandler() {
	var nodeWriteHandler somaNodeWriteHandler
	nodeWriteHandler.input = make(chan somaNodeRequest, 64)
	nodeWriteHandler.shutdown = make(chan bool)
	nodeWriteHandler.conn = conn
	handlerMap["nodeWriteHandler"] = nodeWriteHandler
	go nodeWriteHandler.run()
}

func spawnServerReadHandler() {
	var serverReadHandler somaServerReadHandler
	serverReadHandler.input = make(chan somaServerRequest, 64)
	serverReadHandler.shutdown = make(chan bool)
	serverReadHandler.conn = conn
	handlerMap["serverReadHandler"] = serverReadHandler
	go serverReadHandler.run()
}

func spawnServerWriteHandler() {
	var serverWriteHandler somaServerWriteHandler
	serverWriteHandler.input = make(chan somaServerRequest, 64)
	serverWriteHandler.shutdown = make(chan bool)
	serverWriteHandler.conn = conn
	handlerMap["serverWriteHandler"] = serverWriteHandler
	go serverWriteHandler.run()
}

func spawnUnitReadHandler() {
	var unitReadHandler somaUnitReadHandler
	unitReadHandler.input = make(chan somaUnitRequest, 64)
	unitReadHandler.shutdown = make(chan bool)
	unitReadHandler.conn = conn
	handlerMap["unitReadHandler"] = unitReadHandler
	go unitReadHandler.run()
}

func spawnUnitWriteHandler() {
	var unitWriteHandler somaUnitWriteHandler
	unitWriteHandler.input = make(chan somaUnitRequest, 64)
	unitWriteHandler.shutdown = make(chan bool)
	unitWriteHandler.conn = conn
	handlerMap["unitWriteHandler"] = unitWriteHandler
	go unitWriteHandler.run()
}

func spawnProviderReadHandler() {
	var providerReadHandler somaProviderReadHandler
	providerReadHandler.input = make(chan somaProviderRequest, 64)
	providerReadHandler.shutdown = make(chan bool)
	providerReadHandler.conn = conn
	handlerMap["providerReadHandler"] = providerReadHandler
	go providerReadHandler.run()
}

func spawnProviderWriteHandler() {
	var providerWriteHandler somaProviderWriteHandler
	providerWriteHandler.input = make(chan somaProviderRequest, 64)
	providerWriteHandler.shutdown = make(chan bool)
	providerWriteHandler.conn = conn
	handlerMap["providerWriteHandler"] = providerWriteHandler
	go providerWriteHandler.run()
}

func spawnMetricReadHandler() {
	var metricReadHandler somaMetricReadHandler
	metricReadHandler.input = make(chan somaMetricRequest, 64)
	metricReadHandler.shutdown = make(chan bool)
	metricReadHandler.conn = conn
	handlerMap["metricReadHandler"] = metricReadHandler
	go metricReadHandler.run()
}

func spawnMetricWriteHandler() {
	var metricWriteHandler somaMetricWriteHandler
	metricWriteHandler.input = make(chan somaMetricRequest, 64)
	metricWriteHandler.shutdown = make(chan bool)
	metricWriteHandler.conn = conn
	handlerMap["metricWriteHandler"] = metricWriteHandler
	go metricWriteHandler.run()
}

func spawnModeReadHandler() {
	var modeReadHandler somaModeReadHandler
	modeReadHandler.input = make(chan somaModeRequest, 64)
	modeReadHandler.shutdown = make(chan bool)
	modeReadHandler.conn = conn
	handlerMap["modeReadHandler"] = modeReadHandler
	go modeReadHandler.run()
}

func spawnModeWriteHandler() {
	var modeWriteHandler somaModeWriteHandler
	modeWriteHandler.input = make(chan somaModeRequest, 64)
	modeWriteHandler.shutdown = make(chan bool)
	modeWriteHandler.conn = conn
	handlerMap["modeWriteHandler"] = modeWriteHandler
	go modeWriteHandler.run()
}

func spawnUserReadHandler() {
	var userReadHandler somaUserReadHandler
	userReadHandler.input = make(chan somaUserRequest, 64)
	userReadHandler.shutdown = make(chan bool)
	userReadHandler.conn = conn
	handlerMap["userReadHandler"] = userReadHandler
	go userReadHandler.run()
}

func spawnUserWriteHandler() {
	var userWriteHandler somaUserWriteHandler
	userWriteHandler.input = make(chan somaUserRequest, 64)
	userWriteHandler.shutdown = make(chan bool)
	userWriteHandler.conn = conn
	handlerMap["userWriteHandler"] = userWriteHandler
	go userWriteHandler.run()
}

func spawnMonitoringReadHandler() {
	var monitoringReadHandler somaMonitoringReadHandler
	monitoringReadHandler.input = make(chan somaMonitoringRequest, 64)
	monitoringReadHandler.shutdown = make(chan bool)
	monitoringReadHandler.conn = conn
	handlerMap["monitoringReadHandler"] = monitoringReadHandler
	go monitoringReadHandler.run()
}

func spawnMonitoringWriteHandler() {
	var monitoringWriteHandler somaMonitoringWriteHandler
	monitoringWriteHandler.input = make(chan somaMonitoringRequest, 64)
	monitoringWriteHandler.shutdown = make(chan bool)
	monitoringWriteHandler.conn = conn
	handlerMap["monitoringWriteHandler"] = monitoringWriteHandler
	go monitoringWriteHandler.run()
}

func spawnCapabilityReadHandler() {
	var capabilityReadHandler somaCapabilityReadHandler
	capabilityReadHandler.input = make(chan somaCapabilityRequest, 64)
	capabilityReadHandler.shutdown = make(chan bool)
	capabilityReadHandler.conn = conn
	handlerMap["capabilityReadHandler"] = capabilityReadHandler
	go capabilityReadHandler.run()
}

func spawnCapabilityWriteHandler() {
	var capabilityWriteHandler somaCapabilityWriteHandler
	capabilityWriteHandler.input = make(chan somaCapabilityRequest, 64)
	capabilityWriteHandler.shutdown = make(chan bool)
	capabilityWriteHandler.conn = conn
	handlerMap["capabilityWriteHandler"] = capabilityWriteHandler
	go capabilityWriteHandler.run()
}

func spawnPropertyReadHandler() {
	var propertyReadHandler somaPropertyReadHandler
	propertyReadHandler.input = make(chan somaPropertyRequest, 64)
	propertyReadHandler.shutdown = make(chan bool)
	propertyReadHandler.conn = conn
	handlerMap["propertyReadHandler"] = propertyReadHandler
	go propertyReadHandler.run()
}

func spawnPropertyWriteHandler() {
	var propertyWriteHandler somaPropertyWriteHandler
	propertyWriteHandler.input = make(chan somaPropertyRequest, 64)
	propertyWriteHandler.shutdown = make(chan bool)
	propertyWriteHandler.conn = conn
	handlerMap["propertyWriteHandler"] = propertyWriteHandler
	go propertyWriteHandler.run()
}

func spawnAttributeReadHandler() {
	var attributeReadHandler somaAttributeReadHandler
	attributeReadHandler.input = make(chan somaAttributeRequest, 64)
	attributeReadHandler.shutdown = make(chan bool)
	attributeReadHandler.conn = conn
	handlerMap["attributeReadHandler"] = attributeReadHandler
	go attributeReadHandler.run()
}

func spawnAttributeWriteHandler() {
	var attributeWriteHandler somaAttributeWriteHandler
	attributeWriteHandler.input = make(chan somaAttributeRequest, 64)
	attributeWriteHandler.shutdown = make(chan bool)
	attributeWriteHandler.conn = conn
	handlerMap["attributeWriteHandler"] = attributeWriteHandler
	go attributeWriteHandler.run()
}

func spawnRepositoryReadHandler() {
	var repositoryReadHandler somaRepositoryReadHandler
	repositoryReadHandler.input = make(chan somaRepositoryRequest, 64)
	repositoryReadHandler.shutdown = make(chan bool)
	repositoryReadHandler.conn = conn
	handlerMap["repositoryReadHandler"] = repositoryReadHandler
	go repositoryReadHandler.run()
}

func spawnBucketReadHandler() {
	var bucketReadHandler somaBucketReadHandler
	bucketReadHandler.input = make(chan somaBucketRequest, 64)
	bucketReadHandler.shutdown = make(chan bool)
	bucketReadHandler.conn = conn
	handlerMap["bucketReadHandler"] = bucketReadHandler
	go bucketReadHandler.run()
}

func spawnGroupReadHandler() {
	var groupReadHandler somaGroupReadHandler
	groupReadHandler.input = make(chan somaGroupRequest, 64)
	groupReadHandler.shutdown = make(chan bool)
	groupReadHandler.conn = conn
	handlerMap["groupReadHandler"] = groupReadHandler
	go groupReadHandler.run()
}

func spawnClusterReadHandler() {
	var clusterReadHandler somaClusterReadHandler
	clusterReadHandler.input = make(chan somaClusterRequest, 64)
	clusterReadHandler.shutdown = make(chan bool)
	clusterReadHandler.conn = conn
	handlerMap["clusterReadHandler"] = clusterReadHandler
	go clusterReadHandler.run()
}

func spawnForestCustodian() {
	var fC forestCustodian
	fC.input = make(chan somaRepositoryRequest, 64)
	fC.shutdown = make(chan bool)
	fC.conn = conn
	handlerMap["forestCustodian"] = fC
	go fC.run()
}

func spawnGuidePost() {
	var gP guidePost
	gP.input = make(chan treeRequest, 4096)
	gP.shutdown = make(chan bool)
	gP.conn = conn
	handlerMap["guidePost"] = gP
	go gP.run()
}

func spawnCheckConfigurationReadHandler() {
	var checkConfigurationReadHandler somaCheckConfigurationReadHandler
	checkConfigurationReadHandler.input = make(chan somaCheckConfigRequest, 64)
	checkConfigurationReadHandler.shutdown = make(chan bool)
	checkConfigurationReadHandler.conn = conn
	handlerMap["checkConfigurationReadHandler"] = checkConfigurationReadHandler
	go checkConfigurationReadHandler.run()
}

func spawnLifeCycle() {
	var lifeCycleHandler lifeCycle
	lifeCycleHandler.shutdown = make(chan bool)
	lifeCycleHandler.conn = conn
	handlerMap["lifeCycle"] = lifeCycleHandler
	go lifeCycleHandler.run()
}

func spawnDeploymentHandler() {
	var deploymentHandler somaDeploymentHandler
	deploymentHandler.input = make(chan somaDeploymentRequest, 64)
	deploymentHandler.shutdown = make(chan bool)
	deploymentHandler.conn = conn
	handlerMap["deploymentHandler"] = deploymentHandler
	go deploymentHandler.run()
}

func spawnHostDeploymentHandler() {
	var hostDeploymentHandler somaHostDeploymentHandler
	hostDeploymentHandler.input = make(chan somaHostDeploymentRequest, 64)
	hostDeploymentHandler.shutdown = make(chan bool)
	hostDeploymentHandler.conn = conn
	handlerMap["hostDeploymentHandler"] = hostDeploymentHandler
	go hostDeploymentHandler.run()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
