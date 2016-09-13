package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

)

type grimReaper struct {
	system chan msg.Request
	conn   *sql.DB
}

func (grim *grimReaper) run() {
	// defer calls stack in LIFO order
	defer os.Exit(0)
	defer grim.conn.Close()

runloop:
	for {
		select {
		case req := <-grim.system:
			grim.process(&req)
		}
		break runloop
	}

	time.Sleep(time.Duration(SomaCfg.ShutdownDelay) * time.Second)
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
	ShutdownInProgress = true

	// answer shutdown request
	result.OK()
	q.Reply <- result

	time.Sleep(time.Duration(SomaCfg.ShutdownDelay) * time.Second)

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
		`guidePost`,
		`lifeCycle`,
		`deploymentHandler`,
		`hostDeploymentHandler`,
	} {
		handlerMap[h].(Downer).shutdownNow()
		delete(handlerMap, h)
		log.Printf("grimReaper: shut down %s", h)
	}

	// shutdown supervisor -- needs handling in BasicAuth()
	handlerMap[`supervisor`].(Downer).shutdownNow()
	delete(handlerMap, `supervisor`)
	log.Println(`grimReaper: shut down the supervisor`)

	// log what we have missed
	log.Println(`grimReaper: checking for still running handlers`)
	for name, _ := range handlerMap {
		if name == `grimReaper` {
			continue
		}
		log.Printf("grimReaper: %s is still running\n", name)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
