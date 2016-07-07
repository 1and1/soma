package tree

import (
	"testing"

	"github.com/satori/go.uuid"
)

// Invalid Attach
func TestInvalidRepositoryAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal attach on repository did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

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
		ParentType: `bucket`,
		ParentId:   buckId,
	})
}

func TestInvalidBucketAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal attach on bucket did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

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
		ParentType: `cluster`,
		ParentId:   clrId,
	})
}

func TestInvalidGroupAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal attach on group did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()
	clrId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   clrId,
	})
}

func TestInvalidClusterAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal attach on cluster did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	clrId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   clrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})

}

func TestInvalidNodeAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal attach on node did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
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
		ParentType: `repository`,
		ParentId:   repoId,
	})

}

// Double Attach
func TestDoubleRepositoryAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Double attach on repository did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()

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
	repo.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
}

func TestDoubleBucketAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Double attach on bucket did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

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
	buck.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentId:   repoId,
	})
}

func TestDoubleGroupAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Double attach on group did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()
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
	grp.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})
}

func TestDoubleClusterAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Double attach on cluster did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

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
	clr := NewCluster(ClusterSpec{
		Id:   clrId,
		Name: `testcluster`,
		Team: teamId,
	})
	clr.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})
	clr.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})
}

func TestDoubleNodeAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Double attach on node did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

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
	node.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})
}

// Invalid Destroy
func TestInvalidRepositoryDestroy(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal destroy on repository did not panic`)
		}
	}()

	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Destroy()
}

func TestInvalidBucketDestroy(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal destroy on bucket did not panic`)
		}
	}()

	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Destroy()
}

func TestInvalidGroupDestroy(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal destroy on group did not panic`)
		}
	}()

	teamId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).Destroy()
}

func TestInvalidClusterDestroy(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal destroy on cluster did not panic`)
		}
	}()

	teamId := uuid.NewV4().String()
	clrId := uuid.NewV4().String()

	// create cluster
	NewCluster(ClusterSpec{
		Id:   clrId,
		Name: `testcluster`,
		Team: teamId,
	}).Destroy()
}

func TestInvalidNodeDestroy(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal destroy on node did not panic`)
		}
	}()

	teamId := uuid.NewV4().String()
	nodeId := uuid.NewV4().String()
	servId := uuid.NewV4().String()

	// create new node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: servId,
		Online:   true,
		Deleted:  false,
	}).Destroy()
}

// Invalid Detach
func TestInvalidRepositoryDetach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal detach on repository did not panic`)
		}
	}()

	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `test`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Detach()
}

func TestInvalidBucketDetach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal detach on bucket did not panic`)
		}
	}()

	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Detach()
}

func TestInvalidGroupDetach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal detach on group did not panic`)
		}
	}()

	teamId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).Detach()
}

func TestInvalidClusterDetach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal detach on cluster did not panic`)
		}
	}()

	teamId := uuid.NewV4().String()
	clrId := uuid.NewV4().String()

	// create cluster
	NewCluster(ClusterSpec{
		Id:   clrId,
		Name: `testcluster`,
		Team: teamId,
	}).Detach()
}

func TestInvalidNodeDetach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal detach on node did not panic`)
		}
	}()

	teamId := uuid.NewV4().String()
	nodeId := uuid.NewV4().String()
	servId := uuid.NewV4().String()

	// create new node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: servId,
		Online:   true,
		Deleted:  false,
	}).Detach()
}

// Invalid ReAttach
func TestInvalidGroupReAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal reattach on group did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testgroup`,
		Team: teamId,
	}).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})
}

func TestInvalidClusterReAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal reattach on cluster did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	clrId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create cluster
	clr := NewCluster(ClusterSpec{
		Id:   clrId,
		Name: `testcluster`,
		Team: teamId,
	})
	clr.ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})
}

func TestInvalidNodeReAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal reattach on node did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	nodeId := uuid.NewV4().String()
	servId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_testing`,
		Action: actionC,
	})

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
	node.ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
