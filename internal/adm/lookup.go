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

// LookupGrantIdRef lookup up the UUID of a permission grant from
// the server and fills it into the provided id pointer.
// Error is set if no such grant was found or an error occured.
func LookupGrantIdRef(rcptType, rcptId, permId, cat string,
	id *string) error {
	return grantIdFromServer(rcptType, rcptId, permId, cat, id)
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
