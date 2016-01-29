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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
