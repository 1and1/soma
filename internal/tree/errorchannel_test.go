/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"testing"

	"github.com/satori/go.uuid"
)

func TestErrorChannelNode(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	nodeId := uuid.NewV4().String()
	servId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	repo := NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	})
	repo.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)
	if repo.Fault.Error == nil {
		t.Errorf(`Repository.Fault.Error is nil`)
	} else {
		repo.Fault.Error <- &Error{Action: `testmessage_repo`}
	}

	// create bucket
	buck := NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	})
	buck.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})
	if buck.Fault.Error == nil {
		t.Errorf(`Bucket.Fault.Error is nil`)
	} else {
		buck.Fault.Error <- &Error{Action: `testmessage_bucket`}
	}

	// create new node
	node := NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: servId,
		Online:   true,
		Deleted:  false,
	})

	node.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	if node.Fault.Error == nil {
		t.Errorf(`Node.Fault.Error is nil`)
	} else {
		node.Fault.Error <- &Error{Action: `testmessage_node`}
	}

	close(actionC)
	close(errC)

	if len(errC) != 3 {
		t.Error(len(errC), `elements in error channel`)
	}

	if len(actionC) != 6 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
