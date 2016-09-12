package main

import (
	"fmt"
	"log"
	"os"
	"strings"

)

type grimReaper struct {
	system chan msg.Request
}

func (grim *grimReaper) run() {
	// defer calls stack in LIFO order
	defer os.Exit(0)
	defer conn.Close()

runloop:
	for {
		select {
		case req := <-grim.system:
			grim.process(&req)
		}
		break runloop
	}
}

func (grim *grimReaper) process(q *msg.Request) {
	result := msg.Result{Type: `grimReaper`, Action: q.Action}
	switch q.Action {
	case `shutdown`:
	default:
		result.NotImplemented(
			fmt.Errorf("Unknown requested action: %s",
				q.Action),
		)
		q.Reply <- result
		return
	}

	// tell HTTP handlers to start turning people away
	// TODO 900/ShutdownInProgress
	ShutdownInProgress = true

	// answer shutdown request
	result.OK()
	q.Reply <- result

	// I have awoken.
	log.Println(`GRIM REAPER ACTIVATED. SYSTEM SHUTDOWN INITIATED`)

	// stop all treeKeeper       : /^repository_.*/
	for handler, _ := range handlerMap {
		if strings.HasPrefix(handler, `repository_`) {
			handlerMap[handler].(*treeKeeper).stopchan <- true
		}
	}
	// shutdown all treeKeeper   : /^repository_.*/
	for handler, _ := range handlerMap {
		if strings.HasPrefix(handler, `repository_`) {
			handlerMap[handler].(*treeKeeper).shutdown <- true
			delete(handlerMap, handler)
			log.Printf("grimReaper: shut down %s", handler)
		}
	}
	// shutdown all write handler: /WriteHandler$/
	for handler, _ := range handlerMap {
		if !strings.HasSuffix(handler, `WriteHandler`) {
			continue
		}
		switch handlerMap[handler].(type) {
		case *somaViewWriteHandler:
			handlerMap[handler].(*somaViewWriteHandler).
				shutdown <- true
		case *somaEnvironmentWriteHandler:
			handlerMap[handler].(*somaEnvironmentWriteHandler).
				shutdown <- true
		case *somaObjectStateWriteHandler:
			handlerMap[handler].(*somaObjectStateWriteHandler).
				shutdown <- true
		case *somaObjectTypeWriteHandler:
			handlerMap[handler].(*somaObjectTypeWriteHandler).
				shutdown <- true
		case *somaDatacenterWriteHandler:
			handlerMap[handler].(*somaDatacenterWriteHandler).
				shutdown <- true
		default:
			continue
		}
		delete(handlerMap, handler)
		log.Printf("grimReaper: shut down %s", handler)
	}
	// shutdown all read handler : /ReadHandler$/
	for handler, _ := range handlerMap {
		if !strings.HasSuffix(handler, `ReadHandler`) {
			continue
		}
		switch handlerMap[handler].(type) {
		case *somaViewReadHandler:
			handlerMap[handler].(*somaViewReadHandler).
				shutdown <- true
		case *somaEnvironmentReadHandler:
			handlerMap[handler].(*somaEnvironmentReadHandler).
				shutdown <- true
		case *somaObjectStateReadHandler:
			handlerMap[handler].(*somaObjectStateReadHandler).
				shutdown <- true
		case *somaObjectTypeReadHandler:
			handlerMap[handler].(*somaObjectTypeReadHandler).
				shutdown <- true
		case *somaDatacenterReadHandler:
			handlerMap[handler].(*somaDatacenterReadHandler).
				shutdown <- true
		default:
			continue
		}
		delete(handlerMap, handler)
		log.Printf("grimReaper: shut down %s", handler)
	}
	//                             /_r$/
	// shutdown special handlers
	// shutdown supervisor -- needs handling in BasicAuth()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
