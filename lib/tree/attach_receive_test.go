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

func TestAttachRepository(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    uuid.NewV4().String(),
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 3 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

func TestAttachBucket(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 4 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

func TestAttachGroup(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 5 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

func TestAttachGroupToGroup(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()
	grpId2 := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grpId2,
		Name: `testgroup2`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grpId,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 7 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

func TestAttachCluster(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	clrId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   clrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 5 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

func TestAttachNode(t *testing.T) {
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
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create new node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: servId,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 6 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

func TestAttachNodeToGroup(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	nodeId := uuid.NewV4().String()
	servId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create group
	grp := NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	})
	grp.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create new node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: servId,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grpId,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 7 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

func TestMoveNodeToGroup(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()
	nodeId := uuid.NewV4().String()
	servId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: servId,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// move node to group
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grpId,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 9 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestMoveClusterToGroup(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()
	clrId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   clrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// move cluster to group
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   clrId,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grpId,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 8 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestMoveGroupToGroup(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()
	grp2Id := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grp2Id,
		Name: `testgroup2`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// move group to group
	sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grp2Id,
	}, true).(*Group).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grpId,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 8 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestMoveNodeToCluster(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	nodeId := uuid.NewV4().String()
	clrId := uuid.NewV4().String()
	servId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create new node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: servId,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   clrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// move node to cluster
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   clrId,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 9 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestDetachGroupToBucket(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()
	grp2Id := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grp2Id,
		Name: `testgroup2`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// move group to group
	sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grp2Id,
	}, true).(*Group).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grpId,
	})

	// detach group
	sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grp2Id,
	}, true).(*Group).Detach()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 10 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestDetachClusterToBucket(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()
	clrId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   clrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// move cluster to group
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   clrId,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grpId,
	})

	// detach cluster
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   clrId,
	}, true).(Attacher).Detach()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 10 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestDetachNodeToBucket(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()
	clrId := uuid.NewV4().String()
	nodeId := uuid.NewV4().String()
	servId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create new node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: servId,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   clrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// move cluster to group
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   clrId,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grpId,
	})

	// move node to cluster
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(ClusterAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   clrId,
	})

	// detach node
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(Attacher).Detach()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 15 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestDestroyRepository(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()
	clrId := uuid.NewV4().String()
	nodeId := uuid.NewV4().String()
	servId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create new node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: servId,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   clrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// move cluster to group
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   clrId,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grpId,
	})

	// move node to cluster
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(ClusterAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   clrId,
	})

	// destroy bucket
	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Attacher).Destroy()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 20 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}

	if sTree.Child != nil {
		t.Error(`Destroy failed`)
	}
}

func TestDestroyBucket(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()
	clrId := uuid.NewV4().String()
	nodeId := uuid.NewV4().String()
	servId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create new node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: servId,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   clrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// move cluster to group
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   clrId,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grpId,
	})

	// move node to cluster
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(ClusterAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   clrId,
	})

	// destroy bucket
	sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(Attacher).Destroy()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 18 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestRollbackDetachNodeToBucket(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()
	clrId := uuid.NewV4().String()
	nodeId := uuid.NewV4().String()
	servId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create new node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: servId,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   clrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// move cluster to group
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   clrId,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grpId,
	})

	// move node to cluster
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(ClusterAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   clrId,
	})

	sTree.Begin()

	if sTree.Snap.Id.String() != repoId {
		t.Error(`Clone failure`)
	}
	if sTree.Snap.Children[buckId].(*Bucket).
		Children[grpId].(*Group).
		Children[clrId].(*Cluster).
		Children[nodeId].(*Node).Name != `testnode` {
		t.Error(`Deep clone failure`)
	}

	// detach node
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(Attacher).Detach()

	sTree.Rollback()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 16 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}

	if sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   clrId,
	}, true).(*Cluster).Children[nodeId] != sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true) {
		t.Error(`Bad things`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
