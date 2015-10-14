package main

func startHandlers() {
	spawnViewReadHandler()
	spawnViewWriteHandler()
	spawnEnvironmentReadHandler()
	spawnEnvironmentWriteHandler()
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
