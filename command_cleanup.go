package main

import "fmt"

func commandCleanupChecks(done chan bool, printOnly bool) {
	if !printOnly {
		dbOpen()
	}

	stmts := []string{
		`DELETE FROM check_instance_configurations;`,
		`DELETE FROM check_instance_configuration_dependencies;`,
		`DELETE FROM check_instances;`,
		`DELETE FROM checks;`,
		`DELETE FROM configuration_thresholds;`,
		`DELETE FROM constraints_custom_property;`,
		`DELETE FROM constraints_native_property;`,
		`DELETE FROM constraints_oncall_property;`,
		`DELETE FROM constraints_service_attribute;`,
		`DELETE FROM constraints_service_property;`,
		`DELETE FROM constraints_system_property;`,
		`DELETE FROM check_configurations;`,
	}

	for _, stmt := range stmts {
		if printOnly {
			fmt.Println(stmt)
			continue
		}
		db.Exec(stmt)
	}

	done <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
