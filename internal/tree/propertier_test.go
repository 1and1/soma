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

func TestSetProperty(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`

	propUUID, _ := uuid.FromString(propId)

	// create tree
	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	// set property on repository
	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalue`,
	})

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 14 {
		t.Error(
			`Expected 14 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`repository`, `property_new`},
		[]string{`bucket`, `create`},
		[]string{`bucket`, `property_new`},
		[]string{`group`, `create`},
		[]string{`group`, `property_new`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `property_new`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	if len(sTree.Child.PropertySystem) != 1 {
		t.Error(
			`Exptected 1 system property, got `,
			len(sTree.Child.PropertySystem),
		)
	}

	if _, ok := sTree.Child.PropertySystem[propId]; !ok {
		t.Error(
			`Could not find property under requested id`,
		)
	}

	if len(sTree.Child.Children[buckId].(*Bucket).PropertySystem) != 1 {
		t.Error(
			`Expected bucket to have 1 system property, found`,
			len(sTree.Child.Children[buckId].(*Bucket).PropertySystem),
		)
	}

	for _, p := range sTree.Child.Children[buckId].(*Bucket).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalue` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 1 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 1 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 1 {
		t.Error(
			`Node has wrong system property count`,
		)
	}
}

func TestUpdateProperty(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`

	propUUID, _ := uuid.FromString(propId)

	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalue`,
	})

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).UpdateProperty(&PropertySystem{
		SourceId:     propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalueUPDATED`,
	})
	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 19 {
		t.Error(
			`Expected 19 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`repository`, `property_new`},
		[]string{`bucket`, `create`},
		[]string{`bucket`, `property_new`},
		[]string{`group`, `create`},
		[]string{`group`, `property_new`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `property_new`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
		[]string{`repository`, `property_update`},
		[]string{`bucket`, `property_update`},
		[]string{`group`, `property_update`},
		[]string{`cluster`, `property_update`},
		[]string{`node`, `property_update`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	if len(sTree.Child.PropertySystem) != 1 {
		t.Error(
			`Exptected repository to have 1 system property, found`,
			len(sTree.Child.PropertySystem),
		)
	}

	if _, ok := sTree.Child.PropertySystem[propId]; !ok {
		t.Error(
			`Could not find property under requested id`,
		)
	}

	if len(sTree.Child.Children[buckId].(*Bucket).PropertySystem) != 1 {
		t.Error(
			`Expected bucket to have 1 system property, found`,
			len(sTree.Child.Children[buckId].(*Bucket).PropertySystem),
		)
	}

	for _, p := range sTree.Child.Children[buckId].(*Bucket).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueUPDATED` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 1 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 1 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 1 {
		t.Error(
			`Node has wrong system property count`,
		)
	}
}

func TestDeleteProperty(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`

	propUUID, _ := uuid.FromString(propId)

	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalue`,
	})

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).UpdateProperty(&PropertySystem{
		SourceId:     propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalueUPDATED`,
	})

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).DeleteProperty(&PropertySystem{
		SourceId: propUUID,
		View:     `testview`,
		Key:      `testkey`,
		Value:    `testvalue`,
	})

	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 24 {
		t.Error(
			`Expected 24 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`repository`, `property_new`},
		[]string{`bucket`, `create`},
		[]string{`bucket`, `property_new`},
		[]string{`group`, `create`},
		[]string{`group`, `property_new`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `property_new`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
		[]string{`repository`, `property_update`},
		[]string{`bucket`, `property_update`},
		[]string{`group`, `property_update`},
		[]string{`cluster`, `property_update`},
		[]string{`node`, `property_update`},
		[]string{`repository`, `property_delete`},
		[]string{`bucket`, `property_delete`},
		[]string{`group`, `property_delete`},
		[]string{`cluster`, `property_delete`},
		[]string{`node`, `property_delete`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	if len(sTree.Child.PropertySystem) != 0 {
		t.Error(
			`Exptected repository to have 0 system property, found`,
			len(sTree.Child.PropertySystem),
		)
	}

	if len(sTree.Child.Children[buckId].(*Bucket).PropertySystem) != 0 {
		t.Error(
			`Expected bucket to have 0 system property, found`,
			len(sTree.Child.Children[buckId].(*Bucket).PropertySystem),
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 0 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 0 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 0 {
		t.Error(
			`Node has wrong system property count`,
		)
	}
}

func TestDeletePropertyNoInheritance(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`

	propUUID, _ := uuid.FromString(propId)

	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  false,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalue`,
	})

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).DeleteProperty(&PropertySystem{
		SourceId: propUUID,
		View:     `testview`,
		Key:      `testkey`,
		Value:    `testvalue`,
	})

	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 11 {
		t.Error(
			`Expected 11 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`repository`, `property_new`},
		[]string{`bucket`, `create`},
		[]string{`group`, `create`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`repository`, `property_delete`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	if len(sTree.Child.PropertySystem) != 0 {
		t.Error(
			`Exptected repository to have 0 system property, found`,
			len(sTree.Child.PropertySystem),
		)
	}

	if len(sTree.Child.Children[buckId].(*Bucket).PropertySystem) != 0 {
		t.Error(
			`Expected bucket to have 0 system property, found`,
			len(sTree.Child.Children[buckId].(*Bucket).PropertySystem),
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 0 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 0 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 0 {
		t.Error(
			`Node has wrong system property count`,
		)
	}
}

func TestOverwriteProperty(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`
	overId := `90009000-9000-4000-9000-900090009000`

	propUUID, _ := uuid.FromString(propId)
	overUUID, _ := uuid.FromString(overId)

	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalue`,
	})

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).UpdateProperty(&PropertySystem{
		SourceId:     propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalueUPDATED`,
	})

	sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           overUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `OVERWRITE`,
	})

	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 25 {
		t.Error(
			`Expected 25 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`repository`, `property_new`},
		[]string{`bucket`, `create`},
		[]string{`bucket`, `property_new`},
		[]string{`group`, `create`},
		[]string{`group`, `property_new`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `property_new`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
		[]string{`repository`, `property_update`},
		[]string{`bucket`, `property_update`},
		[]string{`group`, `property_update`},
		[]string{`cluster`, `property_update`},
		[]string{`node`, `property_update`},
		[]string{`group`, `property_delete`},
		[]string{`cluster`, `property_delete`},
		[]string{`node`, `property_delete`},
		[]string{`node`, `property_new`},
		[]string{`cluster`, `property_new`},
		[]string{`group`, `property_new`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	if len(sTree.Child.PropertySystem) != 1 {
		t.Error(
			`Exptected repository to have 1 system property, found`,
			len(sTree.Child.PropertySystem),
		)
	}

	if _, ok := sTree.Child.PropertySystem[propId]; !ok {
		t.Error(
			`Could not find property under requested id`,
		)
	}

	if len(sTree.Child.Children[buckId].(*Bucket).PropertySystem) != 1 {
		t.Error(
			`Expected bucket to have 1 system property, found`,
			len(sTree.Child.Children[buckId].(*Bucket).PropertySystem),
		)
	}

	for _, p := range sTree.Child.Children[buckId].(*Bucket).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueUPDATED` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 1 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 1 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 1 {
		t.Error(
			`Node has wrong system property count`,
		)
	}
}

func TestUpdateAfterOverwriteProperty(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`
	overId := `90009000-9000-4000-9000-900090009000`

	propUUID, _ := uuid.FromString(propId)
	overUUID, _ := uuid.FromString(overId)

	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalue`,
	})

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	// node must have the property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalue` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).UpdateProperty(&PropertySystem{
		SourceId:     propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalueUPDATED`,
	})

	// node must have the updated property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueUPDATED` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           overUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `OVERWRITE`,
	})

	// node must have the overwrite property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem {
		if p.GetSourceInstance() != overId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `OVERWRITE` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).UpdateProperty(&PropertySystem{
		SourceId:     propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `UPDATEAFTEROVERWRITE`,
	})

	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 27 {
		t.Error(
			`Expected 27 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`repository`, `property_new`},
		[]string{`bucket`, `create`},
		[]string{`bucket`, `property_new`},
		[]string{`group`, `create`},
		[]string{`group`, `property_new`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `property_new`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
		[]string{`repository`, `property_update`},
		[]string{`bucket`, `property_update`},
		[]string{`group`, `property_update`},
		[]string{`cluster`, `property_update`},
		[]string{`node`, `property_update`},
		[]string{`group`, `property_delete`},
		[]string{`cluster`, `property_delete`},
		[]string{`node`, `property_delete`},
		[]string{`node`, `property_new`},
		[]string{`cluster`, `property_new`},
		[]string{`group`, `property_new`},
		[]string{`repository`, `property_update`},
		[]string{`bucket`, `property_update`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	if len(sTree.Child.PropertySystem) != 1 {
		t.Error(
			`Exptected repository to have 1 system property, found`,
			len(sTree.Child.PropertySystem),
		)
	}

	if _, ok := sTree.Child.PropertySystem[propId]; !ok {
		t.Error(
			`Could not find property under requested id`,
		)
	}

	if len(sTree.Child.Children[buckId].(*Bucket).PropertySystem) != 1 {
		t.Error(
			`Expected bucket to have 1 system property, found`,
			len(sTree.Child.Children[buckId].(*Bucket).PropertySystem),
		)
	}

	// bucket must have the update after overwrite property
	for _, p := range sTree.Child.Children[buckId].(*Bucket).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `UPDATEAFTEROVERWRITE` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 1 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 1 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 1 {
		t.Error(
			`Node has wrong system property count`,
		)
	}

	// node must still have the overwrite property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem {
		if p.GetSourceInstance() != overId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `OVERWRITE` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}
}

func TestSetAboveSetProperty(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`
	overId := `90009000-9000-4000-9000-900090009000`

	propUUID, _ := uuid.FromString(propId)
	overUUID, _ := uuid.FromString(overId)

	// create tree
	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	// set property on cluster
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testLOWER`,
	})

	// set property on bucket
	sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           overUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testABOVE`,
	})

	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 13 {
		t.Error(
			`Expected 13 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, `create`},
		[]string{`group`, `create`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
		[]string{`cluster`, `property_new`},
		[]string{`group`, `property_new`},
		[]string{`bucket`, `property_new`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	// repo has no property
	if len(sTree.Child.PropertySystem) != 0 {
		t.Error(
			`Exptected 0 system property, got `,
			len(sTree.Child.PropertySystem),
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(*Bucket).PropertySystem) != 1 {
		t.Error(
			`Bucket has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 1 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 1 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 1 {
		t.Error(
			`Node has wrong system property count`,
		)
	}

	// bucket must have the above property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(*Bucket).PropertySystem {
		if p.GetSourceInstance() != overId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testABOVE` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	// group must have the above property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem {
		if p.GetSourceInstance() != overId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testABOVE` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	// cluster must have the lower property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testLOWER` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	// node must have the lower property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testLOWER` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

}

func TestDeleteAboveSetProperty(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`
	overId := `90009000-9000-4000-9000-900090009000`

	propUUID, _ := uuid.FromString(propId)
	overUUID, _ := uuid.FromString(overId)

	// create tree
	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	// set property on cluster
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testLOWER`,
	})

	// set property on bucket
	sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           overUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testABOVE`,
	})

	// delete property on bucket
	sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(Propertier).DeleteProperty(&PropertySystem{
		SourceId: overUUID,
		View:     `testview`,
		Key:      `testkey`,
		Value:    `testABOVE`,
	})

	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 15 {
		t.Error(
			`Expected 15 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, `create`},
		[]string{`group`, `create`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
		[]string{`cluster`, `property_new`},
		[]string{`group`, `property_new`},
		[]string{`bucket`, `property_new`},
		[]string{`bucket`, `property_delete`},
		[]string{`group`, `property_delete`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	// repo has no property
	if len(sTree.Child.PropertySystem) != 0 {
		t.Error(
			`Exptected 0 system property, got `,
			len(sTree.Child.PropertySystem),
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(*Bucket).PropertySystem) != 0 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 0 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 1 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 1 {
		t.Error(
			`Node has wrong system property count`,
		)
	}

	// cluster must have the lower property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testLOWER` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	// node must have the lower property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testLOWER` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

}

func TestUpdatePropertyInheritanceFalseToTrue(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`

	propUUID, _ := uuid.FromString(propId)

	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	// set property with no inheritance
	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  false,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalue`,
	})

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 0 {
		t.Error(
			`Node has wrong system property count`,
		)
	}

	// update property with inheritance
	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).UpdateProperty(&PropertySystem{
		SourceId:     propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalueUPDATED`,
	})
	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}
	for e := range errChan {
		t.Error(
			`Received error via channel: `, e,
		)
	}

	if len(actionChan) != 15 {
		t.Error(
			`Expected 15 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, `create`},
		[]string{`group`, `create`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`repository`, `property_new`},
		[]string{`repository`, `property_update`},
		[]string{`node`, `property_new`},
		[]string{`cluster`, `property_new`},
		[]string{`group`, `property_new`},
		[]string{`bucket`, `property_new`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	if len(sTree.Child.PropertySystem) != 1 {
		t.Error(
			`Exptected repository to have 1 system property, found`,
			len(sTree.Child.PropertySystem),
		)
	}

	if _, ok := sTree.Child.PropertySystem[propId]; !ok {
		t.Error(
			`Could not find property under requested id`,
		)
	}

	if len(sTree.Child.Children[buckId].(*Bucket).PropertySystem) != 1 {
		t.Error(
			`Expected bucket to have 1 system property, found`,
			len(sTree.Child.Children[buckId].(*Bucket).PropertySystem),
		)
	}

	for _, p := range sTree.Child.Children[buckId].(*Bucket).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueUPDATED` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 1 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 1 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 1 {
		t.Error(
			`Node has wrong system property count`,
		)
	}
}

func TestUpdatePropertyInheritanceTrueToFalse(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`

	propUUID, _ := uuid.FromString(propId)

	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	// set property with inheritance
	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalue`,
	})

	// node must inherited property
	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 1 {
		t.Error(
			`Node has wrong system property count`,
		)
	}
	for _, p := range sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalue` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	// update property with no inheritance
	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).UpdateProperty(&PropertySystem{
		SourceId:     propUUID,
		Inheritance:  false,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalueUPDATED`,
	})
	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}
	for e := range errChan {
		t.Error(
			`Received error via channel: `, e.Action,
		)
	}

	if len(actionChan) != 19 {
		t.Error(
			`Expected 19 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, `create`},
		[]string{`group`, `create`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
		[]string{`cluster`, `property_new`},
		[]string{`group`, `property_new`},
		[]string{`bucket`, `property_new`},
		[]string{`repository`, `property_new`},
		[]string{`repository`, `property_update`},
		[]string{`bucket`, `property_delete`},
		[]string{`group`, `property_delete`},
		[]string{`cluster`, `property_delete`},
		[]string{`node`, `property_delete`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	if len(sTree.Child.PropertySystem) != 1 {
		t.Error(
			`Exptected repository to have 1 system property, found`,
			len(sTree.Child.PropertySystem),
		)
	}

	if _, ok := sTree.Child.PropertySystem[propId]; !ok {
		t.Error(
			`Could not find property under requested id`,
		)
	}

	if len(sTree.Child.Children[buckId].(*Bucket).PropertySystem) != 0 {
		t.Error(
			`Expected bucket to have 0 system property, found`,
			len(sTree.Child.Children[buckId].(*Bucket).PropertySystem),
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 0 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 0 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 0 {
		t.Error(
			`Node has wrong system property count`,
		)
	}
}

func TestDeletePropertyAllLocal(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`
	overId := `90009000-9000-4000-9000-900090009000`

	propUUID, _ := uuid.FromString(propId)
	overUUID, _ := uuid.FromString(overId)

	// create tree
	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	// set property on cluster
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkeyLOW`,
		Value:        `testvalueLOW`,
	})

	// set property on bucket
	sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           overUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkeyHIGH`,
		Value:        `testvalueHIGH`,
	})

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 2 {
		t.Error(
			`Node has wrong system property count`,
		)
	}

	// delete locally set properties on cluster
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(Propertier).deletePropertyAllLocal()

	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 17 {
		t.Error(
			`Expected 17 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, `create`},
		[]string{`group`, `create`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
		[]string{`cluster`, `property_new`},
		[]string{`node`, `property_new`},
		[]string{`cluster`, `property_new`},
		[]string{`group`, `property_new`},
		[]string{`bucket`, `property_new`},
		[]string{`cluster`, `property_delete`},
		[]string{`node`, `property_delete`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	// repo has no property
	if len(sTree.Child.PropertySystem) != 0 {
		t.Error(
			`Exptected 0 system property, got `,
			len(sTree.Child.PropertySystem),
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(*Bucket).PropertySystem) != 1 {
		t.Error(
			`Bucket has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 1 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 1 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 1 {
		t.Error(
			`Node has wrong system property count`,
		)
	}

	// bucket must have the high property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(*Bucket).PropertySystem {
		if p.GetSourceInstance() != overId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkeyHIGH` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueHIGH` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	// group must have the high property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem {
		if p.GetSourceInstance() != overId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkeyHIGH` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueHIGH` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	// cluster must have the high property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem {
		if p.GetSourceInstance() != overId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkeyHIGH` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueHIGH` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	// node must have the high property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem {
		if p.GetSourceInstance() != overId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkeyHIGH` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueHIGH` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

}

func TestDeletePropertyAllInherited(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`
	overId := `90009000-9000-4000-9000-900090009000`

	propUUID, _ := uuid.FromString(propId)
	overUUID, _ := uuid.FromString(overId)

	// create tree
	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	// set property on cluster
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkeyLOW`,
		Value:        `testvalueLOW`,
	})

	// set property on bucket
	sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           overUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkeyHIGH`,
		Value:        `testvalueHIGH`,
	})

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 2 {
		t.Error(
			`Node has wrong system property count`,
		)
	}

	// delete inherited properties on cluster
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(Propertier).deletePropertyAllInherited()

	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 17 {
		t.Error(
			`Expected 17 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, `create`},
		[]string{`group`, `create`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
		[]string{`cluster`, `property_new`},
		[]string{`node`, `property_new`},
		[]string{`cluster`, `property_new`},
		[]string{`group`, `property_new`},
		[]string{`bucket`, `property_new`},
		[]string{`cluster`, `property_delete`},
		[]string{`node`, `property_delete`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	// repo has no property
	if len(sTree.Child.PropertySystem) != 0 {
		t.Error(
			`Exptected 0 system property, got `,
			len(sTree.Child.PropertySystem),
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(*Bucket).PropertySystem) != 1 {
		t.Error(
			`Bucket has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 1 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 1 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 1 {
		t.Error(
			`Node has wrong system property count`,
		)
	}

	// bucket must have the high property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(*Bucket).PropertySystem {
		if p.GetSourceInstance() != overId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkeyHIGH` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueHIGH` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	// group must have the high property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem {
		if p.GetSourceInstance() != overId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkeyHIGH` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueHIGH` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	// cluster must have the low property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkeyLOW` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueLOW` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	// node must have the low property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkeyLOW` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueLOW` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

}

func TestDeletePropertyAll(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`
	overId := `90009000-9000-4000-9000-900090009000`

	propUUID, _ := uuid.FromString(propId)
	overUUID, _ := uuid.FromString(overId)

	// create tree
	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	// set property on cluster
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkeyLOW`,
		Value:        `testvalueLOW`,
	})

	// set property on bucket
	sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           overUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkeyHIGH`,
		Value:        `testvalueHIGH`,
	})

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 2 {
		t.Error(
			`Node has wrong system property count`,
		)
	}

	// delete inherited properties on cluster
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(Propertier).deletePropertyAllInherited()
	// delete locally set properties on cluster
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(Propertier).deletePropertyAllLocal()

	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 19 {
		t.Error(
			`Expected 19 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, `create`},
		[]string{`group`, `create`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
		[]string{`cluster`, `property_new`},
		[]string{`node`, `property_new`},
		[]string{`cluster`, `property_new`},
		[]string{`group`, `property_new`},
		[]string{`bucket`, `property_new`},
		[]string{`cluster`, `property_delete`},
		[]string{`node`, `property_delete`},
		[]string{`cluster`, `property_delete`},
		[]string{`node`, `property_delete`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	// repo has no property
	if len(sTree.Child.PropertySystem) != 0 {
		t.Error(
			`Exptected 0 system property, got `,
			len(sTree.Child.PropertySystem),
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(*Bucket).PropertySystem) != 1 {
		t.Error(
			`Bucket has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 1 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 0 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 0 {
		t.Error(
			`Node has wrong system property count`,
		)
	}

	// bucket must have the high property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(*Bucket).PropertySystem {
		if p.GetSourceInstance() != overId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkeyHIGH` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueHIGH` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	// group must have the high property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem {
		if p.GetSourceInstance() != overId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkeyHIGH` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalueHIGH` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}
}

func TestBackflowAfterDeleteSetProperty(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`
	overId := `90009000-9000-4000-9000-900090009000`

	propUUID, _ := uuid.FromString(propId)
	overUUID, _ := uuid.FromString(overId)

	// create tree
	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	// set property on cluster
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testLOWER`,
	})

	// set property on bucket
	sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           overUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testABOVE`,
	})

	// delete property on cluster
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(Propertier).DeleteProperty(&PropertySystem{
		SourceId: propUUID,
		View:     `testview`,
		Key:      `testkey`,
		Value:    `testLOWER`,
	})

	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 17 {
		t.Error(
			`Expected 17 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, `create`},
		[]string{`group`, `create`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
		[]string{`cluster`, `property_new`},
		[]string{`group`, `property_new`},
		[]string{`bucket`, `property_new`},
		[]string{`cluster`, `property_delete`},
		[]string{`node`, `property_delete`},
		[]string{`node`, `property_new`},
		[]string{`cluster`, `property_new`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	// repo has no property
	if len(sTree.Child.PropertySystem) != 0 {
		t.Error(
			`Exptected 0 system property, got `,
			len(sTree.Child.PropertySystem),
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(*Bucket).PropertySystem) != 1 {
		t.Error(
			`Bucket has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 1 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 1 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 1 {
		t.Error(
			`Node has wrong system property count`,
		)
	}

	// cluster must have the above property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem {
		if p.GetSourceInstance() != overId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testABOVE` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	// node must have the above property
	for _, p := range sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem {
		if p.GetSourceInstance() != overId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testABOVE` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

}

func TestCloneAfterProperty(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	grupId := `60006000-6000-4000-6000-600060006000`
	cltrId := `70007000-7000-4000-7000-700070007000`
	nodeId := `80008000-8000-4000-8000-800080008000`

	propUUID, _ := uuid.FromString(propId)

	// create tree
	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})
	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	// set property on repository
	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementId:   repoId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalue`,
	})

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// create group
	NewGroup(GroupSpec{
		Id:   grupId,
		Name: `testgroup`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentId:   buckId,
	})

	// create cluster
	NewCluster(ClusterSpec{
		Id:   cltrId,
		Name: `testcluster`,
		Team: teamId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentId:   grupId,
	})

	// assign node
	NewNode(NodeSpec{
		Id:       nodeId,
		AssetId:  1,
		Name:     `testnode`,
		Team:     teamId,
		ServerId: `00000000-0000-0000-0000-000000000000`,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentId:   cltrId,
	})

	sTree.Begin()
	sTree.Rollback()

	close(actionChan)
	close(errChan)

	if len(errChan) != 0 {
		t.Error(
			`Expected 0 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 15 {
		t.Error(
			`Expected 14 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`repository`, `property_new`},
		[]string{`bucket`, `create`},
		[]string{`bucket`, `property_new`},
		[]string{`group`, `create`},
		[]string{`group`, `property_new`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `property_new`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
		[]string{`errorchannel`, `attached`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	if len(sTree.Child.PropertySystem) != 1 {
		t.Error(
			`Exptected 1 system property, got `,
			len(sTree.Child.PropertySystem),
		)
	}

	if _, ok := sTree.Child.PropertySystem[propId]; !ok {
		t.Error(
			`Could not find property under requested id`,
		)
	}

	if len(sTree.Child.Children[buckId].(*Bucket).PropertySystem) != 1 {
		t.Error(
			`Expected bucket to have 1 system property, found`,
			len(sTree.Child.Children[buckId].(*Bucket).PropertySystem),
		)
	}

	for _, p := range sTree.Child.Children[buckId].(*Bucket).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalue` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}

	if len(sTree.Find(FindRequest{
		ElementType: `group`,
		ElementId:   grupId,
	}, true).(*Group).PropertySystem) != 1 {
		t.Error(
			`Group has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementId:   cltrId,
	}, true).(*Cluster).PropertySystem) != 1 {
		t.Error(
			`Cluster has wrong system property count`,
		)
	}

	if len(sTree.Find(FindRequest{
		ElementType: `node`,
		ElementId:   nodeId,
	}, true).(*Node).PropertySystem) != 1 {
		t.Error(
			`Node has wrong system property count`,
		)
	}
}

func TestSetPropertyDuplicateDetectOnBucket(t *testing.T) {
	actionChan := make(chan *Action, 1024)
	errChan := make(chan *Error, 1024)

	treeId := `10001000-1000-4000-1000-100010001000`
	repoId := `20002000-2000-4000-2000-200020002000`
	propId := `30003000-3000-4000-3000-300030003000`
	teamId := `40004000-4000-4000-4000-400040004000`
	buckId := `50005000-5000-4000-5000-500050005000`
	dupeId := `99999999-9999-4999-9999-999999999999`

	propUUID, _ := uuid.FromString(propId)
	dupeUUID, _ := uuid.FromString(dupeId)

	// create tree
	sTree := New(TreeSpec{
		Id:     treeId,
		Name:   `root_testing`,
		Action: actionChan,
	})

	// create repository
	NewRepository(RepositorySpec{
		Id:      repoId,
		Name:    `testrepo`,
		Team:    teamId,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentId:   treeId,
	})
	sTree.SetError(errChan)

	// create bucket
	NewBucket(BucketSpec{
		Id:          buckId,
		Name:        `testrepo_test`,
		Environment: `testing`,
		Team:        teamId,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoId,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: "repository",
		ParentId:   `repoId`,
		ParentName: `testrepo`,
	})

	// set property on bucket
	sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           propUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalue`,
	})

	// set property on bucket, again
	sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementId:   buckId,
	}, true).(Propertier).SetProperty(&PropertySystem{
		Id:           dupeUUID,
		Inheritance:  true,
		ChildrenOnly: false,
		View:         `testview`,
		Key:          `testkey`,
		Value:        `testvalue`,
	})

	close(actionChan)
	close(errChan)

	if len(errChan) != 1 {
		t.Error(
			`Expected 1 actions in errorChan, got`,
			len(errChan),
		)
	}

	if len(actionChan) != 5 {
		t.Error(
			`Expected 5 actions in actionChan, got`,
			len(actionChan),
		)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, `create`},
		[]string{`bucket`, `property_new`},
		[]string{`group`, `create`},
		[]string{`group`, `property_new`},
		[]string{`group`, `member_new`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `property_new`},
		[]string{`cluster`, `member_new`},
		[]string{`node`, `update`},
		[]string{`node`, `property_new`},
	}
	for a := range actionChan {
		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}

	if len(sTree.Child.PropertySystem) != 0 {
		t.Error(
			`Exptected 0 system property, got `,
			len(sTree.Child.PropertySystem),
		)
	}

	if len(sTree.Child.Children[buckId].(*Bucket).PropertySystem) != 1 {
		t.Error(
			`Expected bucket to have 1 system property, found`,
			len(sTree.Child.Children[buckId].(*Bucket).PropertySystem),
		)
	}

	for _, p := range sTree.Child.Children[buckId].(*Bucket).PropertySystem {
		if p.GetSourceInstance() != propId {
			t.Error(`Wrong source id`, p.GetSourceInstance())
		}
		if p.GetKey() != `testkey` {
			t.Error(`Wrong key:`, p.GetKey())
		}

		if p.GetValue() != `testvalue` {
			t.Error(`Wrong value:`, p.GetValue())
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
