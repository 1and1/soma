package main

import (
	"database/sql"
	"time"

	"github.com/1and1/soma/internal/stmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/resty.v0"
)

type lifeCycle struct {
	shutdown     chan bool
	conn         *sql.DB
	tick         <-chan time.Time
	stmt_unblock *sql.Stmt
	stmt_poke    *sql.Stmt
	stmt_clear   *sql.Stmt
	stmt_delblk  *sql.Stmt
	stmt_delact  *sql.Stmt
	stmt_dead    *sql.Stmt
	appLog       *log.Logger
	reqLog       *log.Logger
	errLog       *log.Logger
}

type PokeMessage struct {
	Uuid string `json:"uuid"`
	// path should be used to tell the client system the basepath
	// where to get it so SOMA + path + item_id === complete_url
	Path string `json:"path"`
}

func (lc *lifeCycle) run() {
	var err error
	lc.tick = time.NewTicker(time.Duration(SomaCfg.LifeCycleTick) * time.Second).C

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.LifecycleActiveUnblockCondition:           lc.stmt_unblock,
		stmt.LifecycleReadyDeployments:                 lc.stmt_poke,
		stmt.LifecycleClearUpdateFlag:                  lc.stmt_clear,
		stmt.LifecycleBlockedConfigsForDeletedInstance: lc.stmt_delblk,
		stmt.LifecycleDeprovisionDeletedActive:         lc.stmt_delact,
		stmt.LifecycleDeadLockResolver:                 lc.stmt_dead,
	} {
		if prepStmt, err = lc.conn.Prepare(statement); err != nil {
			lc.errLog.Fatal(`lifecycle`, err, statement)
		}
		defer prepStmt.Close()
	}

	if SomaCfg.Observer {
		lc.appLog.Println(`LifeCycle entered observer mode`)
		<-lc.shutdown
		goto exit
	}

runloop:
	for {
		select {
		case <-lc.shutdown:
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
			if !SomaCfg.NoPoke {
				lc.poke()
			}
		}
	}
exit:
}

// ghost deletes configurations that that are still in in
// awaiting_rollout and have update_available set, ie. they have not yet
// been sent to the monitoring system
func (lc *lifeCycle) ghost() {
	lc.conn.Exec(stmt.LifecycleDeleteGhosts)
	lc.conn.Exec(stmt.LifecycleDeleteFailedRollouts)
	lc.conn.Exec(stmt.LifecycleDeleteDeprovisioned)
}

// search if there are check instance configurations in status blocked
// for checkinstances that are flagged as deleted. These do not need to
// be rolled out. Delete the dependencies and set the instance
// configurations to awaiting_deletion/none.
func (lc *lifeCycle) discardDeletedBlocked() error {
	var (
		err                          error
		blockedID, blockingID, state string
		tx                           *sql.Tx
		deps                         *sql.Rows
	)

	if deps, err = lc.stmt_delblk.Query(); err != nil {
		lc.errLog.Printf("LifeCycle: %s\n", err.Error())
		return err
	}
	defer deps.Close()

	// open multi-statement transaction. this ensures that we never
	// create a partial discard that awards does not hit our select
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

func (lc *lifeCycle) handleDelete() {
	var (
		rows              *sql.Rows
		err               error
		instCfgId, instId string
		tx                *sql.Tx
	)

	if rows, err = lc.stmt_delact.Query(); err != nil {
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
			&instCfgId,
			&instId,
		); err != nil {
			lc.errLog.Println(err)
			continue cfgloop
		}

		// set instance configuration to awaiting_deprovision
		if _, err = tx.Exec(stmt.LifecycleDeprovisionConfiguration, instCfgId); err != nil {
			lc.errLog.Println(err)
			tx.Rollback()
			return
		}

		// set instance to update_available -> pickup by poke
		if _, err = tx.Exec(stmt.LifecycleUpdateInstance, true, instCfgId, instId); err != nil {
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

func (lc *lifeCycle) unblock() {
	var (
		cfgIds                                                         *sql.Rows
		blockedID, blockingID, instanceID, state, next, nextNG, status string
		err                                                            error
		tx                                                             *sql.Tx
	)

	// lcStmtActiveUnblockCondition
	if cfgIds, err = lc.stmt_unblock.Query(); err != nil {
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
				// tx.Rollback() closes open prepared statements
				tx.Rollback()
				continue idloop
			}
		}

		switch next {
		case "awaiting_rollout":
			nextNG = "rollout_in_progress"
		default:
			lc.errLog.Printf("lifeCycle.unblock() error: blocked: %s, blocking%s, next: %s, instanceID: %s\n",
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
			lc.errLog.Println(`lifeCycle.unblock(moveConfig)`, err.Error())
			tx.Rollback()
			continue idloop
		}
		if _, err = txMap[`instance`].Exec(
			true,
			blockedID,
			instanceID,
		); err != nil {
			lc.errLog.Println(`lifeCycle.unblock(updateInstance)`, err.Error())
			tx.Rollback()
			continue idloop
		}
		if _, err = txMap[`delete`].Exec(
			blockedID,
			blockingID,
			state,
		); err != nil {
			lc.errLog.Println(`lifeCycle.unblock(deleteDependency)`, err.Error())
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

func (lc *lifeCycle) poke() {
	var (
		chkIds                        *sql.Rows
		err                           error
		chkID, monitoringID, callback string
		cl                            *resty.Client
	)

	if chkIds, err = lc.stmt_poke.Query(); err != nil {
		lc.errLog.Println(`lifeCycle.poke()`, err)
	}
	defer chkIds.Close()

	callbacks := map[string]string{}
	pokeIDs := map[string][]string{}

	for chkIds.Next() {
		if err = chkIds.Scan(
			&chkID,
			&monitoringID,
			&callback,
		); err != nil {
			lc.errLog.Println(err)
			continue
		}

		callbacks[monitoringID] = callback
		if pokeIDs[monitoringID] == nil {
			pokeIDs[monitoringID] = []string{}
		}
		pokeIDs[monitoringID] = append(pokeIDs[monitoringID], chkID)
	}

	// do not poke the bear
bearloop:
	for mon, idList := range pokeIDs {
		var i uint64 = 0
		for _, id := range idList {
			cl = resty.New()
			if _, err = cl.SetTimeout(500 * time.Millisecond).R().
				SetBody(PokeMessage{Uuid: id, Path: SomaCfg.PokePath}).
				Post(callbacks[mon]); err != nil {
				lc.errLog.Println(err)
				continue bearloop
			}
			// XXX TODO: MAYBE we should look at the return code. MAYBE.
			lc.appLog.Printf("Poked %s about %s", mon, id)
			lc.stmt_clear.Exec(id)
			i++
			if i == SomaCfg.PokeBatchSize {
				continue bearloop
			}
		}
	}
}

func (lc *lifeCycle) deadlockResolver() {
	var (
		rows                       *sql.Rows
		chkInstID, chkInstConfigID string
		err                        error
	)

	if rows, err = lc.stmt_dead.Query(); err != nil {
		lc.errLog.Println(`lifeCycle.deadLockResolver()`, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&chkInstID,
			&chkInstConfigID,
		); err != nil {
			lc.errLog.Println(`lifeCycle.deadLockResolver()`, err)
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

/* Ops Access
 */
func (lc *lifeCycle) shutdownNow() {
	lc.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
