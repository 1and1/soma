package main

func (tk *treeKeeper) orderDeploymentDetails() {

	/*
		1. GET ALL NEW INSTANCES
		select	check_instance_id,
				check_instance_config_id,
				deployment_details
		from	soma.check_instance_configurations
		where	status = 'computed'

		for loop:
			2. GET PREVIOUS VERSION
			select	check_instance_config_id,
					version,
					status,
					deployment_details
			from	soma.check_instance_configurations
			where	status != 'computed'
			and		check_instance_id = $1::uuid
			order	by version DESC
			limit	1

			if curr.version(deployment_details).DeepCompare(prev.version(deployment_details)):
				no change -> delete cur.version
				continue for loop

			if status == "active":
				prev.version status      => active
							 next_status => awaiting_deprovision
				curr.version status      => blocked
							 next_status => awaiting_rollout

				soma.check_instance_configuration_dependencies:
				blocked_instance_config_id:  curr.version
				blocking_instance_config_id: prev.version
				unblocking_state:            deprovisioned

			if status == "awaiting_deprovision || deprovision_in_progress":
				curr.version status      => blocked
							 next_status => awaiting_rollout

				soma.check_instance_configuration_dependencies:
				blocked_instance_config_id:  curr.version
				blocking_instance_config_id: prev.version
				unblocking_state:            deprovisioned

			if awaiting_rollout:
				prev.version status      => deprovisioned
				             next_status => none
				curr.version status      => awaiting_rollout
				             next_status => rollout_in_progress

			if rollout_in_progress:
				prev.version status      => rollout_in_progress
				             next_status => awaiting_deprovision
				curr.version status      => blocked
				             next_status => awaiting_rollout

				soma.check_instance_configuration_dependencies:
				blocked_instance_config_id:  curr.version
				blocking_instance_config_id: prev.version
				unblocking_state:            deprovisioned

	*/

}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
