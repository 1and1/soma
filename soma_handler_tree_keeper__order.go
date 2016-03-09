package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

)

func (tk *treeKeeper) orderDeploymentDetails() {

	var (
		stmt_GetComputed  *sql.Stmt
		stmt_GetPrevious  *sql.Stmt
		stmt_DelDuplicate *sql.Stmt
		computed          *sql.Rows
		err               error
	)
	if stmt_GetComputed, err = tk.conn.Prepare(tkStmtGetComputedDeployments); err != nil {
		log.Fatal(err)
	}
	defer stmt_GetComputed.Close()

	if stmt_GetPrevious, err = tk.conn.Prepare(tkStmtGetPreviousDeployment); err != nil {
		log.Fatal(err)
	}
	defer stmt_GetPrevious.Close()

	if stmt_DelDuplicate, err = tk.conn.Prepare(tkStmtDeleteDuplicateDetails); err != nil {
		log.Fatal(err)
	}
	defer stmt_DelDuplicate.Close()

	if computed, err = stmt_GetComputed.Query(); err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}

		// fetch previous deployment details for this check_instance_id
		err = stmt_GetPrevious.QueryRow(chkInstanceId).Scan(
			&previousChkInstanceConfigId,
			&previousVersion,
			&previousStatus,
			&previousDeploymentDetailsJSON,
		)
		if err == sql.ErrNoRows {
			noPrevious = true
		} else if err != nil {
			log.Fatal(err)
		}

		/* there is no previous version of this check instance rolled
		 * out on a monitoring system
		 */
		if noPrevious {
			// open multi statement transaction
			if tx, err = tk.conn.Begin(); err != nil {
				log.Fatal(err)
			}
			defer tx.Rollback()

			// prepare statements within transaction
			if txUpdateStatus, err = tx.Prepare(tkStmtUpdateConfigStatus); err != nil {
				log.Fatal()
			}
			defer txUpdateStatus.Close()

			if txUpdateInstance, err = tx.Prepare(tkStmtUpdateCheckInstance); err != nil {
				log.Fatal()
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
		curDetails := somaproto.DeploymentDetails{}
		prvDetails := somaproto.DeploymentDetails{}
		json.Unmarshal([]byte(currentDeploymentDetailsJSON), curDetails)
		json.Unmarshal([]byte(previousDeploymentDetailsJSON), prvDetails)

		if curDetails.DeepCompare(&prvDetails) {
			// there is no change in deployment details, thus no point
			// to sending the new deployment details as an update to the
			// monitoring systems
			stmt_DelDuplicate.Exec(currentChkInstanceConfigId)
			continue deployments
		}

		// UPDATE config status
		// open multi statement transaction
		if tx, err = tk.conn.Begin(); err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback()

		// prepare statements within transaction
		if txUpdateStatus, err = tx.Prepare(tkStmtUpdateConfigStatus); err != nil {
			log.Fatal()
		}
		defer txUpdateStatus.Close()

		if txUpdateInstance, err = tx.Prepare(tkStmtUpdateCheckInstance); err != nil {
			log.Fatal()
		}
		defer txUpdateInstance.Close()

		if txDependency, err = tx.Prepare(tkStmtSetDependency); err != nil {
			log.Fatal()
		}
		defer txDependency.Close()

		switch previousStatus {
		case "awaiting_rollout":
			if _, err = txUpdateStatus.Exec(
				"deprovisioned",
				"none",
				previousChkInstanceConfigId,
			); err != nil {
				goto bailout_withprev
			}
			if _, err = txUpdateStatus.Exec(
				"awaiting_rollout",
				"rollout_in_progress",
				currentChkInstanceConfigId,
			); err != nil {
				goto bailout_withprev
			}
		case "rollout_in_progress":
			if _, err = txUpdateStatus.Exec(
				"rollout_in_progress",
				"awaiting_deprovision",
				previousChkInstanceConfigId,
			); err != nil {
				goto bailout_withprev
			}
			if _, err = txUpdateStatus.Exec(
				"blocked",
				"awaiting_rollout",
				currentChkInstanceConfigId,
			); err != nil {
				goto bailout_withprev
			}
			if _, err = txUpdateInstance.Exec(
				time.Now().UTC(),
				false,
				previousChkInstanceConfigId,
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
		case "active":
			if _, err = txUpdateStatus.Exec(
				"awaiting_deprovision",
				"deprovision_in_progress",
				previousChkInstanceConfigId,
			); err != nil {
				goto bailout_withprev
			}
			if _, err = txUpdateStatus.Exec(
				"blocked",
				"awaiting_rollout",
				currentChkInstanceConfigId,
			); err != nil {
				goto bailout_withprev
			}
			if _, err = txUpdateInstance.Exec(
				time.Now().UTC(),
				true,
				previousChkInstanceConfigId,
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
		case "blocked":
			// TODO: if there are >1 configurations waiting for rollout,
			// these could be compacted by cutting middle versions
			fallthrough
		case "awaiting_deprovision":
			fallthrough
		case "deprovision_in_progress":
			if _, err = txUpdateStatus.Exec(
				"blocked",
				"awaiting_rollout",
				currentChkInstanceConfigId,
			); err != nil {
				goto bailout_withprev
			}
			if _, err = txUpdateInstance.Exec(
				time.Now().UTC(),
				false,
				previousChkInstanceConfigId,
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
		case "deprovisioned":
			fallthrough
		case "awaiting_deletion":
			if _, err = txUpdateStatus.Exec(
				"awaiting_rollout",
				"rollout_in_progress",
				currentChkInstanceConfigId,
			); err != nil {
				goto bailout_withprev
			}
			if _, err = txUpdateInstance.Exec(
				time.Now().UTC(),
				true,
				currentChkInstanceConfigId,
				chkInstanceId,
			); err != nil {
				goto bailout_withprev
			}
		}

		if err = tx.Commit(); err != nil {
			goto bailout_withprev
		}
		continue deployments

	bailout_withprev:
		tx.Rollback()
		continue deployments
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
