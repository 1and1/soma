/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"math"
	"time"

	"github.com/1and1/soma/internal/stmt"
	"github.com/Sirupsen/logrus"
	"gopkg.in/resty.v0"
)

// LifeCycle handles the check rollout workflow
type LifeCycle struct {
	Shutdown          chan struct{}
	conn              *sql.DB
	tick              <-chan time.Time
	stmtUnblock       *sql.Stmt
	stmtPoke          *sql.Stmt
	stmtClear         *sql.Stmt
	stmtDeleteBlocked *sql.Stmt
	stmtDeleteActive  *sql.Stmt
	stmtDeadlock      *sql.Stmt
	stmtReschedule    *sql.Stmt
	stmtSetNotify     *sql.Stmt
	appLog            *logrus.Logger
	reqLog            *logrus.Logger
	errLog            *logrus.Logger
	pokers            map[string]chan string
	soma              *Soma
}

// PokeMessage is the JSON that is sent to notify a monitoring system
// about an available update.
type PokeMessage struct {
	UUID string `json:"uuid"`
	// path should be used to tell the client system the basepath
	// where to get it so SOMA + path + item_id === complete_url
	Path string `json:"path"`
}

// newLifeCycle returns a new LifeCycle handler
func newLifeCycle(s *Soma) (l *LifeCycle) {
	l = &LifeCycle{}
	l.Shutdown = make(chan struct{})
	l.soma = s
	return
}

// register initializes resources provided by the Soma app
func (lc *LifeCycle) register(c *sql.DB, l ...*logrus.Logger) {
	lc.conn = c
	lc.appLog = l[0]
	lc.reqLog = l[1]
	lc.errLog = l[2]
}

// run is the loop for LifeCycle
func (lc *LifeCycle) run() {
	var err error
	lc.pokers = make(map[string]chan string)

	lc.tick = time.NewTicker(
		time.Duration(lc.soma.conf.LifeCycleTick) * time.Second,
	).C

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.LifecycleActiveUnblockCondition:           lc.stmtUnblock,
		stmt.LifecycleReadyDeployments:                 lc.stmtPoke,
		stmt.LifecycleClearUpdateFlag:                  lc.stmtClear,
		stmt.LifecycleBlockedConfigsForDeletedInstance: lc.stmtDeleteBlocked,
		stmt.LifecycleDeprovisionDeletedActive:         lc.stmtDeleteActive,
		stmt.LifecycleDeadLockResolver:                 lc.stmtDeadlock,
		stmt.LifecycleRescheduleDeployments:            lc.stmtReschedule,
		stmt.LifecycleSetNotified:                      lc.stmtSetNotify,
	} {
		if prepStmt, err = lc.conn.Prepare(statement); err != nil {
			lc.errLog.Fatal(`lifecycle`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

	if lc.soma.conf.Observer {
		lc.appLog.Println(`LifeCycle entered observer mode`)
		<-lc.Shutdown
		goto exit
	}

runloop:
	for {
		select {
		case <-lc.Shutdown:
			break runloop
		case <-lc.tick:
			lc.ghost()
			if err = lc.discardDeletedBlocked(); err == nil {
				// skip unblock steps if there was an error to discard
				// deleted blocks
				lc.unblock()
			}
			lc.deadlockResolver()
			lc.handleDelete()
			if !lc.soma.conf.NoPoke {
				lc.poke()
			}
		}
	}
exit:
}

// ghost deletes configurations that that are still in in
// awaiting_rollout and have update_available set, ie. they have not
// yet been sent to the monitoring system
func (lc *LifeCycle) ghost() {
	lc.conn.Exec(stmt.LifecycleDeleteGhosts)
	lc.conn.Exec(stmt.LifecycleDeleteFailedRollouts)
	lc.conn.Exec(stmt.LifecycleDeleteDeprovisioned)
}

// search if there are check instance configurations in status blocked
// for checkinstances that are flagged as deleted. These do not need to
// be rolled out. Delete the dependencies and set the instance
// configurations to awaiting_deletion/none.
func (lc *LifeCycle) discardDeletedBlocked() error {
	var (
		err                          error
		blockedID, blockingID, state string
		tx                           *sql.Tx
		deps                         *sql.Rows
	)

	if deps, err = lc.stmtDeleteBlocked.Query(); err != nil {
		lc.errLog.Printf("LifeCycle: %s\n", err.Error())
		return err
	}
	defer deps.Close()

	// open multi-statement transaction. this ensures that we never
	// create a partial discard that afterwards does not hit our select
	// statement to find it
	if tx, err = lc.conn.Begin(); err != nil {
		lc.errLog.Println(err)
		return err
	}

	for deps.Next() {
		if err = deps.Scan(
			&blockedID,
			&blockingID,
			&state,
		); err != nil {
			lc.errLog.Println(err)
			tx.Rollback()
			return err
		}

		// delete record that blockedID waits on blockingID
		if _, err = tx.Exec(stmt.LifecycleDeleteDependency, blockedID, blockingID, state); err != nil {
			lc.errLog.Println(err)
			tx.Rollback()
			return err
		}

		// set blockedID to awaiting_deletion
		if _, err = tx.Exec(stmt.LifecycleConfigAwaitingDeletion, blockedID); err != nil {
			lc.errLog.Println(err)
			tx.Rollback()
			return err
		}
	}
	if deps.Err() != nil {
		lc.errLog.Println(err)
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		lc.errLog.Println(err)
		tx.Rollback()
		return err
	}
	return nil
}

// unblock checks if there are any blocked instance configurations
// that have their unblock condition met and moves them along in the
// rollout workflow
func (lc *LifeCycle) unblock() {
	var (
		cfgIds                            *sql.Rows
		blockedID, blockingID, instanceID string
		state, next, nextNG, status       string
		err                               error
		tx                                *sql.Tx
	)

	// lcStmtActiveUnblockCondition
	if cfgIds, err = lc.stmtUnblock.Query(); err != nil {
		lc.errLog.Println(err)
		return
	}
	defer cfgIds.Close()

idloop:
	for cfgIds.Next() {
		txMap := map[string]*sql.Stmt{}
		if err = cfgIds.Scan(
			&blockedID,
			&blockingID,
			&state,
			&status,
			&next,
			&instanceID,
		); err != nil {
			lc.errLog.Println(err.Error())
			continue idloop
		}

		if tx, err = lc.conn.Begin(); err != nil {
			lc.errLog.Println(err.Error())
			continue idloop
		}

		for name, statement := range map[string]string{
			`update`:   stmt.LifecycleUpdateConfig,
			`delete`:   stmt.LifecycleDeleteDependency,
			`instance`: stmt.LifecycleUpdateInstance,
		} {
			if txMap[name], err = tx.Prepare(statement); err != nil {
				lc.errLog.Println(`aborting lifecycle transaction`, err, stmt.Name(statement))
				// tx.Rollback() closes open prepared statements
				tx.Rollback()
				continue idloop
			}
		}

		switch next {
		case "awaiting_rollout":
			nextNG = "rollout_in_progress"
		default:
			lc.errLog.Printf("LifeCycle.unblock() error: blocked: %s, blocking %s, next: %s, instanceID: %s\n",
				blockedID, blockingID, next, instanceID)
			tx.Rollback()
			continue idloop
		}
		if _, err = txMap[`update`].Exec(
			next,
			nextNG,
			false,
			blockedID,
		); err != nil {
			lc.errLog.Println(`LifeCycle.unblock(moveConfig)`, err.Error())
			tx.Rollback()
			continue idloop
		}
		if _, err = txMap[`instance`].Exec(
			true,
			blockedID,
			instanceID,
		); err != nil {
			lc.errLog.Println(`LifeCycle.unblock(updateInstance)`, err.Error())
			tx.Rollback()
			continue idloop
		}
		if _, err = txMap[`delete`].Exec(
			blockedID,
			blockingID,
			state,
		); err != nil {
			lc.errLog.Println(`LifeCycle.unblock(deleteDependency)`, err.Error())
			tx.Rollback()
			continue idloop
		}
		if err = tx.Commit(); err != nil {
			lc.errLog.Println(err.Error())
			tx.Rollback()
			continue idloop
		}
	}
}

// deadlockResolver checks if there are check instance configurations
// that block on another check instance configuration that is in a
// stable workflow state, which would cause the new configurations to
// block forever. Resolve this by forcing the blocking configuration
// into the deprovisioning workflow so that LifeCycle picks it up
// on the next tick.
func (lc *LifeCycle) deadlockResolver() {
	var (
		rows                       *sql.Rows
		chkInstID, chkInstConfigID string
		err                        error
	)

	if rows, err = lc.stmtDeadlock.Query(); err != nil {
		lc.errLog.Println(`LifeCycle.deadLockResolver()`, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&chkInstID,
			&chkInstConfigID,
		); err != nil {
			lc.errLog.Println(`LifeCycle.deadLockResolver()`, err)
			return
		}
		lc.conn.Exec(stmt.LifecycleUpdateConfig,
			`awaiting_deprovision`,
			`deprovision_in_progress`,
			false,
			chkInstConfigID,
		)
		lc.conn.Exec(stmt.LifecycleUpdateInstance,
			true,
			chkInstConfigID,
			chkInstID,
		)
	}
}

// handleDelete checks for check instance configurations that are
// currently provisioned whose check instance has been deleted and
// triggers their deprovisioning
func (lc *LifeCycle) handleDelete() {
	var (
		rows              *sql.Rows
		err               error
		instCfgID, instID string
		tx                *sql.Tx
	)

	if rows, err = lc.stmtDeleteActive.Query(); err != nil {
		lc.errLog.Println(err)
		return
	}
	defer rows.Close()

	if tx, err = lc.conn.Begin(); err != nil {
		lc.errLog.Println(err)
		return
	}

cfgloop:
	for rows.Next() {
		if err = rows.Scan(
			&instCfgID,
			&instID,
		); err != nil {
			lc.errLog.Println(err)
			continue cfgloop
		}

		// set instance configuration to awaiting_deprovision
		if _, err = tx.Exec(
			stmt.LifecycleDeprovisionConfiguration, instCfgID,
		); err != nil {
			lc.errLog.Println(err)
			tx.Rollback()
			return
		}

		// set instance to update_available -> pickup by poke
		if _, err = tx.Exec(
			stmt.LifecycleUpdateInstance, true, instCfgID, instID,
		); err != nil {
			lc.errLog.Println(err)
			tx.Rollback()
			return
		}
	}
	if rows.Err() != nil {
		lc.errLog.Println(err)
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		lc.errLog.Println(err)
		tx.Rollback()
	}
	return
}

// poke triggers update notifications to monitoring systems that have
// a configured callback address
func (lc *LifeCycle) poke() {
	var (
		chkIds                        *sql.Rows
		err                           error
		chkID, monitoringID, callback string
	)

	for _, mode := range []string{`reschedule`, `poke`} {
		switch mode {
		case `reschedule`:
			// reschedule picks up configurations that have been
			// previously notified, have update_available unset
			// and have not moved along in > 5 minutes
			if chkIds, err = lc.stmtReschedule.Query(); err != nil {
				lc.errLog.Println(`LifeCycle.reschedule()`, err)
			}
		case `poke`:
			// poke picks up configurations that have update_available
			// set and have not been notified before
			if chkIds, err = lc.stmtPoke.Query(); err != nil {
				lc.errLog.Println(`LifeCycle.poke()`, err)
			}
		}

		for chkIds.Next() {
			if err = chkIds.Scan(
				&chkID,
				&monitoringID,
				&callback,
			); err != nil {
				lc.errLog.Println(err)
				continue
			}

			// there is no goroutine running for the system yet
			if _, ok := lc.pokers[monitoringID]; !ok {
				lc.pokers[monitoringID] = make(chan string, 4096)
				go lc.pokeSystem(callback, lc.pokers[monitoringID])
			}

			// if the channel is full we skip the checkid and pick it
			// up on the next lifecycle tick, otherwise an unresponsive
			// monitoring system could block the lifecycle system
			if len(lc.pokers[monitoringID]) < 4095 {
				lc.pokers[monitoringID] <- chkID
				if mode == `poke` {
					// notify has been triggered,
					// clear update available flag
					lc.stmtClear.Exec(chkID)
				}
			}
		}
		if err = chkIds.Err(); err != nil {
			lc.errLog.Println(err)
		}
	}
}

// pokeSystem sends out the update notifications for check IDs it
// receives via channel in to the callback URL
func (lc *LifeCycle) pokeSystem(callback string, in chan string) {
	client := resty.New()
	for {
		select {
		case chkID := <-in:
			retries := 0
		retry:
			client.SetTimeout(time.Duration(
				lc.soma.conf.PokeTimeout) * time.Millisecond,
			)
			if _, err := client.R().SetBody(
				PokeMessage{
					UUID: chkID,
					Path: lc.soma.conf.PokePath,
				},
			).Post(callback); err != nil {
				// with limit 4 this implements retries with 1, 2, 4
				// and 8 seconds sleeps between them
				if retries < 4 {
					timeout := math.Pow(2, float64(retries))
					time.Sleep(time.Duration(timeout) * time.Second)
					retries++
					goto retry
				}
				lc.errLog.Println(err)
				continue
			}
			lc.appLog.Printf("Poked %s (%s)", callback, chkID)
			lc.stmtSetNotify.Exec(chkID)
		}
	}
}

// shutdownNow signals the handler to shut down
func (lc *LifeCycle) shutdownNow() {
	close(lc.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
