package soma

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
)

func (tk *TreeKeeper) orderDeploymentDetails() {

	var (
		computed *sql.Rows
		err      error
	)
	if computed, err = tk.stmt_GetComputed.Query(tk.meta.repoID); err != nil {
		tk.treeLog.Println("tk.stmt_GetComputed.Query(): ", err)
		tk.status.isBroken = true
		return
	}
	defer computed.Close()

deployments:
	for computed.Next() {
		var (
			chkInstanceId                 string
			currentChkInstanceConfigId    string
			currentDeploymentDetailsJSON  string
			previousChkInstanceConfigId   string
			previousVersion               string
			previousStatus                string
			previousDeploymentDetailsJSON string
			noPrevious                    bool
			tx                            *sql.Tx
		)
		err = computed.Scan(
			&chkInstanceId,
			&currentChkInstanceConfigId,
			&currentDeploymentDetailsJSON,
		)
		if err == sql.ErrNoRows {
			continue deployments
		} else if err != nil {
			tk.treeLog.Println("tk.stmt_GetComputed.Query().Scan(): ", err)
			break deployments
		}

		// fetch previous deployment details for this check_instance_id
		err = tk.stmt_GetPrevious.QueryRow(chkInstanceId).Scan(
			&previousChkInstanceConfigId,
			&previousVersion,
			&previousStatus,
			&previousDeploymentDetailsJSON,
		)
		if err == sql.ErrNoRows {
			noPrevious = true
		} else if err != nil {
			tk.treeLog.Println("tk.stmt_GetPrevious.QueryRow(): ", err)
			break deployments
		}

		/* there is no previous version of this check instance rolled
		 * out on a monitoring system
		 */
		if noPrevious {
			// open multi statement transaction
			txMap := map[string]*sql.Stmt{}
			if tx, err = tk.conn.Begin(); err != nil {
				tk.treeLog.Println("TreeKeeper/Order sql.Begin: ", err)
				break deployments
			}

			// prepare statements within transaction
			for name, statement := range map[string]string{
				`UpdateStatus`:   stmt.TreekeeperUpdateConfigStatus,
				`UpdateInstance`: stmt.TreekeeperUpdateCheckInstance,
			} {
				if txMap[name], err = tx.Prepare(statement); err != nil {
					tk.treeLog.Println(`treekeeper/order/tx`, err, stmt.Name(statement))
					tx.Rollback()
					break deployments
				}
			}

			//
			if _, err = txMap[`UpdateStatus`].Exec(
				"awaiting_rollout",
				"rollout_in_progress",
				currentChkInstanceConfigId,
			); err != nil {
				goto bailout_noprev
			}

			if _, err = txMap[`UpdateInstance`].Exec(
				time.Now().UTC(),
				true,
				currentChkInstanceConfigId,
				chkInstanceId,
			); err != nil {
				goto bailout_noprev
			}

			if err = tx.Commit(); err != nil {
				goto bailout_noprev
			}
			continue deployments

		bailout_noprev:
			tx.Rollback()
			continue deployments
		}
		/* a previous version of this check instance was found
		 */
		curDetails := proto.Deployment{}
		prvDetails := proto.Deployment{}
		err = json.Unmarshal([]byte(currentDeploymentDetailsJSON), &curDetails)
		if err != nil {
			tk.treeLog.Printf("Error unmarshal/deploymentdetails %s: %s",
				currentChkInstanceConfigId,
				err.Error(),
			)
			err = nil
			continue deployments
		}
		err = json.Unmarshal([]byte(previousDeploymentDetailsJSON), &prvDetails)
		if err != nil {
			tk.treeLog.Printf("Error unmarshal/deploymentdetails %s: %s",
				previousChkInstanceConfigId,
				err.Error(),
			)
			err = nil
			continue deployments
		}

		if curDetails.DeepCompare(&prvDetails) {
			// there is no change in deployment details, thus no point
			// to sending the new deployment details as an update to the
			// monitoring systems
			tk.stmt_DelDuplicate.Exec(currentChkInstanceConfigId)
			continue deployments
		}

		// UPDATE config status
		// open multi statement transaction
		txMap := map[string]*sql.Stmt{}
		if tx, err = tk.conn.Begin(); err != nil {
			tk.treeLog.Println("TreeKeeper/Order sql.Begin: ", err)
			break deployments
		}

		// prepare statements within transaction
		for name, statement := range map[string]string{
			`UpdateStatus`:   stmt.TreekeeperUpdateConfigStatus,
			`UpdateExisting`: stmt.TreekeeperUpdateExistingCheckInstance,
			`SetDependency`:  stmt.TreekeeperSetDependency,
		} {
			if txMap[name], err = tx.Prepare(statement); err != nil {
				tk.treeLog.Println(`treekeeper/order/tx`, err, stmt.Name(statement))
				tx.Rollback()
				break deployments
			}
		}

		if _, err = txMap[`UpdateStatus`].Exec(
			"blocked",
			"awaiting_rollout",
			currentChkInstanceConfigId,
		); err != nil {
			goto bailout_withprev
		}
		if _, err = txMap[`UpdateExisting`].Exec(
			time.Now().UTC(),
			true,
			chkInstanceId,
		); err != nil {
			goto bailout_withprev
		}
		if _, err = txMap[`SetDependency`].Exec(
			currentChkInstanceConfigId,
			previousChkInstanceConfigId,
			"deprovisioned",
		); err != nil {
			goto bailout_withprev
		}

		if err = tx.Commit(); err != nil {
			goto bailout_withprev
		}
		continue deployments

	bailout_withprev:
		tx.Rollback()
		continue deployments
	}
	// mark the tree as broken to prevent further data processing
	if err != nil {
		tk.status.isBroken = true
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
