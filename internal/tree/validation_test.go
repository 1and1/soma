/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"strings"
	"testing"

	"github.com/satori/go.uuid"
)

func TestInvalidRepoSpec(t *testing.T) {

	repoId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoName := `repoTest`

	spec := RepositorySpec{
		Id:   ``,
		Name: ``,
		Team: ``,
	}

	if specRepoCheck(spec) {
		t.Errorf(`Empty repositoryID did not error`)
	}
	spec.Id = repoId

	if specRepoCheck(spec) {
		t.Errorf(`Empty repository name did not error`)
	}
	spec.Name = repoName

	if specRepoCheck(spec) {
		t.Errorf(`Empty teamId did not error`)
	}
	spec.Team = teamId

	for i := 1; i < 4; i++ {
		spec.Name = strings.Repeat(`x`, i)
		if specRepoCheck(spec) {
			t.Errorf("Short repo name of length %d did not error", i)
		}
	}

	spec.Name = strings.Repeat(`x`, 129)
	if specRepoCheck(spec) {
		t.Error(`Long repo name of length 129 did not error`)
	}

	spec.Name = repoName
	spec.Id = `invalid`
	if specRepoCheck(spec) {
		t.Error(`Invalid repositoryID uuid did not error`)
	}
	spec.Id = repoId

	spec.Team = `invalid`
	if specRepoCheck(spec) {
		t.Error(`Invalid teamID uuid did not error`)
	}
	spec.Team = teamId

	for i := 4; i < 129; i++ {
		spec.Name = strings.Repeat(`x`, i)
		if !specRepoCheck(spec) {
			t.Errorf("Valid reponame length %d `was not accepted", i)
		}
	}
}

func TestInvalidBucketSpec(t *testing.T) {
	buckId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	bucketName := `bucketTest`
	bucketEnv := `envTest`

	spec := BucketSpec{
		Id:          ``,
		Name:        ``,
		Environment: ``,
		Team:        ``,
		Repository:  ``,
	}

	if specBucketCheck(spec) {
		t.Errorf(`Empty bucketID did not error`)
	}
	spec.Id = buckId

	if specBucketCheck(spec) {
		t.Errorf(`Empty bucketName did not error`)
	}
	spec.Name = bucketName

	if specBucketCheck(spec) {
		t.Errorf(`Empty environment did not error`)
	}
	spec.Environment = bucketEnv

	if specBucketCheck(spec) {
		t.Errorf(`Empty teamID did not error`)
	}
	spec.Team = teamId

	if specBucketCheck(spec) {
		t.Errorf(`Empty repositoryID did not error`)
	}
	spec.Repository = repoId

	for i := 1; i < 4; i++ {
		spec.Name = strings.Repeat(`x`, i)
		if specBucketCheck(spec) {
			t.Errorf("Short bucket name of length %d did not error", i)
		}
	}

	spec.Name = strings.Repeat(`x`, 513)
	if specBucketCheck(spec) {
		t.Errorf(`Long bucket name of length 513 did not error`)
	}
	spec.Name = bucketName

	spec.Id = `invalid`
	if specBucketCheck(spec) {
		t.Errorf(`Invalid bucketId uuid did not error`)
	}
	spec.Id = buckId

	spec.Team = `invalid`
	if specBucketCheck(spec) {
		t.Errorf(`Invalid teamId uuid did not error`)
	}
	spec.Team = teamId

	spec.Repository = `invalid`
	if specBucketCheck(spec) {
		t.Errorf(`Invalid repositoryId uuid did not error`)
	}
	spec.Repository = repoId

	for i := 4; i < 513; i++ {
		spec.Name = strings.Repeat(`x`, i)
		if !specBucketCheck(spec) {
			t.Errorf("Valid bucket name of length %d did error", i)
		}
	}

}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
