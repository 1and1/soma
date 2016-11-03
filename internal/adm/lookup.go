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

// oncallIdByName implements the actual serverside lookup of the
// oncall duty UUID
func oncallIdByName(oncall string) (string, error) {
	req := proto.NewOncallFilter()
	req.Filter.Oncall = &proto.OncallFilter{Name: oncall}

	res, err := fetchFilter(req, `/filter/oncall/`)
	if err != nil {
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
		// application errors
		s := fmt.Sprintf("Request failed: %d - %s",
			result.StatusCode, result.StatusText)
		m := []string{s}

		if result.Errors != nil {
			m = append(m, *result.Errors...)
		}

		return fmt.Errorf(combineStrings(m...))
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
