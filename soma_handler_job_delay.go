package main

import "time"

type jobDelay struct {
	input    chan waitSpec
	shutdown chan bool
	notify   chan string
	waitList map[string][]waitSpec
	jobDone  map[string]time.Time
}

type waitSpec struct {
	JobId string
	RecvT time.Time
	Reply chan bool
}

func (j *jobDelay) run() {
	tock := time.Tick(1 * time.Minute)
	j.jobDone = make(map[string]time.Time)
	j.waitList = make(map[string][]waitSpec)

runloop:
	for {
		select {
		case <-j.shutdown:
			break runloop
		case jid := <-j.notify:
			j.jobDone[jid] = time.Now().UTC()
			for _, ws := range j.waitList[jid] {
				close(ws.Reply)
			}
			delete(j.waitList, jid)
		case ws := <-j.input:
			if _, ok := j.jobDone[ws.JobId]; ok {
				close(ws.Reply)
			} else {
				ws.RecvT = time.Now().UTC()
				if _, ok := j.waitList[ws.JobId]; !ok {
					j.waitList[ws.JobId] = []waitSpec{}
				}
				j.waitList[ws.JobId] = append(j.waitList[ws.JobId], ws)
			}
		case <-tock:
			// clean jobDone cache
			for jid, _ := range j.jobDone {
				if time.Since(j.jobDone[jid]) > (2 * time.Hour) {
					delete(j.jobDone, jid)
				}
			}
			// clean waiters
			for jid, _ := range j.waitList {
				newList := []waitSpec{}
				for i, _ := range j.waitList[jid] {
					if time.Since(j.waitList[jid][i].RecvT) < (5 * time.Minute) {
						newList = append(newList, j.waitList[jid][i])
					} else {
						close(j.waitList[jid][i].Reply)
					}
				}
				if len(newList) > 0 {
					j.waitList[jid] = newList
				} else {
					delete(j.waitList, jid)
				}
			}
		}
	}
}

/* Ops Access
 */
func (j *jobDelay) shutdownNow() {
	j.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
