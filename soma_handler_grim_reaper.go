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
			handlerMap[handler].(Stopper).stopNow()
		}
	}
	// shutdown all treeKeeper   : /^repository_.*/
	for handler, _ := range handlerMap {
		if strings.HasPrefix(handler, `repository_`) {
			handlerMap[handler].(Downer).shutdownNow()
			delete(handlerMap, handler)
			log.Printf("grimReaper: shut down %s", handler)
		}
	}
	// shutdown all write handler: /WriteHandler$/
	for handler, _ := range handlerMap {
		if !strings.HasSuffix(handler, `WriteHandler`) {
			continue
		}
		handlerMap[handler].(Downer).shutdownNow()
		delete(handlerMap, handler)
		log.Printf("grimReaper: shut down %s", handler)
	}
	// shutdown all read handler : /ReadHandler$/
	for handler, _ := range handlerMap {
		if !(strings.HasSuffix(handler, `ReadHandler`) ||
			strings.HasSuffix(handler, `_r`)) {
			continue
		}
		handlerMap[handler].(Downer).shutdownNow()
		delete(handlerMap, handler)
		log.Printf("grimReaper: shut down %s", handler)
	}
	// shutdown special handlers
	for _, h := range []string{
		`jobDelay`,
		`forestCustodian`,
	} {
		handlerMap[h].(Downer).shutdownNow()
		delete(handlerMap, handler)
		log.Printf("grimReaper: shut down %s", handler)
	}

	// shutdown supervisor -- needs handling in BasicAuth()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
