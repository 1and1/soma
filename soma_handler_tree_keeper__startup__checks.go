package main

import "database/sql"

type tkLoaderChecks struct {
	loadChecks     *sql.Stmt
	loadItems      *sql.Stmt
	loadConfig     *sql.Stmt
	loadThresh     *sql.Stmt
	loadCstrCustom *sql.Stmt
	loadCstrNative *sql.Stmt
	loadCstrOncall *sql.Stmt
	loadCstrAttr   *sql.Stmt
	loadCstrSystem *sql.Stmt
	loadInstances  *sql.Stmt
	loadInstConfig *sql.Stmt
}

func (tk *treeKeeper) startupChecks() {

	ld := tkLoaderChecks{}
	// TODO: prepare statements

	for _, typ := range []string{`repository`, `bucket`, `group`, `cluster`, `node`} {
		tk.startupScopedChecks(typ, &ld)
	}
}

func (tk *treeKeeper) startupScopedChecks(typ string, ld *tkLoaderChecks) {
	if tk.broken {
		return
	}

	/*
		Source-Checks laden: tkStmtLoadChecks
		Source-Checks.each:
			Config laden: tkStmtLoadCheckConfiguration
			Thresholds laden: tkStmtLoadCheckThresholds
			Constraints laden:	tkStmtLoadCheckConstraintCustom
								tkStmtLoadCheckConstraintNative
								tkStmtLoadCheckConstraintOncall
								tkStmtLoadCheckConstraintAttribute
								tkStmtLoadCheckConstraintService
								tkStmtLoadCheckConstraintSystem
			Vererbte Checks laden [CheckItem]: tkStmtLoadInheritedChecks
			--> Check anlegen
		!! << NICHT tk.tree.ComputeCheckInstances() AUFRUFEN >> !!
		CheckInstanzen laden: tkStmtLoadCheckInstances
		CheckInstanz.each:
			InstanzConfig laden: tkStmtLoadCheckInstanceConfiguration
			--> Instanz anlegen
			--> TODO: somatree CheckInstanz-Load Interface
	*/
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
