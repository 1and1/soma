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

func TestCheckerAddCheck(t *testing.T) {
	sTree, actionC, errC := testSpawnCheckTree()

	chkConfigId := uuid.NewV4()
	capId := uuid.NewV4()

	chk := Check{
		Id:            uuid.Nil,
		SourceId:      uuid.Nil,
		InheritedFrom: uuid.Nil,
		Inheritance:   true,
		ChildrenOnly:  false,
		Interval:      60,
		ConfigId:      chkConfigId,
		CapabilityId:  capId,
		View:          `any`,
		Thresholds: []CheckThreshold{
			{
				Predicate: `>=`,
				Level:     1,
				Value:     100,
			},
			{
				Predicate: `>=`,
				Level:     3,
				Value:     450,
			},
		},
		Constraints: []CheckConstraint{},
	}

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementName: `checkTest`,
	}, true).SetCheck(chk)

	sTree.ComputeCheckInstances()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, `create`},
		[]string{`group`, `create`},
		[]string{`cluster`, `create`},
		[]string{`bucket`, `node_assignment`}, // NewNode
		[]string{`node`, `update`},
		[]string{`bucket`, `node_assignment`}, // NewNode
		[]string{`node`, `update`},
		[]string{`bucket`, `node_assignment`}, // NewNode
		[]string{`node`, `update`},
		[]string{`group`, `member_new`}, // MoveClusterToGroup
		[]string{`cluster`, `update`},
		[]string{`cluster`, `member_new`}, // MoveNodeToCluster
		[]string{`node`, `update`},
		[]string{`group`, `member_new`}, // MoveNodeToGroup
		[]string{`node`, `update`},
		[]string{`node`, `check_new`}, // SetCheck
		[]string{`cluster`, `check_new`},
		[]string{`node`, `check_new`},
		[]string{`group`, `check_new`},
		[]string{`node`, `check_new`},
		[]string{`bucket`, `check_new`},
		[]string{`repository`, `check_new`},
		[]string{`node`, `check_instance_create`}, // ComputeInstances
		[]string{`node`, `check_instance_create`},
		[]string{`node`, `check_instance_create`},
		[]string{`cluster`, `check_instance_create`},
		[]string{`group`, `check_instance_create`},
	}
	for a := range actionC {
		if elem >= len(actions) {
			t.Error(
				`Received additional action`,
				a.Type, a.Action,
			)
			continue
		}

		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}
}

func TestCheckerDeleteCheck(t *testing.T) {
	sTree, actionC, errC := testSpawnCheckTree()

	chkConfigId := uuid.NewV4()
	capId := uuid.NewV4()
	chkId := uuid.NewV4()

	chk := Check{
		Id:            chkId,
		SourceId:      uuid.Nil,
		InheritedFrom: uuid.Nil,
		Inheritance:   true,
		ChildrenOnly:  false,
		Interval:      60,
		ConfigId:      chkConfigId,
		CapabilityId:  capId,
		View:          `any`,
		Thresholds: []CheckThreshold{
			{
				Predicate: `>=`,
				Level:     1,
				Value:     100,
			},
			{
				Predicate: `>=`,
				Level:     3,
				Value:     450,
			},
		},
		Constraints: []CheckConstraint{},
	}

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementName: `checkTest`,
	}, true).SetCheck(chk)

	sTree.ComputeCheckInstances()

	delChk := Check{
		Id:            uuid.Nil,
		InheritedFrom: uuid.Nil,
		SourceId:      chkId,
		ConfigId:      chkConfigId,
	}

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementName: `checkTest`,
	}, true).DeleteCheck(delChk)

	sTree.ComputeCheckInstances()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, `create`},
		[]string{`group`, `create`},
		[]string{`cluster`, `create`},
		[]string{`bucket`, `node_assignment`}, // NewNode
		[]string{`node`, `update`},
		[]string{`bucket`, `node_assignment`}, // NewNode
		[]string{`node`, `update`},
		[]string{`bucket`, `node_assignment`}, // NewNode
		[]string{`node`, `update`},
		[]string{`group`, `member_new`}, // MoveClusterToGroup
		[]string{`cluster`, `update`},
		[]string{`cluster`, `member_new`}, // MoveNodeToCluster
		[]string{`node`, `update`},
		[]string{`group`, `member_new`}, // MoveNodeToGroup
		[]string{`node`, `update`},
		[]string{`node`, `check_new`}, // SetCheck
		[]string{`cluster`, `check_new`},
		[]string{`node`, `check_new`},
		[]string{`group`, `check_new`},
		[]string{`node`, `check_new`},
		[]string{`bucket`, `check_new`},
		[]string{`repository`, `check_new`},
		[]string{`node`, `check_instance_create`}, // ComputeInstances
		[]string{`node`, `check_instance_create`},
		[]string{`node`, `check_instance_create`},
		[]string{`cluster`, `check_instance_create`},
		[]string{`group`, `check_instance_create`},
		[]string{`node`, `check_removed`}, // DeleteCheck
		[]string{`cluster`, `check_removed`},
		[]string{`node`, `check_removed`},
		[]string{`group`, `check_removed`},
		[]string{`node`, `check_removed`},
		[]string{`bucket`, `check_removed`},
		[]string{`repository`, `check_removed`},
		[]string{`node`, `check_instance_delete`},
		[]string{`node`, `check_instance_delete`},
		[]string{`node`, `check_instance_delete`},
		[]string{`cluster`, `check_instance_delete`},
		[]string{`group`, `check_instance_delete`},
	}
	for a := range actionC {
		if elem >= len(actions) {
			t.Error(
				`Received additional action`,
				elem, a.Type, a.Action,
			)
			elem++
			continue
		}

		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action`, elem, `. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}
}

func testSpawnCheckTree() (*Tree, chan *Action, chan *Error) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootId := uuid.NewV4().String()
	teamId := uuid.NewV4().String()
	repoId := uuid.NewV4().String()
	buckId := uuid.NewV4().String()
	grpId := uuid.NewV4().String()
	clrId := uuid.NewV4().String()
	nod1Id := uuid.NewV4().String()
	srv1Id := uuid.NewV4().String()
	nod2Id := uuid.NewV4().String()
	srv2Id := uuid.NewV4().String()
	nod3Id := uuid.NewV4().String()
	srv3Id := uuid.NewV4().String()

	sTree := New(TreeSpec{
		Id:     rootId,
		Name:   `root_checkTest`,
		Action: actionC,
	})

	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `checkTest`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   rootId,
	})
	sTree.SetError(errC)

	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `checkTest_master`,
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

	NewGroup(GroupSpec{
		Id:   grpId,
		Name: `testGroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	NewCluster(ClusterSpec{
		Id:   clrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	NewNode(NodeSpec{
		Id:       nod1Id,
		AssetId:  1,
		Name:     `testnode1`,
		Team:     teamId,
		ServerId: srv1Id,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	NewNode(NodeSpec{
		Id:       nod2Id,
		AssetId:  2,
		Name:     `testnode2`,
		Team:     teamId,
		ServerId: srv2Id,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	NewNode(NodeSpec{
		Id:       nod3Id,
		AssetId:  3,
		Name:     `testnode3`,
		Team:     teamId,
		ServerId: srv3Id,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   clrId,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grpId,
	})

	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nod1Id,
	}, true).(ClusterAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   clrId,
	})

	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nod2Id,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grpId,
	})

	return sTree, actionC, errC
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
