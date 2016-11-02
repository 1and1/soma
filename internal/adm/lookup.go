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

// oncallIdByName implements the actual serverside lookup of the
// oncall duty UUID
func oncallIdByName(oncall string) (string, error) {
	req := proto.NewOncallFilter()
	req.Filter.Oncall = &proto.OncallFilter{Name: oncall}

	var (
		err  error
		resp *resty.Response
		res  *proto.Result
	)
	if resp, err = PostReqBody(req, `/filter/oncall/`); err != nil {
		// transport errors
		goto abort
	}

	if res, err = decodeResponse(resp); err != nil {
		// http code errors
		goto abort
	}

	if err = checkApplicationError(res); err != nil {
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

	var (
		err  error
		resp *resty.Response
		res  *proto.Result
	)
	if resp, err = PostReqBody(req, `/filter/users/`); err != nil {
		// transport errors
		goto abort
	}

	if res, err = decodeResponse(resp); err != nil {
		// http code errors
		goto abort
	}

	if err = checkApplicationError(res); err != nil {
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

	var (
		err  error
		resp *resty.Response
		res  *proto.Result
	)
	if resp, err = PostReqBody(req, `/filter/teams/`); err != nil {
		// transport errors
		goto abort
	}

	if res, err = decodeResponse(resp); err != nil {
		// http code errors
		goto abort
	}

	if err = checkApplicationError(res); err != nil {
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

	var (
		err  error
		resp *resty.Response
		res  *proto.Result
	)
	if resp, err = PostReqBody(req,
		`/filter/repository/`); err != nil {
		// transport errors
		goto abort
	}

	if res, err = decodeResponse(resp); err != nil {
		// http code errors
		goto abort
	}

	if err = checkApplicationError(res); err != nil {
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

	var (
		err  error
		resp *resty.Response
		res  *proto.Result
	)
	if resp, err = PostReqBody(req, `/filter/buckets/`); err != nil {
		// transport errors
		goto abort
	}

	if res, err = decodeResponse(resp); err != nil {
		// http code errors
		goto abort
	}

	if err = checkApplicationError(res); err != nil {
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
