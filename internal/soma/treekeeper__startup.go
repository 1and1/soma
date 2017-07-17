package soma

import (
	"database/sql"
	"encoding/json"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
)

func (tk *TreeKeeper) startupLoad() {

	stMap := tk.prepareStartupStatements()
	for n, _ := range stMap {
		defer stMap[n].Close()
	}

	tk.startupBuckets(stMap)
	tk.startupGroups(stMap)
	tk.startupGroupMemberGroups(stMap)
	tk.startupGroupedClusters(stMap)
	tk.startupClusters(stMap)
	tk.startupNodes(stMap)

	if len(tk.actions) > 0 {
		tk.startLog.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.meta.repoName)
		tk.status.isBroken = true
		return
	}

	// attach system properties
	tk.startupSystemProperties(stMap)

	if len(tk.actions) > 0 {
		tk.startLog.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.meta.repoName)
		tk.status.isBroken = true
		return
	}

	// attach service properties
	tk.startupServiceProperties(stMap)

	if len(tk.actions) > 0 {
		tk.startLog.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.meta.repoName)
		tk.status.isBroken = true
		return
	}

	// attach custom properties
	tk.startupCustomProperties(stMap)

	if len(tk.actions) > 0 {
		tk.startLog.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.meta.repoName)
		tk.status.isBroken = true
		return
	}

	// attach oncall properties
	tk.startupOncallProperties(stMap)

	if len(tk.actions) > 0 {
		tk.startLog.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.meta.repoName)
		tk.status.isBroken = true
		return
	}

	// attach checks
	tk.startupChecks(stMap)

	if !tk.status.requiresRebuild && len(tk.actions) > 0 {
		tk.startLog.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.meta.repoName)
		tk.status.isBroken = true
		return
	}

	// these run as part of a job, but not inside the job's transaction. If there are leftovers
	// after a crash, fix them up
	if !tk.soma.conf.Observer {
		tk.buildDeploymentDetails()
		tk.orderDeploymentDetails()
	}

	// preload pending/unfinished jobs if not rebuilding the tree or
	// running in observer mode
	if !tk.status.requiresRebuild && !tk.soma.conf.Observer {
		tk.startupJobs(stMap)
	}

	if !tk.status.requiresRebuild && len(tk.actions) > 0 {
		tk.startLog.Printf("TK[%s] ERROR! Stray startup actions pending in action queue!", tk.meta.repoName)
		tk.status.isBroken = true
		return
	}
}

func (tk *TreeKeeper) startupJobs(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		err  error
		rows *sql.Rows
		job  string
	)

	tk.startLog.Printf("TK[%s]: loading pending jobs\n", tk.meta.repoName)
	rows, err = stMap[`LoadJob`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading clusters: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

jobloop:
	for rows.Next() {
		err = rows.Scan(
			&job,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break jobloop
			}
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}

		// XXX BUG
		// REQUIRES MIGRATION TOOL TO CONVERT PENDING JOBS FROM treeRequest
		// TO msg.Request
		tr := msg.Request{}
		err = json.Unmarshal([]byte(job), &tr)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}
		tk.Input <- tr
		tk.startLog.Printf("TK[%s] Loaded job %s (%s)\n", tk.meta.repoName, tr.JobID, tr.Action)
	}
}

func (tk *TreeKeeper) prepareStartupStatements() map[string]*sql.Stmt {
	var err error
	stMap := map[string]*sql.Stmt{}

	for name, statement := range map[string]string{
		`LoadBucket`:             stmt.TkStartLoadBuckets,
		`LoadGroup`:              stmt.TkStartLoadGroups,
		`LoadGroupMbrGroup`:      stmt.TkStartLoadGroupMemberGroups,
		`LoadGroupMbrCluster`:    stmt.TkStartLoadGroupedClusters,
		`LoadCluster`:            stmt.TkStartLoadCluster,
		`LoadNode`:               stmt.TkStartLoadNode,
		`LoadJob`:                stmt.TkStartLoadJob,
		`LoadChecks`:             stmt.TkStartLoadChecks,
		`LoadItems`:              stmt.TkStartLoadInheritedChecks,
		`LoadConfig`:             stmt.TkStartLoadCheckConfiguration,
		`LoadAllConfigsForType`:  stmt.TkStartLoadAllCheckConfigurationsForType,
		`LoadThreshold`:          stmt.TkStartLoadCheckThresholds,
		`LoadCustomCstr`:         stmt.TkStartLoadCheckConstraintCustom,
		`LoadNativeCstr`:         stmt.TkStartLoadCheckConstraintNative,
		`LoadOncallCstr`:         stmt.TkStartLoadCheckConstraintOncall,
		`LoadAttributeCstr`:      stmt.TkStartLoadCheckConstraintAttribute,
		`LoadServiceCstr`:        stmt.TkStartLoadCheckConstraintService,
		`LoadSystemCstr`:         stmt.TkStartLoadCheckConstraintSystem,
		`LoadChecksForType`:      stmt.TkStartLoadChecksForType,
		`LoadInstances`:          stmt.TkStartLoadCheckInstances,
		`LoadInstanceCfg`:        stmt.TkStartLoadCheckInstanceConfiguration,
		`LoadGroupState`:         stmt.TkStartLoadCheckGroupState,
		`LoadGroupRelations`:     stmt.TkStartLoadCheckGroupRelations,
		`CapabilityView`:         stmt.TreekeeperGetViewFromCapability,
		`LoadPropRepoSystem`:     stmt.TkStartLoadRepoSysProp,
		`LoadPropBuckSystem`:     stmt.TkStartLoadBucketSysProp,
		`LoadPropGrpSystem`:      stmt.TkStartLoadGroupSysProp,
		`LoadPropClrSystem`:      stmt.TkStartLoadClusterSysProp,
		`LoadPropNodeSystem`:     stmt.TkStartLoadNodeSysProp,
		`LoadPropSystemInstance`: stmt.TkStartLoadSystemPropInstances,
		`LoadPropRepoCustom`:     stmt.TkStartLoadRepositoryCstProp,
		`LoadPropBuckCustom`:     stmt.TkStartLoadBucketCstProp,
		`LoadPropGrpCustom`:      stmt.TkStartLoadGroupCstProp,
		`LoadPropClrCustom`:      stmt.TkStartLoadClusterCstProp,
		`LoadPropNodeCustom`:     stmt.TkStartLoadNodeCstProp,
		`LoadPropCustomInstance`: stmt.TkStartLoadCustomPropInstances,
		`LoadPropRepoOncall`:     stmt.TkStartLoadRepoOncProp,
		`LoadPropBuckOncall`:     stmt.TkStartLoadBucketOncProp,
		`LoadPropGrpOncall`:      stmt.TkStartLoadGroupOncProp,
		`LoadPropClrOncall`:      stmt.TkStartLoadClusterOncProp,
		`LoadPropNodeOncall`:     stmt.TkStartLoadNodeOncProp,
		`LoadPropOncallInstance`: stmt.TkStartLoadOncallPropInstances,
		`LoadPropRepoService`:    stmt.TkStartLoadRepoSvcProp,
		`LoadPropBuckService`:    stmt.TkStartLoadBucketSvcProp,
		`LoadPropGrpService`:     stmt.TkStartLoadGroupSvcProp,
		`LoadPropClrService`:     stmt.TkStartLoadClusterSvcProp,
		`LoadPropNodeService`:    stmt.TkStartLoadNodeSvcProp,
		`LoadPropRepoSvcAttr`:    stmt.TkStartLoadRepoSvcAttr,
		`LoadPropBuckSvcAttr`:    stmt.TkStartLoadBucketSvcAttr,
		`LoadPropGrpSvcAttr`:     stmt.TkStartLoadGroupSvcAttr,
		`LoadPropClrSvcAttr`:     stmt.TkStartLoadClusterSvcAttr,
		`LoadPropNodeSvcAttr`:    stmt.TkStartLoadNodeSvcAttr,
		`LoadPropSvcInstance`:    stmt.TkStartLoadServicePropInstances,
	} {
		if stMap[name], err = tk.conn.Prepare(statement); err != nil {
			tk.startLog.Println(`treekeeper startup`, err,
				stmt.Name(statement))
			tk.status.isBroken = true
			return map[string]*sql.Stmt{}
		}
	}
	return stMap
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
