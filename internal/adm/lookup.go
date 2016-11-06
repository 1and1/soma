/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package adm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/1and1/soma/lib/proto"
	resty "gopkg.in/resty.v0"
)

// LookupOncallId looks up the UUID for an oncall duty on the
// server with name s. Error is set if no such oncall duty was
// found or an error occured.
// If s is already a UUID, then s is immediately returned.
func LookupOncallId(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return oncallIdByName(s)
}

// LookupOncallId looks up the UUID for a user on the server
// with username s. Error is set if no such user was found
// or an error occured.
// If s is already a UUID, then s is immediately returned.
func LookupUserId(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return userIdByUserName(s)
}

// LookupTeamId looks up the UUID for a team on the server
// with teamname s. Error is set if no such team was found
// or an error occured.
// If s is already a UUID, then s is immediately returned.
func LookupTeamId(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return teamIdByName(s)
}

// LookupTeamByRepo looks up the UUID for the team that is the
// owner of a given repository s, which can be the name or UUID
// of the repository.
func LookupTeamByRepo(s string) (string, error) {
	var (
		bId string
		err error
	)

	if !IsUUID(s) {
		if bId, err = LookupRepoId(s); err != nil {
			return ``, err
		}
	} else {
		bId = s
	}

	return teamIdByRepoId(bId)
}

// LookupTeamByBucket looks up the UUID for the team that is
// the owner of a given bucket s, which can be the name or
// UUID of the bucket.
func LookupTeamByBucket(s string) (string, error) {
	var (
		bId string
		err error
	)

	if !IsUUID(s) {
		if bId, err = LookupBucketId(s); err != nil {
			return ``, err
		}
	} else {
		bId = s
	}

	return teamIdByBucketId(bId)
}

// LookupTeamByNode looks up the UUID for the team that is
// the owner of a given node s, which can be the name or
// UUID of the node.
func LookupTeamByNode(s string) (string, error) {
	var (
		nId string
		err error
	)

	if !IsUUID(s) {
		if nId, err = LookupBucketId(s); err != nil {
			return ``, err
		}
	} else {
		nId = s
	}

	return teamIdByNodeId(nId)
}

// LookupRepoId looks up the UUID for a repository on the server
// with reponame s. Error is set if no such repository was found
// or an error occured.
// If s is already a UUID, then s is immediately returned.
func LookupRepoId(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return repoIdByName(s)
}

// LookupRepoByBucket looks up the UUI for a repository by either
// the UUID or name of a bucket in that repository.
func LookupRepoByBucket(s string) (string, error) {
	var (
		bId string
		err error
	)

	if !IsUUID(s) {
		if bId, err = LookupBucketId(s); err != nil {
			return ``, err
		}
	} else {
		bId = s
	}

	return repoIdByBucketId(bId)
}

// LookupBucketId looks up the UUID for a bucket on the server
// with bucketname s. Error is set if no such bucket was found
// or an error occured.
// If s is already a UUID, then s is immediately returned.
func LookupBucketId(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return bucketIdByName(s)
}

func LookupGroupId(group, bucket string) (string, error) {
	if IsUUID(group) {
		return group, nil
	}
	var (
		bId string
		err error
	)
	if !IsUUID(bucket) {
		if bId, err = LookupBucketId(bucket); err != nil {
			return ``, err
		}
	} else {
		bId = bucket
	}

	return groupIdByName(group, bId)
}

func LookupClusterId(cluster, bucket string) (string, error) {
	if IsUUID(cluster) {
		return cluster, nil
	}
	var (
		bId string
		err error
	)
	if !IsUUID(bucket) {
		if bId, err = LookupBucketId(bucket); err != nil {
			return ``, err
		}
	} else {
		bId = bucket
	}

	return clusterIdByName(cluster, bId)
}

// LookupServerId looks up the UUID for a server either in the
// local cache or on the server. Error is set if no such server
// was found or an error occured.
// If s is already a UUID, then s is immediately returned.
// If s is a Uint64 number, then the serverlookup is by AssetID.
// Otherwise s is the server name.
func LookupServerId(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	if ok, num := isUint64(s); ok {
		return serverIdByAsset(s, num)
	}
	return serverIdByName(s)
}

// LookupPermIdRef looks up the UUID for a permission from
// the server. Error is set if no such permission was found or
// an error occured.
// If s is already a UUID, then is is immediately returned.
func LookupPermIdRef(s string, id *string) error {
	if IsUUID(s) {
		*id = s
		return nil
	}
	return permissionIdByName(s, id)
}

// LookupGrantIdRef looks up the UUID of a permission grant from
// the server and fills it into the provided id pointer.
// Error is set if no such grant was found or an error occured.
func LookupGrantIdRef(rcptType, rcptId, permId, cat string,
	id *string) error {
	return grantIdFromServer(rcptType, rcptId, permId, cat, id)
}

// LookupMonitoringId looks up the UUID of the monitoring system
// with the name s. Returns immediately if s is a UUID.
func LookupMonitoringId(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return monitoringIdByName(s)
}

// LookupNodeId looks up the UUID of the repository the bucket
// given via string s is part of. If s is a UUID, it is used as
// bucketId for the lookup.
func LookupNodeId(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return nodeIdByName(s)
}

// LookupCapabilityId looks up the UUID of the capability with the
// name s. Returns immediately if s is a UUID.
func LookupCapabilityId(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return capabilityIdByName(s)
}

// LookupNodeConfig looks up the node repo/bucket configuration
// given the name or UUID s of the node.
func LookupNodeConfig(s string) (*proto.NodeConfig, error) {
	var (
		nId string
		err error
	)

	if !IsUUID(s) {
		if nId, err = LookupNodeId(s); err != nil {
			return nil, err
		}
	} else {
		nId = s
	}

	return nodeConfigById(nId)
}

// LookupCheckConfigId looks up the UUID of check configuration s
// in Repository repo. Returns immediately if s is a UUID.
func LookupCheckConfigId(s, repo string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	var rId string
	if r, err := LookupRepoId(repo); err != nil {
		return ``, err
	} else {
		rId = r
	}
	return checkConfigIdByName(s, rId)
}

// LookupCustomPropertyId looks up the UUID of a custom property s
// in Repository repo. Returns immediately if s is a UUID.
func LookupCustomPropertyId(s, repo string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	var rId string
	if r, err := LookupRepoId(repo); err != nil {
		return ``, err
	} else {
		rId = r
	}
	return propertyIdByName(`custom`, s, rId)
}

// LookupServicePropertyId looks up the id of a service property s
// of team team.
func LookupServicePropertyId(s, team string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	var tId string
	if t, err := LookupTeamId(team); err != nil {
		return ``, err
	} else {
		tId = t
	}
	return propertyIdByName(`service`, s, tId)
}

// LookupTemplatePropertyId looks up the id of a service template
// property
func LookupTemplatePropertyId(s string) (string, error) {
	if IsUUID(s) {
		return s, nil
	}
	return propertyIdByName(`template`, s, `none`)
}

// LookupLevelName looks up the long name of a level s, where s
// can be the level's long or short name.
func LookupLevelName(s string) (string, error) {
	return levelByName(s)
}

// oncallIdByName implements the actual serverside lookup of the
// oncall duty UUID
func oncallIdByName(oncall string) (string, error) {
	req := proto.NewOncallFilter()
	req.Filter.Oncall = &proto.OncallFilter{Name: oncall}

	res, err := fetchFilter(req, `/filter/oncall/`)
	if err != nil {
		goto abort
	}

	if res.Oncalls == nil || len(*res.Oncalls) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if oncall != (*res.Oncalls)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			oncall, (*res.Oncalls)[0].Name)
		goto abort
	}
	return (*res.Oncalls)[0].Id, nil

abort:
	return ``, fmt.Errorf("OncallId lookup failed: %s", err.Error())
}

// userIdByUserName implements the actual serverside lookup of the
// user's UUID
func userIdByUserName(user string) (string, error) {
	req := proto.NewUserFilter()
	req.Filter.User.UserName = user

	res, err := fetchFilter(req, `/filter/users/`)
	if err != nil {
		goto abort
	}

	if res.Users == nil || len(*res.Users) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if user != (*res.Users)[0].UserName {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			user, (*res.Users)[0].UserName)
		goto abort
	}
	return (*res.Users)[0].Id, nil

abort:
	return ``, fmt.Errorf("UserId lookup failed: %s", err.Error())
}

// teamIdByName implements the actual serverside lookup of the
// team's UUID
func teamIdByName(team string) (string, error) {
	req := proto.NewTeamFilter()
	req.Filter.Team.Name = team

	res, err := fetchFilter(req, `/filter/teams/`)
	if err != nil {
		goto abort
	}

	if res.Teams == nil || len(*res.Teams) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if team != (*res.Teams)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			team, (*res.Teams)[0].Name)
		goto abort
	}
	return (*res.Teams)[0].Id, nil

abort:
	return ``, fmt.Errorf("TeamId lookup failed: %s", err.Error())
}

// teamIdByRepoId implements the actual serverside lookup of
// a repository's TeamId
func teamIdByRepoId(repo string) (string, error) {
	res, err := fetchObjList(fmt.Sprintf("/repository/%s", repo))
	if err != nil {
		goto abort
	}

	if res.Repositories == nil || len(*res.Repositories) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if repo != (*res.Repositories)[0].Id {
		err = fmt.Errorf("RepositoryId mismatch: %s vs %s",
			repo, (*res.Repositories)[0].Id)
		goto abort
	}
	return (*res.Repositories)[0].TeamId, nil

abort:
	return ``, fmt.Errorf("TeamId lookup failed: %s",
		err.Error())
}

// teamIdByBucketId implements the actual serverside lookup of
// a bucket's TeamId
func teamIdByBucketId(bucket string) (string, error) {
	res, err := fetchObjList(fmt.Sprintf("/buckets/%s", bucket))
	if err != nil {
		goto abort
	}

	if res.Buckets == nil || len(*res.Buckets) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if bucket != (*res.Buckets)[0].Id {
		err = fmt.Errorf("BucketId mismatch: %s vs %s",
			bucket, (*res.Buckets)[0].Id)
		goto abort
	}
	return (*res.Buckets)[0].TeamId, nil

abort:
	return ``, fmt.Errorf("TeamId lookup failed: %s",
		err.Error())
}

// teamIdByNodeId implements the actual serverside lookup of a
// node's TeamId
func teamIdByNodeId(node string) (string, error) {
	res, err := fetchObjList(fmt.Sprintf("/nodes/%s", node))
	if err != nil {
		goto abort
	}

	if res.Nodes == nil || len(*res.Nodes) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if node != (*res.Nodes)[0].Id {
		err = fmt.Errorf("NodeId mismatch: %s vs %s",
			node, (*res.Nodes)[0].Id)
		goto abort
	}
	return (*res.Nodes)[0].TeamId, nil

abort:
	return ``, fmt.Errorf("TeamId lookup failed: %s",
		err.Error())
}

// repoIdByName implements the actual serverside lookup of the
// repo's UUID
func repoIdByName(repo string) (string, error) {
	req := proto.NewRepositoryFilter()
	req.Filter.Repository.Name = repo

	res, err := fetchFilter(req, `/filter/repository/`)
	if err != nil {
		goto abort
	}

	if res.Repositories == nil || len(*res.Repositories) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if repo != (*res.Repositories)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			repo, (*res.Repositories)[0].Name)
		goto abort
	}
	return (*res.Repositories)[0].Id, nil

abort:
	return ``, fmt.Errorf("RepositoryId lookup failed: %s",
		err.Error())
}

// repoIdByBucketId implements the actual serverside lookup of the
// repo's UUID
func repoIdByBucketId(bucket string) (string, error) {
	res, err := fetchObjList(fmt.Sprintf("/buckets/%s", bucket))
	if err != nil {
		goto abort
	}

	if res.Buckets == nil || len(*res.Buckets) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if bucket != (*res.Buckets)[0].Id {
		err = fmt.Errorf("BucketId mismatch: %s vs %s",
			bucket, (*res.Buckets)[0].Id)
		goto abort
	}
	return (*res.Buckets)[0].RepositoryId, nil

abort:
	return ``, fmt.Errorf("RepositoryId lookup failed: %s",
		err.Error())
}

// bucketIdByName implements the actual serverside lookup of the
// bucket's UUID
func bucketIdByName(bucket string) (string, error) {
	req := proto.NewBucketFilter()
	req.Filter.Bucket.Name = bucket

	res, err := fetchFilter(req, `/filter/buckets/`)
	if err != nil {
		goto abort
	}

	if res.Buckets == nil || len(*res.Buckets) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if bucket != (*res.Buckets)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			bucket, (*res.Buckets)[0].Name)
		goto abort
	}
	return (*res.Buckets)[0].Id, nil

abort:
	return ``, fmt.Errorf("BucketId lookup failed: %s",
		err.Error())
}

//
func groupIdByName(group, bucketId string) (string, error) {
	req := proto.NewGroupFilter()
	req.Filter.Group.Name = group
	req.Filter.Group.BucketId = bucketId

	res, err := fetchFilter(req, `/filter/groups/`)
	if err != nil {
		goto abort
	}

	if res.Groups == nil || len(*res.Groups) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if group != (*res.Groups)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			group, (*res.Groups)[0].Name)
	}
	return (*res.Groups)[0].Id, nil

abort:
	return ``, fmt.Errorf("GroupId lookup failed: %s",
		err.Error())
}

//
func clusterIdByName(cluster, bucketId string) (string, error) {
	req := proto.NewClusterFilter()
	req.Filter.Cluster.Name = cluster
	req.Filter.Cluster.BucketId = bucketId

	res, err := fetchFilter(req, `/filter/clusters/`)
	if err != nil {
		goto abort
	}

	if res.Clusters == nil || len(*res.Clusters) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if cluster != (*res.Clusters)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			cluster, (*res.Clusters)[0].Name)
	}
	return (*res.Clusters)[0].Id, nil

abort:
	return ``, fmt.Errorf("ClusterId lookup failed: %s",
		err.Error())
}

// serverIdByName implements the actual lookup of the server UUID
// by name
func serverIdByName(s string) (string, error) {
	if m, err := cache.ServerByName(s); err == nil {
		return m[`id`], nil
	}
	req := proto.NewServerFilter()
	req.Filter.Server.Name = s

	res, err := fetchFilter(req, `/filter/servers/`)
	if err != nil {
		goto abort
	}

	if res.Servers == nil || len(*res.Servers) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if s != (*res.Servers)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			s, (*res.Servers)[0].Name)
		goto abort
	}
	// save server in cacheDB
	cache.Server(
		(*res.Servers)[0].Name,
		(*res.Servers)[0].Id,
		strconv.Itoa(int((*res.Servers)[0].AssetId)),
	)
	return (*res.Servers)[0].Id, nil

abort:
	return ``, fmt.Errorf("ServerId lookup failed: %s",
		err.Error())
}

// serverIdByAsset implements the actual lookup of the server UUID
// by numeric AssetID
func serverIdByAsset(s string, aid uint64) (string, error) {
	if m, err := cache.ServerByAsset(s); err == nil {
		return m[`id`], nil
	}
	req := proto.NewServerFilter()
	req.Filter.Server.AssetId = aid

	res, err := fetchFilter(req, `/filter/servers/`)
	if err != nil {
		goto abort
	}

	if res.Servers == nil || len(*res.Servers) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if aid != (*res.Servers)[0].AssetId {
		err = fmt.Errorf("AssetId mismatch: %d vs %d",
			aid, (*res.Servers)[0].AssetId)
		goto abort
	}
	// save server in cacheDB
	cache.Server(
		(*res.Servers)[0].Name,
		(*res.Servers)[0].Id,
		strconv.Itoa(int((*res.Servers)[0].AssetId)),
	)
	return (*res.Servers)[0].Id, nil

abort:
	return ``, fmt.Errorf("ServerId lookup failed: %s",
		err.Error())
}

// permissionIdByName implements the actual lookup of the permission
// UUID by name
func permissionIdByName(perm string, id *string) error {
	req := proto.NewPermissionFilter()
	req.Filter.Permission.Name = perm

	res, err := fetchFilter(req, `/filter/permission/`)
	if err != nil {
		goto abort
	}

	if res.Permissions == nil || len(*res.Permissions) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if perm != (*res.Permissions)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			perm, (*res.Permissions)[0].Name)
		goto abort
	}
	*id = (*res.Permissions)[0].Id
	return nil

abort:
	return fmt.Errorf("PermissionId lookup failed: %s",
		err.Error())
}

// grantIdFromServer implements the actual lookup of the grant UUID
func grantIdFromServer(rcptType, rcptId, permId, cat string,
	id *string) error {
	req := proto.NewGrantFilter()
	req.Filter.Grant.RecipientType = rcptType
	req.Filter.Grant.RecipientId = rcptId
	req.Filter.Grant.PermissionId = permId
	req.Filter.Grant.Category = cat

	res, err := fetchFilter(req, `/filter/grant/`)
	if err != nil {
		goto abort
	}

	if res.Grants == nil || len(*res.Grants) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if permId != (*res.Grants)[0].PermissionId {
		err = fmt.Errorf("PermissionId mismatch: %s vs %s",
			permId, (*res.Grants)[0].PermissionId)
		goto abort
	}
	*id = (*res.Grants)[0].Id

abort:
	return fmt.Errorf("GrantId lookup failed: %s",
		err.Error())
}

// monitoringIdByName implements the actual lookup of the monitoring
// system UUID
func monitoringIdByName(monitoring string) (string, error) {
	req := proto.NewMonitoringFilter()
	req.Filter.Monitoring.Name = monitoring

	res, err := fetchFilter(req, `/filter/monitoring/`)
	if err != nil {
		goto abort
	}

	if res.Monitorings == nil || len(*res.Monitorings) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if monitoring != (*res.Monitorings)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			monitoring, (*res.Monitorings)[0].Name)
		goto abort
	}
	return (*res.Monitorings)[0].Id, nil

abort:
	return ``, fmt.Errorf("MonitoringId lookup failed: %s",
		err.Error())
}

// nodeIdByName implements the actual lookup of the node UUID
func nodeIdByName(node string) (string, error) {
	req := proto.NewNodeFilter()
	req.Filter.Node.Name = node

	res, err := fetchFilter(req, `/filter/nodes/`)
	if err != nil {
		goto abort
	}

	if res.Nodes == nil || len(*res.Nodes) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if node != (*res.Nodes)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			node, (*res.Nodes)[0].Name)
		goto abort
	}
	return (*res.Nodes)[0].Id, nil

abort:
	return ``, fmt.Errorf("NodeId lookup failed: %s",
		err.Error())
}

// nodeConfigById implements the actual lookup of the node's repo
// and bucket assignment information from the server
func nodeConfigById(node string) (*proto.NodeConfig, error) {
	path := fmt.Sprintf("/nodes/%s/config", node)
	var (
		err  error
		resp *resty.Response
		res  *proto.Result
	)
	if resp, err = GetReq(path); err != nil {
		goto abort
	}
	if res, err = decodeResponse(resp); err != nil {
		goto abort
	}
	if res.StatusCode == 404 {
		err = fmt.Errorf(`Node is not assigned to a configuration` +
			` repository yet.`)
		goto abort
	}
	if err = checkApplicationError(res); err != nil {
		goto abort
	}

	if res.Nodes == nil || len(*res.Nodes) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	// check the received record against the input
	if node != (*res.Nodes)[0].Id {
		err = fmt.Errorf("NodeId mismatch: %s vs %s",
			node, (*res.Nodes)[0].Id)
		goto abort
	}
	return (*res.Nodes)[0].Config, nil

abort:
	return nil, fmt.Errorf("NodeConfig lookup failed: %s",
		err.Error())
}

// capabilityIdByName implements the actual lookup of the capability
// UUID from the server
func capabilityIdByName(cap string) (string, error) {
	var err error
	var res *proto.Result
	req := proto.NewCapabilityFilter()

	split := strings.SplitN(cap, ".", 3)
	if len(split) != 3 {
		err = fmt.Errorf(`Capability split failed, name invalid`)
		goto abort
	}
	if req.Filter.Capability.MonitoringId, err = LookupMonitoringId(
		split[0]); err != nil {
		goto abort
	}
	req.Filter.Capability.View = split[1]
	req.Filter.Capability.Metric = split[2]

	if res, err = fetchFilter(req, `/filter/capability/`); err != nil {
		goto abort
	}

	if res.Capabilities == nil || len(*res.Capabilities) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if cap != (*res.Capabilities)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			cap, (*res.Capabilities)[0].Name)
		goto abort
	}
	return (*res.Capabilities)[0].Id, nil

abort:
	return ``, fmt.Errorf("CapabilityId lookup failed: %s",
		err.Error())
}

// checkConfigIdByName implements the actual lookup of the check
// configuration's UUID from the server
func checkConfigIdByName(check, repo string) (string, error) {
	req := proto.NewCheckConfigFilter()
	req.Filter.CheckConfig.Name = check

	res, err := fetchFilter(req, fmt.Sprintf(
		"/filter/checks/%s/", repo))
	if err != nil {
		goto abort
	}

	if res.CheckConfigs == nil || len(*res.CheckConfigs) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if check != (*res.CheckConfigs)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			check, (*res.CheckConfigs)[0].Name)
		goto abort
	}
	return (*res.CheckConfigs)[0].Id, nil

abort:
	return ``, fmt.Errorf("CheckConfigId lookup failed: %s",
		err.Error())
}

// propertyIdByName implements the actual lookup of property ids
// from the server
func propertyIdByName(pType, pName, refId string) (string, error) {
	req := proto.NewPropertyFilter()
	req.Filter.Property.Type = pType
	req.Filter.Property.Name = pName

	var (
		path string
		err  error
		res  *proto.Result
	)

	switch pType {
	case `custom`:
		// custom properties are per-repository
		req.Filter.Property.RepositoryId = refId
		path = fmt.Sprintf("/filter/property/custom/%s/", refId)
	case `service`:
		path = fmt.Sprintf("/filter/property/service/team/%s/",
			refId)
	case `template`:
		path = `/filter/property/service/global/`
	default:
		err = fmt.Errorf("Unknown property type: %s", pType)
		goto abort
	}

	if res, err = fetchFilter(req, path); err != nil {
		goto abort
	}

	if res.Properties == nil || len(*res.Properties) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	switch pType {
	case `custom`:
		if pName != (*res.Properties)[0].Custom.Name {
			err = fmt.Errorf("Name mismatch: %s vs %s",
				pName, (*res.Properties)[0].Custom.Name)
			goto abort
		}
		if refId != (*res.Properties)[0].Custom.RepositoryId {
			err = fmt.Errorf("RepositoryId mismatch: %s vs %s",
				refId, (*res.Properties)[0].Custom.RepositoryId)
			goto abort
		}
		return (*res.Properties)[0].Custom.Id, nil
	case `service`:
		if refId != (*res.Properties)[0].Service.TeamId {
			err = fmt.Errorf("TeamId mismatch: %s vs %s",
				refId, (*res.Properties)[0].Service.TeamId)
			goto abort
		}
		fallthrough
	case `template`:
		if pName != (*res.Properties)[0].Service.Name {
			err = fmt.Errorf("Name mismatch: %s vs %s",
				pName, (*res.Properties)[0].Service.Name)
			goto abort
		}
		return (*res.Properties)[0].Service.Name, nil
	default:
		err = fmt.Errorf("Unknown property type: %s", pType)
	}

abort:
	return ``, fmt.Errorf("PropertyId lookup failed: %s", err.Error())
}

// levelByName implements the actual lookup of the level details
// from the server
func levelByName(lvl string) (string, error) {
	req := proto.NewLevelFilter()
	req.Filter.Level.Name = lvl
	req.Filter.Level.ShortName = lvl

	res, err := fetchFilter(req, `/filter/levels/`)
	if err != nil {
		goto abort
	}

	if res.Levels == nil || len(*res.Levels) == 0 {
		err = fmt.Errorf(`no object returned`)
		goto abort
	}

	if lvl != (*res.Levels)[0].Name &&
		lvl != (*res.Levels)[0].ShortName {
		err = fmt.Errorf("Name mismatch: %s vs %s/%s",
			lvl, (*res.Levels)[0].Name, (*res.Levels)[0].ShortName)
		goto abort
	}
	return (*res.Levels)[0].Name, nil

abort:
	return ``, fmt.Errorf("LevelName lookup failed: %s",
		err.Error())
}

// fetchFilter is a helper used in the ...IdByFoo functions
func fetchFilter(req proto.Request, path string) (*proto.Result, error) {
	var (
		err  error
		resp *resty.Response
		res  *proto.Result
	)
	if resp, err = PostReqBody(req, path); err != nil {
		// transport errors
		return nil, err
	}

	if res, err = decodeResponse(resp); err != nil {
		// http code errors
		return nil, err
	}

	if err = checkApplicationError(res); err != nil {
		return nil, err
	}
	return res, nil
}

// checkApplicationError tests the server result for
// application errors
func checkApplicationError(result *proto.Result) error {
	if result.StatusCode >= 300 {
		var s string
		// application errors
		if result.StatusCode == 404 {
			s = fmt.Sprintf("Object lookup error: %d - %s",
				result.StatusCode, result.StatusText)
		} else {
			s = fmt.Sprintf("Application error: %d - %s",
				result.StatusCode, result.StatusText)
		}
		m := []string{s}

		if result.Errors != nil {
			m = append(m, *result.Errors...)
		}

		return fmt.Errorf(combineStrings(m...))
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
