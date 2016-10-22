package main

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/1and1/soma/lib/proto"
)

func (tk *treeKeeper) orderDeploymentDetails() {

	var (
		computed *sql.Rows
		err      error
	)
	if computed, err = tk.stmt_GetComputed.Query(tk.repoId); err != nil {
		tk.errLog.Println("tk.stmt_GetComputed.Query(): ", err)
		tk.broken = true
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
			txUpdateStatus                *sql.Stmt
			txUpdateInstance              *sql.Stmt
			txUpdateExisting              *sql.Stmt
			txDependency                  *sql.Stmt
		)
		err = computed.Scan(
			&chkInstanceId,
			&currentChkInstanceConfigId,
			&currentDeploymentDetailsJSON,
		)
		if err == sql.ErrNoRows {
			continue deployments
		} else if err != nil {
			tk.errLog.Println("tk.stmt_GetComputed.Query().Scan(): ", err)
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
			tk.errLog.Println("tk.stmt_GetPrevious.QueryRow(): ", err)
			break deployments
		}

		/* there is no previous version of this check instance rolled
		 * out on a monitoring system
		 */
		if noPrevious {
			// open multi statement transaction
			if tx, err = tk.conn.Begin(); err != nil {
				tk.errLog.Println("TreeKeeper/Order sql.Begin: ", err)
				break deployments
			}
			defer tx.Rollback()

			// prepare statements within transaction
			if txUpdateStatus, err = tx.Prepare(tkStmtUpdateConfigStatus); err != nil {
				tk.errLog.Printf("Failed to prepare %s: %s\n",
					`tkStmtUpdateConfigStatus`, err)
				break deployments
			}
			defer txUpdateStatus.Close()

			if txUpdateInstance, err = tx.Prepare(tkStmtUpdateCheckInstance); err != nil {
				tk.errLog.Printf("Failed to prepare %s: %s\n",
					`tkStmtUpdateCheckInstance`, err)
				break deployments
			}
			defer txUpdateInstance.Close()

			//
			if _, err = txUpdateStatus.Exec(
				"awaiting_rollout",
				"rollout_in_progress",
				currentChkInstanceConfigId,
			); err != nil {
				goto bailout_noprev
			}

			if _, err = txUpdateInstance.Exec(
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
			tk.errLog.Printf("Error unmarshal/deploymentdetails %s: %s",
				currentChkInstanceConfigId,
				err.Error(),
			)
			err = nil
			continue deployments
		}
		err = json.Unmarshal([]byte(previousDeploymentDetailsJSON), &prvDetails)
		if err != nil {
			tk.errLog.Printf("Error unmarshal/deploymentdetails %s: %s",
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
		if tx, err = tk.conn.Begin(); err != nil {
			tk.errLog.Println("TreeKeeper/Order sql.Begin: ", err)
			break deployments
		}
		defer tx.Rollback()

		// prepare statements within transaction
		if txUpdateStatus, err = tx.Prepare(tkStmtUpdateConfigStatus); err != nil {
			tk.errLog.Println("Failed to prepare %s: %s\n",
				`tkStmtUpdateConfigStatus`, err)
			break deployments
		}
		defer txUpdateStatus.Close()

		if txUpdateInstance, err = tx.Prepare(tkStmtUpdateCheckInstance); err != nil {
			tk.errLog.Println("Failed to prepare %s: %s\n",
				`tkStmtUpdateCheckInstance`, err)
			break deployments
		}
		defer txUpdateInstance.Close()

		if txUpdateExisting, err = tx.Prepare(tkStmtUpdateExistingCheckInstance); err != nil {
			tk.errLog.Println("Failed to prepare %s: %s\n",
				`tkStmtUpdateExistingCheckInstance`, err)
			break deployments
		}
		defer txUpdateExisting.Close()

		if txDependency, err = tx.Prepare(tkStmtSetDependency); err != nil {
			tk.errLog.Println("Failed to prepare %s: %s\n",
				`tkStmtSetDependency`, err)
			break deployments
		}
		defer txDependency.Close()

		if _, err = txUpdateStatus.Exec(
			"blocked",
			"awaiting_rollout",
			currentChkInstanceConfigId,
		); err != nil {
			goto bailout_withprev
		}
		if _, err = txUpdateExisting.Exec(
			time.Now().UTC(),
			true,
			chkInstanceId,
		); err != nil {
			goto bailout_withprev
		}
		if _, err = txDependency.Exec(
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
		tk.broken = true
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
