package main

import "github.com/codegangsta/cli"

func registerCommands(app cli.App) *cli.App {

	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Initialize local client files",
			Action: cmdClientInit,
		},
		{
			Name:   "experiment",
			Usage:  "Test cli.Action functionality",
			Action: runtime(cmdExperiment),
		},
	}

	app = *registerAttributes(app)
	app = *registerBuckets(app)
	app = *registerCapability(app)
	app = *registerChecks(app)
	app = *registerClusters(app)
	app = *registerDatacenters(app)
	app = *registerEnvironments(app)
	app = *registerGroups(app)
	app = *registerJobs(app)
	app = *registerLevels(app)
	app = *registerMetrics(app)
	app = *registerModes(app)
	app = *registerMonitoring(app)
	app = *registerNodes(app)
	app = *registerOncall(app)
	app = *registerPermissions(app)
	app = *registerPredicates(app)
	app = *registerProperty(app)
	app = *registerProviders(app)
	app = *registerRights(app)
	app = *registerRepository(app)
	app = *registerServers(app)
	app = *registerStates(app)
	app = *registerStatus(app)
	app = *registerTeams(app)
	app = *registerTypes(app)
	app = *registerUnits(app)
	app = *registerUsers(app)
	app = *registerValidity(app)
	app = *registerViews(app)
	app = *registerOps(app)

	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
